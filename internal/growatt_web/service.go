package growatt_web

import (
	"context"
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
	Location                      *time.Location
}

type DurationCalculator interface {
	Initial() (time.Duration, time.Duration)
	Next(lastTimestamp time.Time, retryDuration time.Duration) (time.Duration, time.Time, time.Duration)
}

type GrowattService struct {
	opts             Options
	client           *Client
	health           models.ServiceHealth
	devices          []models.NoahDevicePayload
	endpoint         endpoint.Endpoint
	cancel           context.CancelFunc
	parameterTrigger map[string]chan struct{}
}

func NewGrowattService(options Options) *GrowattService {
	return &GrowattService{
		opts:             options,
		client:           newClient(options.ServerUrl, options.Username, options.Password),
		health:           models.NewServiceHealth(),
		parameterTrigger: make(map[string]chan struct{}),
	}
}

func (g *GrowattService) TriggerParameterPolling(device models.NoahDevicePayload) {
	select {
	case g.parameterTrigger[device.Serial] <- struct{}{}:
	default: // Trigger already pending
	}
}

func (g *GrowattService) Login() error {
	slog.Info("logging in to growatt (web)...")
	if err := g.client.Login(); err != nil {
		return err
	}
	return nil
}

func (g *GrowattService) StartPolling(dc DurationCalculator) {
	g.devices = g.enumerateDevices()
	g.endpoint.SetDevices(g.devices)

	var ctx context.Context
	ctx, g.cancel = context.WithCancel(context.Background())
	for _, device := range g.devices {
		g.parameterTrigger[device.Serial] = make(chan struct{}, 1)
		go g.poll(ctx, device, dc)
	}
}

func (g *GrowattService) StopPolling() {
	g.cancel()
}

func (g *GrowattService) SetEndpoint(e endpoint.Endpoint) {
	g.endpoint = e
	if e != nil {
		e.SetParameterApplier(g)
	}
}

func (g *GrowattService) publishHealth(device models.NoahDevicePayload, err error) {
	if err != nil {
		g.health.UpdateError(device.Serial, err)
	} else {
		g.health.UpdateSuccess(device.Serial)
		g.TriggerParameterPolling(device)
	}
	g.endpoint.PublishHealth(device, &g.health)
}

