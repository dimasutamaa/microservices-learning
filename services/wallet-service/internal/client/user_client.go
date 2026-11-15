package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type UserResponse struct {
	UserID  string `json:"userID"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (a *UserClient) FindUserByID(userID string) (*UserResponse, error) {
	serviceURL := a.baseURL + "/auth/user/" + userID
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, serviceURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, ErrAuthServiceUnavailable
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var result UserResponse
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil

	case http.StatusNotFound:
		return nil, ErrUserNotFound

	default:
		return nil, ErrAuthService
	}
}

var (
	ErrAuthServiceUnavailable = errors.New("Auth service unavailable")
	ErrUserNotFound           = errors.New("User not found")
	ErrAuthService            = errors.New("Auth service error")
)
