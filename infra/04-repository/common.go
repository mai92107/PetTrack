package repo

import (
	bun "PetTrack/infra/00-core/model/bunMachines"
	domainRepo "PetTrack/domain/repo"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type accountRepoImpl struct {
	write *gorm.DB
	read  *gorm.DB
}

func NewAccountRepository(
	db *bun.DB,
) domainRepo.AccountRepository {
	return &accountRepoImpl{write: db.Write, read: db.Read}
}

type memberRepoImpl struct {
	write *gorm.DB
	read  *gorm.DB
}

func NewMemberRepository(
	db *bun.DB,
) domainRepo.MemberRepository {
	return &memberRepoImpl{write: db.Write, read: db.Read}
}

type passwordRepoImpl struct {
	write *gorm.DB
	read  *gorm.DB
}

func NewPasswordRepository(
	db *bun.DB,
) domainRepo.PasswordRepository {
	return &passwordRepoImpl{write: db.Write, read: db.Read}
}

type deviceRepoImpl struct {
	write *gorm.DB
	read  *gorm.DB
}

func NewDeviceRepository(
	db *bun.DB,
) domainRepo.DeviceRepository {
	return &deviceRepoImpl{write: db.Write, read: db.Read}
}

type tripRepoImpl struct {
	write *gorm.DB
	read  *gorm.DB
	mongo *mongo.Database
}

func NewTripRepository(
	db *bun.DB,
	mongo *mongo.Database,
) domainRepo.TripRepository {
	return &tripRepoImpl{write: db.Write, read: db.Read, mongo: mongo}
}

type memberDeviceRepoImpl struct {
	write *gorm.DB
	read  *gorm.DB
}

func NewMemberDeviceRepository(
	db *bun.DB,
) domainRepo.MemberDeviceRepository {
	return &memberDeviceRepoImpl{write: db.Write, read: db.Read}
}
