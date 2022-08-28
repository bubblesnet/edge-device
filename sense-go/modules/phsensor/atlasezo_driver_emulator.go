//go:build darwin || (windows && amd64) || (linux && amd64)
// +build darwin windows,amd64 linux,amd64

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

/// TODO copyright and license inspection - 4/13/22 - recheck - figure out where this code came from

import (
	"fmt"
	"gobot.io/x/gobot/drivers/i2c"
)

type bmp280CalibrationCoefficients struct {
	t1 uint16
	t2 int16
	t3 int16
	p1 uint16
	p2 int16
	p3 int16
	p4 int16
	p5 int16
	p6 int16
	p7 int16
	p8 int16
	p9 int16
}

func clen(b []byte) (n int) {
	return (0)
}

// AtlasEZODriver is a driver for the BMP280 temperature/pressure sensor
type AtlasEZODriver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config

	tpc            *bmp280CalibrationCoefficients
	Connection     func() (err error)
	Name           func() (name string)
	Halt           func() error
	rawPh          func() (float64, error)
	Ph             func() (float64, error)
	initialization func() error
	read           func(byte, int) ([]byte, error)
}

func connection() (err error) {
	return nil
}

func NewAtlasEZODriver(c i2c.Connector, options ...func(i2c.Config)) *AtlasEZODriver {
	driver := AtlasEZODriver{
		name: "test",
		read: func(address byte, n int) ([]byte, error) {
			return []byte{}, nil
		},
		initialization: func() error {
			return nil
		},
		rawPh: func() (float64, error) {
			return 0, nil
		},
		Ph: func() (float64, error) {
			return 0, nil
		},
		Halt: func() error {
			fmt.Printf("Halt")
			return nil
		},
		Name: func() string {
			return "test"
		},
		Connection: func() (err error) {
			fmt.Printf("Connection")
			return nil
		},
	}
	return &driver
}

func (d *AtlasEZODriver) Start() (err error) {
	return nil
}
