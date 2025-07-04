package growatt_app

import (
	"errors"
	"net/http/cookiejar"

	"nexa-mqtt/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupGrowattAppServiceMock(t *testing.T) (*MockHttpClient, *GrowattAppService, models.NoahDevicePayload, *MockEndpoint) {
	mockHttpClient := MockHttpClient{}
	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)

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
		stop:     make(chan bool),
		endpoint: &endpoint,
	}

	device := models.NoahDevicePayload{
		Serial: "serial123",
	}

	return &mockHttpClient, &service, device, &endpoint
}

func TestServiceLogin_Ok(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", nil)

	err := service.Login()

	assert.Nil(t, err)
	assert.True(t, service.loggedIn)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestServiceLogin_Fail(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	err := service.Login()

	assert.NotNil(t, err)

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

func Test_fetchDevices_GetPlantList_Fail(t *testing.T) {
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

func Test_fetchDevices_NoDevices_Fail(t *testing.T) {
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

	assert.True(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func Test_ensureParameterLogin_Ok(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", nil)

	service.loggedIn = false
	result := service.ensureParameterLogin()

	assert.True(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func Test_ensureParameterLogin_Fail(t *testing.T) {
	mockHttpClient, service, _, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.ensureParameterLogin()

	assert.False(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "system_out_put_power", "1", "350", nil)

	service.loggedIn = true
	result := service.SetOutputPowerW(device, models.WorkMode("battery_first"), 350.0)

	assert.True(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetOutputPowerW(device, models.WorkMode("battery_first"), 350.0)

	assert.False(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_InvalidWorkmode(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	service.loggedIn = true
	result := service.SetOutputPowerW(device, models.WorkMode("invalid_mode"), 350.0)

	assert.False(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetOutputPowerW_SetFailed(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "system_out_put_power", "1", "350", errors.New("SetSystemOutputPower fails"))

	service.loggedIn = true
	result := service.SetOutputPowerW(device, models.WorkMode("battery_first"), 350.0)

	assert.False(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetChargingLimits_Ok(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "charging_soc", "75", "25", nil)

	service.loggedIn = true
	result := service.SetChargingLimits(device, 75.0, 25.0)

	assert.True(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetChargingLimits_LoginFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnLogin("user", "secret", errors.New("Login fails"))

	service.loggedIn = false
	result := service.SetChargingLimits(device, 75.0, 25.0)

	assert.False(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}

func TestSetChargingLimits_SetFail(t *testing.T) {
	mockHttpClient, service, device, endpoint := setupGrowattAppServiceMock(t)

	mockHttpClient.OnSet2Params(device.Serial, "charging_soc", "75", "25", errors.New("SetChargingSoc fails"))

	service.loggedIn = true
	result := service.SetChargingLimits(device, 75.0, 25.0)

	assert.False(t, result)

	mockHttpClient.AssertExpectations(t)
	endpoint.AssertExpectations(t)
}
