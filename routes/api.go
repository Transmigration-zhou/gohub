package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controllers/api/v1/auth"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// 判断 Email 是否已注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			// 短信验证码注册
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)
			// 邮件验证码注册
			authGroup.POST("/signup/using-email", suc.SignupUsingEmail)

			vcc := new(auth.VerifyCodeController)
			// 图片验证码
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			// 短信验证码
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			// 邮件验证码
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)

			lgc := new(auth.LoginController)
			// 手机号登录
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			// 密码登录
			authGroup.POST("/login/using-password", lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", lgc.RefreshToken)

			// 重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-phone", pwc.ResetByPhone)
			authGroup.POST("/password-reset/using-email", pwc.ResetByEmail)
		}
	}
}
