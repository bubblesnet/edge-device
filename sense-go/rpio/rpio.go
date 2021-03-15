// +build linux,arm

package rpio

import (
	"github.com/go-playground/log"
	"github.com/stianeikeland/go-rpio"
)

func OpenRpio() {
	log.Info("Calling rpio.open")
	_ = rpio.Open()
	defer func() {
		err := rpio.Close()
		if err != nil {
			log.Errorf("rpio.close %+v", err)
		}
	}()
}