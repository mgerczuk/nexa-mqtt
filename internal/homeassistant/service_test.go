package homeassistant

import (
	"strings"
	"testing"
)

func Test_sendDiscovery(t *testing.T) {

	mockClient := MockMqttClient{}
	service := &Service{
		options: Options{
			MqttClient:  &mockClient,
			TopicPrefix: "homeassistant",
			Version:     "version",
		},
		devices: []DeviceInfo{
			{
				SerialNumber:          "device123",
				StateTopic:            "test/device123",
				ParameterStateTopic:   "test/device123/parameters",
				ParameterCommandTopic: "test/device123/parameters/set",
				Batteries: []BatteryInfo{
					{
						Alias:      "BAT0",
						StateTopic: "test/device123/BAT0",
					},
					{
						Alias:      "BAT1",
						StateTopic: "test/device123/BAT1",
					},
				},
			},
			{
				SerialNumber:          "device234",
				StateTopic:            "test/device234",
				ParameterStateTopic:   "test/device234/parameters",
				ParameterCommandTopic: "test/device234/parameters/set",
				Batteries: []BatteryInfo{
					{
						Alias:      "BAT0",
						StateTopic: "test/device234/BAT0",
					},
				},
			},
		},
	}

	setupTopics(&mockClient, "device123")
	setupSwitchTopics(&mockClient, "device123")
	setupBatteryTopics(&mockClient, "device123", "BAT0")
	setupBatteryTopics(&mockClient, "device123", "BAT1")
	setupTopics(&mockClient, "device234")
	setupSwitchTopics(&mockClient, "device234")
	setupBatteryTopics(&mockClient, "device234", "BAT0")

	service.sendDiscovery()

	mockClient.AssertExpectations(t)
}

func Test_sendDiscoverySwitchAsSelect(t *testing.T) {

	mockClient := MockMqttClient{}
	service := &Service{
		options: Options{
			MqttClient:     &mockClient,
			TopicPrefix:    "homeassistant",
			Version:        "version",
			SwitchAsSelect: true,
		},
		devices: []DeviceInfo{
			{
				SerialNumber:          "device123",
				StateTopic:            "test/device123",
				ParameterStateTopic:   "test/device123/parameters",
				ParameterCommandTopic: "test/device123/parameters/set",
				Batteries: []BatteryInfo{
					{
						Alias:      "BAT0",
						StateTopic: "test/device123/BAT0",
					},
					{
						Alias:      "BAT1",
						StateTopic: "test/device123/BAT1",
					},
				},
			},
			{
				SerialNumber:          "device234",
				StateTopic:            "test/device234",
				ParameterStateTopic:   "test/device234/parameters",
				ParameterCommandTopic: "test/device234/parameters/set",
				Batteries: []BatteryInfo{
					{
						Alias:      "BAT0",
						StateTopic: "test/device234/BAT0",
					},
				},
			},
		},
	}

	setupTopics(&mockClient, "device123")
	setupSwitchTopicsAsSelect(&mockClient, "device123")
	setupBatteryTopics(&mockClient, "device123", "BAT0")
	setupBatteryTopics(&mockClient, "device123", "BAT1")
	setupTopics(&mockClient, "device234")
	setupSwitchTopicsAsSelect(&mockClient, "device234")
	setupBatteryTopics(&mockClient, "device234", "BAT0")

	service.sendDiscovery()

	mockClient.AssertExpectations(t)
}

