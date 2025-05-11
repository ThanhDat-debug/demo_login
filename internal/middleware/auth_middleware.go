package middleware

import (
	"net/http"
	"strings"

	"github.com/Thanhdat-debug/demo_login/internal/services"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware chứa các middleware liên quan đến xác thực
type AuthMiddleware struct {
	authService services.AuthService
	logger      *logger.Logger
}

// NewAuthMiddleware tạo một instance mới của AuthMiddleware
func NewAuthMiddleware(authService services.AuthService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// JWTAuthMiddleware kiểm tra và xác thực JWT token
func (m *AuthMiddleware) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Kiểm tra format của Authorization header
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Lấy token string
		tokenString := parts[1]

		// Xác thực token
		token, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Kiểm tra claims
		claims, ok := token.Claims.(*services.Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Parse userID từ claims
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			c.Abort()
			return
		}

		// Lưu thông tin vào context để các handler có thể sử dụng
		c.Set("userID", userID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// AdminRequired kiểm tra nếu user có role 'admin'
func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
