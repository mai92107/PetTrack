package account

import (
	domain "PetTrack/domain/service"

	"gorm.io/gorm"
)

type AccountServiceImpl struct {
	accountRepo     domain.AccountRepository
	memberRepo      domain.MemberRepository
	passwordHistory domain.PasswordRepository
	db              *gorm.DB
}

func NewAccountService(db *gorm.DB, passwordHistory domain.PasswordRepository, accountRepo domain.AccountRepository, memberRepo domain.MemberRepository) domain.AccountService {
	return &AccountServiceImpl{
		db:              db,
		passwordHistory: passwordHistory,
		accountRepo:     accountRepo,
		memberRepo:      memberRepo,
	}
}
