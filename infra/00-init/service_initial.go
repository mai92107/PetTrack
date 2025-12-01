package initial

import (
	bun "PetTrack/core/model/bunMachines"
	domainService "PetTrack/domain/service"
	service "PetTrack/infra/03-service"
	"PetTrack/infra/03-service/account"
	"PetTrack/infra/03-service/device"
	"PetTrack/infra/03-service/member"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	Account domainService.AccountService
	Device  domainService.DeviceService
	Member  domainService.MemberService
	Common  domainService.CommonService
	Redis   domainService.RedisService
}

func InitServices(repo *Repositories, db *bun.DB, redis *redis.Client) *Services {
	redisService := service.NewRedisRepository(redis, repo.Redis)
	commonService := service.NewCommonService(repo.Device, repo.MemberDevice)

	return &Services{
		Account: account.NewAccountService(db, repo.Password, repo.Account, repo.Member),
		Device:  device.NewDeviceService(repo.Device, repo.Trip, redisService, commonService),
		Member:  member.NewMemberService(repo.Device, repo.MemberDevice),
		Common:  commonService,
		Redis:   redisService,
	}
}
