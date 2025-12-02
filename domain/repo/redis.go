package domainRepo

import "context"

type RedisRepository interface {
	HSetData(ctx context.Context, key string, mapData map[string]interface{}) error
	HSetFieldData(ctx context.Context, key, field, value string) error
	HGetData(ctx context.Context, key, field string) (string, error)
	HGetAllData(ctx context.Context, key string) (map[string]string, error)

	ZAddData(ctx context.Context, key string, score float64, byteData []byte) error
	KeyScan(ctx context.Context, pattern string) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, startTs, endTs int64) ([]string, error)
	ZRemRangeByScore(ctx context.Context, key string, startTs, endTs int64) error
}
