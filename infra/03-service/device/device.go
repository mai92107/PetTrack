package device

import (
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
)

type DeviceServiceImpl struct {
	deviceRepo    domainRepo.DeviceRepository
	tripRepo      domainRepo.TripRepository
	redisService  domainService.RedisService
}

func NewDeviceService(
	deviceRepo domainRepo.DeviceRepository,
	tripRepo domainRepo.TripRepository,
	redisService domainService.RedisService,
) domainService.DeviceService {
	return &DeviceServiceImpl{
		deviceRepo:    deviceRepo,
		tripRepo:      tripRepo,
		redisService:  redisService,
	}
}
