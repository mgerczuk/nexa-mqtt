package endpoint

import "nexa-mqtt/pkg/models"

// Used for triggering parameter polling after setting a parameter. Necessary for web+app mode,
// where the app service is responsible for applying parameters, but the web service is responsible
// for polling parameters.
type ParameterQuery interface {
	TriggerParameterPolling(device models.NoahDevicePayload)
}
