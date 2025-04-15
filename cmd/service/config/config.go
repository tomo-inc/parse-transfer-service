package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

type EvmInfo struct {
	Endpoint     string `mapstructure:"endpoint"`
	SupportDebug bool   `mapstructure:"support_debug"`
}

type Config struct {
	ListenHost    string             `mapstructure:"listen_host"`
	EVMEndpoints  map[string]EvmInfo `mapstructure:"evm"`
	SOlEndpoints  map[string]string  `mapstructure:"sol"`
	TRONEndpoints map[string]string  `mapstructure:"tron"`
	AlertConfig   AlertConfig        `mapstructure:"alert_config"`
}

type AlertConfig struct {
	Interval  uint64 `mapstructure:"interval"` // s
	LarkBotId string `mapstructure:"lark_bot_id"`
}

func GetConfig(cfgFile string) *Config {
	if cfgFile == "" {
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")

		viper.SetConfigName("config")
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("failed to read config", zap.Error(err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("failed to unmarshal config", zap.Error(err))
	}

	return &config
}
