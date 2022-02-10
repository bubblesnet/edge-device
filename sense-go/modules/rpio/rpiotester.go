package rpio

// +build linux,arm

package main

import (
"fmt"
"github.com/stianeikeland/go-rpio"
)

func main() {
	fmt.Println("Calling rpio.open")
	err := rpio.Open()
	if err != nil  {
		fmt.Println("%#v", err )
	} else {
		fmt.Println("Error is nil rpio.open worked")
	}
	defer func() {
		err := rpio.Close()
		if err != nil {
			fmt.Println("rpio.close %+v", err)
		}
	}()
}
