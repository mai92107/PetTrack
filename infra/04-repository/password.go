package repo

import (
	cryptoUtil "PetTrack/infra/00-core/util/crypto"
	domain "PetTrack/domain/repo"
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

func (r *passwordRepoImpl) CreateHistory(tx *gorm.DB, ctx context.Context, accountUuid uuid.UUID, password string) error {
	hashedPassword, _ := cryptoUtil.BcryptHash(password)
	pastPassword := &domain.PasswordHistory{
		AccountUuid: accountUuid,
		Password:    hashedPassword,
	}
	err := tx.WithContext(ctx).Create(&pastPassword).Error
	if err != nil {
		// logafa.Error("建立使用者歷史密碼失敗, error: %+v", err)
		return fmt.Errorf("建立使用者歷史密碼失敗")
	}
	return nil
}
