package model

import (
	"fmt"
	"log"
	"time"

	"github.com/go-ini/ini"
	"gorm.io/driver/mysql"
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

	var tempDB *gorm.DB
	var dsn string

	// 根据数据库类型选择正确的数据库驱动
	switch DatabaseSetting.Type {
	case "mysql":
		// 连接到 MySQL 服务器，不指定数据库
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
			DatabaseSetting.User,
			DatabaseSetting.Password,
			DatabaseSetting.Host)
		tempDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "pgsql":
		// 连接到 PostgreSQL 服务器，不指定数据库
		dsn = fmt.Sprintf("host=%s user=%s password=%s sslmode=disable",
			DatabaseSetting.Host,
			DatabaseSetting.User,
			DatabaseSetting.Password)
		tempDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported database type: %s", DatabaseSetting.Type)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database server: %v", err)
	}

	// 创建数据库（如果不存在）
	var sql string
	switch DatabaseSetting.Type {
	case "mysql":
		sql = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", DatabaseSetting.Name)
	case "pgsql":
		// 检查数据库是否存在
		var exists bool
		err = tempDB.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", DatabaseSetting.Name).Scan(&exists).Error
		if err != nil {
			log.Fatalf("Failed to check if database exists: %v", err)
		}

		// 如果数据库不存在，则创建
		if !exists {
			sql = fmt.Sprintf("CREATE DATABASE %s;", DatabaseSetting.Name)
		} else {
			// 数据库已存在，跳过创建
			sql = ""
		}
	}

	// 执行创建数据库的SQL语句
	if sql != "" {
		err = tempDB.Exec(sql).Error
	}
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// 连接到指定的数据库
	var dialector gorm.Dialector
	switch DatabaseSetting.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DatabaseSetting.User,
			DatabaseSetting.Password,
			DatabaseSetting.Host,
			DatabaseSetting.Name)
		dialector = mysql.Open(dsn)
	case "pgsql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			DatabaseSetting.Host,
			DatabaseSetting.User,
			DatabaseSetting.Password,
			DatabaseSetting.Name)
		dialector = postgres.Open(dsn)
	default:
		log.Fatalf("Unsupported database type: %s", DatabaseSetting.Type)
	}

	db, err = gorm.Open(dialector, &gorm.Config{
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

	// 创建错误相关表
	err = db.AutoMigrate(
		&ErrorDetail{},
		&HttpErrorDetail{},
		&ResourceErrorDetail{},
		&VueErrorDetail{},
		&ReactErrorDetail{},
		&ErrorGroup{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate error tables: %v", err)
	}
}

// 获取数据库连接
func GetDB() *gorm.DB {
	return db
}
