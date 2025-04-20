package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/akinoccc/web-tracing-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取事件列表
// @Description 获取项目的事件列表
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param eventType query string false "事件类型"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} service.EventListResponse "事件列表"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events [get]
func GetEvents(c *gin.Context) {
	var req service.EventQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	events, err := eventService.GetEvents(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary 获取事件详情
// @Description 获取事件详细信息
// @Tags 事件
// @Produce json
// @Param id path int true "事件ID"
// @Success 200 {object} service.EventDetailResponse "事件详情"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "事件不存在"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/{id} [get]
func GetEventDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的事件ID"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	event, err := eventService.GetEventDetail(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// @Summary 获取错误事件统计
// @Description 获取项目的错误事件统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.EventStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats [get]
func GetErrorStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetErrorStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 获取错误事件按浏览器分布统计
// @Description 获取项目的错误事件按浏览器分布统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.DistributionStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats/browser [get]
func GetErrorByBrowserStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetErrorByBrowserStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 获取错误事件按操作系统分布统计
// @Description 获取项目的错误事件按操作系统分布统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.DistributionStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats/os [get]
func GetErrorByOSStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetErrorByOSStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 获取错误事件按设备类型分布统计
// @Description 获取项目的错误事件按设备类型分布统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.DistributionStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats/device [get]
func GetErrorByDeviceStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetErrorByDeviceStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 获取错误事件按错误类型分布统计
// @Description 获取项目的错误事件按错误类型分布统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.DistributionStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats/error-type [get]
func GetErrorByTypeStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetErrorByTypeStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 获取性能指标统计
// @Description 获取项目的性能指标统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.PerformanceStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats/performance [get]
func GetPerformanceStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetPerformanceStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 获取请求错误统计
// @Description 获取项目的请求错误统计
// @Tags 事件
// @Produce json
// @Param projectId query int true "项目ID"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} service.RequestErrorStatsResponse "统计数据"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Security ApiKeyAuth
// @Router /api/events/stats/request-error [get]
func GetRequestErrorStats(c *gin.Context) {
	var req service.EventStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的请求参数"})
		return
	}

	// 获取当前用户 ID
	userID := c.GetUint("userID")

	eventService := service.EventService{}
	stats, err := eventService.GetRequestErrorStats(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary 上报跟踪数据
// @Description SDK 上报跟踪数据接口
// @Tags 跟踪
// @Accept json
// @Produce json
// @Param app_key query string true "应用密钥"
// @Param data body service.TrackWebRequest true "跟踪数据"
// @Success 200 {object} SuccessResponse "上报成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /trackweb [post]
func TrackWeb(c *gin.Context) {
	// appKey := c.Query("app_key")
	// if appKey == "" {
	// 	c.JSON(http.StatusBadRequest, ErrorResponse{Message: "缺少 app_key 参数"})
	// 	return
	// }

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "读取请求体失败"})
		return
	}

	// 解析请求数据
	var req service.TrackWebRequest
	if err := json.Unmarshal(body, &req); err != nil {
		// 尝试兼容 sendbeacon 的传输数据格式
		if len(body) > 0 {
			// 如果请求体不为空但解析失败，可能是其他格式
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "解析请求数据失败"})
			return
		}
		// 如果请求体为空，可能是 GET 请求或其他情况
		req = service.TrackWebRequest{}
	}

	eventService := service.EventService{}
	err = eventService.ProcessTrackData(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "上报成功"})
}

// @Summary 获取所有跟踪数据
// @Description 获取所有跟踪数据
// @Tags 跟踪
// @Produce json
// @Param eventType query string false "事件类型"
// @Param startTime query string false "开始时间，格式：2006-01-02"
// @Param endTime query string false "结束时间，格式：2006-01-02"
// @Param keyword query string false "关键字搜索"
// @Param browser query string false "浏览器类型"
// @Param os query string false "操作系统"
// @Param deviceType query string false "设备类型"
// @Param errorType query string false "错误类型，仅当 eventType=error 时有效"
// @Param responseStatus query string false "响应状态码，仅当 eventType=request 时有效"
// @Success 200 {object} SuccessResponse "获取成功"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /getAllTracingList [get]
func GetAllTracingList(c *gin.Context) {
	// 获取查询参数
	eventType := c.Query("eventType")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	keyword := c.Query("keyword")
	browser := c.Query("browser")
	os := c.Query("os")
	deviceType := c.Query("deviceType")
	errorType := c.Query("errorType")
	responseStatus := c.Query("responseStatus")

	// 构建查询参数映射
	params := make(map[string]interface{})
	if startTime != "" {
		params["startTime"] = startTime
	}
	if endTime != "" {
		params["endTime"] = endTime
	}
	if keyword != "" {
		params["keyword"] = keyword
	}
	if browser != "" {
		params["browser"] = browser
	}
	if os != "" {
		params["os"] = os
	}
	if deviceType != "" {
		params["deviceType"] = deviceType
	}
	if errorType != "" {
		params["errorType"] = errorType
	}
	if responseStatus != "" {
		params["responseStatus"] = responseStatus
	}

	eventService := service.EventService{}
	events, err := eventService.GetAllEvents(eventType, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": events,
	})
}

// @Summary 清除所有跟踪数据
// @Description 清除所有跟踪数据
// @Tags 跟踪
// @Produce json
// @Success 200 {object} SuccessResponse "清除成功"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /cleanTracingList [post]
func CleanTracingList(c *gin.Context) {
	eventService := service.EventService{}
	err := eventService.CleanAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "清除成功！",
	})
}

// @Summary 获取基础信息
// @Description 获取基础信息
// @Tags 跟踪
// @Produce json
// @Success 200 {object} SuccessResponse "获取成功"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /getBaseInfo [get]
func GetBaseInfo(c *gin.Context) {
	eventService := service.EventService{}
	baseInfo, err := eventService.GetBaseInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": baseInfo,
	})
}
