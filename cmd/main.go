package main

import (
	"PetTrack/infra/00-core/cron"
	"PetTrack/infra/00-core/util/logafa"
	initial "PetTrack/infra/00-init"
	router "PetTrack/infra/01-router"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initial.Init()

	srv := initServer()
	gracefulShutdown(srv)
}

func initServer() *http.Server {
	r := gin.Default()
	router.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8000"),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logafa.Error("伺服器啟動失敗", "error", err)
		}
	}()
	return srv

}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // 等待訊號
	logafa.Info("收到終止訊號，開始優雅關閉...")

	cron.CheckIsCronJobsFinished()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logafa.Error("伺服器優雅關閉失敗", "error", err)
	} else {
		logafa.Info("伺服器成功關閉")
	}
}
