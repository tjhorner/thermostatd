package thermostat

import (
	"errors"
	"fmt"
	"time"

	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"

	"github.com/chbmuc/lirc"
)

// Thermostat represents an active thermostat that can receive commands
type Thermostat struct {
	// State is the current state of the thermostat. DO NOT CHANGE THIS DIRECTLY
	State State `json:"state"`
	lircd *lirc.Router
	lcd   *synchronized.SynchronizedLCD
}

// New returns a newly-created Thermostat object
func New(lircd *lirc.Router, lcd *synchronized.SynchronizedLCD) *Thermostat {
	return &Thermostat{NewState(), lircd, lcd}
}

func (t *Thermostat) sendCommand(cmd string) error {
	return t.lircd.Send(fmt.Sprintf("fujitsu_heat_ac %s", cmd))
}

func (t *Thermostat) sendOnState() error {
	if t.State.CurrentMode == ModeHeat {
		return nil
	}

	onCmd, err := t.State.ToLircOn()
	if err != nil {
		return err
	}

	if err = t.sendCommand(onCmd); err != nil {
		return err
	}

	return nil
}

func (t *Thermostat) sendCurrentState() error {
	cmd, err := t.State.ToLirc()
	if err != nil {
		return err
	}

	if err = t.sendCommand(cmd); err != nil {
		return err
	}

	if t.lcd != nil {
		if t.State.PoweredOn {
			t.lcd.WriteLines(
				padForLcd(fmt.Sprintf("Temp: %dF (%s)", t.State.TargetTemperature, t.State.CurrentMode)),
				padForLcd(fmt.Sprintf("Fan: %s", t.State.FanSpeed)),
			)
		} else {
			t.lcd.WriteLines(padForLcd("Thermostat Off"), "")
		}
	}

	return nil
}

// Reset changes the thermostat back to default values
func (t *Thermostat) Reset() error {
	t.State = NewState()
	return t.sendCurrentState()
}

// SetPower turns the thermostat on/off
func (t *Thermostat) SetPower(power bool) error {
	err := t.State.SetPower(power)
	if err != nil {
		return err
	}

	if power {
		if err = t.sendOnState(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return t.sendCurrentState()
}

// SetTargetTemperature sets the target temperature
func (t *Thermostat) SetTargetTemperature(temp int) error {
	err := t.State.SetTargetTemperature(temp)
	if err != nil {
		return err
	}

	return t.sendCurrentState()
}

// SetMode sets the current mode
func (t *Thermostat) SetMode(mode Mode) error {
	err := t.State.SetMode(mode)
	if err != nil {
		return err
	}

	return t.sendCurrentState()
}

// SetFanSpeed sets the current fan speed
func (t *Thermostat) SetFanSpeed(speed FanSpeed) error {
	err := t.State.SetFanSpeed(speed)
	if err != nil {
		return err
	}

	return t.sendCurrentState()
}

// SetState sets the entire state at once
func (t *Thermostat) SetState(state State) error {
	if !state.IsValid() {
		return errors.New("invalid state")
	}

	sendOn := !t.State.PoweredOn && state.PoweredOn
	t.State = state

	if sendOn {
		if err := t.sendOnState(); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return t.sendCurrentState()
}
