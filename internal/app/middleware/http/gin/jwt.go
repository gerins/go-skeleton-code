package http

import (
	"strings"

	"github.com/gin-gonic/gin"

	serverError "go-skeleton-code/pkg/error"
	"go-skeleton-code/pkg/jwt"
	"go-skeleton-code/pkg/log"
	response "go-skeleton-code/pkg/response/gin"
)

// ValidateJwtToken is a Gin middleware to validate JWT tokens from the Authorization header.
func ValidateJwtToken(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the context from the request
		ctx := c.Request.Context()

		// Extract the token from the Authorization header
		var tokenString string
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// Validate the JWT token
		payload, err := jwt.Validate(tokenString, secretKey)
		if err != nil {
			log.Context(ctx).Error(err)
			response.Failed(c, serverError.ErrUnauthorized(err))
			c.Abort() // Abort the request chain if validation fails
			return
		}

		// Save the JWT payload to the context
		newCtx := jwt.SavePayloadToContext(ctx, payload)
		c.Request = c.Request.WithContext(newCtx)

		// Continue to the next handler
		c.Next()
	}
}
