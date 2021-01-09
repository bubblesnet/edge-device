package powerstrip

import (
	"fmt"
	"github.com/go-playground/log"
	"github.com/stianeikeland/go-rpio"
	"bubblesnet/edge-device/sense-go/globals"
	"time"
)

// var pins = [...] rpio.Pin { rpio.Pin(17), rpio.Pin(27), rpio.Pin(22), rpio.Pin(5), rpio.Pin(6), rpio.Pin(13), rpio.Pin(19), rpio.Pin(26) }
var pins [8] rpio.Pin

func InitRpioPins() {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		log.Infof("initing BCM%d controlling the device named %s", globals.Config.ACOutlets[i].BCMPinNumber, globals.Config.ACOutlets[i].Name )
		pins[i] = rpio.Pin(globals.Config.ACOutlets[i].BCMPinNumber)
		pins[i].Output()
	}
}

func TurnAllOn(timeout time.Duration) {
	log.Info("Toggling pins ON")
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		TurnOnOutlet(globals.Config.ACOutlets[i].Index)
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func TurnOffOutletByName( name string ) {
	if isOutletOn(name) {
		log.Info(fmt.Sprintf("TurnOffOutletByName %s", name))
	}
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			TurnOffOutlet(globals.Config.ACOutlets[i].Index)
			return
		}
	}
	log.Error(fmt.Sprintf("error: couldn't find outlet named %s", name ))
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
	if !isOutletOn(name) {
		log.Info(fmt.Sprintf("turnOnOutletByName %s", name))
	}
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			TurnOnOutlet(globals.Config.ACOutlets[i].Index)
			return
		}
	}
	log.Error(fmt.Sprintf("error: couldn't find outlet named %s", name ))
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
			pins[index].Low()
			break
		}
	}
}

func TurnOffOutlet( index int ) {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Index == index {
			globals.Config.ACOutlets[i].PowerOn = false
			pins[index].High()
			break
		}
	}
}

func runPinToggler() {
	log.Info(fmt.Sprintf("pins %v", pins))
	for i := 0; i < 8; i++ {
		log.Debug(fmt.Sprintf("setting up pin[%d] %v", i, pins[i]))
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
		time.Sleep(15 * time.Second)
	}
}
