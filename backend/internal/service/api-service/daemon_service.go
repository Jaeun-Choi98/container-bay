package apiservice

import (
	"fmt"

	"github.com/Jaeun-Choi98/container-bay/internal/redis"
	redismodel "github.com/Jaeun-Choi98/container-bay/internal/redis/redis-model"
	"github.com/Jaeun-Choi98/container-bay/internal/transport/http/rest/request"
)

func (s *ApiService) GetDaemons() ([]*redismodel.DockerDaemon, error) {
	repo := redis.GetRepository[*redismodel.DockerDaemon](redismodel.DockerDaemonKey)
	if repo == nil {
		return nil, fmt.Errorf("[API Service] daemon repository not available")
	}
	daemons, err := repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("[API Service] failed to list daemons: %w", err)
	}
	return daemons, nil
}

func (s *ApiService) AddDaemon(req *request.PostAddDaemonRequest) error {
	repo := redis.GetRepository[*redismodel.DockerDaemon](redismodel.DockerDaemonKey)
	if repo == nil {
		return fmt.Errorf("[API Service] daemon repository not available")
	}
	return repo.Create(&redismodel.DockerDaemon{Host: req.Host, Label: req.Label})
}

func (s *ApiService) RemoveDaemon(id int64) error {
	repo := redis.GetRepository[*redismodel.DockerDaemon](redismodel.DockerDaemonKey)
	if repo == nil {
		return fmt.Errorf("[API Service] daemon repository not available")
	}
	return repo.Delete(id)
}
