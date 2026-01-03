package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Thomas struct {
	Cron string `toml:"cron"`
}

type Appie struct {
	ClientName         string `toml:"client_name"`
	ClientVersion      string `toml:"client_version"`
	UserAgent          string `toml:"user_agent"`
	ClientPlatformType string `toml:"client_platform_type"`
	BonusDay           int    `toml:"bonus_day"`
}

type Config struct {
	Thomas Thomas `toml:"thomas"`
	Appie  Appie  `toml:"appie"`
}

func New() Config {
	return Config{}
}

func (c *Config) ParseConfig() error {
	var conf Config
	if _, err := toml.DecodeFile(os.Getenv("CONFIG_FILE_PATH"), &conf); err != nil {
		return err
	}
	c.Thomas = conf.Thomas
	c.Appie = conf.Appie

	return nil
}
