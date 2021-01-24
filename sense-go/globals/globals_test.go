package globals

import (
	"github.com/go-playground/log"
	"testing"
)

func TestConfigureLogging(t *testing.T) {
	type args struct {
		config        Configuration
		containerName string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{Configuration{LogLevel: "error,warn,info,debug,notice,panic"}, "sense-go"}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConfigureLogging(tt.args.config,tt.args.containerName)
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
		{name: "happy",want: 1,},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSequence(); got < 0 || got >=300000 {
				t.Errorf("GetSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFromPersistentStore(t *testing.T) {
	type args struct {
		storeMountPoint      string
		relativePath         string
		fileName             string
		config               *Configuration
		currentStageSchedule *StageSchedule
	}
	config := Configuration{}
	stageSchedule := StageSchedule{}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "happy",
			args: args{ storeMountPoint: ".", relativePath: "", fileName: "config.json", config: &config, currentStageSchedule: &stageSchedule},
			wantErr: true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadFromPersistentStore(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName, tt.args.config, tt.args.currentStageSchedule); (err != nil) != tt.wantErr {
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
