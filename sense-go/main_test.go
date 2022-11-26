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
package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/accelerometer"
	"bubblesnet/edge-device/sense-go/modules/distancesensor"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"bubblesnet/edge-device/sense-go/modules/phsensor"

	"errors"
	"github.com/go-playground/log"
	"github.com/go-stomp/stomp"
	"testing"
)

func init() {
	globals.MyDeviceID = 70000008
	if err := globals.ReadCompleteSiteFromPersistentStore("./testdata", "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
		log.Errorf("ReadCompleteSiteFromPersistentStore() error = %#v", err)
	}

}
func TestControlHeat(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "all"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHeat(t)
		})
	}
}

/*
globals.CurrentStageSchedule.EnvironmentalTargets.Temperature
globals.LastTemp
globals.ExternalCurrentState.TempAirMiddle
*/
func testHeat(t *testing.T) { //				,
	globals.CurrentStageSchedule.EnvironmentalTargets.Temperature = 80
	globals.ExternalCurrentState.TempAirMiddle = globals.TEMPNOTSET
	globals.MyDevice = &globals.EdgeDevice{DeviceID: 0}
	ControlHeat(true, globals.MyDevice.DeviceID, globals.MyDevice, globals.CurrentStageSchedule.Name, globals.CurrentStageSchedule,
		globals.ExternalCurrentState, &globals.LocalCurrentState, &globals.LastTemp, gpiorelay.GetPowerstripService())

	// all set
	globals.LastTemp = 80
	globals.ExternalCurrentState.TempAirMiddle = 77
	ControlHeat(true, globals.MyDevice.DeviceID, globals.MyDevice, globals.CurrentStageSchedule.Name, globals.CurrentStageSchedule,
		globals.ExternalCurrentState, &globals.LocalCurrentState, &globals.LastTemp, gpiorelay.GetPowerstripService())

	globals.LastTemp = 79
	globals.ExternalCurrentState.TempAirMiddle = 77
	ControlHeat(true, globals.MyDevice.DeviceID, globals.MyDevice, globals.CurrentStageSchedule.Name, globals.CurrentStageSchedule,
		globals.ExternalCurrentState, &globals.LocalCurrentState, &globals.LastTemp, gpiorelay.GetPowerstripService())

	globals.LastTemp = 79
	globals.ExternalCurrentState.TempAirMiddle = 79
	ControlHeat(true, globals.MyDevice.DeviceID, globals.MyDevice, globals.CurrentStageSchedule.Name, globals.CurrentStageSchedule,
		globals.ExternalCurrentState, &globals.LocalCurrentState, &globals.LastTemp, gpiorelay.GetPowerstripService())

	globals.LastTemp = 81
	globals.ExternalCurrentState.TempAirMiddle = 80
	ControlHeat(true, globals.MyDevice.DeviceID, globals.MyDevice, globals.CurrentStageSchedule.Name, globals.CurrentStageSchedule,
		globals.ExternalCurrentState, &globals.LocalCurrentState, &globals.LastTemp, gpiorelay.GetPowerstripService())

	globals.LastTemp = 81
	globals.ExternalCurrentState.TempAirMiddle = 83
	ControlHeat(true, globals.MyDevice.DeviceID, globals.MyDevice, globals.CurrentStageSchedule.Name, globals.CurrentStageSchedule,
		globals.ExternalCurrentState, &globals.LocalCurrentState, &globals.LastTemp, gpiorelay.GetPowerstripService())

}

/*
globals.LastHumidity = globals.ExternalCurrentState.HumidityInternal
globals.CurrentStageSchedule.EnvironmentalTargets.HumidityInternal
*/

/*
globals.MySite.Stage
globals.MySite.LightOnHour
globals.CurrentStageSchedule.HoursOfLight
globals.LocalCurrentState.GrowLightVeg
"germination"
"seedling"
"vegetative"
*/
func initGlobalsLocally(t *testing.T) {
	globals.MyDeviceID = 70000008
	if err := globals.ReadCompleteSiteFromPersistentStore("./testdata", "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
		t.Errorf("getConfigFromServer() error = %#v", err)
	}
	if globals.MyStation == nil {
		t.Error("mystation is nil")
	}
}

