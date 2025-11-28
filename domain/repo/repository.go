package domain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByAccountName(accountName string) (*Account, error)
	UpdateLoginTime(uuid uuid.UUID) error
	Create(tx *gorm.DB, memberId int64, username, password, email string) (uuid.UUID, error)
}

type MemberRepository interface {
	CreateMember(tx *gorm.DB, lastname, firstname, nickname, email string) (int64, error)
}

type PasswordRepository interface {
	CreateHistory(tx *gorm.DB, accountUuid uuid.UUID, password string) error
}

type DeviceRepository interface {
	Create(deviceType string, memberId int64, deviceId string) (string, error)
	GenerateDeviceId(ctx *context.Context) string
}
