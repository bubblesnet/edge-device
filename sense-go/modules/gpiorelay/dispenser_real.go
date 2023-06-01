//go:build (linux && arm) || arm64
// +build linux,arm arm64

/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package gpiorelay

// copyright and license inspection - no issues 4/13/22

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

var dispenserPins [9]rpio.Pin

type RealDispenser struct {
	Real bool
}

var singletonDispenser = RealDispenser{Real: true}

func (m *RealDispenser) SetupDispenserGPIO(MyStation *globals.Station, MyDevice *globals.EdgeDevice) {
	if len((*MyStation).Dispensers) > 0 {
		log.Infof("Relay is attached to device %d", MyDevice.DeviceID)
		GetDispenserService().InitRpioPins(MyStation, MyDevice, globals.RunningOnUnsupportedHardware())
		GetDispenserService().TurnAllOff(MyStation, MyDevice, 1)
		if (*MyStation).AutomaticControl {
			initializeDispensersForAutomation()
		} else {
			initializeDispensersFromConfiguration()
		}
		GetDispenserService().SendDispenserStatusChangeEvent("automaticControl", (*MyStation).AutomaticControl, globals.GetSequence())
	} else {
		log.Infof("There is no dispenser relay attached to device %d", (*MyDevice).DeviceID)
	}
}

func initializeDispensersForAutomation() {
	log.Infof("initializeDispensersForAutomation - NOOP")
}
func initializeDispensersFromConfiguration() {
	log.Infof("initializeDispensersFromConfiguration - NOOP")
}

func controlDispensers(MyStation *globals.Station, MyDevice *globals.EdgeDevice) {
	log.Infof("RealDispenser.controlDispensers ... do nothing")
	//	GetDispenserService().TimedDispenseSynchronous(MyStation, MyDevice, "pH Up", 5000)
}

func (m *RealDispenser) StartDispensing(MyStation *globals.Station, MyDevice *globals.EdgeDevice) (err error) {
	log.Infof("RealDispenser.StartDispensing")
	for {
		controlDispensers(MyStation, MyDevice)
		time.Sleep(60 * time.Second)
	}
}

func DoADispenserTest(pinnum int) {
	fmt.Printf("DoADispenserTest %d\n", pinnum)
	if err := rpio.Open(); err != nil {
		fmt.Printf("error on pin %d %#v\n", pinnum, err)
		return
	}
	pin := rpio.Pin(pinnum)
	pin.Output()
	for i := 0; i < 10; i++ {
		fmt.Printf("set %d lo \n", pinnum)
		pin.Low()
		time.Sleep(5 * time.Second)
		fmt.Printf("set %d hi \n", pinnum)
		pin.High()
	}
}
func GetDispenserService() DispenserService {
	return &singletonDispenser
}

func (r *RealDispenser) SendDispenserStatusChangeEvent(dispenser_name string, on bool, sequence int32) {
	log.Infof("Reporting switch %s status %#v", dispenser_name, on)
	dm := messaging.NewDispenserStatusChangeMessage(dispenser_name, on, messaging.GetNowMillis())
	bytearray, err := json.Marshal(dm)
	message := pb.SensorRequest{Sequence: sequence, TypeId: globals.Grpc_message_typeid_dispenser, Data: string(bytearray)}
	if globals.Client == nil {
		fmt.Printf("No connection to grpc client\n")
	} else {
		_, err = globals.Client.StoreAndForward(context.Background(), &message)
		if err != nil {
			log.Errorf("sendDispenserStatusChangeEvent ERROR %#v", err)
		} else {
			//				log.Debugf("%#v", sensor_reply)
		}
	}
}

func (r *RealDispenser) InitRpioPins(MyStation *globals.Station, MyDevice *globals.EdgeDevice, RunningOnUnsupportedHardware bool) {
	log.Infof("RealDispenser.InitRpioPins")
	for i := 0; i < len((*MyStation).Dispensers); i++ {
		if (*MyStation).Dispensers[i].DeviceID != (*MyDevice).DeviceID {
			log.Infof("RealDispenser.InitRpioPins Not my device %d != %d", (*MyStation).Dispensers[i].DeviceID, (*MyDevice).DeviceID)
			continue
		}
		log.Infof("RealDispenser.InitRpioPins initing BCM%d controlling the device named %s", (*MyStation).Dispensers[i].BCMPinNumber, (*MyStation).Dispensers[i].DispenserName)
		dispenserPins[(*MyStation).Dispensers[i].Index] = rpio.Pin((*MyStation).Dispensers[i].BCMPinNumber)
		log.Infof("RealDispenser.InitRpioPins dispenserPins[%d] = %#v", (*MyStation).Dispensers[i].Index, dispenserPins[(*MyStation).Dispensers[i].Index])
		if RunningOnUnsupportedHardware {
			log.Infof("RealDispenser.InitRpioPins Skipping pin output because we're running on windows")
			continue
		}
		log.Infof("RealDispenser.InitRpioPins Setting BCM%d to output mode", (*MyStation).Dispensers[i].BCMPinNumber)
		dispenserPins[(*MyStation).Dispensers[i].Index].Output()
	}
}

