package main

import (
	"fmt"
	"math"
	"math/rand"

	g "github.com/AllenDang/giu"
)

var (
	linedata  []float64
	linedata2 []float64
	lineTicks []g.PlotTicker
	bardata   []float64
	bardata2  []float64
	bardata3  []float64
)

func loop() {
	g.SingleWindow("Plot Demo").Layout(g.Layout{
		g.Plot("Plot").AxisLimits(0, 100, -1.2, 1.2, g.ConditionOnce).XTicks(lineTicks, false).Plots(g.Plots{
			g.PlotLine("Plot Line", linedata),
			g.PlotLine("Plot Line2", linedata2),
		}),
		g.Plot("Plot Bars").AxisLimits(0, 10, -1.2, 1.2, g.ConditionOnce).Plots(g.Plots{
			g.PlotBar("Plot Bar", bardata),
			g.PlotBar("Plot Bar2", bardata2).Shift(0.2),
			g.PlotBarH("Plot Bar H", bardata3),
		}),
	})
}

func main() {
	var delta = 0.1
	for x := 0.0; x < 10; x += delta {
		linedata = append(linedata, math.Sin(x))
		linedata2 = append(linedata2, math.Cos(x))

	}

	for i := 0; i < 100; i += 5 {
		lineTicks = append(lineTicks, g.PlotTicker{Position: float64(i), Label: fmt.Sprintf("P%d", i)})
	}

	delta = 1
	for x := 0.0; x < 10; x += delta {
		bardata = append(bardata, math.Sin(x))
		bardata2 = append(bardata2, math.Sin(x)-0.2)
		bardata3 = append(bardata3, rand.Float64())
	}

	wnd := g.NewMasterWindow("Plot Demo", 1000, 600, g.MasterWindowFlagsNotResizable, nil)
	wnd.Run(loop)
}
