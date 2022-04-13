package gpiorelay

// copyright and license inspection - no issues 4/13/22

import (
	"bubblesnet/edge-device/sense-go/globals"
	"time"
)

var PowerstripSvc PowerstripService = GetPowerstripService()

type PowerstripService interface {
	SendSwitchStatusChangeEvent(switch_name string, on bool, sequence int32)
	InitRpioPins(MyDevice *globals.EdgeDevice, RunningOnUnsupportedHardware bool)
	TurnAllOn(MyDevice *globals.EdgeDevice, timeout time.Duration)
	TurnOffOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (stateChanged bool)
	IsOutletOn(MyDevice *globals.EdgeDevice, name string) bool
	TurnOnOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (stateChanged bool)
	ReportAll(MyDevice *globals.EdgeDevice, timeout time.Duration)
	TurnAllOff(MyDevice *globals.EdgeDevice, timeout time.Duration)
	TurnOnOutlet(index int)
	TurnOffOutlet(index int)
	RunPinToggler(MyDevice *globals.EdgeDevice, isTest bool)
	IsMySwitch(MyDevice *globals.EdgeDevice, switchName string) bool
}
