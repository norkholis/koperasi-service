package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/model"
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

func (h *SimpananHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	// Only member/admin/super_admin allowed (all roles) so no restriction here unless future change.
	var input struct {
		Type        string  `json:"type" binding:"required,oneof=wajib sukarela"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}
	_ = role // currently unused but kept for future logic
	s := &model.Simpanan{
		UserID:      userID,
		Type:        input.Type,
		Amount:      input.Amount,
		Description: input.Description,
	}
	if err := h.service.Create(s); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.ResponseSuccess("Simpanan created"))
}

func (h *SimpananHandler) List(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	allowAll := role == "admin" || role == "super_admin"
	list, err := h.service.List(userID, allowAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *SimpananHandler) Detail(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	allowAll := role == "admin" || role == "super_admin"
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	item, err := h.service.Get(uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseError("not found"))
		return
	}
	if !allowAll && item.UserID != userID {
		c.JSON(http.StatusForbidden, utils.ResponseError("forbidden"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *SimpananHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	allowAll := role == "admin" || role == "super_admin"
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	var input struct {
		Type        string  `json:"type" binding:"required,oneof=wajib sukarela"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}
	payload := &model.Simpanan{Type: input.Type, Amount: input.Amount, Description: input.Description}
	updated, err := h.service.Update(uint(id64), userID, allowAll, payload)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, utils.ResponseError("forbidden"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (h *SimpananHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")
	allowAll := role == "admin" || role == "super_admin"
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("invalid id"))
		return
	}
	if err := h.service.Delete(uint(id64), userID, allowAll); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, utils.ResponseError("forbidden"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("deleted"))
}
