package redismodel

import "time"

type DockerLoginSession struct {
	Id      int64
	IsLogin bool
	Url     string
}

func (d *DockerLoginSession) GetId() int64      { return d.Id }
func (d *DockerLoginSession) SetId(id int64)    { d.Id = id }
func (d *DockerLoginSession) TableName() string { return "dockerlogin" }
func (d *DockerLoginSession) GetIndexFields() map[string]interface{} {
	return map[string]interface{}{
		"url": d.Url,
	}
}
func (d *DockerLoginSession) GetTTL() time.Duration {
	return 1 * time.Minute
}
