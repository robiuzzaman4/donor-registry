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
	"github.com/robiuzzaman4/donor-registry/internal/config"
	"github.com/robiuzzaman4/donor-registry/internal/repository"
	userhandler "github.com/robiuzzaman4/donor-registry/internal/rest/handler/user"
	"github.com/robiuzzaman4/donor-registry/internal/rest/middleware"
	"github.com/robiuzzaman4/donor-registry/internal/user"
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

	// cors
	router.Use(middleware.CORS())

	// routes
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

		// protected routes
		userhandler.RegisterProtectedRoutes(api, userHandler)
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
