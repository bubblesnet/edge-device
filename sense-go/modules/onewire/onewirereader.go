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

package gonewire

// copyright and license inspection - 4/13/22 - added original copyright into code

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func ReadOneWire() {
	dir := "/sys/bus/w1/devices/"

	gw, err := New(dir)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigs:
				cancel()
				return
			case v := <-gw.Values():
				ftemp, _ := strconv.ParseFloat(v.Value, 64)
				ftemp = ftemp / 1000
				fahrenheit := (ftemp * 1.8000) + 32.00

				direction := ""
				if fahrenheit > float64(globals.LastWaterTemp) {
					direction = "up"
				} else if fahrenheit < float64(globals.LastWaterTemp) {
					direction = "down"
				}
				globals.LastWaterTemp = float32(fahrenheit)

				phm := messaging.NewGenericSensorMessage("thermometer_water", "temp_water", fahrenheit, "F", direction)
				bytearray, err := json.Marshal(phm)
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
				if globals.Client != nil {
					_, err = globals.Client.StoreAndForward(context.Background(), &message)
					if err != nil {
						log.Errorf("RunADCPoller ERROR %#v", err)
					} else {
						//				log.Infof("sensor_reply %#v", sensor_reply)

					}
				} else {
					_ = errors.New("GRPC client is not connected!")
				}
			}
		}
	}()

	gw.OnReadError(func(e error, s *Sensor) {
		fmt.Printf("onReadError\n")
		log.Errorf("blah")
		log.Errorf("[ERR] %s", s.ID())
		log.Errorf("[ERR] %#v", err)
	})

	gw.Start(ctx, 10*time.Second)
}
