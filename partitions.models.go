package main

import "fmt"

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

type AddServerNodeResponse struct {
	Msg       string `json:"msg"`
	Part_id   int    `json:"partition_id"`
	Num_parts int    `json:"number_of_partitions"`
}

type RemoveServerNodeResponse struct {
	Msg       string `json:"msg"`
	Num_parts int    `json:"number_of_partitions"`
}

// StatusMsg is a type alias to allow for 'enum' style type
type StatusMsg int

// Port used by gservice
const (
	port = ":50051"
)

// The constants to represent a Status message
const (
	SUCCESS StatusMsg = iota
	ERROR
)

// Represents the string representation of the 'status' field in some
// responses
var statuses = [...]string{
	"success",
	"error",
}

// String stringifies a ServerNode
func (n *ServerNode) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}
