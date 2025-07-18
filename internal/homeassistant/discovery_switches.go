package homeassistant

import (
	"fmt"
)

func generateSwitchDiscoveryPayload(appVersion string, info DeviceInfo) []Switch {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	switches := []Switch{
		{
			CommonConfig: CommonConfig{
				Name:     "AllowGridCharging",
				UniqueId: fmt.Sprintf("%s_allow_grid_charging", info.SerialNumber),
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.allow_grid_charging }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"allow_grid_charging\": \"{{ value }}\"}",
			},
		},
		{
			CommonConfig: CommonConfig{
				Name:     "GridConnectionControl",
				UniqueId: fmt.Sprintf("%s_grid_connection_control", info.SerialNumber),
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.grid_connection_control }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"grid_connection_control\": \"{{ value }}\"}",
			},
		},
		{
			CommonConfig: CommonConfig{
				Name:     "AcCouplePowerControl",
				UniqueId: fmt.Sprintf("%s_ac_couple_power_control", info.SerialNumber),
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.ac_couple_power_control }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"ac_couple_power_control\": \"{{ value }}\"}",
			},
		},
	}
	return switches
}
