package initial

import handler "PetTrack/infra/02-handler/http"

func InitHandlers(s *Services) {
	handler.InitAccountHandler(s.Account)
	handler.InitDeviceHandler(s.Device)
	handler.InitMemberHandler(s.Member)
}
