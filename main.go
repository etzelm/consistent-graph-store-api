package main

import (
	"os"

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

// Represents the string representation of the 'status' field in some
// responses
var statuses = [...]string{
	"success",
	"error",
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
