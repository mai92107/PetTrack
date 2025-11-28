package device

var deviceService device.DeviceServiceImpl

func InitDeviceHandler(service device.DeviceServiceImpl) {
	deviceService = service
}