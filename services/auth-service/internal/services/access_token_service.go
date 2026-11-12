package services

import (
	"auth-service/internal/repositories"
	"time"
)

type AccessTokenService struct {
	tokenRepo *repositories.AccessTokenRepository
}

func NewAccessTokenService(tokenRepo *repositories.AccessTokenRepository) *AccessTokenService {
	return &AccessTokenService{tokenRepo: tokenRepo}
}

func (s *AccessTokenService) IsTokenUsable(tokenID string) (bool, error) {
	if _, err := s.tokenRepo.IsUsable(tokenID, time.Now()); err != nil {
		return false, err
	}
	return true, nil
}
