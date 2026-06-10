package growatt_app

import (
	"log/slog"
	"nexa-mqtt/pkg/models"
)

func (g *GrowattAppService) pollStatus(device models.NoahDevicePayload) {
	if data, err := g.client.GetSystemStatus(device.Serial); err != nil {
		slog.Error("could not get device data", slog.String("error", err.Error()), slog.String("device", device.Serial))
		g.health.UpdateError(device.Serial, err)
	} else {
		payload := devicePayload(data)
		g.endpoint.PublishDeviceStatus(device, payload)
		g.health.UpdateSuccess(device.Serial)
	}
	g.endpoint.PublishHealth(device, &g.health)
}

func (g *GrowattAppService) pollBatteryDetails(device models.NoahDevicePayload) {
	if data, err := g.client.GetBatteryData(device.Serial); err != nil {
		slog.Error("could not get battery data", slog.String("error", err.Error()), slog.String("device", device.Serial))
		g.health.UpdateError(device.Serial, err)
	} else {
		var batteryPayloads []models.BatteryPayload

		for _, bat := range data.Obj.Batter {
			batteryPayloads = append(batteryPayloads, batteryPayload(&bat))
		}

		g.endpoint.PublishBatteryDetails(device, batteryPayloads)
		g.health.UpdateSuccess(device.Serial)
	}
	g.endpoint.PublishHealth(device, &g.health)
}

func (g *GrowattAppService) pollParameterData(device models.NoahDevicePayload) {
	if data, err := g.client.GetNexaInfoBySn(device.Serial); err != nil {
		slog.Error("could not get parameter data", slog.String("error", err.Error()), slog.String("device", device.Serial))
		g.health.UpdateError(device.Serial, err)
	} else {
		payload := parameterPayload(data)
		g.endpoint.PublishParameterData(device, payload)
		g.health.UpdateSuccess(device.Serial)
	}
	g.endpoint.PublishHealth(device, &g.health)
}
