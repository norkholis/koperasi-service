package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/model"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler exposes user CRUD endpoints.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler returns a new UserHandler.
func NewUserHandler(s *service.UserService) *UserHandler { return &UserHandler{service: s} }

// List users (super_admin only).
func (h *UserHandler) List(c *gin.Context) {
	role := c.GetString("role")
	users, err := h.service.ListUsers(role)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	resp := make([]gin.H, 0, len(users))
	for _, u := range users {
		resp = append(resp, gin.H{"id": u.ID, "email": u.Email, "role_id": u.RoleID})
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Detail returns a single user respecting access control.
func (h *UserHandler) Detail(c *gin.Context) {
	role := c.GetString("role")
	reqID := c.GetUint("userID")
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	u, err := h.service.GetUser(reqID, role, uint(id64))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": u.ID, "email": u.Email, "role_id": u.RoleID}})
}

// Create user (super_admin only).
func (h *UserHandler) Create(c *gin.Context) {
	role := c.GetString("role")
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}
	u := &model.User{Email: input.Email, Password: input.Password, RoleID: input.RoleID}
	if err := h.service.CreateUser(role, u); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.ResponseSuccess("user created"))
}

// Update user (self or any if super_admin).
func (h *UserHandler) Update(c *gin.Context) {
	role := c.GetString("role")
	reqID := c.GetUint("userID")
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	var input struct {
		Email    string  `json:"email"`
		Password *string `json:"password"`
		RoleID   *uint   `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}
	u, err := h.service.UpdateUser(reqID, role, uint(id64), input.Email, input.Password, input.RoleID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": u.ID, "email": u.Email, "role_id": u.RoleID}})
}

// Delete user (self or any if super_admin).
func (h *UserHandler) Delete(c *gin.Context) {
	role := c.GetString("role")
	reqID := c.GetUint("userID")
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	if err := h.service.DeleteUser(reqID, role, uint(id64)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("user deleted"))
}
