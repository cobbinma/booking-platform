package auth0

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type TokenClient struct {
	baseURL  string
	clientID string
	secret   string
	log      *zap.SugaredLogger
	token    *oauth2.Token
}

func NewTokenClient(log *zap.SugaredLogger, domain string) (*TokenClient, error) {
	if len(domain) > 0 && domain[len(domain)-1] != '/' {
		domain = domain + "/"
	}
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("auth0 client id missing")
	}
	secret := os.Getenv("AUTH0_CLIENT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("auth0 secret missing")
	}

	return &TokenClient{
		baseURL:  domain,
		clientID: clientID,
		secret:   secret,
		log:      log,
	}, nil
}

func (tc *TokenClient) GetToken(log *zap.SugaredLogger, audience string) (*oauth2.Token, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%soauth/token", tc.baseURL),
		strings.NewReader(fmt.Sprintf(
			`{"client_id":"%s","client_secret":"%s","audience":"%s","grant_type":"client_credentials"}`, tc.clientID, tc.secret, audience)))
	if err != nil {
		return nil, fmt.Errorf("could not create request : %w", err)
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do request : %w", err)
	}
	defer func(log *zap.SugaredLogger) {
		if err := res.Body.Close(); err != nil {
			log.Errorf("could not close token request body : %s", err)
		}
	}(log)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request : %w", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Errorw("unexpected status code", "status code", res.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code '%v'", res.StatusCode)
	}

	resp := &oauth2.Token{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, fmt.Errorf("could not unmarshall : %w", err)
	}

	return resp, nil
}
