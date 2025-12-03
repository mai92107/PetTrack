package member

import (
	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/util/logafa"	
	handler "PetTrack/infra/02-handler/handler"
	"PetTrack/infra/02-handler/request"
	"net/http"
)

type request02 struct {
	MemberId int64 `json:"memberId"`
}

func MemberDeviceList(ctx request.RequestContext) {
	var req request02
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	deviceIds, err := handler.MemberService.MemberDeviceList(ctx.GetContext(), req.MemberId)
	if err != nil {
		logafa.Error("系統發生錯誤", "error", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(deviceIds)
}
