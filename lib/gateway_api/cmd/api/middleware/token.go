package middleware

import (
	"context"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/auth0"
	"github.com/labstack/echo/v4"
)

func AddTokenToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			headers := c.Request().Header
			token := headers.Get(echo.HeaderAuthorization)
			if token == "" {
				return fmt.Errorf("no token")
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), auth0.TokenCtxKey, token)))
			return next(c)
		}
	}
}
