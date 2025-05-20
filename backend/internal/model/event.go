package model

// 事件类型枚举
const (
	EventTypeError               = "error"
	EventTypePerformancePage     = "performance_page"
	EventTypePerformanceResource = "performance_resource"
	EventTypePV                  = "pv"
	EventTypeClick               = "click"
	EventTypeDwell               = "dwell"
	EventTypeIntersection        = "intersection"
	EventTypeCustom              = "custom"
)

// 错误类型枚举
const (
	ErrorTypeJS       = "js_error"
	ErrorTypePromise  = "promise_rejection"
	ErrorTypeResource = "resource_error"
	ErrorTypeHttp     = "http_error"
	ErrorTypeVue      = "vue_error"
	ErrorTypeReact    = "react_error"
	ErrorTypeCustom   = "custom_error"
)

// 错误子类型枚举
const (
	ErrorSubTypeSyntax       = "syntax_error"
	ErrorSubTypeReference    = "reference_error"
	ErrorSubTypeType         = "type_error"
	ErrorSubTypeRange        = "range_error"
	ErrorSubTypeBadRequest   = "bad_request"
	ErrorSubTypeUnauthorized = "unauthorized"
	ErrorSubTypeForbidden    = "forbidden"
	ErrorSubTypeNotFound     = "not_found"
	ErrorSubTypeServerError  = "server_error"
	ErrorSubTypeTimeout      = "timeout"
	ErrorSubTypeNetworkError = "network_error"
	ErrorSubTypeComponent    = "component_error"
	ErrorSubTypeRender       = "render_error"
	ErrorSubTypeLifecycle    = "lifecycle_error"
	ErrorSubTypeRouter       = "router_error"
	ErrorSubTypeStore        = "store_error"
)

// 错误严重程度枚举
const (
	ErrorSeverityFatal   = "fatal"
	ErrorSeverityError   = "error"
	ErrorSeverityWarning = "warning"
	ErrorSeverityInfo    = "info"
)

// 事件基础信息
type BaseInfo struct {
	Model
	ProjectID      uint   `json:"projectId" gorm:"not null"`
	AppKey         string `json:"appKey" gorm:"size:50;not null"`
	UserID         string `json:"userId" gorm:"size:100"`
	UserUUID       string `json:"userUuid" gorm:"size:100"`
	SessionID      string `json:"sessionId" gorm:"size:100"`
	PageURL        string `json:"pageUrl" gorm:"type:text"`
	Referrer       string `json:"referrer" gorm:"type:text"`
	UserAgent      string `json:"userAgent" gorm:"type:text"`
	IP             string `json:"ip" gorm:"size:50"`
	Browser        string `json:"browser" gorm:"size:50"`
	BrowserVersion string `json:"browserVersion" gorm:"size:50"`
	OS             string `json:"os" gorm:"size:50"`
	OSVersion      string `json:"osVersion" gorm:"size:50"`
	Device         string `json:"device" gorm:"size:50"`
	DeviceType     string `json:"deviceType" gorm:"size:50"`
	Vendor         string `json:"vendor" gorm:"size:50"`
	// 扩展字段
	SDKVersion   string `json:"sdkVersion" gorm:"size:50"`
	SDKUserUUID  string `json:"sdkUserUuid" gorm:"size:100"`
	AppName      string `json:"appName" gorm:"size:100"`
	AppCode      string `json:"appCode" gorm:"size:50"`
	Platform     string `json:"platform" gorm:"size:50"`
	ScreenWidth  int    `json:"screenWidth"`
	ScreenHeight int    `json:"screenHeight"`
	ClientWidth  int    `json:"clientWidth"`
	ClientHeight int    `json:"clientHeight"`
	ColorDepth   int    `json:"colorDepth"`
	PixelDepth   int    `json:"pixelDepth"`
	DeviceID     string `json:"deviceId" gorm:"size:100"`
	PageID       string `json:"pageId" gorm:"size:100"`
	SendTime     int64  `json:"sendTime"`
	Ext          string `json:"ext" gorm:"type:text"`
}

// 事件主表
type EventMain struct {
	Model
	EventID        string    `json:"eventId" gorm:"size:100;not null;unique"`
	EventType      string    `json:"eventType" gorm:"size:50;not null"`
	ProjectID      uint      `json:"projectId" gorm:"not null"`
	BaseInfoID     uint      `json:"baseInfoId" gorm:"not null"`
	BaseInfo       *BaseInfo `json:"baseInfo" gorm:"foreignKey:BaseInfoID"`
	TriggerTime    int64     `json:"triggerTime" gorm:"not null"`
	SendTime       int64     `json:"sendTime" gorm:"not null"`
	TriggerPageURL string    `json:"triggerPageUrl" gorm:"type:text"`
	Title          string    `json:"title" gorm:"size:255"`
	Referer        string    `json:"referer" gorm:"type:text"`
}

