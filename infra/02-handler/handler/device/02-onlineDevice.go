package device

import (
	"PetTrack/infra/00-core/global"
	jwtUtil "PetTrack/infra/00-core/util/jwt"
	handler "PetTrack/infra/02-handler/handler"
	"PetTrack/infra/02-handler/request"
	"net/http"
)

func OnlineDeviceList(ctx request.RequestContext) {
	jwt := ctx.GetJWT()
	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil || userInfo.Identity != "ADMIN" {
		// logafa.Error("身份認證錯誤, error: %+v", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	deviceList, err := handler.DeviceService.OnlineDeviceList(ctx.GetContext())
	if err != nil {
		// logafa.Error("系統發生錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(deviceList)
}
