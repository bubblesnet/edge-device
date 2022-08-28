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

package accelerometer

// copyright and license inspection - no issues 4/13/22
// from the sample code at https://github.com/hybridgroup/gobot/blob/release/examples/firmata_adxl345.go

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"math"
	"time"
)

var singletonTamperDetectorService = RealTamperDetector{Real: true}

type RealTamperDetector struct {
	Real bool
}

func GetTamperDetectorService() TamperDetectorService {
	return &singletonTamperDetectorService
}

func (r *RealTamperDetector) RunTamperDetector(onceOnly bool) {
	log.Infof("adxl345: runTamperDetector Xmove: %f, Ymove %f, Zmove %f", globals.MyStation.TamperSpec.Xmove,
		globals.MyStation.TamperSpec.Ymove, globals.MyStation.TamperSpec.Zmove)
	adxl345Adaptor := raspi.NewAdaptor()
	adxl345 := i2c.NewADXL345Driver(adxl345Adaptor)
	lastx := 0.0
	lasty := 0.0
	lastz := 0.0

	xmove := 0.0
	ymove := 0.0
	zmove := 0.0

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			x, y, z, _ := adxl345.XYZ()
			//			log.Debugf("adxl345: x: %.7f | y: %.7f | z: %.7f \n", x, y, z))
			if lastx == 0.0 {
			} else {
				xmove = math.Abs(lastx - x)
				ymove = math.Abs(lasty - y)
				zmove = math.Abs(lastz - z)
				if xmove > globals.MyStation.TamperSpec.Xmove || ymove > globals.MyStation.TamperSpec.Ymove || zmove > globals.MyStation.TamperSpec.Zmove {
					log.Infof("adxl345: new tamper message !! x: %.3f | y: %.3f | z: %.3f ", xmove, ymove, zmove)
					var tamperMessage = messaging.NewTamperSensorMessage(globals.Sensor_name_tamper_sensor,
						0.0, "", "", xmove, ymove, zmove)
					bytearray, err := json.Marshal(tamperMessage)
					if err != nil {
						fmt.Println(err)
						return
					}
					message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: globals.Grpc_message_typeid_sensor, Data: string(bytearray)}
					_, err = globals.Client.StoreAndForward(context.Background(), &message)
					if err != nil {
						log.Errorf("adxl345: runTamperDetector ERROR %#v", err)
					} else {
						//						log.Debugf("adxl345: %#v", sensor_reply)
					}

				} else {
					// log.Debugf("adxl345: non-tamper movement - x: %.3f | y: %.3f | z: %.3f", xmove, ymove, zmove)
				}
			}
			lastx = x
			lasty = y
			lastz = z
		})
	}

	robot := gobot.NewRobot("adxl345Bot",
		[]gobot.Connection{adxl345Adaptor},
		[]gobot.Device{adxl345},
		work,
	)

	err := robot.Start()
	if err != nil {
		globals.ReportDeviceFailed("adxl345")
		log.Errorf("adxl345: robot start error %#v", err)
	}

	if onceOnly {
		robot.Stop()
	}

	if onceOnly {
		robot.Stop()
	}
}
