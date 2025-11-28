package repo

import (
	domain "PetTrack/domain/repo"
	"fmt"

	"gorm.io/gorm"
)

func (r *memberRepoImpl) CreateMember(tx *gorm.DB, lastName, firstName, nickName, email string) (int64, error) {
	member := &domain.Member{
		LastName:  lastName,
		FirstName: firstName,
		NickName:  nickName,
		Email:     email,
	}
	err := tx.Create(&member).Error
	if err != nil {
		// logafa.Error("建立使用者資料失敗, error: %+v",err)
		return 0, fmt.Errorf("建立使用者失敗")
	}
	return member.Id, nil
}
