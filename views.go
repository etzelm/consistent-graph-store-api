package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Hello is the controller for the "/hello" route. Has to match query 'name' or
// address the user as 'user'.
func Hello(c *gin.Context) {
	name, hasName := c.GetQuery("name")
	log.WithFields(log.Fields{"name": name, "hasName": hasName}).Info("Hello request query string -->")
	if hasName == false {
		c.String(http.StatusOK, "Hello user!")
		return
	}
	c.String(http.StatusOK, "Hello %s!", name)
}

// CheckGet is used for all GET requests to the '/check' path.
func CheckGet(c *gin.Context) {
	c.String(http.StatusOK, "This is a GET request")
}

// CheckPost is used for all POST requests to the '/check' path.
func CheckPost(c *gin.Context) {
	c.String(http.StatusOK, "This is a POST request")
}

// CheckPut is used for all PUT requests to the '/check' path. Needed to explicitly
// issue a 405 (method not allowed) instead of the default 404 (not found).
func CheckPut(c *gin.Context) {
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}

func LandingPage(c *gin.Context) {
	c.String(http.StatusOK, "Hello!")
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
