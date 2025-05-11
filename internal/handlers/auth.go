package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Thanhdat-debug/demo_login/internal/services"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AuthHandler xử lý các yêu cầu liên quan đến xác thực
type AuthHandler struct {
	authService services.AuthService
	logger      *logger.Logger
}

// NewAuthHandler tạo một instance mới của AuthHandler
func NewAuthHandler(authService services.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRequest chứa thông tin đăng ký từ client
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Register xử lý yêu cầu đăng ký
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gọi service để đăng ký
	userResponse, err := h.authService.Register(
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)

	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
			return
		}
		h.logger.Errorf("Register error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully", "user": userResponse})
}

// LoginRequest chứa thông tin đăng nhập từ client
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// Login xử lý yêu cầu đăng nhập
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gọi service để đăng nhập
	token, userResponse, err := h.authService.Login(req.UsernameOrEmail, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username/email or password"})
			return
		}
		h.logger.Errorf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user":    userResponse,
	})
}

// ValidateToken xử lý yêu cầu kiểm tra token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	// Kiểm tra format của Authorization header
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		return
	}

	// Lấy token string
	tokenString := parts[1]

	// Xác thực token
	token, err := h.authService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token is valid"})
}
