package homeassistant

import "fmt"

type DeviceInfo struct {
	SerialNumber string
	Model        string
	Version      string
	Alias        string
	TopicPrefix  string
	Batteries    []BatteryInfo
	PVs          []PVInfo
}

func (d DeviceInfo) StateTopic() string {
	return fmt.Sprintf("%s/%s", d.TopicPrefix, d.SerialNumber)
}

func (d DeviceInfo) ParameterStateTopic() string {
	return fmt.Sprintf("%s/%s/parameters", d.TopicPrefix, d.SerialNumber)
}

func (d DeviceInfo) ParameterCommandTopic() string {
	return fmt.Sprintf("%s/%s/parameters/set", d.TopicPrefix, d.SerialNumber)
}

func (d DeviceInfo) HealthTopic() string {
	return fmt.Sprintf("%s/%s/health", d.TopicPrefix, d.SerialNumber)
}

func (d DeviceInfo) AvailabilityTopic() string {
	return fmt.Sprintf("%s/availability", d.TopicPrefix)
}

type BatteryInfo struct {
	Alias      string
	StateTopic string
}

type PVInfo struct {
	StateTopic string
}
