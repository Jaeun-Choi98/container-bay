package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-ini/ini"
)

type Config struct {
	mu      sync.RWMutex
	cfgFile *ini.File

	ShellName string

	LogFileMaxAge int

	RestIp   string
	RestPort string

	RepoDir   string
	ShellDir  string
	VolumeDir string

	GitId  string
	GitPwd string

	DockerRepoIp   string
	DockerRepoPort string
	DockerRepoId   string
	DockerRepoPwd  string

	BuildSvrPasswd string

	RedisIp      string
	RedisPort    string
	RedisPwd     string
	RedisDB      int
	RedisProtocl int
	RedisTimeout int
}

func NewConfig() (*Config, error) {
	cfgFile, err := ini.Load("env.ini")
	if err != nil {
		return nil, err
	}
	return initConfig(cfgFile)
}

func initConfig(cfgFile *ini.File) (*Config, error) {
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

	volumeDir := cfgFile.Section("DIR").Key("VOLUME_DIR").MustString("/tmp/volume")
	config.VolumeDir = filepath.FromSlash(volumeDir)
	os.MkdirAll(config.VolumeDir, 0755)

	config.GitId = cfgFile.Section("GITAUTH").Key("ID").MustString("")
	config.GitPwd = cfgFile.Section("GITAUTH").Key("PASSWD").MustString("")
	if config.GitId == "" || config.GitPwd == "" {
		return nil, fmt.Errorf("[CONFIG] git auth info is empty")
	}

	config.DockerRepoIp = cfgFile.Section("DOCKER").Key("DOCKER_REPO_IP").MustString("")
	config.DockerRepoPort = cfgFile.Section("DOCKER").Key("DOCKER_REPO_PORT").MustString("")
	if config.DockerRepoIp == "" || config.DockerRepoPort == "" {
		return nil, fmt.Errorf("[CONFIG] docker repo info is empty")
	}

	config.DockerRepoId = cfgFile.Section("DOCKER").Key("DOCKER_REPO_ID").MustString("")
	config.DockerRepoPwd = cfgFile.Section("DOCKER").Key("DOCKER_REPO_PASSWD").MustString("")
	if config.DockerRepoId == "" || config.DockerRepoPwd == "" {
		return nil, fmt.Errorf("[CONFIG] docker repo auth info is empty")
	}

	config.BuildSvrPasswd = cfgFile.Section("BUILD_SERVER").Key("BUILD_SERVER_PASSWD").MustString("")
	if config.BuildSvrPasswd == "" {
		return nil, fmt.Errorf("[CONFIG] build server passwd is empty")
	}

	config.ShellName = cfgFile.Section("SHELL").Key("SHELL_NAME").MustString("bash")

	config.RedisIp = cfgFile.Section("REDIS").Key("IP").MustString("")
	config.RedisPort = cfgFile.Section("REDIS").Key("PORT").MustString("")
	config.RedisPwd = cfgFile.Section("REDIS").Key("PASSWD").MustString("")
	if config.RedisIp == "" || config.RedisPort == "" || config.RedisPwd == "" {
		return nil, fmt.Errorf("[CONFIG] redis info is empty")
	}

	config.RedisDB = cfgFile.Section("REDIS").Key("DB").MustInt(0)
	config.RedisProtocl = cfgFile.Section("REDIS").Key("PROTOCOL").MustInt(2)
	config.RedisTimeout = cfgFile.Section("REDIS").Key("TIMEOUT").MustInt(5)

	return config, nil
}
