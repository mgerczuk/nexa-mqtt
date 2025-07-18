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
	service.SetEndpoint(ep)

	assert.Equal(t, ep, service.endpoint)
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

	dev1 := GrowattNoahList{Datas: []GrowattNoahListData{
		{Sn: "Serial123", PlantID: "1", DeviceModel: "NEXA 2000", Version: "10.10.07.07.4016"},
	}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)
	dev2 := GrowattNoahList{Datas: []GrowattNoahListData{}}
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

	dev1 := GrowattNoahList{Datas: []GrowattNoahListData{
		{Sn: "Serial123", PlantID: "1", DeviceModel: "NEXA 2000", Version: "10.10.07.07.4016"},
	}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)
	dev2 := GrowattNoahList{Datas: []GrowattNoahListData{}}
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

	dev1 := GrowattNoahList{Datas: []GrowattNoahListData{
		{Sn: "Serial123", PlantID: "1", DeviceModel: "NEXA 2000", Version: "10.10.07.07.4016"},
		{Sn: "Serial234", PlantID: "1", DeviceModel: "NEXA 2001"},
	}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)
	dev2 := GrowattNoahList{Datas: []GrowattNoahListData{
		{Sn: "Serial345", PlantID: "2", DeviceModel: "NEXA 2002"},
	}}
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

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{Obj: GrowattNoahStatusObj{
		Pac:                           "-400",
		Ppv:                           "538",
		TotalBatteryPackSoc:           "93",
		TotalBatteryPackChargingPower: "132",
		WorkMode:                      "0",
		Status:                        "6",
	}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{Obj: GrowattNoahTotalsObj{
		EacTotal: "9.6",
		EacToday: "3.3",
	}}, nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			OutputPower:           -400.0,
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

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{Obj: GrowattNoahStatusObj{
		Pac:                           "-400",
		Ppv:                           "538",
		TotalBatteryPackSoc:           "93",
		TotalBatteryPackChargingPower: "-132",
		WorkMode:                      "0",
		Status:                        "6",
	}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{Obj: GrowattNoahTotalsObj{
		EacTotal: "9.6",
		EacToday: "3.3",
	}}, nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			OutputPower:           -400.0,
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

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{Obj: GrowattNoahStatusObj{
		Pac:                           "-400",
		Ppv:                           "538",
		TotalBatteryPackSoc:           "93",
		TotalBatteryPackChargingPower: "132",
		WorkMode:                      "0",
		Status:                        "6",
	}}, nil)
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

	mockHttpClient.OnGetNoahDetails(device.PlantId, device.Serial, GrowattNoahList{Datas: []GrowattNoahListData{
		{
			ChargingSocHighLimit:  "95",
			ChargingSocLowLimit:   "11",
			DefaultMode:           "0",
			DefaultACCouplePower:  "100",
			AllowGridCharging:     "0",
			GridConnectionControl: "1",
			AcCouplePowerControl:  "1",
		},
	}}, nil)

	mockEndpoint.On(
		"PublishParameterData",
		device,
		models.ParameterPayload{
			ChargingLimit:         &chargingLimit,
			DischargeLimit:        &dischargeLimit,
			DefaultACCouplePower:  &defaultACCouplePower,
			DefaultMode:           &defaultMode,
			AllowGridCharging:     allowGridCharging,
			GridConnectionControl: gridConnectionControl,
			AcCouplePowerControl:  acCouplePowerControl,
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
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{Obj: GrowattNoahHistoryObj{Datas: []GrowattNoahHistoryData{
		{
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
		}}}}, nil)

	mockEndpoint.On(
		"PublishBatteryDetails",
		device,
		[]models.BatteryPayload{
			{SerialNumber: "serial124", Soc: 93.0, Temperature: 39.0},
			{SerialNumber: "serial125", Soc: 78.0, Temperature: 41.0},
			{SerialNumber: "serial126", Soc: 82.0, Temperature: 40.0},
			{SerialNumber: "serial127", Soc: 66.0, Temperature: 36.0},
		},
	)

	service.pollBatteryDetails(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_OnGetNoahHistoryFails(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{}, errors.New("OnGetNoahHistory fails"))

	service.pollBatteryDetails(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func Test_pollBatteryDetails_GetNoahHistoryNoData(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)

	today := time.Now().Format("2006-01-02")
	mockHttpClient.OnGetNoahHistory(device.Serial, today, today, GrowattNoahHistory{}, nil)

	service.pollBatteryDetails(device)

	mockHttpClient.AssertExpectations(t)
	mockEndpoint.AssertExpectations(t)
}

func setupPoll(wg *sync.WaitGroup, mockHttpClient *MockHttpClient, device models.NoahDevicePayload, mockEndpoint *MockEndpoint) {
	// ----- enumerateDevices

	plantList := []GrowattPlant{
		{PlantId: "1", PlantName: "plant1"},
	}
	mockHttpClient.OnGetPlantList(plantList, nil)

	dev1 := GrowattNoahList{Datas: []GrowattNoahListData{
		{Sn: device.Serial, PlantID: strconv.Itoa(device.PlantId)},
	}}
	mockHttpClient.OnGetNoahList(1, dev1, nil)

	today := time.Now().Format("2006-01-02")
	// mockHttpClient.OnGetNoahHistory("Serial123", today, today, hist1, nil) set below

	mockEndpoint.On(
		"SetDevices",
		[]models.NoahDevicePayload{device},
	)

	// ----- pollStatus

	mockHttpClient.OnGetNoahStatus(device.PlantId, device.Serial, GrowattNoahStatus{Obj: GrowattNoahStatusObj{
		Pac:                           "-400",
		Ppv:                           "538",
		TotalBatteryPackSoc:           "93",
		TotalBatteryPackChargingPower: "132",
		WorkMode:                      "0",
		Status:                        "6",
	}}, nil)
	mockHttpClient.OnGetNoahTotals(device.PlantId, device.Serial, GrowattNoahTotals{Obj: GrowattNoahTotalsObj{
		EacTotal: "9.6",
		EacToday: "3.3",
	}}, nil)

	mockEndpoint.On(
		"PublishDeviceStatus",
		device,
		models.DevicePayload{
			OutputPower:           -400.0,
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
		}}}}, nil)

	mockEndpoint.On(
		"PublishBatteryDetails",
		device,
		[]models.BatteryPayload{
			{SerialNumber: "serial124", Soc: 93.0, Temperature: 39.0},
			{SerialNumber: "serial125", Soc: 78.0, Temperature: 41.0},
			{SerialNumber: "serial126", Soc: 82.0, Temperature: 40.0},
			{SerialNumber: "serial127", Soc: 66.0, Temperature: 36.0},
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

	mockHttpClient.OnGetNoahDetails(device.PlantId, device.Serial, GrowattNoahList{Datas: []GrowattNoahListData{
		{
			ChargingSocHighLimit:  "95",
			ChargingSocLowLimit:   "11",
			DefaultMode:           "0",
			DefaultACCouplePower:  "100",
			AllowGridCharging:     "1",
			GridConnectionControl: "0",
			AcCouplePowerControl:  "0",
		},
	}}, nil)

	mockEndpoint.On(
		"PublishParameterData",
		device,
		models.ParameterPayload{
			ChargingLimit:         &chargingLimit,
			DischargeLimit:        &dischargeLimit,
			DefaultACCouplePower:  &defaultACCouplePower,
			DefaultMode:           &defaultMode,
			AllowGridCharging:     allowGridCharging,
			GridConnectionControl: gridConnectionControl,
			AcCouplePowerControl:  acCouplePowerControl,
		},
	).Run(func(args mock.Arguments) { wg.Done() })
}

func TestPolling_once(t *testing.T) {
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)
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
	mockHttpClient, service, device, mockEndpoint := setupGrowattServiceMocks(t)
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
