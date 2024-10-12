package main

import (
	"oss-upload-example/core"
	"oss-upload-example/global"
	"oss-upload-example/initialize"
)

func main() {
	initialize.LoggerInit(initialize.WithLevel("debug")) // 日志初始化
	global.VIPER = initialize.Viper()                    // 初始化Viper 配置信息
	core.RunHttpServer()                                 // 启动http
}
