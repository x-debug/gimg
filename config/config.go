package config

import (
	"gimg/logger"

	"github.com/spf13/viper"
)

type ProxyConf struct {
	BaseRemotePath string `mapstructure:"base_remote_path"`
	RequestTimeout int    `mapstructure:"request_timeout"`
}

type EngineConf struct {
	SavePath  string `mapstructure:"save_path"`
	CachePath string `mapstructure:"cache_path"`
	Name      string
}

type ActionConf struct {
	LoadScriptPath string `mapstructure:"load_path"`
}

type AuthConf struct {
	Close    bool   `mapstructure:"close"`
	Type     string `mapstructure:"type"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pwd"`
}

type CacheBrockerConf struct {
	Addr string `mapstructure:"addr"`
	Port int    `mapstructure:"port"`
}

type CacheConf struct {
	Type     string              `mapstructure:"type"`
	LifeTime int64               `mapstructure:"life_time"` //only for memory cache
	Brockers []*CacheBrockerConf `mapstructure:"brockers"`
}

type Config struct {
	Debug  bool
	Port   int
	Engine *EngineConf
	Logger *logger.Config
	Action *ActionConf
	Auth   *AuthConf
	Cache  *CacheConf
	Proxy  *ProxyConf
}

func defaultConfig() *Config {
	return &Config{
		Debug: false,
		Port:  8080,
		Engine: &EngineConf{
			Name:      "imagick",
			SavePath:  "./images",
			CachePath: "./images",
		},
		Action: &ActionConf{
			LoadScriptPath: "./scripts",
		},
		Logger: &logger.Config{Level: logger.DevelopmentLevel},
		Auth:   &AuthConf{User: "test", Password: "123456", Type: "basic", Close: true},
		Cache:  &CacheConf{Type: "memory", Brockers: []*CacheBrockerConf{}, LifeTime: 60},
		Proxy:  &ProxyConf{BaseRemotePath: "", RequestTimeout: 10}, //10 second for default setting
	}
}

func Load(filename string) (*Config, error) {
	var config Config

	defaultConfig := defaultConfig()
	viper.SetDefault("debug", defaultConfig.Debug)
	viper.SetDefault("port", defaultConfig.Port)
	viper.SetDefault("engine.save_path", defaultConfig.Engine.SavePath)
	viper.SetDefault("engine.cache_path", defaultConfig.Engine.CachePath)
	viper.SetDefault("engine.name", defaultConfig.Engine.Name)
	viper.SetDefault("logger", defaultConfig.Logger)
	viper.SetDefault("action.load_path", defaultConfig.Action.LoadScriptPath)
	viper.SetDefault("auth.close", defaultConfig.Auth.Close)
	viper.SetDefault("auth.pwd", defaultConfig.Auth.Password)
	viper.SetDefault("auth.type", defaultConfig.Auth.Type)
	viper.SetDefault("auth.user", defaultConfig.Auth.User)
	viper.SetDefault("cache.life_time", defaultConfig.Cache.LifeTime)
	viper.SetDefault("cache.type", defaultConfig.Cache.Type)
	viper.SetDefault("cache.brockers", defaultConfig.Cache.Brockers)
	viper.SetDefault("proxy.base_remote_path", defaultConfig.Proxy.BaseRemotePath)
	viper.SetDefault("proxy.request_timeout", defaultConfig.Proxy.RequestTimeout)
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
