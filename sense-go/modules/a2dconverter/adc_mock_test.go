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

package a2dconverter

import "testing"

func TestReadAllChannels(t *testing.T) {
	var adcmessage = ADCMessage{
		BusId:         0,
		Address:       0,
		ChannelValues: Channels{},
	}
	type args struct {
		index      int
		adcMessage *ADCMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "happy", args: args{index: 0, adcMessage: &adcmessage}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadAllChannels(tt.args.index, tt.args.adcMessage); (err != nil) != tt.wantErr {
				t.Errorf("ReadAllChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunADCPoller(t *testing.T) {
	type args struct {
		onceOnly      bool
		waitInSeconds int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "happy",
			args:    args{onceOnly: true, waitInSeconds: 10},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunADCPoller(tt.args.onceOnly, tt.args.waitInSeconds); (err != nil) != tt.wantErr {
				t.Errorf("RunADCPoller() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}