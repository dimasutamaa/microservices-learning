package services

import (
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	tokenRepo *repositories.AccessTokenRepository
	rdb       *redis.Client
	cfg       *config.GlobalConfig
}

func NewUserService(userRepo *repositories.UserRepository, tokenRepo *repositories.AccessTokenRepository,
	rdb *redis.Client, cfg *config.GlobalConfig) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		rdb:       rdb,
		cfg:       cfg,
	}
}

func (s *AuthService) Register(request *dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.NewString(),
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashed),
	}

	if _, err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		UserID:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Message: "User registered successfully",
	}, nil
}

func (s *AuthService) Login(request *dto.LoginRequest) (interface{}, error) {
	user, err := s.userRepo.FindByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	tokenID := uuid.NewString()
	expiressAt := time.Now().Add(s.cfg.TokenTTL)
	accessToken := &models.AccessToken{
		ID:        tokenID,
		UserID:    user.ID,
		Revoked:   false,
		ExpiresAt: expiressAt,
	}
	if err := s.tokenRepo.Create(accessToken); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(accessToken, s.cfg.JWTSecret, s.cfg.TokenTTL)
	if err != nil {
		return nil, err
	}

	// cache active token with ttl, redis errors do not prevent login
	if s.rdb != nil {
		_ = s.rdb.Set(context.Background(), cacheKey(tokenID), user.ID, s.cfg.TokenTTL).Err()
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

func (s *AuthService) Logout(tokenID string, userID string) error {
	if err := s.tokenRepo.Revoke(tokenID, userID); err != nil {
		return err
	}
	if s.rdb != nil {
		_ = s.rdb.Del(context.Background(), cacheKey(tokenID)).Err()
	}
	return nil
}

func cacheKey(id string) string {
	return "access-token:" + id
}

var (
	ErrInvalidCredentials = errors.New("Invalid email or password")
)
