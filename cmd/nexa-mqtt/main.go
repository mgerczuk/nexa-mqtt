package main

import (
	"fmt"
	"log/slog"
	"nexa-mqtt/internal/config"
	"nexa-mqtt/internal/endpoint_mqtt"
	"nexa-mqtt/internal/growatt_app"
	"nexa-mqtt/internal/growatt_web"
	"nexa-mqtt/internal/homeassistant"
	"nexa-mqtt/internal/logging"
	"nexa-mqtt/internal/misc"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	version = "local"
	commit  = "none"
)

func main() {
	cfg := config.Get()
	logging.Init(cfg.LogLevel)
	if err := config.Validate(); err != nil {
		slog.Error("couldn't validate config", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	fmt.Fprintf(os.Stdout, "--- nexa-mqtt started (version: %s, commit: %s)\n", version, commit)

	if currentUser, err := user.Current(); err == nil {
		fmt.Fprintf(os.Stdout, "    running as user: %s (uid: %s)\n", currentUser.Username, currentUser.Uid)
	}

	app := NewApp(cfg)
	connectMqtt(cfg.Mqtt, app)

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	slog.Info("Caught signal", slog.Any("signal", sig))
}

type App struct {
	mode              string
	cfg               config.Config
	growattWebService *growatt_web.GrowattService
	growattAppService *growatt_app.GrowattAppService
}

func (a *App) onMqttDisconnect() {
	if a.growattWebService != nil {
		a.growattWebService.StopPolling()
		a.growattWebService.SetEndpoint(nil)
	}
	if a.growattAppService != nil {
		a.growattAppService.StopPolling()
		a.growattAppService.SetEndpoint(nil)
	}
}

func (a *App) onMqttConnect(client mqtt.Client) {
	haService := homeassistant.NewService(homeassistant.Options{
		MqttClient:     client,
		TopicPrefix:    a.cfg.HomeAssistant.TopicPrefix,
		SwitchAsSelect: a.cfg.HomeAssistant.SwitchAsSelect,
		Version:        version,
	})

	mqttEndpoint := endpoint_mqtt.NewEndpoint(endpoint_mqtt.Options{
		MqttClient:  client,
		TopicPrefix: a.cfg.Mqtt.TopicPrefix,
		HaClient:    haService,
	})

	client.Publish(fmt.Sprintf("%s/availability", a.cfg.Mqtt.TopicPrefix), 1, true, "online")

	switch a.mode {
	case "app":
		a.growattAppService.SetEndpoint(mqttEndpoint)
		a.growattAppService.StartPolling()
		mqttEndpoint.SetParameterApplier(a.growattAppService)

	case "web":
		a.growattWebService.SetEndpoint(mqttEndpoint)
		a.growattWebService.StartPolling(growatt_web.NewDefaultDurationCalculator(a.growattWebService))
		mqttEndpoint.SetParameterApplier(a.growattWebService)

	case "web+app":
		a.growattWebService.SetEndpoint(mqttEndpoint)
		a.growattWebService.StartPolling(growatt_web.NewDefaultDurationCalculator(a.growattWebService))
		a.growattAppService.SetEndpoint(mqttEndpoint)
		a.growattAppService.SetParameterQuery(a.growattWebService)
		mqttEndpoint.SetParameterApplier(a.growattAppService)
	}
}

func NewApp(cfg config.Config) *App {

	mode := strings.ToLower(strings.TrimSpace(cfg.Growatt.APIMode))
	switch mode {
	case "app":
		slog.Info("setting mode", slog.String("mode", mode))
		growattApp := growatt_app.NewGrowattAppService(growatt_app.Options{
			ServerUrl:                     cfg.Growatt.ServerUrlApp,
			Username:                      cfg.Growatt.Username,
			Password:                      cfg.Growatt.Password,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
		})

		if err := growattApp.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}
		return &App{
			mode:              mode,
			cfg:               cfg,
			growattAppService: growattApp,
		}

	case "web":
		slog.Info("setting mode", slog.String("mode", mode))
		growattService := growatt_web.NewGrowattService(growatt_web.Options{
			ServerUrl:                     cfg.Growatt.ServerUrlWeb,
			Username:                      cfg.Growatt.Username,
			Password:                      cfg.Growatt.Password,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
			Location:                      cfg.Growatt.Location,
		})

		if err := growattService.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}

		return &App{
			mode:              mode,
			cfg:               cfg,
			growattWebService: growattService,
		}

	case "web+app":
		slog.Info("setting mode", slog.String("mode", mode))
		growattService := growatt_web.NewGrowattService(growatt_web.Options{
			ServerUrl:                     cfg.Growatt.ServerUrlWeb,
			Username:                      cfg.Growatt.Username,
			Password:                      cfg.Growatt.Password,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
			Location:                      cfg.Growatt.Location,
		})

		if err := growattService.Login(); err != nil {
			slog.Error("could not login to growatt account", slog.String("error", err.Error()))
			misc.Panic(err)
		}

		growattApp := growatt_app.NewGrowattAppService(growatt_app.Options{
			ServerUrl:                     cfg.Growatt.ServerUrlApp,
			Username:                      cfg.Growatt.Username,
			Password:                      cfg.Growatt.Password,
			PollingInterval:               cfg.PollingInterval,
			BatteryDetailsPollingInterval: cfg.BatteryDetailsPollingInterval,
			ParameterPollingInterval:      cfg.ParameterPollingInterval,
		})

		return &App{
			mode:              mode,
			cfg:               cfg,
			growattWebService: growattService,
			growattAppService: growattApp,
		}

	default:
		misc.Panic(fmt.Errorf("invalid growatt api type: %s", cfg.Growatt.APIMode))
		return nil
	}
}

func connectMqtt(mqttCfg config.Mqtt, app *App) {
	var brokerUrl string
	if mqttCfg.BrokerURL != "" {
		brokerUrl = mqttCfg.BrokerURL
	} else {
		brokerUrl = fmt.Sprintf("tcp://%s:%d", mqttCfg.Host, mqttCfg.Port)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(brokerUrl).
		SetClientID(mqttCfg.ClientId).
		SetUsername(mqttCfg.Username).
		SetPassword(mqttCfg.Password).
		SetWill(fmt.Sprintf("%s/availability", mqttCfg.TopicPrefix), "offline", 1, true)

	opts.OnConnect = func(client mqtt.Client) {
		slog.Info("connected to mqtt broker")
		app.onMqttConnect(client)
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		slog.Warn("lost connection to mqtt broker", slog.String("error", err.Error()))
		app.onMqttDisconnect()
	}

	c := mqtt.NewClient(opts)
	slog.Info("connecting to mqtt broker", slog.String("brokerUrl", brokerUrl), slog.String("host", mqttCfg.Host), slog.Int("port", mqttCfg.Port), slog.String("clientId", mqttCfg.ClientId), slog.String("username", mqttCfg.Username))
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("could not connect to mqtt broker", slog.String("error", token.Error().Error()))
		misc.Panic(token.Error())
	}
}
