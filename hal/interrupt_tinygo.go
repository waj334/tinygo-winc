//go:build tinygo

package hal

import "machine"

type interruptPin struct {
	machine.Pin
	Falling machine.PinChange

	mask uint32
}

func NewInterruptPin(pin machine.Pin) InterruptPin {
	return &interruptPin{
		Pin:     pin,
		Falling: machine.PinFalling,
	}
}

func (i *interruptPin) Enable(fn func(Pin)) error {
	if err := i.Pin.SetInterrupt(i.Falling, func(p machine.Pin) {
		fn(p)
	}); err != nil {
		return err
	}
	return nil
}

func (i *interruptPin) Disable() error {
	if err := i.Pin.SetInterrupt(i.Falling, nil); err != nil {
		return err
	}
	return nil
}
