package config

import (
	"gimg/logger"
	"github.com/spf13/viper"
)

type EngineConf struct {
	SavePath string `mapstructure:"save_path"`
	Name     string
}

type Config struct {
	Debug  bool
	Port   int
	Engine *EngineConf
	Logger *logger.Config
}

func defaultConfig() *Config {
	return &Config{
		Debug: false,
		Port:  8080,
		Engine: &EngineConf{
			Name:     "imagick",
			SavePath: "./images",
		},
		Logger: &logger.Config{Level: logger.DevelopmentLevel},
	}
}

func Load(filename string) (*Config, error) {
	var config Config

	defaultConfig := defaultConfig()
	viper.SetDefault("debug", defaultConfig.Debug)
	viper.SetDefault("port", defaultConfig.Port)
	viper.SetDefault("engine", defaultConfig.Engine)
	viper.SetDefault("logger", defaultConfig.Logger)
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
