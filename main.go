package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// IpPort is this computers IP address
var IpPort string

// VIEW is a list of current nodes in the view.
var VIEW [][]Node

// View is used as a type alias for a slice of Nodes
type View [][]Node

// Self represents this current node
var SELF Node

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

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	LoadRoutes(server)
	generateTicker()
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
