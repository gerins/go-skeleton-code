package http

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"

	serverError "go-skeleton-code/pkg/error"
	"go-skeleton-code/pkg/jwt"
	"go-skeleton-code/pkg/log"
	"go-skeleton-code/pkg/response/echo"
)

func ValidateJwtToken(secretKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get parent context from Echo Locals
			ctx, ok := c.Get("ctx").(context.Context)
			if !ok {
				ctx = context.Background()
			}

			// Get the Authorization header
			var tokenString string
			authHeader := c.Request().Header.Get("Authorization")

			// Extract the token string from the header
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}

			payload, err := jwt.Validate(tokenString, secretKey)
			if err != nil {
				log.Context(ctx).Error(err)
				return response.Failed(c, serverError.ErrUnauthorized(err))
			}

			ctx = jwt.SavePayloadToContext(ctx, payload)

			c.Set("ctx", ctx)
			return next(c)
		}
	}
}
