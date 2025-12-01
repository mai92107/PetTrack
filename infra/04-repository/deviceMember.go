package repo

import (
	domainRepo "PetTrack/domain/repo"
	"fmt"
)

func (r *memberDeviceRepoImpl) AddDevice(memberId int64, deviceId, deviceName string) error {

	memberDevice := domainRepo.MemberDevice{
		MemberId:   memberId,
		DeviceId:   deviceId,
		DeviceName: deviceName,
	}
	err := r.write.Create(&memberDevice).Error
	if err != nil {
		// logafa.Error("新增使用者裝置失敗, error: %+v", err)
		return fmt.Errorf("新增使用者裝置失敗")
	}
	return nil
}

func (r *memberDeviceRepoImpl) GetMemberDeviceList(memberId int64) ([]string, error) {
	var deviceIds []string
	err := r.read.Model(&domainRepo.MemberDevice{}).
		Pluck("device_id", &deviceIds).
		Where("member_id = ?", memberId).Error
	if err != nil {
		// logafa.Error("查詢所有 deviceIds 失敗, error: %+v", err)
		return nil, fmt.Errorf("裝置ID查詢失敗")
	}
	return deviceIds, nil
}
