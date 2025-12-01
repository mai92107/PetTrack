package initial

import (
	bun "PetTrack/core/model/bunMachines"
	"PetTrack/core/util/logafa"
	initMethod "PetTrack/domain/init"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init() {

	InitLogger()

	db := InitDB()
	_ = InitMongo()
	redis := InitRedis()
	repos := InitRepositories(db, redis)
	services := InitServices(repos, db, redis)
	InitHandlers(services)
}

func InitLogger() *slog.Logger {
	// log.CreateLogFileNow()

	level := slog.LevelDebug
	logHandler := logafa.NewColorHandler(level, os.Stdout)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return logger
}

func InitDB() *bun.DB {
	db, err := initMethod.InitDB("localhost", "3306", "bunbun", "bunbun", "pettrack")
	if err != nil {
		panic(err)
	}
	return &bun.DB{
		Write: db,
		Read:  db,
	}
}

func InitMongo() *mongo.Database {
	mongo, err := initMethod.InitMongo("localhost", "27017", "bunbun", "bunbun")
	if err != nil {
		panic(err)
	}
	return mongo
}

func InitRedis() *redis.Client {
	redis, err := initMethod.InitRedis("localhost", "6379", "", 0)
	if err != nil {
		panic(err)
	}
	return redis
}
