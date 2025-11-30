package service

import (
	"PetTrack/core/global"
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type redisServiceImpl struct {
	redis     *redis.Client
	redisUtil domainRepo.RedisRepository
}

func NewRedisRepository(
	client *redis.Client,
	redis domainRepo.RedisRepository,
) domainService.RedisService {
	return &redisServiceImpl{redis: client, redisUtil: redis}
}

func (s *redisServiceImpl) InitDeviceSequence(ctx context.Context) {
	device_setting := s.redisUtil.HGetAllData(ctx, "device_setting")

	prefix := "AA"
	seq := 0

	if len(device_setting) == 0 {
		s.redisUtil.HSetData(ctx, "device_setting",
			map[string]interface{}{
				"device_prefix":   prefix,
				"device_sequence": seq,
			})
	}
	// logafa.Debug(" ✅ 成功初始化裝置設定")
}

func (s *redisServiceImpl) GenerateDeviceId(ctx context.Context) string {
	prefix := s.redisUtil.HGetData(ctx, "device_setting", "device_prefix")
	sequence, err := s.redis.HIncrBy(ctx, "device_setting", "device_sequence", 1).Result()
	if err != nil {
		// logafa.Error("failed to increment sequence in Redis: %v", err)
		return ""
	}
	return fmt.Sprintf("%s-%06d", prefix, sequence)
}

func (s *redisServiceImpl) GetOnlineDeviceList(ctx context.Context) ([]string, error) {
	keys, err := s.redisUtil.KeyScan(ctx, "device:*")
	if err != nil {
		// logafa.Error("redis 掃描 device:* 失敗: %v", err)
		return nil, fmt.Errorf("%s: redis scan error", global.COMMON_SYSTEM_ERROR)
	}
	deviceIds := make([]string, 0, len(keys))
	for _, key := range keys {
		if !strings.HasPrefix(key, "device:") {
			continue // 防呆
		}
		parts := strings.SplitN(key, ":", 2) // 只切一次
		if len(parts) == 2 {
			deviceIds = append(deviceIds, parts[1])
		}
	}
	// logafa.Info("目前在線裝置數量: %d", len(deviceIds))
	return deviceIds, nil
}
