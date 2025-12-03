package device

import (
	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/util/logafa"
	handler "PetTrack/infra/02-handler/handler"
	"PetTrack/infra/02-handler/request"
	"net/http"
)

func DeviceList(ctx request.RequestContext) {

	deviceIds, err := handler.DeviceService.DeviceList(ctx.GetContext())
	if err != nil {
		logafa.Error("系統發生錯誤", "error", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(deviceIds)
}
