package growatt_app

import (
	"errors"
	"math/rand"
	"net/http/cookiejar"
	"strconv"
	"sync"
	"time"

	"nexa-mqtt/pkg/models"
	"testing"

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

func setupGrowattAppServiceMock(t *testing.T) (*MockHttpClient, *GrowattAppService, models.NoahDevicePayload, *MockEndpoint) {
	mockHttpClient := MockHttpClient{}
	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)

	client := Client{
		client:    &mockHttpClient,
		serverUrl: "https://server-api.growatt.com",
		username:  "user",
		password:  "secret",
		jar:       jar,
	}

	endpoint := MockEndpoint{}

	service := GrowattAppService{
		client:   &client,
		endpoint: &endpoint,
	}

	device := models.NoahDevicePayload{
		Serial:  randSerial(10),
		PlantId: rand.Intn(20000000),
	}

	return &mockHttpClient, &service, device, &endpoint
}

func TestServiceLogin_Ok(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", nil)

	err := service.Login()

	assert.NoError(t, err)
	assert.True(t, service.loggedIn)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestServiceLogin_Fails(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	err := service.Login()

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func Test_fetchDevices_Ok(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetPlantList(PlantListV2{
		PlantList: []struct {
			ID int `json:"id"`
		}{
			{ID: 1},
			{ID: 2},
			{ID: 3},
			{ID: 4},
		},
	}, nil)

	mockHttpClient.OnGetNoahPlantInfo("1", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: "serial234"}, nil)
	mockHttpClient.OnGetNoahPlantInfo("2", NoahPlantInfoObj{IsPlantHaveNexa: false}, nil)
	mockHttpClient.OnGetNoahPlantInfo("3", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: ""}, nil)
	mockHttpClient.OnGetNoahPlantInfo("4", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: "serial235"}, nil)

	devices := service.fetchDevices()

	assert.Equal(t, 2, len(devices))
	assert.Equal(
		t,
		[]models.NoahDevicePayload{
			{PlantId: 1, Serial: "serial234"},
			{PlantId: 4, Serial: "serial235"},
		},
		devices)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func Test_fetchDevices_GetPlantList_Fails(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetPlantList(PlantListV2{}, errors.New("GetPlantList fails"))

	defer func() {
		if r := recover(); r != nil {
			mockHttpClient.AssertExpectations(t)
			endpoint.AssertExpectations(t)
		}
	}()

	service.fetchDevices()
	t.Errorf("Test failed, panic was expected")
}

func Test_fetchDevices_NoDevices_Fails(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetPlantList(PlantListV2{
		PlantList: []struct {
			ID int `json:"id"`
		}{
			{ID: 1},
			{ID: 2},
		},
	}, nil)

	mockHttpClient.OnGetNoahPlantInfo("1", NoahPlantInfoObj{IsPlantHaveNexa: false}, nil)
	mockHttpClient.OnGetNoahPlantInfo("2", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: ""}, nil)

	defer func() {
		if r := recover(); r != nil {
			mockHttpClient.AssertExpectations(t)
			endpoint.AssertExpectations(t)
		}
	}()

	service.fetchDevices()
	t.Errorf("Test failed, panic was expected")
}

