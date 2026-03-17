package infra

import "github.com/zeromicro/go-zero/core/stores/redis"

func Redis(config redis.RedisConf) (*redis.Redis, error) {
	redisClient, err := redis.NewRedis(config)
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
