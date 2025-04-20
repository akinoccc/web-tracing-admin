package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/akinoccc/web-tracing-admin/internal/model"
)

// 事件查询请求
type EventQueryRequest struct {
	ProjectID      uint   `form:"projectId" binding:"required"`
	EventType      string `form:"eventType"`
	StartTime      string `form:"startTime"`      // 开始时间，格式：2006-01-02
	EndTime        string `form:"endTime"`        // 结束时间，格式：2006-01-02
	Keyword        string `form:"keyword"`        // 关键字搜索，搜索标题、URL、错误信息等
	Browser        string `form:"browser"`        // 浏览器类型
	OS             string `form:"os"`             // 操作系统
	DeviceType     string `form:"deviceType"`     // 设备类型
	ErrorType      string `form:"errorType"`      // 错误类型，仅当 eventType=error 时有效
	ResponseStatus string `form:"responseStatus"` // 响应状态码，仅当 eventType=request 时有效
	Page           int    `form:"page,default=1"`
	PageSize       int    `form:"pageSize,default=10"`
}

// 事件统计请求
type EventStatsRequest struct {
	ProjectID uint   `form:"projectId" binding:"required"`
	StartTime string `form:"startTime" binding:"required"`
	EndTime   string `form:"endTime" binding:"required"`
}

// 事件列表响应
type EventListResponse struct {
	Total int64         `json:"total"`
	List  []model.Event `json:"list"`
}

// 事件详情响应
type EventDetailResponse struct {
	Event       *model.Event            `json:"event"`
	ErrorEvent  *model.ErrorEvent       `json:"errorEvent,omitempty"`
	Performance *model.PerformanceEvent `json:"performance,omitempty"`
	Request     *model.RequestEvent     `json:"request,omitempty"`
	Click       *model.ClickEvent       `json:"click,omitempty"`
	Route       *model.RouteEvent       `json:"route,omitempty"`
	Exposure    *model.ExposureEvent    `json:"exposure,omitempty"`
}

// 事件统计响应
type EventStatsResponse struct {
	Dates  []string `json:"dates"`
	Counts []int64  `json:"counts"`
}

// 分布统计响应
type DistributionStatsResponse struct {
	Labels []string `json:"labels"`
	Values []int64  `json:"values"`
}

// 性能统计响应
type PerformanceStatsResponse struct {
	Labels []string  `json:"labels"`
	Values []float64 `json:"values"`
}

// 请求错误统计响应
type RequestErrorStatsResponse struct {
	StatusCodes []int   `json:"statusCodes"`
	Counts      []int64 `json:"counts"`
}

// 跟踪数据请求
type TrackWebRequest struct {
	EventInfo []map[string]interface{} `json:"eventInfo"`
	BaseInfo  map[string]interface{}   `json:"baseInfo"`
}

// 事件服务
type EventService struct{}

