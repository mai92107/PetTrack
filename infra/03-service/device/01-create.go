package device

import (
	"context"
	"fmt"
)

func (s *DeviceServiceImpl) Create(ctx context.Context, deviceType string, memberId int64) (string, error) {
	deviceId := ""
	if err := validateRequest(deviceType); err != nil {
		return deviceId, err
	}
	deviceId = s.redisService.GenerateDeviceId(ctx)
	// 取得用戶資料
	err := s.deviceRepo.Create(ctx, deviceType, memberId, deviceId)
	return deviceId, err
}

func validateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
