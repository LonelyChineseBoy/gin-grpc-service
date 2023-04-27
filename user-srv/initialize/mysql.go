package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"user-srv/global"
	"user-srv/handler"
)

type config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Addr     string `json:"addr"`
	Dbname   string `json:"dbname"`
}

func InitGormConfig() {
	handler.ReadConfigByYamlFile("config/mysql.yaml", &global.MysqlConfigParam)
	client := handler.NewNacosConfigClient()
	mysqlConfig, getConfigErr := client.GetConfig(vo.ConfigParam{
		DataId: global.MysqlConfigParam.DataId,
		Group:  global.MysqlConfigParam.Group,
	})
	if getConfigErr != nil {
		zap.S().Errorf("get mysql.yaml config failed %v", getConfigErr)
	}
	var sql config
	jsonErr := json.Unmarshal([]byte(mysqlConfig), &sql)
	if jsonErr != nil {
		zap.S().Errorf("Unmarshal mysql config faield %v", jsonErr)
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", sql.Username, sql.Password, sql.Addr, sql.Dbname)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Errorf("connecting mysql.yaml seerver failed %v", err)
	}
}
