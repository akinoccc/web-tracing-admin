package model

import (
	"time"
)

// 错误事件详情
type ErrorDetail struct {
	Model
	EventID       uint       `json:"eventId" gorm:"not null"`
	Event         *EventMain `json:"event" gorm:"foreignKey:EventID"`
	ErrorType     string     `json:"errorType" gorm:"size:50;not null"`
	ErrorMessage  string     `json:"errorMessage" gorm:"type:text;not null"`
	ErrorStack    string     `json:"errorStack" gorm:"type:text"`
	FilePath      string     `json:"filePath" gorm:"type:text"`
	LineNumber    int        `json:"lineNumber"`
	ColumnNumber  int        `json:"columnNumber"`
	ComponentName string     `json:"componentName" gorm:"size:100"`
	RecordScreen  string     `json:"recordScreen" gorm:"type:text"`
	// 扩展字段
	Severity    string `json:"severity" gorm:"size:20"`
	Fingerprint string `json:"fingerprint"`
	SubType     string `json:"subType" gorm:"size:50"`
	Context     string `json:"context" gorm:"type:text"`
}

// HTTP 请求错误详情
type HttpErrorDetail struct {
	Model
	EventID      uint       `json:"eventId" gorm:"not null"`
	Event        *EventMain `json:"event" gorm:"foreignKey:EventID"`
	URL          string     `json:"url" gorm:"type:text;not null"`
	Method       string     `json:"method" gorm:"size:20;not null"`
	Status       int        `json:"status" gorm:"not null"`
	StatusText   string     `json:"statusText" gorm:"size:100"`
	RequestData  string     `json:"requestData" gorm:"type:text"`
	ResponseData string     `json:"responseData" gorm:"type:text"`
	Duration     int64      `json:"duration"`
	ErrorType    string     `json:"errorType" gorm:"size:50"`
	ErrorMessage string     `json:"errorMessage" gorm:"type:text"`
}

// 资源错误详情
type ResourceErrorDetail struct {
	Model
	EventID      uint       `json:"eventId" gorm:"not null"`
	Event        *EventMain `json:"event" gorm:"foreignKey:EventID"`
	ResourceURL  string     `json:"resourceUrl" gorm:"type:text;not null"`
	ResourceType string     `json:"resourceType" gorm:"size:50;not null"`
	ErrorType    string     `json:"errorType" gorm:"size:50"`
	ErrorMessage string     `json:"errorMessage" gorm:"type:text"`
	ElementType  string     `json:"elementType" gorm:"size:50"`
}

// Vue 错误详情
type VueErrorDetail struct {
	Model
	EventID       uint       `json:"eventId" gorm:"not null"`
	Event         *EventMain `json:"event" gorm:"foreignKey:EventID"`
	ComponentName string     `json:"componentName" gorm:"size:100"`
	PropsData     string     `json:"propsData" gorm:"type:text"`
	ErrorType     string     `json:"errorType" gorm:"size:50;not null"`
	ErrorMessage  string     `json:"errorMessage" gorm:"type:text;not null"`
	ErrorStack    string     `json:"errorStack" gorm:"type:text"`
	Info          string     `json:"info" gorm:"type:text"`
}

// React 错误详情
type ReactErrorDetail struct {
	Model
	EventID        uint       `json:"eventId" gorm:"not null"`
	Event          *EventMain `json:"event" gorm:"foreignKey:EventID"`
	ComponentName  string     `json:"componentName" gorm:"size:100"`
	ComponentStack string     `json:"componentStack" gorm:"type:text"`
	ErrorType      string     `json:"errorType" gorm:"size:50;not null"`
	ErrorMessage   string     `json:"errorMessage" gorm:"type:text;not null"`
	ErrorStack     string     `json:"errorStack" gorm:"type:text"`
}

// 错误分组
type ErrorGroup struct {
	Model
	Fingerprint   string  `json:"fingerprint" gorm:"size:100;not null;unique"`
	ErrorType     string  `json:"errorType" gorm:"size:50;not null"`
	ErrorMessage  string  `json:"errorMessage" gorm:"type:text;not null"`
	Count         int     `json:"count" gorm:"not null;default:1"`
	FirstSeen     int64   `json:"firstSeen" gorm:"not null"`
	LastSeen      int64   `json:"lastSeen" gorm:"not null"`
	ProjectID     uint    `json:"projectId" gorm:"not null"`
	Project       Project `json:"project" gorm:"foreignKey:ProjectID"`
	SampleEventID uint    `json:"sampleEventId"`
	Status        string  `json:"status" gorm:"size:20;default:'active'"`
	Severity      string  `json:"severity" gorm:"size:20"`
	SubType       string  `json:"subType" gorm:"size:50"`
}

// 创建错误详情
func CreateErrorDetail(detail *ErrorDetail) error {
	return db.Create(detail).Error
}

// 创建 HTTP 错误详情
func CreateHttpErrorDetail(detail *HttpErrorDetail) error {
	return db.Create(detail).Error
}

// 创建资源错误详情
func CreateResourceErrorDetail(detail *ResourceErrorDetail) error {
	return db.Create(detail).Error
}

// 创建 Vue 错误详情
func CreateVueErrorDetail(detail *VueErrorDetail) error {
	return db.Create(detail).Error
}

// 创建 React 错误详情
func CreateReactErrorDetail(detail *ReactErrorDetail) error {
	return db.Create(detail).Error
}

// 创建或更新错误分组
func CreateOrUpdateErrorGroup(fingerprint, errorType, errorMessage string, projectID uint, eventID uint, severity, subType string) (*ErrorGroup, error) {
	var group ErrorGroup
	now := time.Now().Unix()

	// 查找是否已存在相同指纹的错误分组
	err := db.Where("fingerprint = ? AND project_id = ?", fingerprint, projectID).First(&group).Error
	if err != nil {
		// 不存在，创建新分组
		group = ErrorGroup{
			Fingerprint:   fingerprint,
			ErrorType:     errorType,
			ErrorMessage:  errorMessage,
			Count:         1,
			FirstSeen:     now,
			LastSeen:      now,
			ProjectID:     projectID,
			SampleEventID: eventID,
			Status:        "active",
			Severity:      severity,
			SubType:       subType,
		}
		return &group, db.Create(&group).Error
	}

	// 已存在，更新分组
	group.Count++
	group.LastSeen = now
	// 每隔一定次数更新示例事件，以获取最新的上下文
	if group.Count%10 == 0 {
		group.SampleEventID = eventID
	}

	return &group, db.Save(&group).Error
}

// 获取项目的错误分组列表
func GetErrorGroupsByProjectID(projectID uint, limit, offset int) ([]ErrorGroup, int64, error) {
	var groups []ErrorGroup
	var total int64

	// 获取总数
	if err := db.Model(&ErrorGroup{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := db.Where("project_id = ?", projectID).Order("last_seen DESC").Limit(limit).Offset(offset).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// 获取错误分组详情
func GetErrorGroupByID(id uint) (*ErrorGroup, error) {
	var group ErrorGroup
	if err := db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// 获取错误分组的事件列表
func GetErrorEventsByGroupID(groupID uint, limit, offset int) ([]ErrorDetail, int64, error) {
	var group ErrorGroup
	if err := db.First(&group, groupID).Error; err != nil {
		return nil, 0, err
	}

	var events []ErrorDetail
	var total int64

	// 获取总数
	if err := db.Model(&ErrorDetail{}).Where("fingerprint = ?", group.Fingerprint).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := db.Where("fingerprint = ?", group.Fingerprint).Order("created_at DESC").Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
