package device

import "fmt"

func (service *DeviceServiceImpl) Create(deviceType string, memberId int64) (string, error) {

	if err := validateRequest(deviceType); err != nil {
		return "", err
	}

	// 取得用戶資料
	deviceId, err := service.Create(deviceType, memberId)
	if err != nil {
		return "", fmt.Errorf("新增使用者裝置發生錯誤，error: %+v", err)
	}
	return deviceId, nil
}

func validateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
