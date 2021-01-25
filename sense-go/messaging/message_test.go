package messaging

import (
	"bubblesnet/edge-device/sense-go/globals"
	"reflect"
	"testing"
)

func TestNewADCSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		measurement_name string
		value       float64
		units       string
		direction 	string
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
			args: args{sensor_name: "test", value: 99.99, units: "Volts", direction: "up", channel: 0, gain: 1, rate: 2},
			wantPmsg: &ADCSensorMessage{
				DeviceId: globals.DeviceId,
				SampleTimestamp: getNowMillis(),
				ContainerName:     "sense-go",
				MeasurementName:	"",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				Value:             99.99,
				Units:	"Volts",
				Direction: "up",
				ChannelNumber:     0,
				Rate:              2,
				Gain:              1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPmsg := NewADCSensorMessage(tt.args.sensor_name, tt.args.measurement_name, tt.args.value, tt.args.units, tt.args.direction, tt.args.channel, tt.args.gain, tt.args.rate);
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewADCSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func TestNewDistanceSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		measurement_name string
		value       float64
		units       string
		direction 	string
		distanceCm  float64
		distanceIn  float64
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *DistanceSensorMessage
	}{
		{args: args{sensor_name: "test", measurement_name: "test_measurement", value: 99.99, units: "Volts", direction: "up", distanceIn: 2.2, distanceCm: 2.1},
			wantPmsg: &DistanceSensorMessage{
				DeviceId: globals.DeviceId,
				SampleTimestamp:   getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				MeasurementName:	"test_measurement",
				Value:             99.99,
				Units:             "Volts",
				Direction: "up",
				DistanceCm: 2.1,
				DistanceIn: 2.2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPmsg := NewDistanceSensorMessage(tt.args.sensor_name, tt.args.measurement_name, tt.args.value, tt.args.units, tt.args.direction, tt.args.distanceCm, tt.args.distanceIn)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewDistanceSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func TestNewGenericSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		measurement_name string
		value       float64
		units       string
		direction 	string
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *GenericSensorMessage
	}{
		{name: "happy", args: args{sensor_name: "test", measurement_name: "test_measurement", value: 99.99, units: "Volts",direction: "up",},
			wantPmsg: &GenericSensorMessage {
				DeviceId: globals.DeviceId,
				SampleTimestamp:   getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				MeasurementName: 	"test_measurement",
				Value:             99.99,
				Units:             "Volts",
				Direction: "up",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantPmsg.SampleTimestamp = getNowMillis()
			gotPmsg := NewGenericSensorMessage(tt.args.sensor_name, tt.args.measurement_name, tt.args.value, tt.args.units, tt.args.direction)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewGenericSensorMessage() = %v, want %v", gotPmsg, tt.wantPmsg)
			}
		})
	}
}

func TestNewTamperSensorMessage(t *testing.T) {
	type args struct {
		sensor_name string
		measurement_name string
		value       float64
		units       string
		direction 	string
		moveX       float64
		moveY       float64
		moveZ       float64
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *TamperSensorMessage
	}{
		{name: "happy", args: args{sensor_name: "test", value: 99.99, units: "Volts", direction: "", measurement_name: "movement", moveX: 1.1, moveY: 2.2, moveZ: 3.3},
			wantPmsg: &TamperSensorMessage{
				DeviceId: globals.DeviceId,
				SampleTimestamp:   getNowMillis(),
				ContainerName:     "sense-go",
				MessageType:       "measurement",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				MeasurementName: "movement",
				Value:             99.99,
				Units:             "Volts",
				Direction: "",
				XMove: 1.1,
				YMove: 2.2,
				ZMove: 3.3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantPmsg.SampleTimestamp = getNowMillis()
			gotPmsg := NewTamperSensorMessage(tt.args.sensor_name, tt.args.value, tt.args.units, tt.args.direction, tt.args.moveX, tt.args.moveY, tt.args.moveZ)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
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
