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
package main

const TEMPNOTSET float64 = -100.5
const HUMIDITYNOTSET float64 = -100.5
const PRESSURENOTSET float64 = -100.5
const LIGHTNOTSET float64 = -100.5
const HEIGHTNOTSET float64 = -1

type externalState struct {
	PlantHeightIn    float64 `json:"plant_height_in,omitempty"`
	WaterTempF       float64 `json:"water_tempF,omitempty"`
	TempF            float64 `json:"tempF,omitempty"`
	ExternalTempF    float64 `json:"externalTempF,omitempty"`
	Humidity         float64 `json:"humidity,omitempty"`
	PressureInternal float64 `json:"pressure_internal,omitempty"`
	LightInternal    float64 `json:"light_internal,omitempty"`
}

var ExternalCurrentState = externalState{
	PlantHeightIn:    HEIGHTNOTSET,
	WaterTempF:       TEMPNOTSET,
	TempF:            TEMPNOTSET,
	ExternalTempF:    TEMPNOTSET,
	Humidity:         HUMIDITYNOTSET,
	LightInternal:    LIGHTNOTSET,
	PressureInternal: PRESSURENOTSET,
}

type GenericSensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}
