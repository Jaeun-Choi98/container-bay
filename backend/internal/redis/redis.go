package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Jaeun-Choi98/container-bay/internal/config"
	"github.com/Jaeun-Choi98/container-bay/internal/logger"
	redismodel "github.com/Jaeun-Choi98/container-bay/internal/redis/redis-model"
	"github.com/Jaeun-Choi98/modules/orm/redisorm"
)

var redisClient *redisorm.RedisClient
var repositories map[redismodel.RepoKey]any
var redisMu sync.RWMutex
var cleanupCtx context.Context
var cleanupCancel context.CancelFunc

func InitRedis(cfg *config.Config) error {
	redisClient = redisorm.NewRedisClient(
		fmt.Sprintf("%s:%s", cfg.RedisIp, cfg.RedisPort),
		cfg.RedisPwd,
		cfg.RedisDB,
		cfg.RedisProtocl,
		cfg.RedisTimeout,
	)
	if redisClient == nil {
		return fmt.Errorf("[REDIS] failed to connect redis")
	}
	repositories = make(map[redismodel.RepoKey]interface{})
	LoadDefaultRepo()
	go StartCleanupWorker(3 * time.Minute)
	return nil
}

func CloseRedisClient() error {
	if cleanupCancel != nil {
		cleanupCancel()
	}
	return redisClient.Close()
}

func AddRepository[T redisorm.Model](key redismodel.RepoKey, model T) error {
	redisMu.Lock()
	defer redisMu.Unlock()

	if _, exists := repositories[key]; !exists {
		repositories[key] = redisorm.NewRepository(redisClient, model)
	}
	return nil
}

func DeleteRepository(key redismodel.RepoKey) error {
	redisMu.Lock()
	defer redisMu.Unlock()
	delete(repositories, key)
	return nil
}

func GetRepository[T redisorm.Model](key redismodel.RepoKey) *redisorm.Repository[T] {
	redisMu.RLock()
	defer redisMu.RUnlock()

	repo, exists := repositories[key]
	if !exists {
		return nil
	}

	// 타입 단언
	typedRepo, ok := repo.(*redisorm.Repository[T])
	if !ok {
		return nil
	}

	return typedRepo
}

func LoadDefaultRepo() {
	AddRepository(redismodel.DockerLoginSessionKey, &redismodel.DockerLoginSession{})
}

func StartCleanupWorker(interval time.Duration) {
	cleanupCtx, cleanupCancel = context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		logger.Printf("[Redis] Cleanup worker started (interval: %v)", interval)

		for {
			select {
			case <-ticker.C:
				cleanupAllRepositories()
			case <-cleanupCtx.Done():
				logger.Println("[Redis] Cleanup worker routine is terminated")
				return
			}
		}
	}()
}

// 만료된 키의 인덱스 정리가 필요한 경우, 따로 코드를 아래처럼 작성해주어야 함.
func cleanupAllRepositories() {
	logger.Println("[Redis] Starting cleanup...")
	dockerRepo := GetRepository[*redismodel.DockerLoginSession](redismodel.DockerLoginSessionKey)
	if dockerRepo != nil {
		count, err := dockerRepo.CleanupExpired()
		if err != nil {
			logger.Printf("[Redis] Error: %v", err)
		} else if count > 0 {
			logger.Printf("[Redis] Cleaned %d records", count)
		}
	}
}
