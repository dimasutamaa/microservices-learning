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

type TopUpRequest struct {
	Amount int64 `json:"amount"`
}

type TopUpResponse struct {
	NewBalance int64 `json:"newBalance"`
}

type TransferRequest struct {
	ToUserID string `json:"toUserID"`
	Amount   int64  `json:"amount"`
}

func (a *WalletClient) Credit(userID string, amount int64) (int64, error) {
	serviceURL := a.baseURL + "/wallet/credit"

	data := TopUpRequest{Amount: amount}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodPatch, serviceURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-User-ID", userID)

	response, err := a.httpClient.Do(request)
	if err != nil {
		return 0, ErrWalletServiceUnavailable
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, ErrWalletService
	}

	var result TopUpResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.NewBalance, nil
}

var (
	ErrWalletServiceUnavailable = errors.New("Wallet service unavailable")
	ErrWalletService            = errors.New("Wallet service error")
)
