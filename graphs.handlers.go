package main

import "fmt"

type Graph interface {
	// for a little setup maybe
	// Init()
	AddGraphNode(n GraphNode)
	AddEdge(e Edge)
}

type g struct {
}

func (test *g) AddGraphNode(n GraphNode) {

}

func (test *g) AddEdge(e Edge) {

}

func testInterface(g Graph) {
	fmt.Println("interface works")
}
