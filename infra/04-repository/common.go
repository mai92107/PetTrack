package repo

import (
	domain "PetTrack/domain/repo"
	domainRepo "PetTrack/domain/repo"

	"gorm.io/gorm"
)

type accountRepoImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) domain.AccountRepository {
	return &accountRepoImpl{db: db}
}

type memberRepoImpl struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) domain.MemberRepository {
	return &memberRepoImpl{db: db}
}

type passwordRepoImpl struct {
	db *gorm.DB
}

func NewPasswordRepository(db *gorm.DB) domain.PasswordRepository {
	return &passwordRepoImpl{db: db}
}

type deviceRepoImpl struct {
	db *gorm.DB
}

func NewDeviceRepository(
	db *gorm.DB,
) domainRepo.DeviceRepository {
	return &deviceRepoImpl{db: db}
}

type tripRepoImpl struct {
	db *gorm.DB
}

func NewTripRepository(
	db *gorm.DB,
) domainRepo.TripRepository {
	return &tripRepoImpl{db: db}
}

type memberDeviceRepoImpl struct {
	db *gorm.DB
}

func NewMemberDeviceRepository(
	db *gorm.DB,
) domainRepo.MemberDeviceRepository {
	return &memberDeviceRepoImpl{db: db}
}
