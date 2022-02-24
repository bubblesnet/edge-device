// +build linux,arm arm64

package rpio

import (
	"fmt"
	"github.com/go-playground/log"
	"github.com/stianeikeland/go-rpio/v4"
)

func OpenRpio() {
	log.Info("Calling rpio.open")
	err := rpio.Open()
	if err != nil {
		log.Errorf("open rpio error %#v", err)
	}
}

func CloseRpio() {
	fmt.Printf("Called rpio.close - figure out defer bozo")
	err := rpio.Close()
	if err != nil {
		log.Errorf("rpio.close %+v", err)
	}
}
