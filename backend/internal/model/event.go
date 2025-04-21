package model

import (
	"gorm.io/datatypes"
)

// BaseInfo holds the common client-side metadata.
type BaseInfo struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	ClientHeight int            `gorm:"not null"`         // client viewport height
	ClientWidth  int            `gorm:"not null"`         // client viewport width
	ColorDepth   int            `gorm:"not null"`         // screen color depth
	PixelDepth   int            `gorm:"not null"`         // screen pixel depth
	DeviceID     string         `gorm:"size:64;not null"` // unique device identifier
	ScreenWidth  int            `gorm:"not null"`         // physical screen width
	ScreenHeight int            `gorm:"not null"`         // physical screen height
	Vendor       string         `gorm:"size:64;not null"` // browser vendor
	Platform     string         `gorm:"size:64;not null"` // operating system/platform
	UserUUID     string         `gorm:"size:64;not null"` // user UUID from app
	SDKUserUUID  string         `gorm:"size:64;not null"` // SDK‑generated user UUID
	Ext          datatypes.JSON `gorm:"type:jsonb"`       // any extra fields
	AppName      string         `gorm:"size:64;not null"` // application name
	AppCode      string         `gorm:"size:64"`          // optional app code
	PageID       string         `gorm:"size:64;not null"` // page identifier
	SessionID    string         `gorm:"size:64;not null"` // session identifier
	SDKVersion   string         `gorm:"size:32;not null"` // SDK version
	IP           string         `gorm:"size:45"`          // client IP address
	SendTime     int64          `gorm:"not null"`         // original JS timestamp (ms)
	CreatedAt    int64          // record creation time
}

// EventMain holds the common fields for all events.
type EventMain struct {
	ID             uint     `gorm:"primaryKey;autoIncrement"`
	BaseInfoID     uint     `gorm:"not null;index"` // foreign key to BaseInfo
	BaseInfo       BaseInfo `gorm:"constraint:OnDelete:CASCADE"`
	EventType      string   `gorm:"size:32;not null"`   // pv, error, performance, etc.
	EventID        string   `gorm:"size:128;not null"`  // SDK unique event id
	TriggerPageURL string   `gorm:"type:text;not null"` // location.href at trigger time
	SendTime       int64    `gorm:"not null"`           // original JS timestamp (ms)
	CreatedAt      int64    // record creation time
}

// performance_page_detail stores page‐level metrics.
type PerformancePageDetail struct {
	MainID    uint      `gorm:"primaryKey"` // references EventMain.ID
	Main      EventMain `gorm:"constraint:OnDelete:CASCADE"`
	TTI       float64   // Time To Interactive (ms)
	Ready     float64   // Document ready event time (ms)
	LoadOn    float64   // Load event time (ms)
	FirstByte float64   // Time to first byte (ms)
	TTFB      float64   // Time to first paint (ms)
	Trans     float64   // TCP transfer time (ms)
	DOM       float64   // DOM parsing time (ms)
	Res       float64   // Resource load time (ms)
	SSLLINK   float64   // SSL handshake time (ms)
}

// performance_resource_detail stores resource‐level metrics.
type PerformanceResourceDetail struct {
	MainID            uint      `gorm:"primaryKey"`
	Main              EventMain `gorm:"constraint:OnDelete:CASCADE"`
	InitiatorType     string    `gorm:"size:32"` // resource initiatorType
	TransferSize      int       // transferSize (bytes)
	EncodedBodySize   int       // encodedBodySize (bytes)
	DecodedBodySize   int       // decodedBodySize (bytes)
	Duration          float64   // total fetch duration (ms)
	StartTime         float64   // startTime (ms)
	FetchStart        float64   // fetchStart (ms)
	DomainLookupStart float64   // domainLookupStart (ms)
	DomainLookupEnd   float64   // domainLookupEnd (ms)
	ConnectStart      float64   // connectStart (ms)
	ConnectEnd        float64   // connectEnd (ms)
	RequestStart      float64   // requestStart (ms)
	ResponseStart     float64   // responseStart (ms)
	ResponseEnd       float64   // responseEnd (ms)
	RequestURL        string    `gorm:"type:text"` // resource URL
}

// pv_detail stores page view events.
type PVDetail struct {
	MainID      uint      `gorm:"primaryKey"`
	Main        EventMain `gorm:"constraint:OnDelete:CASCADE"`
	Referer     string    `gorm:"type:text"` // document.referrer
	Title       string    `gorm:"type:text"` // document.title
	Action      string    `gorm:"size:32"`   // custom action e.g. navigation
	TriggerTime int64     // JS timestamp of the PV event
}

// click_detail stores click events.
type ClickDetail struct {
	MainID      uint      `gorm:"primaryKey"`
	Main        EventMain `gorm:"constraint:OnDelete:CASCADE"`
	ElementID   string    `gorm:"type:text"` // clicked element selector or id
	TriggerTime int64     // JS timestamp of click
}

// dwell_detail stores page unload (dwell) events.
type DwellDetail struct {
	MainID       uint      `gorm:"primaryKey"`
	Main         EventMain `gorm:"constraint:OnDelete:CASCADE"`
	StayDuration int       // time on page (ms)
	TriggerTime  int64     // JS timestamp of unload
}

// custom_detail stores custom events.
type CustomDetail struct {
	MainID      uint           `gorm:"primaryKey"`
	Main        EventMain      `gorm:"constraint:OnDelete:CASCADE"`
	Data        datatypes.JSON `gorm:"type:jsonb;not null"` // custom payload from ext
	TriggerTime int64          // JS timestamp of custom event
}

// intersection_detail stores exposure events.
type IntersectionDetail struct {
	MainID      uint      `gorm:"primaryKey"`
	Main        EventMain `gorm:"constraint:OnDelete:CASCADE"`
	ElementID   string    `gorm:"type:text"`         // observed element selector or id
	Ratio       float32   `gorm:"type:numeric(5,4)"` // intersectionRatio
	TriggerTime int64     // JS timestamp of intersection
}
