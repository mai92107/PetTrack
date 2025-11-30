package device

import (
	"PetTrack/core/global"
	"PetTrack/infra/02-handler/request"
	"net/http"
)

func DeviceList(ctx request.RequestContext) {

	deviceIds, err := deviceService.DeviceList()
	if err != nil {
		// logafa.Error("系統發生錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(deviceIds)
}
