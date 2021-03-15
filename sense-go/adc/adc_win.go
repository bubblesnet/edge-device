// +build darwin windows,amd64


package adc

import "gobot.io/x/gobot/drivers/i2c"

func RunADCPoller() (err error) {
	return nil
}

func readAllChannels(ads1115 *i2c.ADS1x15Driver, config AdapterConfig, adcMessage *ADCMessage) ( err error ) {
	return nil
}