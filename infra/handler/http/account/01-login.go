package account

import (
	"PetTrack/infra/handler/request"
	"net/http"
)

type request01 struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
}

func Login(ctx request.RequestContext) {
	var req request01
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(http.StatusBadRequest, "請求格式錯誤")
		return
	}
	data, err := accountService.Login(ctx.GetClientIP(), req.UserAccount, req.Password)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Success(data)
}
