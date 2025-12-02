package router

import (
	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/util/logafa"
	"PetTrack/infra/01-router/middleware"
	"PetTrack/infra/02-handler/adapter"
	"PetTrack/infra/02-handler/handler/account"
	"PetTrack/infra/02-handler/handler/device"
	"PetTrack/infra/02-handler/handler/member"
	"PetTrack/infra/02-handler/handler/test"
	"PetTrack/infra/02-handler/handler/trip"
	"PetTrack/infra/02-handler/request"
	"PetTrack/infra/02-handler/response"
	"fmt"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT route
var mqttRoutes = map[string]Route{
	// Guest (無需 JWT)
	"account_login":    {Handler: executeMqtt(account.Login), Permission: PermGuest},
	"account_register": {Handler: executeMqtt(account.Register), Permission: PermGuest},
	"say_hello":        {Handler: executeMqtt(test.SayHello), Permission: PermGuest},
	"system_status":    {Handler: nil, Permission: PermGuest},

	// Admin
	"device_create": {Handler: executeMqtt(device.Create), Permission: PermAdmin},
	"device_online": {Handler: nil, Permission: PermAdmin},
	"device_all":    {Handler: executeMqtt(device.DeviceList), Permission: PermAdmin},

	// Member
	"device_recording": {Handler: nil, Permission: PermMember},
	"member_addDevice": {Handler: executeMqtt(member.AddDevice), Permission: PermMember},
	"device_status":    {Handler: nil, Permission: PermMember},
	"member_devices":   {Handler: executeMqtt(member.MemberDeviceList), Permission: PermMember},
	"trip_list":        {Handler: executeMqtt(trip.DeviceTrips), Permission: PermMember},
	"trip":             {Handler: executeMqtt(trip.TripDetail), Permission: PermMember},
}

type MqttHandler func(cleint mqtt.Client, payload, jwt, clientId, ip string, requestTime time.Time)

type Permission int

const (
	PermGuest Permission = iota
	PermMember
	PermAdmin
)

type Route struct {
	Handler    MqttHandler
	Permission Permission
}

// topic sample : req/action/clientId/jwt/ip
func RouteFunction(client mqtt.Client, action, payload, clientId, jwt, ip string, requestTime time.Time) {
	// 查找路由
	routeInfo := mqttRoutes[action]
	// 權限檢查
	if _, err := middleware.ValidateJWT(jwt, middleware.Permission(routeInfo.Permission)); err != nil {
		logafa.Warn("JWT 驗證失敗", "error", err, "action", action, "user", clientId)
		sendBackErrMsg(client, clientId, "JWT 驗證失敗, err: %+v", err)
		return
	}

	routeInfo.Handler(client, payload, jwt, clientId, ip, requestTime)
}

func OnMessageReceived(client mqtt.Client, msg mqtt.Message) {
	requestTime := global.GetNow()
	payload := string(msg.Payload())
	topic := msg.Topic()

	logafa.Debug("收到 MQTT 訊息！Topic: %s | Payload: %s", topic, payload)

	action, clientId, jwt, ip := extractInfoFromTopic(topic)
	if action == "" || ip == "" {
		logafa.Warn("無法解析 action 或 ip: %s", topic)
		sendBackErrMsg(client, clientId, "無法解析 action 或 ip: %s", topic)
		return
	}

	// 使用 worker pool 執行
	<-global.NormalWorkerPool
	go func() {
		defer func() {
			global.NormalWorkerPool <- struct{}{}
			if r := recover(); r != nil {
				logafa.Error("MQTT handler panic: %v\n on %s", r, topic)
			}
		}()
		RouteFunction(client, action, payload, clientId, jwt, ip, requestTime)
	}()
}

func extractInfoFromTopic(topic string) (action, clientId, jwt, ip string) {
	parts := strings.Split(topic, "/")
	return parts[1], parts[2], parts[3], parts[4]
}

func sendBackErrMsg(client mqtt.Client, clientId, reason string, args ...interface{}) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId
	fullReason := fmt.Sprintf(reason, args...)
	response.ErrorMqtt(client, errTopic, http.StatusBadRequest, requestTime, fullReason)
}

func executeMqtt(handler func(request.RequestContext)) MqttHandler {
	return func(client mqtt.Client, payload, jwt, clientID, ip string, requestTime time.Time) {
		ctx := adapter.NewMQTTContext(client, payload, jwt, clientID, ip, requestTime)
		handler(ctx)
	}
}
