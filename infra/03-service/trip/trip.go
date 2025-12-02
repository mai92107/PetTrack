package trip

import (
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
)

type TripServiceImpl struct {
	tripRepo      domainRepo.TripRepository
	commonService domainService.CommonService
	redisUtil     domainRepo.RedisRepository
}

func NewTripService(
	tripRepo domainRepo.TripRepository,
	commonService domainService.CommonService,
	redisUtil domainRepo.RedisRepository,
) domainService.TripService {
	return &TripServiceImpl{tripRepo: tripRepo, commonService: commonService, redisUtil: redisUtil}
}
