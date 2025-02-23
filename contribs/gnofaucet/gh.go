package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v64/github"
)

func getGithubMiddleware(clientID, secret string, cooldown time.Duration) func(next http.Handler) http.Handler {
	coolDownLimiter := NewCooldownLimiter(cooldown)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// github Oauth flow is enabled
				if secret == "" || clientID == "" {
					// Continue with serving the faucet request
					next.ServeHTTP(w, r)

					return
				}

				code := r.URL.Query().Get("code")
				if code == "" {
					http.Error(w, "missing code", http.StatusBadRequest)
					return
				}

				user, err := exchangeCodeForUser(r.Context(), secret, clientID, code)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				// Just check if given account have asked for faucet before the cooldown period
				if !coolDownLimiter.CheckCooldown(user.GetLogin()) {
					http.Error(w, "user is on cooldown", http.StatusTooManyRequests)
					return
				}

				// Possibility to have more conditions like accountAge, commits, pullRequest etc

				next.ServeHTTP(w, r)
			},
		)
	}
}

type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

var exchangeCodeForUser = func(ctx context.Context, secret, clientID, code string) (*github.User, error) {
	url := "https://github.com/login/oauth/access_token"
	body := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s", clientID, secret, code)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("unable to exchange code for token")
	}

	ghClient := github.NewClient(http.DefaultClient).WithAuthToken(tokenResponse.AccessToken)
	user, _, err := ghClient.Users.Get(ctx, "")
	return user, err
}
