package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/accelerometer"
	"bubblesnet/edge-device/sense-go/modules/distancesensor"
	"bubblesnet/edge-device/sense-go/modules/phsensor"
	"errors"
	"github.com/go-playground/log"
	"github.com/go-stomp/stomp"
	"testing"
)

func init() {
	globals.MyDeviceID = 70000008
	if err := globals.ReadFromPersistentStore("./testdata", "", "config.json",&globals.MySite, &globals.CurrentStageSchedule ); err != nil {
		log.Errorf("ReadFromPersistentStore() error = %v", err)
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
globals.GERMINATION,globals.SEEDLING,globals.VEGETATIVE,globals.IDLE,
}
var growlightstates = []bool {
	true,false,
}
func testLight(t *testing.T) {
	globals.MyStation = &globals.Station{CurrentStage: globals.IDLE}
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


func Test_moduleShouldBeHere(t *testing.T) {
	globals.MyDeviceID = 70000008
	if err := globals.ReadFromPersistentStore("./testdata", "", "config.json",&globals.MySite, &globals.CurrentStageSchedule ); err != nil {
		log.Errorf("ReadFromPersistentStore() error = %v", err)
	}
	type args struct {
		containerName   string
		mydeviceid      int64
		deviceInStation bool
		moduleType      string
	}
	tests := []struct {
		name                string
		args                args
		wantShouldBePresent bool
	}{
		{ name: "happy", wantShouldBePresent: true, args: args{ containerName: "sense-python", mydeviceid: 70000008, deviceInStation: true, moduleType: "bme280"}},
		{ name: "unhappy", wantShouldBePresent: false, args: args{ containerName: "sense-python", mydeviceid: 70000006, deviceInStation: true, moduleType: "bme280"}},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotShouldBePresent := moduleShouldBeHere(tt.args.containerName, tt.args.mydeviceid, tt.args.deviceInStation, tt.args.moduleType); gotShouldBePresent != tt.wantShouldBePresent {
				t.Errorf("moduleShouldBeHere(%v) = %v, want %v", tt.args, gotShouldBePresent, tt.wantShouldBePresent)
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
	globals.MyStation.AutomaticControl = true
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
			if err := phsensor.ReadPh(true); (err != nil) != tt.wantErr {
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
			distancesensor.RunDistanceWatcher(true)
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
			accelerometer.RunTamperDetector(true)
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


func Test_getNowMillis(t *testing.T) {

	tests := []struct {
		name string
		want int64
	}{
		{name:"happy", want: 1000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNowMillis(); got <= 1000000 {
				t.Errorf("getNowMillis() = %v, want >= %v", got, tt.want)
			}
		})
	}
}

func Test_countACOutlets(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
{name: "happy", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countACOutlets(); got != tt.want {
				t.Errorf("countACOutlets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isMySwitch(t *testing.T) {
	globals.MyDeviceID = 70000007
	if err := globals.ReadFromPersistentStore("./testdata", "", "config.json",&globals.MySite, &globals.CurrentStageSchedule ); err != nil {
		t.Errorf("ReadFromPersistentStore() error = %v", err)
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
				t.Errorf("IsMySwitch() = %v, want %v", got, tt.want)
			}
		})
	}
	 */
}

func Test_initializeOutletsForAutomation(t *testing.T) {
	tests := []struct {
		name string
	}{
		{ name: "happy"},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initializeOutletsForAutomation()
		})
	}
}

func Test_initGlobals(t *testing.T) {
	globals.PersistentStoreMountPoint = "./testdata"
	tests := []struct {
		name string
	}{
		{ name: "happy"},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initGlobals()
		})
	}
}

func Test_setupGPIO(t *testing.T) {
	tests := []struct {
		name string
	}{
		{ name: "happy"},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupGPIO()
		})
	}
}

