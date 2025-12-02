package middleware

import (
	"PetTrack/infra/00-core/util/logafa"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WorkerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// // 搶 worker（會阻塞直到有空位）
		// <-global.NormalWorkerPool
		// defer func() {
		// 	global.NormalWorkerPool <- struct{}{}
		// 	// logafa.Debug("工作完畢")
		// 	if r := recover(); r != nil {
		// 		// logafa.Error("Handler panic recovered: %v", r)
		// 		c.JSON(500, gin.H{"error": "internal server error"})
		// 	}
		// }()
		c.Next()
	}
}

// 為每個 request 套用 timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 建立可取消 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 將新的 context 放進 request
		c.Request = c.Request.WithContext(ctx)

		// 建立 channel 用來接 handler 結果
		finished := make(chan struct{})
		panicChan := make(chan interface{})

		// 使用 goroutine 來執行真正的 handler
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()        // 執行 handler / middleware
			close(finished) // handler 正常完成
		}()

		// 用 select 監控三種結果
		select {
		case <-finished:
			// handler 正常完成
			return

		case p := <-panicChan:
			// handler 發生 panic
			logafa.Error("發生Panic", "panic", p)

		case <-ctx.Done():
			// timeout，取消 handler 實作
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error":   "request timeout",
				"timeout": timeout.String(),
			})

			// 終止 request chain
			c.Abort()
			return
		}
	}
}
