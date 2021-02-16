package powerstrip

import (
	"bubblesnet/edge-device/sense-go/globals"
	"github.com/go-playground/log"
	"github.com/stianeikeland/go-rpio"
	"time"
)

// var pins = [...] rpio.Pin { rpio.Pin(17), rpio.Pin(27), rpio.Pin(22), rpio.Pin(5), rpio.Pin(6), rpio.Pin(13), rpio.Pin(19), rpio.Pin(26) }
var pins [8] rpio.Pin

func InitRpioPins() {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		log.Infof("initing BCM%d controlling the device named %s", globals.Config.ACOutlets[i].BCMPinNumber, globals.Config.ACOutlets[i].Name )
		pins[i] = rpio.Pin(globals.Config.ACOutlets[i].BCMPinNumber)
		if globals.RunningOnUnsupportedHardware() {
			log.Infof("Skipping pin output because we're running on windows")
			continue
		}
		pins[i].Output()
	}
}

func TurnAllOn(timeout time.Duration) {
	log.Info("Toggling all pins ON")
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		TurnOnOutlet(globals.Config.ACOutlets[i].Index)
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func TurnOffOutletByName( name string ) {
	log.Infof("TurnOffOutletByName %s", name)
	if !isOutletOn(name) {
		log.Infof(" %s already OFF!!", name)
		return
	}
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			TurnOffOutlet(globals.Config.ACOutlets[i].Index)
			return
		}
	}
	log.Errorf("error: couldn't find outlet named %s", name )
}

func isOutletOn( name string ) bool {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			return globals.Config.ACOutlets[i].PowerOn
		}
	}
	return false
}

func TurnOnOutletByName( name string ) {
	log.Infof("turnOnOutletByName %s", name)
	if isOutletOn(name) {
		log.Debugf("Already ON!!!!")
		return
	}
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			TurnOnOutlet(globals.Config.ACOutlets[i].Index)
			return
		}
	}
	log.Errorf("error: couldn't find outlet named %s", name )
}

func TurnAllOff(timeout time.Duration) {
	print("Toggling pins OFF")
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		TurnOffOutletByName(globals.Config.ACOutlets[i].Name)
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func TurnOnOutlet( index int ) {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Index == index {
			globals.Config.ACOutlets[i].PowerOn = true
			if globals.RunningOnUnsupportedHardware()  {
				log.Infof("Skipping pin LOW because we're running on windows")
				continue
			}
			log.Debugf("TurnOn setting pin LOW for outlet %s",globals.Config.ACOutlets[i].Name)
			pins[index].Low()
			break
		}
	}
}

func TurnOffOutlet( index int ) {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Index == index {
			globals.Config.ACOutlets[i].PowerOn = false
			if globals.RunningOnUnsupportedHardware()  {
				log.Infof("Skipping pin HIGH because we're running on windows")
				continue
			}
			log.Debugf("TurnOff setting pin HIGH for outlet %s",globals.Config.ACOutlets[i].Name)
			pins[index].High()
			break
		}
	}
}

///
func runPinToggler(isTest bool) {
	log.Infof("pins %v", pins)
	for i := 0; i < 8; i++ {
		log.Debugf("setting up pin[%d] %v", i, pins[i])
		if globals.RunningOnUnsupportedHardware() {
			continue
		}
		pins[i].Output() // Output mode
		pins[i].High()   // Output mode 
	}
	PinsOn := true
	for true {
		if PinsOn == true {
			TurnAllOff(1)
			PinsOn = false
		} else {
			TurnAllOn(1)
			PinsOn = true
		}
		if isTest {
			return
		}
		time.Sleep(15 * time.Second)
	}
}
