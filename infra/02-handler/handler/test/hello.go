package test

import (
	"PetTrack/infra/02-handler/request"
)

type request01 struct {
	Msg string `json:"msg"`
}

func SayHello(ctx request.RequestContext) {
	var req request01
	ctx.BindJSON(&req)
	// time.Sleep(10 * time.Second)
	ctx.Success(req.Msg)
}
