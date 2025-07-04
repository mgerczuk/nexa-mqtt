package growatt_app

import (
	"nexa-mqtt/internal/endpoint"
	"nexa-mqtt/pkg/models"

	"github.com/stretchr/testify/mock"
)

// MockEndpoint implements endpoint.Endpoint
type MockEndpoint struct {
	mock.Mock
}

func (e *MockEndpoint) SetParameterApplier(applier endpoint.ParameterApplier) {
	e.Called(applier)
}

func (e *MockEndpoint) SetDevices(devices []models.NoahDevicePayload) {
	e.Called(devices)
}

func (e *MockEndpoint) PublishDeviceStatus(device models.NoahDevicePayload, status models.DevicePayload) {
	e.Called(device, status)
}

func (e *MockEndpoint) PublishBatteryDetails(device models.NoahDevicePayload, details []models.BatteryPayload) {
	e.Called(device, details)
}

func (e *MockEndpoint) PublishParameterData(device models.NoahDevicePayload, param models.ParameterPayload) {
	e.Called(device, param)
}
