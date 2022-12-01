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
 */

package gpiorelay

import (
	"bubblesnet/edge-device/sense-go/globals"
	"reflect"
	"testing"
	"time"
)

func TestGetDispenserService(t *testing.T) {
	dispenserService := GetDispenserService()

	tests := []struct {
		name string
		want DispenserService
	}{
		{name: "happy", want: dispenserService},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDispenserService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDispenserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockDispenser_InitRpioPins(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation                    *globals.Station
		MyDevice                     *globals.EdgeDevice
		RunningOnUnsupportedHardware bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false}, args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice, RunningOnUnsupportedHardware: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.InitRpioPins(tt.args.MyStation, tt.args.MyDevice, tt.args.RunningOnUnsupportedHardware)
		})
	}
}

func TestMockDispenser_IsDispenserOn(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		name      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice,
				name: globals.DISPENSER_NAME_PH_UP}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			if got := m.IsDispenserOn(tt.args.MyStation, tt.args.MyDevice, tt.args.name); got != tt.want {
				t.Errorf("IsDispenserOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockDispenser_IsMyDispenser(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		name      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice,
				name: globals.DISPENSER_NAME_PH_UP}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MockDispenser{
				Real: tt.fields.Real,
			}
			if got := r.IsMyDispenser(tt.args.MyStation, tt.args.MyDevice, tt.args.name); got != tt.want {
				t.Errorf("IsMyDispenser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockDispenser_ReportAll(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		timeout   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false}, args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.ReportAll(tt.args.MyStation, tt.args.MyDevice, tt.args.timeout)
		})
	}
}

func TestMockDispenser_RunPinToggler(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		isTest    bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false}, args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice, isTest: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.RunPinToggler(tt.args.MyStation, tt.args.MyDevice, tt.args.isTest)
		})
	}
}

func TestMockDispenser_SendDispenserStatusChangeEvent(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		switch_name string
		on          bool
		sequence    int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false}, args: args{switch_name: globals.DISPENSER_NAME_PH_UP,
			on: true, sequence: 111}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.SendDispenserStatusChangeEvent(tt.args.switch_name, tt.args.on, tt.args.sequence)
		})
	}
}

func TestMockDispenser_SetupDispenserGPIO(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.SetupDispenserGPIO(tt.args.MyStation, tt.args.MyDevice)
		})
	}
}

func TestMockDispenser_StartDispensing(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "happy", fields: fields{Real: false},
			args:    args{MyStation: globals.MyStation, MyDevice: globals.MyDevice},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			if err := m.StartDispensing(tt.args.MyStation, tt.args.MyDevice); (err != nil) != tt.wantErr {
				t.Errorf("StartDispensing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockDispenser_TimedDispenseSynchronous(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation     *globals.Station
		MyDevice      *globals.EdgeDevice
		dispenserName string
		milliseconds  int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice,
				dispenserName: globals.DISPENSER_NAME_PH_DOWN, milliseconds: 100},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			if err := m.TimedDispenseSynchronous(tt.args.MyStation, tt.args.MyDevice, tt.args.dispenserName, tt.args.milliseconds); (err != nil) != tt.wantErr {
				t.Errorf("TimedDispenseSynchronous() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockDispenser_TurnAllOff(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		timeout   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice,
				timeout: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.TurnAllOff(tt.args.MyStation, tt.args.MyDevice, tt.args.timeout)
		})
	}
}

func TestMockDispenser_TurnOffDispenserByIndex(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{index: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.TurnOffDispenserByIndex(tt.args.index)
		})
	}
}

func TestMockDispenser_TurnOffDispenserByName(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		name      string
		force     bool
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantSomethingChanged bool
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice,
				name:  globals.DISPENSER_NAME_PH_DOWN,
				force: true},
			wantSomethingChanged: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			if gotSomethingChanged := m.TurnOffDispenserByName(tt.args.MyStation, tt.args.MyDevice, tt.args.name, tt.args.force); gotSomethingChanged != tt.wantSomethingChanged {
				t.Errorf("TurnOffDispenserByName() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
			}
		})
	}
}

func TestMockDispenser_TurnOnDispenserByIndex(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{index: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			m.TurnOnDispenserByIndex(tt.args.index)
		})
	}
}

func TestMockDispenser_TurnOnDispenserByName(t *testing.T) {
	type fields struct {
		Real bool
	}
	type args struct {
		MyStation *globals.Station
		MyDevice  *globals.EdgeDevice
		name      string
		force     bool
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantSomethingChanged bool
	}{
		{name: "happy", fields: fields{Real: false},
			args: args{MyStation: globals.MyStation, MyDevice: globals.MyDevice,
				name:  globals.DISPENSER_NAME_PH_DOWN,
				force: true},
			wantSomethingChanged: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDispenser{
				Real: tt.fields.Real,
			}
			if gotSomethingChanged := m.TurnOnDispenserByName(tt.args.MyStation, tt.args.MyDevice, tt.args.name, tt.args.force); gotSomethingChanged != tt.wantSomethingChanged {
				t.Errorf("TurnOnDispenserByName() = %v, want %v", gotSomethingChanged, tt.wantSomethingChanged)
			}
		})
	}
}
