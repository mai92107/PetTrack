package account

import (
	bun "PetTrack/core/model/bunMachines"
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"
)

type AccountServiceImpl struct {
	accountRepo     domainRepo.AccountRepository
	memberRepo      domainRepo.MemberRepository
	passwordHistory domainRepo.PasswordRepository
	db              *bun.DB
}

func NewAccountService(
	db *bun.DB,
	passwordHistory domainRepo.PasswordRepository,
	accountRepo domainRepo.AccountRepository,
	memberRepo domainRepo.MemberRepository,
) domainService.AccountService {
	return &AccountServiceImpl{
		db:              db,
		passwordHistory: passwordHistory,
		accountRepo:     accountRepo,
		memberRepo:      memberRepo,
	}
}