func Test_setupPhMonitor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{ name: "happy"},
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
		{ name: "happy", wantCount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCount := countGoRoutines(); gotCount <= tt.wantCount {
				t.Errorf("countGoRoutines() = %v, want %v", gotCount, tt.wantCount)
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
		{ name: "happy", args: args{ once_only: true}},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startGoRoutines(tt.args.once_only)
		})
	}
}

func Test_testableSubmain(t *testing.T) {
	globals.MyDeviceID = 70000007
	globals.PersistentStoreMountPoint = "./testdata"
	if err := globals.ReadFromPersistentStore(globals.PersistentStoreMountPoint, "", "config.json",&globals.MySite, &globals.CurrentStageSchedule ); err != nil {
		t.Errorf("ReadFromPersistentStore() error = %v", err)
	}
	type args struct {
		isUnitTest bool
	}
	tests := []struct {
		name string
		args args
	}{
		{ name: "happy", args: args{isUnitTest: true}},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testableSubmain(tt.args.isUnitTest)
		})
	}
}

func Test_processCommand(t *testing.T) {
	var uninitMessage stomp.Message

	emptyBody := "{}"
	emptyMessage := stomp.Message {
		Body: []byte(emptyBody),
	}
	myswitchOnBody := "{ \"command\": \"switch\", \"switch_name\": \"heater\", \"on\": true }"
	myswitchOnMessage := stomp.Message {
		Body: []byte(myswitchOnBody),
	}
	myswitchOffBody := "{ \"command\": \"switch\", \"switch_name\": \"heater\", \"on\": false }"
	myswitchOffMessage := stomp.Message {
		Body: []byte(myswitchOffBody),
	}
	pictureBody := "{ \"command\": \"picture\" }"
	pictureMessage := stomp.Message {
		Body: []byte(pictureBody),
	}

	notMyswitchBody := "{ \"command\": \"switch\", \"switch_name\": \"blahblah\", \"on\": true }"
	notMyswitchMessage := stomp.Message {
		Body: []byte(notMyswitchBody),
	}
	autoSwitchBody := "{ \"command\": \"switch\", \"switch_name\": \"automaticControl\", \"on\": true }"
	autoSwitchMessage := stomp.Message {
		Body: []byte(autoSwitchBody),
	}

	messageWithError := stomp.Message {
		Body: []byte(myswitchOnBody),
		Err: errors.New("test error handling"),
	}
	messageWithTimeout := stomp.Message {
		Body: []byte(myswitchOnBody),
		Err: errors.New("timeout"),
	}

	type args struct {
		msg *stomp.Message
	}
	tests := []struct {
		name      string
		args      args
		wantResub bool
		wantErr   bool
	}{
		{ name: "nil_message", args: args{msg: nil}, wantResub: false, wantErr: false},
		{ name: "messageWithError", args: args{msg: &messageWithError}, wantResub: true, wantErr: true},
		{ name: "messageTimeout", args: args{msg: &messageWithTimeout}, wantResub: true, wantErr: true},
		{ name: "uninit_message", args: args{msg: &uninitMessage}, wantResub: false, wantErr: true},
		{ name: "emptyMessage", args: args{msg: &emptyMessage}, wantResub: false, wantErr: false},
		{ name: "myswitchOnMessage", args: args{msg: &myswitchOnMessage}, wantResub: false, wantErr: false},
		{ name: "myswitchOffMessage", args: args{msg: &myswitchOffMessage}, wantResub: false, wantErr: false},
		{ name: "notMyswitchMessage", args: args{msg: &notMyswitchMessage}, wantResub: false, wantErr: false},
		{ name: "autoSwitchMessage", args: args{msg: &autoSwitchMessage}, wantResub: false, wantErr: false},
		{ name: "pictureMessage", args: args{msg: &pictureMessage}, wantResub: false, wantErr: false},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResub, err := processCommand(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("processCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResub != tt.wantResub {
				t.Errorf("processCommand() gotResub = %v, want %v", gotResub, tt.wantResub)
			}
		})
	}
}