package gpiorelay

import (
	"time"
)

var PowerstripSvc PowerstripService = GetPowerstripService()

type PowerstripService interface {
	SendSwitchStatusChangeEvent(switch_name string, on bool)
	InitRpioPins()
	TurnAllOn(timeout time.Duration)
	TurnOffOutletByName(name string, force bool)
	isOutletOn(name string) bool
	TurnOnOutletByName(name string, force bool)
	ReportAll(timeout time.Duration)
	TurnAllOff(timeout time.Duration)
	TurnOnOutlet(index int)
	TurnOffOutlet(index int)
	RunPinToggler(isTest bool)
	IsMySwitch(switchName string) bool
}
