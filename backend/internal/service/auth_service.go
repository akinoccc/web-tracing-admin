package service

import (
	"errors"
	"time"

	"github.com/akinoccc/web-tracing-admin/internal/model"
	"github.com/dgrijalva/jwt-go"
)

// JWT 声明结构
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// 登录响应
type LoginResponse struct {
	Token    string      `json:"token"`
	User     *model.User `json:"user"`
	ExpireAt int64       `json:"expireAt"`
}

// 认证服务
type AuthService struct{}

// 生成 JWT token
func (s *AuthService) GenerateToken(user *model.User) (string, int64, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	expireAt := expireTime.Unix()

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    "web-tracing-admin",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(model.ServerSetting.JwtSecret))

	return token, expireAt, err
}

// 解析 JWT token
func (s *AuthService) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(model.ServerSetting.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}

// 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("密码错误")
	}

	token, expireAt, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    token,
		User:     user,
		ExpireAt: expireAt,
	}, nil
}

// 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := model.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 创建新用户
	user, err := model.CreateUser(req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
