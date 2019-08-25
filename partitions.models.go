package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"

	pb "github.com/etzelm/consistent-graph-store-api/gservice"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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

func (s *server) AddServerNode(c context.Context, vcr *pb.ViewChangeRequest) (*pb.ViewChangeResponse, error) {
	return nil, nil
}

func (s *server) RemoveServerNode(c context.Context, vcr *pb.ViewChangeRequest) (*pb.ViewChangeResponse, error) {
	return nil, nil
}

// GetPResponse is the structure used to return the user requested partition ID
type GetPResponse struct {
	Msg    string `json:"msg"`
	PartID int    `json:"partitionID"`
}

// GetPsResponse is the structure used to return the user requested list of partition IDs
type GetPsResponse struct {
	Msg        string `json:"msg"`
	PartIDList []int  `json:"partitionID_list"`
}

// GetPartResponse is the structure used to return the user requested list of partition members
type GetPartResponse struct {
	Msg         string   `json:"msg"`
	PartMembers []string `json:"partition_members"`
}

// AddServerNodeResponse is the structure used to return the response for adding a server
type AddServerNodeResponse struct {
	Msg      string `json:"msg"`
	PartID   int    `json:"partitionID"`
	NumParts int    `json:"number_of_partitions"`
}

// RemoveServerNodeResponse is the structure used to return the response for removing a server
type RemoveServerNodeResponse struct {
	Msg      string `json:"msg"`
	NumParts int    `json:"number_of_partitions"`
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

// stringifyCausal returns a string version of a CausalMap
func stringifyCausal(m map[string]int64) string {
	b := new(bytes.Buffer)
	ipPorts := make([]string, 0, len(m))
	for ipPort := range m {
		ipPorts = append(ipPorts, ipPort)
	}
	sort.Strings(ipPorts)
	for _, ipPort := range ipPorts {
		fmt.Fprintf(b, "%s.", fmt.Sprintf("%d", m[ipPort]))
	}
	b = bytes.NewBuffer(bytes.Trim(b.Bytes(), "."))
	return b.String()
}

// CompareCausal returns an int value based on Vector Clock comparison
// 0 == greater && 1 == lesser && 2 == concurrent
func CompareCausal(c1 map[string]int64, c2 map[string]int64) int {
	lessSeen := false
	greatSeen := false
	if c1 == nil && c2 != nil {
		return 1
	} else if c1 != nil && c2 == nil {
		return 0
	}
	ipPorts := make([]string, 0, len(c1))
	for ipPort := range c1 {
		ipPorts = append(ipPorts, ipPort)
	}
	for _, ipPort := range ipPorts {
		if c1[ipPort] < c2[ipPort] {
			lessSeen = true
		} else if c1[ipPort] > c2[ipPort] {
			greatSeen = true
		}
	}
	if lessSeen && !greatSeen {
		return 1
	} else if !lessSeen && greatSeen {
		return 0
	}
	return 2
}

// UpdateCausal brings first CausalMap into allignment by taking the later values
// between the two maps
func UpdateCausal(c1 map[string]int64, c2 map[string]int64) map[string]int64 {
	ipPorts := make([]string, 0, len(c1))
	for ipPort := range c1 {
		ipPorts = append(ipPorts, ipPort)
	}
	for _, ipPort := range ipPorts {
		if c1[ipPort] < c2[ipPort] {
			c1[ipPort] = c2[ipPort]
		}
	}
	return c1
}

// GenerateServerNode is used to take an ip/port and return a Node instance
func GenerateServerNode(ip, port string) *ServerNode {
	return &ServerNode{IP: ip, Port: port}
}

// AddServerNode is used to add a node to this server's current given view
func AddServerNode(node ServerNode, view View) (View, bool, int) {
	found := false
	partID := 0
	for ind, part := range view {
		for _, no := range part {
			if reflect.DeepEqual(no, node) {
				found = true
				partID = ind
			}
		}
	}
	if !found {
		if len(view[partitionIter]) < R {
			view[partitionIter] = append(view[partitionIter], node)
			partID = partitionIter
			partitionIter = partitionIter + 1
			if partitionIter == numPartitions {
				partitionIter = 0
			}
			numNodes = numNodes + 1
			return view, false, partID
		}
		for ind := range view {
			if len(view[ind]) < R {
				view[ind] = append(view[ind], node)
				partID = ind
				partitionIter = ind + 1
				if partitionIter == numPartitions {
					partitionIter = 0
				}
				numNodes = numNodes + 1
				return view, false, partID
			}
		}
		log.Info("All partitions full, adding new one...")
		view = append(view, []ServerNode{node})
		partitionIter = numPartitions
		partID = partitionIter
		numPartitions = numPartitions + 1
		numNodes = numNodes + 1
		return view, true, partID
	}
	return view, false, partID
}

// RemoveServerNode removes a node from the given view
func RemoveServerNode(node ServerNode, view View) (View, bool) {
	log.Info("Before Removing ServerNode: ", view)
	newView := make([][]ServerNode, 0)
	for _, part := range view {
		newNodes := make([]ServerNode, 0)
		for _, no := range part {
			if no != node {
				newNodes = append(newNodes, no)
			} else {
				numNodes = numNodes - 1
			}
		}
		newView = append(newView, newNodes)
	}
	deleted := false
	newView2 := make([][]ServerNode, 0)
	for _, part := range newView {
		if len(part) > 0 {
			newView2 = append(newView2, part)
		} else {
			deleted = true
			numPartitions = numPartitions - 1
		}
	}
	holdNodes := make([]ServerNode, 0)
	for _, part := range newView2 {
		for _, no := range part {
			holdNodes = append(holdNodes, no)
		}
	}
	temp := numNodes / R
	if temp != numPartitions {
		numPartitions = temp
	}
	realView := make([][]ServerNode, temp)
	partitionIter = 0
	numNodes = 0
	for _, node := range holdNodes {
		realView, _, _ = AddServerNode(node, realView)
	}
	log.Info("After Removing ServerNode: ", realView)
	return realView, deleted
}

// OpenNodeConnection opens grpc connection with given ServerNode
func OpenNodeConnection(n *ServerNode) (*grpc.ClientConn, error) {
	return grpc.Dial(n.IP+port, grpc.WithInsecure())
}

// GeneratePBView makes the Protocol Buffer version of VIEW
func GeneratePBView() []*pb.View {
	// generate a view
	newView := make([]*pb.View, 0, 0)
	for _, part := range VIEW {
		newNode := make([]*pb.ServerNode, 0, 0)
		for _, no := range part {
			newNode = append(newNode, &pb.ServerNode{
				IP:   no.IP,
				Port: no.Port,
			})
		}
		newView = append(newView, &pb.View{
			CurrentPartition: newNode,
		})
	}
	log.Info(newView)
	return newView
}

// Example ticker
func generateTicker() {
	highest := 200
	for i := 0; i <= 10000; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		antiEntropy := rand.Intn(350) + 200
		if antiEntropy > highest {
			highest = antiEntropy
		}
		log.Info(antiEntropy)
	}
	log.Info("Highest is... ", highest)
	c, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(500 * time.Millisecond)
	go func(ctx context.Context) {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at ", t)
			case <-ctx.Done():
				fmt.Println("exiting goroutine....")
				return
			}
		}
	}(c)
	go func() {
		<-time.After(2 * time.Second)
		cancel()
		fmt.Println("Canceled")
		<-time.After(1 * time.Second)
		ticker.Stop()
		fmt.Println("Ticker stopped.")
	}()
}

func timeSeededRandom() *rand.Rand {
	return rand.New(
		rand.NewSource(time.Now().UnixNano()))
}

// SimpleHash -- This is a copy of how Java generates hashes actually
func SimpleHash(s string) uint {
	h := 0
	for i := 0; i < len(s); i++ {
		h = 31*h + int(s[i])
	}
	return uint(h)
}

// FindPartition is used to figure out which partition a graph belongs to
func FindPartition(s string, partitionnumber uint) uint {
	return SimpleHash(s) % partitionnumber
}

// RouteToNode takes a graph identifier and returns the partition it belongs to
func RouteToNode(graph string) uint {
	return FindPartition(graph, uint(len(VIEW)))
}
