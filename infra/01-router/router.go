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

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	const ADMIN = middleware.PermAdmin
	const MEMBER = middleware.PermMember

	r.Use(middleware.WorkerMiddleware())

	// 註冊路由
	// TODO: 未來需要檢查ip body header 在路徑後加上IdentityRequired檢查
	// 依類別分組
	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/say_hello", executeHttp(test.SayHello))
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
		tripGroup.GET("/trip", identityRequired(MEMBER), executeHttp(trip.TripDetail))
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
