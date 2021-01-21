package messaging

import (
	"reflect"
	"testing"
)

func TestNewADCSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		value       float64
		units       string
		channel     int
		gain        int
		rate        int
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *ADCSensorMessage
	}{
		{name: "happy",
			args: args{sensor_name: "test", value: 99.99, units: "Volts", channel: 0, gain: 1, rate: 2},
			wantPmsg: &ADCSensorMessage{
				SampleTimestamp: getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				Value:             99.99,
				Units:	"Volts",
				ChannelNumber:     0,
				Rate:              2,
				Gain:              1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewADCSensorMessage(tt.args.sensor_name, tt.args.value, tt.args.units, tt.args.channel, tt.args.gain, tt.args.rate); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewADCSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func TestNewDistanceSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		value       float64
		units       string
		distanceCm  float64
		distanceIn  float64
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *DistanceSensorMessage
	}{
		{args: args{sensor_name: "test", value: 99.99, units: "Volts", distanceIn: 2.2, distanceCm: 2.1},
			wantPmsg: &DistanceSensorMessage{
				SampleTimestamp:   getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				Value:             99.99,
				Units:             "Volts",
				DistanceCm: 2.1,
				DistanceIn: 2.2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewDistanceSensorMessage(tt.args.sensor_name, tt.args.value, tt.args.units, tt.args.distanceCm, tt.args.distanceIn); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewDistanceSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func TestNewGenericSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		value       float64
		units       string
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *GenericSensorMessage
	}{
		{name: "happy", args: args{sensor_name: "test", value: 99.99, units: "Volts",},
			wantPmsg: &GenericSensorMessage {
				SampleTimestamp:   getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				Value:             99.99,
				Units:             "Volts",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantPmsg.SampleTimestamp = getNowMillis()
			if gotPmsg := NewGenericSensorMessage(tt.args.sensor_name, tt.args.value, tt.args.units); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewGenericSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func TestNewTamperSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		value       float64
		units       string
		moveX       float64
		moveY       float64
		moveZ       float64
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *TamperSensorMessage
	}{
		{name: "happy", args: args{sensor_name: "test", value: 99.99, units: "Volts", moveX: 1.1, moveY: 2.2, moveZ: 3.3},
			wantPmsg: &TamperSensorMessage{
				SampleTimestamp:   getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				Value:             99.99,
				Units:             "Volts",
				XMove: 1.1,
				YMove: 2.2,
				ZMove: 3.3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantPmsg.SampleTimestamp = getNowMillis()
			if gotPmsg := NewTamperSensorMessage(tt.args.sensor_name, tt.args.value, tt.args.units, tt.args.moveX, tt.args.moveY, tt.args.moveZ); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewTamperSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func Test_getNowMillis(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNowMillis(); got != tt.want {
				t.Errorf("getNowMillis() = %v, want %v", got, tt.want)
			}
		})
	}
}
