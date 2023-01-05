// Package config 负责配置信息：初始化，读取配置文件，设置配置项，读取配置项
package config

import (
	viperlib "github.com/spf13/viper" //自定义包名，避免与内置viper
)

// viper库实例
var viper *viperlib.Viper

// ConfigFunc 动态加载配置信息
type ConfigFunc func() map[string]interface{}

// ConfigFuncs 先加载到此数组，loadConfig 再动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	//1. 初始化 Viper 库
	viper = viperlib.New()

	//2. 配置类型 支持 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	viper.SetConfigType("env")

	//3. 环境变量配置文件查找的路径，相当于main.go
	viper.AddConfigPath(".")

	//4. 设置环境变量前缀，用以区分 Go 的系统环境变量
	viper.SetEnvPrefix("appenv")

	//5. 读取环境变量（支持flags）
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}
