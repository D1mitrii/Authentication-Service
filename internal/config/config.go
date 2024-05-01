package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string     `yaml:"env" env:"ENV" env-required:"true"`
	JWT    JWT        `yaml:"jwt"`
	PG     Postgres   `yaml:"storage"`
	RDB    Redis      `yaml:"redis"`
	HTTP   HTTPServer `yaml:"http"`
	GRPC   GRPC       `yaml:"grpc"`
	Hasher Hasher     `yaml:"hasher"`
}

type HTTPServer struct {
	Port    int           `yaml:"port" env:"HTTP_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
}

type GRPC struct {
	Port int `yaml:"port" env:"GRPC_PORT" env-required:"true"`
}

type Postgres struct {
	URL string `yaml:"url" env:"POSTGRES_URL" env-required:"true"`
}

type Redis struct {
	Host     string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"REDIS_PORT" env-required:"true"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
}

type JWT struct {
	Secret      string        `yaml:"secret_key" env:"JWT_SECRET" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env:"JWT_TOKEN_TTL" env-required:"true"`
	RefreshTime time.Duration `yaml:"refresh_time" env:"JWT_REFRESH" env-required:"true"`
}

type Hasher struct {
	Salt int `yaml:"salt" env:"HASH_SALT" env-default:"10"`
}

func MustLoad() *Config {
	var cfg Config
	path := fetchConfigPath()
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic("config file doesn't exist: " + path)
		}
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic("failed to read config file: " + err.Error())
		}
		return &cfg
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read configuration from env")
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
