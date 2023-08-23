package db

import (
	"cook-robot-middle-platform-go/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var Postgres *gorm.DB

func init() {
	dsn := "user=postgres password=123456 dbname=cook_robot port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		logger.Log.Println("连接postgres失败")
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Println("获取postgres底层数据库实例失败")
		return
	}

	sqlDB.SetMaxOpenConns(10) // 最大打开连接数
	sqlDB.SetMaxIdleConns(5)  // 最大闲置连接数
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Log.Println("连接postgres数据库成功")
	Postgres = db
}
