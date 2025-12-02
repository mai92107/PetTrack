package initial

import handler "PetTrack/infra/02-handler/handler"

func InitHandlers(s *Services) {
	handler.InitAccountHandler(s.Account)
	handler.InitDeviceHandler(s.Device,s.Trip)
	handler.InitMemberHandler(s.Member)
}
