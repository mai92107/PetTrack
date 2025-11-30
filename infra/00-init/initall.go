package initial

import (
	initMethod "PetTrack/domain/init"
	handler "PetTrack/infra/02-handler/http"
	service "PetTrack/infra/03-service"
	"PetTrack/infra/03-service/account"
	"PetTrack/infra/03-service/device"
	"PetTrack/infra/03-service/member"
	repo "PetTrack/infra/04-repository"
	cache "PetTrack/infra/04-repository/redisUtilImpl"
)

func Init() {
	db, err := initMethod.InitDB("localhost", "3306", "bunbun", "bunbun", "pettrack")
	if err != nil {

		return
	}
	redis, err := initMethod.InitRedis("localhost", "6379", "", 0)
	if err != nil {

		return
	}
	// mongo, err := initMethod.InitMongo("localhost", "27017", "bunbun", "bunbun")
	// if err != nil {

	// 	return
	// }
	accountRepo := repo.NewAccountRepository(db)
	memberRepo := repo.NewMemberRepository(db)
	passwordRepo := repo.NewPasswordRepository(db)
	deviceRepo := repo.NewDeviceRepository(db)
	tripRepo := repo.NewTripRepository(db)
	memberDeviceRepo := repo.NewMemberDeviceRepository(db)
	redisRepo := cache.NewRedisClient(redis)

	redisService := service.NewRedisRepository(redis, redisRepo)
	commonService := service.NewCommonService(deviceRepo, memberDeviceRepo)
	accountService := account.NewAccountService(db, passwordRepo, accountRepo, memberRepo)
	deviceService := device.NewDeviceService(deviceRepo, tripRepo, redisService, commonService)
	memberService := member.NewMemberService(deviceRepo, memberDeviceRepo)

	handler.InitAccountHandler(accountService)
	handler.InitDeviceHandler(deviceService)
	handler.InitMemberHandler(memberService)
}
