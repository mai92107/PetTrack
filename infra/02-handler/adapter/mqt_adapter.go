package adapter

import (
	"PetTrack/infra/02-handler/request"
	"PetTrack/infra/02-handler/response"
	"context"
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
)

type MQTTContext struct {
	client      mqtt.Client
	payload     string
	jwt         string
	clientID    string
	ip          string
	requestTime time.Time

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMQTTContext(client mqtt.Client, payload, jwt, clientID, ip string, requestTime time.Time) request.RequestContext {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	return &MQTTContext{
		client:      client,
		payload:     payload,
		jwt:         jwt,
		clientID:    clientID,
		ip:          ip,
		requestTime: requestTime,

		ctx:    ctx,
		cancel: cancel,
	}
}

// Create new context
func (m *MQTTContext) GetContext() context.Context {
	return m.ctx
}
func (m *MQTTContext) Cancel() {
	if m.cancel != nil {
		m.cancel()
	}
}

// BindJSON implements request.RequestContext.
func (m *MQTTContext) BindJSON(obj interface{}) error {
	if m.payload == "" || m.payload == "{}" {
		return errors.New("empty payload")
	}
	return jsoniter.UnmarshalFromString(m.payload, obj)
}

// GetClientID implements request.RequestContext.
func (m *MQTTContext) GetClientID() string {
	return m.clientID
}

// GetClientIP implements request.RequestContext.
func (m *MQTTContext) GetClientIP() string {
	return m.ip
}

// GetJWT implements request.RequestContext.
func (m *MQTTContext) GetJWT() string {
	return m.jwt
}

// GetRequestTime implements request.RequestContext.
func (m *MQTTContext) GetRequestTime() time.Time {
	return m.requestTime
}

// Success implements request.RequestContext.
func (m *MQTTContext) Success(data interface{}) {
	response.SuccessMqtt(m.client, m.getResponseTopic(), m.requestTime, data)
}

// Error implements request.RequestContext.
func (m *MQTTContext) Error(code int, message string) {
	errTopic := "errReq/" + m.clientID
	response.ErrorMqtt(m.client, errTopic, code, m.requestTime, message)
}

func (m *MQTTContext) getResponseTopic() string {
	var temp struct {
		SubscribeTo string `json:"subscribeTo"`
	}
	jsoniter.UnmarshalFromString(m.payload, &temp)
	return temp.SubscribeTo
}
