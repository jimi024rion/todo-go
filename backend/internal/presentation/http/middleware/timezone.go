package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimi024rion/todo-go/backend/internal/config/clock"
)

// Timezone returns a middleware that sets the timezone location in the context
// based on the X-Timezone header.
func Timezone() gin.HandlerFunc {
	return func(c *gin.Context) {
		tz := c.GetHeader("X-Timezone")
		if tz == "" {
			c.Next()
			return
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			// Fallback to UTC if invalid timezone is provided
			c.Next()
			return
		}

		// Update request context with the location
		ctx := clock.ContextWithLocation(c.Request.Context(), loc)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
