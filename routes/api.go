package routes

import (
	"github.com/gin-gonic/gin"
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"
	"gohub/pkg/config"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	var v1 *gin.RouterGroup
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}
	// 全局限流中间件：每小时限流 200 次
	v1.Use(middlewares.LimitIP("200-H"))
	{
		authGroup := v1.Group("/auth")
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			// 注册
			suc := new(auth.SignupController)
			authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), suc.IsEmailExist)
			authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", middlewares.GuestJWT(), suc.SignupUsingEmail)

			// 验证码
			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("20-H"), vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)

			// 登录
			lgc := new(auth.LoginController)
			authGroup.POST("/login/using-phone", middlewares.GuestJWT(), lgc.LoginByPhone)
			authGroup.POST("/login/using-password", middlewares.GuestJWT(), lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), lgc.RefreshToken)

			// 重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), pwc.ResetByEmail)

			// 获取用户
			uc := new(controllers.UsersController)
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

			// 分类
			cgc := new(controllers.CategoriesController)
			cgcGroup := v1.Group("/categories")
			{
				cgcGroup.GET("", cgc.Index)
				cgcGroup.POST("", middlewares.AuthJWT(), cgc.Store)
				cgcGroup.PUT("/:id", middlewares.AuthJWT(), cgc.Update)
				cgcGroup.DELETE("/:id", middlewares.AuthJWT(), cgc.Delete)
			}

			// 话题
			tpc := new(controllers.TopicsController)
			tpcGroup := v1.Group("/topics")
			{
				tpcGroup.GET("", tpc.Index)
				tpcGroup.POST("", middlewares.AuthJWT(), tpc.Store)
				tpcGroup.PUT("/:id", middlewares.AuthJWT(), tpc.Update)
				tpcGroup.DELETE("/:id", middlewares.AuthJWT(), tpc.Delete)
				tpcGroup.GET("/:id", tpc.Show)
			}

			// 友情链接
			lsc := new(controllers.LinksController)
			linksGroup := v1.Group("/links")
			{
				linksGroup.GET("", lsc.Index)
			}
		}
	}
}
