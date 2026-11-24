package dto

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required,min=10000"`
}
