package growatt_web

import (
	"errors"
	"math/rand"
	"net/http/cookiejar"
	"nexa-mqtt/pkg/models"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSerial(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func setupGrowattServiceMocks(t *testing.T) (*MockHttpClient, *GrowattService, models.NoahDevicePayload, *MockEndpoint) {
	mockHttpClient := MockHttpClient{}
	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)

	client := Client{
		client:    &mockHttpClient,
		serverUrl: "https://openapi.growatt.com",
		username:  "user",
		password:  "secret",
		jar:       jar,
	}

	endpoint := MockEndpoint{}

	service := GrowattService{
		client:   &client,
		endpoint: &endpoint,
	}

	device := models.NoahDevicePayload{
		Serial:    randSerial(10),
		PlantId:   rand.Intn(20000000),
		Batteries: []models.NoahDeviceBatteryPayload{{Alias: "BAT0"}, {Alias: "BAT1"}, {Alias: "BAT2"}, {Alias: "BAT3"}},
	}

	return &mockHttpClient, &service, device, &endpoint
}

func TestServiceLogin_Ok(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnLogin("user", "secret", GrowattResult{}, nil)

	err := service.Login()

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestServiceLogin_Fails(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnLogin("user", "secret", GrowattResult{}, errors.New("login failed"))

	err := service.Login()

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetEndpoint(t *testing.T) {
	_, service, _, _ := setupGrowattServiceMocks(t)

	service.SetEndpoint(nil)

	assert.Equal(t, nil, service.endpoint)

	ep := &MockEndpoint{}
	ep.On("SetParameterApplier", service).Return()

	service.SetEndpoint(ep)

	assert.Equal(t, ep, service.endpoint)
	ep.AssertExpectations(t)
}

func TestSetOutputPowerW_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet2Params(device.Serial, "system_out_put_power", "1", "100", nil, SetResponse{Success: true})

	err := service.SetOutputPowerW(device, models.WorkMode(models.WorkModeBatteryFirst), 100)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetOutputPowerW_InvalidWorkmode(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetOutputPowerW(device, models.WorkMode("invalid"), 100)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetOutputPowerW_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet2Params(device.Serial, "system_out_put_power", "1", "100", nil, SetResponse{Success: false})

	err := service.SetOutputPowerW(device, models.WorkMode(models.WorkModeBatteryFirst), 100)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingLimits_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "charging_soc_low_limit", "10", nil, SetResponse{Success: true})
	mockHttpClient.OnSet1Params(device.Serial, "charging_soc_high_limit", "95", nil, SetResponse{Success: true})

	err := service.SetChargingLimits(device, 95, 10)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingLimits_LowFails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "charging_soc_low_limit", "10", nil, SetResponse{Success: false})
	//mockHttpClient.OnSet1Params(device.Serial, "charging_soc_high_limit", "95", nil, SetResponse{Success: true})

	err := service.SetChargingLimits(device, 95, 10)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingLimits_HighFails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "charging_soc_low_limit", "10", nil, SetResponse{Success: true})
	mockHttpClient.OnSet1Params(device.Serial, "charging_soc_high_limit", "95", nil, SetResponse{Success: false})

	err := service.SetChargingLimits(device, 95, 10)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "allow_grid_charging", "1", nil, SetResponse{Success: true})

	err := service.SetAllowGridCharging(device, models.ON)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Invalid(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetAllowGridCharging(device, models.OnOff("invalid"))

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "allow_grid_charging", "1", nil, SetResponse{Success: false})

	err := service.SetAllowGridCharging(device, models.ON)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "grid_connection_control", "0", nil, SetResponse{Success: true})

	err := service.SetGridConnectionControl(device, models.OFF)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Invalid(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetGridConnectionControl(device, models.OnOff("invalid"))

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "grid_connection_control", "0", nil, SetResponse{Success: false})

	err := service.SetGridConnectionControl(device, models.OFF)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetAcCouplePowerControl_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "ac_couple_power_control", "1", nil, SetResponse{Success: true})

	err := service.SetAcCouplePowerControl(device, models.ON)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetAcCouplePowerControl_Invalid(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetAcCouplePowerControl(device, models.OnOff("invalid"))

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetAcCouplePowerControl_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "ac_couple_power_control", "1", nil, SetResponse{Success: false})

	err := service.SetAcCouplePowerControl(device, models.ON)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetLightLoadEnable_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "light_load_enable", "0", nil, SetResponse{Success: true})

	err := service.SetLightLoadEnable(device, models.OFF)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetLightLoadEnable_Invalid(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetLightLoadEnable(device, models.OnOff("invalid"))

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetLightLoadEnable_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "light_load_enable", "0", nil, SetResponse{Success: false})

	err := service.SetLightLoadEnable(device, models.OFF)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetNeverPowerOff_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "never_power_off", "1", nil, SetResponse{Success: true})

	err := service.SetNeverPowerOff(device, models.ON)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetNeverPowerOff_Invalid(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetNeverPowerOff(device, models.OnOff("invalid"))

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetNeverPowerOff_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet1Params(device.Serial, "never_power_off", "1", nil, SetResponse{Success: false})

	err := service.SetNeverPowerOff(device, models.ON)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetBackflow_Ok(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet2Params(device.Serial, "anti_back_flow_setting", "1", "15", nil, SetResponse{Success: true})

	err := service.SetBackflow(device, models.ON, 15.0)

	assert.NoError(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetBackflow_Invalid(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	err := service.SetBackflow(device, models.OnOff("invalid"), 15.0)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func TestSetBackflow_Fails(t *testing.T) {
	mockHttpClient, service, device, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnSet2Params(device.Serial, "anti_back_flow_setting", "1", "15", nil, SetResponse{Success: false})

	err := service.SetBackflow(device, models.ON, 15.0)

	assert.Error(t, err)
	mockHttpClient.AssertExpectations(t)
}

func Test_enumerateDevices_GetPlantListFails(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetPlantList([]GrowattPlant{}, errors.New("GetPlantList fails"))

	result := service.enumerateDevices()

	assert.Equal(t, 0, len(result))
	mockHttpClient.AssertExpectations(t)
}

func Test_enumerateDevices_GetNoahListFails(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	plantList := []GrowattPlant{
		{PlantId: "1", PlantName: "plant1"},
		{PlantId: "2", PlantName: "plant2"},
	}
	mockHttpClient.OnGetPlantList(plantList, nil)

	mockHttpClient.OnGetNoahList(1, GrowattNoahList{}, errors.New("GrowattNoahList fails"))
	mockHttpClient.OnGetNoahList(2, GrowattNoahList{}, errors.New("GrowattNoahList fails"))

	result := service.enumerateDevices()

	assert.Equal(t, 0, len(result))
	mockHttpClient.AssertExpectations(t)
}

func Test_enumerateDevices_GetNoahHistoryFails(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	plantList := []GrowattPlant{
		{PlantId: "1", PlantName: "plant1"},
		{PlantId: "2", PlantName: "plant2"},
	}
	mockHttpClient.OnGetPlantList(plantList, nil)

	dev1 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{Sn: "Serial123", PlantID: "1", DeviceModel: "NEXA 2000", Version: "10.10.07.07.4016"},
	}}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)
	dev2 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{}}}
	mockHttpClient.OnGetNoahList(2, dev2, nil)

	today := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory("Serial123", today, today, GrowattNoahHistory{}, errors.New("GrowattNoahList fails"))

	result := service.enumerateDevices()

	assert.Equal(t, 0, len(result))
	mockHttpClient.AssertExpectations(t)
}

