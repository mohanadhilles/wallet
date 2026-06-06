package handler

import (
	"wallet/internal/middleware"
	"wallet/internal/service"
	"github.com/gin-gonic/gin"
)



type WalletHandler struct {
	service service.WalletServiceInterface
}

func ExportWalletHandler(s service.WalletServiceInterface) *WalletHandler	{
	return &WalletHandler{service: s}
}

func RegisterWalletRoutes(router *gin.RouterGroup, walletService service.WalletServiceInterface) {
	handler := ExportWalletHandler(walletService)
	router.GET("/me", middleware.AuthMiddleware(), handler.GetMyWallet)
}

func (h *WalletHandler) GetMyWallet(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	wallet, err := h.service.GetWalletByUserID(ctx, userID.(uint64))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"wallet": wallet,
	})
}