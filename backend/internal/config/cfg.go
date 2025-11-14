package config

import (
	"sync"

	"github.com/go-ini/ini"
)

type Config struct {
	mu      sync.RWMutex
	cfgFile *ini.File

	LogFileMaxAge int

	RestIp   string
	RestPort string
}

func NewConfig() (*Config, error) {
	cfgFile, err := ini.Load("env.ini")
	if err != nil {
		return nil, err
	}
	return initConfig(cfgFile), nil
}

func initConfig(cfgFile *ini.File) *Config {
	config := &Config{}

	config.cfgFile = cfgFile
	config.LogFileMaxAge = cfgFile.Section("LOG").Key("MAX_AGE").MustInt(-5)
	config.RestIp = cfgFile.Section("REST").Key("IP").MustString("")
	config.RestPort = cfgFile.Section("REST").Key("PORT").MustString("8080")

	return config
}
