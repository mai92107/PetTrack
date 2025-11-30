package device

func (s *DeviceServiceImpl) DeviceList() ([]string, error) {
	deviceIds := []string{}
	deviceIds, err := s.deviceRepo.GetDeviceList()
	if err != nil {
		return deviceIds, err
	}
	return deviceIds, nil
}
