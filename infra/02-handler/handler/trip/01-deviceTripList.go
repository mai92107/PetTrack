package trip

import (
	"PetTrack/core/global"
	"PetTrack/core/model"
	jwtUtil "PetTrack/core/util/jwt"
	handler "PetTrack/infra/02-handler/handler"
	"PetTrack/infra/02-handler/request"
	"PetTrack/infra/02-handler/response"
	"net/http"
)

type request01 struct {
	DeviceId string `json:"deviceId"`
	request.PageInfo
}

func DeviceTrips(ctx request.RequestContext) {
	var req request01
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
	datas, total, totalPages, err := handler.TripService.GetDeviceTrips(*userInfo, ctx.GetContext(), req.DeviceId, model.NewPageable(&req.Page, &req.Size, req.Direction, req.OrderBy))
	if err != nil {
		// logafa.Error("系統發生錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	pageInfo := response.GetPageResponse(req.PageInfo, total, totalPages)
	info := map[string]interface{}{
		"pageInfo": pageInfo,
		"trips":    datas,
	}

	ctx.Success(info)
}
