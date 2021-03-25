package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

//Config ...
type Config struct {
	Host      string `env:"HOST"`
	Port      string `env:"PORT"`
	RedisHost string `env:"REDIS_HOST"`
	RedisPort string `env:"REDIS_PORT"`
	LogLevel  string `env:"LOG_LEVEL"`
}

//New ...
func New(file string) (*Config, error) {

	err := godotenv.Load(file)

	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = env.Parse(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
