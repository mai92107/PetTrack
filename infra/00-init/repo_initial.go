package initial

import (
	domainRepo "PetTrack/domain/repo"
	bun "PetTrack/infra/00-core/model/bunMachines"
	repo "PetTrack/infra/04-repository"
	cache "PetTrack/infra/04-repository/redisUtilImpl"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
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

func InitRepositories(db *bun.DB, redis *redis.Client, mongo *mongo.Database) *Repositories {
	return &Repositories{
		Account:      repo.NewAccountRepository(db),
		Member:       repo.NewMemberRepository(db),
		Password:     repo.NewPasswordRepository(db),
		Device:       repo.NewDeviceRepository(db),
		Trip:         repo.NewTripRepository(db, mongo),
		MemberDevice: repo.NewMemberDeviceRepository(db),
		Redis:        cache.NewRedisClient(redis),
	}
}
