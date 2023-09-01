package logger

import (
	"github.com/gookit/slog/rotatefile"
	"io"
	"log"
	"os"
)

var Log *log.Logger

func init() {
	//file, err := os.OpenFile("controller.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal("Unable to create log file:", err)
	//}
	//defer file.Close()

	config := rotatefile.NewConfigWith(func(c *rotatefile.Config) {
		c.Filepath = "logs/middle-platform.log"
		c.FilePerm = rotatefile.DefaultFilePerm
		c.RotateMode = rotatefile.ModeCreate
		c.RotateTime = rotatefile.EveryDay
		c.BackupNum = 7
		c.BackupTime = 0
	})

	rf, _ := config.Create()

	// 使用io.MultiWriter将日志输出同时写入控制台和文件
	logWriter := io.MultiWriter(os.Stdout, rf)

	flag := log.Ldate | log.Ltime | log.Lmicroseconds // log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile
	Log = log.New(logWriter, "", flag)
}
