package router

import (
	"PetTrack/infra/01-router/middleware"
	"PetTrack/infra/02-handler/adapter"
	"PetTrack/infra/02-handler/handler/account"
	"PetTrack/infra/02-handler/handler/device"
	"PetTrack/infra/02-handler/handler/member"
	"PetTrack/infra/02-handler/handler/test"
	"PetTrack/infra/02-handler/handler/trip"
	"PetTrack/infra/02-handler/request"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	const ADMIN = middleware.PermAdmin
	const MEMBER = middleware.PermMember

	r.Use(
		// TODO: 未來需要檢查
		// ip: 看地區
		// header:
		// body:
		middleware.TimeoutMiddleware(3*time.Second),
		middleware.WorkerMiddleware(),
	)

	// 依類別分組
	homeGroup := r.Group("/home")
	{
		homeGroup.POST("/hello", executeHttp(test.SayHello))
	}

	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", executeHttp(account.Login))
		accountGroup.POST("/register", executeHttp(account.Register))
	}

	deviceGroup := r.Group("/device")
	{
		deviceGroup.POST("/create", identityRequired(ADMIN), executeHttp(device.Create))
		deviceGroup.POST("/recording", identityRequired(MEMBER), nil)
		deviceGroup.GET("/onlineDevice", identityRequired(ADMIN), executeHttp(device.OnlineDeviceList))
		deviceGroup.GET("/:deviceId/status", identityRequired(MEMBER), nil)
		deviceGroup.GET("/all", identityRequired(ADMIN), executeHttp(device.DeviceList))
	}

	tripGroup := r.Group("/trip")
	{
		tripGroup.GET("/tripList", identityRequired(MEMBER), executeHttp(trip.DeviceTrips))
		tripGroup.GET("/detail", identityRequired(MEMBER), executeHttp(trip.TripDetail))
	}

	memberGroup := r.Group("/member")
	{
		memberGroup.POST("/addDevice", identityRequired(MEMBER), executeHttp(member.AddDevice))
		memberGroup.GET("/allDevice", identityRequired(MEMBER), executeHttp(member.MemberDeviceList))
	}

	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/status", identityRequired(MEMBER), nil)
	}
}

func executeHttp(handler func(request.RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := adapter.NewHttpContext(c)
		handler(ctx)
	}
}

func identityRequired(identity middleware.Permission) gin.HandlerFunc {
	return middleware.JWTMiddleware(identity)
}
