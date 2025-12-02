package device

import "context"

func (s *DeviceServiceImpl) DeviceList(ctx context.Context) ([]string, error) {
	deviceIds := []string{}
	deviceIds, err := s.deviceRepo.GetDeviceList(ctx)
	if err != nil {
		return deviceIds, err
	}
	return deviceIds, nil
}
