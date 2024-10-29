package main

import (
	"fmt"
	"math"
	"net/http"

	"image"
	"image/color"
	"image/png"

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

func PlotImage(N int, img []bool, record image.Rectangle) {

	newImg := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{N, N},
	})

	x, y := 0, 0
	for i := range img {
		if img[i] {
			newImg.Set(x, y, color.Black)
		} else {
			newImg.Set(x, y, color.White)
		}
		x++
		if x >= N {
			x = 0
			y++
		}
	}

	highlightColor := color.RGBA{0, 0xff, 0, 0xff}
	AlphaMix := func(a color.RGBA, b color.RGBA, blend float64) (blended color.RGBA) {
		f1, f2 := 1.0-blend, blend
		blended.R = uint8(math.Round(f1*float64(a.R) + f2*float64(b.R)))
		blended.G = uint8(math.Round(f1*float64(a.G) + f2*float64(b.G)))
		blended.B = uint8(math.Round(f1*float64(a.B) + f2*float64(b.B)))
		blended.A = 0xff
		return
	}

	x, y = record.Min.X, record.Min.Y
	for y < record.Max.Y {

		r, g, b, a := newImg.At(x, y).RGBA()
		oldColor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
		newImg.SetRGBA(x, y, AlphaMix(oldColor, highlightColor, 0.5))

		x++
		if x >= record.Max.X {
			x = record.Min.X
			y++
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `
		<!doctype html>
		<html>
			<head>
				<title>Submatrix</title>
			</head>
			<body>
				<div style="width: 100%; height: 100%; display: flex; flex-direction: column; justify-content: center; align-items: center;">
					<h1 style="text-align: center;">Largest contiguous submatrix</h1>
					<img src="/img.png" style="width: 30vw; height: auto; image-rendering: pixelated;" />
				</div>
			</body>
		</html>
		`)
	})

	http.HandleFunc("/img.png", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, newImg)
	})

	fmt.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
