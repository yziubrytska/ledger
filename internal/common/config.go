package common

import "github.com/caarlos0/env"

type Config struct {
	DatabaseCredentials string `env:"DATABASE_URL" envDefault:""`
	ListenAddress       string `env:"LISTEN_ADDRESS" envDefault:":8080"`
	LogLevel            string `env:"LOG_LEVEL" envDefault:"debug"`
}

func NewConfig() (*Config, error) {
	c := new(Config)
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
