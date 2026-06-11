package growatt_app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"nexa-mqtt/internal/endpoint"
	"nexa-mqtt/internal/misc"
	"nexa-mqtt/pkg/models"
	"time"
)

type Options struct {
	ServerUrl                     string
	Username                      string
	Password                      string
	PollingInterval               time.Duration
	BatteryDetailsPollingInterval time.Duration
	ParameterPollingInterval      time.Duration
}
type GrowattAppService struct {
	opts             Options
	client           *Client
	health           models.ServiceHealth
	devices          []models.NoahDevicePayload
	endpoint         endpoint.Endpoint
	loggedIn         bool
	cancel           context.CancelFunc
	parameterTrigger map[string]chan struct{}
	query            endpoint.ParameterQuery
}

func NewGrowattAppService(options Options) *GrowattAppService {
	service := GrowattAppService{
		opts:             options,
		client:           newClient(options.ServerUrl, options.Username, options.Password),
		health:           models.NewServiceHealth(),
		loggedIn:         false,
		parameterTrigger: make(map[string]chan struct{}),
	}
	service.query = &service
	return &service
}

func (g *GrowattAppService) TriggerParameterPolling(device models.NoahDevicePayload) {
	select {
	case g.parameterTrigger[device.Serial] <- struct{}{}:
	default: // Trigger already pending
	}
}

func (g *GrowattAppService) Login() error {
	slog.Info("logging in to growatt (app)...")

	if err := g.client.Login(); err != nil {
		return err
	}
	g.loggedIn = true
	return nil
}

func (g *GrowattAppService) StartPolling() {
	g.enumerateDevices()
	var ctx context.Context
	ctx, g.cancel = context.WithCancel(context.Background())
	for _, device := range g.devices {
		g.parameterTrigger[device.Serial] = make(chan struct{}, 1)
		go g.poll(ctx, device)
	}
}

func (g *GrowattAppService) StopPolling() {
	g.cancel()
}

func (g *GrowattAppService) fetchDevices() []models.NoahDevicePayload {
	slog.Info("fetching plant list")
	list, err := g.client.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	var devices []models.NoahDevicePayload

	for _, plant := range list.PlantList {
		slog.Info("fetch plant details", slog.Int("plantId", plant.ID))
		if info, err := g.client.GetNoahPlantInfo(fmt.Sprintf("%d", plant.ID)); err != nil {
			slog.Error("could not get plant info", slog.Int("plantId", plant.ID), slog.String("error", err.Error()))
		} else {
			if len(info.Obj.DeviceSn) > 0 {
				devices = append(devices, models.NoahDevicePayload{
					PlantId:   plant.ID,
					Serial:    info.Obj.DeviceSn,
					Batteries: nil,
				})
				slog.Info("found device sn", slog.String("deviceSn", info.Obj.DeviceSn), slog.Int("plantId", plant.ID))
			}
		}
	}

	if len(devices) == 0 {
		slog.Error("no nexa devices found")
		misc.Panic(errors.New("no nexa devices found"))
	}

	return devices
}

func (g *GrowattAppService) enumerateDevices() {
	devices := g.fetchDevices()

	for i, device := range devices {
		if data, err := g.client.GetNexaInfoBySn(device.Serial); err != nil {
			slog.Error("could not get nexa status", slog.String("error", err.Error()), slog.String("serialNumber", device.Serial))
		} else {
			batCount := len(data.Obj.Noah.BatSns)
			var batteries []models.NoahDeviceBatteryPayload
			for i := 0; i < batCount; i++ {
				batteries = append(batteries, models.NoahDeviceBatteryPayload{
					Alias: fmt.Sprintf("BAT%d", i),
				})
			}

			devices[i].Model = data.Obj.Noah.Model
			devices[i].Version = data.Obj.Noah.Version
			devices[i].Alias = data.Obj.Noah.Alias
			devices[i].Batteries = batteries
		}
	}

	g.devices = devices

	g.endpoint.SetDevices(devices)
}

func (g *GrowattAppService) SetEndpoint(e endpoint.Endpoint) {
	g.endpoint = e
}

func (g *GrowattAppService) SetParameterQuery(q endpoint.ParameterQuery) {
	g.query = q
}

func (g *GrowattAppService) ensureParameterLogin() error {
	if !g.loggedIn {
		if err := g.Login(); err != nil {
			slog.Error("could not login to growatt account (app)", slog.String("error", err.Error()))
			return err
		}
	}
	return nil
}

func (g *GrowattAppService) publishHealth(device models.NoahDevicePayload, err error) {
	if err != nil {
		g.health.UpdateError(device.Serial, err)
	} else {
		g.health.UpdateSuccess(device.Serial)
		g.query.TriggerParameterPolling(device)
	}
	g.endpoint.PublishHealth(device, &g.health)
}

