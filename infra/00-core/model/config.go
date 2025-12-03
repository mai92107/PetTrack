package model

type Config struct {
	Http         Http         `mapstructure:"http"`
	DeviceConfig DeviceConfig `mapstructure:"device_config"`
	Keys         Keys         `mapstructure:"keys"`
	Machines     Machine      `mapstructure:"machines"`
	DBConfig     DbSetting    `mapstructure:"db_config"`
}

type Http struct {
	Port string `mapstructure:"port"`
}

type DeviceConfig struct {
	DevicePrefix   string `mapstructure:"device_prefix"`
	DeviceSequence string `mapstructure:"device_sequence"`
}

type Keys struct {
	JwtSecretKey string `mapstructure:"jwt_secret_key"`
}

type Machine struct {
	MariaDB    DbSetting    `mapstructure:"maria_db"`
	MongoDB    MongoSetting `mapstructure:"mongo_db"`
	Redis      RedisSetting `mapstructure:"redis_db"`
	MqttBroker MqttConfig   `mapstructure:"mqtt_broker"`
}

type DbSetting struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
}

type MongoSetting struct {
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	TimeoutRange int    `mapstructure:"timeout_range"`
}

type RedisSetting struct {
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type MqttConfig struct {
	HostCloud string   `mapstructure:"host_cloud"`
	HostLocal string   `mapstructure:"host_local"`
	Port      string   `mapstructure:"port"`
	User      string   `mapstructure:"user"`
	Password  string   `mapstructure:"password"`
	Topic     []string `mapstructure:"topic"`
	ClientID  string   `mapstructure:"client_id"`
}
