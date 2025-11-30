package member

import (
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
)

type MemberServiceImpl struct {
	deviceRepo       domainRepo.DeviceRepository
	memberDeviceRepo domainRepo.MemberDeviceRepository
}

func NewMemberService(
	deviceRepo domainRepo.DeviceRepository,
	memberDeviceRepo domainRepo.MemberDeviceRepository,
) domainService.MemberService {
	return &MemberServiceImpl{deviceRepo: deviceRepo, memberDeviceRepo: memberDeviceRepo}
}
