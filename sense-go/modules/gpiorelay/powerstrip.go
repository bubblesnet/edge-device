package gpiorelay

import (
	"time"
)

var PowerstripSvc PowerstripService = NewPowerstripService()

type PowerstripService interface {
	SendSwitchStatusChangeEvent(switch_name string, on bool)
	InitRpioPins()
	TurnAllOn(timeout time.Duration)
	TurnOffOutletByName(name string, force bool)
	isOutletOn(name string) bool
	TurnOnOutletByName(name string, force bool)
	TurnAllOff(timeout time.Duration)
	TurnOnOutlet(index int)
	TurnOffOutlet(index int)
	RunPinToggler(isTest bool)
}


