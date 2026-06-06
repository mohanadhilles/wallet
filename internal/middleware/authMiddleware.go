package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"
	"wallet/provider/jwt"

	"github.com/gin-gonic/gin"
)

/**
AuthMiddleware is a Gin middleware function that validates JWT tokens in the Authorization header of incoming requests. 
It checks for the presence of the Authorization header, extracts the token, and uses the JWTService to validate it. 
If the token is valid, it sets the userID in the context for downstream handlers to use.
 If the token is invalid or missing, it returns a 401 Unauthorized response and aborts the request.
*/
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		jwtService := jwt.JWTServices(os.Getenv("JWT_SECRET_KEY"), 24*time.Hour)
		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
