package member

import (
	"context"
	"fmt"
)

func (s *MemberServiceImpl) AddDevice(ctx context.Context, memberId int64, deviceId, deviceName string) error {
	err := validateRequest(deviceId)
	if err != nil {
		return err
	}
	_, err = s.deviceRepo.FindDeviceById(ctx, deviceId)
	if err != nil {
		return err
	}
	err = s.memberDeviceRepo.AddDevice(ctx, memberId, deviceId, deviceName)
	if err != nil {
		return err
	}
	return nil
}

func validateRequest(deviceId string) error {
	if deviceId == "" {
		return fmt.Errorf("裝置識別碼不可為空")
	}
	return nil
}
