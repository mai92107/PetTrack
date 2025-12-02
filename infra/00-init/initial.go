package initial

import (
	initMethod "PetTrack/domain/init"
	"PetTrack/infra/00-core/global"
	bun "PetTrack/infra/00-core/model/bunMachines"
	"PetTrack/infra/00-core/util/logafa"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init() {
	InitLogger()
	InitWorkers()

	db := InitDB()
	mongo := InitMongo()
	redis := InitRedis()
	repos := InitRepositories(db, redis, mongo)
	services := InitServices(repos, db, redis)

	InitHandlers(services)
	InitCron(services)
}

func InitWorkers() {
	maxPriorWorkers := 20
	maxNormalWorkers := 50
	// 區隔工人 做 故障隔離
	// 高級勞工
	global.PriorWorkerPool = make(chan struct{}, maxPriorWorkers)
	for range maxPriorWorkers {
		global.PriorWorkerPool <- struct{}{}
	}
	logafa.Debug("高級勞工 初始化完成", "數量", maxPriorWorkers)
	// 城市打工人
	global.NormalWorkerPool = make(chan struct{}, maxNormalWorkers)
	for range maxNormalWorkers {
		global.NormalWorkerPool <- struct{}{}
	}
	logafa.Debug("城市打工人 初始化完成", "數量", maxNormalWorkers)
}

func InitLogger() {
	logafa.CreateLogFileNow()

	// 初始化（全專案只需要呼叫一次）
	handler := logafa.NewLogafaHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	slog.SetDefault(slog.New(handler))
	logafa.Debug("Logafa 初始化完成")

}

func InitDB() *bun.DB {
	db, err := initMethod.InitDB("localhost", "3306", "bunbun", "bunbun", "pettrack")
	if err != nil {
		panic(err)
	}
	logafa.Debug("DB 初始化完成")
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
	logafa.Debug("MONGO 初始化完成")
	return mongo
}

func InitRedis() *redis.Client {
	redis, err := initMethod.InitRedis("localhost", "6379", "", 0)
	if err != nil {
		panic(err)
	}
	logafa.Debug("REDIS 初始化完成")
	return redis
}
