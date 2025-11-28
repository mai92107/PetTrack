package router

import (
	"PetTrack/infra/handler/adapter"
	"PetTrack/infra/handler/http/account"
	"PetTrack/infra/handler/request"
	"PetTrack/infra/router/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	const ADMIN = "ADMIN"
	const MEMBER = "MEMBER"

	r.Use(middleware.WorkerMiddleware())

	// 註冊路由
	// TODO: 未來需要檢查ip body header 在路徑後加上middleware檢查
	// 依類別分組
	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/say_hello")
	}

	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", GoTo(account.Login))
		accountGroup.POST("/register", GoTo(account.Register))
	}

	trackGroup := r.Group("/device")
	{
		trackGroup.POST("/create", middleware.JWTValidator(ADMIN), nil)
		trackGroup.POST("/recording", middleware.JWTValidator(MEMBER), nil)
		trackGroup.GET("/onlineDevice", middleware.JWTValidator(ADMIN), nil)
		trackGroup.GET("/:deviceId/status", middleware.JWTValidator(MEMBER), nil)
		trackGroup.GET("/all", middleware.JWTValidator(ADMIN), nil)
	}

	memberGroup := r.Group("/member")
	{
		memberGroup.POST("/addDevice", middleware.JWTValidator(MEMBER), nil)
		memberGroup.GET("/allDevice", middleware.JWTValidator(MEMBER), nil)
	}

	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/status", middleware.JWTValidator(MEMBER), nil)
	}
}

func GoTo(handler func(request.RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := adapter.NewHttpContext(c)
		handler(ctx)
	}
}
