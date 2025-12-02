package handler

import (
	domainService "PetTrack/domain/service"
)

var (
	AccountService domainService.AccountService
	DeviceService  domainService.DeviceService
	MemberService  domainService.MemberService
	TripService	domainService.TripService
)

func InitAccountHandler(service domainService.AccountService) {
	AccountService = service
}
func InitDeviceHandler(deviceService domainService.DeviceService, tripService domainService.TripService) {
	DeviceService = deviceService
	TripService = tripService
}
func InitMemberHandler(service domainService.MemberService) {
	MemberService = service
}