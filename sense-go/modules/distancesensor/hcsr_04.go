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

package distancesensor

// copyright and license inspection - no issues 4/13/22

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"github.com/go-playground/log"
	hc "github.com/jdevelop/golang-rpi-extras/sensor_hcsr04"
	"golang.org/x/net/context"
	"time"
)

var lastDistance = float64(0.0)
var h hc.HCSR04

func RunDistanceWatcher(once_only bool, isUnitTest bool) {
	echoPin := 20
	pingPin := 21
	initHCSR04(echoPin, pingPin)
	RunDistanceWatcher1(once_only, isUnitTest)
}

func initHCSR04(echoPin int, pingPin int) {
	if globals.RunningOnUnsupportedHardware() {
		return
	}
	// Use BCM pin numbering
	// Echo pin
	// Trigger pin
	h = hc.NewHCSR04(echoPin, pingPin)
}

func RunDistanceWatcher1(once_only bool, isUnitTest bool) {
	log.Info("runDistanceWatcher")

	for true {
		distance := h.MeasureDistance()
		nanos := distance * 58000.00
		seconds := nanos / 1000000000.0
		mydistance := (float64)(17150.00 * seconds)
		direction := globals.Directions_none
		if mydistance > lastDistance {
			direction = globals.Directions_up
		} else if mydistance < lastDistance {
			direction = globals.Directions_down
		}
		lastDistance = mydistance
		//		log.Debugf("%.2f inches %.2f distance %.2f nanos %.2f cm\n", distance/2.54, distance, nanos, mydistance))
		dm := messaging.NewDistanceSensorMessage(globals.Sensor_name_height_sensor, globals.Measurement_name_plant_height, mydistance, globals.Distance_units_centimeters, direction, mydistance, mydistance/2.54, messaging.GetNowMillis())
		bytearray, err := json.Marshal(dm)
		if err == nil {
			log.Debugf("sending distance msg %s?", string(bytearray))
			message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: globals.Grpc_message_typeid_sensor, Data: string(bytearray)}
			if !isUnitTest {
				_, err := globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("runDistanceWatcher ERROR %#v", err)
				} else {
					//				log.Debugf("%#v", sensor_reply)
				}
			}
		} else {
			globals.ReportDeviceFailed("hcsr04")
			log.Errorf("rundistancewatcher error = %#v", err)
			break
		}
		if once_only {
			return
		}
		//		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
	}
}
