package service

import (
	"PetTrack/core/model"
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
	"context"
	"fmt"
	"slices"
)

type CommonServiceImpl struct {
	deviceRepo       domainRepo.DeviceRepository
	memberDeviceRepo domainRepo.MemberDeviceRepository
}

func NewCommonService(
	deviceRepo domainRepo.DeviceRepository,
	memberDeviceRepo domainRepo.MemberDeviceRepository,
) domainService.CommonService {
	return &CommonServiceImpl{
		deviceRepo: deviceRepo, memberDeviceRepo: memberDeviceRepo,
	}
}

func (s *CommonServiceImpl) ValidateDeviceOwner(ctx context.Context, deviceId string, member model.Claims) error {
	if member.Identity == "ADMIN" {
		return nil
	}

	deviceIds, err := s.memberDeviceRepo.GetMemberDeviceList(ctx, member.MemberId)
	if err != nil {
		return err
	}
	if !slices.Contains(deviceIds, deviceId) {
		// logafa.Debug("用戶 %v 嘗試讀取裝置 %s 資訊", member.MemberId, deviceId)
		return fmt.Errorf("無權限執行此操作")
	}
	return nil
}
