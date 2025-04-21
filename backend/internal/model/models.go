package model

import (
	"fmt"
	"log"
	"time"

	"github.com/go-ini/ini"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// Model 基础模型
type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 数据库配置
type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

// 服务器配置
type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	JwtSecret    string
}

var DatabaseSetting = &Database{}
var ServerSetting = &Server{}

// 初始化配置
func Setup() {
	var err error
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalf("Failed to parse config/config.ini: %v", err)
	}

	err = cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Failed to map database section: %v", err)
	}

	err = cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Failed to map server section: %v", err)
	}

	// 先连接到 PostgreSQL 服务器，不指定数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		DatabaseSetting.User,
		DatabaseSetting.Password,
		DatabaseSetting.Host)

	tempDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL server: %v", err)
	}

	// 创建数据库（如果不存在）
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", DatabaseSetting.Name)
	err = tempDB.Exec(sql).Error
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// 连接到指定的数据库
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DatabaseSetting.User,
		DatabaseSetting.Password,
		DatabaseSetting.Host,
		DatabaseSetting.Name)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   DatabaseSetting.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库
	// 先创建基础表
	err = db.AutoMigrate(
		&User{},
		&Project{},
		&BaseInfo{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建 Event 表
	err = db.AutoMigrate(&EventMain{})
	if err != nil {
		log.Fatalf("Failed to migrate Event table: %v", err)
	}

	// 创建各种事件表
	err = db.AutoMigrate(
		&PerformancePageDetail{},
		&PerformanceResourceDetail{},
		&PVDetail{},
		&ClickDetail{},
		&DwellDetail{},
		&IntersectionDetail{},
		&CustomDetail{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate event tables: %v", err)
	}
}

// 获取数据库连接
func GetDB() *gorm.DB {
	return db
}
