package model

import (
	"time"
)

// BaseInfo 基础信息模型
type BaseInfo struct {
	Model
	ProjectID      uint    `json:"projectId" gorm:"not null"`
	Project        Project `json:"-" gorm:"foreignKey:ProjectID"`
	UserID         string  `json:"userId" gorm:"size:100"`
	UserAgent      string  `json:"userAgent" gorm:"type:text"`
	IP             string  `json:"ip" gorm:"size:50"`
	Browser        string  `json:"browser" gorm:"size:50"`
	BrowserVersion string  `json:"browserVersion" gorm:"size:50"`
	OS             string  `json:"os" gorm:"size:50"`
	OSVersion      string  `json:"osVersion" gorm:"size:50"`
	DeviceType     string  `json:"deviceType" gorm:"size:50"`
	ScreenWidth    int     `json:"screenWidth"`
	ScreenHeight   int     `json:"screenHeight"`
}

// Event 事件模型
type Event struct {
	Model
	ProjectID      uint     `json:"projectId" gorm:"not null"`
	Project        Project  `json:"-" gorm:"foreignKey:ProjectID"`
	BaseInfoID     uint     `json:"baseInfoId"`
	BaseInfo       BaseInfo `json:"baseInfo" gorm:"foreignKey:BaseInfoID"`
	EventID        string   `json:"eventId" gorm:"size:50;not null"`
	EventType      string   `json:"eventType" gorm:"size:50;not null"`
	TriggerTime    int64    `json:"triggerTime"`
	SendTime       int64    `json:"sendTime"`
	TriggerPageURL string   `json:"triggerPageUrl" gorm:"type:text"`
	Referer        string   `json:"referer" gorm:"type:text"`
	Title          string   `json:"title" gorm:"size:255"`
}

// ErrorEvent 错误事件模型
type ErrorEvent struct {
	Model
	EventID       uint   `json:"eventId" gorm:"not null"`
	Event         Event  `json:"event" gorm:"foreignKey:EventID"`
	ErrorType     string `json:"errorType" gorm:"size:50"`
	ErrorMessage  string `json:"errorMessage" gorm:"type:text"`
	ErrorStack    string `json:"errorStack" gorm:"type:text"`
	ComponentName string `json:"componentName" gorm:"size:100"`
	FilePath      string `json:"filePath" gorm:"type:text"`
	LineNumber    int    `json:"lineNumber"`
	ColumnNumber  int    `json:"columnNumber"`
	RecordScreen  string `json:"recordScreen" gorm:"type:text"`
}

// PerformanceEvent 性能事件模型
type PerformanceEvent struct {
	Model
	EventID         uint    `json:"eventId" gorm:"not null"`
	Event           Event   `json:"event" gorm:"foreignKey:EventID"`
	ResourceType    string  `json:"resourceType" gorm:"size:50"`
	RequestURL      string  `json:"requestUrl" gorm:"type:text"`
	Duration        float64 `json:"duration"`
	ResponseEnd     float64 `json:"responseEnd"`
	TransferSize    int64   `json:"transferSize"`
	DecodedBodySize int64   `json:"decodedBodySize"`
	ResponseStatus  string  `json:"responseStatus" gorm:"size:50"`
}

// RequestEvent 请求事件模型
type RequestEvent struct {
	Model
	EventID        uint    `json:"eventId" gorm:"not null"`
	Event          Event   `json:"event" gorm:"foreignKey:EventID"`
	RequestURL     string  `json:"requestUrl" gorm:"type:text"`
	RequestMethod  string  `json:"requestMethod" gorm:"size:10"`
	RequestParams  string  `json:"requestParams" gorm:"type:text"`
	ResponseStatus int     `json:"responseStatus"`
	ResponseData   string  `json:"responseData" gorm:"type:text"`
	Duration       float64 `json:"duration"`
}

// ClickEvent 点击事件模型
type ClickEvent struct {
	Model
	EventID     uint   `json:"eventId" gorm:"not null"`
	Event       Event  `json:"event" gorm:"foreignKey:EventID"`
	ElementPath string `json:"elementPath" gorm:"type:text"`
	ElementType string `json:"elementType" gorm:"size:50"`
	InnerText   string `json:"innerText" gorm:"type:text"`
}

