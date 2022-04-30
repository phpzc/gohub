package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"
	"gohub/pkg/config"

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

	//new gin
	router := gin.New()

	bootstrap.SetRoute(router)

	err := router.Run(":" + config.Get("app.port"))

	if err != nil {
		fmt.Println(err.Error())
	}
}
