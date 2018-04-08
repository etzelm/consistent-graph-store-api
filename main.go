package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// IP_PORT is this computers IP address
var IP_PORT string

func main() {
	log.Info("Server is starting...")

	IP_PORT := os.Getenv("ip_port")
	if IP_PORT == "" {
		IP_PORT = "127.0.0.1:8080"
	}
	log.Info("IP_PORT: ", IP_PORT)

	server := gin.Default()
	log.WithField("server", server).Info("Default Gin server create.")
	LoadRoutes(server)
	server.Run(IP_PORT)

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