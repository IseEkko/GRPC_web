package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"mxshop-api/user-web/config"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}

func InitConfig() {
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s.pro.yaml", configFilePrefix)
	global.ServerConfig = &config.ServerConfig{}
	v := viper.New()

	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	err := v.Unmarshal(global.ServerConfig)
	if err != nil {
		panic(err)
	}
	zap.S().Info("配置信息：&v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("配置文件产生变化")
		_ = v.ReadRemoteConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Info("配置信息：&v", global.ServerConfig)
	})

}
