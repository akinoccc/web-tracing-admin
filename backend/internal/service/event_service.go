package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akinoccc/web-tracing-admin/internal/model"
)

// TrackRequest SDK上报数据请求
type TrackRequest struct {
	// 事件类别
	Category string `json:"category"`
	// 事件类型
	Type string `json:"type"`
	// 子类型（可选），更细的分类
	SubType string `json:"subType,omitempty"`
	// 事件发生时间戳
	Timestamp int64 `json:"timestamp"`
	// 错误严重程度（错误类事件专用）
	Severity string `json:"severity,omitempty"`
	// 错误指纹（用于错误去重）
	Fingerprint string `json:"fingerprint,omitempty"`
	// 事件数据
	Data json.RawMessage `json:"data"`
	// 环境信息
	Environment string `json:"environment,omitempty"`
	// 版本信息
	Release string `json:"release,omitempty"`
	// 应用标识（从请求头或查询参数获取）
	AppKey string `json:"appKey,omitempty"`
}

// ErrorListResponse 错误列表响应
type ErrorListResponse struct {
	Total int64            `json:"total"`
	List  []ErrorGroupItem `json:"list"`
	Stats ErrorStatsData   `json:"stats"`
}

// ErrorGroupItem 错误分组项
type ErrorGroupItem struct {
	ID           uint   `json:"id"`
	Fingerprint  string `json:"fingerprint"`
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
	Count        int    `json:"count"`
	FirstSeen    int64  `json:"firstSeen"`
	LastSeen     int64  `json:"lastSeen"`
	Status       string `json:"status"`
	Severity     string `json:"severity"`
	SubType      string `json:"subType"`
}

// ErrorDetailResponse 错误详情响应
type ErrorDetailResponse struct {
	Group  ErrorGroupItem   `json:"group"`
	Events []ErrorEventItem `json:"events"`
	Total  int64            `json:"total"`
}

// ErrorEventItem 错误事件项
type ErrorEventItem struct {
	ID           uint   `json:"id"`
	EventID      string `json:"eventId"`
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
	ErrorStack   string `json:"errorStack"`
	FilePath     string `json:"filePath"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`
	TriggerTime  int64  `json:"triggerTime"`
	PageURL      string `json:"pageUrl"`
	Browser      string `json:"browser"`
	OS           string `json:"os"`
	Device       string `json:"device"`
}

// ErrorStatsResponse 错误统计响应
type ErrorStatsResponse struct {
	Stats ErrorStatsData   `json:"stats"`
	Trend []ErrorTrendItem `json:"trend"`
}

// ErrorStatsData 错误统计数据
type ErrorStatsData struct {
	TotalErrors         int64            `json:"totalErrors"`
	AffectedUsers       int64            `json:"affectedUsers"`
	ErrorsToday         int64            `json:"errorsToday"`
	ErrorsYesterday     int64            `json:"errorsYesterday"`
	TypeDistribution    map[string]int64 `json:"typeDistribution"`
	BrowserDistribution map[string]int64 `json:"browserDistribution"`
	OSDistribution      map[string]int64 `json:"osDistribution"`
}

// ErrorTrendItem 错误趋势项
type ErrorTrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// PerformanceListResponse 性能列表响应
type PerformanceListResponse struct {
	Total int64                 `json:"total"`
	List  []PerformancePageItem `json:"list"`
	Stats PerformanceStatsData  `json:"stats"`
}

// PerformancePageItem 性能页面项
type PerformancePageItem struct {
	ID          uint    `json:"id"`
	EventID     string  `json:"eventId"`
	PageURL     string  `json:"pageUrl"`
	TriggerTime int64   `json:"triggerTime"`
	FP          int64   `json:"fp"`
	FCP         int64   `json:"fcp"`
	LCP         int64   `json:"lcp"`
	FID         int64   `json:"fid"`
	CLS         float64 `json:"cls"`
	TTFB        int64   `json:"ttfb"`
	DomReady    int64   `json:"domReady"`
	Load        int64   `json:"load"`
	Browser     string  `json:"browser"`
	OS          string  `json:"os"`
	Device      string  `json:"device"`
}

// ResourcePerformanceListResponse 资源性能列表响应
type ResourcePerformanceListResponse struct {
	Total int64                     `json:"total"`
	List  []ResourcePerformanceItem `json:"list"`
}

// ResourcePerformanceItem 资源性能项
type ResourcePerformanceItem struct {
	ID            uint   `json:"id"`
	EventID       string `json:"eventId"`
	ResourceURL   string `json:"resourceUrl"`
	ResourceType  string `json:"resourceType"`
	InitiatorType string `json:"initiatorType"`
	StartTime     int64  `json:"startTime"`
	Duration      int64  `json:"duration"`
	TransferSize  int64  `json:"transferSize"`
	PageURL       string `json:"pageUrl"`
}

// PerformanceStatsResponse 性能统计响应
type PerformanceStatsResponse struct {
	Stats PerformanceStatsData   `json:"stats"`
	Trend []PerformanceTrendItem `json:"trend"`
}

// PerformanceStatsData 性能统计数据
type PerformanceStatsData struct {
	AvgFP       int64   `json:"avgFP"`
	AvgFCP      int64   `json:"avgFCP"`
	AvgLCP      int64   `json:"avgLCP"`
	AvgFID      int64   `json:"avgFID"`
	AvgCLS      float64 `json:"avgCLS"`
	AvgTTFB     int64   `json:"avgTTFB"`
	AvgDomReady int64   `json:"avgDomReady"`
	AvgLoad     int64   `json:"avgLoad"`
}

// PerformanceTrendItem 性能趋势项
type PerformanceTrendItem struct {
	Date string `json:"date"`
	FP   int64  `json:"fp"`
	FCP  int64  `json:"fcp"`
	LCP  int64  `json:"lcp"`
	TTFB int64  `json:"ttfb"`
}

// PVListResponse 页面访问列表响应
type PVListResponse struct {
	Total int64    `json:"total"`
	List  []PVItem `json:"list"`
}

// PVItem 页面访问项
type PVItem struct {
	ID          uint   `json:"id"`
	EventID     string `json:"eventId"`
	PageURL     string `json:"pageUrl"`
	Title       string `json:"title"`
	Referrer    string `json:"referrer"`
	TriggerTime int64  `json:"triggerTime"`
	StayTime    int64  `json:"stayTime"`
	IsNewVisit  bool   `json:"isNewVisit"`
	Browser     string `json:"browser"`
	OS          string `json:"os"`
	Device      string `json:"device"`
}

// ClickListResponse 点击列表响应
type ClickListResponse struct {
	Total int64       `json:"total"`
	List  []ClickItem `json:"list"`
}

// ClickItem 点击项
type ClickItem struct {
	ID          uint   `json:"id"`
	EventID     string `json:"eventId"`
	ElementPath string `json:"elementPath"`
	ElementType string `json:"elementType"`
	InnerText   string `json:"innerText"`
	TriggerTime int64  `json:"triggerTime"`
	PageURL     string `json:"pageUrl"`
}

