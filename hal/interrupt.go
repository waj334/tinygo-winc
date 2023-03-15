package hal

type InterruptPin interface {
	Pin
	Enable(func(Pin)) error
	Disable() error
}
