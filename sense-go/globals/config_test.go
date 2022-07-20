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

package globals

import "testing"

func TestReadMyDeviceId(t *testing.T) {
	type args struct {
		storeMountPoint string
		relativePath    string
		fileName        string
	}
	tests := []struct {
		name    string
		args    args
		wantId  int64
		wantErr bool
	}{
		{
			name: "happy",
			args: args{
				storeMountPoint: "../testdata",
				relativePath:    "",
				fileName:        "deviceid",
			},
			wantId:  70000008,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := ReadMyDeviceId()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMyDeviceId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("ReadMyDeviceId() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestReadServerHostname(t *testing.T) {
	type args struct {
		storeMountPoint string
		relativePath    string
		fileName        string
	}
	tests := []struct {
		name         string
		args         args
		wantHostname string
		wantErr      bool
	}{
		{
			name: "happy",
			args: args{
				storeMountPoint: "../testdata",
				relativePath:    "",
				fileName:        "hostname",
			},
			wantHostname: "192.168.21.237",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostname, err := ReadMyAPIServerHostname()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMyDeviceId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHostname != tt.wantHostname {
				t.Errorf("ReadMyDeviceId() gotId = %v, want %v", gotHostname, tt.wantHostname)
			}
		})
	}
}

func initGlobalsLocally(t *testing.T) {
	MyDeviceID = 70000008
	if err := ReadCompleteSiteFromPersistentStore("../testdata", "", "config.json", &MySite, &CurrentStageSchedule); err != nil {
		t.Errorf("getConfigFromServer() error = %#v", err)
	}
}

func TestValidateConfigurable(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigurable(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigurable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	initGlobalsLocally(t)
	for _, tt := range tests {
		tt.wantErr = false
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigurable(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigurable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfigured(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigured("test"); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigured() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	initGlobalsLocally(t)
	for _, tt := range tests {
		tt.wantErr = false
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfigured("test"); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigured() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfigFromServer(t *testing.T) {
	type args struct {
		storeMountPoint string
		relativePath    string
		fileName        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "sad", wantErr: true, args: args{
			storeMountPoint: "./testdata",
			relativePath:    "",
			fileName:        "retreivedConfig.json",
		}},
		{name: "less sad", wantErr: true, args: args{
			storeMountPoint: "./testdata",
			relativePath:    "",
			fileName:        "retreivedConfig.json",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfigFromServer(tt.args.storeMountPoint, tt.args.relativePath, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("GetConfigFromServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		initGlobalsLocally(t)
	}
}