func Test_enumerateDevices_GetNoahHistoryNoData(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	plantList := []GrowattPlant{
		{PlantId: "1", PlantName: "plant1"},
		{PlantId: "2", PlantName: "plant2"},
	}
	mockHttpClient.OnGetPlantList(plantList, nil)

	dev1 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{Sn: "Serial123", PlantID: "1", DeviceModel: "NEXA 2000", Version: "10.10.07.07.4016"},
	}}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)
	dev2 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{}}}
	mockHttpClient.OnGetNoahList(2, dev2, nil)

	today := time.Now().Format("2006-01-02")
	hist1 := GrowattNoahHistory{Obj: GrowattNoahHistoryObj{}}
	mockHttpClient.OnGetNoahHistory("Serial123", today, today, hist1, nil)

	result := service.enumerateDevices()

	assert.Equal(t, 0, len(result))
	mockHttpClient.AssertExpectations(t)
}

func Test_enumerateDevices_Ok(t *testing.T) {
	mockHttpClient, service, _, _ := setupGrowattServiceMocks(t)

	plantList := []GrowattPlant{
		{PlantId: "1", PlantName: "plant1"},
		{PlantId: "2", PlantName: "plant2"},
	}
	mockHttpClient.OnGetPlantList(plantList, nil)

	dev1 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{Sn: "Serial123", PlantID: "1", DeviceModel: "NEXA 2000", Version: "10.10.07.07.4016"},
		{Sn: "Serial234", PlantID: "1", DeviceModel: "NEXA 2001"},
	}}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)
	dev2 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{Sn: "Serial345", PlantID: "2", DeviceModel: "NEXA 2002"},
	}}}
	mockHttpClient.OnGetNoahList(2, dev2, nil)

	today := time.Now().Format("2006-01-02")
	hist1 := GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{BatteryPackageQuantity: 1},
	}}}
	mockHttpClient.OnGetNoahHistory("Serial123", today, today, hist1, nil)
	hist2 := GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{BatteryPackageQuantity: 1},
	}}}
	mockHttpClient.OnGetNoahHistory("Serial234", today, today, hist2, nil)
	hist3 := GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{BatteryPackageQuantity: 1},
	}}}
	mockHttpClient.OnGetNoahHistory("Serial345", today, today, hist3, nil)

	result := service.enumerateDevices()

	assert.Equal(t, 3, len(result))

	assert.Equal(t, 1, result[0].PlantId)
	assert.Equal(t, "Serial123", result[0].Serial)
	assert.Equal(t, "NEXA 2000", result[0].Model)
	assert.Equal(t, "10.10.07.07.4016", result[0].Version)

	assert.Equal(t, 1, result[1].PlantId)
	assert.Equal(t, "Serial234", result[1].Serial)
	assert.Equal(t, "NEXA 2001", result[1].Model)

	assert.Equal(t, 2, result[2].PlantId)
	assert.Equal(t, "Serial345", result[2].Serial)
	assert.Equal(t, "NEXA 2002", result[2].Model)
	mockHttpClient.AssertExpectations(t)
}

