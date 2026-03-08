package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robiuzzaman4/donor-registry-backend/internal/config"
	"github.com/robiuzzaman4/donor-registry-backend/internal/repository"
	userhandler "github.com/robiuzzaman4/donor-registry-backend/internal/rest/handler/user"
	"github.com/robiuzzaman4/donor-registry-backend/internal/user"
)

type Server struct {
	cnf    *config.Config
	ctx    context.Context
	dbConn *pgxpool.Pool
}

func NewServer(cnf *config.Config, ctx context.Context, dbConn *pgxpool.Pool) *Server {
	return &Server{
		cnf:    cnf,
		ctx:    ctx,
		dbConn: dbConn,
	}
}

func (s *Server) Start() {

	// initialize repositories
	userRepo := repository.NewUserRepo(s.dbConn)

	// initialize services
	userSvc := user.NewService(s.ctx, userRepo)

	// initialize handlers
	userHandler := userhandler.NewHandler(s.cnf, userSvc)

	// setup gin router
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		// root route
		api.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"path":      "/api/v1",
				"success":   true,
				"message":   "Ok",
				"timestamp": time.Now(),
			})
		})

		// public routes
		userhandler.RegisterPublicRoutes(api, userHandler)
	}

	// server address
	addr := ":" + strconv.Itoa(s.cnf.Port)
	fmt.Println("Server is running on port:", addr)

	// listen server
	err := http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Println("Server run error:", err)
		os.Exit(1)
	}
}
