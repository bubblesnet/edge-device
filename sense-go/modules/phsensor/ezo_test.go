//go:build (linux && arm) || arm64
// +build linux,arm arm64

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

package phsensor

import (
	"github.com/go-playground/log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"reflect"
	"testing"
)

func TestAtlasEZODriver_Connection(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name   string
		fields fields
		want   gobot.Connection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if got := d.Connection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connection() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAtlasEZODriver_Halt(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if err := d.Halt(); (err != nil) != tt.wantErr {
				t.Errorf("Halt() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func TestAtlasEZODriver_Name(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if got := d.Name(); got != tt.want {
				t.Errorf("Name() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAtlasEZODriver_Ph(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantPH  float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			gotPH, err := d.Ph()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ph() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if gotPH != tt.wantPH {
				t.Errorf("Ph() gotPH = %#v, want %#v", gotPH, tt.wantPH)
			}
		})
	}
}

func TestAtlasEZODriver_SetName(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	type args struct {
		n string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			log.Infof("d = %#v", d)
		})
	}
}

func TestAtlasEZODriver_Start(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if err := d.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func TestAtlasEZODriver_initialization(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			if err := d.initialization(); (err != nil) != tt.wantErr {
				t.Errorf("initialization() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func TestAtlasEZODriver_rawPh(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	tests := []struct {
		name    string
		fields  fields
		wantPH  float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			gotPH, err := d.rawPh()
			if (err != nil) != tt.wantErr {
				t.Errorf("rawPh() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if gotPH != tt.wantPH {
				t.Errorf("rawPh() gotPH = %#v, want %#v", gotPH, tt.wantPH)
			}
		})
	}
}

func TestAtlasEZODriver_read(t *testing.T) {
	type fields struct {
		name       string
		connector  i2c.Connector
		connection i2c.Connection
		Config     i2c.Config
		tpc        *bmp280CalibrationCoefficients
	}
	type args struct {
		address byte
		n       int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AtlasEZODriver{
				name:       tt.fields.name,
				connector:  tt.fields.connector,
				connection: tt.fields.connection,
				Config:     tt.fields.Config,
				tpc:        tt.fields.tpc,
			}
			got, err := d.read(tt.args.address, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("read() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("read() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewAtlasEZODriver(t *testing.T) {
	type args struct {
		c       i2c.Connector
		options []func(i2c.Config)
	}
	tests := []struct {
		name string
		args args
		want *AtlasEZODriver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAtlasEZODriver(tt.args.c, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAtlasEZODriver() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
