package growatt_web

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

func (h *MockHttpClient) postForm(url string, data url.Values, responseBody any) error {
	args := h.Called(url, data, responseBody)
	return args.Error(0)
}

// ----- Helpers ------------------------------------------------------------

func (m *MockHttpClient) OnLogin(account string, password string, result GrowattResult, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://openapi.growatt.com/login",
		url.Values{
			"account":  {account},
			"password": {password},
		},
		&GrowattResult{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattResult)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetPlantList(result []GrowattPlant, err error) *mock.Call {
	var expected []GrowattPlant
	call := m.On(
		"postForm",
		"https://openapi.growatt.com/index/getPlantListTitle",
		url.Values{},
		&expected,
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*[]GrowattPlant)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetPlantDevices(plantId string, result GrowattPlantDevices, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://openapi.growatt.com/panel/getDevicesByPlantList",
		url.Values{
			"plantId":  {plantId},
			"currPage": {"1"},
		},
		&GrowattPlantDevices{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattPlantDevices)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahList(plantId int, result GrowattNoahList, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://openapi.growatt.com/device/getNoahList",
		url.Values{
			"plantId":  {fmt.Sprintf("%d", plantId)},
			"currPage": {"1"},
		},
		&GrowattNoahList{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattNoahList)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahDetails(plantId int, serial string, result GrowattNoahList, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://openapi.growatt.com/device/getNoahList",
		url.Values{
			"plantId":  {fmt.Sprintf("%d", plantId)},
			"deviceSn": {serial},
			"currPage": {"1"},
		},
		&GrowattNoahList{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattNoahList)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahHistory(serial string, startDate string, endDate string, result GrowattNoahHistory, err error) *mock.Call {
	call := m.On(
		"postForm",
		"https://openapi.growatt.com/device/getNoahHistory",
		url.Values{
			"deviceSn":  {serial},
			"start":     {"0"},
			"startDate": {startDate},
			"endDate":   {endDate},
		},
		&GrowattNoahHistory{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattNoahHistory)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahStatus(plantId int, serial string, result GrowattNoahStatus, err error) *mock.Call {
	call := m.On(
		"postForm",
		fmt.Sprintf("https://openapi.growatt.com/panel/noah/getNoahStatusData?plantId=%d", plantId),
		url.Values{
			"deviceSn": {serial},
		},
		&GrowattNoahStatus{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattNoahStatus)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}

func (m *MockHttpClient) OnGetNoahTotals(plantId int, serial string, result GrowattNoahTotals, err error) *mock.Call {
	call := m.On(
		"postForm",
		fmt.Sprintf("https://openapi.growatt.com/panel/noah/getNoahTotalData?plantId=%d", plantId),
		url.Values{
			"deviceSn": {serial},
		},
		&GrowattNoahTotals{},
	)

	if err == nil {
		call = call.Run(
			func(args mock.Arguments) {
				responseBody := args.Get(2).(*GrowattNoahTotals)
				*responseBody = result
			},
		)
	}

	return call.Return(err)
}
