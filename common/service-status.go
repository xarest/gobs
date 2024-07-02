package common

type ServiceStatus int

const (
	StatusUninitialized ServiceStatus = iota
	StatusInit          ServiceStatus = iota + 1
	StatusSetup         ServiceStatus = iota + 1
	StatusStart         ServiceStatus = iota + 1
	StatusStop          ServiceStatus = iota + 1
)

func (ss ServiceStatus) String() string {
	switch ss {
	case StatusUninitialized:
		return "Uninitialized"
	case StatusInit:
		return "Init"
	case StatusSetup:
		return "Setup"
	case StatusStart:
		return "Start"
	case StatusStop:
		return "Stop"
	default:
		return "Unknown"
	}
}
