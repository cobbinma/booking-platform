package middleware

import (
	"fmt"
	"github.com/auth0-community/go-auth0"
	auth "github.com/cobbinma/booking-platform/lib/gateway_api/internal/auth0"
	"github.com/labstack/echo/v4"
	"gopkg.in/square/go-jose.v2"
	"net/http"
)

func Auth(domain string, apiIdentifier string) echo.MiddlewareFunc {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: fmt.Sprintf("%s.well-known/jwks.json", domain)}, nil)
	validator := auth0.NewValidator(auth0.NewConfiguration(client, []string{apiIdentifier}, domain, jose.RS256), nil)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			headers := c.Request().Header
			token := headers.Get(echo.HeaderAuthorization)
			if token == "" {
				return c.JSONBlob(http.StatusBadRequest, []byte(`{"error": "token not given"}`))
			}

			_, err := validator.ValidateRequest(c.Request())
			if err != nil {
				return c.JSONBlob(http.StatusUnauthorized, []byte(`{"error": "invalid token"}`))
			}

			c.SetRequest(c.Request().WithContext(auth.AddTokenToCtx(c.Request().Context(), token)))
			return next(c)
		}
	}
}