func (r *RealDispenser) TurnOffDispenserByName(MyStation *globals.Station, MyDevice *globals.EdgeDevice, dispenser_name string, force bool) (stateChanged bool) {
	originallyOn := singletonDispenser.IsDispenserOn(MyStation, MyDevice, dispenser_name)
	if !force && !originallyOn {
		//		log.Infof(" %s already OFF!!", dispenser_name)
		//		SendDispenserStatusChangeEvent(dispenser_name,false)
		return false
	}

	for i := 0; i < len((*MyStation).Dispensers); i++ {
		if (*MyStation).Dispensers[i].DeviceID != (*MyDevice).DeviceID {
			log.Infof("RealDispenser.TurnOffDispenserByName Not my device %d != %d", (*MyStation).Dispensers[i].DeviceID, (*MyDevice).DeviceID)
			continue
		}
		if (*MyStation).Dispensers[i].DispenserName == dispenser_name {
			log.Infof("RealDispenser.TurnOffDispenserByName  %s", dispenser_name)
			log.Infof("RealDispenser.TurnOffDispenserByName offbyname found outlet %s at index %d BCM%d", dispenser_name, MyStation.Dispensers[i].Index, MyStation.Dispensers[i].BCMPinNumber)
			(*MyStation).Dispensers[i].OnOff = false
			singletonDispenser.TurnOffDispenserByIndex((*MyStation).Dispensers[i].Index)
			singletonDispenser.SendDispenserStatusChangeEvent(dispenser_name, false, globals.GetSequence())
			return singletonDispenser.IsDispenserOn(MyStation, MyDevice, dispenser_name) != originallyOn
		}
	}
	return false
	//	log.Warnf("Not my switch %s", dispenser_name)
}

func (r *RealDispenser) IsMyDispenser(MyStation *globals.Station, MyDevice *globals.EdgeDevice, dispenserName string) bool {
	if dispenserName == "automaticControl" {
		return true
	}

	for i := 0; i < len((*MyStation).Dispensers); i++ {
		if (*MyStation).Dispensers[i].DeviceID != (*MyDevice).DeviceID {
			log.Infof("RealDispenser.IsMyDispenser Not my device %d != %d", (*MyStation).Dispensers[i].DeviceID, (*MyDevice).DeviceID)
			continue
		}
		if (*MyStation).Dispensers[i].DispenserName == dispenserName {
			return true
		}
	}
	return false
}

func (r *RealDispenser) IsDispenserOn(MyStation *globals.Station, MyDevice *globals.EdgeDevice, dispenserName string) bool {
	for i := 0; i < len((*MyStation).Dispensers); i++ {
		if (*MyStation).Dispensers[i].DeviceID != (*MyDevice).DeviceID {
			log.Infof("RealDispenser.IsDispenserOn Not my device %d != %d", (*MyStation).Dispensers[i].DeviceID, (*MyDevice).DeviceID)
			continue
		}
		if (*MyStation).Dispensers[i].DispenserName == dispenserName {
			return (*MyStation).Dispensers[i].OnOff
		}
	}
	return false
}

func (r *RealDispenser) TurnOnDispenserByName(MyStation *globals.Station, MyDevice *globals.EdgeDevice, dispenserName string, force bool) (stateChanged bool) {
	log.Infof("RealDispenser.TurnOnDispenserByName %s", dispenserName)
	originallyOn := singletonDispenser.IsDispenserOn(MyStation, MyDevice, dispenserName)
	if !force && singletonDispenser.IsDispenserOn(MyStation, MyDevice, dispenserName) {
		log.Infof("RealDispenser.TurnOnDispenserByName %s Already ON!!!!", dispenserName)
		//		SendDispenserStatusChangeEvent(dispenserName,true)
		return false
	}
	for i := 0; i < len((*MyStation).Dispensers); i++ {
		if (*MyStation).Dispensers[i].DeviceID != (*MyDevice).DeviceID {
			log.Infof("RealDispenser.TurnOnDispenserByName Not my device %d != %d", (*MyStation).Dispensers[i].DeviceID, (*MyDevice).DeviceID)
			continue
		}
		if (*MyStation).Dispensers[i].DispenserName == dispenserName {
			log.Infof("RealDispenser.TurnOnDispenserByName turnOnDispenserByName %s force %#v", dispenserName, force)
			log.Infof("RealDispenser.TurnOnDispenserByName onbyname found outlet %s at index %d BCM%d", dispenserName, (*MyStation).Dispensers[i].Index, (*MyStation).Dispensers[i].BCMPinNumber)
			(*MyStation).Dispensers[i].OnOff = true
			singletonDispenser.TurnOnDispenserByIndex((*MyStation).Dispensers[i].Index)
			singletonDispenser.SendDispenserStatusChangeEvent(dispenserName, true, globals.GetSequence())
			return singletonDispenser.IsDispenserOn(MyStation, MyDevice, dispenserName) == originallyOn
		}
	}
	log.Errorf("RealDispenser.TurnOnDispenserByName couldn't find dispenser named %s in %+v", dispenserName, (*MyStation).Dispensers)
	return false
	//	log.Warnf("Not my switch %s", name)
}

