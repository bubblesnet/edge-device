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

package co2vocmeter

import (
	"bubblesnet/edge-device/sense-go/messaging"
	"periph.io/x/periph/experimental/devices/ccs811"
	"reflect"
	"testing"
)

func TestReadCO2VOC(t *testing.T) {
	type args struct {
		ptemp     *float32
		phumidity *float32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReadCO2VOC(tt.args.ptemp, tt.args.phumidity)
		})
	}
}

func Test_checkErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkErr(tt.args.err)
		})
	}
}

func Test_getCCCS811SensorMessages(t *testing.T) {
	type args struct {
		sensorValues ccs811.SensorValues
	}
	tests := []struct {
		name              string
		args              args
		wantCo2msg        *messaging.CO2SensorMessage
		wantVocmsg        *messaging.VOCSensorMessage
		wantRawcurrentmsg *messaging.CCS811CurrentMessage
		wantRawvoltagemsg *messaging.CCS811VoltageMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCo2msg, gotVocmsg, gotRawcurrentmsg, gotRawvoltagemsg := getCCS811SensorMessages(tt.args.sensorValues)
			if !reflect.DeepEqual(gotCo2msg, tt.wantCo2msg) {
				t.Errorf("getCCS811SensorMessages() gotCo2msg = %v, want %v", gotCo2msg, tt.wantCo2msg)
			}
			if !reflect.DeepEqual(gotVocmsg, tt.wantVocmsg) {
				t.Errorf("getCCS811SensorMessages() gotVocmsg = %v, want %v", gotVocmsg, tt.wantVocmsg)
			}
			if !reflect.DeepEqual(gotRawcurrentmsg, tt.wantRawcurrentmsg) {
				t.Errorf("getCCS811SensorMessages() gotRawcurrentmsg = %v, want %v", gotRawcurrentmsg, tt.wantRawcurrentmsg)
			}
			if !reflect.DeepEqual(gotRawvoltagemsg, tt.wantRawvoltagemsg) {
				t.Errorf("getCCS811SensorMessages() gotRawvoltagemsg = %v, want %v", gotRawvoltagemsg, tt.wantRawvoltagemsg)
			}
		})
	}
}

func Test_loadBaseline(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadBaseline(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadBaseline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reportStatus(t *testing.T) {
	type args struct {
		status byte
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "happy", args: args{status: STATUS_MASK_APPVALID}},
		{name: "happy", args: args{status: STATUS_MASK_APPERASE}},
		{name: "happy", args: args{status: STATUS_MASK_APPVERIFY}},
		{name: "happy", args: args{status: STATUS_MASK_DATAREADY}},
		{name: "happy", args: args{status: STATUS_MASK_ERROR}},
		{name: "happy", args: args{status: STATUS_MASK_FIRMWAREMODE}},
		{name: "happy", args: args{status: STATUS_MASK_RESERVED1}},
		{name: "happy", args: args{status: STATUS_MASK_RESERVED2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := reportStatus(tt.args.status)
			if len(str) == 0 {
				t.Error("bad status string")
			}
		})
	}
}

func Test_saveBaseline(t *testing.T) {
	type args struct {
		baseline []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveBaseline(tt.args.baseline)
		})
	}
}

func Test_toCelsius(t *testing.T) {
	type args struct {
		fahrenheit float32
	}
	tests := []struct {
		name        string
		args        args
		wantCelsius float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCelsius := toCelsius(tt.args.fahrenheit); gotCelsius != tt.wantCelsius {
				t.Errorf("toCelsius() = %v, want %v", gotCelsius, tt.wantCelsius)
			}
		})
	}
}
