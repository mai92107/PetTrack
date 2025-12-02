package repo

import (
	domain "PetTrack/domain/repo"
	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/model"
	"PetTrack/infra/00-core/util/logafa"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *accountRepoImpl) FindByAccountName(ctx context.Context, accountName string) (*domain.Account, error) {
	account := &domain.Account{}
	var err error

	if strings.Contains(accountName, "@") {
		err = r.read.WithContext(ctx).First(&account, "email = ?", accountName).Error
	} else {
		err = r.read.WithContext(ctx).First(&account, "username = ?", accountName).Error
	}

	if err != nil {
		logafa.Error("查詢帳戶發生錯誤", "error", err)
		return nil, fmt.Errorf("查詢帳戶發生錯誤")
	}
	return account, nil
}

func (r *accountRepoImpl) UpdateLoginTime(ctx context.Context, uuid uuid.UUID) error {
	err := r.write.WithContext(ctx).Model(&domain.Account{}).Where("uuid = ?", uuid).Update("last_login_time", gorm.Expr("NOW()")).Error
	if err != nil {
		logafa.Error("更新最後登入時間發生錯誤", "error", err)
		return fmt.Errorf("更新最後登入時間發生錯誤")
	}
	return nil
}

func (r *accountRepoImpl) Create(tx *gorm.DB, ctx context.Context, memberId int64, username, password, email string) (uuid.UUID, error) {
	account := &domain.Account{
		Uuid:          uuid.New(),
		MemberId:      memberId,
		Username:      username,
		Password:      password,
		Email:         email,
		Identity:      model.MEMBER.ToString(),
		LastLoginTime: global.GetNow(),
	}
	err := r.write.WithContext(ctx).Create(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				return uuid.Nil, fmt.Errorf("使用者帳號 %s 已存在", username)
			}
			if strings.Contains(err.Error(), "email") {
				return uuid.Nil, fmt.Errorf("電子郵件 %s 已存在", email)
			}
		}
		logafa.Error("建立帳戶失敗", "error", err)
		return uuid.Nil, fmt.Errorf("建立帳戶失敗")
	}
	return account.Uuid, nil

}
