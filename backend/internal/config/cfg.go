package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/go-ini/ini"
)

type Config struct {
	mu      sync.RWMutex
	cfgFile *ini.File

	LogFileMaxAge int

	RestIp   string
	RestPort string

	RepoDir  string
	ShellDir string

	GitId  string
	GitPwd string
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

	repoDir := cfgFile.Section("DIR").Key("REPO_DIR").MustString("/tmp/repo")
	config.RepoDir = filepath.FromSlash(repoDir)
	os.MkdirAll(config.RepoDir, 0755)

	shellDir := cfgFile.Section("DIR").Key("SHELL_DIR").MustString("/tmp/script")
	config.ShellDir = filepath.FromSlash(shellDir)
	os.MkdirAll(config.ShellDir, 0755)

	config.GitId = cfgFile.Section("GITAUTH").Key("ID").MustString("")
	config.GitPwd = cfgFile.Section("GITAUTH").Key("PASSWD").MustString("")

	return config
}
