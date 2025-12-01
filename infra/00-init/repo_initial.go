package initial

import (
	bun "PetTrack/core/model/bunMachines"
	domainRepo "PetTrack/domain/repo"
	repo "PetTrack/infra/04-repository"
	cache "PetTrack/infra/04-repository/redisUtilImpl"

	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Account      domainRepo.AccountRepository
	Member       domainRepo.MemberRepository
	Password     domainRepo.PasswordRepository
	Device       domainRepo.DeviceRepository
	Trip         domainRepo.TripRepository
	MemberDevice domainRepo.MemberDeviceRepository
	Redis        domainRepo.RedisRepository
}

func InitRepositories(db *bun.DB, redis *redis.Client) *Repositories {
	return &Repositories{
		Account:      repo.NewAccountRepository(db),
		Member:       repo.NewMemberRepository(db),
		Password:     repo.NewPasswordRepository(db),
		Device:       repo.NewDeviceRepository(db),
		Trip:         repo.NewTripRepository(db),
		MemberDevice: repo.NewMemberDeviceRepository(db),
		Redis:        cache.NewRedisClient(redis),
	}
}
