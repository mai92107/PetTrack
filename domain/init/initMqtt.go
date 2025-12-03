package initMethod

import (
	"PetTrack/infra/00-core/global"
	"PetTrack/infra/00-core/model"
	"PetTrack/infra/00-core/util/logafa"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	subscriptionMutex sync.Mutex
	subscribedTopics  = make(map[string]bool)
)

func InitMqtt(cfg model.Config) *mqtt.Client {
	return initMqtt(cfg.Machines.MqttBroker.HostCloud, cfg.Machines.MqttBroker.Port, cfg.Machines.MqttBroker.User, cfg.Machines.MqttBroker.Password, cfg.Machines.MqttBroker.Topic)
}

// InitMosquitto 初始化 MQTT 連線
func initMqtt(host, port, username, password string, topics []string) *mqtt.Client {

	clientId := "bunbun"

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", host, port)).
		SetClientID(fmt.Sprintf("%s-%d", clientId, time.Now().UTC().UnixNano())).
		SetUsername(username).
		SetPassword(password).
		SetKeepAlive(120 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetDefaultPublishHandler(onMessageReceived).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetMaxReconnectInterval(60 * time.Second).
		SetCleanSession(true).
		SetOnConnectHandler(func(c mqtt.Client) {
			logafa.Debug("✅ 已連接到 Mosquitto 伺服器")
			go subscribeVagueTopic(c, topics)
		}).
		SetConnectionLostHandler(onConnectionLost).
		SetReconnectingHandler(func(c mqtt.Client, opts *mqtt.ClientOptions) {
			logafa.Info("🔄 正在重新連接到 Mosquitto 伺服器...")
		})

	client := mqtt.NewClient(opts)

	// 初次連線
	logafa.Debug("🔌 正在連接到 MQTT Broker", "host", host, "port", port)

	// 初次連線（非阻塞）
	if token := client.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
		logafa.Error("初始連線失敗", "error", token.Error())
		return nil
	}
	// 更新連線狀態
	global.IsConnected.Swap(true)
	logafa.Debug("✅ MQTT 客戶端初始化成功")
	return &client
}

func subscribeVagueTopic(client mqtt.Client, vagueTopic []string) {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	for _, topic := range vagueTopic {
		if subscribedTopics[topic] {
			continue
		}
		token := client.Subscribe(topic, 1, nil)
		go func(t string, tok mqtt.Token) {
			if tok.Wait() && tok.Error() != nil {
				logafa.Error("訂閱失敗", "topic", t, "error", tok.Error())
			} else {
				subscriptionMutex.Lock()
				subscribedTopics[t] = true
				subscriptionMutex.Unlock()
				logafa.Debug("系統開始追蹤裝置主題", "topic", t)
			}
		}(topic, token)
	}
}

// onConnectionLost 當連線中斷時觸發
func onConnectionLost(client mqtt.Client, err error) {
	logafa.Error("🚫 Mosquitto 伺服器連線斷開", "error", err)
	subscriptionMutex.Lock()
	subscribedTopics = make(map[string]bool)
	// 更新連線狀態
	global.IsConnected.Swap(false)
	subscriptionMutex.Unlock()
}

func onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	topic := msg.Topic()

	logafa.Debug("收到 MQTT 訊息", "topic", topic, "payload", payload)
	// ProcessMsg(payload, topic, client)
}
