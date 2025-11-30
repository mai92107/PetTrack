package cache

import (
	domainRepo "PetTrack/domain/repo"

	"github.com/redis/go-redis/v9"
)

type redisRepoImpl struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) domainRepo.RedisRepository {
	return &redisRepoImpl{client: client}
}
