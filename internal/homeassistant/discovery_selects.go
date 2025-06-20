package homeassistant

import (
	"fmt"
	"noah-mqtt/pkg/models"
)

func generateSelectDiscoveryPayload(appVersion string, info DeviceInfo) []Select {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	selects := []Select{
		{
			Name:            "Default Mode",
			UniqueId:        fmt.Sprintf("%s_%s", info.SerialNumber, "default_mode"),
			CommandTemplate: "{\"default_mode\": \"{{ value }}\"}",
			CommandTopic:    info.ParameterCommandTopic,
			Device:          device,
			Origin:          origin,
			DeviceClass:     DeviceClassEnum,
			Options:         []string{models.WorkModeLoadFirst, models.WorkModeBatteryFirst},
			StateTopic:      info.ParameterStateTopic,
			ValueTemplate:   "{{ value_json.default_mode }}",
			Component:       "select",
		},
	}

	return selects
}
