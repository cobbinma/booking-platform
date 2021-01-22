package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func AddUserToContext(domain string, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Request().Header
		token := headers.Get(echo.HeaderAuthorization)
		if token == "" {
			return fmt.Errorf("no token")
		}

		baseURL, err := url.Parse(fmt.Sprintf("%suserinfo", domain))
		if err != nil {
			return fmt.Errorf("could not parse url : %w", err)
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		req := &http.Request{
			Method: "GET",
			URL:    baseURL,
			Header: headers,
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("could not make request ; %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code '%v' received", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response ; %w", err)
		}

		user := models.User{}
		if err := json.Unmarshal(body, &user); err != nil {
			return fmt.Errorf("could not unmarshall : %w", err)
		}

		c.Request().WithContext(context.WithValue(c.Request().Context(), models.UserCtxKey, user))
		return nil
	}
}
