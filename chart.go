package main

import (
	"fmt"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func PlotFunction(f []float64, minimum int) {

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		line := charts.NewLine()
		globalOpts := charts.WithTitleOpts(opts.Title{Title: "Local minimum"})
		line.SetGlobalOptions(globalOpts)

		axis := make([]string, len(f))
		data := make([]opts.LineData, len(f))

		for i := range axis {
			axis[i] = fmt.Sprintf("%v", i)
		}

		for i := range data {
			data[i] = opts.LineData{Value: f[i]}
		}

		coord := make([]interface{}, 2, 2)
		coord[0] = minimum
		coord[1] = f[minimum]

		line.SetXAxis(axis).AddSeries("Range", data,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol: "none",
				Smooth: opts.Bool(true),
			})).SetSeriesOptions(
			charts.WithMarkPointNameCoordItemOpts(opts.MarkPointNameCoordItem{
				Name:       "Local Minimum",
				Coordinate: coord,
				Value:      fmt.Sprintf("%.2f", f[minimum]),
				SymbolSize: 80,
			}))
		line.Render(w)
	})

	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
