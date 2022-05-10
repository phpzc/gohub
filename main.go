package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"
	"gohub/pkg/captcha"
	"gohub/pkg/config"
	"gohub/pkg/logger"

	btsConfig "gohub/config"

	"github.com/gin-gonic/gin"
)

func init() {
	btsConfig.Initialize()
}

func main() {

	//解析命令行参数
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()

	fmt.Println("env参数值:" + env)

	config.InitConfig(env)

	//初始化logger
	bootstrap.SetupLogger()

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	//new gin
	router := gin.New()

	//初始化DB
	bootstrap.SetupDB()
	//初始化redis
	bootstrap.SetupRedis()
	//初始化路由绑定
	bootstrap.SetRoute(router)

	logger.Dump(captcha.NewCaptcha().VerifyCaptcha("UFuTESFSUKFhGtdmR8vs", "083804"), "正确的答案")
	logger.Dump(captcha.NewCaptcha().VerifyCaptcha("UFuTESFSUKFhGtdmR8vs", "000000"), "错误的答案")

	err := router.Run(":" + config.Get("app.port"))

	if err != nil {
		fmt.Println(err.Error())
	}
}
