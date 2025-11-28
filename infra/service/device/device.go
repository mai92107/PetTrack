package device

import (
	domain "PetTrack/domain/repo"
)

type DeviceServiceImpl struct {
	deviceRepo domain.DeviceRepository
	redisService domain.
}

func NewDeviceService(deviceRepo domain.DeviceRepository) domain. {
	return &DeviceServiceImpl{
		deviceRepo: deviceRepo,
	}
}
