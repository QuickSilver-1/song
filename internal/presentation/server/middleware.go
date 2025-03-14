package server

import (
	"fmt"
	"song/internal/presentation/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware возвращает middleware, который логирует информацию о запросах
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Logger.Info(fmt.Sprintf("Completed %s %s with %d in %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start)))
	}
}
