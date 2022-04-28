//go:build darwin || windows || (linux && arm)
// +build darwin windows linux,arm

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

package camera

import (
	"bubblesnet/edge-device/sense-go/globals"
	"testing"
)

func Test_uploadFile(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "goodjpg",
			args: args{
				name: "test.jpg",
			},
			wantErr: false,
		},
	}
	globals.MySite.ControllerHostName = "192.168.21.237"
	globals.MySite.ControllerAPIPort = 3003
	globals.MySite.UserID = 90000009
	globals.MyDevice = &globals.EdgeDevice{DeviceID: 70000008}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uploadFile(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("uploadFile() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}
