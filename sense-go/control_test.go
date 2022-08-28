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
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"math/rand"
	"testing"
	"time"
)

var humidifierstates = []bool{true, false}

var stages = []string{
	globals.GERMINATION,
	globals.SEEDLING,
	globals.VEGETATIVE,
	globals.BLOOMING,
	globals.CURING,
	globals.DRYING,
	globals.HARVEST,
	globals.IDLE,
}
var growlightstates = []bool{
	true, false,
}

func Test_inRange(t *testing.T) {
	type args struct {
		starthour    int
		numhours     int
		currenthours int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "insingleday", args: args{starthour: 5, numhours: 5, currenthours: 7}, want: true},
		{name: "outlowsingleday", args: args{starthour: 5, numhours: 5, currenthours: 3}, want: false},
		{name: "outhighsingleday", args: args{starthour: 5, numhours: 5, currenthours: 14}, want: false},
		{name: "inacrossdayfirstday", args: args{starthour: 20, numhours: 10, currenthours: 21}, want: true},
		{name: "inacrossdaysecondday", args: args{starthour: 20, numhours: 10, currenthours: 2}, want: true},
		{name: "outacrossdayfirstday", args: args{starthour: 20, numhours: 10, currenthours: 18}, want: false},
		{name: "outacrossdaysecondday", args: args{starthour: 20, numhours: 10, currenthours: 11}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inRange(tt.args.starthour, tt.args.numhours, tt.args.currenthours); got != tt.want {
				t.Errorf("inRange() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_isRelayAttached(t *testing.T) {
	type args struct {
		deviceid int64
	}
	tests := []struct {
		name                string
		args                args
		wantRelayIsAttached bool
	}{
		{name: "happy", args: args{deviceid: 70000008}, wantRelayIsAttached: true},
		{name: "sad", args: args{deviceid: 70000006}, wantRelayIsAttached: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRelayIsAttached := isRelayAttached(tt.args.deviceid); gotRelayIsAttached != tt.wantRelayIsAttached {
				t.Errorf("isRelayAttached() = %#v, want %#v", gotRelayIsAttached, tt.wantRelayIsAttached)
			}
		})
	}
}

func TestControlLight(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLight(t)
		})
	}
}

func testLight(t *testing.T) {
	for x := 0; x < len(stages); x++ {
		globals.MyStation = &globals.Station{}
		for i := 0; i < len(stages); i++ {
			for n := 1; n <= 24; n++ {
				globals.CurrentStageSchedule.HoursOfLight = n
				for h := 0; h < 24; h++ {
					globals.MyStation.Automation.LightOnStartHour = h
					for k := 0; k < len(growlightstates); k++ {
						globals.LocalCurrentState.GrowLightVeg = growlightstates[k]
						ControlLight(true, globals.MyDevice.DeviceID, globals.MyDevice,
							globals.CurrentStageSchedule.Name,
							*globals.MyStation, globals.CurrentStageSchedule,
							&globals.LocalCurrentState, time.Now(), gpiorelay.GetPowerstripService())
					}
				}
			}
		}
	}
}

func Test_setEnvironmentalControlString(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globals.LocalCurrentState.EnvironmentalControl = setEnvironmentalControlString(&globals.LocalCurrentState)
		})
	}
}

func TestControlHumidity(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "all"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHumidity(t)
		})
	}
}

