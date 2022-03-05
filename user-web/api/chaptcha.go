package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

/**
图片验证码接口
*/

var store = base64Captcha.DefaultMemStore

/**
生成验证码图片base64编码
*/
func GetCaptcha(ctx *gin.Context) {
	//这里是对验证码的配置
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorw("生成验证码错误", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	/**
	这里我们生成来base64编码后，放进img标签src后会生成图图片
	*/
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     id,
		"picPath": b64s,
	})

}
