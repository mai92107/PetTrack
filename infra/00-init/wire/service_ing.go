package wire

import (
	initMethod "PetTrack/domain/init"
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
	"PetTrack/infra/00-core/model"
	bun "PetTrack/infra/00-core/model/bunMachines"
	mqttHandler "PetTrack/infra/01-router/mqtt"
	service "PetTrack/infra/03-service"
	"PetTrack/infra/03-service/account"
	"PetTrack/infra/03-service/device"
	"PetTrack/infra/03-service/member"
	"PetTrack/infra/03-service/trip"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var serviceSet = wire.NewSet(
	provideRedisService,
	provideCommonService,
	provideAccountService,
	provideDeviceService,
	provideMemberService,
	provideTripService,

	provideMqttHandler,
	provideMqtt,

	InitServices,
)

type Services struct {
	Account domainService.AccountService
	Device  domainService.DeviceService
	Member  domainService.MemberService
	Common  domainService.CommonService
	Redis   domainService.RedisService
	Trip    domainService.TripService
	mqtt    *mqtt.Client
}

func InitServices(
	account domainService.AccountService,
	device domainService.DeviceService,
	member domainService.MemberService,
	common domainService.CommonService,
	redis domainService.RedisService,
	trip domainService.TripService,
	mqtt *mqtt.Client,
) *Services {
	return &Services{
		Account: account,
		Device:  device,
		Member:  member,
		Common:  common,
		Redis:   redis,
		Trip:    trip,
		mqtt:    mqtt,
	}
}

func provideRedisService(
	redisRepo domainRepo.RedisRepository,
	redis *redis.Client,
) domainService.RedisService {
	return service.NewRedisService(redis, redisRepo)
}

func provideCommonService(
	deviceRepo domainRepo.DeviceRepository,
	memberDeviceRepo domainRepo.MemberDeviceRepository,
) domainService.CommonService {
	return service.NewCommonService(deviceRepo, memberDeviceRepo)
}

func provideAccountService(
	db *bun.DB,
	passwordRepo domainRepo.PasswordRepository,
	accountRepo domainRepo.AccountRepository,
	memberRepo domainRepo.MemberRepository,
) domainService.AccountService {
	return account.NewAccountService(db, passwordRepo, accountRepo, memberRepo)
}
func provideDeviceService(
	deviceRepo domainRepo.DeviceRepository,
	tripRepo domainRepo.TripRepository,
	redisService domainService.RedisService,
) domainService.DeviceService {
	return device.NewDeviceService(deviceRepo, tripRepo, redisService)
}
func provideMemberService(
	deviceRepo domainRepo.DeviceRepository,
	memberDeviceRepo domainRepo.MemberDeviceRepository,
) domainService.MemberService {
	return member.NewMemberService(deviceRepo, memberDeviceRepo)
}
func provideTripService(
	tripRepo domainRepo.TripRepository,
	commonService domainService.CommonService,
	redisRepo domainRepo.RedisRepository,
) domainService.TripService {
	return trip.NewTripService(tripRepo, commonService, redisRepo)
}
func provideMqttHandler() mqtt.MessageHandler {
	return mqttHandler.NewMQTTMessageHandler()
}

func provideMqtt(cfg model.Config, handler mqtt.MessageHandler) *mqtt.Client {
	return initMethod.InitMqtt(cfg, handler)
}
