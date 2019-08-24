package main

import (
	"net"
	"os"

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

// IPPort is this computers IP address
var IPPort string

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

	log.Info("Launching gRPC Server...")
	go launchGrpcServer()
	log.Info("Finished Launching gRPC Server...")

	IPPort := os.Getenv("ip_port")
	if IPPort == "" {
		IPPort = "127.0.0.1:8080"
	}
	log.Info("IP_PORT: ", IPPort)

	partition_id = 7

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	LoadRoutes(server)
	//generateTicker()
	server.Run(IPPort)
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
