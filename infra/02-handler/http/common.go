package handler

import (
	domainService "PetTrack/domain/service"
)

var (
	AccountService domainService.AccountService
	DeviceService  domainService.DeviceService
	MemberService  domainService.MemberService
)

func InitAccountHandler(service domainService.AccountService) {
	AccountService = service
}
func InitDeviceHandler(service domainService.DeviceService) {
	DeviceService = service
}
func InitMemberHandler(service domainService.MemberService) {
	MemberService = service
}
