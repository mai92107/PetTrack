package initial

import (
	initMethod "PetTrack/domain/init"
	bun "PetTrack/infra/00-core/model/bunMachines"
	"PetTrack/infra/00-core/util/logafa"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init() {

	InitLogger()

	db := InitDB()
	mongo := InitMongo()
	redis := InitRedis()
	repos := InitRepositories(db, redis, mongo)
	services := InitServices(repos, db, redis)
	InitHandlers(services)

	InitCron(services)
}

func InitLogger() {
	logafa.CreateLogFileNow()

	// 初始化（全專案只需要呼叫一次）
	handler := logafa.NewLogafaHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // 關鍵！讓 slog 自動填正確的 caller
	})

	slog.SetDefault(slog.New(handler))
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
