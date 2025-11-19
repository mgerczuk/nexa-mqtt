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
		{
			CommonConfig: CommonConfig{
				Name:     "LightLoadEnable",
				UniqueId: fmt.Sprintf("%s_light_load_enable", info.SerialNumber),
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.light_load_enable }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"light_load_enable\": \"{{ value }}\"}",
			},
		},
		{
			CommonConfig: CommonConfig{
				Name:     "NeverPowerOff",
				UniqueId: fmt.Sprintf("%s_never_power_off", info.SerialNumber),
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.never_power_off }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"never_power_off\": \"{{ value }}\"}",
			},
		},
		{
			CommonConfig: CommonConfig{
				Name:     "AntiBackflowEnable",
				UniqueId: fmt.Sprintf("%s_anti_backflow_enable", info.SerialNumber),
				Device:   device,
				Origin:   origin,
			},
			StateConfig: StateConfig{
				StateTopic:    info.ParameterStateTopic,
				ValueTemplate: "{{ value_json.anti_backflow_enable }}",
			},
			CommandConfig: CommandConfig{
				CommandTopic:    info.ParameterCommandTopic,
				CommandTemplate: "{\"anti_backflow_enable\": \"{{ value }}\"}",
			},
		},
	}
	return switches
}
