package services

import (
	"auth-service/internal/dto"
	"auth-service/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) FindUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		UserID:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Message: "User data retrieved successfully",
	}

	return response, nil
}
