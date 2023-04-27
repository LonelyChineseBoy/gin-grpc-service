package handler

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"reflect"
	"user-srv/global"
)

func NewNacosClientConfig() *constant.ClientConfig {
	nacosClientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NacosClientConfig.NamespaceId,
		TimeoutMs:           global.NacosConfig.NacosClientConfig.TimeoutMs,
		NotLoadCacheAtStart: global.NacosConfig.NacosClientConfig.NotLoadCacheAtStart,
		LogLevel:            global.NacosConfig.NacosClientConfig.LogLevel,
	}
	return &nacosClientConfig
}

func NewNacosServerConfigs() []constant.ServerConfig {
	nacosConfigs := []constant.ServerConfig{}
	for _, value := range global.NacosConfig.NacosServerConfigs {
		serverConfig := constant.ServerConfig{}
		nacosServerConfigValue := reflect.ValueOf(value)
		ServerConfigValue := reflect.ValueOf(&serverConfig)
		for i := 0; i < nacosServerConfigValue.NumField(); i++ {
			fieldName := nacosServerConfigValue.Type().Field(i).Name
			fieldValue := nacosServerConfigValue.Field(i)
			if tagValue, ok := nacosServerConfigValue.Type().Field(i).Tag.Lookup("tagName"); ok {
				fieldName = tagValue
			}
			fieldB := ServerConfigValue.Elem().FieldByName(fieldName)
			if fieldB.IsValid() && fieldB.CanSet() {
				fieldB.Set(fieldValue)
			}
		}
		nacosConfigs = append(nacosConfigs, serverConfig)
	}
	return nacosConfigs
}

// NewNacosNamingClient 创建服务发现客户端
func NewNacosNamingClient() naming_client.INamingClient {
	nacosNamingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  NewNacosClientConfig(),
		ServerConfigs: NewNacosServerConfigs(),
	})
	if err != nil {
		zap.S().Errorf("NewNamingClient failed:%s", err)
	}
	return nacosNamingClient
}

// NewNacosConfigClient 创建动态配置客户端
func NewNacosConfigClient() config_client.IConfigClient {
	nacosConfigClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  NewNacosClientConfig(),
		ServerConfigs: NewNacosServerConfigs(),
	})
	if err != nil {
		zap.S().Errorf("NewConfigClient failed:%s", err)
	}
	return nacosConfigClient
}
