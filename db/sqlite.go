package db

import (
	"cook-robot-middle-platform-go/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var SQLiteDB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("cook.db"), &gorm.Config{})
	if err != nil {
		panic("连接本地数据库失败")
	}
	logger.Log.Println("连接本地数据库成功")
	SQLiteDB = db
}