// 性能页面详情
type PerformancePageDetail struct {
	Model
	EventID         uint       `json:"eventId" gorm:"not null"`
	Event           *EventMain `json:"event" gorm:"foreignKey:EventID"`
	FP              int64      `json:"fp"`
	FCP             int64      `json:"fcp"`
	LCP             int64      `json:"lcp"`
	FID             int64      `json:"fid"`
	CLS             float64    `json:"cls"`
	TTFB            int64      `json:"ttfb"`
	DomReady        int64      `json:"domReady"`
	Load            int64      `json:"load"`
	FirstByte       int64      `json:"firstByte"`
	DNS             int64      `json:"dns"`
	TCP             int64      `json:"tcp"`
	SSL             int64      `json:"ssl"`
	TTFB2           int64      `json:"ttfb2"`
	Trans           int64      `json:"trans"`
	DomParse        int64      `json:"domParse"`
	ResourceLoad    int64      `json:"resourceLoad"`
	DomContentLoad  int64      `json:"domContentLoad"`
	FirstScreenTime int64      `json:"firstScreenTime"`
}

// 性能资源详情
type PerformanceResourceDetail struct {
	Model
	EventID         uint       `json:"eventId" gorm:"not null"`
	Event           *EventMain `json:"event" gorm:"foreignKey:EventID"`
	InitiatorType   string     `json:"initiatorType" gorm:"size:50;not null"`
	ResourceType    string     `json:"resourceType" gorm:"size:50;not null"`
	ResourceURL     string     `json:"resourceUrl" gorm:"type:text;not null"`
	ResponseStatus  string     `json:"responseStatus" gorm:"size:50"`
	StartTime       int64      `json:"startTime" gorm:"not null"`
	Duration        int64      `json:"duration" gorm:"not null"`
	TransferSize    int64      `json:"transferSize"`
	EncodedBodySize int64      `json:"encodedBodySize"`
	DecodedBodySize int64      `json:"decodedBodySize"`
	DNSTime         int64      `json:"dnsTime"`
	TCPTime         int64      `json:"tcpTime"`
	SSLTime         int64      `json:"sslTime"`
	TTFB            int64      `json:"ttfb"`
	DownloadTime    int64      `json:"downloadTime"`
	FromCache       bool       `json:"fromCache"`
}

// 页面浏览详情
type PVDetail struct {
	Model
	EventID      uint       `json:"eventId" gorm:"not null"`
	Event        *EventMain `json:"event" gorm:"foreignKey:EventID"`
	PageURL      string     `json:"pageUrl" gorm:"type:text;not null"`
	Title        string     `json:"title" gorm:"size:255"`
	Referrer     string     `json:"referrer" gorm:"type:text"`
	StayTime     int64      `json:"stayTime"`
	IsNewVisit   bool       `json:"isNewVisit"`
	IsNewSession bool       `json:"isNewSession"`
}

// 点击详情
type ClickDetail struct {
	Model
	EventID     uint       `json:"eventId" gorm:"not null"`
	Event       *EventMain `json:"event" gorm:"foreignKey:EventID"`
	ElementPath string     `json:"elementPath" gorm:"type:text"`
	ElementType string     `json:"elementType" gorm:"size:50"`
	InnerText   string     `json:"innerText" gorm:"type:text"`
}

// 停留详情
type DwellDetail struct {
	Model
	EventID  uint       `json:"eventId" gorm:"not null"`
	Event    *EventMain `json:"event" gorm:"foreignKey:EventID"`
	PageURL  string     `json:"pageUrl" gorm:"type:text;not null"`
	Title    string     `json:"title" gorm:"size:255"`
	StayTime int64      `json:"stayTime" gorm:"not null"`
}

// 曝光详情
type IntersectionDetail struct {
	Model
	EventID     uint       `json:"eventId" gorm:"not null"`
	Event       *EventMain `json:"event" gorm:"foreignKey:EventID"`
	ElementPath string     `json:"elementPath" gorm:"type:text"`
	ElementType string     `json:"elementType" gorm:"size:50"`
	InnerText   string     `json:"innerText" gorm:"type:text"`
}

// 自定义事件详情
type CustomDetail struct {
	Model
	EventID     uint       `json:"eventId" gorm:"not null"`
	Event       *EventMain `json:"event" gorm:"foreignKey:EventID"`
	EventName   string     `json:"eventName" gorm:"size:100;not null"`
	EventParams string     `json:"eventParams" gorm:"type:text"`
	Data        string     `json:"data" gorm:"type:text"`
}
