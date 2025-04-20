package model

import (
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	Model
	Username string `json:"username" gorm:"size:50;not null;unique"`
	Password string `json:"-" gorm:"size:255;not null"`
	Email    string `json:"email" gorm:"size:100;not null;unique"`
}

// 创建用户
func CreateUser(username, password, email string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
