// +build darwin

package powerstrip

import (
	"github.com/go-playground/log"
	"time"
)

// var pins = [...] rpio.Pin { rpio.Pin(17), rpio.Pin(27), rpio.Pin(22), rpio.Pin(5), rpio.Pin(6), rpio.Pin(13), rpio.Pin(19), rpio.Pin(26) }
// var pins [8] rpio.Pin

func SendSwitchStatusChangeEvent(switch_name string, on bool) {
	log.Infof("Reporting switch %s status %v", switch_name, on)
}

func InitRpioPins() {
}

func TurnAllOn(timeout time.Duration) {
	log.Info("Toggling all pins ON")
}

func TurnOffOutletByName( name string, force bool ) {
}

func isOutletOn( name string ) bool {
	return false
}

func TurnOnOutletByName( name string, force bool ) {
}

func TurnAllOff(timeout time.Duration) {
	print("Toggling pins OFF")
}

func TurnOnOutlet( index int ) {
}

func TurnOffOutlet( index int ) {
}

///
func runPinToggler(isTest bool) {
}