// 获取事件列表
func (s *EventService) GetEvents(req *EventQueryRequest, userID uint) (*EventListResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 构建查询参数
	params := make(map[string]interface{})

	// 添加时间范围参数
	if req.StartTime != "" {
		params["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		params["endTime"] = req.EndTime
	}

	// 添加关键字搜索参数
	if req.Keyword != "" {
		params["keyword"] = req.Keyword
	}

	// 添加浏览器、操作系统、设备类型参数
	if req.Browser != "" {
		params["browser"] = req.Browser
	}
	if req.OS != "" {
		params["os"] = req.OS
	}
	if req.DeviceType != "" {
		params["deviceType"] = req.DeviceType
	}

	// 添加错误类型和响应状态码参数
	if req.ErrorType != "" {
		params["errorType"] = req.ErrorType
	}
	if req.ResponseStatus != "" {
		params["responseStatus"] = req.ResponseStatus
	}

	events, total, err := model.GetEventsByProjectID(req.ProjectID, req.EventType, req.Page, req.PageSize, params)
	if err != nil {
		return nil, err
	}

	return &EventListResponse{
		Total: total,
		List:  events,
	}, nil
}

// 获取事件详情
func (s *EventService) GetEventDetail(id uint, userID uint) (*EventDetailResponse, error) {
	event, err := model.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(event.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该事件")
	}

	response := &EventDetailResponse{
		Event: event,
	}

	// 根据事件类型获取详细信息
	switch event.EventType {
	case "error":
		errorEvent, err := model.GetErrorEventByEventID(event.ID)
		if err == nil {
			response.ErrorEvent = errorEvent
		}
	case "performance":
		performanceEvent, err := model.GetPerformanceEventByEventID(event.ID)
		if err == nil {
			response.Performance = performanceEvent
		}
	case "click":
		clickEvent, err := model.GetClickEventByEventID(event.ID)
		if err == nil {
			response.Click = clickEvent
		}
	case "pv":
		routeEvent, err := model.GetRouteEventByEventID(event.ID)
		if err == nil {
			response.Route = routeEvent
		}
	case "intersection":
		exposureEvent, err := model.GetExposureEventByEventID(event.ID)
		if err == nil {
			response.Exposure = exposureEvent
		}
	}

	return response, nil
}

// 获取错误事件统计
func (s *EventService) GetErrorStats(req *EventStatsRequest, userID uint) (*EventStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetErrorEventStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &EventStatsResponse{
		Dates:  make([]string, 0),
		Counts: make([]int64, 0),
	}

	// 填充日期范围内的所有日期
	for d := startTime; !d.After(endTime); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		response.Dates = append(response.Dates, dateStr)
		count, exists := stats[dateStr]
		if exists {
			response.Counts = append(response.Counts, count)
		} else {
			response.Counts = append(response.Counts, 0)
		}
	}

	return response, nil
}

// 获取错误事件按浏览器分布统计
func (s *EventService) GetErrorByBrowserStats(req *EventStatsRequest, userID uint) (*DistributionStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetErrorEventByBrowserStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &DistributionStatsResponse{
		Labels: make([]string, 0),
		Values: make([]int64, 0),
	}

	// 将映射转换为切片
	for browser, count := range stats {
		response.Labels = append(response.Labels, browser)
		response.Values = append(response.Values, count)
	}

	return response, nil
}

// 获取错误事件按操作系统分布统计
func (s *EventService) GetErrorByOSStats(req *EventStatsRequest, userID uint) (*DistributionStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetErrorEventByOSStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &DistributionStatsResponse{
		Labels: make([]string, 0),
		Values: make([]int64, 0),
	}

	// 将映射转换为切片
	for os, count := range stats {
		response.Labels = append(response.Labels, os)
		response.Values = append(response.Values, count)
	}

	return response, nil
}

// 获取错误事件按设备类型分布统计
func (s *EventService) GetErrorByDeviceStats(req *EventStatsRequest, userID uint) (*DistributionStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetErrorEventByDeviceStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &DistributionStatsResponse{
		Labels: make([]string, 0),
		Values: make([]int64, 0),
	}

	// 将映射转换为切片
	for deviceType, count := range stats {
		response.Labels = append(response.Labels, deviceType)
		response.Values = append(response.Values, count)
	}

	return response, nil
}

// 获取错误事件按错误类型分布统计
func (s *EventService) GetErrorByTypeStats(req *EventStatsRequest, userID uint) (*DistributionStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetErrorEventByTypeStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &DistributionStatsResponse{
		Labels: make([]string, 0),
		Values: make([]int64, 0),
	}

	// 将映射转换为切片
	for errorType, count := range stats {
		response.Labels = append(response.Labels, errorType)
		response.Values = append(response.Values, count)
	}

	return response, nil
}

// 获取性能指标统计
func (s *EventService) GetPerformanceStats(req *EventStatsRequest, userID uint) (*PerformanceStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetPerformanceStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &PerformanceStatsResponse{
		Labels: make([]string, 0),
		Values: make([]float64, 0),
	}

	// 将映射转换为切片
	for resourceType, avgDuration := range stats {
		response.Labels = append(response.Labels, resourceType)
		response.Values = append(response.Values, avgDuration)
	}

	return response, nil
}

