package repo

import (
	domain "PetTrack/domain/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *deviceRepoImpl) Create(deviceType string, memberId int64, deviceId string) (string, error) {
	device := domain.Device{
		Uuid:           uuid.New(),
		DeviceId:       deviceId,
		DeviceType:     deviceType,
		CreateByMember: memberId,
	}
	err := r.db.Table("device").Create(&device).Error
	if err != nil {
		// logafa.Error("建立裝置資料失敗, error: %+v", err)
		return "", fmt.Errorf("建立裝置資料失敗")
	}
	return device.DeviceId, nil
}

func (r *deviceRepoImpl) GenerateDeviceId(ctx *context.Context) string {
	prefix := r.HGetData("device_setting", "device_prefix")
	sequence, err := r.redis.HIncrBy(ctx, "device_setting", "device_sequence", 1).Result()
	if err != nil {
		// logafa.Error("failed to increment sequence in Redis: %v", err)
		return ""
	}
	return fmt.Sprintf("%s-%06d", prefix, sequence)
}
