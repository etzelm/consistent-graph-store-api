package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// StatusMsg is a type alias to allow for 'enum' style type
type StatusMsg int

// The constants to represent a Status message
const (
	SUCCESS StatusMsg = iota
	ERROR
)

// IpPort is this computers IP address
var IpPort string

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

// Represents the string representation of the 'status' field in some
// responses
var statuses = [...]string{
	"success",
	"error",
}

type GetPResponse struct {
	Msg     string `json:"msg"`
	Part_id int    `json:"partition_id"`
}

type ServerNode struct {
	IP   string
	Port string
}

func (n *ServerNode) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

// GenerateNode is used to take an ip/port and return a Node instance
func GenerateServerNode(ip, port string) *ServerNode {
	return &ServerNode{IP: ip, Port: port}
}

// AddServerNode is used to add a node to the given view
func AddServerNode(node ServerNode, view View) (View, bool, int) {
	found := false
	part_id := 0
	for ind, part := range view {
		for _, no := range part {
			if reflect.DeepEqual(no, node) {
				found = true
				part_id = ind
			}
		}
	}
	if !found {
		if len(view[partition_it]) < R {
			view[partition_it] = append(view[partition_it], node)
			part_id = partition_it
			partition_it = partition_it + 1
			if partition_it == num_partitions {
				partition_it = 0
			}
			num_nodes = num_nodes + 1
			return view, false, part_id
		}
		for ind := range view {
			if len(view[ind]) < R {
				view[ind] = append(view[ind], node)
				part_id = ind
				partition_it = ind + 1
				if partition_it == num_partitions {
					partition_it = 0
				}
				num_nodes = num_nodes + 1
				return view, false, part_id
			}
		}
		log.Info("All partitions full, adding new one...")
		view = append(view, []ServerNode{node})
		partition_it = num_partitions
		part_id = partition_it
		num_partitions = num_partitions + 1
		num_nodes = num_nodes + 1
		return view, true, part_id
	}
	return view, false, part_id
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
				num_nodes = num_nodes - 1
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
			num_partitions = num_partitions - 1
		}
	}
	holdNodes := make([]ServerNode, 0)
	for _, part := range newView2 {
		for _, no := range part {
			holdNodes = append(holdNodes, no)
		}
	}
	temp := num_nodes / R
	if temp != num_partitions {
		num_partitions = temp
	}
	realView := make([][]ServerNode, temp)
	partition_it = 0
	num_nodes = 0
	for _, node := range holdNodes {
		realView, _, _ = AddServerNode(node, realView)
	}
	log.Info("After Removing ServerNode: ", realView)
	return realView, deleted
}

func main() {

	// testing:
	test := new(g)
	testInterface(test)

	log.Info("Server is starting...")

	IpPort := os.Getenv("ip_port")
	if IpPort == "" {
		IpPort = "127.0.0.1:8080"
	}
	log.Info("IP_PORT: ", IpPort)

	partition_id = 7

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	LoadRoutes(server)
	//generateTicker()
	server.Run(IpPort)
}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/", LandingPage)

	server.GET("/hello", Hello)

	// All '/check' routes are grouped for convenience/clarity.
	check := server.Group("/check")
	{
		check.GET("", CheckGet)
		check.POST("", CheckPost)
		check.PUT("", CheckPut)
	}

	// All '/gs' routes are grouped for convenience/clarity.
	gs := server.Group("/gs")
	{
		gs.GET("/partition", GetPartition)
	}

	return server
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
