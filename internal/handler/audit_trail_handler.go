package handler

import (
	"net/http"
	"strconv"
	"time"

	"koperasi-service/internal/repository"
	"koperasi-service/internal/service"
	"koperasi-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuditTrailHandler struct {
	auditService       service.AuditTrailService
	transactionService service.TransactionHistoryService
}

func NewAuditTrailHandler(auditService service.AuditTrailService, transactionService service.TransactionHistoryService) *AuditTrailHandler {
	return &AuditTrailHandler{
		auditService:       auditService,
		transactionService: transactionService,
	}
}

// GetAuditTrails handles listing audit trails with filters
func (h *AuditTrailHandler) GetAuditTrails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	// Parse filters from query parameters
	filters := repository.AuditTrailFilters{}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uidUint := uint(uid)
			filters.UserID = &uidUint
		}
	}

	filters.Action = c.Query("action")
	filters.EntityTable = c.Query("entity_table")
	filters.IPAddress = c.Query("ip_address")

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = &endDate
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	} else {
		filters.Limit = 50 // Default limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	audits, total, err := h.auditService.GetAuditTrails(userID.(uint), filters)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Audit trails retrieved successfully",
		"data":    audits,
		"total":   total,
		"limit":   filters.Limit,
		"offset":  filters.Offset,
	})
}

// GetAuditTrailDetail handles getting specific audit trail by ID
func (h *AuditTrailHandler) GetAuditTrailDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid audit trail ID"))
		return
	}

	audit, err := h.auditService.GetAuditTrailByID(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Audit trail retrieved successfully",
		"data":    audit,
	})
}

// GetUserActivity handles getting audit trails for a specific user
func (h *AuditTrailHandler) GetUserActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid user ID"))
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default: last month
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		endDate = time.Now() // Default: now
	}

	audits, err := h.auditService.GetUserActivity(userID.(uint), uint(targetUserID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User activity retrieved successfully",
		"data":       audits,
		"user_id":    targetUserID,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"count":      len(audits),
	})
}

// GetSystemActivity handles getting system-wide audit trails
func (h *AuditTrailHandler) GetSystemActivity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -7) // Default: last week
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		endDate = time.Now() // Default: now
	}

	audits, err := h.auditService.GetSystemActivity(userID.(uint), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "System activity retrieved successfully",
		"data":       audits,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"count":      len(audits),
	})
}

// GetAuditSummary handles getting audit trail summary and analytics
func (h *AuditTrailHandler) GetAuditSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // Default: last month
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		endDate = time.Now() // Default: now
	}

	summary, err := h.auditService.GetAuditSummary(userID.(uint), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Audit summary retrieved successfully",
		"data":    summary,
	})
}

// GetTransactionHistory handles listing transaction history with filters
func (h *AuditTrailHandler) GetTransactionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	// Parse filters from query parameters
	filters := repository.TransactionHistoryFilters{}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uidUint := uint(uid)
			filters.UserID = &uidUint
		}
	}

	filters.TransactionType = c.Query("transaction_type")
	filters.Status = c.Query("status")

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = &endDate
		}
	}

	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil {
			filters.MinAmount = &minAmount
		}
	}

	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseFloat(maxAmountStr, 64); err == nil {
			filters.MaxAmount = &maxAmount
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	} else {
		filters.Limit = 50 // Default limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	transactions, total, err := h.transactionService.GetTransactionHistory(userID.(uint), filters)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction history retrieved successfully",
		"data":    transactions,
		"total":   total,
		"limit":   filters.Limit,
		"offset":  filters.Offset,
	})
}

// GetTransactionDetail handles getting specific transaction by ID
func (h *AuditTrailHandler) GetTransactionDetail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid transaction ID"))
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction retrieved successfully",
		"data":    transaction,
	})
}

// GetUserTransactions handles getting transactions for a specific user
func (h *AuditTrailHandler) GetUserTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid user ID"))
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -3, 0) // Default: last 3 months
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		endDate = time.Now() // Default: now
	}

	transactions, err := h.transactionService.GetUserTransactions(userID.(uint), uint(targetUserID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User transactions retrieved successfully",
		"data":       transactions,
		"user_id":    targetUserID,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"count":      len(transactions),
	})
}

// GetFinancialSummary handles getting financial summary and analytics
func (h *AuditTrailHandler) GetFinancialSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid start date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		startDate = time.Now().AddDate(-1, 0, 0) // Default: last year
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid end date format. Use YYYY-MM-DD"))
			return
		}
	} else {
		endDate = time.Now() // Default: now
	}

	summary, err := h.transactionService.GetFinancialSummary(userID.(uint), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Financial summary retrieved successfully",
		"data":    summary,
	})
}

// GenerateFinancialReport handles generating comprehensive financial reports
func (h *AuditTrailHandler) GenerateFinancialReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not authenticated"))
		return
	}

	type ReportRequest struct {
		ReportType string `json:"report_type" binding:"required"`
		StartDate  string `json:"start_date" binding:"required"`
		EndDate    string `json:"end_date" binding:"required"`
	}

	var req ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid start date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid end date format. Use YYYY-MM-DD"))
		return
	}

	report, err := h.transactionService.GenerateFinancialReport(userID.(uint), req.ReportType, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusForbidden, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Financial report generated successfully",
		"data":    report,
	})
}
