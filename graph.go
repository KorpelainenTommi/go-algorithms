package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// A node has a value and neighbours
type Node[T any] struct {
	value      T
	neighbours EdgeSet[T]
}

// Connection between nodes
type Edge[T any] struct {
	a, b *Node[T]
}

// Store nodes and edges in a hash set
type EdgeSet[T any] map[Edge[T]]bool
type NodeSet[T any] map[*Node[T]]bool

// Graph contains nodes and edges
type Graph[T any] struct {
	nodes NodeSet[T]
	edges EdgeSet[T]
	root  *Node[T]
}

// Check for edge (a, b) or (b, a)
func (e EdgeSet[T]) Check(eg Edge[T]) bool {
	if eg.a == eg.b {
		return true
	}
	return e[Edge[T]{eg.a, eg.b}] || e[Edge[T]{eg.b, eg.a}]
}

// Add edge (a, b) if it doesn't exist
func (e EdgeSet[T]) Add(eg Edge[T]) {
	if !e.Check(eg) {
		e[eg] = true
	}
}

// Add a new edge to a graph
func (g *Graph[T]) Connect(eg Edge[T]) {
	g.nodes[eg.a] = true
	g.nodes[eg.b] = true
	eg.a.neighbours.Add(eg)
	eg.b.neighbours.Add(eg)
	g.edges.Add(eg)
}

// Remove an edge from the graph
func (g *Graph[T]) Disconnect(eg Edge[T]) {
	eg2 := Edge[T]{eg.b, eg.a}
	delete(g.edges, eg)
	delete(g.edges, eg2)
	delete(eg.a.neighbours, eg)
	delete(eg.a.neighbours, eg2)
	delete(eg.b.neighbours, eg)
	delete(eg.b.neighbours, eg2)
}

// Initialize an empty graph
func EmptyGraph[T any]() *Graph[T] {
	g := Graph[T]{make(NodeSet[T]), make(EdgeSet[T]), nil}
	return &g
}

// Initialize a graph with N nodes
func NewGraph[T any](N int, value T) *Graph[T] {
	g := Graph[T]{make(NodeSet[T], N), make(EdgeSet[T]), nil}
	g.AddNodes(N, value)
	return &g
}

// Add N new nodes to the graph
func (g *Graph[T]) AddNodes(N int, value T) []*Node[T] {
	newNodes := make([]*Node[T], N)
	for i := range newNodes {
		node := &Node[T]{value, make(EdgeSet[T])}
		newNodes[i] = node
		g.nodes[node] = true
		if g.root == nil {
			g.root = node
		}
	}
	return newNodes
}

// Add N new nodes and connect them to a parent node
func (g *Graph[T]) AddNodesToParent(N int, value T, parent *Node[T]) []*Node[T] {
	newNodes := g.AddNodes(N, value)
	for _, v := range newNodes {
		g.Connect(Edge[T]{parent, v})
	}
	return newNodes
}

// Format a node to string
func (n *Node[T]) String() string {
	return fmt.Sprintf("(%v)", n.value)
}

// Format an edge to string
func (e *Edge[T]) String() string {
	return fmt.Sprintf("%v %v 0", e.a, e.b)
}

// Format a node set to string
func (n NodeSet[T]) String() string {
	strs := make([]string, len(n))
	i := 0
	for k := range n {
		strs[i] = fmt.Sprintf("%d%v", i+1, k)
		i++
	}
	return strings.Join(strs, "\n")
}

// Format an edge set to string
func (e EdgeSet[T]) String() string {

	nodeIds := make(map[*Node[T]]int)
	strs := make([]string, len(e))

	// Assign nodes proper indexes
	i := 1
	for k := range e {
		if _, ok := nodeIds[k.a]; !ok {
			nodeIds[k.a] = i
			i++
		}
		if _, ok := nodeIds[k.b]; !ok {
			nodeIds[k.b] = i
			i++
		}
	}

	i = 0
	for k := range e {
		strs[i] = fmt.Sprintf("%d%v %d%v 0", nodeIds[k.a], k.a, nodeIds[k.b], k.b)
		i++
	}
	return strings.Join(strs, "\n")
}

