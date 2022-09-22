package config

import (
	"github.com/MedmeFord/RestAPItu/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	MongoDB struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		Database    string `json:"database"`
		Auth_db     string `json:"auth_db"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Collections string `json:"collections"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		err := cleanenv.ReadConfig("./pkg/config.yaml", instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