// 获取请求错误统计
func (s *EventService) GetRequestErrorStats(req *EventStatsRequest, userID uint) (*RequestErrorStatsResponse, error) {
	// 检查项目是否属于该用户
	project, err := model.GetProjectByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问该项目")
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02", req.StartTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	endTime, err := time.Parse("2006-01-02", req.EndTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	// 将结束时间设置为当天的最后一秒
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 获取统计数据
	stats, err := model.GetRequestErrorStats(req.ProjectID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &RequestErrorStatsResponse{
		StatusCodes: make([]int, 0),
		Counts:      make([]int64, 0),
	}

	// 将映射转换为切片
	for statusCode, count := range stats {
		response.StatusCodes = append(response.StatusCodes, statusCode)
		response.Counts = append(response.Counts, count)
	}

	return response, nil
}

// 处理跟踪数据
func (s *EventService) ProcessTrackData(req *TrackWebRequest) error {
	// 通过 AppKey 获取项目
	project, err := model.GetProjectByAppKey(req.BaseInfo["appCode"].(string))
	if err != nil {
		return errors.New("无效的 AppKey")
	}

	// 处理基础信息
	var baseInfo *model.BaseInfo
	if req.BaseInfo != nil {
		baseInfo = &model.BaseInfo{
			ProjectID: project.ID,
		}

		// 填充基础信息字段
		if userID, ok := req.BaseInfo["userId"].(string); ok {
			baseInfo.UserID = userID
		}
		if userAgent, ok := req.BaseInfo["userAgent"].(string); ok {
			baseInfo.UserAgent = userAgent
		}
		if ip, ok := req.BaseInfo["ip"].(string); ok {
			baseInfo.IP = ip
		}
		if browser, ok := req.BaseInfo["browser"].(string); ok {
			baseInfo.Browser = browser
		}
		if browserVersion, ok := req.BaseInfo["browserVersion"].(string); ok {
			baseInfo.BrowserVersion = browserVersion
		}
		if os, ok := req.BaseInfo["os"].(string); ok {
			baseInfo.OS = os
		}
		if osVersion, ok := req.BaseInfo["osVersion"].(string); ok {
			baseInfo.OSVersion = osVersion
		}
		if deviceType, ok := req.BaseInfo["deviceType"].(string); ok {
			baseInfo.DeviceType = deviceType
		}
		if screenWidth, ok := req.BaseInfo["screenWidth"].(float64); ok {
			baseInfo.ScreenWidth = int(screenWidth)
		}
		if screenHeight, ok := req.BaseInfo["screenHeight"].(float64); ok {
			baseInfo.ScreenHeight = int(screenHeight)
		}

		// 保存基础信息
		baseInfo, err = model.CreateBaseInfo(baseInfo)
		if err != nil {
			return err
		}
	}

	// 处理事件信息
	for _, eventData := range req.EventInfo {
		// 创建事件
		event := &model.Event{
			ProjectID: project.ID,
		}

		if baseInfo != nil {
			event.BaseInfoID = baseInfo.ID
		}

		// 填充事件字段
		if eventID, ok := eventData["eventId"].(string); ok {
			event.EventID = eventID
		}
		if eventType, ok := eventData["eventType"].(string); ok {
			event.EventType = eventType
		}
		if triggerTime, ok := eventData["triggerTime"].(float64); ok {
			event.TriggerTime = int64(triggerTime)
		}
		if sendTime, ok := eventData["sendTime"].(float64); ok {
			event.SendTime = int64(sendTime)
		}
		if triggerPageURL, ok := eventData["triggerPageUrl"].(string); ok {
			event.TriggerPageURL = triggerPageURL
		}
		if referer, ok := eventData["referer"].(string); ok {
			event.Referer = referer
		}
		if title, ok := eventData["title"].(string); ok {
			event.Title = title
		}

		// 保存事件
		event, err = model.CreateEvent(event)
		if err != nil {
			return err
		}

		// 根据事件类型处理详细信息
		switch event.EventType {
		case "error":
			s.processErrorEvent(event.ID, eventData)
		case "performance":
			s.processPerformanceEvent(event.ID, eventData)
		case "click":
			s.processClickEvent(event.ID, eventData)
		case "pv":
			s.processRouteEvent(event.ID, eventData)
		case "intersection":
			s.processExposureEvent(event.ID, eventData)
		}
	}

	return nil
}

// 处理错误事件
func (s *EventService) processErrorEvent(eventID uint, eventData map[string]interface{}) error {
	errorEvent := &model.ErrorEvent{
		EventID: eventID,
	}

	// 填充错误事件字段
	if errorType, ok := eventData["errorType"].(string); ok {
		errorEvent.ErrorType = errorType
	}
	if errorMessage, ok := eventData["errorMessage"].(string); ok {
		errorEvent.ErrorMessage = errorMessage
	}
	if errorStack, ok := eventData["errorStack"].(string); ok {
		errorEvent.ErrorStack = errorStack
	}
	if componentName, ok := eventData["componentName"].(string); ok {
		errorEvent.ComponentName = componentName
	}
	if filePath, ok := eventData["filePath"].(string); ok {
		errorEvent.FilePath = filePath
	}
	if lineNumber, ok := eventData["lineNumber"].(float64); ok {
		errorEvent.LineNumber = int(lineNumber)
	}
	if columnNumber, ok := eventData["columnNumber"].(float64); ok {
		errorEvent.ColumnNumber = int(columnNumber)
	}
	if recordScreen, ok := eventData["recordScreen"].(string); ok {
		errorEvent.RecordScreen = recordScreen
	} else if recordScreenData, ok := eventData["recordScreen"]; ok {
		// 如果是对象，转换为 JSON 字符串
		recordScreenBytes, err := json.Marshal(recordScreenData)
		if err == nil {
			errorEvent.RecordScreen = string(recordScreenBytes)
		}
	}

	_, err := model.CreateErrorEvent(errorEvent)
	return err
}

// 处理性能事件
func (s *EventService) processPerformanceEvent(eventID uint, eventData map[string]interface{}) error {
	performanceEvent := &model.PerformanceEvent{
		EventID: eventID,
	}

	// 填充性能事件字段
	if resourceType, ok := eventData["resourceType"].(string); ok {
		performanceEvent.ResourceType = resourceType
	}
	if requestURL, ok := eventData["requestUrl"].(string); ok {
		performanceEvent.RequestURL = requestURL
	}
	if duration, ok := eventData["duration"].(float64); ok {
		performanceEvent.Duration = duration
	}
	if responseEnd, ok := eventData["responseEnd"].(float64); ok {
		performanceEvent.ResponseEnd = responseEnd
	}
	if transferSize, ok := eventData["transferSize"].(float64); ok {
		performanceEvent.TransferSize = int64(transferSize)
	}
	if decodedBodySize, ok := eventData["decodedBodySize"].(float64); ok {
		performanceEvent.DecodedBodySize = int64(decodedBodySize)
	}
	if responseStatus, ok := eventData["responseStatus"].(string); ok {
		performanceEvent.ResponseStatus = responseStatus
	}

	_, err := model.CreatePerformanceEvent(performanceEvent)
	return err
}

// 处理请求事件
func (s *EventService) processRequestEvent(eventID uint, eventData map[string]interface{}) error {
	requestEvent := &model.RequestEvent{
		EventID: eventID,
	}

	// 填充请求事件字段
	if requestURL, ok := eventData["requestUrl"].(string); ok {
		requestEvent.RequestURL = requestURL
	}
	if requestMethod, ok := eventData["requestMethod"].(string); ok {
		requestEvent.RequestMethod = requestMethod
	}
	if requestParams, ok := eventData["requestParams"].(string); ok {
		requestEvent.RequestParams = requestParams
	} else if paramsData, ok := eventData["requestParams"]; ok {
		// 如果是对象，转换为 JSON 字符串
		paramsBytes, err := json.Marshal(paramsData)
		if err == nil {
			requestEvent.RequestParams = string(paramsBytes)
		}
	}
	if responseStatus, ok := eventData["responseStatus"].(float64); ok {
		requestEvent.ResponseStatus = int(responseStatus)
	}
	if responseData, ok := eventData["responseData"].(string); ok {
		requestEvent.ResponseData = responseData
	} else if respData, ok := eventData["responseData"]; ok {
		// 如果是对象，转换为 JSON 字符串
		respBytes, err := json.Marshal(respData)
		if err == nil {
			requestEvent.ResponseData = string(respBytes)
		}
	}
	if duration, ok := eventData["duration"].(float64); ok {
		requestEvent.Duration = duration
	}

	_, err := model.CreateRequestEvent(requestEvent)
	return err
}

// 处理点击事件
func (s *EventService) processClickEvent(eventID uint, eventData map[string]interface{}) error {
	clickEvent := &model.ClickEvent{
		EventID: eventID,
	}

	// 填充点击事件字段
	if elementPath, ok := eventData["elementPath"].(string); ok {
		clickEvent.ElementPath = elementPath
	}
	if elementType, ok := eventData["elementType"].(string); ok {
		clickEvent.ElementType = elementType
	}
	if innerText, ok := eventData["innerText"].(string); ok {
		clickEvent.InnerText = innerText
	}

	_, err := model.CreateClickEvent(clickEvent)
	return err
}

// 处理路由事件
func (s *EventService) processRouteEvent(eventID uint, eventData map[string]interface{}) error {
	routeEvent := &model.RouteEvent{
		EventID: eventID,
	}

	// 填充路由事件字段
	if action, ok := eventData["action"].(string); ok {
		routeEvent.Action = action
	}

	_, err := model.CreateRouteEvent(routeEvent)
	return err
}

// 处理曝光事件
func (s *EventService) processExposureEvent(eventID uint, eventData map[string]interface{}) error {
	exposureEvent := &model.ExposureEvent{
		EventID: eventID,
	}

	// 填充曝光事件字段
	if elementPath, ok := eventData["elementPath"].(string); ok {
		exposureEvent.ElementPath = elementPath
	}
	if elementType, ok := eventData["elementType"].(string); ok {
		exposureEvent.ElementType = elementType
	}
	if innerText, ok := eventData["innerText"].(string); ok {
		exposureEvent.InnerText = innerText
	}

	_, err := model.CreateExposureEvent(exposureEvent)
	return err
}

// 获取所有事件
func (s *EventService) GetAllEvents(eventType string, params map[string]interface{}) ([]model.Event, error) {
	// 默认限制返回100条数据
	return model.GetAllEvents(eventType, 100, params)
}

// 清除所有事件数据
func (s *EventService) CleanAllEvents() error {
	return model.CleanAllEvents()
}

// 获取基础信息
func (s *EventService) GetBaseInfo() (map[string]interface{}, error) {
	// 获取最新的一条基础信息
	events, err := model.GetAllEvents("", 1, map[string]interface{}{})
	if err != nil || len(events) == 0 {
		return map[string]interface{}{}, nil
	}

	// 如果有事件，返回其基础信息
	if events[0].BaseInfoID > 0 {
		baseInfo := events[0].BaseInfo
		return map[string]interface{}{
			"userId":         baseInfo.UserID,
			"userAgent":      baseInfo.UserAgent,
			"ip":             baseInfo.IP,
			"browser":        baseInfo.Browser,
			"browserVersion": baseInfo.BrowserVersion,
			"os":             baseInfo.OS,
			"osVersion":      baseInfo.OSVersion,
			"deviceType":     baseInfo.DeviceType,
			"screenWidth":    baseInfo.ScreenWidth,
			"screenHeight":   baseInfo.ScreenHeight,
		}, nil
	}

	return map[string]interface{}{}, nil
}
