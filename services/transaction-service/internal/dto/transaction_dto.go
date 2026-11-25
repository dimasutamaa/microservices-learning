package dto

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required,min=10000"`
}

type TransferRequest struct {
	ToUserID string `json:"toUserID" binding:"required,uuid"`
	Amount   int64  `json:"amount" binding:"required,min=1"`
}

type TransferResponse struct {
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