func Test_pollStatus_OkCharge(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{
		Response: Response[GrowattNoahStatusObj]{
			Obj: GrowattNoahStatusObj{
				Pac:                           "-400",
				Ppv:                           "538",
				TotalBatteryPackSoc:           "93",
				TotalBatteryPackChargingPower: "132",
				WorkMode:                      "0",
				Status:                        "6",
			}}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{
		Response: Response[GrowattNoahTotalsObj]{
			Obj: GrowattNoahTotalsObj{
				EacTotal: "9.6",
				EacToday: "3.3",
			}}}, nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			ACPower:               -400.0,
			SolarPower:            538.0,
			Soc:                   93.0,
			ChargePower:           132.0,
			DischargePower:        0.0,
			BatteryNum:            4,
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

func Test_pollStatus_OkDischarge(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{
		Response: Response[GrowattNoahStatusObj]{
			Obj: GrowattNoahStatusObj{
				Pac:                           "-400",
				Ppv:                           "538",
				TotalBatteryPackSoc:           "93",
				TotalBatteryPackChargingPower: "-132",
				WorkMode:                      "0",
				Status:                        "6",
			}}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{
		Response: Response[GrowattNoahTotalsObj]{
			Obj: GrowattNoahTotalsObj{
				EacTotal: "9.6",
				EacToday: "3.3",
			}}}, nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			ACPower:               -400.0,
			SolarPower:            538.0,
			Soc:                   93.0,
			ChargePower:           0.0,
			DischargePower:        132.0,
			BatteryNum:            4,
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

func Test_pollStatus_GetNoahStatusFails(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{}, errors.New("GetNoahStatus fails"))

	service.pollStatus(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollStatus_GetNoahTotalsFails(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{
		Response: Response[GrowattNoahStatusObj]{
			Obj: GrowattNoahStatusObj{
				Pac:                           "-400",
				Ppv:                           "538",
				TotalBatteryPackSoc:           "93",
				TotalBatteryPackChargingPower: "132",
				WorkMode:                      "0",
				Status:                        "6",
			}}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{}, errors.New("GetNoahTotals fails"))

	service.pollStatus(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollParameterData_Ok(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	chargingLimit := 95.0
	dischargeLimit := 11.0
	defaultACCouplePower := 100.0
	defaultMode := models.WorkMode("load_first")
	allowGridCharging := models.OFF
	gridConnectionControl := models.ON
	acCouplePowerControl := models.ON
	lightLoadEnable := models.OFF
	neverPowerOff := models.ON
	antiBackflowEnable := models.ON
	antiBackflowPowerPercentage := 37.0

	mockHttpClient.OnGetNoahDetails(device.PlantId, device.Serial, GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{
			ChargingSocHighLimit:        "95",
			ChargingSocLowLimit:         "11",
			DefaultMode:                 "0",
			DefaultACCouplePower:        "100",
			AllowGridCharging:           "0",
			GridConnectionControl:       "1",
			AcCouplePowerControl:        "1",
			LightLoadEnable:             "0",
			NeverPowerOff:               "1",
			AntiBackflowEnable:          "1",
			AntiBackflowPowerPercentage: "37",
		},
	}}}, nil)

	mockEndpoint.On(
		"PublishParameterData",
		device,
		models.ParameterPayload{
			ChargingLimit:               &chargingLimit,
			DischargeLimit:              &dischargeLimit,
			DefaultACCouplePower:        &defaultACCouplePower,
			DefaultMode:                 &defaultMode,
			AllowGridCharging:           allowGridCharging,
			GridConnectionControl:       gridConnectionControl,
			AcCouplePowerControl:        acCouplePowerControl,
			LightLoadEnable:             lightLoadEnable,
			NeverPowerOff:               neverPowerOff,
			AntiBackflowEnable:          antiBackflowEnable,
			AntiBackflowPowerPercentage: &antiBackflowPowerPercentage,
		},
	)

	service.pollParameterData(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollParameterData_GetNoahDetailsFails(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetNoahDetails(device.PlantId, device.Serial, GrowattNoahList{}, errors.New("GetNoahDetails fails"))

	service.pollParameterData(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollParameterData_GetNoahDetailsNoData(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	mockHttpClient.OnGetNoahDetails(device.PlantId, device.Serial, GrowattNoahList{}, nil)

	service.pollParameterData(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_Ok(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	tm := time.Now().Add(-3 * time.Minute).Truncate(time.Second)
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{
			Time:                   tm.Format("2006-01-02 15:04:05"),
			BatteryPackageQuantity: 4,

			Battery1SerialNum: "serial124",
			Battery1Soc:       93,
			Battery1Temp:      39.0,

			Battery2SerialNum: "serial125",
			Battery2Soc:       78,
			Battery2Temp:      41.0,

			Battery3SerialNum: "serial126",
			Battery3Soc:       82,
			Battery3Temp:      40.0,

			Battery4SerialNum: "serial127",
			Battery4Soc:       66,
			Battery4Temp:      36.0,

			Pv1Voltage: 33.15,
			Pv1Current: 1.72,
			Pv1Temp:    23,

			Pv2Voltage: 30.35,
			Pv2Current: 1.77,
			Pv2Temp:    23,

			Pv3Voltage: 7.13,
			Pv3Current: 0,
			Pv3Temp:    20.5,

			Pv4Voltage: 7.09,
			Pv4Current: 0.03,
			Pv4Temp:    20.5,
		}}}}, nil)

	mockEndpoint.On(
		"PublishBatteryDetails",
		device,
		[]models.BatteryPayload{
			{Time: tm, SerialNumber: "serial124", Soc: 93.0, Temperature: 39.0},
			{Time: tm, SerialNumber: "serial125", Soc: 78.0, Temperature: 41.0},
			{Time: tm, SerialNumber: "serial126", Soc: 82.0, Temperature: 40.0},
			{Time: tm, SerialNumber: "serial127", Soc: 66.0, Temperature: 36.0},
		},
	)

	mockEndpoint.On(
		"PublishPvDetails",
		device,
		[]models.PvPayload{
			{Time: tm, Voltage: 33.15, Current: 1.72, Temp: 23.0},
			{Time: tm, Voltage: 30.35, Current: 1.77, Temp: 23.0},
			{Time: tm, Voltage: 7.13, Current: 0.0, Temp: 20.5},
			{Time: tm, Voltage: 7.09, Current: 0.03, Temp: 20.5},
		},
	)

	lastTimestamp := service.pollBatteryDetails(device, tm.Add(-3*time.Minute))

	assert.True(t, lastTimestamp.Equal(tm))
	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_OnGetNoahHistoryFails(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{}, errors.New("OnGetNoahHistory fails"))

	lastTimestamp := service.pollBatteryDetails(device, time.Time{})

	assert.True(t, lastTimestamp.IsZero())
	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_GetNoahHistoryNoData(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{}, nil)

	lastTimestamp := service.pollBatteryDetails(device, time.Time{})

	assert.True(t, lastTimestamp.IsZero())
	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_InvalidDate(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{
			Time:                   "20023-13-32 25:61:61",
			BatteryPackageQuantity: 4,

			Battery1SerialNum: "serial124",
			Battery1Soc:       93,
			Battery1Temp:      39.0,

			Battery2SerialNum: "serial125",
			Battery2Soc:       78,
			Battery2Temp:      41.0,

			Battery3SerialNum: "serial126",
			Battery3Soc:       82,
			Battery3Temp:      40.0,

			Battery4SerialNum: "serial127",
			Battery4Soc:       66,
			Battery4Temp:      36.0,

			Pv1Voltage: 33.15,
			Pv1Current: 1.72,
			Pv1Temp:    23,

			Pv2Voltage: 30.35,
			Pv2Current: 1.77,
			Pv2Temp:    23,

			Pv3Voltage: 7.13,
			Pv3Current: 0,
			Pv3Temp:    20.5,

			Pv4Voltage: 7.09,
			Pv4Current: 0.03,
			Pv4Temp:    20.5,
		}}}}, nil)

	mockEndpoint.On(
		"PublishBatteryDetails",
		device,
		[]models.BatteryPayload{
			{Time: time.Time{}, SerialNumber: "serial124", Soc: 93.0, Temperature: 39.0},
			{Time: time.Time{}, SerialNumber: "serial125", Soc: 78.0, Temperature: 41.0},
			{Time: time.Time{}, SerialNumber: "serial126", Soc: 82.0, Temperature: 40.0},
			{Time: time.Time{}, SerialNumber: "serial127", Soc: 66.0, Temperature: 36.0},
		},
	)

	mockEndpoint.On(
		"PublishPvDetails",
		device,
		[]models.PvPayload{
			{Time: time.Time{}, Voltage: 33.15, Current: 1.72, Temp: 23.0},
			{Time: time.Time{}, Voltage: 30.35, Current: 1.77, Temp: 23.0},
			{Time: time.Time{}, Voltage: 7.13, Current: 0.0, Temp: 20.5},
			{Time: time.Time{}, Voltage: 7.09, Current: 0.03, Temp: 20.5},
		},
	)

	lastTimestamp := service.pollBatteryDetails(device, time.Time{})

	assert.True(t, lastTimestamp.IsZero())
	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_NoNewDate(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	tm := time.Now().Truncate(time.Second)
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{
			Time:                   tm.Format("2006-01-02 15:04:05"),
			BatteryPackageQuantity: 4,

			Battery1SerialNum: "serial124",
			Battery1Soc:       93,
			Battery1Temp:      39.0,

			Battery2SerialNum: "serial125",
			Battery2Soc:       78,
			Battery2Temp:      41.0,

			Battery3SerialNum: "serial126",
			Battery3Soc:       82,
			Battery3Temp:      40.0,

			Battery4SerialNum: "serial127",
			Battery4Soc:       66,
			Battery4Temp:      36.0,

			Pv1Voltage: 33.15,
			Pv1Current: 1.72,
			Pv1Temp:    23,

			Pv2Voltage: 30.35,
			Pv2Current: 1.77,
			Pv2Temp:    23,

			Pv3Voltage: 7.13,
			Pv3Current: 0,
			Pv3Temp:    20.5,

			Pv4Voltage: 7.09,
			Pv4Current: 0.03,
			Pv4Temp:    20.5,
		}}}}, nil)

	lastTimestamp := service.pollBatteryDetails(device, tm.Add(3*time.Second))

	assert.True(t, lastTimestamp.Equal(tm))
	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func setupPoll(wg *sync.WaitGroup, mockHttpClient *MockHttpClient, device models.NoahDevicePayload, mockEndpoint *MockEndpoint) {
	// ----- enumerateDevices

	plantList := []GrowattPlant{
		{PlantId: "1", PlantName: "plant1"},
	}
	mockHttpClient.OnGetPlantList(plantList, nil)

	dev1 := GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{Sn: device.Serial, PlantID: strconv.Itoa(device.PlantId)},
	}}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)

	today := time.Now().Format("2006-01-02")
	tm := time.Now().Add(-3 * time.Minute).Truncate(time.Second)
	// mockHttpClient.OnGetNoahHistory("Serial123", today, today, hist1, nil) set below

	mockEndpoint.On(
		"SetDevices",
		[]models.NoahDevicePayload{device},
	)

	// ----- pollStatus

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{
		Response: Response[GrowattNoahStatusObj]{
			Obj: GrowattNoahStatusObj{
				Pac:                           "-400",
				Ppv:                           "538",
				TotalBatteryPackSoc:           "93",
				TotalBatteryPackChargingPower: "132",
				WorkMode:                      "0",
				Status:                        "6",
			}}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{
		Response: Response[GrowattNoahTotalsObj]{
			Obj: GrowattNoahTotalsObj{
				EacTotal: "9.6",
				EacToday: "3.3",
			}}}, nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			ACPower:               -400.0,
			SolarPower:            538.0,
			Soc:                   93.0,
			ChargePower:           132.0,
			DischargePower:        0.0,
			BatteryNum:            4,
			GenerationTotalEnergy: 9.6,
			GenerationTodayEnergy: 3.3,
			WorkMode:              models.WorkMode("load_first"),
			Status:                "on_grid",
		},
	).Run(func(args mock.Arguments) { wg.Done() })

	// ----- pollBatteryDetails

	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{
			Time:                   tm.Format("2006-01-02 15:04:05"),
			BatteryPackageQuantity: 4,

			Battery1SerialNum: "serial124",
			Battery1Soc:       93,
			Battery1Temp:      39.0,

			Battery2SerialNum: "serial125",
			Battery2Soc:       78,
			Battery2Temp:      41.0,

			Battery3SerialNum: "serial126",
			Battery3Soc:       82,
			Battery3Temp:      40.0,

			Battery4SerialNum: "serial127",
			Battery4Soc:       66,
			Battery4Temp:      36.0,

			Pv1Voltage: 33.15,
			Pv1Current: 1.72,
			Pv1Temp:    23,

			Pv2Voltage: 30.35,
			Pv2Current: 1.77,
			Pv2Temp:    23,

			Pv3Voltage: 7.13,
			Pv3Current: 0,
			Pv3Temp:    20.5,

			Pv4Voltage: 7.09,
			Pv4Current: 0.03,
			Pv4Temp:    20.5,
		}}}}, nil)

	mockEndpoint.On(
		"PublishBatteryDetails",
		device,
		[]models.BatteryPayload{
			{Time: tm, SerialNumber: "serial124", Soc: 93.0, Temperature: 39.0},
			{Time: tm, SerialNumber: "serial125", Soc: 78.0, Temperature: 41.0},
			{Time: tm, SerialNumber: "serial126", Soc: 82.0, Temperature: 40.0},
			{Time: tm, SerialNumber: "serial127", Soc: 66.0, Temperature: 36.0},
		},
	).Run(func(args mock.Arguments) { wg.Done() })

	mockEndpoint.On(
		"PublishPvDetails",
		device,
		[]models.PvPayload{
			{Time: tm, Voltage: 33.15, Current: 1.72, Temp: 23.0},
			{Time: tm, Voltage: 30.35, Current: 1.77, Temp: 23.0},
			{Time: tm, Voltage: 7.13, Current: 0.0, Temp: 20.5},
			{Time: tm, Voltage: 7.09, Current: 0.03, Temp: 20.5},
		},
	).Run(func(args mock.Arguments) { wg.Done() })

	// ----- pollParameterData

	chargingLimit := 95.0
	dischargeLimit := 11.0
	defaultACCouplePower := 100.0
	defaultMode := models.WorkMode("load_first")
	allowGridCharging := models.ON
	gridConnectionControl := models.OFF
	acCouplePowerControl := models.OFF
	lightLoadEnable := models.ON
	neverPowerOff := models.OFF
	antiBackflowEnable := models.ON
	antiBackflowPowerPercentage := 45.0

	mockHttpClient.OnGetNoahDetails(device.PlantId, device.Serial, GrowattNoahList{PageResponse[GrowattNoahListData]{Datas: []GrowattNoahListData{
		{
			ChargingSocHighLimit:        "95",
			ChargingSocLowLimit:         "11",
			DefaultMode:                 "0",
			DefaultACCouplePower:        "100",
			AllowGridCharging:           "1",
			GridConnectionControl:       "0",
			AcCouplePowerControl:        "0",
			LightLoadEnable:             "1",
			NeverPowerOff:               "0",
			AntiBackflowEnable:          "1",
			AntiBackflowPowerPercentage: "45",
		},
	}}}, nil)

	mockEndpoint.On(
		"PublishParameterData",
		device,
		models.ParameterPayload{
			ChargingLimit:               &chargingLimit,
			DischargeLimit:              &dischargeLimit,
			DefaultACCouplePower:        &defaultACCouplePower,
			DefaultMode:                 &defaultMode,
			AllowGridCharging:           allowGridCharging,
			GridConnectionControl:       gridConnectionControl,
			AcCouplePowerControl:        acCouplePowerControl,
			LightLoadEnable:             lightLoadEnable,
			NeverPowerOff:               neverPowerOff,
			AntiBackflowEnable:          antiBackflowEnable,
			AntiBackflowPowerPercentage: &antiBackflowPowerPercentage,
		},
	).Run(func(args mock.Arguments) { wg.Done() })
}

