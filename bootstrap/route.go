package bootstrap

import (
	"gohub/app/http/middlwares"
	"gohub/routes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//SetupRoute 路由初始化
func SetRoute(router *gin.Engine) {

	//注册全局中间件
	registerGlobalMiddlWare(router)
	//注册API路由
	routes.RegisterAPIRoutes(router)
	//配置404路由
	setup404Handler(router)
}

func registerGlobalMiddlWare(router *gin.Engine) {

	router.Use(
		//gin.Logger(),
		middlwares.Logger(),
		gin.Recovery(),
	)
}

func setup404Handler(router *gin.Engine) {

	router.NoRoute(func(c *gin.Context) {
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			c.String(http.StatusNotFound, "页面返回404")
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}
