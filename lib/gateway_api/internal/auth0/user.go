package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cobbinma/booking-platform/lib/gateway_api/models"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"time"
)

type ctxKey string

const TokenCtxKey ctxKey = "token-ctx-key"

type userService struct {
	baseURL string
	client  *http.Client
}

func NewUserService(baseURL string) *userService {
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] != '/' {
		baseURL = baseURL + "/"
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &userService{baseURL: baseURL, client: client}
}

func (us *userService) GetUser(ctx context.Context) (*models.User, error) {
	token, ok := ctx.Value(TokenCtxKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not in context")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%suserinfo", us.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("could not construct request : %w", err)
	}
	req.Header.Add(echo.HeaderAuthorization, token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request ; %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code '%v' received", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response ; %w", err)
	}

	user := models.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("could not unmarshall : %w", err)
	}

	return &user, nil
}
