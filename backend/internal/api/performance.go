package api

import (
	"net/http"

	"github.com/akinoccc/web-tracing-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取性能数据
// @Description 获取项目的性能数据
// @Tags 性能监控
// @Produce json
// @Param projectId query int true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Param type query string false "性能类型" Enums(page, resource)
// @Success 200 {object} service.PerformanceListResponse "性能数据列表"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/performance [get]
func GetPerformance(c *gin.Context) {
	projectID := c.Query("projectId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	perfType := c.Query("type")

	eventService := service.EventService{}
	resp, err := eventService.GetPerformanceList(projectID, page, pageSize, startTime, endTime, perfType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 获取性能统计信息
// @Description 获取项目的性能统计信息
// @Tags 性能监控
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Success 200 {object} service.PerformanceStatsResponse "性能统计信息"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/performance/stats [get]
func GetPerformanceStats(c *gin.Context) {
	projectID := c.Query("projectId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	eventService := service.EventService{}
	resp, err := eventService.GetPerformanceStats(projectID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 获取资源性能数据
// @Description 获取项目的资源性能数据
// @Tags 性能监控
// @Produce json
// @Param projectId query int true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Param resourceType query string false "资源类型"
// @Success 200 {object} service.ResourcePerformanceListResponse "资源性能数据列表"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/performance/resources [get]
func GetResourcePerformance(c *gin.Context) {
	projectID := c.Query("projectId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	resourceType := c.Query("resourceType")

	eventService := service.EventService{}
	resp, err := eventService.GetResourcePerformanceList(projectID, page, pageSize, startTime, endTime, resourceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
