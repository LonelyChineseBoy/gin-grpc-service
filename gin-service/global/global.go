package global

import (
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Router      *gin.Engine
	Logger      *zap.Logger
	DB          *gorm.DB
	NacosConfig NacosClientAndServerConfig
)

type NacosClientConfig struct {
	TimeoutMs            uint64                            `mapstructure:"TimeoutMs"`            // timeout for requesting Nacos server, default value is 10000ms
	ListenInterval       uint64                            `mapstructure:"ListenInterval"`       // Deprecated
	BeatInterval         int64                             `mapstructure:"BeatInterval"`         // the time interval for sending beat to server,default value is 5000ms
	NamespaceId          string                            `mapstructure:"NamespaceId"`          // the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	AppName              string                            `mapstructure:"AppName"`              // the appName
	AppKey               string                            `mapstructure:"AppKey"`               // the client identity information
	Endpoint             string                            `mapstructure:"Endpoint"`             // the endpoint for get Nacos server addresses
	RegionId             string                            `mapstructure:"RegionId"`             // the regionId for kms
	AccessKey            string                            `mapstructure:"AccessKey"`            // the AccessKey for kms
	SecretKey            string                            `mapstructure:"SecretKey"`            // the SecretKey for kms
	OpenKMS              bool                              `mapstructure:"OpenKMS"`              // it's to open kms,default is false. https://help.aliyun.com/product/28933.html
	CacheDir             string                            `mapstructure:"CacheDir"`             // the directory for persist nacos service info,default value is current path
	DisableUseSnapShot   bool                              `mapstructure:"DisableUseSnapShot"`   // It's a switch, default is false, means that when get remote config fail, use local cache file instead
	UpdateThreadNum      int                               `mapstructure:"UpdateThreadNum"`      // the number of goroutine for update nacos service info,default value is 20
	NotLoadCacheAtStart  bool                              `mapstructure:"NotLoadCacheAtStart"`  // not to load persistent nacos service info in CacheDir at start time
	UpdateCacheWhenEmpty bool                              `mapstructure:"UpdateCacheWhenEmpty"` // update cache when get empty service instance from server
	Username             string                            `mapstructure:"Username"`             // the username for nacos auth
	Password             string                            `mapstructure:"Password"`             // the password for nacos auth
	LogDir               string                            `mapstructure:"LogDir"`               // the directory for log, default is current path
	LogLevel             string                            `mapstructure:"LogLevel"`             // the level of log, it's must be debug,info,warn,error, default value is info
	ContextPath          string                            `mapstructure:"ContextPath"`          // the nacos server contextpath
	AppendToStdout       bool                              `mapstructure:"AppendToStdout"`       // if append log to stdout
	LogSampling          *constant.ClientLogSamplingConfig // the sampling config of log
	LogRollingConfig     *constant.ClientLogRollingConfig  // log rolling config
	TLSCfg               constant.TLSConfig                // tls Config
}

type NacosServerConfig struct {
	Scheme      string `mapstructure:"Scheme"`      // the nacos server scheme,default=http,this is not required in 2.0
	ContextPath string `mapstructure:"ContextPath"` // the nacos server contextpath,default=/nacos,this is not required in 2.0
	IpAddr      string `mapstructure:"IpAddr"`      // the nacos server address
	Port        uint64 `mapstructure:"Port"`        // nacos server port
	GrpcPort    uint64 `mapstructure:"GrpcPort"`    // nacos server grpc port, default=server port + 1000, this is not required
}

type NacosClientAndServerConfig struct {
	NacosClientConfig `mapstructure:"client-config"`
	NacosServerConfig `mapstructure:"server-config"`
}
