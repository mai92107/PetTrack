package test

import "PetTrack/infra/02-handler/request"

func SayHello(ctx request.RequestContext) {
	ctx.Success("hello")
}
