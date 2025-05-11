package repository

import (
	"errors"

	"github.com/Thanhdat-debug/demo_login/internal/models"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository định nghĩa interface cho các phương thức thao tác với User
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	List(page, size int) ([]models.User, int64, error)
}

// userRepository struct triển khai UserRepository interface
type userRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewUserRepository tạo một instance mới của UserRepository
func NewUserRepository(db *gorm.DB, logger *logger.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// Create tạo một user mới trong database
func (r *userRepository) Create(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		r.logger.Errorf("Error creating user: %v", err)
		return err
	}
	return nil
}

// FindByID tìm user theo ID
func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Không tìm thấy user, trả về nil, nil
		}
		r.logger.Errorf("Error finding user by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

// FindByUsername tìm user theo username
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorf("Error finding user by username: %v", err)
		return nil, err
	}
	return &user, nil
}

// FindByEmail tìm user theo email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorf("Error finding user by email: %v", err)
		return nil, err
	}
	return &user, nil
}

// Update cập nhật thông tin user
func (r *userRepository) Update(user *models.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		r.logger.Errorf("Error updating user: %v", err)
		return err
	}
	return nil
}

// Delete xóa user theo ID
func (r *userRepository) Delete(id uuid.UUID) error {
	err := r.db.Delete(&models.User{}, id).Error
	if err != nil {
		r.logger.Errorf("Error deleting user: %v", err)
		return err
	}
	return nil
}

// List lấy danh sách user với phân trang
func (r *userRepository) List(page, size int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Đếm tổng số user
	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		r.logger.Errorf("Error counting users: %v", err)
		return nil, 0, err
	}

	// Lấy danh sách user với phân trang
	offset := (page - 1) * size
	err = r.db.Offset(offset).Limit(size).Find(&users).Error
	if err != nil {
		r.logger.Errorf("Error listing users: %v", err)
		return nil, 0, err
	}

	return users, total, nil
}
