package user

import (
	"github.com/gin-gonic/gin"
	"github.com/robiuzzaman4/donor-registry/internal/rest/middleware"
)

func RegisterPublicRoutes(r *gin.RouterGroup, h *Handler) {

	public := r.Group("/users")
	{
		public.POST("/register", h.Register)
		public.POST("/login", h.Login)
		public.POST("/logout", h.Logout)
		public.POST("/refresh", h.Refresh)
		public.GET("/", h.List)

		public.GET("/:id", h.GetByID)
	}
}

func RegisterProtectedRoutes(r *gin.RouterGroup, h *Handler) {
	protected := r.Group("/users")
	protected.Use(middleware.AuthGuard())
	{
		protected.GET("/me", h.Me)
	}
}
