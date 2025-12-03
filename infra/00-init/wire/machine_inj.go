package wire

import (
	initMethod "PetTrack/domain/init"
	"PetTrack/infra/00-core/model"
	bun "PetTrack/infra/00-core/model/bunMachines"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var machineSet = wire.NewSet(
	provideDB,
	provideMongo,
	provideRedis,
	provideMqtt,
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

func provideMqtt(cfg model.Config) *mqtt.Client {
	return initMethod.InitMqtt(cfg)
}
