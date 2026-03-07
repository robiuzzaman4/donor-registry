package rest

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robiuzzaman4/donor-registry-backend/internal/config"
)

type Server struct {
	cnf *config.Config
}

func NewServer(cnf *config.Config) *Server {
	return &Server{
		cnf: cnf,
	}
}

func (server *Server) Start() {
	// setup gin router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"path":      "/",
			"success":   true,
			"message":   "system operational",
			"timestamp": time.Now(),
		})
	})

	// server address
	addr := ":" + strconv.Itoa(server.cnf.Port)
	fmt.Println("Server is running on port:", addr)

	// listen server
	err := http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Println("Server run error:", err)
		os.Exit(1)
	}
}
