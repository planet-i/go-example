package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")        // 设置配置文件格式
	viper.SetConfigName("config.yaml") // 设置配置文件名
	viper.AddConfigPath(".")           // 添加配置文件寻找路径
	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("file not found.")
			return
		}
		panic(err)
	}
}