func testHumidity(t *testing.T) {

	for x := 0; x < len(stages); x++ {
		globals.MyStation = &globals.Station{}
		globals.CurrentStageSchedule.EnvironmentalTargets.Humidity = 60

		for i := 0; i < len(humidifierstates); i++ {
			globals.LastHumidity = 59
			globals.ExternalCurrentState.HumidityInternal = 50
			ControlHumidity(true, globals.MyDevice.DeviceID,
				globals.MyDevice,
				globals.CurrentStageSchedule,
				globals.CurrentStageSchedule.Name,
				globals.ExternalCurrentState,
				&globals.LocalCurrentState,
				&globals.LastHumidity,
				gpiorelay.GetPowerstripService())
		}

		for i := 0; i < len(humidifierstates); i++ {
			globals.LastHumidity = 61
			globals.ExternalCurrentState.HumidityInternal = 67
			ControlHumidity(true, globals.MyDevice.DeviceID,
				globals.MyDevice,
				globals.CurrentStageSchedule,
				globals.CurrentStageSchedule.Name,
				globals.ExternalCurrentState,
				&globals.LocalCurrentState,
				&globals.LastHumidity, gpiorelay.GetPowerstripService())
		}

		for i := 0; i < len(humidifierstates); i++ {
			globals.LastHumidity = 60
			globals.ExternalCurrentState.HumidityInternal = 60
			ControlHumidity(true, globals.MyDevice.DeviceID,
				globals.MyDevice,
				globals.CurrentStageSchedule,
				globals.CurrentStageSchedule.Name,
				globals.ExternalCurrentState,
				&globals.LocalCurrentState,
				&globals.LastHumidity, gpiorelay.GetPowerstripService())
		}
	}
}

func TestControlOxygenation(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOxygenation(t)
		})
	}
}

func testOxygenation(t *testing.T) {
	type args struct {
		force        bool
		DeviceID     int64
		MyDevice     *globals.EdgeDevice
		CurrentStage string
		Powerstrip   gpiorelay.PowerstripService
	}
	tests := []struct {
		name                 string
		args                 args
		wantSomethingChanged bool
	}{
		{
			name:                 "happy",
			args:                 args{force: false, DeviceID: globals.MyDeviceID, MyDevice: globals.MyDevice, CurrentStage: "", Powerstrip: gpiorelay.GetPowerstripService()},
			wantSomethingChanged: false,
		},
	}

	for x := 0; x < len(stages); x++ {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotSomethingChanged := ControlOxygenation(tt.args.force, tt.args.DeviceID, tt.args.MyDevice, stages[x], tt.args.Powerstrip); gotSomethingChanged != tt.wantSomethingChanged {
					t.Errorf("ControlOxygenation() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
				}
			})
		}
	}
}

func TestControlRootWater(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRootWater(t)
		})
	}
}

func testRootWater(t *testing.T) {
	type args struct {
		force        bool
		DeviceID     int64
		MyDevice     *globals.EdgeDevice
		CurrentStage string
		Powerstrip   gpiorelay.PowerstripService
	}
	tests := []struct {
		name                 string
		args                 args
		wantSomethingChanged bool
	}{
		{
			name: "happy",
			args: args{
				force:        false,
				DeviceID:     globals.MyDeviceID,
				MyDevice:     globals.MyDevice,
				CurrentStage: "",
				Powerstrip:   gpiorelay.GetPowerstripService(),
			},
			wantSomethingChanged: false,
		},
	}
	for x := 0; x < len(stages); x++ {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotSomethingChanged := ControlRootWater(tt.args.force, tt.args.DeviceID, tt.args.MyDevice, stages[x], tt.args.Powerstrip); gotSomethingChanged != tt.wantSomethingChanged {
					t.Errorf("ControlRootWater() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
				}
			})
		}
	}
}

func TestControlAirflow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAirflow(t)
		})
	}
}

func testAirflow(t *testing.T) {
	type args struct {
		force        bool
		DeviceID     int64
		MyDevice     *globals.EdgeDevice
		CurrentStage string
		Powerstrip   gpiorelay.PowerstripService
	}
	tests := []struct {
		name                 string
		args                 args
		wantSomethingChanged bool
	}{
		{
			name: "happy",
			args: args{
				force:        false,
				DeviceID:     globals.MyDeviceID,
				MyDevice:     globals.MyDevice,
				CurrentStage: "",
				Powerstrip:   gpiorelay.GetPowerstripService(),
			},
			wantSomethingChanged: false,
		},
	}
	for x := 0; x < len(stages); x++ {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotSomethingChanged := ControlAirflow(tt.args.force, tt.args.DeviceID, tt.args.MyDevice, stages[x], tt.args.Powerstrip); gotSomethingChanged != tt.wantSomethingChanged {
					t.Errorf("ControlAirflow() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
				}
			})
		}
	}
}

func TestControlWaterTemp(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWaterTemp(t)
		})
	}
}

