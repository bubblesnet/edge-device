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
package a2dconverter

// copyright and license inspection - no issues 4/13/22

const MinInches = 0.0
const MaxInches = 12.5
const MaxVoltage = 2.65
const MinVoltage = 1.65

const ohmRange = 1600.0
const MinOhms = 400.0
const MaxOhms = 2400.0

const Etape_slope = 11.37795
const Etape_y_intercept = -17.28562

func etapeInchesToGallons(MaxInches float64, MaxGallons float64, inches float64) (gallons float64) {
	return inches * (MaxGallons / MaxInches)
}

func etapeInchesFromVolts(voltage float64, slope float64, yintercept float64) (inches float64) {
	// INCHES = VOLTS * -12.5 + 33.125
	inches = voltage*slope + yintercept
	return inches
}
