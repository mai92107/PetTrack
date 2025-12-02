package repo

import (
	domain "PetTrack/domain/repo"
	domainRepo "PetTrack/domain/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *deviceRepoImpl) Create(ctx context.Context, deviceType string, memberId int64, deviceId string) error {
	device := &domain.Device{
		Uuid:           uuid.New(),
		DeviceId:       deviceId,
		DeviceType:     deviceType,
		CreateByMember: memberId,
	}
	err := r.write.WithContext(ctx).Table("device").Create(&device).Error
	if err != nil {
		// logafa.Error("建立裝置資料失敗, error: %+v", err)
		return fmt.Errorf("建立裝置資料失敗")
	}
	return nil
}

func (r *deviceRepoImpl) GetDeviceList(ctx context.Context) ([]string, error) {
	var deviceIds []string
	err := r.read.WithContext(ctx).Model(&domainRepo.Device{}).
		Pluck("device_id", &deviceIds).Error
	if err != nil {
		// logafa.Error("查詢所有 deviceIds 失敗, error: %+v", err)
		return nil, fmt.Errorf("裝置ID查詢失敗")
	}
	return deviceIds, nil
}

func (r *deviceRepoImpl) FindDeviceById(ctx context.Context, deviceId string) (domainRepo.Device, error) {
	device := domainRepo.Device{}
	err := r.read.WithContext(ctx).First(&device, "device_id = ?", deviceId).Error
	if err != nil {
		// logafa.Error("查無此裝置, error: %+v", err)
		return device, fmt.Errorf("查無此裝置")
	}
	return device, nil
}
