// +build darwin windows linux,amd64

package gpiorelay

import (
	"github.com/go-playground/log"
	"time"
)

type MockPowerStrip struct {
	Real bool
}

var singletonPowerstrip = MockPowerStrip{Real: true}
func NewPowerstripService() PowerstripService {
	return &singletonPowerstrip
}

func (m *MockPowerStrip)SendSwitchStatusChangeEvent(switch_name string, on bool) {
	log.Infof("Reporting switch %s status %v", switch_name, on)
}

func (m *MockPowerStrip)InitRpioPins() {
}

func (m *MockPowerStrip)TurnAllOn(timeout time.Duration) {
	log.Info("Toggling all pins ON")
}

func (m *MockPowerStrip)TurnOffOutletByName( name string, force bool ) {
}

func (m *MockPowerStrip)isOutletOn( name string ) bool {
	return false
}

func (m *MockPowerStrip)TurnOnOutletByName( name string, force bool ) {
}

func (m *MockPowerStrip)ReportAll(timeout time.Duration) {
	print("Reporting ALl")
}

func (m *MockPowerStrip)TurnAllOff(timeout time.Duration) {
	print("Toggling pins OFF")
}

func (m *MockPowerStrip)TurnOnOutlet( index int ) {
}

func (m *MockPowerStrip)TurnOffOutlet( index int ) {
}

func (m *MockPowerStrip)RunPinToggler(isTest bool) {
}
