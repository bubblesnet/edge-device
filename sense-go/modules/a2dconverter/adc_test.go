// +build linux,arm arm64 windows,amd64

package a2dconverter

import (
	"github.com/go-playground/log"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"testing"
	"time"
)

func TestRunADCPoller(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "happy", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//			if err := RunADCPoller(); (err != nil) != tt.wantErr {
			//				t.Errorf("RunADCPoller() error = %#v, wantErr %#v", err, tt.wantErr)
			//			}
		})
	}
}

func initIt() (err error) {
	adcAdaptor := raspi.NewAdaptor() // optional bus/address

	ads1115s[0] = i2c.NewADS1115Driver(adcAdaptor,
		i2c.WithBus(a0.bus_id),
		i2c.WithAddress(a0.address))
	err = ads1115s[0].Start()
	if err != nil {
		log.Errorf("error starting interface %#v", err)
		return err
	}

	ads1115s[1] = i2c.NewADS1115Driver(adcAdaptor,
		i2c.WithBus(a1.bus_id),
		i2c.WithAddress(a1.address))
	err = ads1115s[1].Start()
	if err != nil {
		log.Errorf("error starting interface %#v", err)
		return err
	}
	return nil
}

func _Loop(t *testing.T) {
	for {
		Test_ReadAllChannels(t)
		time.Sleep(1 * time.Second)
	}
}

func Test_ReadAllChannels(t *testing.T) {
	adcM := ADCMessage{}

	type args struct {
		boardIndex int
		adcMessage *ADCMessage
	}
	args0 := args{boardIndex: 0, adcMessage: &adcM}
	args1 := args{boardIndex: 1, adcMessage: &adcM}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Board 0", args: args0, wantErr: false},
		{name: "Board 1", args: args1, wantErr: false},
	}
	initIt()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadAllChannels(tt.args.boardIndex, tt.args.adcMessage); (err != nil) != tt.wantErr {
				t.Errorf("readAllChannels() error = %#v, wantErr %#v", err, tt.wantErr)
			} else {
				t.Logf("Board %d %f/%f/%f/%f", tt.args.boardIndex,
					tt.args.adcMessage.ChannelValues[0].Voltage,
					tt.args.adcMessage.ChannelValues[1].Voltage,
					tt.args.adcMessage.ChannelValues[2].Voltage,
					tt.args.adcMessage.ChannelValues[3].Voltage)
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
			if err := RunADCPoller(true, 10); (err != nil) != tt.wantErr {
				t.Errorf("RunADCPoller() error = %#v, wantErr %#v", err, tt.wantErr)
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
				t.Errorf("readAllChannels() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}
