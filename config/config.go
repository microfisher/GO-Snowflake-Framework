package config

import (
	"bytes"
	"os"

	"snowflake/log"

	"github.com/spf13/viper"
)

var v *viper.Viper

// GetViper 获取之后只作读配置操作,不允许修改
func GetViper() *viper.Viper {
	return v
}

// ReplaceViper 配置热更，替换为新配置对象
func ReplaceViper(newConf *viper.Viper) {
	v = newConf
}

// ReloadConfig 重新加载配置
func ReloadConfig(cfgType string) error {
	// 回调时，日志已初始化过
	confJSON, err := os.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		log.Warnf("event:%s msg:read config file err %s", "event_conf_err", err.Error())
		return err
	}
	v := viper.New()
	v.SetConfigType(cfgType)
	err = v.ReadConfig(bytes.NewBuffer(confJSON))
	if err != nil {
		// EventConfErr
		log.Errorf("event:%s msg:viper parse config err %s", "event_conf_err", err.Error())
		return err
	}
	ReplaceViper(v) // 新配置替换
	// UpConf
	log.Infof("event:%s msg:update config success", "up_conf")
	return nil
}
