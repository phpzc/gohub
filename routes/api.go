package routes

import (
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"
	"gohub/pkg/config"

	"github.com/gin-gonic/gin"
)

//注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	//测试一个V1的路由组 我们所有的V1版本的路由将存放到这里
	//v1 := r.Group("/v1")

	//如果配置了 API 域名那请求应该是：https://api.domain.com/v1/users
	//如果未配置 API 域名，请求是：https://domain.com/api/v1/users
	var v1 *gin.RouterGroup
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}

	//全局限流中间件，每小时限流 这里是所有API 根据IP 请求加起来
	// 作为参考 Github API 每小时最多 60 个请求（根据 IP）。
	// 测试时，可以调高一点。
	v1.Use(middlewares.LimitIP("200-H"))

	{
		authGroup := v1.Group("/auth")

		// 限流中间件：每小时限流，作为参考 Github API 每小时最多 60 个请求（根据 IP）
		// 测试时，可以调高一点
		authGroup.Use(middlewares.LimitIP("1000-H"))

		{
			suc := new(auth.SignupController)

			//判断手机是否已注册
			authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsPhoneExist)

			//判断 Email 是否已注册
			authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsEmailExist)

			authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", middlewares.GuestJWT(), suc.SignupUsingEmail)

			//发送验证码
			vcc := new(auth.VerifyCodeController)
			//图片验证码 需要加限流
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), vcc.ShowCaptcha)

			//发送邮件验证码
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)

			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)

			lgc := new(auth.LoginController)
			//使用手机号短信验证码进行登录
			authGroup.POST("/login/using-phone", middlewares.GuestJWT(), lgc.LoginByPhone)
			// 支持手机号，Email 和 用户名
			authGroup.POST("/login/using-password", middlewares.GuestJWT(), lgc.LoginByPassword)
			//刷新token
			authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), lgc.RefreshToken)

			//重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), pwc.ResetByEmail)
		}

		uc := new(controllers.UsersController)
		//获取当前用户
		v1.GET("/user", middlewares.AuthJWT(), uc.CurrentUser)
		usersGroup := v1.Group("/users")
		{
			usersGroup.GET("", uc.Index)
			usersGroup.PUT("", middlewares.AuthJWT(), uc.UpdateProfile)
			usersGroup.PUT("/email", middlewares.AuthJWT(), uc.UpdateEmail)
			usersGroup.PUT("/phone", middlewares.AuthJWT(), uc.UpdatePhone)
			usersGroup.PUT("/password", middlewares.AuthJWT(), uc.UpdatePassword)
			usersGroup.PUT("/avatar", middlewares.AuthJWT(), uc.UpdateAvatar)
		}

		cgc := new(controllers.CategoriesController)
		cgcGroup := v1.Group("/categories")
		{
			//分类列表 分页
			cgcGroup.GET("", cgc.Index)
			//创建分类
			cgcGroup.POST("", middlewares.AuthJWT(), cgc.Store)
			//更新分类
			cgcGroup.PUT("/:id", middlewares.AuthJWT(), cgc.Update)
			//删除分类
			cgcGroup.DELETE("/:id", middlewares.AuthJWT(), cgc.Delete)
		}

		tpc := new(controllers.TopicsController)
		tpcGroup := v1.Group("/topics")
		{
			//话题列表
			tpcGroup.GET("", tpc.Index)
			//创建话题
			tpcGroup.POST("", middlewares.AuthJWT(), tpc.Store)
			//更新话题
			tpcGroup.PUT("/:id", middlewares.AuthJWT(), tpc.Update)
			//删除话题
			tpcGroup.DELETE("/:id", middlewares.AuthJWT(), tpc.Delete)
			//显示话题
			tpcGroup.GET("/:id", tpc.Show)
		}

		lsc := new(controllers.LinksController)
		linksGroup := v1.Group("/links")
		{
			linksGroup.GET("", lsc.Index)
		}
	}
}
