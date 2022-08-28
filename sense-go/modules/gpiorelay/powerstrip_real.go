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

var pins [8]rpio.Pin

type RealPowerStrip struct {
	Real bool
}

var singletonPowerstrip = RealPowerStrip{Real: true}

func DoATest(pinnum int) {
	fmt.Printf("DoATest %d\n", pinnum)
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
func GetPowerstripService() PowerstripService {
	return &singletonPowerstrip
}

func (r *RealPowerStrip) SendSwitchStatusChangeEvent(switch_name string, on bool, sequence int32) {
	log.Infof("Reporting switch %s status %#v", switch_name, on)
	dm := messaging.NewSwitchStatusChangeMessage(switch_name, on)
	bytearray, err := json.Marshal(dm)
	message := pb.SensorRequest{Sequence: sequence, TypeId: globals.Grpc_message_typeid_switch, Data: string(bytearray)}
	if globals.Client == nil {
		fmt.Printf("No connection to grpc client\n")
	} else {
		_, err = globals.Client.StoreAndForward(context.Background(), &message)
		if err != nil {
			log.Errorf("sendSwitchStatusChangeEvent ERROR %#v", err)
		} else {
			//				log.Debugf("%#v", sensor_reply)
		}
	}
}

func (r *RealPowerStrip) InitRpioPins(MyDevice *globals.EdgeDevice, RunningOnUnsupportedHardware bool) {
	log.Infof("InitRpioPins")
	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		log.Infof("initing BCM%d controlling the device named %s", (*MyDevice).ACOutlets[i].BCMPinNumber, (*MyDevice).ACOutlets[i].Name)
		pins[(*MyDevice).ACOutlets[i].Index] = rpio.Pin((*MyDevice).ACOutlets[i].BCMPinNumber)
		log.Infof("pins[%d] = %#v", (*MyDevice).ACOutlets[i].Index, pins[(*MyDevice).ACOutlets[i].Index])
		if RunningOnUnsupportedHardware {
			log.Infof("Skipping pin output because we're running on windows")
			continue
		}
		log.Debugf("Setting BCM%d to output mode", (*MyDevice).ACOutlets[i].BCMPinNumber)
		pins[(*MyDevice).ACOutlets[i].Index].Output()
	}
}

func (r *RealPowerStrip) TurnAllOn(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	log.Info("Toggling all pins ON")
	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		log.Infof("TurnAllOn Turning on outlet %s", (*MyDevice).ACOutlets[i].Name)
		(*MyDevice).ACOutlets[i].PowerOn = true
		singletonPowerstrip.TurnOnOutletByIndex((*MyDevice).ACOutlets[i].Index)
		singletonPowerstrip.SendSwitchStatusChangeEvent((*MyDevice).ACOutlets[i].Name, true, globals.GetSequence())
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func (r *RealPowerStrip) TurnOffOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (stateChanged bool) {
	originallyOn := singletonPowerstrip.IsOutletOn(MyDevice, name)
	if !force && !originallyOn {
		//		log.Infof(" %s already OFF!!", name)
		//		SendSwitchStatusChangeEvent(name,false)
		return false
	}

	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		if MyDevice.ACOutlets[i].Name == name {
			log.Infof("TurnOffOutletByName %s", name)
			log.Infof("offbyname found outlet %s at index %d BCM%d", name, MyDevice.ACOutlets[i].Index, MyDevice.ACOutlets[i].BCMPinNumber)
			(*MyDevice).ACOutlets[i].PowerOn = false
			singletonPowerstrip.TurnOffOutletByIndex((*MyDevice).ACOutlets[i].Index)
			singletonPowerstrip.SendSwitchStatusChangeEvent(name, false, globals.GetSequence())
			return singletonPowerstrip.IsOutletOn(MyDevice, name) != originallyOn
		}
	}
	return false
	//	log.Warnf("Not my switch %s", name)
}

