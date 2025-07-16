package growatt_app

import (
	"errors"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----- Test functions -----------------------------------------------------

func setupMocks(t *testing.T) (*MockHttpClient, *Client) {
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

	return &mockHttpClient, &client
}

func Test_postForm_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(nil)

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_AnyError(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(errors.New("some error"))

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_RetryLogin(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(errors.New("invalid character '<' looking for beginning of value"))

	mockHttpClient.On_loginGetToken("SHINEuser", "secret", TokenResponse{Token: "THE_TOKEN"}, nil)

	loginResult := LoginResult{}
	loginResult.Back.Success = true
	mockHttpClient.On_newTwoLoginAPIV2("THE_TOKEN", "user", "secret", loginResult, nil)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"THE_TOKEN",
		url.Values{},
		nil,
	).Return(nil)

	err := client.postForm("http://someurl", url.Values{}, nil)
	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func Test_postForm_RetryLoginFails(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On(
		"postForm",
		"http://someurl",
		"",
		url.Values{},
		nil,
	).Return(errors.New("invalid character '<' looking for beginning of value"))

	mockHttpClient.On_loginGetToken("SHINEuser", "secret", TokenResponse{Token: "THE_TOKEN"}, nil)

	mockHttpClient.On_newTwoLoginAPIV2("THE_TOKEN", "user", "secret", LoginResult{}, errors.New("login failed"))

	defer func() {
		if r := recover(); r != nil {
			mockHttpClient.AssertExpectations(t)
		}
	}()

	client.postForm("http://someurl", url.Values{}, nil)
	t.Errorf("Test failed, panic was expected")
}

func TestLogin_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On_loginGetToken("SHINEuser", "secret", TokenResponse{Token: "THE_TOKEN"}, nil)

	loginResult := LoginResult{}
	loginResult.Back.Success = true
	mockHttpClient.On_newTwoLoginAPIV2("THE_TOKEN", "user", "secret", loginResult, nil)

	err := client.Login()
	assert.NoError(t, err)
	assert.Equal(t, "THE_TOKEN", client.token)

	mockHttpClient.AssertExpectations(t)
}

func TestLogin_NoSuccess(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On_loginGetToken("SHINEuser", "secret", TokenResponse{Token: "THE_TOKEN"}, nil)

	loginResult := LoginResult{}
	loginResult.Back.Success = false
	mockHttpClient.On_newTwoLoginAPIV2("THE_TOKEN", "user", "secret", loginResult, nil)

	err := client.Login()
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestLogin_loginGetToken_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On_loginGetToken("SHINEuser", "secret", TokenResponse{}, errors.New("login error"))

	err := client.Login()
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestLogin_newTwoLoginFail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.On_loginGetToken("SHINEuser", "secret", TokenResponse{Token: "THE_TOKEN"}, nil)

	mockHttpClient.On_newTwoLoginAPIV2("THE_TOKEN", "user", "secret", LoginResult{}, errors.New("newTwoLoginAPIV2 fail"))

	err := client.Login()
	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestGetPlantList_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	expectedPlantList := PlantListV2{
		PlantList: []struct {
			ID int `json:"id"`
		}{
			{ID: 1},
			{ID: 2},
		},
	}
	mockHttpClient.OnGetPlantList(expectedPlantList, nil)

	data, err := client.GetPlantList()

	assert.NoError(t, err)
	assert.Equal(t, expectedPlantList, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetPlantList_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetPlantList(PlantListV2{},
		errors.New("newTwoPlantAPI fail"))

	data, err := client.GetPlantList()
	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahPlantInfo_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahPlantInfo("2", NoahPlantInfoObj{IsPlantHaveNexa: true}, nil)

	data, err := client.GetNoahPlantInfo("2")

	assert.NoError(t, err)
	assert.Equal(t, NoahPlantInfo{
		ResponseContainerV2: ResponseContainerV2[NoahPlantInfoObj]{
			Msg:    "",
			Result: 0,
			Obj: NoahPlantInfoObj{
				IsPlantHaveNexa: true,
			},
		},
	}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahPlantInfo_NoNexa(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahPlantInfo("2", NoahPlantInfoObj{IsPlantHaveNexa: false}, nil)

	data, err := client.GetNoahPlantInfo("2")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahPlantInfo_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahPlantInfo("2", NoahPlantInfoObj{}, errors.New("isPlantNoahSystem fail"))

	data, err := client.GetNoahPlantInfo("2")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahStatus_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahStatus("serial123", NoahStatusObj{}, nil)

	data, err := client.GetNoahStatus("serial123")

	assert.NoError(t, err)
	assert.Equal(t, NoahStatus{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahStatus_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahStatus("serial123", NoahStatusObj{}, errors.New("getSystemStatus fail"))

	data, err := client.GetNoahStatus("serial123")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahInfo_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahInfo("serial123", NexaInfoObj{}, nil)

	data, err := client.GetNoahInfo("serial123")

	assert.NoError(t, err)
	assert.Equal(t, NexaInfo{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetNoahInfo_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetNoahInfo("serial123", NexaInfoObj{}, errors.New("getSystemStatus fail"))

	data, err := client.GetNoahInfo("serial123")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetBatteryData_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetBatteryData("serial123", BatteryInfoObj{}, nil)

	data, err := client.GetBatteryData("serial123")

	assert.NoError(t, err)
	assert.Equal(t, BatteryInfo{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestGetBatteryData_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnGetBatteryData("serial123", BatteryInfoObj{}, errors.New("getBatteryData fail"))

	data, err := client.GetBatteryData("serial123")

	assert.Error(t, err)
	assert.Nil(t, data)

	mockHttpClient.AssertExpectations(t)
}

func TestSetSystemOutputPower_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet2Params("serial123", "system_out_put_power", "0", "200", nil)

	err := client.SetSystemOutputPower("serial123", 0, 200)

	assert.NoError(t, err)
	//	assert.Equal(t, BatteryInfo{}, *data)

	mockHttpClient.AssertExpectations(t)
}

func TestSetSystemOutputPower_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet2Params("serial123", "system_out_put_power", "0", "200", errors.New("noahDeviceApi/nexa/set system_out_put_power fail"))

	err := client.SetSystemOutputPower("serial123", 0, 200)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingSoc_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet2Params("serial123", "charging_soc", "85", "15", nil)

	err := client.SetChargingSoc("serial123", 85, 15)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetChargingSoc_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet2Params("serial123", "charging_soc", "85", "15", errors.New("noahDeviceApi/nexa/set charging_soc fail"))

	err := client.SetChargingSoc("serial123", 85, 15)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet1Param("serial123", "allow_grid_charging", "1", nil)

	err := client.SetAllowGridCharging("serial123", 1)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetAllowGridCharging_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet1Param("serial123", "allow_grid_charging", "1", errors.New("noahDeviceApi/nexa/set allow_grid_charging fail"))

	err := client.SetAllowGridCharging("serial123", 1)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet1Param("serial123", "grid_connection_control", "1", nil)

	err := client.SetGridConnectionControl("serial123", 1)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetGridConnectionControl_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet1Param("serial123", "grid_connection_control", "1", errors.New("noahDeviceApi/nexa/set grid_connection_control fail"))

	err := client.SetGridConnectionControl("serial123", 1)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetACCouplePowerControl_Ok(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet1Param("serial123", "ac_couple_power_control", "1", nil)

	err := client.SetACCouplePowerControl("serial123", 1)

	assert.NoError(t, err)

	mockHttpClient.AssertExpectations(t)
}

func TestSetACCouplePowerControl_Fail(t *testing.T) {
	mockHttpClient, client := setupMocks(t)

	mockHttpClient.OnSet1Param("serial123", "ac_couple_power_control", "1", errors.New("noahDeviceApi/nexa/set ac_couple_power_control fail"))

	err := client.SetACCouplePowerControl("serial123", 1)

	assert.Error(t, err)

	mockHttpClient.AssertExpectations(t)
}
