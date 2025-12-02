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
	deviceId, err := s.redisService.GenerateDeviceId(ctx)
	if err != nil {
		return deviceId, err
	}
	// 取得用戶資料
	err = s.deviceRepo.Create(ctx, deviceType, memberId, deviceId)
	if err != nil {
		return deviceId, err
	}
	return deviceId, nil
}

func validateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
