package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/***
短信验证码
*/
func SendSms(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
