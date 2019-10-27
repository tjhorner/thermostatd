package thermostat

import (
	"errors"
	"fmt"
)

// NewState returns a State with sane default values
func NewState() State {
	return State{
		PoweredOn:          false,
		CurrentMode:        ModeCool,
		FanSpeed:           FanSpeedAuto,
		TargetTemperature:  72,
		CurrentTemperature: 0,
	}
}

// State is a thermostat state
type State struct {
	// PoweredOn is true if the thermostat is on
	PoweredOn bool `json:"powered_on"`
	// CurrentMode is the current mode the thermostat is in
	CurrentMode Mode `json:"current_mode"`
	// FanSpeed is the current fan speed
	FanSpeed FanSpeed `json:"fan_speed"`
	// TargetTemperature is the target temperature in Fahrenheit
	TargetTemperature int `json:"target_temperature"`
	// CurrentTemperature is the current temperature in Fahrenheit as represented by the thermometer (currently not in use)
	CurrentTemperature int `json:"current_temperature"`
}

// Mode is a mode that the thermostat can be in
type Mode string

const (
	// ModeCool is the cooling mode
	ModeCool Mode = "COOL"
	// ModeDry is the "dry" mode
	ModeDry Mode = "DRY"
	// ModeHeat is the heating mode
	ModeHeat Mode = "HEAT"
	// ModeFan is the fan-only mode
	ModeFan Mode = "FAN"
)

// IsValid returns true if the mode is valid
func (m Mode) IsValid() bool {
	switch m {
	case ModeCool, ModeDry, ModeHeat, ModeFan:
		return true
	}

	return false
}

// FanSpeed describes a speed at which the fan can run at
type FanSpeed string

const (
	// FanSpeedAuto will make the A/C determine what speed the fan should run at automatically
	FanSpeedAuto FanSpeed = "AUTO"
	// FanSpeedQuiet is the lowest fan speed (1/4)
	FanSpeedQuiet FanSpeed = "QUIET"
	// FanSpeedLow is the second lowest fan speed (2/4)
	FanSpeedLow FanSpeed = "LOW"
	// FanSpeedMedium is the second highest fan speed (3/4)
	FanSpeedMedium FanSpeed = "MEDIUM"
	// FanSpeedHigh is the highest fan speed (4/4)
	FanSpeedHigh FanSpeed = "HIGH"
)

// IsValid returns true if the fan speed is valid
func (f FanSpeed) IsValid() bool {
	switch f {
	case FanSpeedAuto, FanSpeedQuiet, FanSpeedLow, FanSpeedMedium, FanSpeedHigh:
		return true
	}

	return false
}

// IsValid returns true if the state is valid
func (s State) IsValid() bool {
	if !s.CurrentMode.IsValid() || !s.FanSpeed.IsValid() {
		return false
	}

	tempFloor := 64
	tempCeil := 88
	if s.CurrentMode == ModeHeat {
		tempFloor = 60 // what??????????
		tempCeil = 76  // heat only supports up to 76???
	}

	// Check that the temperature isn't outside the range, and that it's an even number
	if s.TargetTemperature < tempFloor || s.TargetTemperature > tempCeil || s.TargetTemperature%2 != 0 {
		return false
	}

	return true
}

var modeToLirc = map[Mode]string{
	ModeHeat: "heat",
	ModeCool: "cool",
	ModeDry:  "dry",
	ModeFan:  "fan",
}

var fanSpeedToLirc = map[FanSpeed]string{
	FanSpeedAuto:   "auto",
	FanSpeedHigh:   "high",
	FanSpeedMedium: "medium",
	FanSpeedLow:    "low",
	FanSpeedQuiet:  "quiet",
}

// ToLirc turns the state into a LIRC command
func (s State) ToLirc() (string, error) {
	if !s.IsValid() {
		return "", errors.New("invalid state")
	}

	if !s.PoweredOn {
		return "turn-off", nil
	}

	cmd := fmt.Sprintf("%s-%s", modeToLirc[s.CurrentMode], fanSpeedToLirc[s.FanSpeed])
	if s.CurrentMode != ModeFan {
		cmd += fmt.Sprintf("-%dF", s.TargetTemperature)
	}

	return cmd, nil
}

// ToLircOn turns the state into a LIRC *-on command, to be sent before the normal command is sent
func (s State) ToLircOn() (string, error) {
	if !s.IsValid() {
		return "", errors.New("invalid state")
	}

	if !s.PoweredOn {
		return "turn-off", nil
	}

	return fmt.Sprintf("%s-on", modeToLirc[s.CurrentMode]), nil
}

// SetPower sets the power and validates it
func (s *State) SetPower(power bool) error {
	draft := *s
	draft.PoweredOn = power
	if !draft.IsValid() {
		return errors.New("invalid state")
	}

	s.PoweredOn = power
	return nil
}

// SetMode sets the current mode and validates it
func (s *State) SetMode(mode Mode) error {
	draft := *s
	draft.CurrentMode = mode
	if !draft.IsValid() {
		return errors.New("invalid state")
	}

	s.CurrentMode = mode
	return nil
}

// SetFanSpeed sets the fan speed and validates it
func (s *State) SetFanSpeed(speed FanSpeed) error {
	draft := *s
	draft.FanSpeed = speed
	if !draft.IsValid() {
		return errors.New("invalid state")
	}

	s.FanSpeed = speed
	return nil
}

// SetTargetTemperature sets the temp and validates it
func (s *State) SetTargetTemperature(temp int) error {
	draft := *s
	draft.TargetTemperature = temp
	if !draft.IsValid() {
		return errors.New("invalid state")
	}

	s.TargetTemperature = temp
	return nil
}
