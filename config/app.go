package config

import (
	"github.com/spf13/viper"
)

// AppConfig represent the data-struct for configuration
type AppConfig struct {
	// another stuff , may be needed by configuration
}

func init() {
	// load and read config file
	setConfigFile()
}

func setConfigFile() {
	// find environment file
	viper.SetConfigFile(`.env`)
	// error handling for specific case
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(".env file not found!, please copy .env.example and paste as .env")
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}
}

func (config AppConfig) GetAppName() string {
	return viper.GetString("APP_NAME")
}

func (config AppConfig) GetAppVersion() string {
	return viper.GetString("APP_VERSION")
}

// InitAppConfig initialize the app configuration
func InitAppConfig() *AppConfig {
	return &AppConfig{}
}
