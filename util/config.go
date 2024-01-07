package util

import (
	"time"

	"github.com/spf13/viper"
)

const (
	configName = "app" // configuration file name.
	configType = "env" // configuration file type.
)

// Config stores all configuration parameters required to start
// the application. The values are read by viper from a configuration
// file.
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBAddr              string        `mapstructure:"DB_ADDRESS"`
	DBUser              string        `mapstructure:"DB_USER"`
	DBName              string        `mapstructure:"DB_NAME"`
	ServerAddr          string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
