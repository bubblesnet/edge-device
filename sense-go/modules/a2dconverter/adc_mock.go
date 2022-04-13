//go:build darwin || (windows && amd64) || (linux && amd64)
// +build darwin windows,amd64 linux,amd64

package a2dconverter

// copyright and license inspection - no issues 4/13/22

import "fmt"

func RunADCPoller(onceOnly bool, waitInSeconds int) (err error) {
	fmt.Printf("mock RunADCPoller")
	return nil
}

func ReadAllChannels(index int, adcMessage *ADCMessage) (err error) {
	fmt.Printf("mock ReadAllChannels")
	return nil
}
