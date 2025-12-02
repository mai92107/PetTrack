package account

import (
	"PetTrack/core/global"
	handler "PetTrack/infra/02-handler/handler"
	"PetTrack/infra/02-handler/request"
	"net/http"
)

type request01 struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
}

func Login(ctx request.RequestContext) {
	var req request01
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	data, err := handler.AccountService.Login(ctx.GetContext(), ctx.GetClientIP(), req.UserAccount, req.Password)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(data)
}
