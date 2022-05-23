package config

import (
	"gimg/logger"

	"github.com/spf13/viper"
)

type EngineConf struct {
	SavePath string `mapstructure:"save_path"`
	Name     string
}

type ActionConf struct {
	LoadScriptPath string `mapstructure:"load_path"`
}

type AuthConf struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pwd"`
}

type Config struct {
	Debug  bool
	Port   int
	Engine *EngineConf
	Logger *logger.Config
	Action *ActionConf
	Auth   *AuthConf
}

func defaultConfig() *Config {
	return &Config{
		Debug: false,
		Port:  8080,
		Engine: &EngineConf{
			Name:     "imagick",
			SavePath: "./images",
		},
		Action: &ActionConf{
			LoadScriptPath: "./scripts",
		},
		Logger: &logger.Config{Level: logger.DevelopmentLevel},
		Auth:   &AuthConf{User: "test", Password: "123456"},
	}
}

func Load(filename string) (*Config, error) {
	var config Config

	defaultConfig := defaultConfig()
	viper.SetDefault("debug", defaultConfig.Debug)
	viper.SetDefault("port", defaultConfig.Port)
	viper.SetDefault("engine", defaultConfig.Engine)
	viper.SetDefault("logger", defaultConfig.Logger)
	viper.SetDefault("action", defaultConfig.Action)
	viper.SetDefault("auth", defaultConfig.Auth)
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
