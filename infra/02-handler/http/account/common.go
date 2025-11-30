package account

import "PetTrack/infra/03-service/account"

var accountService account.AccountServiceImpl

func InitAccountHandler(service account.AccountServiceImpl) {
	accountService = service
}