// BehaviorStatsResponse 用户行为统计响应
type BehaviorStatsResponse struct {
	PVStats    PVStatsData    `json:"pvStats"`
	ClickStats ClickStatsData `json:"clickStats"`
	PVTrend    []PVTrendItem  `json:"pvTrend"`
}

// PVStatsData 页面访问统计数据
type PVStatsData struct {
	TotalPV     int64            `json:"totalPV"`
	TotalUV     int64            `json:"totalUV"`
	PVToday     int64            `json:"pvToday"`
	UVToday     int64            `json:"uvToday"`
	AvgStayTime int64            `json:"avgStayTime"`
	BounceRate  float64          `json:"bounceRate"`
	TopPages    map[string]int64 `json:"topPages"`
}

// ClickStatsData 点击统计数据
type ClickStatsData struct {
	TotalClicks int64            `json:"totalClicks"`
	ClicksToday int64            `json:"clicksToday"`
	TopElements map[string]int64 `json:"topElements"`
}

// PVTrendItem 页面访问趋势项
type PVTrendItem struct {
	Date string `json:"date"`
	PV   int64  `json:"pv"`
	UV   int64  `json:"uv"`
}

type EventService struct{}

// 生成事件ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + strconv.FormatInt(time.Now().UnixNano()%1000000, 10)
}

// ProcessTrackData 处理上报数据
func (s *EventService) ProcessTrackData(req *TrackRequest) error {
	// 创建基础信息
	baseInfo := model.BaseInfo{
		AppKey:   req.AppKey,
		SendTime: req.Timestamp,
		// 其他字段将从事件数据中提取
	}

	// 查找项目
	var project model.Project
	db := model.GetDB()
	if err := db.Where("app_key = ?", req.AppKey).First(&project).Error; err != nil {
		return errors.New("项目不存在")
	}

	// 设置项目ID
	baseInfo.ProjectID = project.ID

	// 从事件数据中提取通用信息
	if req.Data != nil {
		// 尝试从数据中提取URL、用户信息等
		var dataMap map[string]interface{}
		if err := json.Unmarshal(req.Data, &dataMap); err == nil {
			// 提取页面URL
			if url, ok := dataMap["url"].(string); ok {
				baseInfo.PageURL = url
			}

			// 提取用户ID
			if userId, ok := dataMap["userId"].(string); ok {
				baseInfo.UserID = userId
			}

			// 提取会话ID
			if sessionId, ok := dataMap["sessionId"].(string); ok {
				baseInfo.SessionID = sessionId
			}

			// 提取引用页
			if referrer, ok := dataMap["referrer"].(string); ok {
				baseInfo.Referrer = referrer
			}
		}
	}

	// 保存基础信息
	if err := db.Create(&baseInfo).Error; err != nil {
		return err
	}

	// 将SDK事件类型映射到后端事件类型
	var eventType string
	switch req.Category {
	case "error":
		eventType = model.EventTypeError
	case "performance":
		if req.Type == "web_vitals" || req.Type == "page_load" {
			eventType = model.EventTypePerformancePage
		} else if req.Type == "resource_load" {
			eventType = model.EventTypePerformanceResource
		} else {
			eventType = model.EventTypePerformancePage // 默认
		}
	case "user":
		if req.Type == "page_view" {
			eventType = model.EventTypePV
		} else if req.Type == "click" {
			eventType = model.EventTypeClick
		} else if req.Type == "stay_time" {
			eventType = model.EventTypeDwell
		} else {
			eventType = model.EventTypePV // 默认
		}
	case "custom":
		eventType = model.EventTypeCustom
	case "system":
		if req.Type == "batch_report" {
			// 处理批量上报
			var batchData struct {
				Events []json.RawMessage `json:"events"`
				Count  int               `json:"count"`
			}
			if err := json.Unmarshal(req.Data, &batchData); err == nil {
				// 处理批量事件
				for _, eventData := range batchData.Events {
					var batchEvent TrackRequest
					if err := json.Unmarshal(eventData, &batchEvent); err == nil {
						// 递归处理每个事件
						if err := s.ProcessTrackData(&batchEvent); err != nil {
							// 记录错误但继续处理
							fmt.Printf("批量处理事件错误: %v\n", err)
						}
					}
				}
				return nil
			}
		}
		// 其他系统事件暂不处理
		return nil
	default:
		return errors.New("不支持的事件类别")
	}

	// 创建事件主记录
	eventMain := model.EventMain{
		EventID:        generateEventID(),
		EventType:      eventType,
		ProjectID:      project.ID,
		BaseInfoID:     baseInfo.ID,
		TriggerTime:    req.Timestamp,
		SendTime:       time.Now().Unix(),
		TriggerPageURL: baseInfo.PageURL,
		Title:          "",
		Referer:        baseInfo.Referrer,
	}

	if err := db.Create(&eventMain).Error; err != nil {
		return err
	}

	// 根据事件类型处理详情
	switch eventType {
	case model.EventTypeError:
		return s.processErrorEventFromSDK(req, eventMain.ID, project.ID)
	case model.EventTypePerformancePage:
		return s.processPerformancePageEventFromSDK(req, eventMain.ID)
	case model.EventTypePerformanceResource:
		return s.processPerformanceResourceEventFromSDK(req, eventMain.ID)
	case model.EventTypePV:
		return s.processPVEventFromSDK(req, eventMain.ID)
	case model.EventTypeClick:
		return s.processClickEventFromSDK(req, eventMain.ID)
	case model.EventTypeDwell:
		return s.processDwellEventFromSDK(req, eventMain.ID)
	case model.EventTypeCustom:
		return s.processCustomEventFromSDK(req, eventMain.ID)
	default:
		return errors.New("不支持的事件类型")
	}
}

// 处理错误事件
func (s *EventService) processErrorEvent(detail json.RawMessage, eventID uint, projectID uint) error {
	var errorDetail model.ErrorDetail
	if err := json.Unmarshal(detail, &errorDetail); err != nil {
		return err
	}

	errorDetail.EventID = eventID
	db := model.GetDB()
	if err := db.Create(&errorDetail).Error; err != nil {
		return err
	}

	// 创建或更新错误分组
	_, err := model.CreateOrUpdateErrorGroup(
		errorDetail.Fingerprint,
		errorDetail.ErrorType,
		errorDetail.ErrorMessage,
		projectID,
		eventID,
		errorDetail.Severity,
		errorDetail.SubType,
	)
	return err
}

// 处理性能页面事件
func (s *EventService) processPerformancePageEvent(detail json.RawMessage, eventID uint) error {
	var perfDetail model.PerformancePageDetail
	if err := json.Unmarshal(detail, &perfDetail); err != nil {
		return err
	}

	perfDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&perfDetail).Error
}

