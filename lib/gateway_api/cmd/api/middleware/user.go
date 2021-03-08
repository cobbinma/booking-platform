package middleware

import (
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/labstack/echo/v4"
)

func User(service models.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			user, err := service.GetUser(ctx)
			if err != nil {
				return fmt.Errorf("could not get user from service : %w", err)
			}

			c.SetRequest(c.Request().WithContext(models.AddUserToContext(ctx, *user)))
			return next(c)
		}
	}
}
