package test

import (
	"PetTrack/infra/02-handler/request"
)

func SayHello(ctx request.RequestContext) {
	// logafa.Debug("say hello", "user", "unknown")
	// time.Sleep(10 * time.Second)
	ctx.Success("Helllllllo")
}