// 处理性能资源事件
func (s *EventService) processPerformanceResourceEvent(detail json.RawMessage, eventID uint) error {
	var resourceDetail model.PerformanceResourceDetail
	if err := json.Unmarshal(detail, &resourceDetail); err != nil {
		return err
	}

	resourceDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&resourceDetail).Error
}

// 处理页面访问事件
func (s *EventService) processPVEvent(detail json.RawMessage, eventID uint) error {
	var pvDetail model.PVDetail
	if err := json.Unmarshal(detail, &pvDetail); err != nil {
		return err
	}

	pvDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&pvDetail).Error
}

// 处理点击事件
func (s *EventService) processClickEvent(detail json.RawMessage, eventID uint) error {
	var clickDetail model.ClickDetail
	if err := json.Unmarshal(detail, &clickDetail); err != nil {
		return err
	}

	clickDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&clickDetail).Error
}

// 处理停留事件
func (s *EventService) processDwellEvent(detail json.RawMessage, eventID uint) error {
	var dwellDetail model.DwellDetail
	if err := json.Unmarshal(detail, &dwellDetail); err != nil {
		return err
	}

	dwellDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&dwellDetail).Error
}

// 处理曝光事件
func (s *EventService) processIntersectionEvent(detail json.RawMessage, eventID uint) error {
	var intersectionDetail model.IntersectionDetail
	if err := json.Unmarshal(detail, &intersectionDetail); err != nil {
		return err
	}

	intersectionDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&intersectionDetail).Error
}

// 处理自定义事件
func (s *EventService) processCustomEvent(detail json.RawMessage, eventID uint) error {
	var customDetail model.CustomDetail
	if err := json.Unmarshal(detail, &customDetail); err != nil {
		return err
	}

	customDetail.EventID = eventID
	db := model.GetDB()
	return db.Create(&customDetail).Error
}

// 从SDK错误事件处理
func (s *EventService) processErrorEventFromSDK(req *TrackRequest, eventID uint, projectID uint) error {
	// 创建错误详情
	errorDetail := model.ErrorDetail{
		EventID:     eventID,
		ErrorType:   req.Type,
		SubType:     req.SubType,
		Severity:    req.Severity,
		Fingerprint: req.Fingerprint,
	}

	// 从事件数据中提取错误信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// 提取错误消息
		if message, ok := dataMap["message"].(string); ok {
			errorDetail.ErrorMessage = message
		}

		// 提取错误堆栈
		if stack, ok := dataMap["stack"].(string); ok {
			errorDetail.ErrorStack = stack
		}

		// 提取文件路径
		if filename, ok := dataMap["filename"].(string); ok {
			errorDetail.FilePath = filename
		}

		// 提取行号
		if lineno, ok := dataMap["lineno"].(float64); ok {
			errorDetail.LineNumber = int(lineno)
		}

		// 提取列号
		if colno, ok := dataMap["colno"].(float64); ok {
			errorDetail.ColumnNumber = int(colno)
		}

		// 提取组件名称（Vue/React错误）
		if componentName, ok := dataMap["componentName"].(string); ok {
			errorDetail.ComponentName = componentName
		}
	}

	// 保存错误详情
	db := model.GetDB()
	if err := db.Create(&errorDetail).Error; err != nil {
		return err
	}

	// 创建或更新错误分组
	_, err := model.CreateOrUpdateErrorGroup(
		errorDetail.Fingerprint,
		errorDetail.ErrorType,
		errorDetail.ErrorMessage,
		projectID,
		eventID,
		errorDetail.Severity,
		errorDetail.SubType,
	)
	return err
}

// 从SDK性能页面事件处理
func (s *EventService) processPerformancePageEventFromSDK(req *TrackRequest, eventID uint) error {
	// 创建性能页面详情
	perfDetail := model.PerformancePageDetail{
		EventID: eventID,
	}

	// 从事件数据中提取性能信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// Web Vitals事件
		if req.Type == "web_vitals" {
			// 提取指标名称和值
			if name, ok := dataMap["name"].(string); ok {
				switch name {
				case "FP":
					if value, ok := dataMap["value"].(float64); ok {
						perfDetail.FP = int64(value)
					}
				case "FCP":
					if value, ok := dataMap["value"].(float64); ok {
						perfDetail.FCP = int64(value)
					}
				case "LCP":
					if value, ok := dataMap["value"].(float64); ok {
						perfDetail.LCP = int64(value)
					}
				case "FID":
					if value, ok := dataMap["value"].(float64); ok {
						perfDetail.FID = int64(value)
					}
				case "CLS":
					if value, ok := dataMap["value"].(float64); ok {
						perfDetail.CLS = value
					}
				case "TTFB":
					if value, ok := dataMap["value"].(float64); ok {
						perfDetail.TTFB = int64(value)
					}
				}
			}
		} else if req.Type == "page_load" {
			// 页面加载事件
			// 提取导航计时信息
			if paintTiming, ok := dataMap["paintTiming"].(map[string]interface{}); ok {
				if fp, ok := paintTiming["FP"].(float64); ok {
					perfDetail.FP = int64(fp)
				}
				if fcp, ok := paintTiming["FCP"].(float64); ok {
					perfDetail.FCP = int64(fcp)
				}
			}

			// 提取加载时间
			if loadTime, ok := dataMap["loadTime"].(float64); ok {
				perfDetail.Load = int64(loadTime)
			}

			// 提取DOM内容加载时间
			if domContentLoadedTime, ok := dataMap["domContentLoadedTime"].(float64); ok {
				perfDetail.DomReady = int64(domContentLoadedTime)
			}
		}
	}

	// 保存性能页面详情
	db := model.GetDB()
	return db.Create(&perfDetail).Error
}

// 从SDK性能资源事件处理
func (s *EventService) processPerformanceResourceEventFromSDK(req *TrackRequest, eventID uint) error {
	// 创建性能资源详情
	resourceDetail := model.PerformanceResourceDetail{
		EventID: eventID,
	}

	// 从事件数据中提取资源信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// 提取资源URL
		if url, ok := dataMap["url"].(string); ok {
			resourceDetail.ResourceURL = url
		}

		// 提取资源类型
		if initiatorType, ok := dataMap["initiatorType"].(string); ok {
			resourceDetail.InitiatorType = initiatorType
		}

		// 提取资源子类型
		resourceDetail.ResourceType = req.SubType

		// 提取持续时间
		if duration, ok := dataMap["duration"].(float64); ok {
			resourceDetail.Duration = int64(duration)
		}

		// 提取传输大小
		if transferSize, ok := dataMap["transferSize"].(float64); ok {
			resourceDetail.TransferSize = int64(transferSize)
		}

		// 提取解码体大小
		if decodedBodySize, ok := dataMap["decodedBodySize"].(float64); ok {
			resourceDetail.DecodedBodySize = int64(decodedBodySize)
		}
	}

	// 保存性能资源详情
	db := model.GetDB()
	return db.Create(&resourceDetail).Error
}