func (g *GrowattAppService) SetOutputPowerW(device models.NoahDevicePayload, mode models.WorkMode, power float64) error {
	slog.Info("trying to set default system output power (app)", slog.String("device", device.Serial), slog.String("mode", string(mode)), slog.Float64("power", power))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set default system output power (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	modeAsInt := models.IntFromWorkMode(mode)
	if modeAsInt < 0 {
		slog.Error("unable to set default system output power (app). Invalid mode", slog.String("device", device.Serial), slog.String("mode", string(mode)))
		err := fmt.Errorf("invalid work mode: %s", mode)
		g.publishHealth(device, err)
		return err
	}

	slog.Info("set default system output power (app)", slog.String("device", device.Serial), slog.Int("mode", modeAsInt), slog.Float64("power", power))
	err := g.client.SetSystemOutputPower(device.Serial, modeAsInt, power)
	if err != nil {
		slog.Error("unable to set default system output power (app)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetChargingLimits(device models.NoahDevicePayload, chargingLimit float64, dischargeLimit float64) error {
	slog.Info("trying to set charging limits (app)", slog.String("device", device.Serial), slog.Float64("chargingLimit", chargingLimit), slog.Float64("dischargeLimit", dischargeLimit))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set charging limits (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	slog.Info("set charging limit (app)", slog.String("device", device.Serial), slog.Float64("chargingLimit", chargingLimit), slog.Float64("dischargeLimit", dischargeLimit))
	err := g.client.SetChargingSoc(device.Serial, chargingLimit, dischargeLimit)
	if err != nil {
		slog.Error("unable to set charging limits (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetAllowGridCharging(device models.NoahDevicePayload, allow models.OnOff) error {
	slog.Info("trying to set allow charging  (app)", slog.String("device", device.Serial), slog.String("allow", string(allow)))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set allow charging (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	err := g.client.SetAllowGridCharging(device.Serial, misc.OnOffToInt(allow))
	if err != nil {
		slog.Error("unable to set allow charging (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetGridConnectionControl(device models.NoahDevicePayload, offlineEnable models.OnOff) error {
	slog.Info("trying to set grid connection  (app)", slog.String("device", device.Serial), slog.String("offlineEnable", string(offlineEnable)))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set grid connection (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	err := g.client.SetGridConnectionControl(device.Serial, misc.OnOffToInt(offlineEnable))
	if err != nil {
		slog.Error("unable to set grid connection (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetAcCouplePowerControl(device models.NoahDevicePayload, _1000WEnable models.OnOff) error {
	slog.Info("trying to set ac couple power control  (app)", slog.String("device", device.Serial), slog.String("1000WEnable", string(_1000WEnable)))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set ac couple power control (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	err := g.client.SetACCouplePowerControl(device.Serial, misc.OnOffToInt(_1000WEnable))
	if err != nil {
		slog.Error("unable to set ac couple power control (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetLightLoadEnable(device models.NoahDevicePayload, enable models.OnOff) error {
	slog.Info("trying to set light load enable (app)", slog.String("device", device.Serial), slog.String("enable", string(enable)))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set light load enable (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	err := g.client.SetLightLoadEnable(device.Serial, misc.OnOffToInt(enable))
	if err != nil {
		slog.Error("unable to set light load enable (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetNeverPowerOff(device models.NoahDevicePayload, enable models.OnOff) error {
	slog.Info("trying to set never power off (app)", slog.String("device", device.Serial), slog.String("enable", string(enable)))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set set never power off (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	err := g.client.SetNeverPowerOff(device.Serial, misc.OnOffToInt(enable))
	if err != nil {
		slog.Error("unable to set set never power off (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) SetBackflow(device models.NoahDevicePayload, enableLimit models.OnOff, powerSettingPercent float64) error {
	slog.Info("trying to set backflow (app)", slog.String("device", device.Serial), slog.String("enableLimit", string(enableLimit)), slog.Float64("powerSettingPercent", powerSettingPercent))
	if err := g.ensureParameterLogin(); err != nil {
		slog.Error("unable to set backflow (app)", slog.String("device", device.Serial))
		g.publishHealth(device, err)
		return err
	}

	err := g.client.SetBackflow(device.Serial, misc.OnOffToInt(enableLimit), powerSettingPercent)
	if err != nil {
		slog.Error("unable to set backflow (app)", slog.String("error", err.Error()))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattAppService) poll(ctx context.Context, device models.NoahDevicePayload) {
	slog.Info("start polling growatt (app)",
		slog.Int("interval", int(g.opts.PollingInterval/time.Second)),
		slog.Int("battery-details-interval", int(g.opts.BatteryDetailsPollingInterval/time.Second)),
		slog.Int("parameter-interval", int(g.opts.ParameterPollingInterval/time.Second)))

	tickerPolling := time.NewTicker(g.opts.PollingInterval)
	defer tickerPolling.Stop()
	tickerBatteryDetails := time.NewTicker(g.opts.BatteryDetailsPollingInterval)
	defer tickerBatteryDetails.Stop()
	tickerParameter := time.NewTicker(g.opts.ParameterPollingInterval)
	defer tickerParameter.Stop()

	g.pollStatus(device)
	g.pollBatteryDetails(device)
	g.pollParameterData(device)

	for {
		select {
		case <-tickerPolling.C:
			g.pollStatus(device)

		case <-tickerBatteryDetails.C:
			g.pollBatteryDetails(device)

		case <-tickerParameter.C:
			g.pollParameterData(device)

		case <-g.parameterTrigger[device.Serial]:
			g.pollParameterData(device)

			tickerParameter.Stop()
			select {
			case <-tickerParameter.C:
			default:
			}
			tickerParameter.Reset(g.opts.ParameterPollingInterval)

		case <-ctx.Done():
			slog.Info("stop polling growatt (app)")
			return
		}
	}
}
