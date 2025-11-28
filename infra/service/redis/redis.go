package cache

import (
	"PetTrack/domain"

	"github.com/redis/go-redis/v9"
)

type redisServiceImpl struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) domain.RedisService {
	return &redisServiceImpl{client: client}
}
