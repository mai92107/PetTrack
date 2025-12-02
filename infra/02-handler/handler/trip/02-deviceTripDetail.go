package trip

import (
	"PetTrack/infra/00-core/global"
	jwtUtil "PetTrack/infra/00-core/util/jwt"
	handler "PetTrack/infra/02-handler/handler"
	"PetTrack/infra/02-handler/request"
	"net/http"
)

type request02 struct {
	DeviceId string `json:"deviceId"`
	TripUuid string `json:"tripUuid"`
}

func TripDetail(ctx request.RequestContext) {
	var req request02
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	jwt := ctx.GetJWT()
	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		// logafa.Error("身份認證錯誤, error: %+v", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	info, err := handler.TripService.GetTripDetail(ctx.GetContext(), *userInfo, req.DeviceId, req.TripUuid)
	if err != nil {
		// logafa.Error("系統發生錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(info)
}
