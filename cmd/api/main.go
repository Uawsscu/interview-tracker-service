package main

// @title Interview Tracker API
// @version 1.0
// @description API for Interview Tracker
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
import (
	"fmt"
	"interview-tracker/internal/config"
	"interview-tracker/internal/pkg/logs"
	"interview-tracker/internal/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.elastic.co/apm/module/apmgin"
)

func CORSAllow() gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	cfg.AllowAllOrigins = true
	return cors.New(cfg)
}

func main() {
	gin.SetMode(gin.ReleaseMode) // if not set releaseMode will be used debug mode

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	config.LoadConfig()
	logs.LoggerConfig()
	config.ConnectDB()
	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/interview-tracker/health"),
		gin.Recovery(),
		apmgin.Middleware(router),
	)

	router.Use(config.LogStartTimeMiddleware())

	config.SwaggerConfig(router)
	config.NewRedis()

	router.Use(CORSAllow())
	routes.Listen(router)
	router.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")))
}
