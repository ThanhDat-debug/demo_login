package services

import (
	"errors"
	"time"

	"github.com/Thanhdat-debug/demo_login/internal/config"
	"github.com/Thanhdat-debug/demo_login/internal/models"
	"github.com/Thanhdat-debug/demo_login/internal/repository"
	"github.com/Thanhdat-debug/demo_login/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

// Định nghĩa các lỗi
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// AuthService định nghĩa interface cho các phương thức xác thực
type AuthService interface {
	Register(username, email, password, firstName, lastName string) (*models.UserResponse, error)
	Login(usernameOrEmail, password string) (string, *models.UserResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// authService struct triển khai AuthService interface
type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
	logger   *logger.Logger
}

// NewAuthService tạo một instance mới của AuthService
func NewAuthService(userRepo repository.UserRepository, config *config.Config, logger *logger.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
		logger:   logger,
	}
}

// Claims là custom claims cho JWT
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register đăng ký user mới
func (s *authService) Register(username, email, password, firstName, lastName string) (*models.UserResponse, error) {
	// Kiểm tra username đã tồn tại chưa
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Kiểm tra email đã tồn tại chưa
	existingUser, err = s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Tạo user mới
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user", // Mặc định là "user"
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// Lưu user vào database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Trả về thông tin user đã đăng ký (không bao gồm password)
	userResponse := user.ToUserResponse()
	return &userResponse, nil
}

// Login xác thực người dùng và tạo JWT token
func (s *authService) Login(usernameOrEmail, password string) (string, *models.UserResponse, error) {
	var user *models.User
	var err error

	// Kiểm tra nếu đầu vào là email
	user, err = s.userRepo.FindByEmail(usernameOrEmail)
	if err != nil {
		return "", nil, err
	}

	// Nếu không tìm thấy bằng email, thử tìm bằng username
	if user == nil {
		user, err = s.userRepo.FindByUsername(usernameOrEmail)
		if err != nil {
			return "", nil, err
		}
	}

	// Kiểm tra nếu user không tồn tại
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	// Kiểm tra password
	if !user.CheckPassword(password) {
		return "", nil, ErrInvalidCredentials
	}

	// Tạo JWT token
	claims := Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token hết hạn sau 24h
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		s.logger.Errorf("Error signing token: %v", err)
		return "", nil, err
	}

	// Trả về token và thông tin user
	userResponse := user.ToUserResponse()
	return tokenString, &userResponse, nil
}

// ValidateToken kiểm tra JWT token có hợp lệ không
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán ký
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
