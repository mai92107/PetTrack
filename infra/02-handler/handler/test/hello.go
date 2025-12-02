package test

import (
	"PetTrack/core/util/logafa"
	"PetTrack/infra/02-handler/request"
)

func SayHello(ctx request.RequestContext) {
	logafa.Info("say hello", "user", "unknown")
	ctx.Success("hello")
}
