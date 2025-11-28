package account

import (
	"PetTrack/infra/handler/request"
	"net/http"
)

type request02 struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	NickName  string `json:"nickName"`
}

func Register(ctx request.RequestContext) {
	var req request02
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(http.StatusBadRequest, "Json 格式錯誤")
		return
	}
	ip := ctx.GetClientIP()
	loginInfo, err := accountService.Register(ip, req.Username, req.Password, req.Email, req.LastName, req.FirstName, req.NickName)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "註冊發生錯誤")
		return
	}
	ctx.Success(loginInfo)
}
