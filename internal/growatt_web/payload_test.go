package growatt_web

import (
	"nexa-mqtt/pkg/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_devicePayload(t *testing.T) {
	device := models.NoahDevicePayload{
		Batteries: []models.NoahDeviceBatteryPayload{{}, {}},
	}
	status := GrowattNoahStatusObj{
		Pac:                           "-330",
		Ppv:                           "178",
		TotalBatteryPackSoc:           "72",
		TotalBatteryPackChargingPower: "248",
		WorkMode:                      "1",
		Status:                        "2",
	}
	totals := GrowattNoahTotalsObj{
		EacTotal: "328.4",
		EacToday: "4.2",
	}

	payload := devicePayload(device, status, totals)

	assert.Equal(t, -330.0, payload.ACPower)
	assert.Equal(t, 178.0, payload.SolarPower)
	assert.Equal(t, 72.0, payload.Soc)
	assert.Equal(t, 248.0, payload.ChargePower)
	assert.Equal(t, 0.0, payload.DischargePower)
	assert.Equal(t, 2, payload.BatteryNum)
	assert.Equal(t, 328.4, payload.GenerationTotalEnergy)
	assert.Equal(t, 4.2, payload.GenerationTodayEnergy)
	assert.Equal(t, models.WorkMode(models.WorkModeBatteryFirst), payload.WorkMode)
	assert.Equal(t, models.SmartSelfUse, payload.Status)
}

func Test_batteryPayload(t *testing.T) {
	historyData := GrowattNoahHistoryData{
		Battery1SerialNum: "Serial123",
		Battery1Soc:       44,
		Battery1Temp:      35.0,

		Battery2SerialNum: "Serial223",
		Battery2Soc:       55,
		Battery2Temp:      36.0,

		Battery3SerialNum: "Serial323",
		Battery3Soc:       66,
		Battery3Temp:      37.0,

		Battery4SerialNum: "Serial423",
		Battery4Soc:       77,
		Battery4Temp:      38.0,
	}
	tm := time.Now().Truncate(time.Second)

	payload := batteryPayload(historyData, tm, 0)

	assert.Equal(t, models.BatteryPayload{
		Time:         tm,
		SerialNumber: "Serial123",
		Soc:          44.0,
		Temperature:  35.0,
	}, payload)

	payload = batteryPayload(historyData, tm, 1)

	assert.Equal(t, models.BatteryPayload{
		Time:         tm,
		SerialNumber: "Serial223",
		Soc:          55.0,
		Temperature:  36.0,
	}, payload)

	payload = batteryPayload(historyData, tm, 2)

	assert.Equal(t, models.BatteryPayload{
		Time:         tm,
		SerialNumber: "Serial323",
		Soc:          66.0,
		Temperature:  37.0,
	}, payload)

	payload = batteryPayload(historyData, tm, 3)

	assert.Equal(t, models.BatteryPayload{
		Time:         tm,
		SerialNumber: "Serial423",
		Soc:          77.0,
		Temperature:  38.0,
	}, payload)

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	payload = batteryPayload(historyData, tm, 4)
	t.Errorf("Test failed, panic was expected")
}

func Test_parameterPayload(t *testing.T) {
	detailsData := GrowattNoahListData{
		ChargingSocHighLimit:        "95",
		ChargingSocLowLimit:         "11",
		DefaultACCouplePower:        "270",
		DefaultMode:                 "0",
		AllowGridCharging:           "1",
		GridConnectionControl:       "0",
		AcCouplePowerControl:        "1",
		LightLoadEnable:             "1",
		NeverPowerOff:               "0",
		AntiBackflowEnable:          "1",
		AntiBackflowPowerPercentage: "77",
	}

	payload := parameterPayload(detailsData)

	assert.Equal(t, 95.0, *payload.ChargingLimit)
	assert.Equal(t, 11.0, *payload.DischargeLimit)
	assert.Equal(t, 270.0, *payload.DefaultACCouplePower)
	assert.Equal(t, models.WorkMode(models.WorkModeLoadFirst), *payload.DefaultMode)
	assert.Equal(t, models.ON, payload.AllowGridCharging)
	assert.Equal(t, models.OFF, payload.GridConnectionControl)
	assert.Equal(t, models.ON, payload.AcCouplePowerControl)
	assert.Equal(t, models.ON, payload.LightLoadEnable)
	assert.Equal(t, models.OFF, payload.NeverPowerOff)
	assert.Equal(t, models.ON, payload.AntiBackflowEnable)
	assert.Equal(t, 77.0, *payload.AntiBackflowPowerPercentage)
}