// 从SDK页面访问事件处理
func (s *EventService) processPVEventFromSDK(req *TrackRequest, eventID uint) error {
	// 创建页面访问详情
	pvDetail := model.PVDetail{
		EventID: eventID,
	}

	// 从事件数据中提取页面信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// 提取页面URL
		if url, ok := dataMap["url"].(string); ok {
			pvDetail.PageURL = url
		}

		// 提取页面标题
		if title, ok := dataMap["title"].(string); ok {
			pvDetail.Title = title
		}

		// 提取引用页
		if referrer, ok := dataMap["referrer"].(string); ok {
			pvDetail.Referrer = referrer
		}

		// 提取停留时间（如果是停留时间事件）
		if req.Type == "stay_time" {
			if duration, ok := dataMap["duration"].(float64); ok {
				pvDetail.StayTime = int64(duration)
			}
		}
	}

	// 保存页面访问详情
	db := model.GetDB()
	return db.Create(&pvDetail).Error
}

// 从SDK点击事件处理
func (s *EventService) processClickEventFromSDK(req *TrackRequest, eventID uint) error {
	// 创建点击详情
	clickDetail := model.ClickDetail{
		EventID: eventID,
	}

	// 从事件数据中提取点击信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// 提取元素路径
		if path, ok := dataMap["path"].([]interface{}); ok && len(path) > 0 {
			// 将路径数组转换为字符串
			pathStr := ""
			for i, p := range path {
				if i > 0 {
					pathStr += " > "
				}
				pathStr += fmt.Sprintf("%v", p)
			}
			clickDetail.ElementPath = pathStr
		}

		// 提取元素类型
		if tagName, ok := dataMap["tagName"].(string); ok {
			clickDetail.ElementType = tagName
		}

		// 提取内部文本
		if innerText, ok := dataMap["innerText"].(string); ok {
			clickDetail.InnerText = innerText
		}
	}

	// 保存点击详情
	db := model.GetDB()
	return db.Create(&clickDetail).Error
}

// 从SDK停留事件处理
func (s *EventService) processDwellEventFromSDK(req *TrackRequest, eventID uint) error {
	// 创建停留详情
	dwellDetail := model.DwellDetail{
		EventID: eventID,
	}

	// 从事件数据中提取停留信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// 提取页面URL
		if url, ok := dataMap["url"].(string); ok {
			dwellDetail.PageURL = url
		}

		// 提取页面标题
		if title, ok := dataMap["title"].(string); ok {
			dwellDetail.Title = title
		}

		// 提取停留时间
		if duration, ok := dataMap["duration"].(float64); ok {
			dwellDetail.StayTime = int64(duration)
		}
	}

	// 保存停留详情
	db := model.GetDB()
	return db.Create(&dwellDetail).Error
}

// 从SDK自定义事件处理
func (s *EventService) processCustomEventFromSDK(req *TrackRequest, eventID uint) error {
	// 创建自定义事件详情
	customDetail := model.CustomDetail{
		EventID: eventID,
	}

	// 从事件数据中提取自定义信息
	var dataMap map[string]interface{}
	if err := json.Unmarshal(req.Data, &dataMap); err == nil {
		// 提取事件名称
		if name, ok := dataMap["name"].(string); ok {
			customDetail.EventName = name
		}

		// 将整个数据保存为JSON字符串
		dataJSON, _ := json.Marshal(dataMap)
		customDetail.Data = string(dataJSON)
	}

	// 保存自定义事件详情
	db := model.GetDB()
	return db.Create(&customDetail).Error
}

// GetErrorList 获取错误列表
func (s *EventService) GetErrorList(projectIDStr, pageStr, pageSizeStr, startTimeStr, endTimeStr, errorType, severity string) (*ErrorListResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建查询条件
	query := model.GetDB().Model(&model.ErrorGroup{}).Where("project_id = ?", projectID)

	// 添加时间范围过滤
	if startTimeStr != "" {
		startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			query = query.Where("last_seen >= ?", startTime)
		}
	}

	if endTimeStr != "" {
		endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			query = query.Where("last_seen <= ?", endTime)
		}
	}

	// 添加错误类型过滤
	if errorType != "" {
		query = query.Where("error_type = ?", errorType)
	}

	// 添加严重程度过滤
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	// 获取错误分组列表
	var groups []model.ErrorGroup
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Order("last_seen DESC").Limit(pageSize).Offset(offset).Find(&groups).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]ErrorGroupItem, 0, len(groups))
	for _, group := range groups {
		list = append(list, ErrorGroupItem{
			ID:           group.ID,
			Fingerprint:  group.Fingerprint,
			ErrorType:    group.ErrorType,
			ErrorMessage: group.ErrorMessage,
			Count:        group.Count,
			FirstSeen:    group.FirstSeen,
			LastSeen:     group.LastSeen,
			Status:       group.Status,
			Severity:     group.Severity,
			SubType:      group.SubType,
		})
	}

	// 获取统计数据
	stats, err := s.getErrorStatsData(uint(projectID))
	if err != nil {
		return nil, err
	}

	return &ErrorListResponse{
		Total: total,
		List:  list,
		Stats: stats,
	}, nil
}

// GetErrorDetail 获取错误详情
func (s *EventService) GetErrorDetail(idStr string) (*ErrorDetailResponse, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的错误ID")
	}

	// 获取错误分组
	var group model.ErrorGroup
	if err := model.GetDB().First(&group, id).Error; err != nil {
		return nil, errors.New("错误不存在")
	}

	// 获取错误事件列表
	var errorDetails []model.ErrorDetail
	if err := model.GetDB().Where("fingerprint = ?", group.Fingerprint).
		Order("created_at DESC").
		Limit(10).
		Find(&errorDetails).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	groupItem := ErrorGroupItem{
		ID:           group.ID,
		Fingerprint:  group.Fingerprint,
		ErrorType:    group.ErrorType,
		ErrorMessage: group.ErrorMessage,
		Count:        group.Count,
		FirstSeen:    group.FirstSeen,
		LastSeen:     group.LastSeen,
		Status:       group.Status,
		Severity:     group.Severity,
		SubType:      group.SubType,
	}

	events := make([]ErrorEventItem, 0, len(errorDetails))
	for _, detail := range errorDetails {
		// 获取事件主信息
		var eventMain model.EventMain
		if err := model.GetDB().First(&eventMain, detail.EventID).Error; err != nil {
			continue
		}

		// 获取基础信息
		var baseInfo model.BaseInfo
		if err := model.GetDB().First(&baseInfo, eventMain.BaseInfoID).Error; err != nil {
			continue
		}

		events = append(events, ErrorEventItem{
			ID:           detail.ID,
			EventID:      eventMain.EventID,
			ErrorType:    detail.ErrorType,
			ErrorMessage: detail.ErrorMessage,
			ErrorStack:   detail.ErrorStack,
			FilePath:     detail.FilePath,
			LineNumber:   detail.LineNumber,
			ColumnNumber: detail.ColumnNumber,
			TriggerTime:  eventMain.TriggerTime,
			PageURL:      eventMain.TriggerPageURL,
			Browser:      baseInfo.Browser,
			OS:           baseInfo.OS,
			Device:       baseInfo.Device,
		})
	}

	return &ErrorDetailResponse{
		Group:  groupItem,
		Events: events,
		Total:  int64(len(events)),
	}, nil
}

