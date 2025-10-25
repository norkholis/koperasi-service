package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/model"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SHUHandler exposes SHU calculation and management endpoints
type SHUHandler struct {
	service *service.SHUService
}

// NewSHUHandler returns a new SHUHandler
func NewSHUHandler(s *service.SHUService) *SHUHandler {
	return &SHUHandler{service: s}
}

// GenerateReport handles SHU calculation and report generation
func (h *SHUHandler) GenerateReport(c *gin.Context) {
	role := c.GetString("role")

	var input struct {
		Tahun            int     `json:"tahun" binding:"required,min=2000,max=2100"`
		TotalSHUKoperasi float64 `json:"total_shu_koperasi" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	report, err := h.service.GenerateReport(role, input.Tahun, input.TotalSHUKoperasi)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SHU report generated successfully",
		"data":    report,
	})
}

// SaveSHU saves the calculated SHU as a record
func (h *SHUHandler) SaveSHU(c *gin.Context) {
	role := c.GetString("role")

	var input struct {
		Tahun    int     `json:"tahun" binding:"required,min=2000,max=2100"`
		TotalSHU float64 `json:"total_shu" binding:"required,gt=0"`
		Status   string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	shu, err := h.service.SaveSHU(role, input.Tahun, input.TotalSHU, input.Status)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "SHU for this year already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "SHU record saved successfully",
		"data":    shu,
	})
}

// List returns all SHU records
func (h *SHUHandler) List(c *gin.Context) {
	role := c.GetString("role")

	list, err := h.service.List(role)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Detail returns a single SHU record
func (h *SHUHandler) Detail(c *gin.Context) {
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	item, err := h.service.Get(role, uint(id64))
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

// Update modifies an existing SHU record
func (h *SHUHandler) Update(c *gin.Context) {
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	var input struct {
		TotalSHU float64 `json:"total_shu"`
		Status   string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	// Validate status if provided
	if input.Status != "" {
		if input.Status != "draft" && input.Status != "final" {
			c.JSON(http.StatusBadRequest, utils.ResponseError("status must be 'draft' or 'final'"))
			return
		}
	}

	payload := &model.SHUTahunan{
		TotalSHU: input.TotalSHU,
		Status:   input.Status,
	}

	updated, err := h.service.Update(role, uint(id64), payload)
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

// Delete removes a SHU record
func (h *SHUHandler) Delete(c *gin.Context) {
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	if err := h.service.Delete(role, uint(id64)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("SHU record deleted"))
}

// GetByTahun returns SHU record by year
func (h *SHUHandler) GetByTahun(c *gin.Context) {
	role := c.GetString("role")
	tahunParam := c.Param("tahun")

	tahun64, err := strconv.ParseInt(tahunParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid tahun"))
		return
	}

	item, err := h.service.GetByTahun(role, int(tahun64))
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

// GenerateUserSHU handles SHU calculation for a specific user
func (h *SHUHandler) GenerateUserSHU(c *gin.Context) {
	role := c.GetString("role")
	requestorUserID := c.GetUint("userID")
	userIDParam := c.Param("user_id")

	userID64, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid user_id"))
		return
	}

	var input struct {
		Tahun int `json:"tahun" binding:"required,min=2000,max=2100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	userSHU, err := h.service.GenerateUserSHU(role, requestorUserID, uint(userID64), input.Tahun)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "SHU record not found for the specified year. Please generate SHU report first" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User SHU calculated successfully",
		"data":    userSHU,
	})
}

// GenerateReportWithExpenses handles automated SHU calculation based on income and expenses
func (h *SHUHandler) GenerateReportWithExpenses(c *gin.Context) {
	role := c.GetString("role")

	var input struct {
		Tahun               int     `json:"tahun" binding:"required,min=2000,max=2100"`
		BebanOperasional    float64 `json:"beban_operasional" binding:"min=0"`
		BebanNonOperasional float64 `json:"beban_non_operasional" binding:"min=0"`
		BebanPajak          float64 `json:"beban_pajak" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	report, err := h.service.GenerateReportWithExpenses(role, input.Tahun, input.BebanOperasional, input.BebanNonOperasional, input.BebanPajak)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Automated SHU report generated successfully",
		"data":    report,
	})
}

// SaveSHUWithExpenses saves the automated SHU calculation with detailed income and expense information
func (h *SHUHandler) SaveSHUWithExpenses(c *gin.Context) {
	role := c.GetString("role")

	var input struct {
		Tahun                    int     `json:"tahun" binding:"required,min=2000,max=2100"`
		PendapatanOperasional    float64 `json:"pendapatan_operasional" binding:"min=0"`
		PendapatanNonOperasional float64 `json:"pendapatan_non_operasional" binding:"min=0"`
		BebanOperasional         float64 `json:"beban_operasional" binding:"min=0"`
		BebanNonOperasional      float64 `json:"beban_non_operasional" binding:"min=0"`
		BebanPajak               float64 `json:"beban_pajak" binding:"min=0"`
		TotalSHU                 float64 `json:"total_shu" binding:"required,gt=0"`
		Status                   string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	shu, err := h.service.SaveSHUWithExpenses(role, input.Tahun, input.PendapatanOperasional, input.PendapatanNonOperasional, input.BebanOperasional, input.BebanNonOperasional, input.BebanPajak, input.TotalSHU, input.Status)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "SHU for this year already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Automated SHU record saved successfully",
		"data":    shu,
	})
}