func testWaterTemp(t *testing.T) {
	type args struct {
		force                bool
		DeviceID             int64
		MyDevice             *globals.EdgeDevice
		CurrentStage         string
		StageSchedule        globals.StageSchedule
		Powerstrip           gpiorelay.PowerstripService
		LocalCurrentState    *globals.LocalState
		ExternalCurrentState globals.ExternalState
		LastWaterTemp        *float32
	}
	tests := []struct {
		name                 string
		args                 args
		wantSomethingChanged bool
	}{
		{
			name: "happy",
			args: args{
				force:                false,
				DeviceID:             globals.MyDeviceID,
				MyDevice:             globals.MyDevice,
				CurrentStage:         "",
				StageSchedule:        globals.CurrentStageSchedule,
				Powerstrip:           gpiorelay.GetPowerstripService(),
				LocalCurrentState:    &globals.LocalCurrentState,
				ExternalCurrentState: globals.ExternalCurrentState,
				LastWaterTemp:        &globals.LastWaterTemp,
			},
			wantSomethingChanged: false,
		},
	}
	for x := 0; x < len(stages); x++ {
		for _, tt := range tests {
			tt.args.ExternalCurrentState.TempAirMiddle = rand.Float32() * 100.0
			t.Run(tt.name, func(t *testing.T) {
				if gotSomethingChanged := ControlWaterTemp(
					tt.args.force,
					tt.args.DeviceID,
					tt.args.MyDevice,
					tt.args.StageSchedule,
					stages[x],
					tt.args.ExternalCurrentState,
					tt.args.LocalCurrentState,
					tt.args.LastWaterTemp,
					tt.args.Powerstrip,
				); gotSomethingChanged != tt.wantSomethingChanged {
					t.Errorf("ControlWaterTemp() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
				}
			})
		}
	}
}

func TestControlHeat1(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHeat1(t)
		})
	}
}

func testHeat1(t *testing.T) {
	type args struct {
		force                bool
		DeviceID             int64
		MyDevice             *globals.EdgeDevice
		CurrentStage         string
		StageSchedule        globals.StageSchedule
		Powerstrip           gpiorelay.PowerstripService
		LocalCurrentState    *globals.LocalState
		ExternalCurrentState globals.ExternalState
		LastWaterTemp        *float32
	}
	tests := []struct {
		name                 string
		args                 args
		wantSomethingChanged bool
	}{
		{
			name: "happy",
			args: args{
				force:                false,
				DeviceID:             globals.MyDeviceID,
				MyDevice:             globals.MyDevice,
				CurrentStage:         "",
				StageSchedule:        globals.CurrentStageSchedule,
				Powerstrip:           gpiorelay.GetPowerstripService(),
				LocalCurrentState:    &globals.LocalCurrentState,
				ExternalCurrentState: globals.ExternalCurrentState,
				LastWaterTemp:        &globals.LastWaterTemp,
			},
			wantSomethingChanged: false,
		},
	}
	for x := 0; x < len(stages); x++ {
		for _, tt := range tests {
			tt.args.ExternalCurrentState.TempAirMiddle = 100.0 // Too high
			t.Run(tt.name, func(t *testing.T) {
				if gotSomethingChanged := ControlHeat(
					tt.args.force,
					tt.args.DeviceID,
					tt.args.MyDevice,
					stages[x],
					tt.args.StageSchedule,
					tt.args.ExternalCurrentState,
					tt.args.LocalCurrentState,
					tt.args.LastWaterTemp,
					tt.args.Powerstrip,
				); gotSomethingChanged != tt.wantSomethingChanged {
					t.Errorf("ControlWaterTemp() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
				}
			})

			tt.args.ExternalCurrentState.TempAirMiddle = -1.0 // Too low
			t.Run(tt.name, func(t *testing.T) {
				if gotSomethingChanged := ControlHeat(
					tt.args.force,
					tt.args.DeviceID,
					tt.args.MyDevice,
					stages[x],
					tt.args.StageSchedule,
					tt.args.ExternalCurrentState,
					tt.args.LocalCurrentState,
					tt.args.LastWaterTemp,
					tt.args.Powerstrip,
				); gotSomethingChanged != tt.wantSomethingChanged {
					t.Errorf("ControlWaterTemp() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
				}
			})
		}
	}
}
