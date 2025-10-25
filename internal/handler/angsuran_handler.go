package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/model"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AngsuranHandler exposes angsuran CRUD endpoints with verification
type AngsuranHandler struct {
	service *service.AngsuranService
}

// NewAngsuranHandler returns a new AngsuranHandler
func NewAngsuranHandler(s *service.AngsuranService) *AngsuranHandler {
	return &AngsuranHandler{service: s}
}

// Create handles installment payment creation
func (h *AngsuranHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var input struct {
		PinjamanID uint    `json:"pinjaman_id" binding:"required"`
		AngsuranKe int     `json:"angsuran_ke"` // Made optional - will be auto-generated if not provided
		Pokok      float64 `json:"pokok" binding:"required,gt=0"`
		Bunga      float64 `json:"bunga" binding:"required,gte=0"`
		Denda      float64 `json:"denda"`
		TotalBayar float64 `json:"total_bayar"`
		UserID     uint    `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	a := &model.Angsuran{
		PinjamanID: input.PinjamanID,
		AngsuranKe: input.AngsuranKe,
		Pokok:      input.Pokok,
		Bunga:      input.Bunga,
		Denda:      input.Denda,
		TotalBayar: input.TotalBayar,
		UserID:     input.UserID,
	}

	if err := h.service.Create(userID, role, a); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "pinjaman not found" {
			status = http.StatusBadRequest
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseSuccess("Angsuran created"))
}

// List returns filtered list of installments based on role and optional pinjaman filter
func (h *AngsuranHandler) List(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	// Optional pinjaman_id filter
	var pinjamanID uint
	if pinjamanIDStr := c.Query("pinjaman_id"); pinjamanIDStr != "" {
		if id64, err := strconv.ParseUint(pinjamanIDStr, 10, 64); err == nil {
			pinjamanID = uint(id64)
		}
	}

	list, err := h.service.List(userID, role, pinjamanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Detail returns a single installment with access control
func (h *AngsuranHandler) Detail(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	item, err := h.service.Get(userID, role, uint(id64))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Update modifies an existing installment
func (h *AngsuranHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	var input struct {
		Pokok      float64 `json:"pokok"`
		Bunga      float64 `json:"bunga"`
		Denda      float64 `json:"denda"`
		TotalBayar float64 `json:"total_bayar"`
		Status     string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	// Validate status if provided
	if input.Status != "" {
		validStatuses := map[string]bool{"proses": true, "verified": true, "kurang": true, "lebih": true}
		if !validStatuses[input.Status] {
			c.JSON(http.StatusBadRequest, utils.ResponseError("invalid status"))
			return
		}
	}

	payload := &model.Angsuran{
		Pokok:      input.Pokok,
		Bunga:      input.Bunga,
		Denda:      input.Denda,
		TotalBayar: input.TotalBayar,
		Status:     input.Status,
	}

	updated, err := h.service.Update(userID, role, uint(id64), payload)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// Delete removes an installment
func (h *AngsuranHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	if err := h.service.Delete(userID, role, uint(id64)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Angsuran deleted"))
}

// Verify allows admin to verify payment and change status
func (h *AngsuranHandler) Verify(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	var input struct {
		Status string `json:"status" binding:"required,oneof=verified kurang lebih"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	verified, err := h.service.VerifyPayment(userID, role, uint(id64), input.Status)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "invalid status for verification" {
			status = http.StatusBadRequest
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment verification updated",
		"data":    verified,
	})
}

// GetPendingPayments returns installments awaiting verification (admin only)
func (h *AngsuranHandler) GetPendingPayments(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	pending, err := h.service.GetPendingPayments(userID, role)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pending})
}
