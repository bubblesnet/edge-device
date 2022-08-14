//go:build darwin || windows || (linux && amd64)
// +build darwin windows linux,amd64

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

package gpiorelay

import (
	"bubblesnet/edge-device/sense-go/globals"
	"github.com/go-playground/log"
	"time"
)

type MockPowerStrip struct {
	Real bool
}

var singletonPowerstrip = MockPowerStrip{Real: true}

func GetPowerstripService() PowerstripService {
	return &singletonPowerstrip
}

func (r *MockPowerStrip) IsMySwitch(MyDevice *globals.EdgeDevice, switchName string) bool {
	return true
}

func (m *MockPowerStrip) SendSwitchStatusChangeEvent(switch_name string, on bool, sequence int32) {
	log.Infof("Reporting switch %s status %#v", switch_name, on)
}

func (m *MockPowerStrip) InitRpioPins(MyDevice *globals.EdgeDevice, RunningOnUnsupportedHardware bool) {
}

func (m *MockPowerStrip) TurnAllOn(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	log.Info("Toggling all pins ON")
}

func (m *MockPowerStrip) TurnOffOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (somethingChanged bool) {
	return (false)
}

func (m *MockPowerStrip) IsOutletOn(MyDevice *globals.EdgeDevice, name string) bool {
	return false
}

func (m *MockPowerStrip) TurnOnOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (somethingChanged bool) {
	return (false)
}

func (m *MockPowerStrip) ReportAll(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("Reporting ALL")
}

func (m *MockPowerStrip) TurnAllOff(MyDevice *globals.EdgeDevice, timeout time.Duration) {
	print("Toggling pins OFF")
}

func (m *MockPowerStrip) TurnOnOutletByIndex(index int) {
}

func (m *MockPowerStrip) TurnOffOutletByIndex(index int) {
}

func (m *MockPowerStrip) RunPinToggler(MyDevice *globals.EdgeDevice, isTest bool) {
}
