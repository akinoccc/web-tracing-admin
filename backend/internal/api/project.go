package api

import (
	"net/http"
	"strconv"

	"github.com/akinoccc/web-tracing-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary 创建项目
// @Description 创建新项目
// @Tags 项目
// @Accept json
// @Produce json
// @Param data body service.CreateProjectRequest true "项目信息"
// @Success 200 {object} model.Project "创建成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/projects [post]
func CreateProject(c *gin.Context) {
	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	projectService := service.ProjectService{}
	project, err := projectService.CreateProject(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// @Summary 获取项目列表
// @Description 获取当前用户的所有项目
// @Tags 项目
// @Produce json
// @Success 200 {array} model.Project "项目列表"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/projects [get]
func GetProjects(c *gin.Context) {
	// 获取当前用户 ID
	userID := c.GetUint("userID")

	projectService := service.ProjectService{}
	projects, err := projectService.GetUserProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// @Summary 获取项目详情
// @Description 获取项目详细信息
// @Tags 项目
// @Produce json
// @Param id path int true "项目ID"
// @Success 200 {object} model.Project "项目详情"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "项目不存在"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/projects/{id} [get]
func GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的项目ID"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	projectService := service.ProjectService{}
	project, err := projectService.GetProject(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// @Summary 更新项目
// @Description 更新项目信息
// @Tags 项目
// @Accept json
// @Produce json
// @Param id path int true "项目ID"
// @Param data body service.UpdateProjectRequest true "项目信息"
// @Success 200 {object} model.Project "更新成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "项目不存在"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/projects/{id} [put]
func UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的项目ID"})
		return
	}

	var req service.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	projectService := service.ProjectService{}
	project, err := projectService.UpdateProject(uint(id), &req, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// @Summary 删除项目
// @Description 删除项目
// @Tags 项目
// @Produce json
// @Param id path int true "项目ID"
// @Success 200 {object} SuccessResponse "删除成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "项目不存在"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/projects/{id} [delete]
func DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的项目ID"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	projectService := service.ProjectService{}
	err = projectService.DeleteProject(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "删除成功"})
}
