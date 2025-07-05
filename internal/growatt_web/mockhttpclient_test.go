package growatt_web

import (
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

func (m *MockHttpClient) OnLogin(account string, password string, err error) *mock.Call {
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
				*responseBody = GrowattResult{}
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
