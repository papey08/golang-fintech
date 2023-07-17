package httpgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware is my own middleware for logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestTime := time.Now().UTC()
		c.Next()
		path := c.Request.URL.Path
		method := c.Request.Method
		status := c.Writer.Status()
		log.Printf("%v | %d | %-10s | %s\n",
			requestTime.Format("2006/01/02 - 15:04:05"),
			status,
			method,
			path)
	}
}

// PanicMiddleware is my own middleware for panic
func PanicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"data": nil, "error": fmt.Sprintf("%v", err)})
			}
		}()
		c.Next()
	}
}
