package account

import (
	cryptoUtil "PetTrack/core/util/crypto"
	jwtUtil "PetTrack/core/util/jwt"
	"fmt"
	"time"
)

func (service *AccountServiceImpl) Login(ip, accountName, password string) (map[string]interface{}, error) {
	// 參數驗證
	if err := validateLogin(accountName, password); err != nil {
		return nil, err
	}

	var err error
	data := map[string]interface{}{}

	// 驗證帳號
	userAccount, err := service.accountRepo.FindByAccountName(accountName)
	if err != nil {
		return data, err
	}

	// 驗證密碼
	if !cryptoUtil.BcryptCompare(userAccount.Password, password) {
		return data, fmt.Errorf("密碼錯誤")
	}

	err = service.accountRepo.UpdateLoginTime(userAccount.Uuid)
	if err != nil {
		return data, err
	}

	now := time.Now().UTC()
	expireTime := 24 * time.Hour
	token, err := jwtUtil.GenerateJwt(accountName, userAccount.Identity, userAccount.MemberId, ip, now, expireTime)
	if err != nil {
		return data, err
	}
	data = map[string]interface{}{
		"token":     token,
		"identity":  userAccount.Identity,
		"loginTime": now,
	}
	return data, nil
}

func validateLogin(accountName, password string) error {
	if accountName == "" || password == "" {
		return fmt.Errorf("帳號或密碼不可為空")
	}
	return nil
}
