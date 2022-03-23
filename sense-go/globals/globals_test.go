package globals

import (
	"fmt"
	"github.com/go-playground/log"
	"os"
	"runtime"
	"strings"
	"testing"
)

func init_config() {
	MyDeviceID = 70000008
}

func TestSlugs(t *testing.T) {
	RunningOnUnsupportedHardware()
	Sequence = 180000
	for i := 0; i < 20500; i++ {
		GetSequence()
	}
}
func TestConfigureLogging(t *testing.T) {
	type args struct {
		site          Site
		containerName string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy1", args: args{Site{LogLevel: "error,warn,info,debug,notice,panic"}, "sense-go"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConfigureLogging(tt.args.site, tt.args.containerName)
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
			t.Logf("%#v", c)
		})
	}
}

func TestGetSequence(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		// TODO: Add test cases.
		{name: "happy", want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSequence(); got < 0 || got >= 300000 {
				t.Errorf("GetSequence() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func nextGlobalMisConfig() {
	if MySite.ControllerAPIPort == 3003 && MySite.ControllerHostName == "localhost" && MySite.UserID == 90000009 {
		MySite.ControllerAPIPort = 0
	} else {
		if MySite.ControllerAPIPort == 0 {
			MySite.ControllerAPIPort = 3003
			MySite.ControllerHostName = "localhost"
			MySite.UserID = -1
		} else {
			MySite.ControllerAPIPort = 3003
			MySite.ControllerHostName = "blahblah"
			MySite.UserID = 90000009
		}
	}
}

func Test_getConfigFromServer(t *testing.T) {
	MySite.ControllerAPIPort = 3003
	MySite.ControllerHostName = "localhost"
	MySite.UserID = 90000009
	MyDevice = &EdgeDevice{DeviceID: int64(70000008)}

	ci := false
	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" { /// TODO this is AWFUL CI hack
		ci = true
	} else {
		return
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: ci},
		{name: "bad_port", wantErr: true},
		{name: "bad_user", wantErr: true},
		{name: "bad_host", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfigFromServer("../testdata", "", "config.json"); (err != nil) != tt.wantErr {
				t.Errorf("getConfigFromServer() error = %#v, wantErr %#v", err, tt.wantErr)

			}
			nextGlobalMisConfig()
		})
	}
}

func TestReadFromPersistentStore(t *testing.T) {
	init_config()

	currentWorkingDirectory, _ := os.Getwd()
	fmt.Printf("cwd = %s\n", currentWorkingDirectory)
	type args struct {
		storeMountPoint      string
		relativePath         string
		fileName             string
		site                 *Site
		currentStageSchedule *StageSchedule
	}
	config := Site{}
	stageSchedule := StageSchedule{}
	MyStation = &Station{CurrentStage: IDLE}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Read valid config file with plausible data",
			args:    args{storeMountPoint: "../testdata", relativePath: "", fileName: "config.json", site: &config, currentStageSchedule: &stageSchedule},
			wantErr: false},
		{name: "Read non-existent config file",
			args:    args{storeMountPoint: "/notavaliddirectoryname", relativePath: "", fileName: "config.json", site: &config, currentStageSchedule: &stageSchedule},
			wantErr: true},
	}

	MyStation.CurrentStage = IDLE

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadCompleteSiteFromPersistentStore(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName, tt.args.site, tt.args.currentStageSchedule); (err != nil) != tt.wantErr {
				if strings.Contains(err.Error(), "not found") {

				} else {
					t.Errorf("ReadCompleteSiteFromPersistentStore() error = %#v, wantErr %#v", err, tt.wantErr)
				}
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
		{name: "happy", args: args{devicename: "testdevice"}},
		{name: "devicefailed", args: args{devicename: "testdevice"}},
	}
	DevicesFailed = []string{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReportDeviceFailed(tt.args.devicename)
			if len(DevicesFailed) == 0 {
				t.Errorf("DevicesFailed length %d", len(DevicesFailed))
			}
		})
	}
}
