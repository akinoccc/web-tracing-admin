package service

import (
	"errors"

	"github.com/akinoccc/web-tracing-admin/internal/model"
)

// 项目创建请求
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// 项目更新请求
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// 项目服务
type ProjectService struct{}

// 创建项目
func (s *ProjectService) CreateProject(req *CreateProjectRequest, userID uint) (*model.Project, error) {
	project, err := model.CreateProject(req.Name, req.Description, userID)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// 获取用户的所有项目
func (s *ProjectService) GetUserProjects(userID uint) ([]model.Project, error) {
	projects, err := model.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// 获取项目详情
func (s *ProjectService) GetProject(id uint, userID uint) (*model.Project, error) {
	project, err := model.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	// 检查项目是否属于该用户
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	return project, nil
}

// 更新项目
func (s *ProjectService) UpdateProject(id uint, req *UpdateProjectRequest, userID uint) (*model.Project, error) {
	project, err := model.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	// 检查项目是否属于该用户
	if project.UserID != userID {
		return nil, errors.New("无权修改该项目")
	}

	project, err = model.UpdateProject(id, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// 删除项目
func (s *ProjectService) DeleteProject(id uint, userID uint) error {
	project, err := model.GetProjectByID(id)
	if err != nil {
		return err
	}

	// 检查项目是否属于该用户
	if project.UserID != userID {
		return errors.New("无权删除该项目")
	}

	return model.DeleteProject(id)
}
