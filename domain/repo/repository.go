package domainRepo

import (
	"PetTrack/core/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByAccountName(ctx context.Context, accountName string) (*Account, error)
	UpdateLoginTime(ctx context.Context, uuid uuid.UUID) error
	Create(tx *gorm.DB, ctx context.Context, memberId int64, username, password, email string) (uuid.UUID, error)
}

type MemberRepository interface {
	CreateMember(tx *gorm.DB, ctx context.Context, lastname, firstname, nickname, email string) (int64, error)
}

type PasswordRepository interface {
	CreateHistory(tx *gorm.DB, ctx context.Context, accountUuid uuid.UUID, password string) error
}

type DeviceRepository interface {
	Create(ctx context.Context, deviceType string, memberId int64, deviceId string) error
	GetDeviceList(ctx context.Context) ([]string, error)
	FindDeviceById(ctx context.Context, deviceId string) (Device, error)
}

type TripRepository interface {
	GetDeviceTrips(deviceId string, pageable model.Pageable) ([]TripSummary, int64, int64, error)
	GetTripDetail(tripUuid string) (TripSummary, error)
	SaveLocationToDB(records []DeviceLocation) error
	ReadTripFromMongo(timeDuration int) ([]TripSummary, error)
	SaveTripToDB(results []TripSummary) error
}

type MemberDeviceRepository interface {
	AddDevice(ctx context.Context, memberId int64, deviceId, deviceName string) error
	GetMemberDeviceList(ctx context.Context, memberId int64) ([]string, error)
}
