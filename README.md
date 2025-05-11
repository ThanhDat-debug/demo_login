# Demo Login - Go API

Một API đơn giản được xây dựng bằng Go, Gin và GORM để xử lý xác thực người dùng và quản lý tài khoản.

## Cấu trúc dự án

```
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   ├── auth.go
│   │   └── user.go
│   ├── models/
│   │   └── user.go
│   ├── repository/
│   │   └── user_repository.go
│   ├── services/
│   │   ├── auth_service.go
│   │   └── user_service.go
│   └── middleware/
│       └── auth_middleware.go
├── pkg/
│   ├── database/
│   │   └── database.go
│   └── logger/
│       └── logger.go
├── go.mod
├── go.sum
├── README.md
└── .env

## Yêu cầu

- Go (version 1.18+)
- MySQL

## Cài đặt và Chạy

1. Clone repository

git clone https://github.com/Thanhdat-debug/demo_login.git
cd demo_login


2. Cài đặt các dependency

go mod download


3. Tạo cơ sở dữ liệu MySQL

```sql
CREATE DATABASE demo_login CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. Cập nhật file `.env` với thông tin cấu hình của bạn

DB_USER=root
DB_PASS=admin
DB_NAME=demo_login
DB_HOST=127.0.0.1
DB_PORT=3306
JWT_SECRET=mysecretkey


5. Build và chạy ứng dụng

```bash
go build -o ./bin/api ./cmd/api
./bin/api
```

Hoặc chạy trực tiếp:

```bash
go run ./cmd/api/main.go
```

Sau khi chạy, API sẽ khả dụng tại `http://localhost:8080`.

## API Endpoints

### Xác thực

- `POST /api/auth/register` - Đăng ký người dùng mới
- `POST /api/auth/login` - Đăng nhập và lấy token JWT
- `GET /api/auth/validate` - Kiểm tra token JWT

### Quản lý người dùng (cần xác thực)

- `GET /api/users/profile` - Lấy thông tin cá nhân
- `PUT /api/users/profile` - Cập nhật thông tin cá nhân
- `PUT /api/users/change-password` - Thay đổi mật khẩu
- `DELETE /api/users/account` - Xóa tài khoản

### Quản lý Admin (cần quyền admin)

- `GET /api/admin/users` - Lấy danh sách người dùng

## Ví dụ Request

### Đăng ký người dùng mới

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### Đăng nhập

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "testuser",
    "password": "password123"
  }'
```

### Lấy thông tin cá nhân (sử dụng token)

```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer your_token_here"
```

## Tài liệu tham khảo

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT Go](https://github.com/golang-jwt/jwt)