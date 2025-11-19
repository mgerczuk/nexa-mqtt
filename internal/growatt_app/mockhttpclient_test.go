package growatt_app

import (
	"fmt"
	"net/url"

	"github.com/stretchr/testify/mock"
)

// ----- Mock ---------------------------------------------------------------

// MockHttpClient implements HttpClient
type MockHttpClient struct {
	mock.Mock
}

func (h *MockHttpClient) postForm(url string, token string, data url.Values, responseBody any) error {
	args := h.Called(url, token, data, responseBody)
	return args.Error(0)
}

// ----- Helpers ------------------------------------------------------------

func (m *MockHttpClient) On_loginGetToken(userid string, password string, result TokenResponse, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://evcharge.growatt.com/ocpp/user",
		"",
		url.Values{
			"cmd":      {"shineLogin"},
			"userId":   {userid},
			"password": {password},
			"lan":      {"1"},
		},
		&TokenResponse{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(3).(*TokenResponse)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) On_newTwoLoginAPIV2(token string, username string, password string, result LoginResult, err error) *mock.Call {
	expectedMatch := mock.MatchedBy(func(m url.Values) bool {
		return m.Get("userName") == username &&
			m.Get("password") == password &&
			m.Get("newLogin") == "1" &&
			m.Get("appType") == "ShinePhone"
	})

	call := m.On(
		"postForm",
		"https://server-api.growatt.com/newTwoLoginAPIV2.do",
		token,
		expectedMatch,
		&LoginResult{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(3).(*LoginResult)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnLogin(username string, password string, err error) {
	tokenResponse := TokenResponse{}
	m.On_loginGetToken(fmt.Sprintf("SHINE%s", username), password, tokenResponse, nil)

	loginResult := LoginResult{}
	loginResult.Back.Success = true
	m.On_newTwoLoginAPIV2(tokenResponse.Token, username, password, loginResult, err)
}

func (m *MockHttpClient) OnGetPlantList(result PlantListV2, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/newTwoPlantAPI.do?op=getAllPlantListTwo",
		"",
		url.Values{
			"plantStatus": {""},
			"pageSize":    {"20"},
			"language":    {"1"},
			"toPageNum":   {"1"},
			"order":       {"1"},
		},
		&PlantListV2{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(3).(*PlantListV2)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahPlantInfo(plantId string, result NoahPlantInfoObj, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/noah/isPlantNoahSystem",
		"",
		url.Values{
			"plantId": {plantId},
		},
		&NoahPlantInfo{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(3).(*NoahPlantInfo)
				*responseBody = NoahPlantInfo{
					ResponseContainerV2[NoahPlantInfoObj]{
						Obj: result,
					},
				}
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahStatus(serialNumber string, result NoahStatusObj, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getSystemStatus",
		"",
		url.Values{
			"deviceSn": {serialNumber},
		},
		&NoahStatus{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(3).(*NoahStatus)
				*responseBody = NoahStatus{
					ResponseContainerV2[NoahStatusObj]{
						Obj: result,
					},
				}
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahInfo(serialNumber string, result NexaInfoObj, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getNexaInfoBySn",
		"",
		url.Values{
			"language": {"1"},
			"deviceSn": {serialNumber},
		},
		&NexaInfo{},
	)

	if err == nil {
		call = call.Run(func(args mock.Arguments) {
			responseBody := args.Get(3).(*NexaInfo)
			*responseBody = NexaInfo{
				ResponseContainerV2[NexaInfoObj]{
					Obj: result,
				},
			}
		})
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetBatteryData(serialNumber string, result BatteryInfoObj, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/getBatteryData",
		"",
		url.Values{
			"deviceSn": {serialNumber},
		},
		&BatteryInfo{},
	)

	if err == nil {
		call = call.Run(func(args mock.Arguments) {
			responseBody := args.Get(3).(*BatteryInfo)
			*responseBody = BatteryInfo{
				ResponseContainerV2[BatteryInfoObj]{
					Obj: result,
				},
			}
		})
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnSet1Param(serialNumber string, typ string, param1 string, err error, response SetResponse) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {serialNumber},
			"type":      {typ},
			"param1":    {param1},
		},
		&SetResponse{},
	)

	if err == nil {
		call = call.Run(func(args mock.Arguments) {
			responseBody := args.Get(3).(*SetResponse)
			*responseBody = response
		})
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnSet2Params(serialNumber string, typ string, param1 string, param2 string, err error, response SetResponse) *mock.Call {
	call := m.On(
		"postForm",
		"https://server-api.growatt.com/noahDeviceApi/nexa/set",
		"",
		url.Values{
			"serialNum": {serialNumber},
			"type":      {typ},
			"param1":    {param1},
			"param2":    {param2},
		},
		&SetResponse{},
	)

	if err == nil {
		call = call.Run(func(args mock.Arguments) {
			responseBody := args.Get(3).(*SetResponse)
			*responseBody = response
		})
	}

	return call.Return(err)
}