func Test_moduleShouldBeHere(t *testing.T) {
	initGlobalsLocally(t)

	type args struct {
		containerName   string
		mydeviceid      int64
		myStation       *globals.Station
		deviceInStation bool
		moduleType      string
	}
	tests := []struct {
		name                string
		args                args
		wantShouldBePresent bool
	}{
		{name: "happy",
			wantShouldBePresent: true,
			args: args{
				containerName:   "sense-python",
				mydeviceid:      globals.MyDeviceID,
				myStation:       globals.MyStation,
				deviceInStation: true,
				moduleType:      "bme280",
			},
		},
		{name: "unhappy",
			wantShouldBePresent: false,
			args: args{
				containerName:   "sense-python",
				mydeviceid:      70000006,
				myStation:       globals.MyStation,
				deviceInStation: true,
				moduleType:      "bme280",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotShouldBePresent := moduleShouldBeHere(tt.args.containerName, tt.args.myStation, tt.args.mydeviceid, tt.args.deviceInStation, tt.args.moduleType); gotShouldBePresent != tt.wantShouldBePresent {
				t.Errorf("moduleShouldBeHere(%#v)", tt.args)
				t.Errorf("moduleShouldBeHere got  %#v", gotShouldBePresent)
				t.Errorf("moduleShouldBeHere want %#v", tt.wantShouldBePresent)
			}
		})
	}
}

func Test_makeControlDecisions(t *testing.T) {
	globals.MyStation.AutomaticControl = true
	var tests []struct {
		name string
	}

	if globals.Client == nil {
		t.Logf("globals.Client is nil - won't work")
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			makeControlDecisions(true)
		})
	}
}

func Test_readPh(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := phsensor.ReadPh(true); (err != nil) != tt.wantErr {
				t.Errorf("ReadPh() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_reportVersion(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportVersion()
		})
	}
}

func Test_runDistanceWatcher(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distancesensor.RunDistanceWatcher(true, true)
		})
	}
}

func Test_runLocalStateWatcher(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLocalStateWatcher()
		})
	}
}

func Test_runTamperDetector(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accelerometer.GetTamperDetectorService().RunTamperDetector(true)
		})
	}
}

func Test_getNowMillis(t *testing.T) {

	tests := []struct {
		name string
		want int64
	}{
		{name: "happy", want: 1000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNowMillis(); got <= 1000000 {
				t.Errorf("getNowMillis() = %#v, want >= %#v", got, tt.want)
			}
		})
	}
}

func Test_countACOutlets(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "happy", want: 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countACOutlets(); got != tt.want {
				t.Errorf("countACOutlets() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_isMySwitch(t *testing.T) {
	globals.MyDeviceID = 70000008
	if err := globals.ReadCompleteSiteFromPersistentStore("./testdata", "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
		t.Errorf("ReadCompleteSiteFromPersistentStore() error = %#v", err)
	}
	/*
		type args struct {
			switchName string
		}
		tests := []struct {
			name string
			args args
			want bool
		}{
			{ name: "nonsense", want: false, args: args{switchName: "blah"}},
			{ name: "auto", want: true, args: args{switchName: "automaticControl"}},
			{ name: "heater", want: true, args: args{switchName: "heater"}},
		}
			for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := IsMySwitch(tt.args.switchName); got != tt.want {
					t.Errorf("IsMySwitch() = %#v, want %#v", got, tt.want)
				}
			})
		}
	*/
}

func Test_initializeOutletsForAutomation(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initializePowerstripForAutomation()
		})
	}
}

func Test_initGlobals(t *testing.T) {
	globals.PersistentStoreMountPoint = "./testdata"
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initGlobals(true)
		})
	}
}

func Test_setupGPIO(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPowerstripGPIO(globals.MyStation, globals.MyDevice, gpiorelay.GetPowerstripService())
		})
	}
}

func Test_setupPhMonitor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPhMonitor()
		})
	}
}

func Test_countGoRoutines(t *testing.T) {
	tests := []struct {
		name      string
		wantCount int
	}{
		{name: "happy", wantCount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCount := countGoRoutines(); gotCount <= tt.wantCount {
				t.Errorf("countGoRoutines() = %#v, want %#v", gotCount, tt.wantCount)
			}
		})
	}
}

func Test_startGoRoutines(t *testing.T) {
	type args struct {
		once_only bool
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{once_only: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startGoRoutines(tt.args.once_only)
		})
	}
}

func Test_testableSubmain(t *testing.T) {
	/*
		globals.MyDeviceID = 70000008
		globals.PersistentStoreMountPoint = "./testdata"
		if err := globals.ReadCompleteSiteFromPersistentStore(globals.PersistentStoreMountPoint, "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
			t.Errorf("ReadCompleteSiteFromPersistentStore() error = %#v", err)
		}
		type args struct {
			isUnitTest bool
		}
		tests := []struct {
			name string
			args args
		}{
			{name: "happy", args: args{isUnitTest: true}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testableSubmain(tt.args.isUnitTest)
			})
		}
	*/
}