// GetErrorStats 获取错误统计信息
func (s *EventService) GetErrorStats(projectIDStr, startTimeStr, endTimeStr string) (*ErrorStatsResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	// 获取统计数据
	stats, err := s.getErrorStatsData(uint(projectID))
	if err != nil {
		return nil, err
	}

	// 获取趋势数据
	trend, err := s.getErrorTrendData(uint(projectID), startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	return &ErrorStatsResponse{
		Stats: stats,
		Trend: trend,
	}, nil
}

// 获取错误统计数据
func (s *EventService) getErrorStatsData(projectID uint) (ErrorStatsData, error) {
	db := model.GetDB()
	var stats ErrorStatsData

	// 获取总错误数
	var totalErrors struct {
		TotalErrors sql.NullInt64 `gorm:"column:total_errors"`
	}
	if err := db.Model(&model.ErrorGroup{}).Where("project_id = ?", projectID).
		Select("COALESCE(SUM(count), 0) as total_errors").
		Scan(&totalErrors).Error; err != nil {
		return stats, err
	}

	if totalErrors.TotalErrors.Valid {
		stats.TotalErrors = totalErrors.TotalErrors.Int64
	} else {
		stats.TotalErrors = 0
	}

	// 获取受影响用户数
	var affectedUsers int64
	if err := db.Model(&model.BaseInfo{}).
		Joins("JOIN wt_event_main ON wt_event_main.base_info_id = wt_base_info.id").
		Joins("JOIN wt_error_detail ON wt_error_detail.event_id = wt_event_main.id").
		Where("wt_event_main.project_id = ?", projectID).
		Distinct("wt_base_info.user_uuid").
		Count(&affectedUsers).Error; err != nil {
		return stats, err
	}
	stats.AffectedUsers = affectedUsers

	// 获取今天的错误数
	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	if err := db.Model(&model.ErrorDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_error_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ?", projectID, todayStart).
		Count(&stats.ErrorsToday).Error; err != nil {
		return stats, err
	}

	// 获取昨天的错误数
	yesterdayStart := todayStart - 86400
	if err := db.Model(&model.ErrorDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_error_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ? AND wt_event_main.trigger_time < ?",
			projectID, yesterdayStart, todayStart).
		Count(&stats.ErrorsYesterday).Error; err != nil {
		return stats, err
	}

	// 获取错误类型分布
	var typeDistribution []struct {
		ErrorType string
		Count     int64
	}
	if err := db.Model(&model.ErrorDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_error_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("wt_error_detail.error_type, COUNT(*) as count").
		Group("wt_error_detail.error_type").
		Scan(&typeDistribution).Error; err != nil {
		return stats, err
	}

	stats.TypeDistribution = make(map[string]int64)
	for _, item := range typeDistribution {
		stats.TypeDistribution[item.ErrorType] = item.Count
	}

	// 获取浏览器分布
	var browserDistribution []struct {
		Browser string
		Count   int64
	}
	if err := db.Model(&model.BaseInfo{}).
		Joins("JOIN wt_event_main ON wt_event_main.base_info_id = wt_base_info.id").
		Joins("JOIN wt_error_detail ON wt_error_detail.event_id = wt_event_main.id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("wt_base_info.browser, COUNT(*) as count").
		Group("wt_base_info.browser").
		Scan(&browserDistribution).Error; err != nil {
		return stats, err
	}

	stats.BrowserDistribution = make(map[string]int64)
	for _, item := range browserDistribution {
		stats.BrowserDistribution[item.Browser] = item.Count
	}

	// 获取操作系统分布
	var osDistribution []struct {
		OS    string
		Count int64
	}
	if err := db.Model(&model.BaseInfo{}).
		Joins("JOIN wt_event_main ON wt_event_main.base_info_id = wt_base_info.id").
		Joins("JOIN wt_error_detail ON wt_error_detail.event_id = wt_event_main.id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("wt_base_info.os, COUNT(*) as count").
		Group("wt_base_info.os").
		Scan(&osDistribution).Error; err != nil {
		return stats, err
	}

	stats.OSDistribution = make(map[string]int64)
	for _, item := range osDistribution {
		stats.OSDistribution[item.OS] = item.Count
	}

	return stats, nil
}

// 获取错误趋势数据
func (s *EventService) getErrorTrendData(projectID uint, startTimeStr, endTimeStr string) ([]ErrorTrendItem, error) {
	db := model.GetDB()
	var trend []ErrorTrendItem

	// 设置默认时间范围为最近7天
	endTime := time.Now().Unix()
	startTime := endTime - 7*86400

	// 解析用户提供的时间范围
	if startTimeStr != "" {
		parsedStartTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			startTime = parsedStartTime
		}
	}

	if endTimeStr != "" {
		parsedEndTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			endTime = parsedEndTime
		}
	}

	// 按天分组查询错误数量
	var results []struct {
		Date  string
		Count int64
	}

	// 使用SQL函数将时间戳转换为日期字符串
	query := db.Model(&model.EventMain{}).
		Joins("JOIN wt_error_detail ON wt_error_detail.event_id = wt_event_main.id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ? AND wt_event_main.trigger_time <= ?",
			projectID, startTime, endTime).
		Select("DATE_FORMAT(FROM_UNIXTIME(wt_event_main.trigger_time), '%Y-%m-%d') as date, COUNT(*) as count").
		Group("date").
		Order("date")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	for _, result := range results {
		trend = append(trend, ErrorTrendItem{
			Date:  result.Date,
			Count: result.Count,
		})
	}

	return trend, nil
}

// GetPerformanceList 获取性能列表
func (s *EventService) GetPerformanceList(projectIDStr, pageStr, pageSizeStr, startTimeStr, endTimeStr, perfType string) (*PerformanceListResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建查询条件
	db := model.GetDB()
	query := db.Model(&model.PerformancePageDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_performance_page_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID)

	// 添加时间范围过滤
	if startTimeStr != "" {
		startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time >= ?", startTime)
		}
	}

	if endTimeStr != "" {
		endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time <= ?", endTime)
		}
	}

	// 获取性能数据列表
	// 使用中间结构体处理MySQL返回的浮点数字符串
	var perfDetails []struct {
		ID          uint
		EventID     string
		PageURL     string
		TriggerTime int64
		FP          sql.NullFloat64 `gorm:"column:fp"`
		FCP         sql.NullFloat64 `gorm:"column:fcp"`
		LCP         sql.NullFloat64 `gorm:"column:lcp"`
		FID         sql.NullFloat64 `gorm:"column:fid"`
		CLS         sql.NullFloat64 `gorm:"column:cls"`
		TTFB        sql.NullFloat64 `gorm:"column:ttfb"`
		DomReady    sql.NullFloat64 `gorm:"column:dom_ready"`
		Load        sql.NullFloat64 `gorm:"column:load"`
		Browser     string
		OS          string
		Device      string
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Select("wt_performance_page_detail.id, wt_event_main.event_id, wt_event_main.trigger_page_url as page_url, wt_event_main.trigger_time, " +
			"wt_performance_page_detail.fp, wt_performance_page_detail.fcp, wt_performance_page_detail.lcp, " +
			"0 as fid, " + // 使用固定值0替代不存在的字段
			"wt_performance_page_detail.cls, wt_performance_page_detail.ttfb, wt_performance_page_detail.dom_ready, wt_performance_page_detail.load, " +
			"wt_base_info.browser, wt_base_info.os, wt_base_info.device").
		Joins("JOIN wt_base_info ON wt_base_info.id = wt_event_main.base_info_id").
		Order("wt_event_main.trigger_time DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&perfDetails).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]PerformancePageItem, 0, len(perfDetails))
	for _, detail := range perfDetails {
		item := PerformancePageItem{
			ID:          detail.ID,
			EventID:     detail.EventID,
			PageURL:     detail.PageURL,
			TriggerTime: detail.TriggerTime,
			Browser:     detail.Browser,
			OS:          detail.OS,
			Device:      detail.Device,
		}

		// 将浮点数转换为整数
		if detail.FP.Valid {
			item.FP = int64(detail.FP.Float64)
		}
		if detail.FCP.Valid {
			item.FCP = int64(detail.FCP.Float64)
		}
		if detail.LCP.Valid {
			item.LCP = int64(detail.LCP.Float64)
		}
		if detail.FID.Valid {
			item.FID = int64(detail.FID.Float64)
		}
		if detail.CLS.Valid {
			item.CLS = detail.CLS.Float64
		}
		if detail.TTFB.Valid {
			item.TTFB = int64(detail.TTFB.Float64)
		}
		if detail.DomReady.Valid {
			item.DomReady = int64(detail.DomReady.Float64)
		}
		if detail.Load.Valid {
			item.Load = int64(detail.Load.Float64)
		}

		list = append(list, item)
	}

	// 获取统计数据
	stats, err := s.getPerformanceStatsData(uint(projectID))
	if err != nil {
		return nil, err
	}

	return &PerformanceListResponse{
		Total: total,
		List:  list,
		Stats: stats,
	}, nil
}

// GetPerformanceStats 获取性能统计信息
func (s *EventService) GetPerformanceStats(projectIDStr, startTimeStr, endTimeStr string) (*PerformanceStatsResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	// 获取统计数据
	stats, err := s.getPerformanceStatsData(uint(projectID))
	if err != nil {
		return nil, err
	}

	// 获取趋势数据
	trend, err := s.getPerformanceTrendData(uint(projectID), startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	return &PerformanceStatsResponse{
		Stats: stats,
		Trend: trend,
	}, nil
}

// 获取性能统计数据
func (s *EventService) getPerformanceStatsData(projectID uint) (PerformanceStatsData, error) {
	db := model.GetDB()
	var stats PerformanceStatsData

	// 使用中间结构体处理MySQL返回的浮点数字符串
	var rawStats struct {
		AvgFP       sql.NullFloat64 `gorm:"column:avg_fp"`
		AvgFCP      sql.NullFloat64 `gorm:"column:avg_fcp"`
		AvgLCP      sql.NullFloat64 `gorm:"column:avg_lcp"`
		AvgFID      sql.NullFloat64 `gorm:"column:avg_fid"`
		AvgCLS      sql.NullFloat64 `gorm:"column:avg_cls"`
		AvgTTFB     sql.NullFloat64 `gorm:"column:avg_ttfb"`
		AvgDomReady sql.NullFloat64 `gorm:"column:avg_dom_ready"`
		AvgLoad     sql.NullFloat64 `gorm:"column:avg_load"`
	}

	// 获取平均性能指标
	query := db.Model(&model.PerformancePageDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_performance_page_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("AVG(wt_performance_page_detail.fp) as avg_fp, " +
			"AVG(wt_performance_page_detail.fcp) as avg_fcp, " +
			"AVG(wt_performance_page_detail.lcp) as avg_lcp, " +
			"0 as avg_fid, " + // 使用固定值0替代不存在的字段
			"AVG(wt_performance_page_detail.cls) as avg_cls, " +
			"AVG(wt_performance_page_detail.ttfb) as avg_ttfb, " +
			"AVG(wt_performance_page_detail.dom_ready) as avg_dom_ready, " +
			"AVG(wt_performance_page_detail.load) as avg_load")

	if err := query.Scan(&rawStats).Error; err != nil {
		return stats, err
	}

	// 将浮点数转换为整数
	if rawStats.AvgFP.Valid {
		stats.AvgFP = int64(rawStats.AvgFP.Float64)
	}
	if rawStats.AvgFCP.Valid {
		stats.AvgFCP = int64(rawStats.AvgFCP.Float64)
	}
	if rawStats.AvgLCP.Valid {
		stats.AvgLCP = int64(rawStats.AvgLCP.Float64)
	}
	if rawStats.AvgFID.Valid {
		stats.AvgFID = int64(rawStats.AvgFID.Float64)
	}
	if rawStats.AvgCLS.Valid {
		stats.AvgCLS = rawStats.AvgCLS.Float64
	}
	if rawStats.AvgTTFB.Valid {
		stats.AvgTTFB = int64(rawStats.AvgTTFB.Float64)
	}
	if rawStats.AvgDomReady.Valid {
		stats.AvgDomReady = int64(rawStats.AvgDomReady.Float64)
	}
	if rawStats.AvgLoad.Valid {
		stats.AvgLoad = int64(rawStats.AvgLoad.Float64)
	}

	return stats, nil
}

// 获取性能趋势数据
func (s *EventService) getPerformanceTrendData(projectID uint, startTimeStr, endTimeStr string) ([]PerformanceTrendItem, error) {
	db := model.GetDB()
	var trend []PerformanceTrendItem

	// 设置默认时间范围为最近7天
	endTime := time.Now().Unix()
	startTime := endTime - 7*86400

	// 解析用户提供的时间范围
	if startTimeStr != "" {
		parsedStartTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			startTime = parsedStartTime
		}
	}

	if endTimeStr != "" {
		parsedEndTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			endTime = parsedEndTime
		}
	}

	// 按天分组查询性能指标
	// 使用中间结构体处理MySQL返回的浮点数字符串
	var rawResults []struct {
		Date string
		FP   sql.NullFloat64 `gorm:"column:fp"`
		FCP  sql.NullFloat64 `gorm:"column:fcp"`
		LCP  sql.NullFloat64 `gorm:"column:lcp"`
		TTFB sql.NullFloat64 `gorm:"column:ttfb"`
	}

	// 使用SQL函数将时间戳转换为日期字符串
	query := db.Model(&model.PerformancePageDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_performance_page_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ? AND wt_event_main.trigger_time <= ?",
			projectID, startTime, endTime).
		Select("DATE_FORMAT(FROM_UNIXTIME(wt_event_main.trigger_time), '%Y-%m-%d') as date, " +
			"AVG(wt_performance_page_detail.fp) as fp, " +
			"AVG(wt_performance_page_detail.fcp) as fcp, " +
			"AVG(wt_performance_page_detail.lcp) as lcp, " +
			"AVG(wt_performance_page_detail.ttfb) as ttfb").
		Group("date").
		Order("date")

	if err := query.Scan(&rawResults).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	for _, result := range rawResults {
		item := PerformanceTrendItem{
			Date: result.Date,
		}

		// 将浮点数转换为整数
		if result.FP.Valid {
			item.FP = int64(result.FP.Float64)
		}
		if result.FCP.Valid {
			item.FCP = int64(result.FCP.Float64)
		}
		if result.LCP.Valid {
			item.LCP = int64(result.LCP.Float64)
		}
		if result.TTFB.Valid {
			item.TTFB = int64(result.TTFB.Float64)
		}

		trend = append(trend, item)
	}

	return trend, nil
}

// GetResourcePerformanceList 获取资源性能列表
func (s *EventService) GetResourcePerformanceList(projectIDStr, pageStr, pageSizeStr, startTimeStr, endTimeStr, resourceType string) (*ResourcePerformanceListResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建查询条件
	db := model.GetDB()
	query := db.Model(&model.PerformanceResourceDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_performance_resource_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID)

	// 添加时间范围过滤
	if startTimeStr != "" {
		startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time >= ?", startTime)
		}
	}

	if endTimeStr != "" {
		endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time <= ?", endTime)
		}
	}

	// 添加资源类型过滤
	if resourceType != "" {
		query = query.Where("wt_performance_resource_detail.resource_type = ?", resourceType)
	}

	// 获取资源性能数据列表
	var resourceDetails []struct {
		ID            uint
		EventID       string
		ResourceURL   string
		ResourceType  string
		InitiatorType string
		StartTime     int64
		Duration      int64
		TransferSize  int64
		PageURL       string
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Select("wt_performance_resource_detail.id, wt_event_main.event_id, wt_performance_resource_detail.resource_url, " +
			"wt_performance_resource_detail.resource_type, wt_performance_resource_detail.initiator_type, " +
			"wt_performance_resource_detail.start_time, wt_performance_resource_detail.duration, " +
			"wt_performance_resource_detail.transfer_size, wt_event_main.trigger_page_url as page_url").
		Order("wt_event_main.trigger_time DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&resourceDetails).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]ResourcePerformanceItem, 0, len(resourceDetails))
	for _, detail := range resourceDetails {
		list = append(list, ResourcePerformanceItem{
			ID:            detail.ID,
			EventID:       detail.EventID,
			ResourceURL:   detail.ResourceURL,
			ResourceType:  detail.ResourceType,
			InitiatorType: detail.InitiatorType,
			StartTime:     detail.StartTime,
			Duration:      detail.Duration,
			TransferSize:  detail.TransferSize,
			PageURL:       detail.PageURL,
		})
	}

	return &ResourcePerformanceListResponse{
		Total: total,
		List:  list,
	}, nil
}

// GetPageViewList 获取页面访问列表
func (s *EventService) GetPageViewList(projectIDStr, pageStr, pageSizeStr, startTimeStr, endTimeStr string) (*PVListResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建查询条件
	db := model.GetDB()
	query := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID)

	// 添加时间范围过滤
	if startTimeStr != "" {
		startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time >= ?", startTime)
		}
	}

	if endTimeStr != "" {
		endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time <= ?", endTime)
		}
	}

	// 获取页面访问数据列表
	var pvDetails []struct {
		ID          uint
		EventID     string
		PageURL     string
		Title       string
		Referrer    string
		TriggerTime int64
		StayTime    int64
		IsNewVisit  bool
		Browser     string
		OS          string
		Device      string
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Select("wt_pv_detail.id, wt_event_main.event_id, wt_pv_detail.page_url, wt_pv_detail.title, " +
			"wt_pv_detail.referrer, wt_event_main.trigger_time, wt_pv_detail.stay_time, wt_pv_detail.is_new_visit, " +
			"wt_base_info.browser, wt_base_info.os, wt_base_info.device").
		Joins("JOIN wt_base_info ON wt_base_info.id = wt_event_main.base_info_id").
		Order("wt_event_main.trigger_time DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&pvDetails).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]PVItem, 0, len(pvDetails))
	for _, detail := range pvDetails {
		list = append(list, PVItem{
			ID:          detail.ID,
			EventID:     detail.EventID,
			PageURL:     detail.PageURL,
			Title:       detail.Title,
			Referrer:    detail.Referrer,
			TriggerTime: detail.TriggerTime,
			StayTime:    detail.StayTime,
			IsNewVisit:  detail.IsNewVisit,
			Browser:     detail.Browser,
			OS:          detail.OS,
			Device:      detail.Device,
		})
	}

	return &PVListResponse{
		Total: total,
		List:  list,
	}, nil
}

// GetClickList 获取点击列表
func (s *EventService) GetClickList(projectIDStr, pageStr, pageSizeStr, startTimeStr, endTimeStr string) (*ClickListResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建查询条件
	db := model.GetDB()
	query := db.Model(&model.ClickDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_click_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID)

	// 添加时间范围过滤
	if startTimeStr != "" {
		startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time >= ?", startTime)
		}
	}

	if endTimeStr != "" {
		endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			query = query.Where("wt_event_main.trigger_time <= ?", endTime)
		}
	}

	// 获取点击数据列表
	var clickDetails []struct {
		ID          uint
		EventID     string
		ElementPath string
		ElementType string
		InnerText   string
		TriggerTime int64
		PageURL     string
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Select("wt_click_detail.id, wt_event_main.event_id, wt_click_detail.element_path, " +
			"wt_click_detail.element_type, wt_click_detail.inner_text, wt_event_main.trigger_time, " +
			"wt_event_main.trigger_page_url as page_url").
		Order("wt_event_main.trigger_time DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&clickDetails).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]ClickItem, 0, len(clickDetails))
	for _, detail := range clickDetails {
		list = append(list, ClickItem{
			ID:          detail.ID,
			EventID:     detail.EventID,
			ElementPath: detail.ElementPath,
			ElementType: detail.ElementType,
			InnerText:   detail.InnerText,
			TriggerTime: detail.TriggerTime,
			PageURL:     detail.PageURL,
		})
	}

	return &ClickListResponse{
		Total: total,
		List:  list,
	}, nil
}

// GetBehaviorStats 获取用户行为统计信息
func (s *EventService) GetBehaviorStats(projectIDStr, startTimeStr, endTimeStr string) (*BehaviorStatsResponse, error) {
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	// 获取PV统计数据
	pvStats, err := s.getPVStatsData(uint(projectID))
	if err != nil {
		return nil, err
	}

	// 获取点击统计数据
	clickStats, err := s.getClickStatsData(uint(projectID))
	if err != nil {
		return nil, err
	}

	// 获取PV趋势数据
	pvTrend, err := s.getPVTrendData(uint(projectID), startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	return &BehaviorStatsResponse{
		PVStats:    pvStats,
		ClickStats: clickStats,
		PVTrend:    pvTrend,
	}, nil
}

// 获取PV统计数据
func (s *EventService) getPVStatsData(projectID uint) (PVStatsData, error) {
	db := model.GetDB()
	var stats PVStatsData

	// 获取总PV数
	if err := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Count(&stats.TotalPV).Error; err != nil {
		return stats, err
	}

	// 获取总UV数
	if err := db.Model(&model.BaseInfo{}).
		Joins("JOIN wt_event_main ON wt_event_main.base_info_id = wt_base_info.id").
		Joins("JOIN wt_pv_detail ON wt_pv_detail.event_id = wt_event_main.id").
		Where("wt_event_main.project_id = ?", projectID).
		Distinct("wt_base_info.user_uuid").
		Count(&stats.TotalUV).Error; err != nil {
		return stats, err
	}

	// 获取今天的PV数
	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	if err := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ?", projectID, todayStart).
		Count(&stats.PVToday).Error; err != nil {
		return stats, err
	}

	// 获取今天的UV数
	if err := db.Model(&model.BaseInfo{}).
		Joins("JOIN wt_event_main ON wt_event_main.base_info_id = wt_base_info.id").
		Joins("JOIN wt_pv_detail ON wt_pv_detail.event_id = wt_event_main.id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ?", projectID, todayStart).
		Distinct("wt_base_info.user_uuid").
		Count(&stats.UVToday).Error; err != nil {
		return stats, err
	}

	// 获取平均停留时间
	var avgStayTime struct {
		AvgStayTime int64
	}
	if err := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("AVG(wt_pv_detail.stay_time) as avg_stay_time").
		Scan(&avgStayTime).Error; err != nil {
		return stats, err
	}
	stats.AvgStayTime = avgStayTime.AvgStayTime

	// 获取跳出率
	var bounceCount int64
	if err := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_pv_detail.stay_time < 10", projectID).
		Count(&bounceCount).Error; err != nil {
		return stats, err
	}
	if stats.TotalPV > 0 {
		stats.BounceRate = float64(bounceCount) / float64(stats.TotalPV) * 100
	}

	// 获取热门页面
	var topPages []struct {
		PageURL string
		Count   int64
	}
	if err := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("wt_pv_detail.page_url, COUNT(*) as count").
		Group("wt_pv_detail.page_url").
		Order("count DESC").
		Limit(10).
		Scan(&topPages).Error; err != nil {
		return stats, err
	}

	stats.TopPages = make(map[string]int64)
	for _, page := range topPages {
		stats.TopPages[page.PageURL] = page.Count
	}

	return stats, nil
}

// 获取点击统计数据
func (s *EventService) getClickStatsData(projectID uint) (ClickStatsData, error) {
	db := model.GetDB()
	var stats ClickStatsData

	// 获取总点击数
	if err := db.Model(&model.ClickDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_click_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Count(&stats.TotalClicks).Error; err != nil {
		return stats, err
	}

	// 获取今天的点击数
	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	if err := db.Model(&model.ClickDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_click_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ?", projectID, todayStart).
		Count(&stats.ClicksToday).Error; err != nil {
		return stats, err
	}

	// 获取热门元素
	var topElements []struct {
		ElementPath string
		Count       int64
	}
	if err := db.Model(&model.ClickDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_click_detail.event_id").
		Where("wt_event_main.project_id = ?", projectID).
		Select("wt_click_detail.element_path, COUNT(*) as count").
		Group("wt_click_detail.element_path").
		Order("count DESC").
		Limit(10).
		Scan(&topElements).Error; err != nil {
		return stats, err
	}

	stats.TopElements = make(map[string]int64)
	for _, element := range topElements {
		stats.TopElements[element.ElementPath] = element.Count
	}

	return stats, nil
}

// 获取PV趋势数据
func (s *EventService) getPVTrendData(projectID uint, startTimeStr, endTimeStr string) ([]PVTrendItem, error) {
	db := model.GetDB()
	var trend []PVTrendItem

	// 设置默认时间范围为最近7天
	endTime := time.Now().Unix()
	startTime := endTime - 7*86400

	// 解析用户提供的时间范围
	if startTimeStr != "" {
		parsedStartTime, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err == nil {
			startTime = parsedStartTime
		}
	}

	if endTimeStr != "" {
		parsedEndTime, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err == nil {
			endTime = parsedEndTime
		}
	}

	// 按天分组查询PV和UV
	var results []struct {
		Date string
		PV   int64
		UV   int64
	}

	// 获取PV数据
	pvQuery := db.Model(&model.PVDetail{}).
		Joins("JOIN wt_event_main ON wt_event_main.id = wt_pv_detail.event_id").
		Where("wt_event_main.project_id = ? AND wt_event_main.trigger_time >= ? AND wt_event_main.trigger_time <= ?",
			projectID, startTime, endTime).
		Select("DATE_FORMAT(FROM_UNIXTIME(wt_event_main.trigger_time), '%Y-%m-%d') as date, COUNT(*) as pv").
		Group("date").
		Order("date")

	if err := pvQuery.Scan(&results).Error; err != nil {
		return nil, err
	}

	// 获取UV数据
	for i, result := range results {
		var uv int64
		uvQuery := db.Model(&model.BaseInfo{}).
			Joins("JOIN wt_event_main ON wt_event_main.base_info_id = wt_base_info.id").
			Joins("JOIN wt_pv_detail ON wt_pv_detail.event_id = wt_event_main.id").
			Where("wt_event_main.project_id = ? AND DATE_FORMAT(FROM_UNIXTIME(wt_event_main.trigger_time), '%Y-%m-%d') = ?",
				projectID, result.Date).
			Distinct("wt_base_info.user_uuid").
			Count(&uv)

		if uvQuery.Error != nil {
			return nil, uvQuery.Error
		}

		results[i].UV = uv
	}

	// 转换为响应格式
	for _, result := range results {
		trend = append(trend, PVTrendItem{
			Date: result.Date,
			PV:   result.PV,
			UV:   result.UV,
		})
	}

	return trend, nil
}
