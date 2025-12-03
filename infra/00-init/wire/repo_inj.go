package wire

import (
	domainRepo "PetTrack/domain/repo"
	bun "PetTrack/infra/00-core/model/bunMachines"
	repo "PetTrack/infra/04-repository"
	cache "PetTrack/infra/04-repository/redisUtilImpl"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var repoSet = wire.NewSet(
	provideAccountRepo,
	provideMemberRepo,
	providePasswordRepo,
	provideDeviceRepo,
	provideTripRepo,
	provideMemberDeviceRepo,
	provideRedisRepo,
)

func provideAccountRepo(db *bun.DB) domainRepo.AccountRepository {
	return repo.NewAccountRepository(db)
}
func provideMemberRepo(db *bun.DB) domainRepo.MemberRepository {
	return repo.NewMemberRepository(db)
}
func providePasswordRepo(db *bun.DB) domainRepo.PasswordRepository {
	return repo.NewPasswordRepository(db)
}
func provideDeviceRepo(db *bun.DB) domainRepo.DeviceRepository {
	return repo.NewDeviceRepository(db)
}
func provideTripRepo(db *bun.DB, mongo *mongo.Database) domainRepo.TripRepository {
	return repo.NewTripRepository(db, mongo)
}
func provideMemberDeviceRepo(db *bun.DB) domainRepo.MemberDeviceRepository {
	return repo.NewMemberDeviceRepository(db)
}
func provideRedisRepo(redis *redis.Client) domainRepo.RedisRepository {
	return cache.NewRedisClient(redis)
}
