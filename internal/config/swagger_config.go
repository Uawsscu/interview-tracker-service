package config

import (
	"interview-tracker/docs"
	"interview-tracker/internal/pkg/logs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Interview Tracker
// @description Interview Tracker APIs document
func SwaggerConfig(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Interview Tracker"
	docs.SwaggerInfo.Description = "Interview Tracker APIs document"

	router.GET("/interview-tracker/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/interview-tracker/swagger/doc.json"),
	))

	logs.Logger.Printf("swagger| url: /interview-tracker/swagger/index.html")
}
