package model

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

// Project 项目模型
type Project struct {
	Model
	Name        string `json:"name" gorm:"size:100;not null"`
	AppKey      string `json:"appKey" gorm:"size:50;not null;unique"`
	Description string `json:"description" gorm:"type:text"`
	UserID      uint   `json:"userId"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
}

// 创建项目
func CreateProject(name, description string, userID uint) (*Project, error) {
	// 生成 AppKey
	h := md5.New()
	h.Write([]byte(name + time.Now().String()))
	appKey := hex.EncodeToString(h.Sum(nil))

	project := Project{
		Name:        name,
		AppKey:      appKey,
		Description: description,
		UserID:      userID,
	}

	if err := db.Create(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

// 获取用户的所有项目
func GetProjectsByUserID(userID uint) ([]Project, error) {
	var projects []Project
	if err := db.Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// 通过 ID 获取项目
func GetProjectByID(id uint) (*Project, error) {
	var project Project
	if err := db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// 通过 AppKey 获取项目
func GetProjectByAppKey(appKey string) (*Project, error) {
	var project Project
	if err := db.Where("app_key = ?", appKey).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// 更新项目
func UpdateProject(id uint, name, description string) (*Project, error) {
	project, err := GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	project.Name = name
	project.Description = description

	if err := db.Save(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

// 删除项目
func DeleteProject(id uint) error {
	return db.Delete(&Project{}, id).Error
}
