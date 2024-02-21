package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	JWT        `yaml:"jwt"`
	Storage    `yaml:"storage"`
	HTTPServer `yaml:"http"`
}

type HTTPServer struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Storage struct {
	URL  string `yaml:"url" env-required:"true"`
	Name string `yaml:"name" env-required:"true"`
}

type JWT struct {
	SecretKey   string        `yaml:"secret_key" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	RefreshTime time.Duration `yaml:"refresh_time" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config file: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var result string
	flag.StringVar(&result, "config", "", "path to config file")
	flag.Parse()
	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}
	return result
}