func Test_processCommand(t *testing.T) {
	var uninitMessage stomp.Message

	emptyBody := "{}"
	emptyMessage := stomp.Message{
		Body: []byte(emptyBody),
	}
	myswitchOnBody := "{ \"command\": \"switch\", \"switch_name\": \"heater\", \"on\": true }"
	myswitchOnMessage := stomp.Message{
		Body: []byte(myswitchOnBody),
	}
	myswitchOffBody := "{ \"command\": \"switch\", \"switch_name\": \"heater\", \"on\": false }"
	myswitchOffMessage := stomp.Message{
		Body: []byte(myswitchOffBody),
	}
	pictureBody := "{ \"command\": \"picture\" }"
	pictureMessage := stomp.Message{
		Body: []byte(pictureBody),
	}

	notMyswitchBody := "{ \"command\": \"switch\", \"switch_name\": \"blahblah\", \"on\": true }"
	notMyswitchMessage := stomp.Message{
		Body: []byte(notMyswitchBody),
	}
	autoSwitchBody := "{ \"command\": \"switch\", \"switch_name\": \"automaticControl\", \"on\": true }"
	autoSwitchMessage := stomp.Message{
		Body: []byte(autoSwitchBody),
	}

	messageWithError := stomp.Message{
		Body: []byte(myswitchOnBody),
		Err:  errors.New("test error handling"),
	}

	statusmessageBody := "{ \"command\": \"status\" }"
	statusMessage := stomp.Message{
		Body: []byte(statusmessageBody),
	}

	dispenseBody := "{ \"command\": \"dispense\", \"dispenser_name\": \"pH Up\", \"milliseconds\": 100 }"
	dispenseMessage := stomp.Message{
		Body: []byte(dispenseBody),
	}

	stageBody := "{ \"command\": \"stage\", \"stage\": \"IDLE\"}"
	stageMessage := stomp.Message{
		Body: []byte(stageBody),
	}

	//	messageWithTimeout := stomp.Message{
	//		Body: []byte(myswitchOnBody),
	//		Err:  errors.New("timeout"),
	//	}

	type args struct {
		msg *stomp.Message
	}
	tests := []struct {
		name      string
		args      args
		wantResub bool
		wantErr   bool
	}{
		{name: "nil_message", args: args{msg: nil}, wantResub: false, wantErr: false},
		{name: "messageWithError", args: args{msg: &messageWithError}, wantResub: true, wantErr: true},
		//		{name: "messageTimeout", args: args{msg: &messageWithTimeout}, wantResub: true, wantErr: true},
		{name: "uninit_message", args: args{msg: &uninitMessage}, wantResub: false, wantErr: true},
		{name: "emptyMessage", args: args{msg: &emptyMessage}, wantResub: false, wantErr: false},
		{name: "myswitchOnMessage", args: args{msg: &myswitchOnMessage}, wantResub: false, wantErr: false},
		{name: "myswitchOffMessage", args: args{msg: &myswitchOffMessage}, wantResub: false, wantErr: false},
		{name: "notMyswitchMessage", args: args{msg: &notMyswitchMessage}, wantResub: false, wantErr: false},
		{name: "autoSwitchMessage", args: args{msg: &autoSwitchMessage}, wantResub: false, wantErr: false},
		{name: "pictureMessage", args: args{msg: &pictureMessage}, wantResub: false, wantErr: false},
		{name: "dispenseMessage", args: args{msg: &dispenseMessage}, wantResub: false, wantErr: false},
		{name: "statusMessage", args: args{msg: &statusMessage}, wantResub: false, wantErr: false},
		{name: "stageMessage", args: args{msg: &stageMessage}, wantResub: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResub, err := processCommand(tt.args.msg, gpiorelay.GetPowerstripService())
			if (err != nil) != tt.wantErr {
				t.Errorf("processCommand() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if gotResub != tt.wantResub {
				t.Errorf("processCommand() gotResub = %#v, want %#v", gotResub, tt.wantResub)
			}
		})
	}
}

func Test_makeControlDecisions1(t *testing.T) {
	type args struct {
		once_only bool
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{once_only: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//			makeControlDecisions(tt.args.once_only)
		})
	}
}
