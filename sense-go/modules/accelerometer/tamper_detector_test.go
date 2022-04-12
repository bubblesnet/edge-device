//go:build (linux && arm) || arm64
// +build linux,arm arm64

package accelerometer

import "testing"

func Test_didWeMove(t *testing.T) {
	type args struct {
		x          int32
		y          int32
		z          int32
		isUnitTest bool
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{x: 10.0, y: 10.0, z: 10.0, isUnitTest: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DidWeMove(tt.args.x, tt.args.y, tt.args.z, tt.args.isUnitTest)
		})
	}
}