func (r *RealPowerStrip) IsMySwitch(MyDevice *globals.EdgeDevice, switchName string) bool {
	if switchName == "automaticControl" {
		return true
	}

	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		if (*MyDevice).ACOutlets[i].Name == switchName {
			return true
		}
	}
	return false
}

func (r *RealPowerStrip) IsOutletOn(MyDevice *globals.EdgeDevice, name string) bool {
	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		if (*MyDevice).ACOutlets[i].Name == name {
			return (*MyDevice).ACOutlets[i].PowerOn
		}
	}
	return false
}

func (r *RealPowerStrip) TurnOnOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (stateChanged bool) {
	originallyOn := singletonPowerstrip.IsOutletOn(MyDevice, name)
	if !force && singletonPowerstrip.IsOutletOn(MyDevice, name) {
		//		log.Debugf("Already ON!!!!")
		//		SendSwitchStatusChangeEvent(name,true)
		return false
	}
	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		if MyDevice.ACOutlets[i].Name == name {
			log.Infof("turnOnOutletByName %s force %#v", name, force)
			log.Infof("onbyname found outlet %s at index %d BCM%d", name, (*MyDevice).ACOutlets[i].Index, (*MyDevice).ACOutlets[i].BCMPinNumber)
			MyDevice.ACOutlets[i].PowerOn = true
			singletonPowerstrip.TurnOnOutletByIndex(MyDevice.ACOutlets[i].Index)
			singletonPowerstrip.SendSwitchStatusChangeEvent(name, true, globals.GetSequence())
			return singletonPowerstrip.IsOutletOn(MyDevice, name) == originallyOn
		}
	}
	return false
	//	log.Warnf("Not my switch %s", name)
}

func (r *RealPowerStrip) ReportAll(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	fmt.Printf("Reporting all switch statuses [%d]\n", len((*MyDevice).ACOutlets))
	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		fmt.Printf("ReportAll outlet %s\n", (*MyDevice).ACOutlets[i].Name)
		singletonPowerstrip.SendSwitchStatusChangeEvent((*MyDevice).ACOutlets[i].Name, (*MyDevice).ACOutlets[i].PowerOn, globals.GetSequence())
		if timeout > 0 {
			time.Sleep(timeout)
		}
	}
}

func (r *RealPowerStrip) TurnAllOff(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("Toggling pins OFF")
	for i := 0; i < len((*MyDevice).ACOutlets); i++ {
		fmt.Printf("TurnAllOff Turning off outlet %s\n", (*MyDevice).ACOutlets[i].Name)
		(*MyDevice).ACOutlets[i].PowerOn = false
		singletonPowerstrip.TurnOffOutletByIndex((*MyDevice).ACOutlets[i].Index)
		//		fmt.Printf("TurnAllOff 1 after\n")
		singletonPowerstrip.SendSwitchStatusChangeEvent((*MyDevice).ACOutlets[i].Name, false, globals.GetSequence())
		//		fmt.Printf("TurnAllOff 2 after\n")
		if timeout > 0 {
			time.Sleep(timeout * time.Second)
		}
	}
}

func (r *RealPowerStrip) TurnOnOutletByIndex(index int) {
	pins[index].High()
}

func (r *RealPowerStrip) TurnOffOutletByIndex(index int) {
	pins[index].Low()
}

func (r *RealPowerStrip) RunPinToggler(MyDevice *globals.EdgeDevice, isTest bool) {
	log.Infof("pins %#v", pins)
	for i := 0; i < len(pins); i++ {
		log.Debugf("setting up pin[%d] %#v", i, pins[i])
		if globals.RunningOnUnsupportedHardware() {
			continue
		}
		pins[i].Output() // Output mode
		pins[i].High()   // Output mode
	}
	PinsOn := true
	for true {
		if PinsOn == true {
			singletonPowerstrip.TurnAllOff(MyDevice, 1)
			PinsOn = false
		} else {
			singletonPowerstrip.TurnAllOn(MyDevice, 1)
			PinsOn = true
		}
		if isTest {
			return
		}
		time.Sleep(15 * time.Second)
	}
}
