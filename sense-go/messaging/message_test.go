/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package messaging

import (
	"bubblesnet/edge-device/sense-go/globals"
	"reflect"
	"testing"
)

func TestNewADCSensorMessage(t *testing.T) {
	initGlobalsLocally(t)
	globals.MyDevice = &globals.EdgeDevice{DeviceID: 90000009}
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
		channel          int
		gain             int
		rate             int
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *ADCSensorMessage
	}{
		{name: "happy",
			args: args{sensor_name: "test", value: 99.99, units: globals.UNIT_VOLTS, direction: globals.Directions_up, channel: 0, gain: 1, rate: 2},
			wantPmsg: &ADCSensorMessage{
				DeviceId:  globals.MyDevice.DeviceID,
				SiteId:    1,
				StationId: 1,

				SampleTimestamp:   getNowMillis(),
				ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
				MeasurementName:   "",
				MessageType:       globals.Message_type_measurement,
				ExecutableVersion: "..  ",
				SensorName:        "test",
				Value:             99.99,
				FloatValue:        99.99,
				Units:             globals.UNIT_VOLTS,
				Direction:         globals.Directions_up,
				ChannelNumber:     0,
				Rate:              2,
				Gain:              1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPmsg := NewADCSensorMessage(tt.args.sensor_name, tt.args.measurement_name, tt.args.value, tt.args.units, tt.args.direction, tt.args.channel, tt.args.gain, tt.args.rate)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewADCSensorMessage() got  %#v", gotPmsg)
				t.Errorf("NewADCSensorMessage() want %#v", tt.wantPmsg)
			}
		})
	}
}

func TestNewDistanceSensorMessage(t *testing.T) {
	initGlobalsLocally(t)
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
		distanceCm       float64
		distanceIn       float64
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *DistanceSensorMessage
	}{
		{args: args{sensor_name: "test", measurement_name: "test_measurement", value: 99.99,
			units: globals.UNIT_VOLTS, direction: globals.Directions_up, distanceIn: 2.2, distanceCm: 2.1},
			wantPmsg: &DistanceSensorMessage{
				SiteId:            1,
				StationId:         1,
				DeviceId:          globals.MyDevice.DeviceID,
				SampleTimestamp:   getNowMillis(),
				ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
				MessageType:       globals.Message_type_measurement,
				ExecutableVersion: "..  ",
				SensorName:        "test",
				MeasurementName:   "test_measurement",
				Value:             99.99,
				FloatValue:        99.99,
				Units:             globals.UNIT_VOLTS,
				Direction:         globals.Directions_up,
				DistanceCm:        2.1,
				DistanceIn:        2.2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPmsg := NewDistanceSensorMessage(tt.args.sensor_name, tt.args.measurement_name, tt.args.value, tt.args.units, tt.args.direction, tt.args.distanceCm, tt.args.distanceIn)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewDistanceSensorMessage() got  %#v", gotPmsg)
				t.Errorf("NewDistanceSensorMessage() want %#v", tt.wantPmsg)
			}
		})
	}
}

func TestNewGenericSensorMessage(t *testing.T) {
	initGlobalsLocally(t)
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *GenericSensorMessage
	}{
		{name: "happy", args: args{sensor_name: "test", measurement_name: "test_measurement", value: 99.99,
			units: globals.UNIT_VOLTS, direction: globals.Directions_up},
			wantPmsg: &GenericSensorMessage{
				SiteId:            1,
				StationId:         1,
				DeviceId:          globals.MyDevice.DeviceID,
				SampleTimestamp:   getNowMillis(),
				ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
				MessageType:       globals.Message_type_measurement,
				ExecutableVersion: "..  ",
				SensorName:        "test",
				MeasurementName:   "test_measurement",
				Value:             99.99,
				FloatValue:        99.99,
				Units:             globals.UNIT_VOLTS,
				Direction:         globals.Directions_up,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantPmsg.SampleTimestamp = getNowMillis()
			gotPmsg := NewGenericSensorMessage(tt.args.sensor_name, tt.args.measurement_name, tt.args.value, tt.args.units, tt.args.direction)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewGenericSensorMessage() got  %#v", gotPmsg)
				t.Errorf("NewGenericSensorMessage() want %#v", tt.wantPmsg)
			}
		})
	}
}
func initGlobalsLocally(t *testing.T) {
	globals.MyDeviceID = 70000008
	if err := globals.ReadCompleteSiteFromPersistentStore("../testdata", "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
		t.Errorf("getConfigFromServer() error = %#v", err)
	}
}

