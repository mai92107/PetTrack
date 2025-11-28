package domain

type AccountService interface {
	Login(ip, accountName, password string) (map[string]interface{}, error)
	Register(ip, username, password, email, lastName, firstName, nickName string) (map[string]interface{}, error)
}

type DeviceService interface {
	Create(deviceType string, memberId int64) (string, error)
}
