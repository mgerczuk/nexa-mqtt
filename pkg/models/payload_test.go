package models

import "testing"

func TestParameterPayloadUpdateFrom(t *testing.T) {
	chargingLimit := 12.5
	dischargeLimit := 7.5
	defaultACCouplePower := 450.0
	var defaultMode WorkMode = WorkModeBatteryFirst
	antiBackflowPowerPercentage := 80.0

	from := ParameterPayload{
		ChargingLimit:               &chargingLimit,
		DischargeLimit:              &dischargeLimit,
		DefaultACCouplePower:        &defaultACCouplePower,
		DefaultMode:                 &defaultMode,
		AllowGridCharging:           ON,
		GridConnectionControl:       ON,
		AcCouplePowerControl:        ON,
		AcCoupleEnable:              ON,
		LightLoadEnable:             ON,
		NeverPowerOff:               ON,
		AntiBackflowEnable:          ON,
		AntiBackflowPowerPercentage: &antiBackflowPowerPercentage,
	}

	to := EmptyParameterPayload()
	to.UpdateFrom(from)

	if to.ChargingLimit == nil || *to.ChargingLimit != chargingLimit {
		t.Fatalf("ChargingLimit = %v, want %v", to.ChargingLimit, chargingLimit)
	}
	if to.DischargeLimit == nil || *to.DischargeLimit != dischargeLimit {
		t.Fatalf("DischargeLimit = %v, want %v", to.DischargeLimit, dischargeLimit)
	}
	if to.DefaultACCouplePower == nil || *to.DefaultACCouplePower != defaultACCouplePower {
		t.Fatalf("DefaultACCouplePower = %v, want %v", to.DefaultACCouplePower, defaultACCouplePower)
	}
	if to.DefaultMode == nil || *to.DefaultMode != defaultMode {
		t.Fatalf("DefaultMode = %v, want %v", to.DefaultMode, defaultMode)
	}
	if to.AllowGridCharging != ON {
		t.Fatalf("AllowGridCharging = %q, want %q", to.AllowGridCharging, ON)
	}
	if to.GridConnectionControl != ON {
		t.Fatalf("GridConnectionControl = %q, want %q", to.GridConnectionControl, ON)
	}
	if to.AcCouplePowerControl != ON {
		t.Fatalf("AcCouplePowerControl = %q, want %q", to.AcCouplePowerControl, ON)
	}
	if to.AcCoupleEnable != ON {
		t.Fatalf("AcCoupleEnable = %q, want %q", to.AcCoupleEnable, ON)
	}
	if to.LightLoadEnable != ON {
		t.Fatalf("LightLoadEnable = %q, want %q", to.LightLoadEnable, ON)
	}
	if to.NeverPowerOff != ON {
		t.Fatalf("NeverPowerOff = %q, want %q", to.NeverPowerOff, ON)
	}
	if to.AntiBackflowEnable != ON {
		t.Fatalf("AntiBackflowEnable = %q, want %q", to.AntiBackflowEnable, ON)
	}
	if to.AntiBackflowPowerPercentage == nil || *to.AntiBackflowPowerPercentage != antiBackflowPowerPercentage {
		t.Fatalf("AntiBackflowPowerPercentage = %v, want %v", to.AntiBackflowPowerPercentage, antiBackflowPowerPercentage)
	}
}

func TestParameterPayloadUpdateFrom_LeavesZeroValuesUntouched(t *testing.T) {
	to := ParameterPayload{
		ChargingLimit:               nil,
		DischargeLimit:              nil,
		DefaultACCouplePower:        nil,
		DefaultMode:                 nil,
		AllowGridCharging:           OFF,
		GridConnectionControl:       OFF,
		AcCouplePowerControl:        OFF,
		AcCoupleEnable:              OFF,
		LightLoadEnable:             OFF,
		NeverPowerOff:               OFF,
		AntiBackflowEnable:          OFF,
		AntiBackflowPowerPercentage: nil,
	}

	from := ParameterPayload{}
	to.UpdateFrom(from)

	if to.ChargingLimit != nil {
		t.Fatalf("ChargingLimit should remain nil")
	}
	if to.DischargeLimit != nil {
		t.Fatalf("DischargeLimit should remain nil")
	}
	if to.DefaultACCouplePower != nil {
		t.Fatalf("DefaultACCouplePower should remain nil")
	}
	if to.DefaultMode != nil {
		t.Fatalf("DefaultMode should remain nil")
	}
	if to.AllowGridCharging != OFF {
		t.Fatalf("AllowGridCharging = %q, want %q", to.AllowGridCharging, OFF)
	}
	if to.GridConnectionControl != OFF {
		t.Fatalf("GridConnectionControl = %q, want %q", to.GridConnectionControl, OFF)
	}
	if to.AcCouplePowerControl != OFF {
		t.Fatalf("AcCouplePowerControl = %q, want %q", to.AcCouplePowerControl, OFF)
	}
	if to.AcCoupleEnable != OFF {
		t.Fatalf("AcCoupleEnable = %q, want %q", to.AcCoupleEnable, OFF)
	}
	if to.LightLoadEnable != OFF {
		t.Fatalf("LightLoadEnable = %q, want %q", to.LightLoadEnable, OFF)
	}
	if to.NeverPowerOff != OFF {
		t.Fatalf("NeverPowerOff = %q, want %q", to.NeverPowerOff, OFF)
	}
	if to.AntiBackflowEnable != OFF {
		t.Fatalf("AntiBackflowEnable = %q, want %q", to.AntiBackflowEnable, OFF)
	}
	if to.AntiBackflowPowerPercentage != nil {
		t.Fatalf("AntiBackflowPowerPercentage should remain nil")
	}
}
