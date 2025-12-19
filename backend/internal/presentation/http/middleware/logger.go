package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/logger"
)

// Logger is a Gin middleware for logging requests using zerolog.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		l := logger.NewLogger(c.Request.Context())
		l.InfoEvent().
			Str("method", c.Request.Method).
			Str("full_path", c.FullPath()).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("remote_ip", c.RemoteIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Request")
	}
}
