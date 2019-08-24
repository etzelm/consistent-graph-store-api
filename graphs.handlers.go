package main

// Graph interface for all graph functions
type Graph interface {
	// for a little setup maybe
	// Init()
	AddGraphNode(n GraphNode)
	AddEdge(e Edge)
}

type graph struct {
}

func (test *graph) AddGraphNode(n GraphNode) {

}

func (test *graph) AddEdge(e Edge) {

}
