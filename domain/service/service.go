package domainService

import (
	"PetTrack/core/model"
	"context"
)

type AccountService interface {
	Login(ip, accountName, password string) (map[string]interface{}, error)
	Register(ip, username, password, email, lastName, firstName, nickName string) (map[string]interface{}, error)
}

type DeviceService interface {
	Create(deviceType string, memberId int64) (string, error)
	OnlineDeviceList() ([]string, error)
	DeviceList() ([]string, error)
}

type TripService interface {
	GetDeviceTrips(member model.Claims, deviceId string, pageable model.Pageable) ([]map[string]interface{}, int64, int64, error)
	GetTripDetail(member model.Claims, deviceId string, tripUuid string) (map[string]interface{}, error)
	SaveGpsFmRdsToMongo()
	FlushTripFmMongoToMaria(timeDuration int)
}

type MemberService interface {
	AddDevice(memberId int64, deviceId, deviceName string) error
	MemberDeviceList(memberId int64) ([]string, error)
}

type CommonService interface {
	ValidateDeviceOwner(deviceId string, member model.Claims) error
}

type RedisService interface {
	InitDeviceSequence(ctx context.Context)
	GenerateDeviceId(ctx context.Context) (string, error)
	GetOnlineDeviceList(ctx context.Context) ([]string, error)
}