func (g *GrowattService) SetOutputPowerW(device models.NoahDevicePayload, mode models.WorkMode, power float64) error {
	slog.Info("trying to set default system output power (web)", slog.String("device", device.Serial), slog.Any("mode", mode), slog.Float64("power", power))

	modeAsInt := models.IntFromWorkMode(mode)
	if modeAsInt < 0 {
		slog.Error("unable to set default system output power (web). Invalid mode", slog.String("device", device.Serial), slog.String("mode", string(mode)))
		err := fmt.Errorf("invalid work mode: %s", mode)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set default system output power (web)", slog.String("device", device.Serial), slog.Int("mode", modeAsInt), slog.Float64("power", power))
	err := g.client.SetSystemOutputPower(device.Serial, modeAsInt, power)
	if err != nil {
		slog.Error("unable to set default system output power (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetChargingLimits(device models.NoahDevicePayload, chargingLimit float64, dischargeLimit float64) error {
	slog.Info("trying to set charging limits (web)", slog.String("device", device.Serial), slog.Float64("chargingLimit", chargingLimit), slog.Float64("dischargeLimit", dischargeLimit))

	slog.Debug("set charging limit low (web)", slog.String("device", device.Serial), slog.Float64("dischargeLimit", dischargeLimit))
	err := g.client.SetChargingSocLowLimit(device.Serial, dischargeLimit)
	if err != nil {
		slog.Error("unable to set charging limit low (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {

		slog.Debug("set charging limit high (web)", slog.String("device", device.Serial), slog.Float64("chargingLimit", chargingLimit))
		err = g.client.SetChargingSocHighLimit(device.Serial, chargingLimit)
		if err != nil {
			slog.Error("unable to set charging limit high (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		}
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetAllowGridCharging(device models.NoahDevicePayload, allow models.OnOff) error {
	slog.Info("trying to set allow grid charging (web)", slog.String("device", device.Serial), slog.Any("allow", allow))

	on := misc.OnOffToInt(allow)
	if on < 0 {
		slog.Error("unable to set allow grid charging (web). Invalid allow value", slog.String("device", device.Serial), slog.String("allow", string(allow)))
		err := fmt.Errorf("invalid ON/OFF value: %s", allow)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set allow grid charging (web)", slog.String("device", device.Serial), slog.Int("allow", on))
	err := g.client.SetAllowGridCharging(device.Serial, on)
	if err != nil {
		slog.Error("unable to set allow grid charging (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetGridConnectionControl(device models.NoahDevicePayload, offlineEnable models.OnOff) error {
	slog.Info("trying to set grid connection control (web)", slog.String("device", device.Serial), slog.Any("offlineEnable", offlineEnable))

	on := misc.OnOffToInt(offlineEnable)
	if on < 0 {
		slog.Error("unable to set grid connection control (web). Invalid offlineEnable value", slog.String("device", device.Serial), slog.String("offlineEnable", string(offlineEnable)))
		err := fmt.Errorf("invalid ON/OFF value: %s", offlineEnable)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set grid connection control (web)", slog.String("device", device.Serial), slog.Int("offlineEnable", on))
	err := g.client.SetGridConnectionControl(device.Serial, on)
	if err != nil {
		slog.Error("unable to set grid connection (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetAcCouplePowerControl(device models.NoahDevicePayload, _1000WEnable models.OnOff) error {
	slog.Info("trying to set ac couple power control (web)", slog.String("device", device.Serial), slog.Any("_1000WEnable", _1000WEnable))

	on := misc.OnOffToInt(_1000WEnable)
	if on < 0 {
		slog.Error("unable to set ac couple power control (web). Invalid _1000WEnable value", slog.String("device", device.Serial), slog.String("_1000WEnable", string(_1000WEnable)))
		err := fmt.Errorf("invalid ON/OFF value: %s", _1000WEnable)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set ac couple power control (web)", slog.String("device", device.Serial), slog.Int("_1000WEnable", on))
	err := g.client.SetACCouplePowerControl(device.Serial, on)
	if err != nil {
		slog.Error("unable to set ac couple power control (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetLightLoadEnable(device models.NoahDevicePayload, enable models.OnOff) error {
	slog.Info("trying to set light load enable (web)", slog.String("device", device.Serial), slog.Any("enable", enable))

	on := misc.OnOffToInt(enable)
	if on < 0 {
		slog.Error("unable to set light load enable (web). Invalid enable value", slog.String("device", device.Serial), slog.String("enable", string(enable)))
		err := fmt.Errorf("invalid ON/OFF value: %s", enable)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set light load enable (web)", slog.String("device", device.Serial), slog.Int("enable", misc.OnOffToInt(enable)))
	err := g.client.SetLightLoadEnable(device.Serial, misc.OnOffToInt(enable))
	if err != nil {
		slog.Error("unable to set light load enable (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetNeverPowerOff(device models.NoahDevicePayload, enable models.OnOff) error {
	slog.Info("trying to set never power off (web)", slog.String("device", device.Serial), slog.Any("enable", enable))

	on := misc.OnOffToInt(enable)
	if on < 0 {
		slog.Error("unable to set never power off (web). Invalid enable value", slog.String("device", device.Serial), slog.String("enable", string(enable)))
		err := fmt.Errorf("invalid ON/OFF value: %s", enable)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set never power off (web)", slog.String("device", device.Serial), slog.Int("enable", on))
	err := g.client.SetNeverPowerOff(device.Serial, on)
	if err != nil {
		slog.Error("unable to set never power off (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) SetBackflow(device models.NoahDevicePayload, enableLimit models.OnOff, powerSettingPercent float64) error {
	slog.Info("trying to set backflow (web)", slog.String("device", device.Serial), slog.Any("enableLimit", enableLimit), slog.Float64("powerSettingPercent", powerSettingPercent))

	on := misc.OnOffToInt(enableLimit)
	if on < 0 {
		slog.Error("unable to set backflow (web). Invalid enable value", slog.String("device", device.Serial), slog.String("enable", string(enableLimit)))
		err := fmt.Errorf("invalid ON/OFF value: %s", enableLimit)
		g.publishHealth(device, err)
		return err
	}

	slog.Debug("set backflow (web)", slog.String("device", device.Serial), slog.Int("enableLimit", on), slog.Float64("powerSettingPercent", powerSettingPercent))
	err := g.client.SetAntiBackflowSetting(device.Serial, on, powerSettingPercent)
	if err != nil {
		slog.Error("unable to set backflow (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
	}

	g.publishHealth(device, err)
	return err
}

func (g *GrowattService) enumerateDevices() []models.NoahDevicePayload {
	var enumeratedDevices []models.NoahDevicePayload

	plantList, err := g.client.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		return enumeratedDevices
	}

	for _, plant := range plantList {
		if devices, err := g.client.GetNoahList(misc.S2i(plant.PlantId)); err != nil {
			slog.Error("could not get plant devices", slog.String("plantId", plant.PlantId), slog.String("error", err.Error()))
		} else {
			for _, dev := range devices.Datas {

				if history, err := g.client.GetNoahHistory(dev.Sn, "", ""); err != nil {
					slog.Error("could not get device history", slog.String("device", dev.Sn), slog.String("error", err.Error()))
				} else {
					if len(history.Obj.Datas) == 0 {
						slog.Info("could not get device history, data empty", slog.String("device", dev.Sn))
					} else {
						var batCount = history.Obj.Datas[0].BatteryPackageQuantity
						var batteries []models.NoahDeviceBatteryPayload
						for i := 0; i < batCount; i++ {
							batteries = append(batteries, models.NoahDeviceBatteryPayload{
								Alias: fmt.Sprintf("BAT%d", i),
							})
						}
						d := models.NoahDevicePayload{
							PlantId:   misc.S2i(dev.PlantID),
							Serial:    dev.Sn,
							Model:     dev.DeviceModel,
							Version:   dev.Version,
							Alias:     dev.Alias,
							Batteries: batteries,
						}

						enumeratedDevices = append(enumeratedDevices, d)
					}
				}

			}
		}
	}

	return enumeratedDevices
}

func (g *GrowattService) poll(ctx context.Context, device models.NoahDevicePayload, dc DurationCalculator) {
	slog.Info("start polling growatt (web)",
		slog.Int("interval", int(g.opts.PollingInterval/time.Second)),
		slog.Int("battery-details-interval", int(g.opts.BatteryDetailsPollingInterval/time.Second)),
		slog.Int("parameter-interval", int(g.opts.ParameterPollingInterval/time.Second)))

	durationToWait, retryDuration := dc.Initial()

	tickerPolling := time.NewTicker(g.opts.PollingInterval)
	defer tickerPolling.Stop()
	timerBatteryPolling := time.NewTimer(durationToWait)
	defer timerBatteryPolling.Stop()
	tickerParameter := time.NewTicker(g.opts.ParameterPollingInterval)
	defer tickerParameter.Stop()

	g.pollStatus(device)
	g.pollParameterData(device)

	var lastTimestamp time.Time
	for {
		select {
		case <-tickerPolling.C:
			g.pollStatus(device)

		case <-timerBatteryPolling.C:
			lastTimestamp = g.pollBatteryDetails(device, lastTimestamp)
			durationToWait, lastTimestamp, retryDuration = dc.Next(lastTimestamp, retryDuration)
			timerBatteryPolling.Reset(durationToWait)

			slog.Debug("next battery & pv polling in", slog.String("durationToWait", durationToWait.String()), slog.Time("lastTimestamp", lastTimestamp), slog.String("retryDuration", retryDuration.String()))

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
			slog.Info("stop polling growatt (web)")
			return
		}
	}
}

func (g *GrowattService) pollStatus(device models.NoahDevicePayload) {
	if status, err := g.client.GetNoahStatus(device.PlantId, device.Serial); err != nil {
		slog.Error("could not get device data", slog.String("error", err.Error()), slog.String("device", device.Serial))
		g.health.UpdateError(device.Serial, err)
	} else {
		if totals, err := g.client.GetNoahTotals(device.PlantId, device.Serial); err != nil {
			slog.Error("could not get device totals", slog.String("error", err.Error()), slog.String("device", device.Serial))
			g.health.UpdateError(device.Serial, err)
		} else {
			payload := devicePayload(device, status.Obj, totals.Obj)
			g.endpoint.PublishDeviceStatus(device, payload)
			g.health.UpdateSuccess(device.Serial)
		}
	}
	g.endpoint.PublishHealth(device, &g.health)
}

func (g *GrowattService) pollParameterData(device models.NoahDevicePayload) {
	if details, err := g.client.GetNoahDetails(device.PlantId, device.Serial); err != nil {
		slog.Error("could not get device details data", slog.String("error", err.Error()))
		g.health.UpdateError(device.Serial, err)
	} else {
		if len(details.Datas) != 1 {
			slog.Error("could not get device details data", slog.String("device", device.Serial))
			g.health.UpdateError(device.Serial, fmt.Errorf("no devices available"))
		} else {
			paramPayload := parameterPayload(details.Datas[0])

			g.endpoint.PublishParameterData(device, paramPayload)
			g.health.UpdateSuccess(device.Serial)
		}
	}
	g.endpoint.PublishHealth(device, &g.health)
}

func (g *GrowattService) pollBatteryDetails(device models.NoahDevicePayload, lastTimestamp time.Time) time.Time {

	if history, err := g.client.GetNoahHistory(device.Serial, "", ""); err != nil {
		slog.Error("could not get device history", slog.String("error", err.Error()), slog.String("device", device.Serial))
		g.health.UpdateError(device.Serial, err)
		g.endpoint.PublishHealth(device, &g.health)
	} else {
		if len(history.Obj.Datas) == 0 {
			slog.Info("could not get device history, data empty", slog.String("device", device.Serial))
		} else {
			historyData := history.Obj.Datas[0]

			tm, err := time.ParseInLocation("2006-01-02 15:04:05", historyData.Time, g.opts.Location)
			if err != nil {
				slog.Error("GrowattNoahHistoryData.Time invalid time format", "historyData.Time", historyData.Time, "error", err.Error())
				tm = time.Time{}
			}

			if !lastTimestamp.IsZero() && !tm.IsZero() && !tm.After(lastTimestamp.Add(time.Second)) {
				slog.Info("no new battery details data, last timestamp", slog.Time("lastTimestamp", lastTimestamp), slog.Time("currentDataTimestamp", tm))
				return tm
			}

			var batteries []models.BatteryPayload
			for i := 0; i < len(device.Batteries); i++ {
				batteries = append(batteries, batteryPayload(historyData, tm, i))
			}

			g.endpoint.PublishBatteryDetails(device, batteries)

			var pvs []models.PvPayload
			for i := range 4 {
				pvs = append(pvs, pvPayload(historyData, tm, i))
			}

			g.endpoint.PublishPvDetails(device, pvs)

			g.health.UpdateSuccess(device.Serial)
			g.endpoint.PublishHealth(device, &g.health)
			return tm
		}
	}
	return time.Time{}
}

type defaultDurationCalculator struct {
	defaultDuration time.Duration
}

func NewDefaultDurationCalculator(g *GrowattService) DurationCalculator {
	return &defaultDurationCalculator{
		defaultDuration: g.opts.BatteryDetailsPollingInterval,
	}
}

func (d *defaultDurationCalculator) Initial() (time.Duration, time.Duration) {
	return 0, time.Second * 5
}

func (d *defaultDurationCalculator) Next(lastTimestamp time.Time, retryDuration time.Duration) (time.Duration, time.Time, time.Duration) {
	// Usually a new value in the Growatt historic data is added every 3 minutes. It takes a few
	// seconds until the new value is available after the timestamp of the new value. So we wait
	// until 185 seconds after the last timestamp to be sure that the new value is available.
	// Every now and then the Growatt API waits another few seconds before adding the new value, so
	// we retry once after 5 seconds. If that still gives no new value, we wait the default duration
	// until the next polling.
	//
	// If there is no valid last timestamp at all, we also wait the default duration until the next
	// polling.

	if lastTimestamp.IsZero() {
		// no valid last timestamp, use default duration, reset retry duration
		return d.defaultDuration, lastTimestamp, time.Second * 5
	} else {
		durationToWait := time.Until(lastTimestamp.Add(185 * time.Second))
		if durationToWait < 0 {
			// last timestamp more than 185 seconds ago, use retry duration, set retry duration to default duration for next time
			return retryDuration, lastTimestamp, d.defaultDuration
		} else {
			// should be the normal case, use 185 seconds after last timestamp, reset retry duration
			return durationToWait, lastTimestamp, time.Second * 5
		}
	}
}
