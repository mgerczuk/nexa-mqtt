package misc

import (
	"log/slog"
	"nexa-mqtt/pkg/models"
	"strconv"
)

func S2i(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		return -1
	} else {
		return i
	}
}

func ParseFloat(s string) float64 {
	if s, err := strconv.ParseFloat(s, 64); err == nil {
		return s
	} else {
		return 0
	}
}

func IntStringToOnOff(s string) models.OnOff {
	if s == "0" {
		return models.OFF
	}
	if s == "1" {
		return models.ON
	}

	slog.Error("Invalid 0/1 string", slog.String("s", s))
	return models.OnOff("")
}

func OnOffToInt(s models.OnOff) int {
	if s == models.ON {
		return 1
	}
	if s == models.OFF {
		return 0
	}

	slog.Error("Invalid ON/OFF value", slog.String("s", string(s)))
	return -1
}
