package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/model"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PinjamanHandler exposes pinjaman CRUD endpoints
type PinjamanHandler struct {
	service *service.PinjamanService
}

// NewPinjamanHandler returns a new PinjamanHandler
func NewPinjamanHandler(s *service.PinjamanService) *PinjamanHandler {
	return &PinjamanHandler{service: s}
}

// Create handles loan creation
func (h *PinjamanHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var input struct {
		KodePinjaman        string  `json:"kode_pinjaman"`
		UserID              uint    `json:"user_id"`
		JumlahPinjaman      float64 `json:"jumlah_pinjaman" binding:"required,gt=0"`
		BungaOptionID       *uint   `json:"bunga_option_id"`              // New: Select from available options
		BungaPersen         float64 `json:"bunga_persen" binding:"gte=0"` // Optional: Will be auto-filled if bunga_option_id is provided
		LamaBulan           int     `json:"lama_bulan" binding:"required,gt=0"`
		JumlahAngsuran      float64 `json:"jumlah_angsuran" binding:"required,gt=0"`
		Status              string  `json:"status"`
		NoRekeningPencairan string  `json:"no_rekening_pencairan"` // Account number for disbursement
		BankName            string  `json:"bank_name"`             // Bank name for disbursement
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	// Validate that either bunga_option_id or bunga_persen is provided
	if input.BungaOptionID == nil && input.BungaPersen == 0 {
		c.JSON(http.StatusBadRequest, utils.ResponseError("either bunga_option_id or bunga_persen must be provided"))
		return
	}

	// If UserID not specified, default to requestor's ID
	if input.UserID == 0 {
		input.UserID = userID
	}

	p := &model.Pinjaman{
		KodePinjaman:        input.KodePinjaman,
		UserID:              input.UserID,
		JumlahPinjaman:      input.JumlahPinjaman,
		BungaOptionID:       input.BungaOptionID,
		BungaPersen:         input.BungaPersen,
		LamaBulan:           input.LamaBulan,
		JumlahAngsuran:      input.JumlahAngsuran,
		Status:              input.Status,
		NoRekeningPencairan: input.NoRekeningPencairan,
		BankName:            input.BankName,
	}

	if err := h.service.Create(userID, role, p); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.ResponseSuccess("Pinjaman created"))
}

// List returns filtered list of loans based on role
func (h *PinjamanHandler) List(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	list, err := h.service.List(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Detail returns a single loan with access control
func (h *PinjamanHandler) Detail(c *gin.Context) {
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

// Update modifies an existing loan
func (h *PinjamanHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	var input struct {
		JumlahPinjaman float64 `json:"jumlah_pinjaman"`
		BungaPersen    float64 `json:"bunga_persen"`
		LamaBulan      int     `json:"lama_bulan"`
		JumlahAngsuran float64 `json:"jumlah_angsuran"`
		SisaAngsuran   int     `json:"sisa_angsuran"`
		Status         string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	// Validate status if provided
	if input.Status != "" {
		validStatuses := map[string]bool{"proses": true, "disetujui": true, "ditolak": true, "lunas": true, "macet": true}
		if !validStatuses[input.Status] {
			c.JSON(http.StatusBadRequest, utils.ResponseError("invalid status"))
			return
		}
	}

	payload := &model.Pinjaman{
		JumlahPinjaman: input.JumlahPinjaman,
		BungaPersen:    input.BungaPersen,
		LamaBulan:      input.LamaBulan,
		JumlahAngsuran: input.JumlahAngsuran,
		SisaAngsuran:   input.SisaAngsuran,
		Status:         input.Status,
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

// Delete removes a loan
func (h *PinjamanHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, utils.ResponseSuccess("Pinjaman deleted"))
}

// Approve approves a loan (admin can approve loans from their registered users)
func (h *PinjamanHandler) Approve(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	var input struct {
		Reason string `json:"reason"` // Optional reason for approval
	}

	// Parse input (optional)
	c.ShouldBindJSON(&input)

	if err := h.service.Approve(userID, role, uint(id64), input.Reason); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "pinjaman not found" || err.Error() == "invalid pinjaman status" {
			status = http.StatusBadRequest
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Pinjaman approved successfully"))
}

// Reject rejects a loan (admin can reject loans from their registered users)
func (h *PinjamanHandler) Reject(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	var input struct {
		Reason string `json:"alasan" binding:"required"` // Required reason for rejection
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	if err := h.service.Reject(userID, role, uint(id64), input.Reason); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "pinjaman not found" || err.Error() == "invalid pinjaman status" {
			status = http.StatusBadRequest
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Pinjaman rejected successfully"))
}
