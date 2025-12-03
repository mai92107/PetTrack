// infra/01-router/mqtt.go
package mqtt

import (
	"net/http"
	"strings"
	"time"

	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/util/logafa"
	router "PetTrack/infra/01-router"
	"PetTrack/infra/02-handler/adapter"
	"PetTrack/infra/02-handler/handler/account"
	"PetTrack/infra/02-handler/handler/device"
	"PetTrack/infra/02-handler/handler/member"
	"PetTrack/infra/02-handler/handler/test"
	"PetTrack/infra/02-handler/handler/trip"
	"PetTrack/infra/02-handler/request"
	"PetTrack/infra/02-handler/response"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTMessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		topic := msg.Topic()

		logafa.Debug("收到 MQTT 訊息", "topic", topic, "payload", payload)

		action, clientId, jwt, ip := extractInfoFromTopic(topic)
		if action == "" || clientId == "" || ip == "" {
			logafa.Warn("MQTT topic 解析失敗", "topic", topic)
			sendBackErrMsg(client, clientIdOrEmpty(clientId), "topic 格式錯誤")
			return
		}

		// 使用 worker pool，避免阻塞 MQTT 客戶端
		select {
		case <-global.NormalWorkerPool:
			go func() {
				defer func() { global.NormalWorkerPool <- struct{}{} }()
				handleMQTTMessage(client, action, payload, clientId, jwt, ip)
			}()
		default:
			logafa.Warn("Worker pool 已滿，丟棄 MQTT 訊息", "action", action, "clientId", clientId)
			sendBackErrMsg(client, clientId, "系統繁忙，請稍後再試")
		}
	}
}

func handleMQTTMessage(client mqtt.Client, action, payload, clientId, jwt, ip string) {
	now := global.GetNow()

	// 路由分發 + 權限檢查
	routeInfo, exists := mqttRoutes[action]
	if !exists || routeInfo.Handler == nil {
		logafa.Warn("查無此 MQTT 路徑", "action", action)
		sendBackErrMsg(client, clientId, "此功能暫未開放")
		return
	}

	// JWT 權限驗證
	if _, err := router.ValidateJWT(jwt, router.Permission(routeInfo.Permission)); err != nil {
		logafa.Warn("MQTT JWT 驗證失敗", "action", action, "clientId", clientId, "error", err)
		sendBackErrMsg(client, clientId, "登入驗證失敗")
		return
	}

	// 執行對應的 handler
	routeInfo.Handler(client, payload, jwt, clientId, ip, now)
}

var mqttRoutes = map[string]Route{
	"home_hello": {Handler: executeMqtt(test.SayHello), Permission: PermGuest},

	// account
	"account_login":    {Handler: executeMqtt(account.Login), Permission: PermGuest},
	"account_register": {Handler: executeMqtt(account.Register), Permission: PermGuest},

	// device
	"device_create": {Handler: executeMqtt(device.Create), Permission: PermAdmin},
	"device_online": {Handler: executeMqtt(device.OnlineDeviceList), Permission: PermAdmin},
	"device_all":    {Handler: executeMqtt(device.DeviceList), Permission: PermAdmin},

	// trip
	"trips": {Handler: executeMqtt(trip.DeviceTrips), Permission: PermMember},
	"trip":  {Handler: executeMqtt(trip.TripDetail), Permission: PermMember},

	// member
	"member_addDevice": {Handler: executeMqtt(member.AddDevice), Permission: PermMember},
	"member_devices":   {Handler: executeMqtt(member.MemberDeviceList), Permission: PermMember},

	// system
	"system_status": {Handler: nil, Permission: PermGuest},
}


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

type MqttHandler func(mqtt.Client, string, string, string, string, time.Time)

// 解析 topic: req/action/clientId/jwt/ip
func extractInfoFromTopic(topic string) (action, clientId, jwt, ip string) {
	parts := strings.Split(topic, "/")
	if len(parts) < 5 {
		return "", "", "", ""
	}
	return parts[1], parts[2], parts[3], parts[4]
}

func clientIdOrEmpty(id string) string {
	if id == "" {
		return "unknown"
	}
	return id
}

func sendBackErrMsg(client mqtt.Client, clientId, reason string, args ...interface{}) {
	errTopic := "errReq/" + clientId
	response.ErrorMqtt(client, errTopic, http.StatusBadRequest, global.GetNow(), fmt.Sprintf(reason, args...))
}

// executeMqtt：把普通 handler 轉成 MqttHandler，並自動處理 cancel + panic
func executeMqtt(handler func(request.RequestContext)) MqttHandler {
	return func(client mqtt.Client, payload, jwt, clientID, ip string, requestTime time.Time) {
		ctx := adapter.NewMQTTContext(client, payload, jwt, clientID, ip, requestTime)

		defer ctx.Cancel() // 自動 cancel，永不洩漏！
		defer func() {
			if r := recover(); r != nil {
				logafa.Error("MQTT handler panic", "error", r, "action", payload[:min(len(payload), 100)])
				ctx.Error(500, "系統錯誤")
			}
		}()

		handler(ctx)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
