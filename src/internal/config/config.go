package config

import (
	"flag"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port    int           `yaml:"port" env-default:"8000"`
		Timeout time.Duration `yaml:"timeout" env-default:"10s"`
	} `yaml:"server"`

	WebSocket struct {
		Timeout    time.Duration `yaml:"timeout" env-default:"10s"`
		PongWait   time.Duration `yaml:"pong_wait" env-default:"30s"`
		PingPeriod time.Duration `yaml:"ping_period" env-default:"15s"`
	} `yaml:"websocket"`

	Database struct {
		Addr     string `yaml:"addr" env-required:"true"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"database"`

	Redis struct {
		Addr     string `yaml:"addr" env-required:"true"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db" env-default:"0"`
	} `yaml:"redis"`

	Broker struct {
		Addr string `yaml:"addr" env-required:"true"`
	} `yaml:"broker"`

	CORS struct {
		Origins        []string `yaml:"origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
	} `yaml:"cors"`

	Media struct {
		Url string `yaml:"url" env-required:"true"`
	} `yaml:"media"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		path := getConfigPath()

		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "../config/config.yaml", "set config file")

	envPath := os.Getenv("CONFIG_PATH")

	if len(envPath) > 0 {
		path = envPath
	}

	return path
}
