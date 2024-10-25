// Implementation of graph coloring algorithm

// Build a two-colored spanning tree, then conflict resolve until the graph is colored
// or error if all options are exhausted

package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Color a graph G with max c colors
func (g *Graph[T]) Color(c int, colors []T) error {

	if g.root == nil {
		return errors.New("Graph is empty")
	}

	if connected, _ := g.Connected(); !connected {
		return errors.New("Graph is not connected")
	}

	if len(colors) < c {
		return errors.New("Not enough colors to color the graph")
	}

	ordering := make(map[*Node[T]]int) // Store the ordering of nodes in the constructed tree
	conflicts := make(EdgeSet[T])      // Store the conflict edges that are not part of the tree

	// Walk the graph, coloring it with alternating colors
	twoColor := func(node *Node[T], i int, visited NodeSet[T]) {

		// Assign the node an ordering and color it based on the numbers parity
		ordering[node] = i
		if i%2 == 0 {
			node.value = colors[0]
		} else {
			node.value = colors[1]
		}

		// Iterate edges
		for e := range node.neighbours {
			var n *Node[T]
			if node != e.a {
				n = e.a
			} else {
				n = e.b
			}

			// This is the parent node so skip it
			if o, ok := ordering[n]; ok && o == i-1 {
				continue
			}

			// Edge to a visited non-parent node, possible conflict
			if visited[n] {
				conflicts[e] = true
			}
		}
	}

	// Tree-building step
	g.root.Walk(make(NodeSet[T]), 0, twoColor)

	// Type to keep a backlog of color changes
	type ColorStep struct {
		node     *Node[T]
		current  T
		original T
		possible map[T]struct{}
	}

	// Track possible color permutations
	backlog := make([]ColorStep, 0)

	// Assign a non-conflicting color to a node
	chooseColor := func(node *Node[T]) int {

		// All choices
		possible := make(map[T]struct{})
		for _, color := range colors {
			possible[color] = struct{}{}
		}

		// Remove choices conflicting with neighbours
		for e := range node.neighbours {
			var n *Node[T]
			if node != e.a {
				n = e.a
			} else {
				n = e.b
			}
			delete(possible, n.value)
		}

		// Already satisfies coloring
		if _, ok := possible[node.value]; ok {
			return 0
		}

		// No options
		if len(possible) < 1 {
			return 1
		}

		// Choose first possible option
		for p := range possible {
			// Track the choice
			backlog = append(backlog, ColorStep{node, p, node.value, possible})
			node.value = p
			return 0
		}
		return 1
	}


	// Backtrack to a multiple color choice node and remove the used color option
	backtrack := func() bool {
		for i := 0; i < len(backlog); i++ {
			s := &backlog[i]
			if len(s.possible) > 1 {
				delete(s.possible, s.current) // This choice lead to a dead end, remove it
				for p := range s.possible {
					s.current = p // Next option
					s.node.value = p
					break
				}
				// Clear backlog after i
				for j := i+1; j < len(backlog); j++ {
					backlog[j].node.value = backlog[j].original
				}
				// Restart backlog keeping from i
				backlog = backlog[:i+2]

				return false
			}
		}
		return true // No more options to backtrack
	}

	// Conflict resolution step
	incorrect := len(conflicts)
	backlen := len(backlog)

	for incorrect > 0 {
		incorrect = 0

		// Pass over conflicts, assigning colors
		for e := range conflicts {
			incorrect += chooseColor(e.a)
			incorrect += chooseColor(e.b)
		}

		// All resolved
		if incorrect < 1 {
			break
		}

		changes := len(backlog) - backlen
		if changes < 1 {
			// The pass made no color changes, the configuration is stuck so backtrack
			if backtrack() {
				// All options exhausted
				return errors.New("Exhausted all configurations, graph might not be colorable")
			}
		} else {
			backlen = len(backlog)
		}
	}

	return nil
}

// Program to demonstrate the graph coloring algorithm
func RunGraphColor() {
	var N int        // Number of nodes for graph
	var D float64    // Average node degree
	var Out string   // Output file
	var noPrint bool // Use this for large N to prevent filling the terminal

	var seed = time.Now().UnixNano()
	var intSeed int

	// Extract args
	for i, v := range os.Args {
		switch v {
		case "-N", "--N":
			SScanInt(os.Args[i+1], &N, "N")
		case "-D", "--D":
			SScanFloat(os.Args[i+1], &D, "D")
		case "-S", "--S":
			SScanInt(os.Args[i+1], &intSeed, "seed")
		case "-O", "--O":
			fmt.Sscanf(os.Args[i+1], "%s", &Out)
		case "--noprint":
			noPrint = true
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

	err := graph.Color(5, []string{"Red", "Green", "Blue", "Cyan", "Pink"})

	if err != nil {
		fmt.Println("Unable to color the graph:", err)
	}

	isConnected, count := graph.Connected()
	fmt.Println("Connected:", isConnected, "\nReachable:", count)

	isColored, ncolors, conflicts := graph.Colored(5)
	fmt.Println("Colored:", isColored, "\nColors used:", ncolors)

	if !isColored && conflicts != nil {
		fmt.Println("There were conflicts with the coloring:\n", conflicts)
	} else {
		fmt.Println("Successfully colored the graph")
	}

	if noPrint {
		fmt.Println("--noprint specified, skip graph printing")
	} else {
		fmt.Println("Printing the graph", graph)
	}

	if len(Out) < 1 {
		saveGraph := "N"
		fmt.Println("Save graph to file? (Y/N)")
		fmt.Scanf("%s\n", &saveGraph)
		if saveGraph == "Y" {
			fmt.Print("Filename: ")
			fmt.Scanf("%s\n", &Out)
		}
	}

	if len(Out) > 0 {
		if err := graph.SaveJson(Out); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Saved to %s\n", Out)
		}
	}

	// Loading graph from JSON
	// loadedGraph, err := LoadGraphJson[string]("./out.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // fmt.Println(loadedGraph)
	// loadedGraph.Stats()

}

// Excellent rendering tool for visualizing the graph data
// https://mikhad.github.io/graph-builder/#2023
