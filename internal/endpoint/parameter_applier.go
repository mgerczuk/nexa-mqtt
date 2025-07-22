package endpoint

import "nexa-mqtt/pkg/models"

type ParameterApplier interface {
	SetOutputPowerW(device models.NoahDevicePayload, mode models.WorkMode, power float64) error
	SetChargingLimits(device models.NoahDevicePayload, chargingLimit float64, dischargeLimit float64) error
	SetAllowGridCharging(device models.NoahDevicePayload, allow models.OnOff) error
	SetGridConnectionControl(device models.NoahDevicePayload, offlineEnable models.OnOff) error
	SetAcCouplePowerControl(device models.NoahDevicePayload, _1000WEnable models.OnOff) error
}
