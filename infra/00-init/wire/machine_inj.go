package wire

import (
	initMethod "PetTrack/domain/init"
	"PetTrack/infra/00-core/model"
	bun "PetTrack/infra/00-core/model/bunMachines"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var machineSet = wire.NewSet(
	provideDB,
	provideMongo,
	provideRedis,
)

func provideDB(cfg model.Config) *bun.DB {
	return initMethod.InitDB(cfg)
}

func provideMongo(cfg model.Config) *mongo.Database {
	return initMethod.InitMongo(cfg)
}

func provideRedis(cfg model.Config) *redis.Client {
	return initMethod.InitRedis(cfg)
}
