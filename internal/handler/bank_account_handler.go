package handler

import (
	"koperasi-service/internal/model"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankAccountHandler struct {
	service *service.BankAccountService
}

func NewBankAccountHandler(service *service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{service: service}
}

// Create handles bank account creation (admin only)
func (h *BankAccountHandler) Create(c *gin.Context) {
	var input struct {
		BankName      string `json:"bank_name" binding:"required"`
		AccountNumber string `json:"account_number" binding:"required"`
		AccountName   string `json:"account_name" binding:"required"`
		Description   string `json:"description"`
		IsActive      *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	account := &model.BankAccount{
		BankName:      input.BankName,
		AccountNumber: input.AccountNumber,
		AccountName:   input.AccountName,
		Description:   input.Description,
		IsActive:      true,
	}

	if input.IsActive != nil {
		account.IsActive = *input.IsActive
	}

	if err := h.service.Create(account); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bank account created successfully", "data": account})
}

// Update handles bank account update (admin only)
func (h *BankAccountHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	var input struct {
		BankName      string `json:"bank_name"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		Description   string `json:"description"`
		IsActive      bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	account := &model.BankAccount{
		BankName:      input.BankName,
		AccountNumber: input.AccountNumber,
		AccountName:   input.AccountName,
		Description:   input.Description,
		IsActive:      input.IsActive,
	}

	if err := h.service.Update(uint(id), account); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank account not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Bank account updated successfully"))
}

// Delete handles bank account deletion (admin only)
func (h *BankAccountHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bank account not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Bank account deleted successfully"))
}

// GetByID handles retrieving a single bank account
func (h *BankAccountHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	account, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseError("Bank account not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank account retrieved successfully", "data": account})
}

// List handles retrieving all bank accounts (admin sees all, users see only active)
func (h *BankAccountHandler) List(c *gin.Context) {
	role := c.GetString("role")

	var accounts []model.BankAccount
	var err error

	if role == "super_admin" || role == "admin" {
		accounts, err = h.service.GetAll()
	} else {
		accounts, err = h.service.GetActive()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank accounts retrieved successfully", "data": accounts})
}

// GetActive handles retrieving only active bank accounts (for users to see where to send money)
func (h *BankAccountHandler) GetActive(c *gin.Context) {
	accounts, err := h.service.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Active bank accounts retrieved successfully", "data": accounts})
}
