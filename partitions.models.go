package main

// IP is this computers IP address
var IP string

// Port is this computers PORT address
var Port string

// ServerNode is used generically to hold server ip/port combos
type ServerNode struct {
	IP   string
	Port string
}

// VIEW is a list of current nodes in the view.
var VIEW [][]ServerNode

// View is used as a type alias for a slice of Nodes
type View [][]ServerNode

// SELF represents this current node
var SELF ServerNode

// R holds the size that our partitions are supposed to be
var R int

// partitionIter holds the place we want to put the next node we add
var partitionIter int

// partitionID holds the id of the parition the Node belongs too
var partitionID int

// numPartitions holds the number of partitions we currently have
var numPartitions int

// numNodes holds the number of partitions we currently have
var numNodes int

// causalMap holds the latest causal_payload we have currently seen on this server
var causalMap map[string]int64

type server struct{}

type GetPResponse struct {
	Msg     string `json:"msg"`
	Part_id int    `json:"partitionID"`
}

type GetPsResponse struct {
	Msg          string `json:"msg"`
	Part_id_list []int  `json:"partitionID_list"`
}

type GetPartResponse struct {
	Msg          string   `json:"msg"`
	Part_members []string `json:"partition_members"`
}
