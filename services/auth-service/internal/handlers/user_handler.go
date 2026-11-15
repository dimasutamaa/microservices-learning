package handlers

import (
	"auth-service/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /user/{id}
func (h *UserHandler) FindUserByID(c *gin.Context) {
	user, err := h.userService.FindUserByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, "Unexpected error occured")
		return
	}

	c.JSON(http.StatusOK, user)
}