const millisRate = 10

type MockDurationCalculator struct {
}

func (m *MockDurationCalculator) Initial() (time.Duration, time.Duration) {
	return 0, time.Second * 5
}

func (m *MockDurationCalculator) Next(lastTimestamp time.Time, retryDuration time.Duration) (time.Duration, time.Time, time.Duration) {
	return millisRate * time.Millisecond, lastTimestamp.Add(-2 * time.Second), time.Second * 5
}

func TestPolling_once(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)
	var wg sync.WaitGroup

	service.opts.PollingInterval = time.Duration(5 * time.Millisecond)
	service.opts.BatteryDetailsPollingInterval = time.Duration(5 * time.Millisecond)
	service.opts.ParameterPollingInterval = time.Duration(5 * time.Millisecond)

	setupPoll(&wg, mockHttpClient, device, mockEndpoint)

	wg.Add(4)

	service.StartPolling(&MockDurationCalculator{})
	time.Sleep(1 * time.Millisecond)
	service.StopPolling()
	time.Sleep(10 * time.Millisecond)

	wg.Wait()

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
	mockEndpoint.AssertNumberOfCalls(t, "PublishDeviceStatus", 1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishBatteryDetails", 1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishPvDetails", 1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishParameterData", 1)
}

func TestPolling_multipleTimes(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)
	var wg sync.WaitGroup

	service.opts.PollingInterval = time.Duration(millisRate * time.Millisecond)
	service.opts.BatteryDetailsPollingInterval = time.Duration(millisRate * time.Millisecond)
	service.opts.ParameterPollingInterval = time.Duration(millisRate * time.Millisecond)

	setupPoll(&wg, mockHttpClient, device, mockEndpoint)

	nLoops := 3
	wg.Add((nLoops + 1) * 4)

	service.StartPolling(&MockDurationCalculator{})
	time.Sleep(time.Duration(nLoops*millisRate+millisRate/2) * time.Millisecond)
	service.StopPolling()

	wg.Wait()

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
	mockEndpoint.AssertNumberOfCalls(t, "PublishDeviceStatus", nLoops+1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishBatteryDetails", nLoops+1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishPvDetails", nLoops+1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishParameterData", nLoops+1)
}