func Test_enumerateDevices_Ok(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnGetPlantList(PlantListV2{
		PlantList: []struct {
			ID int `json:"id"`
		}{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		},
	}, nil)

	mockHttpClient.OnGetNoahPlantInfo("1", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: "serial234"}, nil)
	mockHttpClient.OnGetNoahPlantInfo("2", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: "serial235"}, nil)
	mockHttpClient.OnGetNoahPlantInfo("3", NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: "serial236"}, nil)

	nexaInfoObj1 := NexaInfoObj{}
	nexaInfoObj1.Noah.Model = "NEXA 2000"
	nexaInfoObj1.Noah.Alias = "NEXA 2000"
	nexaInfoObj1.Noah.Version = "09.05.05.04.9000.4014"
	nexaInfoObj1.Noah.BatSns = []string{"0XXX00XX00XX0000"}
	mockHttpClient.OnGetNoahInfo("serial234", nexaInfoObj1, nil)

	nexaInfoObj2 := NexaInfoObj{}
	nexaInfoObj2.Noah.Model = "NEXA 2001"
	nexaInfoObj2.Noah.Alias = "NEXA 2001"
	nexaInfoObj2.Noah.Version = "09.05.05.04.9000.4033"
	nexaInfoObj2.Noah.BatSns = []string{"0XXX00XX00XX0001", "0XXX00XX00XX0002"}
	mockHttpClient.OnGetNoahInfo("serial235", nexaInfoObj2, nil)

	mockHttpClient.OnGetNoahInfo("serial236", NexaInfoObj{}, errors.New("One GetNoahInfo fails"))

	expectedDevices := []models.NoahDevicePayload{
		{
			PlantId: 1,
			Serial:  "serial234",
			Model:   "NEXA 2000",
			Alias:   "NEXA 2000",
			Version: "09.05.05.04.9000.4014",
			Batteries: []models.NoahDeviceBatteryPayload{
				{Alias: "BAT0"},
			},
		},
		{
			PlantId: 2,
			Serial:  "serial235",
			Model:   "NEXA 2001",
			Alias:   "NEXA 2001",
			Version: "09.05.05.04.9000.4033",
			Batteries: []models.NoahDeviceBatteryPayload{
				{Alias: "BAT0"},
				{Alias: "BAT1"},
			},
		},
		{
			PlantId: 3,
			Serial:  "serial236",
		},
	}

	endpoint.On(
		"SetDevices",
		expectedDevices,
	)

	service.enumerateDevices()

	assert.Equal(t, 3, len(service.devices))
	assert.Equal(
		t,
		expectedDevices,
		service.devices)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetEndpoint(t *testing.T) {
	_, service, _, _ := setupGrowattAppServiceMock(t)

	service.SetEndpoint(nil)

	assert.Equal(t, nil, service.endpoint)

	ep := &MockEndpoint{}
	service.SetEndpoint(ep)

	assert.Equal(t, ep, service.endpoint)
}

