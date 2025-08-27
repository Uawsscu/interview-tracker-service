package routers

import (
	"interview-tracker/internal/adapters/handlers"

	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	r.GET("/health", handlers.GetHealth)
}