// RouteEvent 路由事件模型
type RouteEvent struct {
	Model
	EventID uint   `json:"eventId" gorm:"not null"`
	Event   Event  `json:"event" gorm:"foreignKey:EventID"`
	Action  string `json:"action" gorm:"size:50"`
}

// ExposureEvent 曝光事件模型
type ExposureEvent struct {
	Model
	EventID     uint   `json:"eventId" gorm:"not null"`
	Event       Event  `json:"event" gorm:"foreignKey:EventID"`
	ElementPath string `json:"elementPath" gorm:"type:text"`
	ElementType string `json:"elementType" gorm:"size:50"`
	InnerText   string `json:"innerText" gorm:"type:text"`
}

// 创建基础信息
func CreateBaseInfo(baseInfo *BaseInfo) (*BaseInfo, error) {
	if err := db.Create(baseInfo).Error; err != nil {
		return nil, err
	}
	return baseInfo, nil
}

// 创建事件
func CreateEvent(event *Event) (*Event, error) {
	if err := db.Create(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

// 创建错误事件
func CreateErrorEvent(errorEvent *ErrorEvent) (*ErrorEvent, error) {
	if err := db.Create(errorEvent).Error; err != nil {
		return nil, err
	}
	return errorEvent, nil
}

// 创建性能事件
func CreatePerformanceEvent(performanceEvent *PerformanceEvent) (*PerformanceEvent, error) {
	if err := db.Create(performanceEvent).Error; err != nil {
		return nil, err
	}
	return performanceEvent, nil
}

// 创建请求事件
func CreateRequestEvent(requestEvent *RequestEvent) (*RequestEvent, error) {
	if err := db.Create(requestEvent).Error; err != nil {
		return nil, err
	}
	return requestEvent, nil
}

// 创建点击事件
func CreateClickEvent(clickEvent *ClickEvent) (*ClickEvent, error) {
	if err := db.Create(clickEvent).Error; err != nil {
		return nil, err
	}
	return clickEvent, nil
}

// 创建路由事件
func CreateRouteEvent(routeEvent *RouteEvent) (*RouteEvent, error) {
	if err := db.Create(routeEvent).Error; err != nil {
		return nil, err
	}
	return routeEvent, nil
}

// 创建曝光事件
func CreateExposureEvent(exposureEvent *ExposureEvent) (*ExposureEvent, error) {
	if err := db.Create(exposureEvent).Error; err != nil {
		return nil, err
	}
	return exposureEvent, nil
}

// 获取项目的所有事件
func GetEventsByProjectID(projectID uint, eventType string, page, pageSize int, params map[string]interface{}) ([]Event, int64, error) {
	var events []Event
	var total int64

	// 基本查询条件
	query := db.Where("project_id = ?", projectID)
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	// 时间范围筛选
	if startTimeStr, ok := params["startTime"].(string); ok && startTimeStr != "" {
		startTime, err := time.Parse("2006-01-02", startTimeStr)
		if err == nil {
			query = query.Where("trigger_time >= ?", startTime.Unix()*1000)
		}
	}

	if endTimeStr, ok := params["endTime"].(string); ok && endTimeStr != "" {
		endTime, err := time.Parse("2006-01-02", endTimeStr)
		if err == nil {
			// 将结束时间设置为当天的最后一秒
			endTime = endTime.Add(24*time.Hour - time.Second)
			query = query.Where("trigger_time <= ?", endTime.Unix()*1000)
		}
	}

	// 关键字搜索
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		keywordLike := "%" + keyword + "%"
		query = query.Where(
			"title LIKE ? OR trigger_page_url LIKE ? OR referer LIKE ?",
			keywordLike, keywordLike, keywordLike,
		)
	}

	// 浏览器筛选
	if browser, ok := params["browser"].(string); ok && browser != "" {
		// 连接 BaseInfo 表进行查询
		query = query.Joins("JOIN wt_base_info ON wt_event.base_info_id = wt_base_info.id")
		query = query.Where("wt_base_info.browser = ?", browser)
	}

	// 操作系统筛选
	if os, ok := params["os"].(string); ok && os != "" {
		// 如果还没有连接 BaseInfo 表，则进行连接
		if _, exists := params["browser"].(string); !exists || params["browser"].(string) == "" {
			query = query.Joins("JOIN wt_base_info ON wt_event.base_info_id = wt_base_info.id")
		}
		query = query.Where("wt_base_info.os = ?", os)
	}

	// 设备类型筛选
	if deviceType, ok := params["deviceType"].(string); ok && deviceType != "" {
		// 如果还没有连接 BaseInfo 表，则进行连接
		existsBrowser := false
		if _, ok := params["browser"].(string); ok && params["browser"].(string) != "" {
			existsBrowser = true
		}

		existsOS := false
		if _, ok := params["os"].(string); ok && params["os"].(string) != "" {
			existsOS = true
		}

		if !existsBrowser && !existsOS {
			query = query.Joins("JOIN wt_base_info ON wt_event.base_info_id = wt_base_info.id")
		}
		query = query.Where("wt_base_info.device_type = ?", deviceType)
	}

	// 错误类型筛选，仅当 eventType=error 时有效
	if errorType, ok := params["errorType"].(string); ok && errorType != "" && eventType == "error" {
		query = query.Joins("JOIN wt_error_event ON wt_event.id = wt_error_event.event_id")
		query = query.Where("wt_error_event.error_type = ?", errorType)
	}

	// 响应状态码筛选，仅当 eventType=request 时有效
	if responseStatus, ok := params["responseStatus"].(string); ok && responseStatus != "" && eventType == "request" {
		query = query.Joins("JOIN wt_request_event ON wt_event.id = wt_request_event.event_id")
		query = query.Where("wt_request_event.response_status = ?", responseStatus)
	}

	// 获取总数
	if err := query.Model(&Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("BaseInfo").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// 获取事件详情
func GetEventByID(id uint) (*Event, error) {
	var event Event
	if err := db.Preload("BaseInfo").First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// 获取错误事件详情
func GetErrorEventByEventID(eventID uint) (*ErrorEvent, error) {
	var errorEvent ErrorEvent
	if err := db.Where("event_id = ?", eventID).First(&errorEvent).Error; err != nil {
		return nil, err
	}
	return &errorEvent, nil
}

// 获取性能事件详情
func GetPerformanceEventByEventID(eventID uint) (*PerformanceEvent, error) {
	var performanceEvent PerformanceEvent
	if err := db.Where("event_id = ?", eventID).First(&performanceEvent).Error; err != nil {
		return nil, err
	}
	return &performanceEvent, nil
}

// 获取请求事件详情
func GetRequestEventByEventID(eventID uint) (*RequestEvent, error) {
	var requestEvent RequestEvent
	if err := db.Where("event_id = ?", eventID).First(&requestEvent).Error; err != nil {
		return nil, err
	}
	return &requestEvent, nil
}

// 获取点击事件详情
func GetClickEventByEventID(eventID uint) (*ClickEvent, error) {
	var clickEvent ClickEvent
	if err := db.Where("event_id = ?", eventID).First(&clickEvent).Error; err != nil {
		return nil, err
	}
	return &clickEvent, nil
}

// 获取路由事件详情
func GetRouteEventByEventID(eventID uint) (*RouteEvent, error) {
	var routeEvent RouteEvent
	if err := db.Where("event_id = ?", eventID).First(&routeEvent).Error; err != nil {
		return nil, err
	}
	return &routeEvent, nil
}

// 获取曝光事件详情
func GetExposureEventByEventID(eventID uint) (*ExposureEvent, error) {
	var exposureEvent ExposureEvent
	if err := db.Where("event_id = ?", eventID).First(&exposureEvent).Error; err != nil {
		return nil, err
	}
	return &exposureEvent, nil
}

// 获取所有事件
func GetAllEvents(eventType string, limit int, params map[string]interface{}) ([]Event, error) {
	var events []Event
	query := db

	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	// 时间范围筛选
	if startTimeStr, ok := params["startTime"].(string); ok && startTimeStr != "" {
		startTime, err := time.Parse("2006-01-02", startTimeStr)
		if err == nil {
			query = query.Where("trigger_time >= ?", startTime.Unix()*1000)
		}
	}

	if endTimeStr, ok := params["endTime"].(string); ok && endTimeStr != "" {
		endTime, err := time.Parse("2006-01-02", endTimeStr)
		if err == nil {
			// 将结束时间设置为当天的最后一秒
			endTime = endTime.Add(24*time.Hour - time.Second)
			query = query.Where("trigger_time <= ?", endTime.Unix()*1000)
		}
	}

	// 关键字搜索
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		keywordLike := "%" + keyword + "%"
		query = query.Where(
			"title LIKE ? OR trigger_page_url LIKE ? OR referer LIKE ?",
			keywordLike, keywordLike, keywordLike,
		)
	}

	// 浏览器筛选
	if browser, ok := params["browser"].(string); ok && browser != "" {
		// 连接 BaseInfo 表进行查询
		query = query.Joins("JOIN wt_base_info ON wt_event.base_info_id = wt_base_info.id")
		query = query.Where("wt_base_info.browser = ?", browser)
	}

	// 操作系统筛选
	if os, ok := params["os"].(string); ok && os != "" {
		// 如果还没有连接 BaseInfo 表，则进行连接
		if _, exists := params["browser"].(string); !exists || params["browser"].(string) == "" {
			query = query.Joins("JOIN wt_base_info ON wt_event.base_info_id = wt_base_info.id")
		}
		query = query.Where("wt_base_info.os = ?", os)
	}

	// 设备类型筛选
	if deviceType, ok := params["deviceType"].(string); ok && deviceType != "" {
		// 如果还没有连接 BaseInfo 表，则进行连接
		existsBrowser := false
		if _, ok := params["browser"].(string); ok && params["browser"].(string) != "" {
			existsBrowser = true
		}

		existsOS := false
		if _, ok := params["os"].(string); ok && params["os"].(string) != "" {
			existsOS = true
		}

		if !existsBrowser && !existsOS {
			query = query.Joins("JOIN wt_base_info ON wt_event.base_info_id = wt_base_info.id")
		}
		query = query.Where("wt_base_info.device_type = ?", deviceType)
	}

	// 错误类型筛选，仅当 eventType=error 时有效
	if errorType, ok := params["errorType"].(string); ok && errorType != "" && eventType == "error" {
		query = query.Joins("JOIN wt_error_event ON wt_event.id = wt_error_event.event_id")
		query = query.Where("wt_error_event.error_type = ?", errorType)
	}

	// 响应状态码筛选，仅当 eventType=request 时有效
	if responseStatus, ok := params["responseStatus"].(string); ok && responseStatus != "" && eventType == "request" {
		query = query.Joins("JOIN wt_request_event ON wt_event.id = wt_request_event.event_id")
		query = query.Where("wt_request_event.response_status = ?", responseStatus)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Preload("BaseInfo").Order("created_at DESC").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// 清除所有事件数据
func CleanAllEvents() error {
	// 先删除所有子表数据
	if err := db.Exec("DELETE FROM wt_error_event").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM wt_performance_event").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM wt_request_event").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM wt_click_event").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM wt_route_event").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM wt_exposure_event").Error; err != nil {
		return err
	}

	// 删除事件表数据
	if err := db.Exec("DELETE FROM wt_event").Error; err != nil {
		return err
	}

	// 删除基础信息表数据
	if err := db.Exec("DELETE FROM wt_base_info").Error; err != nil {
		return err
	}

	return nil
}

// 获取项目的错误事件统计
func GetErrorEventStats(projectID uint, startTime, endTime time.Time) (map[string]int64, error) {
	var results []struct {
		Date  string
		Count int64
	}

	query := `
		SELECT
			DATE(FROM_UNIXTIME(e.trigger_time/1000)) as date,
			COUNT(*) as count
		FROM
			wt_event e
		JOIN
			wt_error_event ee ON e.id = ee.event_id
		WHERE
			e.project_id = ?
			AND e.event_type = 'error'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
		GROUP BY
			date
		ORDER BY
			date
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.Date] = result.Count
	}

	return stats, nil
}

// 获取错误事件按浏览器分布统计
func GetErrorEventByBrowserStats(projectID uint, startTime, endTime time.Time) (map[string]int64, error) {
	var results []struct {
		Browser string
		Count   int64
	}

	query := `
		SELECT
			b.browser,
			COUNT(*) as count
		FROM
			wt_event e
		JOIN
			wt_error_event ee ON e.id = ee.event_id
		JOIN
			wt_base_info b ON e.base_info_id = b.id
		WHERE
			e.project_id = ?
			AND e.event_type = 'error'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
		GROUP BY
			b.browser
		ORDER BY
			count DESC
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.Browser] = result.Count
	}

	return stats, nil
}

