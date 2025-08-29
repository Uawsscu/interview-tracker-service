package routes

import (
	"interview-tracker/internal/middleware"
	"interview-tracker/internal/routes/routers"
	"time"

	"github.com/gin-gonic/gin"
)

func Listen(eg *gin.Engine) {
	interviewTrackerGroup := eg.Group("/interview-tracker")
	interviewTrackerGroup.Use(middleware.RateLimitPerMinute(60, 20, 10*time.Minute))

	routers.Health(interviewTrackerGroup)
	routers.User(interviewTrackerGroup)
	routers.Auth(interviewTrackerGroup)
	routers.Card(interviewTrackerGroup)
}
