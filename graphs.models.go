package main

// http://karras.rutgers.edu/hexastore.pdf
// This is a cool idea ^^^^^ but probably originally we should just make a naive implementation
// then we can always optimize. We'll see.

// GraphNode interface for the graphs we store
type GraphNode interface {
	ID() int
	String() string

	// ? for holding the actual data potentially. Chances are holding it in some sort of fast binary
	// format is probably more efficient that using interface{} and reflection.
	// Value() map[string][]byte

	// alternatively:
	Value() map[string]interface{}
}

// Edge inferface for the edges we store
type Edge interface {
	Source() GraphNode
	Target() GraphNode
	Weight() float64
	String() string

	// ? Should roughly represent a NoSQL datastore / json datastore in structure for the main data
	// because that is a really easy to use format.
	Value() map[string]interface{}
}
