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
	deviceId, err := s.redisService.GenerateDeviceId(ctx)
	if err != nil {
		return "", err
	}
	err = s.deviceRepo.Create(deviceType, memberId, deviceId)
	if err != nil {
		return "", err
	}
	return deviceId, nil
}

func validateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
