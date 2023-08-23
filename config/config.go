package config

import (
	"cook-robot-middle-platform-go/logger"
	"io"
	"log"
	"os"
)

func init() {
	file, err := os.OpenFile("middlePlatform.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Unable to create log file:", err)
	}
	//defer file.Close()

	// 使用io.MultiWriter将日志输出同时写入控制台和文件
	logWriter := io.MultiWriter(os.Stdout, file)
	logger.Log = log.New(logWriter, "", log.Lmicroseconds)

	App.Reload()
}
