package endpoint_mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
)

// MockToken implements mqtt.Token
type MockToken struct {
	mock.Mock
	done chan struct{}
}

func NewMockToken() *MockToken {
	done := make(chan struct{})
	close(done) // sofort abgeschlossen
	return &MockToken{done: done}
}

func (m *MockToken) Wait() bool                     { return true }
func (m *MockToken) WaitTimeout(time.Duration) bool { return true }
func (t *MockToken) Done() <-chan struct{}          { return t.done }
func (t *MockToken) Error() error {
	args := t.Called("Error")
	return args.Error(0)
}

// MockMqttClient implements mqtt.Client
type MockMqttClient struct {
	mock.Mock
	mqtt.Client
}

func (m *MockMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	args := m.Called(topic, qos, retained, payload)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	args := m.Called(topic, qos, callback)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	ifaceArgs := make([]interface{}, len(topics))
	for i, v := range topics {
		ifaceArgs[i] = v
	}
	args := m.Called(ifaceArgs...)
	return args.Get(0).(mqtt.Token)
}
