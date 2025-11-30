package device

import (
	"context"
	"fmt"
)

func (s *DeviceServiceImpl) Create(deviceType string, memberId int64) (string, error) {
	deviceId := ""
	if err := validateRequest(deviceType); err != nil {
		return deviceId, err
	}
	// TODO: 暫時用background
	ctx := context.Background()
	deviceId = s.redisService.GenerateDeviceId(ctx)
	// 取得用戶資料
	err := s.deviceRepo.Create(deviceType, memberId, deviceId)
	return deviceId, err
}

func validateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
