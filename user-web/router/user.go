package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
	"net/http"
)

func InitUserRouter(Router *gin.RouterGroup) {

	UserRoter := Router.Group("user") //.Use(middlewares.JWTAuth())//使用Use添加中间件，然后这里的就是要登录验证了过后才能进入后面的逻辑
	zap.S().Info("配置用户相关的url")
	{
		UserRoter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserlist) //使用了中间件,在使用中间件的时候，我们需要做的是中间件的执行顺序
		UserRoter.POST("pwd_login", api.PasswordLogin)
		UserRoter.POST("register", api.Register)

		UserRoter.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		UserRoter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
		UserRoter.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "health",
			})
		})
	}

}
