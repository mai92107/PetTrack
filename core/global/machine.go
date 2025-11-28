package global

import (
	"PetTrack/core/model"
	"sync"
	"sync/atomic"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// var (
// 	ConfigSetting jsonModal.Config
// 	Repository    *Repo
// )

var (
	ActiveDevices     = make(map[string]model.DeviceStatus) // 儲存所有裝置和 狀態
	ActiveDevicesLock sync.Mutex                            // 互斥鎖確保併發安全
	GlobalBroker      mqtt.Client                           // 全域 MQTT 客戶端
	IsConnected       atomic.Bool                           // 確認目前連線狀態
)

var (
	PriorWorkerPool  chan struct{}
	NormalWorkerPool chan struct{}
)

// type Repo struct {
// 	DB    *DataBase
// 	Cache *Cache
// }
// type DataBase struct {
// 	MariaDb *SqlDB
// 	MongoDb *NoSqlDB
// }
// type SqlDB struct {
// 	Reading *gorm.DB
// 	Writing *gorm.DB
// }
// type NoSqlDB struct {
// 	Reading *mongo.Database
// 	Writing *mongo.Database
// }
// type Cache struct {
// 	Reading *redis.Client
// 	Writing *redis.Client
// 	CTX     context.Context
// }
