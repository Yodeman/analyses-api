package util

import "github.com/spf13/viper"

// Config stores all configuration of the application.
// The values are read by viper from a config file.
type Config struct {
    DBDriver            string          `mapstructure:"DB_DRIVER"`
    DBAddr              string          `mapstructure:"DB_ADDRESS"`
    DBUser              string          `mapstructure:"DB_USER"`
    DBName              string          `mapstructure:"DB_NAME"`
    ServerAddr          string          `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}
