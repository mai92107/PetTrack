package initMethod

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis(host, port, password string, db int) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	// TODO: 增加 retry
	if err := client.Ping(context.Background()).Err(); err != nil {
		// logafa.Error(" ❌ Redis Read Client 連線失敗: %v", err)
		return nil, err
	}

	// logafa.Debug(" ✅ Redis 資料庫連接成功")
	return client, nil
}
