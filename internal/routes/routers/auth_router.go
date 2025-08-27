package routers

import (
	"interview-tracker/internal/adapters/handlers"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.RouterGroup) {
	r.GET("/v1/auth", handlers.GetHealth)
}
