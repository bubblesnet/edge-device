package distancesensor

import "testing"

func TestRunDistanceWatcher(t *testing.T) {
	type args struct {
		once_only  bool
		isUnitTest bool
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{once_only: true, isUnitTest: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunDistanceWatcher(tt.args.once_only, tt.args.isUnitTest)
		})
	}
}
