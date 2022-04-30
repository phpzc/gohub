package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	//测试一个V1的路由组 我们所有的V1版本的路由将存放到这里
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"Hello": "World!",
			})
		})
	}
}
