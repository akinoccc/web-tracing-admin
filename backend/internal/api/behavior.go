package api

import (
	"net/http"

	"github.com/akinoccc/web-tracing-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取页面访问数据
// @Description 获取项目的页面访问数据
// @Tags 用户行为
// @Produce json
// @Param projectId query int true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Success 200 {object} service.PVListResponse "页面访问数据列表"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/behavior/pv [get]
func GetPageViews(c *gin.Context) {
	projectID := c.Query("projectId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	eventService := service.EventService{}
	resp, err := eventService.GetPageViewList(projectID, page, pageSize, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 获取用户点击数据
// @Description 获取项目的用户点击数据
// @Tags 用户行为
// @Produce json
// @Param projectId query int true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Success 200 {object} service.ClickListResponse "用户点击数据列表"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/behavior/clicks [get]
func GetClicks(c *gin.Context) {
	projectID := c.Query("projectId")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	eventService := service.EventService{}
	resp, err := eventService.GetClickList(projectID, page, pageSize, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary 获取用户行为统计信息
// @Description 获取项目的用户行为统计信息
// @Tags 用户行为
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query int false "开始时间戳"
// @Param endTime query int false "结束时间戳"
// @Success 200 {object} service.BehaviorStatsResponse "用户行为统计信息"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/behavior/stats [get]
func GetBehaviorStats(c *gin.Context) {
	projectID := c.Query("projectId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	eventService := service.EventService{}
	resp, err := eventService.GetBehaviorStats(projectID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
