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

// copyright and license inspection - no issues 4/13/22

import (
	"bubblesnet/edge-device/sense-go/globals"
	"time"
)

var PowerstripSvc PowerstripService = GetPowerstripService()

type PowerstripService interface {
	SendSwitchStatusChangeEvent(switch_name string, on bool, sequence int32)
	InitRpioPins(MyDevice *globals.EdgeDevice, RunningOnUnsupportedHardware bool)
	TurnAllOn(MyDevice *globals.EdgeDevice, timeout time.Duration)
	TurnOffOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (stateChanged bool)
	IsOutletOn(MyDevice *globals.EdgeDevice, name string) bool
	TurnOnOutletByName(MyDevice *globals.EdgeDevice, name string, force bool) (stateChanged bool)
	ReportAll(MyDevice *globals.EdgeDevice, timeout time.Duration)
	TurnAllOff(MyDevice *globals.EdgeDevice, timeout time.Duration)
	TurnOnOutletByIndex(index int)
	TurnOffOutletByIndex(index int)
	RunPinToggler(MyDevice *globals.EdgeDevice, isTest bool)
	IsMySwitch(MyDevice *globals.EdgeDevice, switchName string) bool
}
