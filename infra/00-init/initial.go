package initial

import (
	"PetTrack/infra/00-core/cron"
	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/model"
	"PetTrack/infra/00-core/util/logafa"
	"PetTrack/infra/00-init/wire"
	router "PetTrack/infra/01-router"
	handler "PetTrack/infra/02-handler/handler"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {
	InitLogger()
	InitWorkers()

	// load config
	config, err := ReadFromEnv()
	if err != nil || config == nil {
		logafa.Error("讀取Config失敗", "error", err)
		return
	}

	// init system
	services, _ := wire.InitService(*config)
	InitHandlers(services)
	InitCron(services)

	// start http server
	srv := InitServer(config.Http.Port)
	GracefulShutdown(srv)
}

func InitLogger() {
	logafa.CreateLogFileNow()

	handler := logafa.NewLogafaHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	slog.SetDefault(slog.New(handler))
	logafa.Debug("Logafa 初始化完成")
}

func ReadFromEnv() (*model.Config, error) {

	// 載入 .env（如果檔案不存在就跳過，不會錯）
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := &model.Config{}

	config.Http.Port = viper.GetString("HTTP_PORT")

	config.DeviceConfig.DevicePrefix = viper.GetString("DEVICE_PREFIX")
	config.DeviceConfig.DeviceSequence = viper.GetString("DEVICE_SEQUENCE")

	config.Keys.JwtSecretKey = viper.GetString("JWT_SECRET_KEY")

	config.Machines.MariaDB.User = viper.GetString("DB_USER")
	config.Machines.MariaDB.Password = viper.GetString("DB_PASSWORD")
	config.Machines.MariaDB.Host = viper.GetString("DB_HOST")
	config.Machines.MariaDB.Port = viper.GetString("DB_PORT")
	config.Machines.MariaDB.Name = viper.GetString("DB_NAME")

	config.Machines.MongoDB.User = viper.GetString("MONGO_USER")
	config.Machines.MongoDB.Password = viper.GetString("MONGO_PASSWORD")
	config.Machines.MongoDB.Host = viper.GetString("MONGO_HOST")
	config.Machines.MongoDB.Port = viper.GetString("MONGO_PORT")
	config.Machines.MongoDB.Name = viper.GetString("MONGO_NAME")
	config.Machines.MongoDB.TimeoutRange = viper.GetInt("MONGO_TIMEOUT_RANGE")

	config.Machines.Redis.Password = viper.GetString("REDIS_PASSWORD")
	config.Machines.Redis.Host = viper.GetString("REDIS_HOST")
	config.Machines.Redis.Port = viper.GetString("REDIS_PORT")

	// MQTT
	config.Machines.MqttBroker.HostCloud = viper.GetString("MQTT_HOST_CLOUD")
	config.Machines.MqttBroker.HostLocal = viper.GetString("MQTT_HOST_LOCAL")
	config.Machines.MqttBroker.Port = viper.GetString("MQTT_PORT")
	config.Machines.MqttBroker.User = viper.GetString("MQTT_USER")
	config.Machines.MqttBroker.Password = viper.GetString("MQTT_PASSWORD")
	config.Machines.MqttBroker.ClientID = viper.GetString("MQTT_CLIENT_ID")
	config.Machines.MqttBroker.Topic = strings.Split(viper.GetString("MQTT_TOPIC"), ",")

	return config, nil
}

func InitWorkers() {
	maxPriorWorkers := 20
	maxNormalWorkers := 50
	// 區隔工人 做 故障隔離
	// 高級勞工
	global.PriorWorkerPool = make(chan struct{}, maxPriorWorkers)
	for range maxPriorWorkers {
		global.PriorWorkerPool <- struct{}{}
	}
	logafa.Debug("高級勞工 初始化完成", "數量", maxPriorWorkers)
	// 城市打工人
	global.NormalWorkerPool = make(chan struct{}, maxNormalWorkers)
	for range maxNormalWorkers {
		global.NormalWorkerPool <- struct{}{}
	}
	logafa.Debug("城市打工人 初始化完成", "數量", maxNormalWorkers)
}

func InitHandlers(s *wire.Services) {
	handler.InitAccountHandler(s.Account)
	handler.InitDeviceHandler(s.Device, s.Trip)
	handler.InitMemberHandler(s.Member)
}

func InitCron(s *wire.Services) {
	cron.NewScheduler(
		s.Trip,
	).CronStart()
}

func InitServer(port string) *http.Server {
	r := gin.Default()
	router.RegisterRoutes(r)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           r,
		MaxHeaderBytes:    8 * 1024, // 8K
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       0,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logafa.Error("伺服器啟動失敗", "error", err)
		}
	}()
	return srv

}

func GracefulShutdown(srv *http.Server) {
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
