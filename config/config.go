package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

func init() {
	config = viper.New()
	config.SetConfigType("yaml")
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("can not resolve config")
	}
}

func GetConfig() *viper.Viper {
	return config
}
