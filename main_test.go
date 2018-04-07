package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoadRoutes(t *testing.T) {
	server := gin.New()
	server.GET("/", LandingPage)
	server.GET("/hello", Hello)

	check := server.Group("/check")
	{
		check.GET("", CheckGet)
		check.POST("", CheckPost)
		check.PUT("", CheckPut)
	}

	type args struct {
		server *gin.Engine
	}
	tests := []struct {
		name string
		args args
		want *gin.Engine
	}{
		{
			name: "basic - same handlers",
			args: args{
				server: gin.New(),
			},
			want: server,
		},
	}
	e := gin.New()
	fmt.Println(LoadRoutes(e))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadRoutes(tt.args.server); !reflect.DeepEqual(got.Handlers, tt.want.Handlers) {
				t.Errorf("LoadRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}
