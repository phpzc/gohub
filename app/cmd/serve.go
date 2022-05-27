package cmd

import (
	"gohub/bootstrap"
	"gohub/pkg/config"
	"gohub/pkg/console"
	"gohub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// CmdServe 表示可用的web子命令.
var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {

	gin.SetMode(gin.ReleaseMode)

	//gin实例
	router := gin.New()

	//初始化路由绑定
	bootstrap.SetupRoute(router)

	//运行服务器
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server,error:" + err.Error())
	}
}
