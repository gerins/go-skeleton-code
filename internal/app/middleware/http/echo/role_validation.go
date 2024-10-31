package http

import (
	"context"

	"go-skeleton-code/pkg/log"

	"github.com/labstack/echo/v4"

	serverError "go-skeleton-code/pkg/error"
	"go-skeleton-code/pkg/jwt"
	"go-skeleton-code/pkg/response"
)

func ValidateRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get parent context from Echo Locals
			ctx, ok := c.Get("ctx").(context.Context)
			if !ok {
				ctx = context.Background()
			}

			jwtPayload := jwt.GetPayloadFromContext(ctx)

			for _, role := range roles {
				if role == jwtPayload.Role {
					return next(c)
				}
			}

			log.Context(ctx).Warn("unauthorized role")
			return response.Failed(c, serverError.ErrUnauthorized(nil))
		}
	}
}