func (r *RealDispenser) ReportAll(MyStation *globals.Station, MyDevice *globals.EdgeDevice, timeout time.Duration) {
	fmt.Printf("RealDispenser.ReportAll Reporting all switch statuses [%d]\n", len((*MyStation).Dispensers))
	for i := 0; i < len((*MyStation).Dispensers); i++ {
		if (*MyStation).Dispensers[i].DeviceID != (*MyDevice).DeviceID {
			log.Infof("RealDispenser.ReportAll Not my device %d != %d", (*MyStation).Dispensers[i].DeviceID, (*MyDevice).DeviceID)
			continue
		}
		fmt.Printf("RealDispenser.ReportAll outlet %s\n", (*MyStation).Dispensers[i].DispenserName)
		singletonDispenser.SendDispenserStatusChangeEvent((*MyStation).Dispensers[i].DispenserName, (*MyStation).Dispensers[i].OnOff, globals.GetSequence())
		if timeout > 0 {
			time.Sleep(timeout)
		}
	}
}

func (r *RealDispenser) TurnAllOff(MyStation *globals.Station, MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("RealDispenser.TurnAllOff Toggling dispenserPins OFF")
	for i := 0; i < len((*MyStation).Dispensers); i++ {
		fmt.Printf("RealDispenser.TurnAllOff  Turning off outlet %s\n", (*MyStation).Dispensers[i].DispenserName)
		(*MyStation).Dispensers[i].OnOff = false
		singletonDispenser.TurnOffDispenserByIndex((*MyStation).Dispensers[i].Index)
		//		fmt.Printf("TurnAllOff 1 after\n")
		singletonDispenser.SendDispenserStatusChangeEvent((*MyStation).Dispensers[i].DispenserName, false, globals.GetSequence())
		//		fmt.Printf("TurnAllOff 2 after\n")
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func (r *RealDispenser) TurnOnDispenserByIndex(index int) {
	dispenserPins[index].High()
}

func (r *RealDispenser) TurnOffDispenserByIndex(index int) {
	dispenserPins[index].Low()
}

func (r *RealDispenser) RunPinToggler(MyStation *globals.Station, MyDevice *globals.EdgeDevice, isTest bool) {
	log.Infof("RealDispenser.RunPinToggler dispenserPins %#v", dispenserPins)
	for i := 0; i < len(dispenserPins); i++ {
		log.Debugf("RealDispenser.RunPinTogglersetting up pin[%d] %#v", i, dispenserPins[i])
		if globals.RunningOnUnsupportedHardware() {
			continue
		}
		dispenserPins[i].Output() // Output mode
		dispenserPins[i].High()   // Output mode
	}
	PinsOn := true
	for true {
		if PinsOn == true {
			singletonDispenser.TurnAllOff(MyStation, MyDevice, 1)
			PinsOn = false
		} else {
			turnAllOn(MyStation, MyDevice, 1)
			PinsOn = true
		}
		if isTest {
			return
		}
		time.Sleep(15 * time.Second)
	}
}

/*
turnAllOn is a PRIVATE function so that it doesn't get called by anything other than the test function "toggler"
*/
func turnAllOn(MyStation *globals.Station, MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("RealDispenser.turnAllOn Toggling dispenserPins ON")
	for i := 0; i < len((*MyStation).Dispensers); i++ {
		fmt.Printf("RealDispenser.turnAllOn  Turning on outlet %s\n", (*MyStation).Dispensers[i].DispenserName)
		(*MyStation).Dispensers[i].OnOff = true
		singletonDispenser.TurnOnDispenserByIndex((*MyStation).Dispensers[i].Index)
		//		fmt.Printf("turnAllOn 1 after\n")
		singletonDispenser.SendDispenserStatusChangeEvent((*MyStation).Dispensers[i].DispenserName, true, globals.GetSequence())
		//		fmt.Printf("turnAllOn 2 after\n")
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func (m *RealDispenser) TimedDispenseSynchronous(MyStation *globals.Station, MyDevice *globals.EdgeDevice, dispenserName string, milliseconds int32) (err error) {
	log.Infof("RealDispenser.TimedDispenseSynchronous %s %d turning on", dispenserName, milliseconds)
	GetDispenserService().TurnOnDispenserByName(MyStation, MyDevice, dispenserName, true)
	time.Sleep(time.Millisecond * time.Duration(milliseconds))
	GetDispenserService().TurnOffDispenserByName(MyStation, MyDevice, dispenserName, true)
	log.Infof("RealDispenser.TimedDispenseSynchronous %s %d turned off", dispenserName, milliseconds)
	return nil
}
