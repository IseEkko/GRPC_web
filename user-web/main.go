package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mxshop-api/user-web/utiles"
	myvalidator "mxshop-api/user-web/validator"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {

	//1.初始化router
	//2.初始化logger
	initialize.Logger()
	initialize.InitConfig()
	Router := initialize.Routers()
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
		return
	}
	//初始化srv连接
	initialize.InitSrvConn2()
	//这里实现的是动态的端口注册
	port, err := utiles.GetFreePort()
	if err != nil {
		global.ServerConfig.Port = port
	}

	//注册验证器
	/***
	这里我们还定义了手机号码验证的中文翻译
	*/
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	/**
	1. s()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
	日志是分级别的,debug,info,warn,error,fetal
	从左往右级别依次升高
	*/
	//这里对端口进行了设置动态的端口
	zap.S().Infof("启动服务器，端口：%d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
