package growatt_app

import (
	"errors"
	"nexa-mqtt/pkg/models"
	"testing"
)

// ----- Test functions -----------------------------------------------------

func Test_pollStatus_Ok(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetNoahStatus(
		device.Serial,
		NoahStatusObj{
			LoadPower:     "400",
			GridPower:     "0",
			ChargePower:   "132",
			GroplugPower:  "0",
			WorkMode:      "0",
			Soc:           "93",
			EastronStatus: "-1",
			//AssociatedInvSn: nil,
			BatteryNum:     "1",
			ProfitToday:    "0",
			PlantID:        "10421077",
			DisChargePower: "0",
			EacTotal:       "9.6",
			EacToday:       "3.3",
			IsHaveCt:       "false",
			OnOffGrid:      "0",
			Pac:            "-400",
			Ppv:            "538",
			Alias:          "NEXA 2000",
			ProfitTotal:    "0",
			MoneyUnit:      "â¬",
			GroplugNum:     "0",
			OtherPower:     "-400",
			Status:         "6",
		},
		nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			OutputPower:           -400.0,
			SolarPower:            538.0,
			Soc:                   93.0,
			ChargePower:           132.0,
			DischargePower:        0.0,
			BatteryNum:            1,
			GenerationTotalEnergy: 9.6,
			GenerationTodayEnergy: 3.3,
			WorkMode:              models.WorkMode("load_first"),
			Status:                "on_grid",
		},
	)

	service.pollStatus(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollStatus_Fail(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetNoahStatus(
		device.Serial,
		NoahStatusObj{},
		errors.New("pollStatus fails"))

	service.pollStatus(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_Ok(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetBatteryData(
		device.Serial,
		BatteryInfoObj{
			Batter: []BatteryDetails{
				{
					Temp:      "39",
					SerialNum: "serial124",
					Soc:       "93",
				},
				{
					Temp:      "41",
					SerialNum: "serial125",
					Soc:       "78",
				},
			},
		},
		nil,
	)

	mockEndpoint.On(
		"PublishBatteryDetails",
		device,
		[]models.BatteryPayload{
			{SerialNumber: "serial124", Soc: 93.0, Temperature: 39.0},
			{SerialNumber: "serial125", Soc: 78.0, Temperature: 41.0},
		},
	)

	service.pollBatteryDetails(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_Fail(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetBatteryData(
		device.Serial,
		BatteryInfoObj{},
		errors.New("pollBatteryDetails fails"),
	)

	service.pollBatteryDetails(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollParameterData_Ok(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)

	nexaInfo := NexaInfoObj{}

	nexaInfo.Noah.ChargingSocHighLimit = "95"
	nexaInfo.Noah.DefaultMode = "0"
	nexaInfo.Noah.DefaultACCouplePower = "100"
	nexaInfo.Noah.ChargingSocLowLimit = "11"
	nexaInfo.Noah.AllowGridCharging = "0"
	nexaInfo.Noah.GridConnectionControl = "1"
	nexaInfo.Noah.AcCouplePowerControl = "1"

	chargingLimit := 95.0
	dischargeLimit := 11.0
	defaultACCouplePower := 100.0
	defaultMode := models.WorkMode("load_first")
	allowGridCharging := false
	gridConnectionControl := true
	acCouplePowerControl := true

	mockHttpClient.OnGetNoahInfo(
		device.Serial,
		nexaInfo,
		nil,
	)

	mockEndpoint.On(
		"PublishParameterData",
		device,
		models.ParameterPayload{
			ChargingLimit:         &chargingLimit,
			DischargeLimit:        &dischargeLimit,
			DefaultACCouplePower:  &defaultACCouplePower,
			DefaultMode:           &defaultMode,
			AllowGridCharging:     &allowGridCharging,
			GridConnectionControl: &gridConnectionControl,
			AcCouplePowerControl:  &acCouplePowerControl,
		},
	)

	service.pollParameterData(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollParameterData_Fail(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetNoahInfo(
		device.Serial,
		NexaInfoObj{},
		errors.New("pollParameterData fails"),
	)

	service.pollParameterData(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}
