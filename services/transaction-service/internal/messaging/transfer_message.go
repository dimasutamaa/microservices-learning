package messaging

type TransferMessage struct {
	TransactionID string `json:"transactionID"`
	FromUserID    string `json:"fromUserID"`
	ToUserID      string `json:"toUserID"`
	Amount        int64  `json:"amount"`
}

type TransferResult struct {
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
}
