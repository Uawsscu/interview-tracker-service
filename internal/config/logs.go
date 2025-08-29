package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interview-tracker/internal/pkg/logs"
	"io"

	"github.com/gin-gonic/gin"
)

func LogStartTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// before request
		if c.FullPath() != "/interview-tracker/health" {
			logApi(c)
		}
		// after request
		c.Next()
	}
}

func logApi(c *gin.Context) {
	logs.LoggerInfo(fmt.Sprintf("Request URL: %s %s", c.Request.Method, c.FullPath()))
	headers, _ := json.Marshal(c.Request.Header)
	logs.LoggerInfo(fmt.Sprintf("Headers: %s", string(headers)))

	if len(c.Request.URL.Query()) > 0 {
		queryParams, _ := json.Marshal(c.Request.URL.Query())
		logs.LoggerInfo(fmt.Sprintf("Query Parameters: %s", string(queryParams)))
	}

	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			var formattedBody bytes.Buffer
			if json.Compact(&formattedBody, bodyBytes) == nil {
				logs.LoggerInfo(fmt.Sprintf("Request Body: %s", formattedBody.String()))
			} else {
				logs.LoggerInfo(fmt.Sprintf("Request Body (raw): %s", string(bodyBytes)))
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		} else {
			logs.LoggerInfo(fmt.Sprintf("Failed to read request body: %v", err))
		}
	}
}
