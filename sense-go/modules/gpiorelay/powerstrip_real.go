// +build linux,arm

package gpiorelay

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/stianeikeland/go-rpio/v4"
	"golang.org/x/net/context"
	"time"
)

var pins [8]rpio.Pin

type RealPowerStrip struct {
	Real bool
}

var singletonPowerstrip = RealPowerStrip{Real: true}

func DoATest(pinnum int) {
	fmt.Printf("DoATest %d\n", pinnum)
	if err := rpio.Open(); err != nil {
		fmt.Printf("error on pin %d %v\n", pinnum,err)
		return
	}
	pin:= rpio.Pin(pinnum)
	pin.Output()
	for i := 0; i < 10; i++  {
		fmt.Printf("set %d lo \n", pinnum)
		pin.Low()
		time.Sleep(5 * time.Second)
		fmt.Printf("set %d hi \n", pinnum)
		pin.High()
	}
}
func NewPowerstripService() PowerstripService {
	return &singletonPowerstrip
}

func (r *RealPowerStrip) SendSwitchStatusChangeEvent(switch_name string, on bool) {
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

func (r *RealPowerStrip) InitRpioPins() {
	log.Infof("InitRpioPins")
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		log.Infof("initing BCM%d controlling the device named %s", globals.MyDevice.ACOutlets[i].BCMPinNumber, globals.MyDevice.ACOutlets[i].Name)
		pins[globals.MyDevice.ACOutlets[i].Index] = rpio.Pin(globals.MyDevice.ACOutlets[i].BCMPinNumber)
		log.Infof("pins[%d] = %v", globals.MyDevice.ACOutlets[i].Index,pins[globals.MyDevice.ACOutlets[i].Index])
		if globals.RunningOnUnsupportedHardware() {
			log.Infof("Skipping pin output because we're running on windows")
			continue
		}
		log.Debugf("Setting BCM%d to output mode",globals.MyDevice.ACOutlets[i].BCMPinNumber)
		pins[globals.MyDevice.ACOutlets[i].Index].Output()
	}
}

func (r *RealPowerStrip) TurnAllOn(timeout time.Duration) {
	log.Info("Toggling all pins ON")
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		log.Infof("TurnAllOn Turning on outlet %s", globals.MyDevice.ACOutlets[i].Name)
		singletonPowerstrip.TurnOnOutlet(globals.MyDevice.ACOutlets[i].Index)
		singletonPowerstrip.SendSwitchStatusChangeEvent(globals.MyDevice.ACOutlets[i].Name, true)
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func (r *RealPowerStrip) TurnOffOutletByName(name string, force bool) {
	if !force && !singletonPowerstrip.isOutletOn(name) {
		//		log.Infof(" %s already OFF!!", name)
		//		SendSwitchStatusChangeEvent(name,false)
		return
	}
	log.Infof("TurnOffOutletByName %s", name)
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].Name == name {
			log.Infof("offbyname found outlet %s at index %d BCM%d", name, globals.MyDevice.ACOutlets[i].Index,globals.MyDevice.ACOutlets[i].BCMPinNumber)
			singletonPowerstrip.TurnOffOutlet(globals.MyDevice.ACOutlets[i].Index)
			singletonPowerstrip.SendSwitchStatusChangeEvent(name, false)
			return
		}
	}
	log.Errorf("error: couldn't find outlet named %s", name)
}

func (r *RealPowerStrip) isOutletOn(name string) bool {
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].Name == name {
			return globals.MyDevice.ACOutlets[i].PowerOn
		}
	}
	return false
}

func (r *RealPowerStrip) TurnOnOutletByName(name string, force bool) {
	if !force && singletonPowerstrip.isOutletOn(name) {
		//		log.Debugf("Already ON!!!!")
		//		SendSwitchStatusChangeEvent(name,true)
		return
	}
	log.Infof("turnOnOutletByName %s force %v", name, force)
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].Name == name {
			log.Infof("onbyname found outlet %s at index %d BCM%d", name, globals.MyDevice.ACOutlets[i].Index,globals.MyDevice.ACOutlets[i].BCMPinNumber)
			singletonPowerstrip.TurnOnOutlet(globals.MyDevice.ACOutlets[i].Index)
			singletonPowerstrip.SendSwitchStatusChangeEvent(name, true)
			return
		}
	}
	log.Errorf("error: couldn't find outlet named %s", name)
}

func (r *RealPowerStrip) TurnAllOff(timeout time.Duration) {
	print("Toggling pins OFF")
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		fmt.Printf("TurnAllOff Turning off outlet %s\n", globals.MyDevice.ACOutlets[i].Name)
		singletonPowerstrip.TurnOffOutlet(globals.MyDevice.ACOutlets[i].Index)
		fmt.Printf("TurnAllOff 1 after\n")
		singletonPowerstrip.SendSwitchStatusChangeEvent(globals.MyDevice.ACOutlets[i].Name, false)
		fmt.Printf("TurnAllOff 2 after\n")
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func (r *RealPowerStrip) TurnOnOutlet(index int) {
	pins[index].Low()
}

func (r *RealPowerStrip) TurnOffOutlet(index int) {
	pins[index].High()
}

func (r *RealPowerStrip) RunPinToggler(isTest bool) {
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
			singletonPowerstrip.TurnAllOff(1)
			PinsOn = false
		} else {
			singletonPowerstrip.TurnAllOn(1)
			PinsOn = true
		}
		if isTest {
			return
		}
		time.Sleep(15 * time.Second)
	}
}
