package device

import (
	jwtUtil "PetTrack/core/util/jwt"
	"PetTrack/infra/handler/request"
	"net/http"
)

type request01 struct {
	DeviceType string `json:"deviceType"`
}

func Create(ctx request.RequestContext) {
	var req request01
	if err := ctx.BindJSON(&req); err != nil {
		// logafa.Error("Json 格式錯誤, error: %+v", err)
		ctx.Error(http.StatusBadRequest, "Json 格式錯誤")
		return
	}
	userData, err := jwtUtil.GetUserDataFromJwt(ctx.GetJWT())
	if err != nil || userData.Identity != "ADMIN" {
		// logafa.Error("身份認證錯誤, error: %+v", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	deviceId, err := service.Create(req.DeviceType, userData.MemberId)
	if err != nil {
		// logafa.Error("裝置新增失敗，請稍後嘗試, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, "裝置新增失敗，請稍後嘗試")
		return
	}
	ctx.Success(deviceId)
}
