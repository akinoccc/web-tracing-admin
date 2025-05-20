package api

import (
	"net/http"

	"github.com/akinoccc/web-tracing-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary 接收SDK上报数据
// @Description 接收SDK上报的错误和性能数据
// @Tags 数据上报
// @Accept json
// @Produce json
// @Param data body service.TrackRequest true "上报数据"
// @Success 200 {object} SuccessResponse "上报成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /trackweb [post]
func TrackWeb(c *gin.Context) {
	var req service.TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	eventService := service.EventService{}
	err := eventService.ProcessTrackData(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "上报成功"})
}

// @Summary 获取错误列表
// @Description 获取项目的错误列表
// @Tags 错误监控
// @Produce json
// @Param projectId query int true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Param errorType query string false "错误类型"
// @Param severity query string false "严重程度"
// @Success 200 {object} service.ErrorListResponse "错误列表"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/errors [get]
func GetErrors(c *gin.Context) {
	projectID := c.Query("projectId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	errorType := c.Query("errorType")
	severity := c.Query("severity")

	eventService := service.EventService{}
	resp, err := eventService.GetErrorList(projectID, page, pageSize, startTime, endTime, errorType, severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 获取错误详情
// @Description 获取错误详细信息
// @Tags 错误监控
// @Produce json
// @Param id path int true "错误ID"
// @Success 200 {object} service.ErrorDetailResponse "错误详情"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "错误不存在"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/errors/{id} [get]
func GetErrorDetail(c *gin.Context) {
	id := c.Param("id")

	eventService := service.EventService{}
	resp, err := eventService.GetErrorDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 获取错误统计信息
// @Description 获取项目的错误统计信息
// @Tags 错误监控
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Success 200 {object} service.ErrorStatsResponse "错误统计信息"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/errors/stats [get]
func GetErrorStats(c *gin.Context) {
	projectID := c.Query("projectId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	eventService := service.EventService{}
	resp, err := eventService.GetErrorStats(projectID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
