package growatt_app

import (
	"nexa-mqtt/internal/misc"
	"nexa-mqtt/pkg/models"
)

func devicePayload(n *NoahStatus) models.DevicePayload {
	return models.DevicePayload{
		ACPower:               misc.ParseFloat(n.Obj.Pac),
		SolarPower:            misc.ParseFloat(n.Obj.Ppv),
		Soc:                   misc.ParseFloat(n.Obj.Soc),
		ChargePower:           misc.ParseFloat(n.Obj.ChargePower),
		DischargePower:        misc.ParseFloat(n.Obj.DisChargePower),
		BatteryNum:            int(misc.ParseFloat(n.Obj.BatteryNum)),
		GenerationTotalEnergy: misc.ParseFloat(n.Obj.EacTotal),
		GenerationTodayEnergy: misc.ParseFloat(n.Obj.EacToday),
		WorkMode:              models.WorkModeFromString(n.Obj.WorkMode),
		Status:                models.StatusFromString(n.Obj.Status),
	}
}

func batteryPayload(n *BatteryDetails) models.BatteryPayload {
	return models.BatteryPayload{
		SerialNumber: n.SerialNum,
		Soc:          misc.ParseFloat(n.Soc),
		Temperature:  misc.ParseFloat(n.Temp),
	}
}

func parameterPayload(n *NexaInfo) models.ParameterPayload {
	chargingLimit := misc.ParseFloat(n.Obj.Noah.ChargingSocHighLimit)
	dischargeLimit := misc.ParseFloat(n.Obj.Noah.ChargingSocLowLimit)
	defaultACCouplePower := misc.ParseFloat(n.Obj.Noah.DefaultACCouplePower)
	defaultMode := models.WorkModeFromString(n.Obj.Noah.DefaultMode)
	antiBackflowPowerPercentage := misc.ParseFloat(n.Obj.Noah.AntiBackflowPowerPercentage)

	return models.ParameterPayload{
		ChargingLimit:               &chargingLimit,
		DischargeLimit:              &dischargeLimit,
		DefaultACCouplePower:        &defaultACCouplePower,
		DefaultMode:                 &defaultMode,
		AllowGridCharging:           misc.IntStringToOnOff(n.Obj.Noah.AllowGridCharging),
		GridConnectionControl:       misc.IntStringToOnOff(n.Obj.Noah.GridConnectionControl),
		AcCouplePowerControl:        misc.IntStringToOnOff(n.Obj.Noah.AcCouplePowerControl),
		AcCoupleEnable:              misc.IntStringToOnOff(n.Obj.Noah.AcCoupleEnable),
		LightLoadEnable:             misc.IntStringToOnOff(n.Obj.Noah.LightLoadEnable),
		NeverPowerOff:               misc.IntStringToOnOff(n.Obj.Noah.NeverPowerOff),
		AntiBackflowEnable:          misc.IntStringToOnOff(n.Obj.Noah.AntiBackflowEnable),
		AntiBackflowPowerPercentage: &antiBackflowPowerPercentage,
	}
}
