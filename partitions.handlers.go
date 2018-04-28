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

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetPartition(c *gin.Context) {
	//server_causal[SELF.String()] = server_causal[SELF.String()] + 2
	c.JSON(200, GetPResponse{statuses[SUCCESS], partition_id})
	return
}

func GetAllPartitions(c *gin.Context) {
	ret := make([]int, 0, 0)
	x := 0
	for x < num_partitions {
		ret = append(ret, x)
		x = x + 1
	}
	server_causal[SELF.String()] = server_causal[SELF.String()] + 2
	c.JSON(200, GetPsResponse{statuses[SUCCESS], ret})
	return
}

func GetPartitionMembers(c *gin.Context) {
	part_id, has_part_id := c.GetQuery("partition_id")
	server_causal[SELF.String()] = server_causal[SELF.String()] + 2
	if has_part_id == false {
		c.AbortWithStatusJSON(404, map[string]string{
			"msg":   statuses[ERROR],
			"error": "Input Not Given",
		})
		return
	}
	i, _ := strconv.Atoi(part_id)
	ret := make([]string, 0, 0)
	for _, no := range VIEW[i] {
		ret = append(ret, no.String())
	}
	c.JSON(200, GetPartResponse{statuses[SUCCESS], ret})
	return
}

func stringifyCausal(m map[string]int64) string {
	b := new(bytes.Buffer)
	ip_ports := make([]string, 0, len(m))
	for ip_port := range m {
		ip_ports = append(ip_ports, ip_port)
	}
	sort.Strings(ip_ports)
	for _, ip_port := range ip_ports {
		fmt.Fprintf(b, "%s.", fmt.Sprintf("%d", m[ip_port]))
	}
	b = bytes.NewBuffer(bytes.Trim(b.Bytes(), "."))
	return b.String()
}

//0 == greater && 1 == lesser && 2 == concurrent
func CompareCausal(c1 map[string]int64, c2 map[string]int64) int {
	less_seen := false
	great_seen := false
	if c1 == nil && c2 != nil {
		return 1
	} else if c1 != nil && c2 == nil {
		return 0
	}
	ip_ports := make([]string, 0, len(c1))
	for ip_port := range c1 {
		ip_ports = append(ip_ports, ip_port)
	}
	for _, ip_port := range ip_ports {
		if c1[ip_port] < c2[ip_port] {
			less_seen = true
		} else if c1[ip_port] > c2[ip_port] {
			great_seen = true
		}
	}
	if less_seen && !great_seen {
		return 1
	} else if !less_seen && great_seen {
		return 0
	}
	return 2
}

func UpdateCausal(c1 map[string]int64, c2 map[string]int64) map[string]int64 {
	ip_ports := make([]string, 0, len(c1))
	for ip_port := range c1 {
		ip_ports = append(ip_ports, ip_port)
	}
	for _, ip_port := range ip_ports {
		if c1[ip_port] < c2[ip_port] {
			c1[ip_port] = c2[ip_port]
		}
	}
	return c1
}

func (n *ServerNode) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

// GenerateNode is used to take an ip/port and return a Node instance
func GenerateServerNode(ip, port string) *ServerNode {
	return &ServerNode{IP: ip, Port: port}
}

func UpdateView(c *gin.Context) {
	operation := c.PostForm("type")
	n := c.PostForm("ip_port")
	nodestr := strings.Split(n, ":")
	node := GenerateServerNode(nodestr[0], nodestr[1])
	switch operation {
	case "add":
		fmt.Println("View change -- Add: ", node)
		/* if err, part_id := AddNodesView(node); err != nil {
			server_causal[SELF.String()] = server_causal[SELF.String()] + 2
			c.JSON(405, map[string]string{
				"msg": err.Error(),
			})
			return
		} else {
			server_causal[SELF.String()] = server_causal[SELF.String()] + 2
			c.JSON(200, AddNodeResponse{statuses[SUCCESS], part_id, num_partitions})
			return
		} */

	case "remove":
		fmt.Println("View change -- Remove: ", node)
		/* if err := RemoveNodesView(node); err != nil {
			c.JSON(405, map[string]string{
				"msg": err.Error(),
			})
			return
		} else {
			c.JSON(200, RemoveNodeResponse{statuses[SUCCESS], num_partitions})
			return
		} */
	}

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
