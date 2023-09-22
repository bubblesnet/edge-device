//go:build (linux && arm) || arm64
// +build linux,arm arm64

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

// copyright and license inspection - no issues 4/13/22

import (
	"bubblesnet/edge-device/sense-go/globals"
	"fmt"
	"github.com/dhowden/raspicam"
	"github.com/go-playground/log"
	"os"
	"time"
)

func IsEnoughLightForCamera() bool {
	if globals.ExternalCurrentState.LightInternal < 2.0 {
		return false
	}
	return true
}

func WaitForLightToRegister() bool {
	log.Debugf("pictureTaker WaitForLightToRegister")
	// Wait through 3 sensor iterations if necessary
	for i := 0; i < 3*10; i++ {
		if IsEnoughLightForCamera() {
			if i > 0 {
				log.Infof("pictureTaker WaitForLightToRegister returns true for light off to on after %d waits", i)
				return true
			}
			log.Debugf("pictureTaker WaitForLightToRegister returns false for light already on")
			return false
		}
		log.Debugf("pictureTaker sleeping %d seconds for light to register", globals.MyDevice.TimeBetweenSensorPollingInSeconds/10)
		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds/10) * time.Second)
	}
	log.Info("pictureTaker WaitForLightToRegister returns false for light STILL OFF!!!!")
	return false
}

func TakeAPicture() {

	// Is there any light?  - check the bh1750
	if globals.MyStation.LightSensorInternal && !IsEnoughLightForCamera() {
		log.Infof("LightInternal value is available AND too low to take picture %f", globals.ExternalCurrentState.LightInternal)
		return
	}

	log.Infof("takeAPicture() with light reading %f", globals.ExternalCurrentState.LightInternal)
	t := time.Now()
	filename := fmt.Sprintf("%8.8d_%8.8d_%4.4d%2.2d%2.2d_%2.2d%2.2d_%2.2d.jpg", globals.MySite.UserID, globals.MyDevice.DeviceID, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	//	log.Debugf("Creating file %s", filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file: %#v", err)

		return
	}
	defer f.Close()

	//	log.Debugf("NewStill")
	s := raspicam.NewStill()
	errCh := make(chan error)
	go func() {
		for x := range errCh {
			log.Debugf("CAPTURE ERROR %#v", x)
		}
	}()
	//	log.Debugf("Capturing image...")
	raspicam.Capture(s, f, errCh)
	log.Debugf("Uploading picture %s", f.Name())
	err = uploadFile(f.Name())
	if err != nil {
		log.Errorf("os.Upload failed for %s", f.Name())
	}
	err = os.Remove(f.Name())
	if err != nil {
		log.Errorf("os.Remove failed for %s", f.Name())
	}
	SendPictureTakenEvent(filename, t.UnixMilli())

}
