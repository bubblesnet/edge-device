// +build windows darwin linux,amd64

package rpio

import (
	"github.com/go-playground/log"
)

func OpenRpio() {
	log.Info("Calling rpio.open")
}