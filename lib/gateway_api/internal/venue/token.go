package venue

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

type tokenClient struct {
	baseURL  string
	clientID string
	secret   string
	log      *zap.SugaredLogger
}

func token(log *zap.SugaredLogger) (*oauth2.Token, error) {
	baseURL := os.Getenv("AUTH0_DOMAIN")
	if baseURL == "" {
		return nil, fmt.Errorf("auth0 domain url missing")
	}
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] != '/' {
		baseURL = baseURL + "/"
	}
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("auth0 client id missing")
	}
	secret := os.Getenv("AUTH0_CLIENT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("auth0 secret missing")
	}

	tk := &tokenClient{
		baseURL:  baseURL,
		clientID: clientID,
		secret:   secret,
		log:      log,
	}

	return tk.getToken()
}

func (tc *tokenClient) getToken() (*oauth2.Token, error) {
	url := fmt.Sprintf("%soauth/token", tc.baseURL)

	payload := strings.NewReader(fmt.Sprintf(`{"client_id":"%s","client_secret":"%s","audience":"http://venue","grant_type":"client_credentials"}`, tc.clientID, tc.secret))

	req, err := http.NewRequest("POST", url, payload)
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
	}(tc.log)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request : %w", err)
	}

	if res.StatusCode != http.StatusOK {
		tc.log.Errorw("unexpected status code", "status code", res.StatusCode, "body", string(body))
		return nil, fmt.Errorf("unexpected status code '%v'", res.StatusCode)
	}

	resp := &oauth2.Token{}

	if err := json.Unmarshal(body, resp); err != nil {
		return nil, fmt.Errorf("could not unmarshall : %w", err)
	}

	return resp, nil
}
