package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}
}

type AppConfig struct {
	Logger LoggerConfig
	Azure  AzureConfig
	Gin    GinConfig
}

func GetAppConfig(validate bool) AppConfig {
	loadEnvironmentVariables()
	LOGGER_CONF := assembleLoggerConfig()
	AZURE_CONF := assembleAzureConfig()
	GIN_CONF := assembleGinConfig()
	APP_CONF := AppConfig{
		Logger: LOGGER_CONF,
		Azure:  AZURE_CONF,
		Gin:    GIN_CONF,
	}
	return APP_CONF
}
