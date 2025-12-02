package cache

import (
	"PetTrack/core/util/logafa"
	"context"
)

func (myCache *redisRepoImpl) HSetData(ctx context.Context, key string, mapData map[string]interface{}) error {
	err := myCache.client.HSet(ctx, key, mapData).Err()
	if err != nil {
		logafa.Error("Redis HSet 寫入失敗", "key", key, "data", mapData, "error", err)
		return err
	}
	return err
}

func (myCache *redisRepoImpl) HSetFieldData(ctx context.Context, key, field, value string) error {
	err := myCache.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		logafa.Error("Redis HSetFieldData 寫入失敗", "key", key, "field", field, "value", value, "error", err)
		return err
	}
	return nil
}

func (myCache *redisRepoImpl) HGetData(ctx context.Context, key, field string) (string, error) {
	value, err := myCache.client.HGet(ctx, key, field).Result()
	if err != nil {
		logafa.Error("Redis HGet 讀取失敗", "key", key, "field", field, "error", err)
		return "", err
	}
	return value, nil
}

func (myCache *redisRepoImpl) HGetAllData(ctx context.Context, key string) (map[string]string, error) {
	value, err := myCache.client.HGetAll(ctx, key).Result()
	if err != nil {
		logafa.Error("Redis HGetAll 讀取失敗", "key", key, "error", err)
		return nil, err
	}
	return value, nil
}
