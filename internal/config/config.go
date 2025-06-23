package config

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	LogLevel                      string
	PollingInterval               time.Duration
	BatteryDetailsPollingInterval time.Duration
	ParameterPollingInterval      time.Duration
	Growatt                       Growatt
	Mqtt                          Mqtt
	HomeAssistant                 HomeAssistant
}

type Growatt struct {
	APIMode      string
	ServerUrlWeb string
	ServerUrlApp string
	Username     string
	Password     string
}

type Mqtt struct {
	Host        string
	Port        int
	ClientId    string
	Username    string
	Password    string
	TopicPrefix string
}

type HomeAssistant struct {
	TopicPrefix string
}

var _config Config
var _once sync.Once

func Get() Config {
	_once.Do(func() {
		_config = Config{
			LogLevel:                      getEnv("LOG_LEVEL", "info"),
			PollingInterval:               time.Duration(s2i(getEnv("POLLING_INTERVAL", "30"))) * time.Second,
			BatteryDetailsPollingInterval: time.Duration(s2i(getEnv("BATTERY_DETAILS_POLLING_INTERVAL", "180"))) * time.Second,
			ParameterPollingInterval:      time.Duration(s2i(getEnv("PARAMETER_POLLING_INTERVAL", "180"))) * time.Second,
			Growatt: Growatt{
				APIMode:      getEnv("GROWATT_API_MODE", "web+app"),
				ServerUrlWeb: getEnv("GROWATT_SERVER_URL_WEB", "https://openapi.growatt.com"),
				ServerUrlApp: getEnv("GROWATT_SERVER_URL_APP", "https://server-api.growatt.com"),
				Username:     getEnv("GROWATT_USERNAME", ""),
				Password:     getEnv("GROWATT_PASSWORD", ""),
			},
			Mqtt: Mqtt{
				Host:        getEnv("MQTT_HOST", ""),
				Port:        s2i(getEnv("MQTT_PORT", "1883")),
				ClientId:    getEnv("MQTT_CLIENT_ID", "nexa-mqtt"),
				Username:    getEnv("MQTT_USERNAME", ""),
				Password:    getEnv("MQTT_PASSWORD", ""),
				TopicPrefix: getEnv("MQTT_TOPIC_PREFIX", "noah2mqtt"),
			},
			HomeAssistant: HomeAssistant{
				TopicPrefix: getEnv("HOMEASSISTANT_TOPIC_PREFIX", "homeassistant"),
			},
		}
	})
	return _config
}

func Validate() error {
	if len(Get().Mqtt.Host) == 0 {
		return errors.New("MQTT_HOST is required")
	}
	if len(Get().Growatt.Username) == 0 {
		return errors.New("GROWATT_USERNAME is required")
	}
	if len(Get().Growatt.Password) == 0 {
		return errors.New("GROWATT_PASSWORD is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func s2i(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
