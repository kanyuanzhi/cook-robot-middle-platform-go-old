package config

import (
	"cook-robot-middle-platform-go/utils"
)

var App = &AppConfig{}

type AppConfig struct {
	DebugMode  bool
	HTTP       HTTPConfig       `mapstructure:"http"`
	Controller ControllerConfig `mapstructure:"controller"`
	Updater    Updater          `mapstructure:"updater"`
}

type ControllerConfig struct {
	GRPCHost string `mapstructure:"grpcHost"`
	GRPCPort uint16 `mapstructure:"grpcPort"`
}

type HTTPConfig struct {
	Host    string `mapstructure:"host"`
	Port    uint16 `mapstructure:"port"`
	UseSSL  bool   `mapstructure:"useSSL"`
	SSLDir  string `mapstructure:"sslDir"`
	CerFile string `mapstructure:"cerFile"`
	KeyFile string `mapstructure:"keyFile"`
}

type Updater struct {
	Host                   string `mapstructure:"host"`
	GRPCPort               uint16 `mapstructure:"grpcPort"`
	FileServerPort         uint16 `mapstructure:"fileServerPort"`
	SavePath               string `mapstructure:"savePath"`
	UnzipPath              string `mapstructure:"unzipPath"`
	UIFolderName           string `mapstructure:"uiFolderName"`
	MiddlePlatformFilename string `mapstructure:"middlePlatformFilename"`
	ControllerFilename     string `mapstructure:"controllerFilename"`
}

func (m *AppConfig) Reload() {
	utils.Reload("middlePlatformConfig", App)
}
