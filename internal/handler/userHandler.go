package handler

import (
	"net/http"
	"wallet/internal/middleware"
	"wallet/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

// UserHandler provides HTTP handlers for user-related operations, such as generating OTPs, verifying OTPs, and retrieving user information.
//  It interacts with the UserService to perform these operations and returns appropriate HTTP responses

// UserHandler provides HTTP handlers for user-related operations, such as generating OTPs, verifying OTPs, and retrieving user information.
//  It interacts with the UserService to perform these operations and returns appropriate HTTP responses based on the service results.

func RegisterUserRoutes(router *gin.RouterGroup, userService service.UserServiceInterface) {
	handler := &UserHandler{userService: userService}
	router.POST("/otp", handler.RequestOTP)
	router.POST("/otp/verify", handler.VerifyOTP)
	router.GET("/me", middleware.AuthMiddleware(), handler.Me)
}


func (h *UserHandler) Me(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := h.userService.Me(ctx, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) RequestOTP(c *gin.Context) {
	ctx := c.Request.Context()
	var otpRequest struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&otpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	_, err := h.userService.GenerateOTP(ctx, otpRequest.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP generated successfully"})
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	ctx := c.Request.Context()
	var confirmOTPRequest struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&confirmOTPRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.userService.GetUserByUsername(ctx, confirmOTPRequest.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	userID := user.ID
	token, err := h.userService.VerifyOTP(ctx, userID, confirmOTPRequest.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
