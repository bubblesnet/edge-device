package video

import (
	"bubblesnet/edge-device/sense-go/globals"
	"testing"
)

func Test_uploadFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "goodjpg",
			args: args{
				name: "test.jpg",
			},
			wantErr: false,
		},
	}
	globals.Config.ControllerHostName = "192.168.21.237"
		globals.Config.ControllerAPIPort = 3003
		globals.Config.UserID = 999999
		globals.Config.DeviceID = 70000007

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uploadFile(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("uploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
