package device

import "context"

func (s *DeviceServiceImpl) OnlineDeviceList() ([]string, error) {
	deviceIds := []string{}
	// TODO: 暫時用background
	ctx := context.Background()
	deviceIds, err := s.redisService.GetOnlineDeviceList(ctx)
	return deviceIds, err
}
