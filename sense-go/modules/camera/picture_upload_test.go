//go:build darwin || windows || (linux && arm)
// +build darwin windows linux,arm

package camera

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
	globals.MySite.ControllerHostName = "192.168.21.237"
	globals.MySite.ControllerAPIPort = 3003
	globals.MySite.UserID = 90000009
	globals.MyDevice = &globals.EdgeDevice{DeviceID: 70000008}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uploadFile(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("uploadFile() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}
