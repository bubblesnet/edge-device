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

package phsensor

// copyright and license inspection - no issues 4/13/22

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"errors"
	"github.com/go-playground/log"
	"gobot.io/x/gobot/platforms/raspi"
	"golang.org/x/net/context"
	"time"
)

func StartEzoDriver() {
	log.Info("ezoph: Starting Atlas EZO driver")
	ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
	err := ezoDriver.Start()
	if err != nil {
		globals.ReportDeviceFailed("ezoph")
		log.Errorf("ezoph: ezo start error %#v", err)
	}
}

func StartEzo(once_only bool) {
	log.Info("ezoph: RootPhSensor should be connected to this device, starting EZO reader")
	go func() {
		if err := ReadPh(once_only); err != nil {
			log.Errorf("ezoph: ReadPh %+v", err)
		}
	}()
}

var lastPh = float64(0.0)

var calibrationAdjustment = 0.0

func applyCalibration(raw float64) (calibrated float64) {
	return (raw + calibrationAdjustment)
}

func ReadPh(once_only bool) error {
	ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
	err := ezoDriver.Start()
	if err != nil {
		log.Errorf("ezoph: ezoDriver.Start returned ph device error %#v", err)
		return err
	}
	var e error = nil

	for {
		ph, err := ezoDriver.Ph()
		if err != nil {
			log.Errorf("ezoph: ReadPh error %#v", err)
			e = err
			break
		} else {
			ph = applyCalibration(ph)
			direction := globals.Directions_none
			if ph > lastPh {
				direction = globals.Directions_up
			} else if ph < lastPh {
				direction = globals.Directions_down
			}
			lastPh = ph
			phm := messaging.NewGenericSensorMessage(globals.Sensor_name_root_ph_sensor, globals.Measurement_name_root_ph, ph, globals.Ph_units_default, direction, messaging.GetNowMillis())
			bytearray, err := json.Marshal(phm)
			message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: globals.Grpc_message_typeid_sensor, Data: string(bytearray)}
			if globals.Client != nil {
				_, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("ezoph: RunADCPoller ERROR %#v", err)
				} else {
					//				log.Infof("sensor_reply %#v", sensor_reply)

				}
			} else {
				e = errors.New("ezoph: GRPC client is not connected!")
			}
		}
		if once_only {
			break
		}
		//		x := globals.MyDevice.TimeBetweenSensorPollingInSeconds

		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
	}
	log.Debugf("ezoph: returning %#v from readph", e)
	return e
}
