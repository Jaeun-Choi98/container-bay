package apiservice

import (
	"context"
	"fmt"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/config"
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	"github.com/Jaeun-Choi98/container-bay/internal/redis"
	redismodel "github.com/Jaeun-Choi98/container-bay/internal/redis/redis-model"
	"github.com/Jaeun-Choi98/modules/shell"
)

type ApiService struct {
	Cfg *config.Config
}

func NewApiService(cfg *config.Config) *ApiService {
	if err := Init(cfg); err != nil {
		logger.Printf("[API Service] failed to init: %v", err)
	}
	return &ApiService{
		Cfg: cfg,
	}
}

func Init(cfg *config.Config) error {

	dockerLoginSessionRepo := redis.GetRepository[*redismodel.DockerLoginSession](redismodel.DockerLoginSessionKey)
	if dockerLoginSessionRepo == nil {
		return fmt.Errorf("failed to get docker login session info")
	}

	dockerLoginSession, err := dockerLoginSessionRepo.
		FindByIndex("url", fmt.Sprintf("%s:%s", cfg.DockerRepoIp, cfg.DockerRepoPort))

	// 해당 도서 사설 레포의 로그인 세션이 없다면, 로그인 이후에 레디스에 데이터 갱신
	if err != nil || dockerLoginSession == nil {
		// if err := DockerRepoLogin(cfg); err != nil {
		// 	return fmt.Errorf("fail docker login")
		// }
		if err := dockerLoginSessionRepo.Create(&redismodel.DockerLoginSession{
			Url:     fmt.Sprintf("%s:%s", cfg.DockerRepoIp, cfg.DockerRepoPort),
			IsLogin: true,
		}); err != nil {
			return fmt.Errorf("failed to create docker login session: %v", err)
		}
		loginSession, err := dockerLoginSessionRepo.FindByIndex("url", fmt.Sprintf("%s:%s", cfg.DockerRepoIp, cfg.DockerRepoPort))
		if err != nil {
			return fmt.Errorf("failed to find login session after created docker login session: %v", err)
		}
		exp, _ := dockerLoginSessionRepo.GetTTL(loginSession.GetId())
		logger.Printf("[API Service] create docker login session in redis, expire: %v", exp)
	} else {
		exp, _ := dockerLoginSessionRepo.GetTTL(dockerLoginSession.GetId())
		logger.Printf("[API Service] already exist docker login session in redis, expire: %v", exp)
	}

	return nil
}

func DockerRepoLogin(cfg *config.Config) error {
	executor := shell.NewScriptExecutor(cfg.ShellDir)

	script := fmt.Sprintf(`
		echo %s | sudo -S sh -c 'echo %s | docker login %s:%s -u %s --password-stdin'
	`, cfg.BuildSvrPasswd, cfg.DockerRepoPwd, cfg.DockerRepoIp, cfg.DockerRepoPort, cfg.DockerRepoId)

	timeout, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	result, err := executor.Execute(timeout, cfg.ShellName, "docker_login.sh", script, true)

	if err != nil {
		return fmt.Errorf("[API Service] failed to execute script(docker_login.sh): %v", err)
	}
	if result.Error != nil {
		return fmt.Errorf("[API Service] failed to execute script(docker_login.sh): %v", err)
	}
	if len(result.Stderr) != 0 {
		logger.Println("=== STDERR ===")
		for _, line := range result.Stderr {
			logger.Println(line)
		}
	}
	return nil
}
