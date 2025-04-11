// trade-service/utils/token_client.go
package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	mu          sync.Mutex
	cachedToken string
	expiry      time.Time
)

func GetMachineToken() (string, error) {
	clientID := os.Getenv("TRADE_SERVICE_CLIENT_ID")
	clientSecret := os.Getenv("TRADE_SERVICE_CLIENT_SECRET")
	tokenURL := os.Getenv("AUTH_URL") + "/oauth/introspect"

	mu.Lock()
	defer mu.Unlock()
	if time.Now().Before(expiry) && cachedToken != "" {
		return cachedToken, nil
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("bad status: " + resp.Status)
	}
	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	cachedToken = body.AccessToken
	expiry = time.Now().Add(time.Duration(body.ExpiresIn-10) * time.Second)
	return cachedToken, nil
}
