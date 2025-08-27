package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealth godoc
// @Summary      Health check
// @Description  Returns success if the service is healthy
// @Tags         health
// @Produce      json
// @Success      200 {string} string "success"
// @Router       /interview-tracker/health [get]
func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "success")
}
