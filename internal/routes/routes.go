package routes

import (
	"interview-tracker/internal/routes/routers"

	"github.com/gin-gonic/gin"
)

func Listen(eg *gin.Engine) {
	interviewTrackerGroup := eg.Group("/interview-tracker")

	routers.Health(interviewTrackerGroup)
	routers.User(interviewTrackerGroup)
	routers.Auth(interviewTrackerGroup)
}