// 获取错误事件按操作系统分布统计
func GetErrorEventByOSStats(projectID uint, startTime, endTime time.Time) (map[string]int64, error) {
	var results []struct {
		OS    string
		Count int64
	}

	query := `
		SELECT
			b.os,
			COUNT(*) as count
		FROM
			wt_event e
		JOIN
			wt_error_event ee ON e.id = ee.event_id
		JOIN
			wt_base_info b ON e.base_info_id = b.id
		WHERE
			e.project_id = ?
			AND e.event_type = 'error'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
		GROUP BY
			b.os
		ORDER BY
			count DESC
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.OS] = result.Count
	}

	return stats, nil
}

// 获取错误事件按设备类型分布统计
func GetErrorEventByDeviceStats(projectID uint, startTime, endTime time.Time) (map[string]int64, error) {
	var results []struct {
		DeviceType string
		Count      int64
	}

	query := `
		SELECT
			b.device_type,
			COUNT(*) as count
		FROM
			wt_event e
		JOIN
			wt_error_event ee ON e.id = ee.event_id
		JOIN
			wt_base_info b ON e.base_info_id = b.id
		WHERE
			e.project_id = ?
			AND e.event_type = 'error'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
		GROUP BY
			b.device_type
		ORDER BY
			count DESC
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.DeviceType] = result.Count
	}

	return stats, nil
}

