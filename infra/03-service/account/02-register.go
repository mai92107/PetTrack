package account

import (
	"PetTrack/core/global"
	"fmt"
)

func (service *AccountServiceImpl) Register(ip, username, password, email, lastName, firstName, nickName string) (map[string]interface{}, error) {
	// 參數驗證

	err := validateRegister(email, nickName, username, password)
	if err != nil {
		return nil, err
	}

	tx := service.db.Write.Begin()

	memberId, err := service.memberRepo.CreateMember(tx, lastName, firstName, nickName, email)
	if err != nil {
		return nil, err
	}

	accountUuid, err := service.accountRepo.Create(tx, memberId, username, password, email)
	if err != nil {
		return nil, err
	}

	err = service.passwordHistory.CreateHistory(tx, accountUuid, password)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	return service.Login(ip, username, password)
}

func validateRegister(email, nickName, username, password string) error {
	if username == "" {
		return fmt.Errorf("使用者帳號不可為空")
	}
	if password == "" {
		return fmt.Errorf("使用者密碼不可為空")
	}
	if email == "" {
		return fmt.Errorf("電子信箱不可為空")
	}
	if nickName == "" {
		return fmt.Errorf("使用者名稱不可為空")
	}
	return nil
}
