package main

import (
	"fmt"
	"gohub/app/cmd"
	"gohub/app/cmd/make"
	"gohub/bootstrap"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"os"

	btsConfig "gohub/config"

	"github.com/spf13/cobra"
)

func init() {
	btsConfig.Initialize()
}

func main() {

	//应用的主入口 默认调用 cmd.CmdServe命令
	var rootCmd = &cobra.Command{
		Use:   config.Get("app.name"),
		Short: "A Simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		//rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			//配置初始化 依赖命令行 --env参数
			config.InitConfig(cmd.Env)

			//初始化logger
			bootstrap.SetupLogger()

			//初始化数据库
			bootstrap.SetupDB()

			//初始化redis
			bootstrap.SetupRedis()

			//初始化缓存
		},
	}

	//注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
		cmd.CmdKey,
		cmd.CmdPlay,
		make.CmdMake,
	)

	//配置默认运行web服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	//注册全局参数 --env
	cmd.RegisterGlobalFlags(rootCmd)

	//执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}

	//解析命令行参数
	// var env string
	// flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	// flag.Parse()

	// fmt.Println("env参数值:" + env)

	// config.InitConfig(env)

	// //初始化logger
	// bootstrap.SetupLogger()

	// // 设置 gin 的运行模式，支持 debug, release, test
	// // release 会屏蔽调试信息，官方建议生产环境中使用
	// // 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// // 故此设置为 release，有特殊情况手动改为 debug 即可
	// gin.SetMode(gin.ReleaseMode)

	// //new gin
	// router := gin.New()

	// //初始化DB
	// bootstrap.SetupDB()
	// //初始化redis
	// bootstrap.SetupRedis()
	// //初始化路由绑定
	// bootstrap.SetupRoute(router)

	// err := router.Run(":" + config.Get("app.port"))

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}
