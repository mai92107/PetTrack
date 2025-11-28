package cache

import "context"

func (myCache *redisServiceImpl) HSetData(ctx context.Context, key string, mapData map[string]interface{}) error {
	err := myCache.client.HSet(ctx, key, mapData).Err()
	if err != nil {
		// logafa.Error("Redis HSet 寫入失敗, key: %s, data: %+v", key, mapData)
	}
	return err
}

func (myCache *redisServiceImpl) HSetFieldData(ctx context.Context, key, field, value string) error {
	err := myCache.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		// logafa.Error("Redis HSetFieldData 寫入失敗, key: %s, field: %s, value: %s", key, field, value)
	}
	return err
}

func (myCache *redisServiceImpl) HGetData(ctx context.Context, key, field string) string {
	value, err := myCache.client.HGet(ctx, key, field).Result()
	if err != nil {
		// logafa.Error("Redis HGet 讀取失敗, key: %s, field: %s, error: %+v", key, field, err)
	}
	return value
}

func (myCache *redisServiceImpl) HGetAllData(ctx context.Context, key string) map[string]string {
	value, err := myCache.client.HGetAll(ctx, key).Result()
	if err != nil {
		// logafa.Error("Redis HGetAll 讀取失敗, key: %s, error: %+v", key, err)
	}
	return value
}
