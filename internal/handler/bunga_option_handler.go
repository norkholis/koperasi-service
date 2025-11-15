package handler

import (
	"net/http"
	"strconv"

	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BungaOptionHandler struct {
	bungaOptionService service.BungaOptionService
}

func NewBungaOptionHandler(bungaOptionService service.BungaOptionService) *BungaOptionHandler {
	return &BungaOptionHandler{bungaOptionService: bungaOptionService}
}

type CreateBungaOptionRequest struct {
	Nama      string  `json:"nama" binding:"required"`
	Persen    float64 `json:"persen" binding:"required"`
	Deskripsi string  `json:"deskripsi"`
}

type UpdateBungaOptionRequest struct {
	Nama      string  `json:"nama" binding:"required"`
	Persen    float64 `json:"persen" binding:"required"`
	Deskripsi string  `json:"deskripsi"`
}

type SetActiveRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *BungaOptionHandler) Create(c *gin.Context) {
	var req CreateBungaOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	bungaOption, err := h.bungaOptionService.CreateBungaOption(userID.(uint), req.Nama, req.Persen, req.Deskripsi)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bunga option created successfully",
		"data":    bungaOption,
	})
}

func (h *BungaOptionHandler) List(c *gin.Context) {
	// Check if we want only active options
	activeOnly := c.Query("active") == "true"

	var bungaOptions interface{}
	var err error

	if activeOnly {
		bungaOptions, err = h.bungaOptionService.GetActiveBungaOptions()
	} else {
		bungaOptions, err = h.bungaOptionService.GetAllBungaOptions()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bunga options retrieved successfully",
		"data":    bungaOptions,
	})
}

func (h *BungaOptionHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	bungaOption, err := h.bungaOptionService.GetBungaOptionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ResponseError("Bunga option not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bunga option retrieved successfully",
		"data":    bungaOption,
	})
}

func (h *BungaOptionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	var req UpdateBungaOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	err = h.bungaOptionService.UpdateBungaOption(uint(id), userID.(uint), req.Nama, req.Persen, req.Deskripsi)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Bunga option updated successfully"))
}

func (h *BungaOptionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	err = h.bungaOptionService.DeleteBungaOption(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Bunga option deleted successfully"))
}

func (h *BungaOptionHandler) SetActive(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid ID"))
		return
	}

	var req SetActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	err = h.bungaOptionService.SetBungaOptionActive(uint(id), userID.(uint), req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	status := "deactivated"
	if req.IsActive {
		status = "activated"
	}

	c.JSON(http.StatusOK, utils.ResponseSuccess("Bunga option "+status+" successfully"))
}