func (g *Graph[T]) Stats() {
	fmt.Printf("Number of nodes: %d, number of edges: %d, average degree: %.2f\n", len(g.nodes), len(g.edges), 2.0*float64(len(g.edges))/float64(len(g.nodes)))
}

func (g *Graph[T]) String() string {
	return fmt.Sprintf("Number of nodes: %d, number of edges: %d, average degree: %.2f\n%v", len(g.nodes), len(g.edges), 2.0*float64(len(g.edges))/float64(len(g.nodes)), g.edges)
}

func RandomGraph[T any](N int, value T, maxEdges int, seed int64) *Graph[T] {

	edgeCount := 0
	r := rand.New(rand.NewSource(seed))

	graph := NewGraph(N, value)
	group := make([]*Node[T], 0, N)
	connected := make(NodeSet[T])

	connected[graph.root] = true
	group = append(group, graph.root)

	for k := range graph.nodes {
		choice := group[r.Intn(len(group))]
		graph.Connect(Edge[T]{k, choice})
		edgeCount++
		if !connected[k] {
			connected[k] = true
			group = append(group, k)
		}
		if !connected[choice] {
			connected[choice] = true
			group = append(group, choice)
		}
	}

	if edgeCount >= maxEdges {
		fmt.Println("Minimally connected tree")
	} else {
		for iterations := 0; edgeCount < maxEdges && iterations < 100000; {
			a := group[r.Intn(len(group))]
			b := group[r.Intn(len(group))]
			if !graph.edges.Check(Edge[T]{a, b}) {
				graph.Connect(Edge[T]{a, b})
				edgeCount++
			} else {
				iterations++
			}
		}
	}

	return graph
}

type JsonEdge[T any] struct {
	A    int
	B    int
	AVal T
	BVal T
}

func (g *Graph[T]) SaveJson(path string) error {

	nodeIds := make(map[*Node[T]]int)
	jsonEdges := make([]JsonEdge[T], len(g.edges))

	// Assign nodes proper indexes
	i := 1
	for k := range g.edges {
		if _, ok := nodeIds[k.a]; !ok {
			nodeIds[k.a] = i
			i++
		}
		if _, ok := nodeIds[k.b]; !ok {
			nodeIds[k.b] = i
			i++
		}
	}

	i = 0
	for k := range g.edges {
		jsonEdges[i] = JsonEdge[T]{nodeIds[k.a], nodeIds[k.b], k.a.value, k.b.value}
		i++
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to open", path)
		return err
	}

	data, err := json.Marshal(jsonEdges)
	if err != nil {
		fmt.Println("Failed to json encode")
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Failed to write")
		return err
	}

	return nil
}

func LoadGraphJson[T any](path string) (*Graph[T], error) {

	nodes := make(map[int]*Node[T])
	jsonEdges := new([]JsonEdge[T])
	bytes, err := os.ReadFile(path)
	graph := EmptyGraph[T]()

	if err != nil {
		fmt.Printf("Failed to read file %s\n", path)
		return nil, err
	}

	if err := json.Unmarshal(bytes, jsonEdges); err != nil {
		fmt.Println("Failed to decode json")
		return nil, err
	}

	for _, v := range *jsonEdges {
		if _, ok := nodes[v.A]; !ok {
			nodes[v.A] = graph.AddNodes(1, v.AVal)[0]
		}
		if _, ok := nodes[v.B]; !ok {
			nodes[v.B] = graph.AddNodes(1, v.BVal)[0]
		}
		graph.Connect(Edge[T]{nodes[v.A], nodes[v.B]})
	}
	return graph, nil
}