// 获取错误事件按错误类型分布统计
func GetErrorEventByTypeStats(projectID uint, startTime, endTime time.Time) (map[string]int64, error) {
	var results []struct {
		ErrorType string
		Count     int64
	}

	query := `
		SELECT
			ee.error_type,
			COUNT(*) as count
		FROM
			wt_event e
		JOIN
			wt_error_event ee ON e.id = ee.event_id
		WHERE
			e.project_id = ?
			AND e.event_type = 'error'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
		GROUP BY
			ee.error_type
		ORDER BY
			count DESC
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.ErrorType] = result.Count
	}

	return stats, nil
}

// 获取性能指标统计
func GetPerformanceStats(projectID uint, startTime, endTime time.Time) (map[string]float64, error) {
	var results []struct {
		ResourceType string
		AvgDuration  float64
	}

	query := `
		SELECT
			pe.resource_type,
			AVG(pe.duration) as avg_duration
		FROM
			wt_event e
		JOIN
			wt_performance_event pe ON e.id = pe.event_id
		WHERE
			e.project_id = ?
			AND e.event_type = 'performance'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
		GROUP BY
			pe.resource_type
		ORDER BY
			avg_duration DESC
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]float64)
	for _, result := range results {
		stats[result.ResourceType] = result.AvgDuration
	}

	return stats, nil
}

// 获取请求错误统计
func GetRequestErrorStats(projectID uint, startTime, endTime time.Time) (map[int]int64, error) {
	var results []struct {
		ResponseStatus int
		Count          int64
	}

	query := `
		SELECT
			re.response_status,
			COUNT(*) as count
		FROM
			wt_event e
		JOIN
			wt_request_event re ON e.id = re.event_id
		WHERE
			e.project_id = ?
			AND e.event_type = 'request'
			AND e.trigger_time >= ?
			AND e.trigger_time <= ?
			AND re.response_status >= 400
		GROUP BY
			re.response_status
		ORDER BY
			count DESC
	`

	if err := db.Raw(query, projectID, startTime.Unix()*1000, endTime.Unix()*1000).Scan(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[int]int64)
	for _, result := range results {
		stats[result.ResponseStatus] = result.Count
	}

	return stats, nil
}
