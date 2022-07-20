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

import (
	"context"
	"fmt"
	"os"
	"time"
)

// copyright and license inspection - 4/13/22 - added original copyright into code

const (
	defaultDirectory = "/sys/bus/w1/devices/"
	minFrequency     = 5000 * time.Millisecond
	chanSize         = 8
)

type Value struct {
	ID    string
	Value string
	Type  string
}

type Gonewire struct {
	valueChannel chan Value
	directory    string
	sensormap    map[string]*Sensor
	errCballback func(err error, sensor *Sensor)
}

func New(dir string) (*Gonewire, error) {
	if dir == "" {
		dir = defaultDirectory
	}

	gw := &Gonewire{
		directory:    dir,
		valueChannel: make(chan Value, chanSize),
		sensormap:    make(map[string]*Sensor),
		errCballback: defaultErrorCallback,
	}
	if err := gw.readFolder(); err != nil {
		return nil, err
	}
	return gw, nil
}

func (gw *Gonewire) Values() chan Value {
	return gw.valueChannel
}

func (gw *Gonewire) Start(ctx context.Context, frequency time.Duration) {
	if frequency < minFrequency {
		frequency = minFrequency
	}
	for {
		select {
		case <-time.After(frequency):
			for _, sensor := range gw.sensormap {
				if err := sensor.parseValue(); err != nil {
					gw.errCballback(err, sensor)
					continue
				}
				fmt.Printf("Curentvalue = %s\n", sensor.currentValue)
				gw.valueChannel <- Value{
					ID:    sensor.id,
					Value: sensor.currentValue,
					Type:  sensor.typeString,
				}
			}
		case <-ctx.Done():
			for _, sensor := range gw.sensormap {
				sensor.close()
			}
			return
		}
	}
}

func (gw *Gonewire) OnReadError(fn func(error, *Sensor)) {
	gw.errCballback = fn
}

func (gw *Gonewire) readFolder() error {
	dir, err := os.Open(gw.directory)
	if err != nil {
		return fmt.Errorf("could not read dir: %w", err)
	}
	defer dir.Close()

	subdirs, err := dir.Readdirnames(0)
	if err != nil {
		return fmt.Errorf("could not read sub dirs: %w", err)
	}

	for _, subdir := range subdirs {
		if subdir == "w1_bus_master1" {
			continue
		}
		s, err := newSensor(gw.directory, subdir)
		if err != nil {
			fmt.Println("could not init sensor: %w", err)
			continue
		}
		gw.sensormap[s.id] = s
	}

	return nil
}

func defaultErrorCallback(err error, sensor *Sensor) {
	return
}
