package initMethod

import (
	"PetTrack/infra/00-core/model"
	"PetTrack/infra/00-core/util/logafa"
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg model.Config) *redis.Client {
	redis, err := initRedis(
		cfg.Machines.Redis.Host,
		cfg.Machines.Redis.Port,
		cfg.Machines.Redis.Password,
		0,
	)
	if err != nil {
		return nil
	}
	return redis
}
func initRedis(host, port, password string, db int) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	// TODO: 增加 retry
	if err := client.Ping(context.Background()).Err(); err != nil {
		logafa.Error(" ❌ Redis Read Client 連線失敗", "error", err)
		return nil, err
	}

	// logafa.Debug(" ✅ Redis 資料庫連接成功")
	return client, nil
}
