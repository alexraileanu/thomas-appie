package config

import (
	"os"

	"github.com/BurntSushi/toml"

	"github.com/alexraileanu/thomas-appie/pkg/logger"
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

	V2           bool   `toml:"v2"`
	XApplication string `toml:"x_application"`
}

type Config struct {
	Thomas Thomas `toml:"thomas"`
	Appie  Appie  `toml:"appie"`
}

func New() Config {
	return Config{}
}

func (c *Config) ParseConfig(loggerService *logger.Service) error {
	path := os.Getenv("CONFIG_FILE_PATH")
	loggerService.Debug("Parsing config file", map[string]interface{}{"path": path})

	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return err
	}
	c.Thomas = conf.Thomas
	c.Appie = conf.Appie

	loggerService.Debug("Parsed config", map[string]interface{}{
		"cron":     c.Thomas.Cron,
		"v2":       c.Appie.V2,
		"bonus_day": c.Appie.BonusDay,
	})
	return nil
}
