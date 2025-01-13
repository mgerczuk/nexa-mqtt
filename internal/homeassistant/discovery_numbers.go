package homeassistant

import "fmt"

func generateNumberDiscoveryPayload(appVersion string, info DeviceInfo) []Number {
	device := generateDevice(info)
	origin := generateOrigin(appVersion)

	numbers := []Number{
		{
			Name:              "System Output Power",
			UniqueId:          fmt.Sprintf("%s_system_output_power", info.SerialNumber),
			CommandTemplate:   "{\"output_power_w\": {{ value }}}",
			CommandTopic:      info.ParameterCommandTopic,
			Device:            device,
			Origin:            origin,
			Icon:              "",
			DeviceClass:       DeviceClassPower,
			StateTopic:        info.ParameterStateTopic,
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               0,
			Max:               800,
			UnitOfMeasurement: UnitWatt,
			ValueTemplate:     "{{ value_json.output_power_w }}",
		},
		{
			Name:              "Charging Limit",
			UniqueId:          fmt.Sprintf("%s_charging_limit", info.SerialNumber),
			CommandTemplate:   "{\"charging_limit\": {{ value }}}",
			CommandTopic:      info.ParameterCommandTopic,
			Device:            device,
			Origin:            origin,
			Icon:              IconBatteryArrowUpOutline,
			StateTopic:        info.ParameterStateTopic,
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               70,
			Max:               100,
			UnitOfMeasurement: UnitPercent,
			ValueTemplate:     "{{ value_json.charging_limit }}",
		},
		{
			Name:              "Discharge Limit",
			UniqueId:          fmt.Sprintf("%s_discharge_limit", info.SerialNumber),
			CommandTemplate:   "{\"discharge_limit\": {{ value }}}",
			CommandTopic:      info.ParameterCommandTopic,
			Device:            device,
			Origin:            origin,
			Icon:              IconBatteryArrowDownOutline,
			StateTopic:        info.ParameterStateTopic,
			StateClass:        StateClassMeasurement,
			Mode:              ModeSlider,
			Step:              1,
			Min:               0,
			Max:               30,
			UnitOfMeasurement: UnitPercent,
			ValueTemplate:     "{{ value_json.discharge_limit }}",
		},
	}

	return numbers
}
