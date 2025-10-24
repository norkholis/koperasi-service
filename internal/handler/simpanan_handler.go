package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SimpananHandler struct {
	service *service.SimpananService
}

func NewSimpananHandler(svc *service.SimpananService) *SimpananHandler {
	return &SimpananHandler{service: svc}
}

// GetWallets returns user wallets (all 3 types)
func (h *SimpananHandler) GetWallets(c *gin.Context) {
	requestorID := c.GetUint("userID")
	requestorRole := c.GetString("role")

	// Get user ID from query param or use requestor's ID
	userIDParam := c.DefaultQuery("user_id", "")
	var userID uint

	if userIDParam != "" {
		id64, err := strconv.ParseUint(userIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("invalid user_id"))
			return
		}
		userID = uint(id64)
	} else {
		userID = requestorID
	}

	wallets, err := h.service.GetUserWallets(userID, requestorID, requestorRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallets})
}

// GetAllWallets returns all wallets (admin only)
func (h *SimpananHandler) GetAllWallets(c *gin.Context) {
	requestorRole := c.GetString("role")

	wallets, err := h.service.GetAllWallets(requestorRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallets})
}

// TopupWallet creates a pending top-up transaction
func (h *SimpananHandler) TopupWallet(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Type        string  `json:"type" binding:"required,oneof=pokok wajib sukarela"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	if err := h.service.TopupWallet(userID, input.Type, input.Amount, input.Description); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseSuccess("Top-up request created, waiting for admin verification"))
}

// GetWalletDetail returns detailed wallet information
func (h *SimpananHandler) GetWalletDetail(c *gin.Context) {
	requestorID := c.GetUint("userID")
	requestorRole := c.GetString("role")

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	wallet, err := h.service.GetWalletDetail(uint(id64), requestorID, requestorRole)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

// GetWalletTransactions returns transaction history for a wallet
func (h *SimpananHandler) GetWalletTransactions(c *gin.Context) {
	requestorID := c.GetUint("userID")
	requestorRole := c.GetString("role")

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid wallet id"))
		return
	}

	transactions, err := h.service.GetWalletTransactions(uint(id64), requestorID, requestorRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

// VerifyTransaction verifies a pending transaction (admin only)
func (h *SimpananHandler) VerifyTransaction(c *gin.Context) {
	adminID := c.GetUint("userID")
	adminRole := c.GetString("role")

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid transaction id"))
		return
	}

	var input struct {
		Approve bool `json:"approve" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	if err := h.service.VerifyTransaction(uint(id64), adminID, adminRole, input.Approve); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	message := "Transaction approved"
	if !input.Approve {
		message = "Transaction rejected"
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess(message))
}

// AdjustWallet allows admin to directly adjust wallet balance
func (h *SimpananHandler) AdjustWallet(c *gin.Context) {
	adminID := c.GetUint("userID")
	adminRole := c.GetString("role")

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid wallet id"))
		return
	}

	var input struct {
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	if err := h.service.AdjustWalletBalance(uint(id64), input.Amount, input.Description, adminID, adminRole); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Wallet balance adjusted"))
}

// GetPendingTransactions returns all pending transactions (admin only)
func (h *SimpananHandler) GetPendingTransactions(c *gin.Context) {
	adminRole := c.GetString("role")

	transactions, err := h.service.GetPendingTransactions(adminRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
