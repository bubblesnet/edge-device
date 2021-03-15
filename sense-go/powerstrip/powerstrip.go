// +build linux, arm

package powerstrip

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/stianeikeland/go-rpio"
	"golang.org/x/net/context"
	"time"
)

// var pins = [...] rpio.Pin { rpio.Pin(17), rpio.Pin(27), rpio.Pin(22), rpio.Pin(5), rpio.Pin(6), rpio.Pin(13), rpio.Pin(19), rpio.Pin(26) }
var pins [8] rpio.Pin

func SendSwitchStatusChangeEvent(switch_name string, on bool) {
	log.Infof("Reporting switch %s status %v", switch_name, on)
	dm := messaging.NewSwitchStatusChangeMessage(switch_name, on)
	bytearray, err := json.Marshal(dm)
	message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "switch", Data: string(bytearray)}
	if globals.Client == nil {
		fmt.Printf("No connection to grpc client\n")
	} else {
		_, err = globals.Client.StoreAndForward(context.Background(), &message)
		if err != nil {
			log.Errorf("sendSwitchStatusChangeEvent ERROR %v", err)
		} else {
			//				log.Debugf("%v", sensor_reply)
		}
	}
}


func InitRpioPins() {
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		log.Infof("initing BCM%d controlling the device named %s", globals.Config.ACOutlets[i].BCMPinNumber, globals.Config.ACOutlets[i].Name )
		pins[i] = rpio.Pin(globals.Config.ACOutlets[i].BCMPinNumber)
		log.Infof("pins[i] = %v", pins[i])
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
		log.Infof("TurnAllOn Turning on outlet %s", globals.Config.ACOutlets[i].Name)
		TurnOnOutlet(globals.Config.ACOutlets[i].Index)
		SendSwitchStatusChangeEvent(globals.Config.ACOutlets[i].Name,true)
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func TurnOffOutletByName( name string, force bool ) {
	if !force && !isOutletOn(name) {
//		log.Infof(" %s already OFF!!", name)
//		SendSwitchStatusChangeEvent(name,false)
		return
	}
	log.Infof("TurnOffOutletByName %s", name)
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			TurnOffOutlet(globals.Config.ACOutlets[i].Index)
			SendSwitchStatusChangeEvent(name,false)
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

func TurnOnOutletByName( name string, force bool ) {
	if !force && isOutletOn(name) {
//		log.Debugf("Already ON!!!!")
//		SendSwitchStatusChangeEvent(name,true)
		return
	}
	log.Infof("turnOnOutletByName %s force %v", name, force)
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].Name == name {
			TurnOnOutlet(globals.Config.ACOutlets[i].Index)
			SendSwitchStatusChangeEvent(name,true)
			return
		}
	}
	log.Errorf("error: couldn't find outlet named %s", name )
}

func TurnAllOff(timeout time.Duration) {
	print("Toggling pins OFF")
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		fmt.Printf("TurnAllOff Turning off outlet %s\n", globals.Config.ACOutlets[i].Name)
		TurnOffOutlet(globals.Config.ACOutlets[i].Index)
		fmt.Printf("TurnAllOff 1 after\n")
		SendSwitchStatusChangeEvent(globals.Config.ACOutlets[i].Name,false)
		fmt.Printf("TurnAllOff 2 after\n")
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
