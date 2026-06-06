package handler

import (
	"fmt"
	"strconv"
	"wallet/internal/middleware"
	"wallet/internal/service"
	"wallet/internal/store"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionServiceInterface
}
type TransactionRequest struct {
	ReceiverID uint64  `json:"receiver_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
}


// TransactionHandler provides HTTP handlers for transaction-related operations, such as creating transactions,
//  retrieving transaction details by ID, and listing transactions for a specific user with pagination support.
//  It interacts with the TransactionService to perform these operations and returns appropriate HTTP responses based on the service results.

// CreateTransaction(ctx context.Context, transaction *store.Transaction) (*store.Transaction, error)
// 	GetTransactionByID(ctx context.Context, id uint64) (*store.Transaction, error)
// 	GetTransactionsByUserID(ctx context.Context, userID uint64, page int, pageSize int) ([]store.Transaction, error)
// TransactionHandler provides HTTP handlers for transaction-related operations, such as creating transactions,
//  retrieving transaction details by ID, and listing transactions for a specific user with pagination support.
//  It interacts with the TransactionService to perform these operations and returns appropriate HTTP responses based on the service results.

func RegisterTransactionRoutes(router *gin.RouterGroup, transactionService service.TransactionServiceInterface) {
	handler := &TransactionHandler{transactionService: transactionService}
	router.POST("/", middleware.AuthMiddleware(), handler.CreateTransaction)
	router.GET("/:id", middleware.AuthMiddleware(), handler.GetTransactionByID)
	router.GET("/user/:userID", middleware.AuthMiddleware(), handler.GetTransactionsByUserID)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Implementation for creating a transaction
	ctx := c.Request.Context()
	var transactionRequest TransactionRequest
	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if transactionRequest.Currency == "" {
		transactionRequest.Currency = "USD"
	}

	transaction := &store.Transaction{
		SenderID:   userID.(uint64),
		Currency:   transactionRequest.Currency,
		ReceiverID: transactionRequest.ReceiverID,
		Amount:     transactionRequest.Amount,
	}

	result, err := h.transactionService.CreateTransaction(ctx, transaction)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)

}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	// Implementation for retrieving a transaction by ID
	ctx := c.Request.Context()
	idParam := c.Param("id")
	var id uint64
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"transaction": transaction,
	})
}

func (h *TransactionHandler) GetTransactionsByUserID(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	phone := c.Query("phone")
	email := c.Query("email")
	username := c.Query("username")

	ctx := c.Request.Context()
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		c.JSON(400, gin.H{"error": "Invalid page size"})
		return
	}

	// Implementation for listing transactions for a specific user with pagination support
	userIDParam := c.Param("userID")
	var userID uint64
	_, err = fmt.Sscanf(userIDParam, "%d", &userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	transactions, err := h.transactionService.GetTransactionsByUserID(ctx, userID, pageInt, pageSizeInt, phone, email, username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"transactions": transactions,
		"page":         pageInt,
		"pageSize":     pageSizeInt,
	})
}
