package globals

import (
	"github.com/go-playground/log"
	"testing"
)

func init_config() {
	MyDeviceID = 70000007
}

func TestConfigureLogging(t *testing.T) {
	type args struct {
		farm          Farm
		containerName string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy1", args: args{Farm{LogLevel: "error,warn,info,debug,notice,panic"}, "sense-go"}},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConfigureLogging(tt.args.farm,tt.args.containerName)
		})
	}
}

func TestCustomHandler_Log(t *testing.T) {
	type args struct {
		e log.Entry
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CustomHandler{}
			t.Logf("%v", c)
		})
	}
}

func TestGetSequence(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		// TODO: Add test cases.
		{name: "happy",want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSequence(); got < 0 || got >=300000 {
				t.Errorf("GetSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func nextGlobalMisConfig() {
	if MyFarm.ControllerAPIPort == 3003 && MyFarm.ControllerHostName == "localhost" && MyFarm.UserID == 90000009 {
		MyFarm.ControllerAPIPort = 0
	} else {
		if MyFarm.ControllerAPIPort == 0 {
			MyFarm.ControllerAPIPort = 3003
			MyFarm.ControllerHostName = "localhost"
			MyFarm.UserID = -1
		} else {
			MyFarm.ControllerAPIPort = 3003
			MyFarm.ControllerHostName = "blahblah"
			MyFarm.UserID = 90000009
		}
	}
}

func Test_getConfigFromServer(t *testing.T) {
	MyFarm.ControllerAPIPort = 3003
	MyFarm.ControllerHostName = "localhost"
	MyFarm.UserID = 90000009
	MyDevice = &EdgeDevice{ DeviceID: int64(70000007) }

	tests := []struct {
		name    string
		wantErr bool
	}{
		{ name: "happy", wantErr: false},
		{ name: "bad_port", wantErr: true},
		{ name: "bad_user", wantErr: true},
		{ name: "bad_host", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfigFromServer(".","", "config.json"); (err != nil) != tt.wantErr {
				t.Errorf("getConfigFromServer() error = %v, wantErr %v", err, tt.wantErr)
			}
			nextGlobalMisConfig()
		})
	}
}

func TestReadFromPersistentStore(t *testing.T) {
	init_config()

	type args struct {
		storeMountPoint      string
		relativePath         string
		fileName             string
		farm                 *Farm
		currentStageSchedule *StageSchedule
	}
	config := Farm{}
	stageSchedule := StageSchedule{}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Read valid config file with plausible data",
			args: args{ storeMountPoint: ".", relativePath: "", fileName: "config.json", farm: &config, currentStageSchedule: &stageSchedule},
			wantErr: false},
		{name: "Read non-existent config file",
			args: args{ storeMountPoint: "/notavaliddirectoryname", relativePath: "", fileName: "config.json", farm: &config, currentStageSchedule: &stageSchedule},
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadFromPersistentStore(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName, tt.args.farm, tt.args.currentStageSchedule); (err != nil) != tt.wantErr {
				t.Errorf("ReadFromPersistentStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReportDeviceFailed(t *testing.T) {
	type args struct {
		devicename string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "happy", args: args{devicename:"testdevice"},},
	}
	for _, tt := range tests {
		DevicesFailed = []string{}
		t.Run(tt.name, func(t *testing.T) {
			ReportDeviceFailed(tt.args.devicename)
			if len(DevicesFailed) == 0 {
				t.Errorf("DevicesFailed length %d", len(DevicesFailed))
			}
		})
	}
}
