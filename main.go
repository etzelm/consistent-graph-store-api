package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	pb "github.com/etzelm/consistent-graph-store-api/gservice"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

func main() {

	log.Info("Testing server before start-up...")

	test := new(graph)
	testInterface(test)

	log.Info("Server is starting...")

	parseCommandLineArgs()

	log.Info("Launching gRPC Server...")
	go launchGrpcServer()
	log.Info("Finished Launching gRPC Server...")

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	LoadRoutes(server)
	//generateTicker()
	server.Run(IP + ":" + Port)
}

func parseCommandLineArgs() {
	IP = os.Getenv("IP")
	if IP == "" {
		IP = "127.0.0.1"
	}
	log.Info("IP: ", IP)

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "80"
	}
	log.Info("PORT: ", Port)

	SELF.IP = IP
	SELF.Port = Port

	r := os.Getenv("R")
	R, _ = strconv.Atoi(r)
	if R == 0 {
		R = 2
	}
	log.Info("R: ", R)

	view := os.Getenv("SERVERS")
	log.Info("SERVERS: ", view)
	if view != "" {
		viewStrings := strings.Split(view, ",")
		numPartitions = len(viewStrings) / R
		if numPartitions == 0 {
			numPartitions = 1
		}
		realView := make([][]Node, numPartitions)
		partitionIter = 0
		numNodes = 0
		for _, v := range viewStrings {
			n := strings.Split(v, ":")
			partID := 0
			serverNode := *GenerateServerNode(n[0], n[1])
			realView, _, partID = AddServerNode(serverNode, realView)
			if serverNode == SELF {
				partitionID = partID
			}
		}
		log.Info("Number of Nodes: ", numNodes)
		causalMap = make(map[string]int64, numNodes)
		for _, part := range realView {
			for _, no := range part {
				causalMap[no.String()] = 0
			}
		}
		fmt.Println("CAUSAL: ", causalMap)
		fmt.Println("VIEW: ", realView)
		VIEW = realView
	}
}

func launchGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStoreServer(s, &server{})
	// Register reflection service on gRPC server.

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// LoadRoutes does exactly that... loads all routes for the server.
func LoadRoutes(server *gin.Engine) *gin.Engine {
	server.GET("/gs", LandingPage)

	server.GET("/gs/hello", Hello)

	// All '/check' routes are grouped for convenience/clarity.
	check := server.Group("/gs/check")
	{
		check.GET("", CheckGet)
		check.POST("", CheckPost)
		check.PUT("", CheckPut)
	}

	// All '/gs' routes are grouped for convenience/clarity.
	gs := server.Group("/gs")
	{
		gs.PUT("/change_view", UpdateView)
		gs.GET("/partition", GetPartition)
		gs.GET("/all_partitions", GetAllPartitions)
		gs.GET("/partition_members", GetPartitionMembers)
	}

	return server
}
