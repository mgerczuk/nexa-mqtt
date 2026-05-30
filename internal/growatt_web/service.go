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
}

type DurationCalculator interface {
	Initial() (time.Duration, time.Duration)
	Next(lastTimestamp time.Time, retryDuration time.Duration) (time.Duration, time.Time, time.Duration)
}

type GrowattService struct {
	opts     Options
	client   *Client
	devices  []models.NoahDevicePayload
	endpoint endpoint.Endpoint
	cancel   context.CancelFunc
}

func NewGrowattService(options Options) *GrowattService {
	return &GrowattService{
		opts:   options,
		client: newClient(options.ServerUrl, options.Username, options.Password),
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

func (g *GrowattService) SetOutputPowerW(device models.NoahDevicePayload, mode models.WorkMode, power float64) error {
	slog.Info("trying to set default system output power (web)", slog.String("device", device.Serial), slog.Any("mode", mode), slog.Float64("power", power))

	modeAsInt := models.IntFromWorkMode(mode)
	if modeAsInt < 0 {
		slog.Error("unable to set default system output power (web). Invalid mode", slog.String("device", device.Serial), slog.String("mode", string(mode)))
		return fmt.Errorf("invalid work mode %s", mode)
	}

	slog.Debug("set default system output power (web)", slog.String("device", device.Serial), slog.Int("mode", modeAsInt), slog.Float64("power", power))
	if err := g.client.SetSystemOutputPower(device.Serial, modeAsInt, power); err != nil {
		slog.Error("unable to set default system output power (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
}

func (g *GrowattService) SetChargingLimits(device models.NoahDevicePayload, chargingLimit float64, dischargeLimit float64) error {
	slog.Info("trying to set charging limits (web)", slog.String("device", device.Serial), slog.Float64("chargingLimit", chargingLimit), slog.Float64("dischargeLimit", dischargeLimit))

	slog.Debug("set charging limit low (web)", slog.String("device", device.Serial), slog.Float64("dischargeLimit", dischargeLimit))
	if err := g.client.SetChargingSocLowLimit(device.Serial, dischargeLimit); err != nil {
		slog.Error("unable to set charging limit low (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	slog.Debug("set charging limit high (web)", slog.String("device", device.Serial), slog.Float64("chargingLimit", chargingLimit))
	if err := g.client.SetChargingSocHighLimit(device.Serial, chargingLimit); err != nil {
		slog.Error("unable to set charging limit high (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
}

func (g *GrowattService) SetAllowGridCharging(device models.NoahDevicePayload, allow models.OnOff) error {
	slog.Info("trying to set allow grid charging (web)", slog.String("device", device.Serial), slog.Any("allow", allow))

	slog.Debug("set allow grid charging (web)", slog.String("device", device.Serial), slog.Int("allow", misc.OnOffToInt(allow)))
	if err := g.client.SetAllowGridCharging(device.Serial, misc.OnOffToInt(allow)); err != nil {
		slog.Error("unable to set allow grid charging (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}
	return nil
}

func (g *GrowattService) SetGridConnectionControl(device models.NoahDevicePayload, offlineEnable models.OnOff) error {
	slog.Info("trying to set grid connection control (web)", slog.String("device", device.Serial), slog.Any("offlineEnable", offlineEnable))

	slog.Debug("set grid connection control (web)", slog.String("device", device.Serial), slog.Int("offlineEnable", misc.OnOffToInt(offlineEnable)))
	if err := g.client.SetGridConnectionControl(device.Serial, misc.OnOffToInt(offlineEnable)); err != nil {
		slog.Error("unable to set grid connection (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
}

func (g *GrowattService) SetAcCouplePowerControl(device models.NoahDevicePayload, _1000WEnable models.OnOff) error {
	slog.Info("trying to set ac couple power control (web)", slog.String("device", device.Serial), slog.Any("_1000WEnable", _1000WEnable))

	slog.Debug("set ac couple power control (web)", slog.String("device", device.Serial), slog.Int("_1000WEnable", misc.OnOffToInt(_1000WEnable)))
	if err := g.client.SetACCouplePowerControl(device.Serial, misc.OnOffToInt(_1000WEnable)); err != nil {
		slog.Error("unable to set ac couple power control (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
}

func (g *GrowattService) SetLightLoadEnable(device models.NoahDevicePayload, enable models.OnOff) error {
	slog.Info("trying to set light load enable (web)", slog.String("device", device.Serial), slog.Any("enable", enable))

	slog.Debug("set light load enable (web)", slog.String("device", device.Serial), slog.Int("enable", misc.OnOffToInt(enable)))
	if err := g.client.SetLightLoadEnable(device.Serial, misc.OnOffToInt(enable)); err != nil {
		slog.Error("unable to set light load enable (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
}

func (g *GrowattService) SetNeverPowerOff(device models.NoahDevicePayload, enable models.OnOff) error {
	slog.Info("trying to set never power off (web)", slog.String("device", device.Serial), slog.Any("enable", enable))

	slog.Debug("set never power off (web)", slog.String("device", device.Serial), slog.Int("enable", misc.OnOffToInt(enable)))
	if err := g.client.SetNeverPowerOff(device.Serial, misc.OnOffToInt(enable)); err != nil {
		slog.Error("unable to set never power off (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
}

func (g *GrowattService) SetBackflow(device models.NoahDevicePayload, enableLimit models.OnOff, powerSettingPercent float64) error {
	slog.Info("trying to set backflow (web)", slog.String("device", device.Serial), slog.Any("enableLimit", enableLimit), slog.Float64("powerSettingPercent", powerSettingPercent))

	slog.Debug("set backflow (web)", slog.String("device", device.Serial), slog.Int("enableLimit", misc.OnOffToInt(enableLimit)), slog.Float64("powerSettingPercent", powerSettingPercent))
	if err := g.client.SetAntiBackflowSetting(device.Serial, misc.OnOffToInt(enableLimit), powerSettingPercent); err != nil {
		slog.Error("unable to set backflow (web)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return err
	}

	return nil
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

		case <-ctx.Done():
			slog.Info("stop polling growatt (web)")
			return
		}
	}
}

func (g *GrowattService) pollStatus(device models.NoahDevicePayload) {
	if status, err := g.client.GetNoahStatus(device.PlantId, device.Serial); err != nil {
		slog.Error("could not get device data", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		if totals, err := g.client.GetNoahTotals(device.PlantId, device.Serial); err != nil {
			slog.Error("could not get device totals", slog.String("error", err.Error()), slog.String("device", device.Serial))
		} else {
			payload := devicePayload(device, status.Obj, totals.Obj)
			g.endpoint.PublishDeviceStatus(device, payload)
		}
	}
}

func (g *GrowattService) pollParameterData(device models.NoahDevicePayload) {
	if details, err := g.client.GetNoahDetails(device.PlantId, device.Serial); err != nil {
		slog.Error("could not get device details data", slog.String("error", err.Error()))
	} else {
		if len(details.Datas) != 1 {
			slog.Error("could not get device details data", slog.String("device", device.Serial))
		} else {
			paramPayload := parameterPayload(details.Datas[0])

			g.endpoint.PublishParameterData(device, paramPayload)
		}
	}
}

func (g *GrowattService) pollBatteryDetails(device models.NoahDevicePayload, lastTimestamp time.Time) time.Time {

	if history, err := g.client.GetNoahHistory(device.Serial, "", ""); err != nil {
		slog.Error("could not get device history", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		if len(history.Obj.Datas) == 0 {
			slog.Info("could not get device history, data empty", slog.String("device", device.Serial))
		} else {
			historyData := history.Obj.Datas[0]

			tm, err := time.ParseInLocation("2006-01-02 15:04:05", historyData.Time, time.Local)
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
