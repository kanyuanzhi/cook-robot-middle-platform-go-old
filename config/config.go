package config

import (
	"cook-robot-middle-platform-go/logger"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

var App *AppConfig

type AppConfig struct {
	DebugMode      bool
	GRPC           GRPCConfig     `mapstructure:"grpc"`
	HTTP           HTTPConfig     `mapstructure:"http"`
	SoftwareUpdate SoftwareUpdate `mapstructure:"softwareUpdate"`
}

type GRPCConfig struct {
	TargetHost string `mapstructure:"targetHost"`
	TargetPort uint16 `mapstructure:"targetPort"`
}

type HTTPConfig struct {
	Host    string `mapstructure:"host"`
	Port    uint16 `mapstructure:"port"`
	UseSSL  bool   `mapstructure:"useSSL"`
	SSLDir  string `mapstructure:"sslDir"`
	CerFile string `mapstructure:"cerFile"`
	KeyFile string `mapstructure:"keyFile"`
}

type SoftwareUpdate struct {
	ServerHost             string `mapstructure:"serverHost"`
	ServerPort             uint16 `mapstructure:"serverPort"`
	Filename               string `mapstructure:"filename"`
	SavePath               string `mapstructure:"savePath"`
	UnzipPath              string `mapstructure:"unzipPath"`
	UIFolderName           string `mapstructure:"uiFolderName"`
	MiddlePlatformFilename string `mapstructure:"middlePlatformFilename"`
	ControllerFilename     string `mapstructure:"controllerFilename"`
}

func (m *AppConfig) Reload() {
	viper.SetConfigName("middlePlatformConfig")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Println("无法读取配置文件:", err)
		return
	}

	err = viper.Unmarshal(App)
	if err != nil {
		logger.Log.Println("解析配置文件失败:", err)
		return
	}
}

func init() {
	file, err := os.OpenFile("middlePlatform.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Unable to create log file:", err)
	}
	//defer file.Close()

	// 使用io.MultiWriter将日志输出同时写入控制台和文件
	logWriter := io.MultiWriter(os.Stdout, file)

	logger.Log = log.New(logWriter, "", log.Lmicroseconds)

	App = &AppConfig{
		DebugMode: false,
		GRPC: GRPCConfig{
			TargetHost: "localhost",
			TargetPort: 50051,
		},
		HTTP: HTTPConfig{
			Host:    "0.0.0.0",
			Port:    8889,
			UseSSL:  false,
			SSLDir:  "ssl",
			CerFile: "certificate.crt",
			KeyFile: "private.key",
		},
		SoftwareUpdate: SoftwareUpdate{
			ServerHost:             "http://124.71.146.83",
			ServerPort:             12306,
			Filename:               "software.zip",
			SavePath:               ".",
			UnzipPath:              ".\\unzip",
			UIFolderName:           "cookRobot-linux-arm64",
			MiddlePlatformFilename: "cook-robot-middle-platform-go",
			ControllerFilename:     "cook-robot-controller-go",
		},
	}
	App.Reload()
}
