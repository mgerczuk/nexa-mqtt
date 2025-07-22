package growatt_web

import (
	"fmt"
	"nexa-mqtt/internal/misc"
	"nexa-mqtt/pkg/models"
)

func devicePayload(device models.NoahDevicePayload, status GrowattNoahStatusObj, totals GrowattNoahTotalsObj) models.DevicePayload {
	batteryPower := misc.ParseFloat(status.TotalBatteryPackChargingPower)

	chargePower := 0.0
	dischargePower := 0.0
	if batteryPower < 0 {
		dischargePower = -batteryPower
	} else {
		chargePower = batteryPower
	}

	payload := models.DevicePayload{
		ACPower:               misc.ParseFloat(status.Pac),
		SolarPower:            misc.ParseFloat(status.Ppv),
		Soc:                   misc.ParseFloat(status.TotalBatteryPackSoc),
		ChargePower:           chargePower,
		DischargePower:        dischargePower,
		BatteryNum:            len(device.Batteries),
		GenerationTotalEnergy: misc.ParseFloat(totals.EacTotal),
		GenerationTodayEnergy: misc.ParseFloat(totals.EacToday),
		WorkMode:              models.WorkModeFromString(status.WorkMode),
		Status:                models.StatusFromString(status.Status),
	}
	return payload
}

func batteryPayload(historyData GrowattNoahHistoryData, i int) models.BatteryPayload {
	switch i {
	case 0:
		return models.BatteryPayload{
			SerialNumber: historyData.Battery1SerialNum,
			Soc:          float64(historyData.Battery1Soc),
			Temperature:  historyData.Battery1Temp,
		}
	case 1:
		return models.BatteryPayload{
			SerialNumber: historyData.Battery2SerialNum,
			Soc:          float64(historyData.Battery2Soc),
			Temperature:  historyData.Battery2Temp,
		}
	case 2:
		return models.BatteryPayload{
			SerialNumber: historyData.Battery3SerialNum,
			Soc:          float64(historyData.Battery3Soc),
			Temperature:  historyData.Battery3Temp,
		}
	case 3:
		return models.BatteryPayload{
			SerialNumber: historyData.Battery4SerialNum,
			Soc:          float64(historyData.Battery4Soc),
			Temperature:  historyData.Battery4Temp,
		}
	}

	panic(fmt.Errorf("growatt_web.batteryPayload: invalid index %d", i))
}

func parameterPayload(detailsData GrowattNoahListData) models.ParameterPayload {
	cl := misc.ParseFloat(detailsData.ChargingSocHighLimit)
	dl := misc.ParseFloat(detailsData.ChargingSocLowLimit)
	op := misc.ParseFloat(detailsData.DefaultACCouplePower)
	mode := models.WorkModeFromString(detailsData.DefaultMode)
	agc := misc.IntStringToOnOff(detailsData.AllowGridCharging)
	gcc := misc.IntStringToOnOff(detailsData.GridConnectionControl)
	acpc := misc.IntStringToOnOff(detailsData.AcCouplePowerControl)
	paramPayload := models.ParameterPayload{
		ChargingLimit:         &cl,
		DischargeLimit:        &dl,
		DefaultACCouplePower:  &op,
		DefaultMode:           &mode,
		AllowGridCharging:     agc,
		GridConnectionControl: gcc,
		AcCouplePowerControl:  acpc,
	}
	return paramPayload
}
