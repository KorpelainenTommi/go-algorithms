package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	var N int     // Number of nodes for graph
	var D float64 // Average node degree
	var seed = time.Now().UnixNano()
	var intSeed int

	for i, v := range os.Args {
		switch v {
		case "-N", "--N":
			SScanInt(os.Args[i+1], &N, "N")
		case "-D", "--D":
			SScanFloat(os.Args[i+1], &D, "D")
		case "-S", "--S":
			SScanInt(os.Args[i+1], &intSeed, "seed")
		}
	}

	fmt.Println("---- Graph coloring program ----")

	if N == 0 {
		fmt.Print("Input the number of nodes N: ")
		ScanInt(&N, "N")
	}

	if D == 0 {
		fmt.Print("Input average node degree D: ")
		ScanFloat(&D, "D")
	}

	if intSeed > 0 {
		seed = int64(intSeed)
	}

	fmt.Printf("---- Creating graph of size %d with average node degree < %.2f ----\n", N, D)
	fmt.Printf("Seed: %v\n", seed)

	// Add an average of N*D / 2 random edges, but assuring the graph is connected
	maxEdges := int(float64(N) * float64(D) / 2.0)

	if N < 2 {
		fmt.Println("Cannot guarantee connectedness")
		os.Exit(1)
	}

	graph := RandomGraph(N, "C", maxEdges, seed)
	fmt.Println(graph)



	// if err := graph.SaveJson("./out.json"); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	if cwd, err := os.Getwd(); err == nil {
	// 		fmt.Printf("Saved to %s/out.json\n", cwd)
	// 	} else {
	// 		fmt.Println("File saved")
	// 	}
	// }

	loadedGraph, err := LoadGraphJson[string]("./out.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(loadedGraph)
	loadedGraph.Stats()


}

// Excellent rendering tool for visualizing the graph data
// https://mikhad.github.io/graph-builder/#2023
