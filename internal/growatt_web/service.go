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

func (g *GrowattService) StartPolling() {
	g.devices = g.enumerateDevices()
	g.endpoint.SetDevices(g.devices)

	var ctx context.Context
	ctx, g.cancel = context.WithCancel(context.Background())
	go g.poll(ctx)
}

func (g *GrowattService) StopPolling() {
	g.cancel()
}

func (g *GrowattService) SetEndpoint(e endpoint.Endpoint) {
	g.endpoint = e
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

func (g *GrowattService) poll(ctx context.Context) {
	slog.Info("start polling growatt (web)",
		slog.Int("interval", int(g.opts.PollingInterval/time.Second)),
		slog.Int("battery-details-interval", int(g.opts.BatteryDetailsPollingInterval/time.Second)),
		slog.Int("parameter-interval", int(g.opts.ParameterPollingInterval/time.Second)))

	tickerPolling := time.NewTicker(g.opts.PollingInterval)
	defer tickerPolling.Stop()
	tickerBatteryDetails := time.NewTicker(g.opts.BatteryDetailsPollingInterval)
	defer tickerBatteryDetails.Stop()
	tickerParameter := time.NewTicker(g.opts.ParameterPollingInterval)
	defer tickerParameter.Stop()

	for _, device := range g.devices {
		g.pollStatus(device)
		g.pollBatteryDetails(device)
		g.pollParameterData(device)
	}

	for {
		select {
		case <-tickerPolling.C:
			for _, device := range g.devices {
				g.pollStatus(device)
			}

		case <-tickerBatteryDetails.C:
			for _, device := range g.devices {
				g.pollBatteryDetails(device)
			}

		case <-tickerParameter.C:
			for _, device := range g.devices {
				g.pollParameterData(device)
			}

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

func (g *GrowattService) pollBatteryDetails(device models.NoahDevicePayload) {

	if history, err := g.client.GetNoahHistory(device.Serial, "", ""); err != nil {
		slog.Error("could not get device history", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		if len(history.Obj.Datas) == 0 {
			slog.Info("could not get device history, data empty", slog.String("device", device.Serial))
		} else {
			historyData := history.Obj.Datas[0]

			var batteries []models.BatteryPayload
			for i := 0; i < len(device.Batteries); i++ {
				batteries = append(batteries, batteryPayload(historyData, i))
			}

			g.endpoint.PublishBatteryDetails(device, batteries)
		}
	}
}
