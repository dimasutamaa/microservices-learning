package dto

type WalletRequest struct {
	UserID string `json:"userID" binding:"required"`
}

type WalletResponse struct {
	UserID  string `json:"userID"`
	Balance int64  `json:"balance"`
}

type CreditWalletRequest struct {
	Amount int64 `json:"amount" binding:"required,min=10000"`
}
