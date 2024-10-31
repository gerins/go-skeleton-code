package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// ValidateJwtToken is a Gin middleware to validate JWT tokens from the Authorization header.
func SetAPITimeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the context from the request
		ctx := c.Request.Context()

		newCtx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()

		// Apply new context
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
