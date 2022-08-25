package giu

import (
	"github.com/AllenDang/cimgui-go"
	imgui "github.com/AllenDang/cimgui-go"
)

// SplitDirection represents a direction (vertical/horizontal) of splitting layout.
type SplitDirection uint8

const (
	// DirectionHorizontal is a horizontal line.
	DirectionHorizontal SplitDirection = 1 << iota
	// DirectionVertical is a vertical line.
	DirectionVertical
)

var _ Disposable = &splitLayoutState{}

type splitLayoutState struct {
	delta float32
}

// Dispose implements disposable interface.
func (s *splitLayoutState) Dispose() {
	// noop
}

// SplitLayoutWidget creates two childs with a line between them.
// This line can be moved by the user to adjust child sizes.
type SplitLayoutWidget struct {
	id                  string
	direction           SplitDirection
	layout1             Widget
	layout2             Widget
	originItemSpacingX  float32
	originItemSpacingY  float32
	originFramePaddingX float32
	originFramePaddingY float32
	sashPos             *float32
	border              bool
}

// SplitLayout creates split layout widget.
func SplitLayout(direction SplitDirection, sashPos *float32, layout1, layout2 Widget) *SplitLayoutWidget {
	return &SplitLayoutWidget{
		direction: direction,
		sashPos:   sashPos,
		layout1:   layout1,
		layout2:   layout2,
		border:    true,
		id:        GenAutoID("SplitLayout"),
	}
}

// Border sets if childs should have borders.
func (s *SplitLayoutWidget) Border(b bool) *SplitLayoutWidget {
	s.border = b
	return s
}

// ID allows to manually set splitter's id.
func (s *SplitLayoutWidget) ID(id string) *SplitLayoutWidget {
	s.id = id
	return s
}

// Build implements widget interface.
func (s *SplitLayoutWidget) Build() {
	splitLayoutState := s.getState()
	s.originItemSpacingX, s.originItemSpacingY = GetItemInnerSpacing()
	s.originFramePaddingX, s.originFramePaddingY = GetFramePadding()

	var layout Layout

	*s.sashPos += splitLayoutState.delta
	if *s.sashPos < 1 {
		*s.sashPos = 1
	}

	switch s.direction {
	case DirectionHorizontal:
		_, availableH := GetAvailableRegion()
		if *s.sashPos >= availableH {
			*s.sashPos = availableH
		}

		layout = Layout{
			Column(
				s.buildChild(Auto, *s.sashPos, s.layout1),
				HSplitter(&(splitLayoutState.delta)).Size(0, s.originItemSpacingY),
				s.buildChild(Auto, Auto, s.layout2),
			),
		}
	case DirectionVertical:
		availableW, _ := GetAvailableRegion()
		if *s.sashPos >= availableW {
			*s.sashPos = availableW
		}
		layout = Layout{
			Row(
				s.buildChild(*s.sashPos, 0, s.layout1),
				VSplitter(&(splitLayoutState.delta)).Size(s.originItemSpacingX, 0),
				s.buildChild(Auto, Auto, s.layout2),
			),
		}
	}

	imgui.PushStyleVar_Vec2(imgui.ImGuiStyleVar_ItemSpacing, imgui.NewImVec2(0, 0))
	layout.Build()
	imgui.PopStyleVar(1)
}

func (s *SplitLayoutWidget) restoreItemSpacing(layout Widget) Layout {
	return Layout{
		Custom(func() {
			imgui.PushStyleVar_Vec2(imgui.ImGuiStyleVar_ItemSpacing, imgui.NewImVec2(s.originItemSpacingX, s.originItemSpacingY))
			imgui.PushStyleVar_Vec2(imgui.ImGuiStyleVar_FramePadding, imgui.NewImVec2(s.originFramePaddingX, s.originFramePaddingY))
			// Restore Child bg color
			bgColor := imgui.GetStyleColorVec4(imgui.ImGuiCol_ChildBg)
			imgui.PushStyleColor_Vec4(imgui.ImGuiCol_ChildBg, bgColor)
		}),
		layout,
		Custom(func() {
			imgui.PopStyleColor(1)
			imgui.PopStyleVar(2)
		}),
	}
}

// Build Child panel. If layout is a SplitLayout, set the frame padding to zero.
func (s *SplitLayoutWidget) buildChild(width, height float32, layout Widget) Widget {
	return Layout{
		Custom(func() {
			_, isSplitLayoutWidget := layout.(*SplitLayoutWidget)
			hasFramePadding := isSplitLayoutWidget || !s.border
			hasBorder := !isSplitLayoutWidget && s.border

			if hasFramePadding {
				imgui.PushStyleVar_Vec2(imgui.ImGuiStyleVar_FramePadding, imgui.NewImVec2(0, 0))
			}

			imgui.PushStyleColor_Vec4(cimgui.ImGuiCol_ChildBg, imgui.NewImVec4(0, 0, 0, 0))

			Child().
				Border(hasBorder).
				Size(width, height).
				Layout(s.restoreItemSpacing(layout)).
				Build()

			imgui.PopStyleColor(1)

			if hasFramePadding {
				imgui.PopStyleVar(1)
			}
		}),
	}
}

func (s *SplitLayoutWidget) getState() (state *splitLayoutState) {
	if state = GetState[splitLayoutState](&Context, s.id); state == nil {
		state = &splitLayoutState{delta: 0.0}
		SetState(&Context, s.id, state)
	}
	return state
}
