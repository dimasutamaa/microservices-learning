package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type WalletClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewWalletClient(baseURL string) *WalletClient {
	return &WalletClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type CreateWalletRequest struct {
	UserID string `json:"userID"`
}

func (a *WalletClient) CreateWallet(userID string) error {
	serviceURL := a.baseURL + "/wallet/create"

	data := CreateWalletRequest{UserID: userID}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, serviceURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := a.httpClient.Do(request)
	if err != nil {
		return ErrWalletServiceUnavailable
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return ErrWalletService
	}

	return nil
}

var (
	ErrWalletServiceUnavailable = errors.New("Wallet service unavailable")
	ErrWalletService            = errors.New("Wallet service error")
)
