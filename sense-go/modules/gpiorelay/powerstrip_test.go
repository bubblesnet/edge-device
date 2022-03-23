package gpiorelay

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/rpio"
	"fmt"
	"testing"
	"time"
)

func ginit() {
	rpio.OpenRpio()
	globals.MyDevice = &globals.EdgeDevice{ACOutlets: []globals.ACOutlet{{}, {}, {}, {}, {}, {}, {}, {}, {}}}
	// globals.MyDevice.ACOutlets = [8]globals.ACOutlet{}
	for i := 0; i < 8; i++ {
		globals.MyDevice.ACOutlets[i].Name = fmt.Sprintf("test%d", i)
		globals.MyDevice.ACOutlets[i].BCMPinNumber = 17
	}

}

func TestInitRpioPins(t *testing.T) {
	ginit()
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	//	globals.Config.ACOutlets = [8]globals.ACOutlet{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.InitRpioPins(globals.MyDevice, false)
		})
	}
}

func TestTurnAllOff(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.TurnAllOff(globals.MyDevice, tt.args.timeout)
		})
	}
}

func TestSendSwitchStatusChangeEvent(t *testing.T) {
	type args struct {
		switchName string
		on         bool
		sequence   int32
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{switchName: "heater", on: true, sequence: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.SendSwitchStatusChangeEvent(tt.args.switchName, tt.args.on, tt.args.sequence)
		})
	}
}

func TestTurnAllOn(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{timeout: 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.TurnAllOn(globals.MyDevice, tt.args.timeout)
		})
	}
}

func TestTurnOffOutlet(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.TurnOffOutlet(tt.args.index)
		})
	}
}

func TestTurnOffOutletByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{"blah"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.TurnOffOutletByName(globals.MyDevice, tt.args.name, false)
		})
	}
}

func TestTurnOnOutlet(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.TurnOnOutlet(tt.args.index)
		})
	}
}

func TestTurnOnOutletByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{"blah"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.TurnOnOutletByName(globals.MyDevice, tt.args.name, false)
		})
	}
}

func Test_isOutletOn(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "happy", args: args{"blah"},
			want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PowerstripSvc.isOutletOn(globals.MyDevice, tt.args.name); got != tt.want {
				t.Errorf("isOutletOn() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_runPinToggler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "happy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PowerstripSvc.RunPinToggler(globals.MyDevice, true)
		})
	}
}
