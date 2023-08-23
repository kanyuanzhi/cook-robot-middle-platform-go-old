package utils

import (
	"github.com/spf13/viper"
	"log"
)

// Reload 解析配置文件configName到target
func Reload(configName string, target interface{}) {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("无法读取配置文件:", err)
		return
	}

	err = viper.Unmarshal(target)
	if err != nil {
		log.Println("解析配置文件失败:", err)
		return
	}
}
