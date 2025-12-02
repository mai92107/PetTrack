package adapter

import (
	"PetTrack/infra/02-handler/request"
	"PetTrack/infra/02-handler/response"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPContext struct {
	ctx         *gin.Context
	requestTime time.Time
}

func NewHttpContext(c *gin.Context) request.RequestContext {
	return &HTTPContext{
		ctx:         c,
		requestTime: time.Now(),
	}
}

// Create new context
func (h *HTTPContext) GetContext() (context.Context) {
	return h.ctx.Request.Context()
}
func (h *HTTPContext) Cancel() {
}

// BindJSON implements request.RequestContext.
func (h *HTTPContext) BindJSON(obj interface{}) error {
	return h.ctx.ShouldBindJSON(obj)
}

// GetClientID implements request.RequestContext.
func (h *HTTPContext) GetClientID() string {
	return ""
}

// GetClientIP implements request.RequestContext.
func (h *HTTPContext) GetClientIP() string {
	return h.ctx.ClientIP()
}

// GetJWT implements request.RequestContext.
func (h *HTTPContext) GetJWT() string {
	return h.ctx.GetHeader("jwt")
}

// GetRequestTime implements request.RequestContext.
func (h *HTTPContext) GetRequestTime() time.Time {
	return h.requestTime
}

// Success implements request.RequestContext.
func (h *HTTPContext) Success(data interface{}) {
	response.Success(h.ctx, h.requestTime, data)
}

// Error implements request.RequestContext.
func (h *HTTPContext) Error(code int, message string) {
	response.Error(h.ctx, code, h.requestTime, message)
}
