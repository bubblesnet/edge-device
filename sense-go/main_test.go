package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"testing"
)

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
globals.Lasttemp
globals.ExternalCurrentState.TempF
 */
func testHeat( t *testing.T) {
	globals.CurrentStageSchedule.EnvironmentalTargets.Temperature = 80
	globals.ExternalCurrentState.TempF = globals.TEMPNOTSET
	globals.MyDevice = &globals.EdgeDevice{DeviceID: 0}
	ControlHeat(true)

	// all set
	globals.Lasttemp = 80
	globals.ExternalCurrentState.TempF = 77
	ControlHeat(true)

	globals.Lasttemp = 79
	globals.ExternalCurrentState.TempF = 77
	ControlHeat(true)

	globals.Lasttemp = 79
	globals.ExternalCurrentState.TempF = 79
	ControlHeat(true)

	globals.Lasttemp = 81
	globals.ExternalCurrentState.TempF = 80
	ControlHeat(true)

	globals.Lasttemp = 81
	globals.ExternalCurrentState.TempF = 83
	ControlHeat(true)

}
func TestControlHumidity(t *testing.T) {
	tests := []struct {
		name string
	}{
		{ name: "all" },
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHumidity(t)
		})
	}
}
/*
globals.Lasthumidity = globals.ExternalCurrentState.Humidity
globals.CurrentStageSchedule.EnvironmentalTargets.Humidity
 */
var humidifierstates = []bool{true,false}
func testHumidity(t *testing.T) {

	globals.CurrentStageSchedule.EnvironmentalTargets.Humidity = 60


	for i := 0; i < len(humidifierstates); i++ {
		globals.Lasthumidity = 59
		globals.ExternalCurrentState.Humidity = 50
		ControlHumidity(true)
	}

	for i := 0; i < len(humidifierstates); i++ {
		globals.Lasthumidity = 61
		globals.ExternalCurrentState.Humidity = 67
		ControlHumidity(true)
	}

	for i := 0; i < len(humidifierstates); i++ {
		globals.Lasthumidity = 60
		globals.ExternalCurrentState.Humidity = 60
		ControlHumidity(true)
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

/*
globals.MySite.Stage
globals.MySite.LightOnHour
globals.CurrentStageSchedule.HoursOfLight
globals.LocalCurrentState.GrowLightVeg
"germination"
"seedling"
"vegetative"
*/
var stages = []string {
"germination","seedling","vegetative","idle",
}
var growlightstates = []bool {
	true,false,
}
func testLight(t *testing.T) {
	globals.MyStation = &globals.Station{CurrentStage: "idle"}
	for i := 0; i < len(stages); i++  {
		globals.MyStation.CurrentStage = stages[i]
		for n := 1; n <= 24; n++ {
			globals.CurrentStageSchedule.HoursOfLight = n
			for h := 0; h < 24; h++ {
				globals.MyStation.LightOnHour = h
				for k := 0; k < len(growlightstates); k++ {
					globals.LocalCurrentState.GrowLightVeg = growlightstates[k]
					ControlLight(true)
				}
			}
		}
	}
}


func Test_deviceShouldBeHere(t *testing.T) {
	type args struct {
		containerName   string
		mydeviceid      int64
		deviceInStation bool
		deviceType      string
	}
	tests := []struct {
		name                string
		args                args
		wantShouldBePresent bool
	}{
		{ name: "happy", args: args{ containerName: "sense-go", mydeviceid: 70000007, deviceInStation: true, deviceType: "test"}},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotShouldBePresent := moduleShouldBeHere(tt.args.containerName, tt.args.mydeviceid, tt.args.deviceInStation, tt.args.deviceType); gotShouldBePresent != tt.wantShouldBePresent {
				t.Errorf("moduleShouldBeHere() = %v, want %v", gotShouldBePresent, tt.wantShouldBePresent)
			}
		})
	}
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
				t.Errorf("inRange() = %v, want %v", got, tt.want)
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
		{name: "happy", args: args{deviceid: 70000007}, wantRelayIsAttached: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRelayIsAttached := isRelayAttached(tt.args.deviceid); gotRelayIsAttached != tt.wantRelayIsAttached {
				t.Errorf("isRelayAttached() = %v, want %v", gotRelayIsAttached, tt.wantRelayIsAttached)
			}
		})
	}
}

func Test_makeControlDecisions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
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
			if err := ReadPh(true); (err != nil) != tt.wantErr {
				t.Errorf("ReadPh() error = %v, wantErr %v", err, tt.wantErr)
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
			RunDistanceWatcher()
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
			RunTamperDetector()
		})
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
			setEnvironmentalControlString()
		})
	}
}
