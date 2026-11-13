package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type ValidateTokenResponse struct {
	Valid bool `json:"valid"`
}

func (a *AuthClient) IsTokenUsable(c context.Context, tokenID string) (bool, error) {
	checkTokenURL := a.baseURL + "auth/token/" + tokenID + "/usability"
	request, err := http.NewRequestWithContext(c, http.MethodGet, checkTokenURL, nil)
	if err != nil {
		return false, err
	}

	response, err := a.httpClient.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("auth service returned status %d", response.StatusCode)
	}

	var result ValidateTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Valid, nil
}
