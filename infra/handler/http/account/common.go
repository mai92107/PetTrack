package account

import "PetTrack/infra/service/account"

var accountService account.AccountServiceImpl

func InitAccountHandler(service account.AccountServiceImpl) {
	accountService = service
}