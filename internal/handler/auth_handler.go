package handler

import (
	"koperasi-service/config"
	"koperasi-service/internal/model"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
	config  *config.Config
}

func NewAuthHandler(svc *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{service: svc, config: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		Name        string `json:"name" binding:"required"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phone_number"`
		NIK         string `json:"nik" binding:"required"`
		RoleID      uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	user := &model.User{
		Email:       input.Email,
		Password:    input.Password,
		Name:        input.Name,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
		NIK:         input.NIK,
		RoleID:      input.RoleID,
	}

	if err := h.service.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseSuccess("User registered"))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	token, err := h.service.Login(input.Email, input.Password, h.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id") // dari middleware JWT
	user, err := h.service.GetUserWithRole(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseError("user not found"))
		return
	}
	response := gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"address":      user.Address,
		"phone_number": user.PhoneNumber,
		"nik":          user.NIK,
		"role":         gin.H{"id": user.Role.ID, "name": user.Role.Name},
	}
	if user.AdminID != nil {
		response["admin_id"] = *user.AdminID
	}
	c.JSON(http.StatusOK, response)
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id") // from JWT middleware

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required,min=6"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	// Validate password confirmation
	if input.NewPassword != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, utils.ResponseError("new password and confirm password do not match"))
		return
	}

	// Change password
	if err := h.service.ChangePassword(userID, input.CurrentPassword, input.NewPassword); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "current password is incorrect" {
			status = http.StatusBadRequest
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Password changed successfully"))
}

// ForgotPassword handles password reset requests
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input struct {
		Email           string `json:"email" binding:"required,email"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	// Validate password confirmation
	if input.NewPassword != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, utils.ResponseError("new password and confirm password do not match"))
		return
	}

	// Reset password
	if err := h.service.ResetPassword(input.Email, input.NewPassword); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Password reset successfully"))
}
