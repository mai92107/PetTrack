package account

import (
	domainRepo "PetTrack/domain/repo"
	domainService "PetTrack/domain/service"

	"gorm.io/gorm"
)

type AccountServiceImpl struct {
	accountRepo     domainRepo.AccountRepository
	memberRepo      domainRepo.MemberRepository
	passwordHistory domainRepo.PasswordRepository
	db              *gorm.DB
}

func NewAccountService(
	db *gorm.DB, 
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
