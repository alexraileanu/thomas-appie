package config

import "github.com/BurntSushi/toml"

type Thomas struct {
	Cron string `toml:"cron"`
}

type Appie struct {
	ClientName    string `toml:"client_name"`
	ClientVersion string `toml:"client_version"`
	UserAgent     string `toml:"user_agent"`
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
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		return err
	}
	c.Thomas = conf.Thomas
	c.Appie = conf.Appie

	return nil
}
