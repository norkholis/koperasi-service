package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SHUAnggotaHandler exposes SHU Anggota endpoints
type SHUAnggotaHandler struct {
	service *service.SHUAnggotaService
}

// NewSHUAnggotaHandler returns a new SHUAnggotaHandler
func NewSHUAnggotaHandler(s *service.SHUAnggotaService) *SHUAnggotaHandler {
	return &SHUAnggotaHandler{service: s}
}

// SaveUserSHU saves the calculated SHU for a specific user
func (h *SHUAnggotaHandler) SaveUserSHU(c *gin.Context) {
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

	shuAnggota, err := h.service.SaveUserSHU(role, requestorUserID, uint(userID64), input.Tahun)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "SHU record not found for the specified year" {
			status = http.StatusNotFound
		} else if err.Error() == "SHU record already exists for this user and year" {
			status = http.StatusConflict
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User SHU saved successfully",
		"data":    shuAnggota,
	})
}

// GetUserSHU retrieves saved SHU data for a user
func (h *SHUAnggotaHandler) GetUserSHU(c *gin.Context) {
	role := c.GetString("role")
	requestorUserID := c.GetUint("userID")
	userIDParam := c.Param("user_id")
	tahunParam := c.Param("tahun")

	userID64, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid user_id"))
		return
	}

	tahun64, err := strconv.ParseInt(tahunParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid tahun"))
		return
	}

	shuAnggota, err := h.service.GetUserSHU(role, requestorUserID, uint(userID64), int(tahun64))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "SHU record not found for the specified year" || err.Error() == "user SHU record not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": shuAnggota,
	})
}

// GetUserSHUHistory retrieves all SHU history for a user
func (h *SHUAnggotaHandler) GetUserSHUHistory(c *gin.Context) {
	role := c.GetString("role")
	requestorUserID := c.GetUint("userID")
	userIDParam := c.Param("user_id")

	userID64, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid user_id"))
		return
	}

	shuHistory, err := h.service.GetUserSHUHistory(role, requestorUserID, uint(userID64))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": shuHistory,
	})
}

// List retrieves all SHU Anggota records (admin only)
func (h *SHUAnggotaHandler) List(c *gin.Context) {
	role := c.GetString("role")

	shuAnggotas, err := h.service.List(role)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": shuAnggotas,
	})
}

// GetBySHUID retrieves all user SHU records for a specific SHU (admin only)
func (h *SHUAnggotaHandler) GetBySHUID(c *gin.Context) {
	role := c.GetString("role")
	shuIDParam := c.Param("shu_id")

	shuID64, err := strconv.ParseUint(shuIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid shu_id"))
		return
	}

	shuAnggotas, err := h.service.GetBySHUID(role, uint(shuID64))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": shuAnggotas,
	})
}

// Delete removes a SHU Anggota record (admin only)
func (h *SHUAnggotaHandler) Delete(c *gin.Context) {
	role := c.GetString("role")
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}

	if err := h.service.Delete(role, uint(id64)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		} else if err.Error() == "SHU Anggota record not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("SHU Anggota record deleted"))
}