func TestDefaultDurationCalculator_Initial(t *testing.T) {
	calculator := &defaultDurationCalculator{}

	initDur, initRetryDur := calculator.Initial()
	assert.Equal(t, time.Duration(0), initDur)
	assert.Equal(t, time.Second*5, initRetryDur)
}

func TestDefaultDurationCalculator_Next(t *testing.T) {
	calculator := &defaultDurationCalculator{}

	lastTimeStamp := time.Now()
	retryDur := time.Second * 5

	nextDur, nextTimeStamp, nextRetryDur := calculator.Next(lastTimeStamp, retryDur)

	assert.Equal(t, time.Duration(185*time.Second), nextDur.Round(time.Millisecond))
	assert.Equal(t, lastTimeStamp, nextTimeStamp)
	assert.Equal(t, time.Second*5, nextRetryDur)
}

func TestDefaultDurationCalculator_Next2(t *testing.T) {
	calculator := &defaultDurationCalculator{defaultDuration: time.Second * 180}

	lastTimeStamp := time.Now().Add(-200 * time.Second)
	retryDur := time.Second * 5

	nextDur, nextTimeStamp, nextRetryDur := calculator.Next(lastTimeStamp, retryDur)

	assert.Equal(t, retryDur, nextDur)
	assert.Equal(t, lastTimeStamp, nextTimeStamp)
	assert.Equal(t, calculator.defaultDuration, nextRetryDur)
}

func TestDefaultDurationCalculator_Next3(t *testing.T) {
	calculator := &defaultDurationCalculator{defaultDuration: time.Second * 180}

	lastTimeStamp := time.Time{}
	retryDur := time.Second * 5

	nextDur, nextTimeStamp, nextRetryDur := calculator.Next(lastTimeStamp, retryDur)

	assert.Equal(t, calculator.defaultDuration, nextDur)
	assert.Equal(t, lastTimeStamp, nextTimeStamp)
	assert.Equal(t, time.Second*5, nextRetryDur)
}
