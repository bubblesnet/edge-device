//go:build darwin || windows || (linux && amd64)
// +build darwin windows linux,amd64

package gpiorelay

import (
	"bubblesnet/edge-device/sense-go/globals"
	"github.com/go-playground/log"
	"time"
)

type MockPowerStrip struct {
	Real bool
}

var singletonPowerstrip = MockPowerStrip{Real: true}

func GetPowerstripService() PowerstripService {
	return &singletonPowerstrip
}

func (r *MockPowerStrip) IsMySwitch(MyDevice *globals.EdgeDevice, switchName string) bool {
	return true
}

func (m *MockPowerStrip) SendSwitchStatusChangeEvent(switch_name string, on bool, sequence int32) {
	log.Infof("Reporting switch %s status %#v", switch_name, on)
}

func (m *MockPowerStrip) InitRpioPins(MyDevice *globals.EdgeDevice, RunningOnUnsupportedHardware bool) {
}

func (m *MockPowerStrip) TurnAllOn(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	log.Info("Toggling all pins ON")
}

func (m *MockPowerStrip) TurnOffOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (somethingChanged bool) {
	return (false)
}

func (m *MockPowerStrip) isOutletOn(MyDevice *globals.EdgeDevice, name string) bool {
	return false
}

func (m *MockPowerStrip) TurnOnOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (somethingChanged bool) {
	return (false)
}

func (m *MockPowerStrip) ReportAll(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("Reporting ALL")
}

func (m *MockPowerStrip) TurnAllOff(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("Toggling pins OFF")
}

func (m *MockPowerStrip) TurnOnOutlet(index int) {
}

func (m *MockPowerStrip) TurnOffOutlet(index int) {
}

func (m *MockPowerStrip) RunPinToggler(MyDevice *globals.EdgeDevice, isTest bool) {
}
