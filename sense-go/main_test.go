package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"github.com/go-playground/log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"reflect"
	"testing"
)

func TestAtlasEZODriver_Connection(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name   string
		fields fields
		want   gobot.Connection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if got := d.Connection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasEZODriver_Halt(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if err := d.Halt(); (err != nil) != tt.wantErr {
				t.Errorf("Halt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAtlasEZODriver_Name(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if got := d.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasEZODriver_Ph(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantPH  float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			gotPH, err := d.Ph()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ph() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPH != tt.wantPH {
				t.Errorf("Ph() gotPH = %v, want %v", gotPH, tt.wantPH)
			}
		})
	}
}

func TestAtlasEZODriver_SetName(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	type args struct {
		n string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			log.Infof("d = %v", d)
		})
	}
}

func TestAtlasEZODriver_Start(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if err := d.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAtlasEZODriver_initialization(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if err := d.initialization(); (err != nil) != tt.wantErr {
				t.Errorf("initialization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAtlasEZODriver_rawPh(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantPH  float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			gotPH, err := d.rawPh()
			if (err != nil) != tt.wantErr {
				t.Errorf("rawPh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPH != tt.wantPH {
				t.Errorf("rawPh() gotPH = %v, want %v", gotPH, tt.wantPH)
			}
		})
	}
}

func TestAtlasEZODriver_read(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	type args struct {
		address byte
		n       int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			got, err := d.read(tt.args.address, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("read() got = %v, want %v", got, tt.want)
			}
		})
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
func TestNewAtlasEZODriver(t *testing.T) {
	type args struct {
		c       i2c.Connector
		options []func(i2c.Config)
	}
	tests := []struct {
		name string
		args args
		want *AtlasEZODriver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAtlasEZODriver(tt.args.c, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAtlasEZODriver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clen(t *testing.T) {
	type args struct {
		n []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clen(tt.args.n); got != tt.want {
				t.Errorf("clen() = %v, want %v", got, tt.want)
			}
		})
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
		{name: "happy", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := readPh(); (err != nil) != tt.wantErr {
				t.Errorf("readPh() error = %v, wantErr %v", err, tt.wantErr)
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
			runDistanceWatcher()
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
			runTamperDetector()
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
