package redismodel

import "time"

type DockerLoginSession struct {
	Id      int64  `json:"id"`
	IsLogin bool   `json:"is_login"`
	Url     string `json:"url"`
}

func (d *DockerLoginSession) GetId() int64      { return d.Id }
func (d *DockerLoginSession) SetId(id int64)    { d.Id = id }
func (d *DockerLoginSession) TableName() string { return "dockerlogin" }
func (d *DockerLoginSession) GetIndexFields() map[string]any {
	return map[string]any{
		"url": d.Url,
	}
}
func (d *DockerLoginSession) GetTTL() time.Duration {
	return 1 * time.Minute
}
