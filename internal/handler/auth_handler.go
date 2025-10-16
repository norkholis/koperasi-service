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
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RoleID   uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	user := &model.User{
		Email:    input.Email,
		Password: input.Password,
		RoleID:   input.RoleID,
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
	userID := c.GetUint("userID") // dari middleware JWT
	user, err := h.service.GetUserWithRole(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseError("user not found"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  gin.H{"id": user.Role.ID, "name": user.Role.Name},
	})
}
