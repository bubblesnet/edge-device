// +build linux,arm windows,amd64

package adc

import (
	"testing"
)

func TestRunADCPoller(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{ name: "happy", wantErr: false},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
//			if err := RunADCPoller(); (err != nil) != tt.wantErr {
//				t.Errorf("RunADCPoller() error = %v, wantErr %v", err, tt.wantErr)
//			}
		})
	}
}

func Test_ReadAllChannels(t *testing.T) {
	type args struct {
		config     AdapterConfig
		adcMessage *ADCMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadAllChannels(0, tt.args.adcMessage); (err != nil) != tt.wantErr {
				t.Errorf("readAllChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunADCPoller1(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunADCPoller(true); (err != nil) != tt.wantErr {
				t.Errorf("RunADCPoller() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readAllChannels1(t *testing.T) {
	type args struct {
		adcMessage *ADCMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadAllChannels(1, tt.args.adcMessage); (err != nil) != tt.wantErr {
				t.Errorf("readAllChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}