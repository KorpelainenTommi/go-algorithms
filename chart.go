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
		line.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{
			Width:  "60vw",
			Height: "30vw",
		}))
		line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
			Title: "Local minimum",
		}))

		axis := make([]string, len(f))
		data := make([]opts.LineData, len(f))

		for i := range axis {
			axis[i] = fmt.Sprint(i)
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

func PlotGraph(g *Graph[string]) {

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		graph := charts.NewGraph()
		graph.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{
			Width:  "80vw",
			Height: "80vw",
		}))
		graph.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
			Title: "K-Coloring",
		}))

		nodes := make([]opts.GraphNode, len(g.nodes))
		links := make([]opts.GraphLink, len(g.edges))
		numbering := make(map[*Node[string]]int)

		i := 0
		for node := range g.nodes {
			nodes[i] = opts.GraphNode{
				Name:      fmt.Sprintf("%d%v", i+1, node),
				ItemStyle: &opts.ItemStyle{Color: node.value},
			}
			numbering[node] = i + 1
			i++
		}

		i = 0
		for edge := range g.edges {
			links[i] = opts.GraphLink{
				Source: fmt.Sprintf("%d%v", numbering[edge.a], edge.a),
				Target: fmt.Sprintf("%d%v", numbering[edge.b], edge.b),
				LineStyle: &opts.LineStyle{
					Color: "black",
					Width: 2,
				},
			}
			i++
		}

		graph.AddSeries("Colored", nodes, links, charts.WithGraphChartOpts(opts.GraphChart{
			Force:     &opts.GraphForce{Repulsion: 800},
			Layout:    "force",
			Roam:      opts.Bool(true),
			Draggable: opts.Bool(true),
		}))
		graph.Render(w)
	})

	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
