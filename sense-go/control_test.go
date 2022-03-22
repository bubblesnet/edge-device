package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"testing"
	"time"
)

var humidifierstates = []bool{true, false}

var stages = []string{
	globals.GERMINATION, globals.SEEDLING, globals.VEGETATIVE, globals.IDLE,
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
func testLight(t *testing.T) {
	globals.MyStation = &globals.Station{CurrentStage: globals.IDLE}
	for i := 0; i < len(stages); i++ {
		globals.MyStation.CurrentStage = stages[i]
		for n := 1; n <= 24; n++ {
			globals.CurrentStageSchedule.HoursOfLight = n
			for h := 0; h < 24; h++ {
				globals.MyStation.LightOnHour = h
				for k := 0; k < len(growlightstates); k++ {
					globals.LocalCurrentState.GrowLightVeg = growlightstates[k]
					ControlLight(true, globals.MyDevice.DeviceID, globals.MyStation.CurrentStage,
						*globals.MyStation, globals.CurrentStageSchedule,
						&globals.LocalCurrentState, time.Now(), gpiorelay.GetPowerstripService())
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

	globals.CurrentStageSchedule.EnvironmentalTargets.Humidity = 60

	for i := 0; i < len(humidifierstates); i++ {
		globals.LastHumidity = 59
		globals.ExternalCurrentState.Humidity = 50
		ControlHumidity(true, globals.MyDevice.DeviceID,
			globals.CurrentStageSchedule,
			globals.MyStation.CurrentStage,
			globals.ExternalCurrentState,
			&globals.LocalCurrentState,
			&globals.LastHumidity, gpiorelay.GetPowerstripService())
	}

	for i := 0; i < len(humidifierstates); i++ {
		globals.LastHumidity = 61
		globals.ExternalCurrentState.Humidity = 67
		ControlHumidity(true, globals.MyDevice.DeviceID,
			globals.CurrentStageSchedule,
			globals.MyStation.CurrentStage,
			globals.ExternalCurrentState,
			&globals.LocalCurrentState,
			&globals.LastHumidity, gpiorelay.GetPowerstripService())
	}

	for i := 0; i < len(humidifierstates); i++ {
		globals.LastHumidity = 60
		globals.ExternalCurrentState.Humidity = 60
		ControlHumidity(true, globals.MyDevice.DeviceID,
			globals.CurrentStageSchedule,
			globals.MyStation.CurrentStage,
			globals.ExternalCurrentState,
			&globals.LocalCurrentState,
			&globals.LastHumidity, gpiorelay.GetPowerstripService())
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
