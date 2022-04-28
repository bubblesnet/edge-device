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

type ChannelConfig struct {
	gain int
	rate int
}

type AdapterConfig struct {
	bus_id            int
	address           int
	channelConfig     [4]ChannelConfig
	channelWaitMillis int
}

type ChannelValue struct {
	ChannelNumber    int     `json:"channel_number"`
	Voltage          float64 `json:"voltage,omitempty"`
	Gain             int     `json:"gain"`
	Rate             int     `json:"rate"`
	SensorName       string  `json:"sensor_name"`
	MeasurementName  string  `json:"measurement_name"`
	MeasurementUnits string  `json:"measurement_units"`
	Slope            float64 `json:"slope"`
	Yintercept       float64 `json:"yintercept"`
}

type LinearChannelTranslations struct {
	SensorName       string  `json:"sensor_name"`
	MeasurementName  string  `json:"measurement_name"`
	MeasurementUnits string  `json:"measurement_units"`
	Slope            float64 `json:"slope"`
	Yintercept       float64 `json:"yintercept"`
}

type Channels [4]ChannelValue
type Translations [4]LinearChannelTranslations

var a0 = AdapterConfig{
	bus_id:  1,
	address: 0x48,
	channelConfig: [4]ChannelConfig{
		{gain: 1,
			rate: 860},
		{gain: 1,
			rate: 860},
		{gain: 1,
			rate: 860},
		{gain: 1,
			rate: 860},
	},
}

var a1 = AdapterConfig{
	bus_id:  1,
	address: 0x49,
	channelConfig: [4]ChannelConfig{
		{gain: 1,
			rate: 860},
		{gain: 1,
			rate: 860},
		{gain: 1,
			rate: 860},
		{gain: 1,
			rate: 860},
	},
}
var daps = []AdapterConfig{a0, a1}

type ADCMessage struct {
	BusId         int      `json:"bus_id"`
	Address       int      `json:"address"`
	ChannelValues Channels `json:"channel_values"`
}