func Test_ensureParameterLogin_AlreadyLoggedIn(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	service.loggedIn = true
	result := service.ensureParameterLogin()

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func Test_ensureParameterLogin_Ok(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", nil)

	service.loggedIn = false
	result := service.ensureParameterLogin()

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func Test_ensureParameterLogin_Fails(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.ensureParameterLogin()

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "system_out_put_power", "1", "350", nil, SetResponse{ResponseContainerV2[any]{Result: 1}})

	service.loggedIn = true
	result := service.SetOutputPowerW(device, models.WorkMode("battery_first"), 350.0)

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetOutputPowerW(device, models.WorkMode("battery_first"), 350.0)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_InvalidWorkmode(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	service.loggedIn = true
	result := service.SetOutputPowerW(device, models.WorkMode("invalid_mode"), 350.0)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_SetFailed(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "system_out_put_power", "1", "350", errors.New("SetSystemOutputPower fails"), SetResponse{})

	service.loggedIn = true
	result := service.SetOutputPowerW(device, models.WorkMode("battery_first"), 350.0)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetChargingLimits_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "charging_soc", "75", "25", nil, SetResponse{ResponseContainerV2[any]{Result: 1}})

	service.loggedIn = true
	result := service.SetChargingLimits(device, 75.0, 25.0)

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetChargingLimits_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetChargingLimits(device, 75.0, 25.0)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetChargingLimits_SetFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "charging_soc", "75", "25", errors.New("SetChargingSoc fails"), SetResponse{})

	service.loggedIn = true
	result := service.SetChargingLimits(device, 75.0, 25.0)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetAllowGridChargingService_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet1Param(device.Serial, "allow_grid_charging", "1", nil, SetResponse{ResponseContainerV2[any]{Result: 1}})

	service.loggedIn = true
	result := service.SetAllowGridCharging(device, models.ON)

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetAllowGridChargingService_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetAllowGridCharging(device, models.ON)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetAllowGridChargingService_Fails(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet1Param(device.Serial, "allow_grid_charging", "0", errors.New("SetAllowGridCharging fails"), SetResponse{ResponseContainerV2[any]{Result: 1}})

	service.loggedIn = true
	result := service.SetAllowGridCharging(device, models.OFF)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetGridConnectionControlService_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet1Param(device.Serial, "grid_connection_control", "0", nil, SetResponse{ResponseContainerV2[any]{Result: 1}})

	service.loggedIn = true
	result := service.SetGridConnectionControl(device, models.OFF)

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetGridConnectionControl_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetGridConnectionControl(device, models.ON)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetGridConnectionControl_Fails(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet1Param(device.Serial, "grid_connection_control", "1", errors.New("SetAllowGridCharging fails"), SetResponse{})

	service.loggedIn = true
	result := service.SetGridConnectionControl(device, models.ON)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetAcCouplePowerControl_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet1Param(device.Serial, "ac_couple_power_control", "0", nil, SetResponse{ResponseContainerV2[any]{Result: 1}})

	service.loggedIn = true
	result := service.SetAcCouplePowerControl(device, models.OFF)

	assert.NoError(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetAcCouplePowerControl_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetAcCouplePowerControl(device, models.ON)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetAcCouplePowerControl_Fails(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet1Param(device.Serial, "ac_couple_power_control", "1", errors.New("SetAllowGridCharging fails"), SetResponse{})

	service.loggedIn = true
	result := service.SetAcCouplePowerControl(device, models.ON)

	assert.Error(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func setupPoll(wg *sync.WaitGroup, mockHttpClient *MockHttpClient, device models.NoahDevicePayload, mockEndpoint *MockEndpoint) {
	nexaInfo := NexaInfoObj{}

	nexaInfo.Noah.ChargingSocHighLimit = "95"
	nexaInfo.Noah.DefaultMode = "0"
	nexaInfo.Noah.DefaultACCouplePower = "100"
	nexaInfo.Noah.ChargingSocLowLimit = "11"

	chargingLimit := 95.0
	dischargeLimit := 11.0
	defaultACCouplePower := 100.0
	defaultMode := models.WorkMode("load_first")

	// ----- enumerateDevices

	mockHttpClient.OnGetPlantList(PlantListV2{
		PlantList: []struct {
			ID int `json:"id"`
		}{
			{ID: device.PlantId},
		},
	}, nil)
	mockHttpClient.OnGetNoahPlantInfo(strconv.Itoa(device.PlantId), NoahPlantInfoObj{IsPlantHaveNexa: true, DeviceSn: device.Serial}, nil)
	mockHttpClient.OnGetNoahInfo(device.Serial, nexaInfo, nil)

	expectedDevices := []models.NoahDevicePayload{
		{
			PlantId: device.PlantId,
			Serial:  device.Serial,
		},
	}
	mockEndpoint.On(
		"SetDevices",
		expectedDevices,
	)

	// ----- pollStatus

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
			ACPower:               -400.0,
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
	).Run(func(args mock.Arguments) { wg.Done() })

	// ----- pollBatteryDetails

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
	).Run(func(args mock.Arguments) { wg.Done() })

	// ----- pollParameterData

	// mockHttpClient.OnGetNoahInfo already set

	mockEndpoint.On(
		"PublishParameterData",
		device,
		models.ParameterPayload{
			ChargingLimit:        &chargingLimit,
			DischargeLimit:       &dischargeLimit,
			DefaultACCouplePower: &defaultACCouplePower,
			DefaultMode:          &defaultMode,
		},
	).Run(func(args mock.Arguments) { wg.Done() })
}

func TestPolling_once(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)
	var wg sync.WaitGroup

	service.opts.PollingInterval = time.Duration(5 * time.Millisecond)
	service.opts.BatteryDetailsPollingInterval = time.Duration(5 * time.Millisecond)
	service.opts.ParameterPollingInterval = time.Duration(5 * time.Millisecond)

	setupPoll(&wg, mockHttpClient, device, mockEndpoint)

	wg.Add(3)

	service.StartPolling()
	service.StopPolling()
	time.Sleep(10 * time.Millisecond)

	wg.Wait()

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
	mockEndpoint.AssertNumberOfCalls(t, "PublishDeviceStatus", 1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishBatteryDetails", 1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishParameterData", 1)
}

func TestPolling_multipleTimes(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattAppServiceMock(t)
	var wg sync.WaitGroup

	service.opts.PollingInterval = time.Duration(5 * time.Millisecond)
	service.opts.BatteryDetailsPollingInterval = time.Duration(5 * time.Millisecond)
	service.opts.ParameterPollingInterval = time.Duration(5 * time.Millisecond)

	setupPoll(&wg, mockHttpClient, device, mockEndpoint)

	nLoops := 2
	wg.Add((nLoops + 1) * 3)

	service.StartPolling()
	time.Sleep(time.Duration(nLoops*5+3) * time.Millisecond)
	service.StopPolling()

	wg.Wait()

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
	mockEndpoint.AssertNumberOfCalls(t, "PublishDeviceStatus", nLoops+1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishBatteryDetails", nLoops+1)
	mockEndpoint.AssertNumberOfCalls(t, "PublishParameterData", nLoops+1)
}
