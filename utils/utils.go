package utils

import "github.com/spf13/viper"

func GetUrl() string {
	return viper.GetString("url")
}
