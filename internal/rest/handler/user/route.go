package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(r *gin.RouterGroup, h *Handler) {

	public := r.Group("/users")
	{
		public.POST("/register", h.Register)
		public.POST("/login", h.Login)
		public.GET("/", h.List)
	}
}
