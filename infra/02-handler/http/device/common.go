package device

import "PetTrack/infra/03-service/device"

var deviceService device.DeviceServiceImpl

func InitDeviceHandler(service device.DeviceServiceImpl) {
	deviceService = service
}