func TestNewTamperSensorMessage(t *testing.T) {
	initGlobalsLocally(t)
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
		moveX            float64
		moveY            float64
		moveZ            float64
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *TamperEventMessage
	}{
		{name: "happy", args: args{sensor_name: "test", value: 99.99, units: globals.UNIT_VOLTS, direction: globals.Directions_none, measurement_name: "movement", moveX: 1.1, moveY: 2.2, moveZ: 3.3},
			wantPmsg: &TamperEventMessage{
				SiteId:            1,
				StationId:         1,
				DeviceId:          globals.MyDevice.DeviceID,
				SampleTimestamp:   getNowMillis(),
				ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
				MessageType:       "event",
				ExecutableVersion: "..  ",
				SensorName:        "test",
				MeasurementName:   "tamper",
				Value:             99.99,
				FloatValue:        99.99,
				Units:             globals.UNIT_VOLTS,
				Direction:         globals.Directions_none,
				XMove:             1.1,
				YMove:             2.2,
				ZMove:             3.3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantPmsg.SampleTimestamp = getNowMillis()
			gotPmsg := NewTamperSensorMessage(tt.args.sensor_name, tt.args.value, tt.args.units, tt.args.direction, tt.args.moveX, tt.args.moveY, tt.args.moveZ)
			tt.wantPmsg.SampleTimestamp = gotPmsg.SampleTimestamp
			if !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewTamperSensorMessage() = got  %#v", gotPmsg)
				t.Errorf("NewTamperSensorMessage() = want %#v", tt.wantPmsg)
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
				t.Errorf("getNowMillis() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewSwitchStatusChangeMessage(t *testing.T) {
	initGlobalsLocally(t)
	testmsg := SwitchStatusChangeMessage{
		DeviceId:          70000008,
		SiteId:            1,
		StationId:         1,
		ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
		ExecutableVersion: "..  ",
		EventTimestamp:    getNowMillis(),
		MessageType:       "switch_event",
		SwitchName:        "testswitch",
		On:                true,
	}
	type args struct {
		switch_name string
		on          bool
	}
	tests := []struct {
		name     string
		args     args
		wantPmsg *SwitchStatusChangeMessage
	}{
		{name: "happy", args: args{switch_name: "testswitch", on: true}, wantPmsg: &testmsg},
	}
	for _, tt := range tests {
		testmsg.EventTimestamp = getNowMillis()
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewSwitchStatusChangeMessage(tt.args.switch_name, tt.args.on); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				// This timestamp hack will cause false negatives if timestamp error coincides with other error
				if gotPmsg.EventTimestamp == tt.wantPmsg.EventTimestamp {
					t.Errorf("NewSwitchStatusChangeMessage() = got  %#v", gotPmsg)
					t.Errorf("NewSwitchStatusChangeMessage() = want %#v", tt.wantPmsg)
				}
			}
		})
	}
}

func TestNewVOCSensorMessage(t *testing.T) {
	initGlobalsLocally(t)
	//	sensor_name string, measurement_name string, value float64, units string, direction string
	testmsg := VOCSensorMessage{
		DeviceId:          70000008,
		SiteId:            1,
		StationId:         1,
		ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
		ExecutableVersion: "..  ",
		MessageType:       globals.Message_type_measurement,
		SensorName:        globals.SENSOR_NAME_VOC,
		MeasurementName:   globals.MEASUREMENT_NAME_VOC,
		Value:             1.1,
		FloatValue:        1.1,
		Units:             globals.UNITS_PARTS_PER_BILLION,
		Direction:         globals.Directions_none,
	}
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
	}

	tests := []struct {
		name     string
		args     args
		wantPmsg *VOCSensorMessage
	}{
		{name: "happy", args: args{sensor_name: globals.SENSOR_NAME_VOC,
			measurement_name: globals.MEASUREMENT_NAME_VOC, value: 1.1, units: globals.UNITS_PARTS_PER_BILLION,
			direction: globals.Directions_none}, wantPmsg: &testmsg},
	}
	for _, tt := range tests {
		testmsg.SampleTimestamp = getNowMillis()
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewVOCSensorMessage(tt.args.sensor_name, tt.args.measurement_name,
				tt.args.value, tt.args.units, tt.args.direction); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("TestNewVOCSensorMessage() = got  %#v", gotPmsg)
				t.Errorf("TestNewVOCSensorMessage() = want %#v", tt.wantPmsg)
			}
		})
	}
}

func TestNewCO2SensorMessage(t *testing.T) {
	initGlobalsLocally(t)
	//	sensor_name string, measurement_name string, value float64, units string, direction string
	testmsg := CO2SensorMessage{
		DeviceId:          70000008,
		SiteId:            1,
		StationId:         1,
		ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
		ExecutableVersion: "..  ",
		MessageType:       globals.Message_type_measurement,
		SensorName:        globals.SENSOR_NAME_CO2,
		MeasurementName:   globals.MEASUREMENT_NAME_CO2,
		Value:             1.1,
		FloatValue:        1.1,
		Units:             globals.UNITS_PARTS_PER_MILLION,
		Direction:         globals.Directions_none,
	}
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
	}

	tests := []struct {
		name     string
		args     args
		wantPmsg *CO2SensorMessage
	}{
		{name: "happy", args: args{sensor_name: globals.SENSOR_NAME_CO2,
			measurement_name: globals.MEASUREMENT_NAME_CO2, value: 1.1, units: globals.UNITS_PARTS_PER_MILLION,
			direction: globals.Directions_none}, wantPmsg: &testmsg},
	}
	for _, tt := range tests {
		testmsg.SampleTimestamp = getNowMillis()
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewCO2SensorMessage(tt.args.sensor_name, tt.args.measurement_name,
				tt.args.value, tt.args.units, tt.args.direction); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("TestNewCO2SensorMessage() = got  %#v", gotPmsg)
				t.Errorf("TestNewCO2SensorMessage() = want %#v", tt.wantPmsg)
			}
		})
	}
}

