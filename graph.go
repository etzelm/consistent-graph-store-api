package main

import "fmt"

// http://karras.rutgers.edu/hexastore.pdf
// This is a cool idea ^^^^^ but probably originally we should just make a naive implementation
// then we can always optimize. We'll see.

type Node interface {
	ID() int
	String() string

	// ? for holding the actual data potentially. Chances are holding it in some sort of fast binary
	// format is probably more efficient that using interface{} and reflection.
	// Value() map[string][]byte

	// alternatively:
	Value() map[string]interface{}
}

type Edge interface {
	Source() Node
	Target() Node
	Weight() float64
	String() string

	// ? Should roughly represent a NoSQL datastore / json datastore in structure for the main data
	// because that is a really easy to use format.
	Value() map[string]interface{}
}

type Graph interface {
	// for a little setup maybe
	// Init()
	AddNode(n Node)
	AddEdge(e Edge)
}

type g struct {
}

func (test *g) AddNode(n Node) {

}

func (test *g) AddEdge(e Edge) {

}

func testInterface(g Graph) {
	fmt.Println("interface works")
}
