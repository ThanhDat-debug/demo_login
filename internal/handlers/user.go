package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Thanhdat-debug/demo_login/internal/services"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler xử lý các yêu cầu liên quan đến user
type UserHandler struct {
	userService services.UserService
	logger      *logger.Logger
}

// NewUserHandler tạo một instance mới của UserHandler
func NewUserHandler(userService services.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile xử lý yêu cầu lấy thông tin cá nhân của user đã đăng nhập
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Lấy userID từ context (đã được set bởi JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Gọi service để lấy thông tin user
	userResponse, err := h.userService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Errorf("GetProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
}

// UpdateProfileRequest chứa thông tin cập nhật hồ sơ từ client
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateProfile xử lý yêu cầu cập nhật thông tin cá nhân
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Lấy userID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse request body
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gọi service để cập nhật user
	userResponse, err := h.userService.UpdateUser(userID.(uuid.UUID), req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Errorf("UpdateProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully", "user": userResponse})
}

// ChangePasswordRequest chứa thông tin đổi mật khẩu từ client
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword xử lý yêu cầu thay đổi mật khẩu
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Lấy userID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse request body
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gọi service để thay đổi mật khẩu
	err := h.userService.ChangePassword(userID.(uuid.UUID), req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect old password"})
			return
		}
		h.logger.Errorf("ChangePassword error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// DeleteAccount xử lý yêu cầu xóa tài khoản
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	// Lấy userID từ context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Gọi service để xóa tài khoản
	err := h.userService.DeleteUser(userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Errorf("DeleteAccount error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}

// GetUsersList xử lý yêu cầu lấy danh sách user (admin only)
func (h *UserHandler) GetUsersList(c *gin.Context) {
	// Xử lý tham số phân trang
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// Gọi service để lấy danh sách user
	users, total, err := h.userService.ListUsers(page, size)
	if err != nil {
		h.logger.Errorf("GetUsersList error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":      users,
		"total":      total,
		"page":       page,
		"size":       size,
		"total_page": (total + int64(size) - 1) / int64(size),
	})
}
