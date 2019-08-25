package main

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	pb "github.com/etzelm/consistent-graph-store-api/gservice"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

// UpdateView handles the main logic for adding/removing ServerNodes from api calls
func UpdateView(c *gin.Context) {
	operation := c.PostForm("type")
	n := c.PostForm("ipPort")
	nodestr := strings.Split(n, ":")
	node := GenerateServerNode(nodestr[0], nodestr[1])
	switch operation {
	case "add":
		fmt.Println("View change -- Add: ", node)
		partID, err := AddServerNodeView(node)
		if err != nil {
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

// AddServerNodeView handles the outer logic of updating all the other server's VIEWs
func AddServerNodeView(node *ServerNode) (int, error) {
	newView, partitionCreated, partID := AddServerNode(*node, VIEW)
	if *node == SELF {
		partitionID = partID
	}
	VIEW = newView
	if causalMap != nil {
		causalMap[node.String()] = 0
	}
	fmt.Println("New view: ", VIEW)
	for _, part := range VIEW {
		for _, no := range part {
			fmt.Println(no)
			if !reflect.DeepEqual(no, SELF) {
				causalMap[SELF.String()] = causalMap[SELF.String()] + 2
				causalMap[no.String()] = causalMap[no.String()] + 2
				conn, err := OpenNodeConnection(&no)
				if err != nil {
					return partID, err
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

	fmt.Println("Partition created: ", partitionCreated)

	return partID, nil
}

// RemoveServerNodeView handles the outer logic of updating all the other server's VIEWs
func RemoveServerNodeView(node *ServerNode) error {
	newView, partitionDeleted := RemoveServerNode(*node, VIEW)
	for _, part := range VIEW {
		for _, n := range part {
			if !reflect.DeepEqual(n, SELF) {
				causalMap[SELF.String()] = causalMap[SELF.String()] + 2
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
		partitionID = -1
	}

	// make sure that partitionID does not equal the same one
	if deleted {
		log.Info("THIS NODE WAS DELETED!!!!!!")
	}

	fmt.Println(newView)
	VIEW = newView

	return nil
}
