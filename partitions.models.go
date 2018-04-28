package main

type ServerNode struct {
	IP   string
	Port string
}

// VIEW is a list of current nodes in the view.
var VIEW [][]ServerNode

// View is used as a type alias for a slice of Nodes
type View [][]ServerNode

// Self represents this current node
var SELF ServerNode

// K holds the size that our partitions are supposed to be
var R int

// partition_it holds the place we want to put the next node we add
var partition_it int

// partition_id holds the id of the parition the Node belongs too
var partition_id int

// num_partitions holds the number of partitions we currently have
var num_partitions int

// num_nodes holds the number of partitions we currently have
var num_nodes int

// server_causal holds the latest causal_payload we have currently seen on this server
var server_causal map[string]int64

type GetPResponse struct {
	Msg     string `json:"msg"`
	Part_id int    `json:"partition_id"`
}

type GetPsResponse struct {
	Msg          string `json:"msg"`
	Part_id_list []int  `json:"partition_id_list"`
}

type GetPartResponse struct {
	Msg          string   `json:"msg"`
	Part_members []string `json:"partition_members"`
}
