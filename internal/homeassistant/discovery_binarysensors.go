package homeassistant

import (
	"fmt"
	"nexa-mqtt/pkg/models"
)

func generateBinarySensorDiscoveryPayload(appVersion string, info DeviceInfo) []BinarySensor {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	binarySensors := []BinarySensor{
		{
			CommonConfig: CommonConfig{
				Name:        "Connectivity",
				UniqueId:    fmt.Sprintf("%s_connectivity", info.SerialNumber),
				Icon:        "",
				DeviceClass: DeviceClassConnectivity,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic: info.StateTopic(),
				// value_json.status may be "online", "offline", "heating", ...
				ValueTemplate: fmt.Sprintf("{{ 'offline' if value_json.status == '%s' else 'online' }}", models.Offline),
			},
			PayloadOff: "offline",
			PayloadOn:  "online",
		},
		{
			CommonConfig: CommonConfig{
				Name:        "Heating",
				UniqueId:    fmt.Sprintf("%s_heating", info.SerialNumber),
				Icon:        IconHeatWave,
				DeviceClass: DeviceClassNone,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic: info.StateTopic(),
				// value_json.status may be "online", "offline", "heating", ...
				ValueTemplate: fmt.Sprintf("{{ 'heating' if value_json.status == '%s' else 'not-heating' }}", models.Heating),
			},
			PayloadOff: "not-heating",
			PayloadOn:  "heating",
		},
		{
			CommonConfig: CommonConfig{
				Name:        "API Health",
				UniqueId:    fmt.Sprintf("%s_api_health", info.SerialNumber),
				DeviceClass: DeviceClassConnectivity,
				Device:      device,
				Origin:      origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.HealthTopic(),
				ValueTemplate: "{{ value_json.status }}",
			},
			PayloadOff: "error",
			PayloadOn:  "ok",
		},
	}

	return binarySensors
}
