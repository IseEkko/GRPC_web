package initialize

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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

	err := v.Unmarshal(global.NacosConfig)
	if err != nil {
		panic(err)
	}
	zap.S().Info("配置信息：&v", global.NacosConfig)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("配置文件产生变化")
		_ = v.ReadRemoteConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Info("配置信息：&v", global.ServerConfig)
	})
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache", //这里配置的是缓存
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		panic(err)
	}
	content, errs := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if errs != nil {
		panic(errs)
	}
	fmt.Println(content)
	//上面就是获取配置文件

	//我这里将json转换成struct
	serverconfig := config.ServerConfig{}
	err = json.Unmarshal([]byte(content), &serverconfig)
	if err != nil {
		zap.S().Fatal("读取nacos配置失败", err.Error())
	}
	fmt.Println(serverconfig)

}
