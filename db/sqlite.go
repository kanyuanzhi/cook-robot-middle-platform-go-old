package db

import (
	"cook-robot-middle-platform-go/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var SQLiteDB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("cook.db"), &gorm.Config{})
	if err != nil {
		panic("连接sqlite数据库失败")
	}
	logger.Log.Println("连接sqlite数据库成功")
	SQLiteDB = db
}
