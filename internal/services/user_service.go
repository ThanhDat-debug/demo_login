package services

import (
	"errors"

	"github.com/Thanhdat-debug/demo_login/internal/models"
	"github.com/Thanhdat-debug/demo_login/internal/repository"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/google/uuid"
)

// Định nghĩa các lỗi
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserService định nghĩa interface cho các phương thức quản lý user
type UserService interface {
	GetUserByID(id uuid.UUID) (*models.UserResponse, error)
	UpdateUser(id uuid.UUID, firstName, lastName string) (*models.UserResponse, error)
	ChangePassword(id uuid.UUID, oldPassword, newPassword string) error
	DeleteUser(id uuid.UUID) error
	ListUsers(page, size int) ([]models.UserResponse, int64, error)
}

// userService struct triển khai UserService interface
type userService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewUserService tạo một instance mới của UserService
func NewUserService(userRepo repository.UserRepository, logger *logger.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUserByID lấy thông tin user theo ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	userResponse := user.ToUserResponse()
	return &userResponse, nil
}

// UpdateUser cập nhật thông tin user
func (s *userService) UpdateUser(id uuid.UUID, firstName, lastName string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Cập nhật thông tin
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}

	// Lưu vào database
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	userResponse := user.ToUserResponse()
	return &userResponse, nil
}

// ChangePassword thay đổi mật khẩu của user
func (s *userService) ChangePassword(id uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Kiểm tra mật khẩu cũ
	if !user.CheckPassword(oldPassword) {
		return ErrInvalidCredentials
	}

	// Cập nhật mật khẩu mới
	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return err
	}

	// Lưu vào database
	return s.userRepo.Update(user)
}

// DeleteUser xóa user theo ID
func (s *userService) DeleteUser(id uuid.UUID) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(id)
}

// ListUsers lấy danh sách users với phân trang
func (s *userService) ListUsers(page, size int) ([]models.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	users, total, err := s.userRepo.List(page, size)
	if err != nil {
		return nil, 0, err
	}

	// Chuyển đổi sang UserResponse
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToUserResponse()
	}

	return userResponses, total, nil
}
