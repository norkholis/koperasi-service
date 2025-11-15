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

// List users (super_admin can see all, admin can see their registered users)
func (h *UserHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	users, err := h.service.ListUsers(userID, role)
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
		user := gin.H{
			"id":           u.ID,
			"email":        u.Email,
			"name":         u.Name,
			"address":      u.Address,
			"phone_number": u.PhoneNumber,
			"nik":          u.NIK,
			"role_id":      u.RoleID,
		}
		if u.AdminID != nil {
			user["admin_id"] = *u.AdminID
		}
		resp = append(resp, user)
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// Detail gets user details
func (h *UserHandler) Detail(c *gin.Context) {
	reqID := c.GetUint("user_id")
	role := c.GetString("role")

	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.service.GetUser(reqID, role, uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	response := gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"address":      user.Address,
		"phone_number": user.PhoneNumber,
		"nik":          user.NIK,
		"role_id":      user.RoleID,
	}
	if user.AdminID != nil {
		response["admin_id"] = *user.AdminID
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// Create user (super_admin and admin)
func (h *UserHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		Name        string `json:"name" binding:"required"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phone_number"`
		NIK         string `json:"nik" binding:"required"`
		RoleID      uint   `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}
	u := &model.User{
		Email:       input.Email,
		Password:    input.Password,
		Name:        input.Name,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
		NIK:         input.NIK,
		RoleID:      input.RoleID,
	}
	if err := h.service.CreateUser(userID, role, u); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	response := gin.H{
		"id":           u.ID,
		"email":        u.Email,
		"name":         u.Name,
		"address":      u.Address,
		"phone_number": u.PhoneNumber,
		"nik":          u.NIK,
		"role_id":      u.RoleID,
	}
	if u.AdminID != nil {
		response["admin_id"] = *u.AdminID
	}
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// Update user (self or any if super_admin).
func (h *UserHandler) Update(c *gin.Context) {
	role := c.GetString("role")
	reqID := c.GetUint("user_id")
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	var input struct {
		Email       string  `json:"email"`
		Name        string  `json:"name"`
		Address     string  `json:"address"`
		PhoneNumber string  `json:"phone_number"`
		NIK         string  `json:"nik"`
		Password    *string `json:"password"`
		RoleID      *uint   `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}
	u, err := h.service.UpdateUser(reqID, role, uint(id64), input.Email, input.Name, input.Address, input.PhoneNumber, input.NIK, input.Password, input.RoleID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}
	response := gin.H{
		"id":           u.ID,
		"email":        u.Email,
		"name":         u.Name,
		"address":      u.Address,
		"phone_number": u.PhoneNumber,
		"nik":          u.NIK,
		"role_id":      u.RoleID,
	}
	if u.AdminID != nil {
		response["admin_id"] = *u.AdminID
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// Delete user (self or any if super_admin).
func (h *UserHandler) Delete(c *gin.Context) {
	role := c.GetString("role")
	reqID := c.GetUint("user_id")
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
