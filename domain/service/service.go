package domainService

import (
	"PetTrack/core/model"
	"context"
)

type AccountService interface {
	Login(ctx context.Context, ip, accountName, password string) (map[string]interface{}, error)
	Register(ctx context.Context, ip, username, password, email, lastName, firstName, nickName string) (map[string]interface{}, error)
}

type DeviceService interface {
	Create(ctx context.Context, deviceType string, memberId int64) (string, error)
	OnlineDeviceList(ctx context.Context) ([]string, error)
	DeviceList(ctx context.Context) ([]string, error)
}

type TripService interface {
	GetDeviceTrips(ctx context.Context, member model.Claims, deviceId string, pageable model.Pageable) ([]map[string]interface{}, int64, int64, error)
	GetTripDetail(ctx context.Context, member model.Claims, deviceId string, tripUuid string) (map[string]interface{}, error)
	SaveGpsFmRdsToMongo(ctx context.Context)
	FlushTripFmMongoToMaria(ctx context.Context, timeDuration int)
}

type MemberService interface {
	AddDevice(ctx context.Context, memberId int64, deviceId, deviceName string) error
	MemberDeviceList(ctx context.Context, memberId int64) ([]string, error)
}

type CommonService interface {
	ValidateDeviceOwner(ctx context.Context, deviceId string, member model.Claims) error
}

type RedisService interface {
	InitDeviceSequence(ctx context.Context)
	GenerateDeviceId(ctx context.Context) (string, error)
	GetOnlineDeviceList(ctx context.Context) ([]string, error)
}