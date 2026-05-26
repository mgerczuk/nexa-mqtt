package homeassistant

type DeviceInfo struct {
	SerialNumber          string
	Model                 string
	Version               string
	Alias                 string
	StateTopic            string
	ParameterStateTopic   string
	ParameterCommandTopic string
	Batteries             []BatteryInfo
	PVs                   []PVInfo
}

type BatteryInfo struct {
	Alias      string
	StateTopic string
}

type PVInfo struct {
	StateTopic string
}
