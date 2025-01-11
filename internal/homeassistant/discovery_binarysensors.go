package homeassistant

import (
	"fmt"
	"noah-mqtt/pkg/models"
)

func generateBinarySensorDiscoveryPayload(appVersion string, info DeviceInfo) []BinarySensor {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	binarySensors := []BinarySensor{
		{
			Name:          "Connectivity",
			Icon:          "",
			DeviceClass:   DeviceClassConnectivity,
			ValueTemplate: fmt.Sprintf("{{ 'offline' if value_json.status == '%s' else 'online' }}", models.Offline),
			PayloadOff:    "offline",
			PayloadOn:     "online",
			UniqueId:      fmt.Sprintf("%s_connectivity", info.SerialNumber),
			StateTopic:    info.StateTopic,
			Device:        device,
			Origin:        origin,
		},
		{
			Name:          "Heating",
			Icon:          IconHeatWave,
			DeviceClass:   DeviceClassNone,
			ValueTemplate: fmt.Sprintf("{{ 'heating' if value_json.status == '%s' else 'not-heating' }}", models.Heating),
			PayloadOff:    "not-heating",
			PayloadOn:     "heating",
			UniqueId:      fmt.Sprintf("%s_heating", info.SerialNumber),
			StateTopic:    info.StateTopic,
			Device:        device,
			Origin:        origin,
		},
	}

	return binarySensors
}
