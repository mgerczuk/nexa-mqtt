package polling

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"noah-mqtt/pkg/models"
)

func (s *Service) parametersSubscription(sn string) func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		changedSomething := false
		var payload models.ParameterPayload
		if err := json.Unmarshal(message.Payload(), &payload); err != nil {
			slog.Error("unable to unmarshal parameter command payload", slog.String("error", err.Error()))
		}

		if payload.OutputPower != nil {
			changedSomething = changedSomething || s.setParameterOutputPower(sn, *payload.OutputPower)
		}

		if payload.ChargingLimit != nil || payload.DischargeLimit != nil {
			changedSomething = changedSomething || s.setParameterChargingLimits(sn, payload.ChargingLimit, payload.DischargeLimit)
		}

		if changedSomething {
			s.pollParameterData(sn)
		}
	}
}

func (s *Service) setParameterOutputPower(sn string, outputPower float64) bool {
	slog.Info("trying to set default power", slog.String("device", sn), slog.Int("power", int(outputPower)))
	if err := s.options.GrowattClient.SetDefaultPower(sn, outputPower); err != nil {
		slog.Error("unable to set default power", slog.String("error", err.Error()), slog.String("device", sn))
		return false
	} else {
		slog.Info("set default power", slog.String("device", sn), slog.Int("power", int(outputPower)))
		return true
	}
}

func (s *Service) setParameterChargingLimits(sn string, chargingLimit *float64, dischargeLimit *float64) bool {
	if chargingLimit == nil || dischargeLimit == nil {
		slog.Info("charging limit or discharge limit value not provided, trying to fetch them from server", slog.String("device", sn))
		if data, err := s.options.GrowattClient.GetNoahInfo(sn); err != nil {
			slog.Error("unable to get missing charging limit / discharge limit parameter", slog.String("error", err.Error()))
			return false
		} else {
			if chargingLimit == nil {
				cl := parseFloat(data.Obj.Noah.ChargingSocHighLimit)
				chargingLimit = &cl
			}
			if dischargeLimit == nil {
				dl := parseFloat(data.Obj.Noah.ChargingSocLowLimit)
				dischargeLimit = &dl
			}
		}
	}

	slog.Info("trying to set charging/discharge limit", slog.String("device", sn), slog.Float64("chargingLimit", *chargingLimit), slog.Float64("dischargeLimit", *dischargeLimit))
	if err := s.options.GrowattClient.SetSocLimit(sn, *chargingLimit, *dischargeLimit); err != nil {
		slog.Error("unable to set charging/discharge limit", slog.String("error", err.Error()))
		return false
	} else {
		slog.Info("set charging/discharge limit", slog.String("device", sn), slog.Float64("chargingLimit", *chargingLimit), slog.Float64("dischargeLimit", *dischargeLimit))
		return true
	}
}