func TestNewCCS811CurrentMessage(t *testing.T) {
	initGlobalsLocally(t)
	//	sensor_name string, measurement_name string, value float64, units string, direction string
	testmsg := CCS811CurrentMessage{
		DeviceId:          70000008,
		SiteId:            1,
		StationId:         1,
		ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
		ExecutableVersion: "..  ",
		MessageType:       globals.Message_type_measurement,
		SensorName:        globals.SENSOR_NAME_CCS811_CURRENT,
		MeasurementName:   globals.MEASUREMENT_NAME_CCS811_RAW_CURRENT,
		Value:             1.1,
		FloatValue:        1.1,
		Units:             globals.UNITS_MICRO_AMPS,
		Direction:         globals.Directions_none,
	}
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
	}

	tests := []struct {
		name     string
		args     args
		wantPmsg *CCS811CurrentMessage
	}{
		{name: "happy", args: args{sensor_name: globals.SENSOR_NAME_CCS811_CURRENT,
			measurement_name: globals.MEASUREMENT_NAME_CCS811_RAW_CURRENT, value: 1.1,
			units:     globals.UNITS_MICRO_AMPS,
			direction: globals.Directions_none}, wantPmsg: &testmsg},
	}
	for _, tt := range tests {
		testmsg.SampleTimestamp = getNowMillis()
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewCCS811CurrentMessage(tt.args.sensor_name, tt.args.measurement_name,
				tt.args.value, tt.args.units, tt.args.direction); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewCCS811CurrentMessage() = got  %#v", gotPmsg)
				t.Errorf("NewCCS811CurrentMessage() = want %#v", tt.wantPmsg)
			}
		})
	}
}

func TestNewCCS811VoltageMessage(t *testing.T) {
	initGlobalsLocally(t)
	//	sensor_name string, measurement_name string, value float64, units string, direction string
	testmsg := CCS811VoltageMessage{
		DeviceId:          70000008,
		SiteId:            1,
		StationId:         1,
		ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
		ExecutableVersion: "..  ",
		MessageType:       globals.Message_type_measurement,
		SensorName:        globals.SENSOR_NAME_CCS811_VOLTAGE,
		MeasurementName:   globals.MEASUREMENT_NAME_CCS811_RAW_VOLTAGE,
		Value:             1.1,
		FloatValue:        1.1,
		Units:             globals.UNITS_MICRO_VOLTS,
		Direction:         globals.Directions_none,
	}
	type args struct {
		sensor_name      string
		measurement_name string
		value            float64
		units            string
		direction        string
	}

	tests := []struct {
		name     string
		args     args
		wantPmsg *CCS811VoltageMessage
	}{
		{name: "happy", args: args{sensor_name: globals.SENSOR_NAME_CCS811_VOLTAGE,
			measurement_name: globals.MEASUREMENT_NAME_CCS811_RAW_VOLTAGE, value: 1.1,
			units:     globals.UNITS_MICRO_VOLTS,
			direction: globals.Directions_none}, wantPmsg: &testmsg},
	}
	for _, tt := range tests {
		testmsg.SampleTimestamp = getNowMillis()
		t.Run(tt.name, func(t *testing.T) {
			if gotPmsg := NewCCS811VoltageMessage(tt.args.sensor_name, tt.args.measurement_name,
				tt.args.value, tt.args.units, tt.args.direction); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewCCS811VoltageMessage() = got  %#v", gotPmsg)
				t.Errorf("NewCCS811VoltageMessage() = want %#v", tt.wantPmsg)
			}
		})
	}
}

func TestNewPictureTakenMessage(t *testing.T) {
	initGlobalsLocally(t)
	testmsg := PictureTakenMessage{
		SiteId:            1,
		StationId:         1,
		DeviceId:          70000008,
		ContainerName:     globals.CONTAINER_NAME_SENSE_GO,
		ExecutableVersion: "..  ",
		EventTimestamp:    getNowMillis(),
		MessageType:       "picture_event",
		PictureFilename:   "blah.jpg",
	}
	tests := []struct {
		name     string
		wantPmsg *PictureTakenMessage
	}{
		{name: "happy", wantPmsg: &testmsg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testmsg.EventTimestamp = getNowMillis()
			if gotPmsg := NewPictureTakenMessage("blah.jpg", 0); !reflect.DeepEqual(gotPmsg, tt.wantPmsg) {
				t.Errorf("NewPictureTakenMessage() got  %#v", gotPmsg)
				t.Errorf("NewPictureTakenMessage() want %#v", tt.wantPmsg)
			}
		})
	}
}
