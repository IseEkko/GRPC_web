package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BseRouter := Router.Group("base")
	{
		BseRouter.GET("captcha", api.GetCaptcha)
		BseRouter.POST("send_sms")
	}
}
