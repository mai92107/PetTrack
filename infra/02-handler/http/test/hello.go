package test

import (
	"PetTrack/core/util/logafa"
	"PetTrack/infra/02-handler/request"
	"time"
)

func SayHello(ctx request.RequestContext) {
	logafa.Info("say hello", "user", "unknown")
	time.Sleep(10 * time.Second)
	ctx.Success("Helllllllo")
}
