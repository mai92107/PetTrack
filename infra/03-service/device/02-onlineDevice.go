package device

import "context"

func (s *DeviceServiceImpl) OnlineDeviceList(ctx context.Context) ([]string, error) {
	deviceIds := []string{}
	deviceIds, err := s.redisService.GetOnlineDeviceList(ctx)
	return deviceIds, err
}
