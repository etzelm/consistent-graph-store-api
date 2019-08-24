package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	pb "github.com/etzelm/consistent-graph-store-api/gservice"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GetPartition returns id for which partition this server currently belongs to
func GetPartition(c *gin.Context) {
	causalMap[SELF.String()] = causalMap[SELF.String()] + 2
	c.JSON(200, GetPResponse{statuses[SUCCESS], partitionID})
	return
}

// GetAllPartitions returns list of ids for all valid partitions in the system
func GetAllPartitions(c *gin.Context) {
	ret := make([]int, 0, 0)
	x := 0
	for x < numPartitions {
		ret = append(ret, x)
		x = x + 1
	}
	causalMap[SELF.String()] = causalMap[SELF.String()] + 2
	c.JSON(200, GetPsResponse{statuses[SUCCESS], ret})
	return
}

// GetPartitionMembers returns a list of all servers in the requested partition
func GetPartitionMembers(c *gin.Context) {
	partID, hasPartID := c.GetQuery("partitionID")
	causalMap[SELF.String()] = causalMap[SELF.String()] + 2
	if hasPartID == false {
		c.AbortWithStatusJSON(404, map[string]string{
			"msg":   statuses[ERROR],
			"error": "Input Not Given",
		})
		return
	}
	i, _ := strconv.Atoi(partID)
	ret := make([]string, 0, 0)
	for _, no := range VIEW[i] {
		ret = append(ret, no.String())
	}
	c.JSON(200, GetPartResponse{statuses[SUCCESS], ret})
	return
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

// UpdateView handles the main logic for adding/removing ServerNodes from api calls
func UpdateView(c *gin.Context) {
	operation := c.PostForm("type")
	n := c.PostForm("ipPort")
	nodestr := strings.Split(n, ":")
	node := GenerateServerNode(nodestr[0], nodestr[1])
	switch operation {
	case "add":
		fmt.Println("View change -- Add: ", node)
		if partID, err := AddServerNodeView(node); err != nil {
			causalMap[SELF.String()] = causalMap[SELF.String()] + 2
			c.JSON(405, map[string]string{
				"msg": err.Error(),
			})
			return
		}
		causalMap[SELF.String()] = causalMap[SELF.String()] + 2
		c.JSON(200, AddServerNodeResponse{statuses[SUCCESS], partID, numPartitions})
		return

	case "remove":
		fmt.Println("View change -- Remove: ", node)
		if err := RemoveServerNodeView(node); err != nil {
			c.JSON(405, map[string]string{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(200, RemoveServerNodeResponse{statuses[SUCCESS], numPartitions})
		return

	}

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

// AddServerNodeView handles the outer logic of updating all the other server's VIEWs
func AddServerNodeView(node *ServerNode) (int, error) {
	newView, partitionCreated, partID := AddServerNode(*node, VIEW)
	if *node == SELF {
		partition_id = partID
	}
	VIEW = newView
	if serverCausal != nil {
		serverCausal[node.String()] = 0
	}
	fmt.Println("New view: ", VIEW)
	for _, part := range VIEW {
		for _, no := range part {
			fmt.Println(no)
			if !reflect.DeepEqual(no, SELF) {
				serverCausal[SELF.String()] = serverCausal[SELF.String()] + 2
				serverCausal[no.String()] = serverCausal[no.String()] + 2
				conn, err := OpenNodeConnection(&no)
				if err != nil {
					return err, partID
				}
				defer conn.Close()
				c := pb.NewStoreClient(conn)

				newView := GeneratePBView()

				c.AddServerNode(context.Background(), &pb.ViewChangeRequest{
					RequestID: 1,
					ServerNode: &pb.ServerNode{
						IP:   node.IP,
						Port: node.Port,
					},
					Type:        pb.ViewChangeRequest_ADD_NODE,
					CurrentView: newView,
				})
			}
		}
	}

	return nil, partID
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

// RemoveServerNodeView handles the outer logic of updating all the other server's VIEWs
func RemoveServerNodeView(node *Node) error {
	newView, partitionDeleted := RemoveServerNode(*node, VIEW)
	for _, part := range VIEW {
		for _, n := range part {
			if !reflect.DeepEqual(n, SELF) {
				serverCausal[SELF.String()] = serverCausal[SELF.String()] + 2
				conn, err := OpenNodeConnection(&n)
				if err != nil {
					return err
				}
				defer conn.Close()
				c := pb.NewStoreClient(conn)

				newV := GeneratePBView()

				c.RemoveServerNode(context.Background(), &pb.ViewChangeRequest{
					RequestID: 1,
					ServerNode: &pb.ServerNode{
						IP:   node.IP,
						Port: node.Port,
					},
					Type:        pb.ViewChangeRequest_REMOVE_NODE,
					CurrentView: newV,
				})
			}
		}
	}

	if partitionDeleted {
		log.Info("partition deleted")
	}

	// before view assignment
	deleted := false
	log.Infof("newView[%d] <---> VIEW[%d], *node==SELF => %t", len(VIEW), len(newView), reflect.DeepEqual(*node, SELF))
	if len(newView) < len(VIEW) && reflect.DeepEqual(*node, SELF) {
		deleted = true
		partition_id = -1
	}

	// make sure that partition_id does not equal the same one
	if deleted {
		log.Info("THIS NODE WAS DELETED!!!!!!")
	}

	fmt.Println(newView)
	VIEW = newView

	return nil
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