func setupTopics(mockClient *MockMqttClient, serial string) {
	r := strings.NewReplacer("$SERIAL", serial)
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/OutputPower/config"),
		r.Replace(`{"name":"Output Power","unique_id":"$SERIAL_output_power","device_class":"power","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.output_w }}","state_class":"measurement","unit_of_measurement":"W"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/SolarPower/config"),
		r.Replace(`{"name":"Solar Power","unique_id":"$SERIAL_solar_power","icon":"mdi:solar-power","device_class":"power","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.solar_w }}","state_class":"measurement","unit_of_measurement":"W"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/SoC/config"),
		r.Replace(`{"name":"SoC","unique_id":"$SERIAL_soc","device_class":"battery","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.soc }}","state_class":"measurement","unit_of_measurement":"%"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/ChargingPower/config"),
		r.Replace(`{"name":"Charging Power","unique_id":"$SERIAL_charging_power","icon":"mdi:battery-plus","device_class":"power","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.charge_w }}","state_class":"measurement","unit_of_measurement":"W"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/DischargePower/config"),
		r.Replace(`{"name":"Discharge Power","unique_id":"$SERIAL_discharge_power","icon":"mdi:battery-minus","device_class":"power","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.discharge_w }}","state_class":"measurement","unit_of_measurement":"W"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/NumberOfBatteries/config"),
		r.Replace(`{"name":"Number Of Batteries","unique_id":"$SERIAL_battery_num","icon":"mdi:car-battery","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.battery_num }}","state_class":"measurement"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/GenerationTotal/config"),
		r.Replace(`{"name":"Generation Total","unique_id":"$SERIAL_generation_total","device_class":"energy","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.generation_total_kwh }}","state_class":"total_increasing","unit_of_measurement":"kWh"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/GenerationToday/config"),
		r.Replace(`{"name":"Generation Today","unique_id":"$SERIAL_generation_today","device_class":"energy","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.generation_today_kwh }}","state_class":"total_increasing","unit_of_measurement":"kWh"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/WorkingMode/config"),
		r.Replace(`{"name":"Working Mode","unique_id":"$SERIAL_work_mode","device_class":"enum","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.work_mode }}","options":["load_first","battery_first"]}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/Status/config"),
		r.Replace(`{"name":"Status","unique_id":"$SERIAL_status","device_class":"enum","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ value_json.status }}","options":["offline","load_first","battery_first","smart_self_use","fault","heating","on_grid","off_grid"]}`))

	mockClient.OnPublish(
		r.Replace("homeassistant/number/nexa_$SERIAL/ChargingLimit/config"),
		r.Replace(`{"name":"Charging Limit","unique_id":"$SERIAL_charging_limit","icon":"mdi:battery-arrow-up-outline","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.charging_limit }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"charging_limit\": {{ value }}}","state_class":"measurement","unit_of_measurement":"%","mode":"slider","step":1,"min":70,"max":100}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/number/nexa_$SERIAL/DischargeLimit/config"),
		r.Replace(`{"name":"Discharge Limit","unique_id":"$SERIAL_discharge_limit","icon":"mdi:battery-arrow-down-outline","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.discharge_limit }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"discharge_limit\": {{ value }}}","state_class":"measurement","unit_of_measurement":"%","mode":"slider","step":1,"max":30}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/number/nexa_$SERIAL/DefaultACOutputPower/config"),
		r.Replace(`{"name":"Default AC Output Power","unique_id":"$SERIAL_default_output_w","device_class":"power","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.default_output_w }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"default_output_w\": {{ value }}}","state_class":"measurement","unit_of_measurement":"W","mode":"slider","step":1,"max":800}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/select/nexa_$SERIAL/DefaultMode/config"),
		r.Replace(`{"name":"Default Mode","unique_id":"$SERIAL_default_mode","device_class":"enum","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.default_mode }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"default_mode\": \"{{ value }}\"}","options":["load_first","battery_first"],"component":"select"}`))

	mockClient.OnPublish(
		r.Replace("homeassistant/binary_sensor/nexa_$SERIAL/Connectivity/config"),
		r.Replace(`{"name":"Connectivity","unique_id":"$SERIAL_connectivity","device_class":"connectivity","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ 'offline' if value_json.status == 'offline' else 'online' }}","payload_off":"offline","payload_on":"online"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/binary_sensor/nexa_$SERIAL/Heating/config"),
		r.Replace(`{"name":"Heating","unique_id":"$SERIAL_heating","icon":"mdi:heat-wave","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL","value_template":"{{ 'heating' if value_json.status == 'heating' else 'not-heating' }}","payload_off":"not-heating","payload_on":"heating"}`))
}

func setupSwitchTopics(mockClient *MockMqttClient, serial string) {
	r := strings.NewReplacer("$SERIAL", serial)
	mockClient.OnPublish(
		r.Replace("homeassistant/switch/nexa_$SERIAL/AllowGridCharging/config"),
		r.Replace(`{"name":"AllowGridCharging","unique_id":"$SERIAL_allow_grid_charging","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.allow_grid_charging }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"allow_grid_charging\": \"{{ value }}\"}"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/switch/nexa_$SERIAL/GridConnectionControl/config"),
		r.Replace(`{"name":"GridConnectionControl","unique_id":"$SERIAL_grid_connection_control","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.grid_connection_control }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"grid_connection_control\": \"{{ value }}\"}"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/switch/nexa_$SERIAL/AcCouplePowerControl/config"),
		r.Replace(`{"name":"AcCouplePowerControl","unique_id":"$SERIAL_ac_couple_power_control","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.ac_couple_power_control }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"ac_couple_power_control\": \"{{ value }}\"}"}`))
}

func setupSwitchTopicsAsSelect(mockClient *MockMqttClient, serial string) {
	r := strings.NewReplacer("$SERIAL", serial)
	mockClient.OnPublish(
		r.Replace("homeassistant/select/nexa_$SERIAL/AllowGridCharging/config"),
		r.Replace(`{"name":"AllowGridCharging","unique_id":"$SERIAL_allow_grid_charging","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.allow_grid_charging }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"allow_grid_charging\": \"{{ value }}\"}","options":["OFF","ON"]}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/select/nexa_$SERIAL/GridConnectionControl/config"),
		r.Replace(`{"name":"GridConnectionControl","unique_id":"$SERIAL_grid_connection_control","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.grid_connection_control }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"grid_connection_control\": \"{{ value }}\"}","options":["OFF","ON"]}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/select/nexa_$SERIAL/AcCouplePowerControl/config"),
		r.Replace(`{"name":"AcCouplePowerControl","unique_id":"$SERIAL_ac_couple_power_control","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/parameters","value_template":"{{ value_json.ac_couple_power_control }}","command_topic":"test/$SERIAL/parameters/set","command_template":"{\"ac_couple_power_control\": \"{{ value }}\"}","options":["OFF","ON"]}`))
}

func setupBatteryTopics(mockClient *MockMqttClient, serial string, name string) {
	r := strings.NewReplacer("$BAT", name, "$SERIAL", serial)
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/$BATSoC/config"),
		r.Replace(`{"name":"$BAT SoC","unique_id":"$SERIAL_$BAT_soc","device_class":"battery","device":{"identifiers":["nexa_$SERIAL"],"manufacturer":"Growatt","serial_number":"$SERIAL"},"origin":{"name":"nexa-mqtt","sw_version":"version","support_url":"https://github.com/mgerczuk/nexa-mqtt"},"state_topic":"test/$SERIAL/$BAT","value_template":"{{ value_json.soc }}","state_class":"measurement","unit_of_measurement":"%"}`))
	mockClient.OnPublish(
		r.Replace("homeassistant/sensor/nexa_$SERIAL/$BATTemperature/config"),
		r.Replace("{\"name\":\"$BAT Temperature\",\"unique_id\":\"$SERIAL_$BAT_temp\",\"device_class\":\"temperature\",\"device\":{\"identifiers\":[\"nexa_$SERIAL\"],\"manufacturer\":\"Growatt\",\"serial_number\":\"$SERIAL\"},\"origin\":{\"name\":\"nexa-mqtt\",\"sw_version\":\"version\",\"support_url\":\"https://github.com/mgerczuk/nexa-mqtt\"},\"state_topic\":\"test/$SERIAL/$BAT\",\"value_template\":\"{{ value_json.temp }}\",\"state_class\":\"measurement\",\"unit_of_measurement\":\"°C\"}"))
}
