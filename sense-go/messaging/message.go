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
package messaging

// copyright and license inspection - no issues 4/13/22

import (
	"bubblesnet/edge-device/sense-go/globals"
	"fmt"
	"time"
)

func getNowMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

type ADCSensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	ChannelNumber     int     `json:"channel_number"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
	Gain              int     `json:"gain"`
	Rate              int     `json:"rate"`
}

type GenericSensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}

type VOCSensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}

type CCS811CurrentMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}

type CCS811VoltageMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}

type CO2SensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
}

type DistanceSensorMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
	ContainerName     string  `json:"container_name"`
	ExecutableVersion string  `json:"executable_version"`
	SampleTimestamp   int64   `json:"sample_timestamp,omitempty"`
	MessageType       string  `json:"message_type"`
	SensorName        string  `json:"sensor_name"`
	MeasurementName   string  `json:"measurement_name"`
	Value             float64 `json:"value,omitempty"`
	FloatValue        float64 `json:"floatvalue,omitempty"`
	Units             string  `json:"units"`
	Direction         string  `json:"direction"`
	DistanceCm        float64 `json:"distance_cm,omitempty"`
	DistanceIn        float64 `json:"distance_in,omitempty"`
}

type PhMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
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

type TamperEventMessage struct {
	DeviceId          int64   `json:"deviceid"`
	StationId         int64   `json:"stationid"`
	SiteId            int64   `json:"siteid"`
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
	XMove             float64 `json:"xmove,omitempty"`
	YMove             float64 `json:"ymove,omitempty"`
	ZMove             float64 `json:"zmove,omitempty"`
}

type SwitchStatusChangeMessage struct {
	DeviceId          int64  `json:"deviceid"`
	StationId         int64  `json:"stationid"`
	SiteId            int64  `json:"siteid"`
	ContainerName     string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	EventTimestamp    int64  `json:"event_timestamp"`
	MessageType       string `json:"message_type"`
	SwitchName        string `json:"switch_name"`
	On                bool   `json:"on"`
}

type DispenserStatusChangeMessage struct {
	DeviceId          int64  `json:"deviceid"`
	StationId         int64  `json:"stationid"`
	SiteId            int64  `json:"siteid"`
	ContainerName     string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	EventTimestamp    int64  `json:"event_timestamp"`
	MessageType       string `json:"message_type"`
	DispenserName     string `json:"dispenser_name"`
	On                bool   `json:"on"`
}

type PictureTakenMessage struct {
	DeviceId              int64  `json:"deviceid"`
	StationId             int64  `json:"stationid"`
	SiteId                int64  `json:"siteid"`
	ContainerName         string `json:"container_name"`
	ExecutableVersion     string `json:"executable_version"`
	EventTimestamp        int64  `json:"event_timestamp,omitempty"`
	PictureFilename       string `json:"picture_filename"`
	PictureDateTimeMillis int64  `json:"picture_datetime_millis"`
	MessageType           string `json:"message_type"`
}

func NewSwitchStatusChangeMessage(switch_name string, on bool) (pmsg *SwitchStatusChangeMessage) {
	msg := SwitchStatusChangeMessage{
		DeviceId:      globals.MyDevice.DeviceID,
		StationId:     globals.MyStation.StationID,
		SiteId:        globals.MySite.SiteID,
		ContainerName: globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		EventTimestamp: getNowMillis(),
		MessageType:    "switch_event",
		SwitchName:     switch_name,
		On:             on}
	return &msg

}

func NewDispenserStatusChangeMessage(dispenser_name string, on bool) (pmsg *DispenserStatusChangeMessage) {
	msg := DispenserStatusChangeMessage{
		DeviceId:      globals.MyDevice.DeviceID,
		StationId:     globals.MyStation.StationID,
		SiteId:        globals.MySite.SiteID,
		ContainerName: globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		EventTimestamp: getNowMillis(),
		MessageType:    "dispenser_event",
		DispenserName:  dispenser_name,
		On:             on}
	return &msg

}

func NewPictureTakenMessage(PictureFilename string, PictureDatetimeMillis int64) (pmsg *PictureTakenMessage) {
	msg := PictureTakenMessage{
		DeviceId:      globals.MyDevice.DeviceID,
		StationId:     globals.MyStation.StationID,
		SiteId:        globals.MySite.SiteID,
		ContainerName: globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		EventTimestamp:        getNowMillis(),
		MessageType:           "picture_event",
		PictureFilename:       PictureFilename,
		PictureDateTimeMillis: PictureDatetimeMillis,
	}
	return &msg

}

func NewVOCSensorMessage(sensor_name string, measurement_name string, value float64, units string, direction string) (pmsg *VOCSensorMessage) {
	msg := VOCSensorMessage{
		DeviceId:        globals.MyDevice.DeviceID,
		StationId:       globals.MyStation.StationID,
		SiteId:          globals.MySite.SiteID,
		ContainerName:   globals.ContainerName,
		MeasurementName: measurement_name,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SensorName:      sensor_name,
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
	}

	return &msg
}

func NewCO2SensorMessage(sensor_name string, measurement_name string, value float64, units string, direction string) (pmsg *CO2SensorMessage) {
	msg := CO2SensorMessage{
		DeviceId:        globals.MyDevice.DeviceID,
		StationId:       globals.MyStation.StationID,
		SiteId:          globals.MySite.SiteID,
		ContainerName:   globals.ContainerName,
		MeasurementName: measurement_name,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SensorName:      sensor_name,
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
	}

	return &msg
}

func NewCCS811CurrentMessage(sensor_name string, measurement_name string, value float64, units string, direction string) (pmsg *CCS811CurrentMessage) {
	msg := CCS811CurrentMessage{
		DeviceId:        globals.MyDevice.DeviceID,
		StationId:       globals.MyStation.StationID,
		SiteId:          globals.MySite.SiteID,
		ContainerName:   globals.ContainerName,
		MeasurementName: measurement_name,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SensorName:      sensor_name,
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
	}
	return &msg
}

func NewCCS811VoltageMessage(sensor_name string, measurement_name string, value float64, units string, direction string) (pmsg *CCS811VoltageMessage) {
	msg := CCS811VoltageMessage{
		DeviceId:        globals.MyDevice.DeviceID,
		StationId:       globals.MyStation.StationID,
		SiteId:          globals.MySite.SiteID,
		ContainerName:   globals.ContainerName,
		MeasurementName: measurement_name,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SensorName:      sensor_name,
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
	}
	return &msg
}

func NewGenericSensorMessage(sensor_name string, measurement_name string, value float64, units string, direction string) (pmsg *GenericSensorMessage) {
	msg := GenericSensorMessage{
		DeviceId:        globals.MyDevice.DeviceID,
		StationId:       globals.MyStation.StationID,
		SiteId:          globals.MySite.SiteID,
		ContainerName:   globals.ContainerName,
		MeasurementName: measurement_name,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SensorName:      sensor_name,
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
	}

	return &msg
}

func NewADCSensorMessage(sensor_name string, measurement_name string, value float64, units string, direction string, channel int, gain int, rate int) (pmsg *ADCSensorMessage) {
	msg := ADCSensorMessage{
		DeviceId:      globals.MyDevice.DeviceID,
		StationId:     globals.MyStation.StationID,
		SiteId:        globals.MySite.SiteID,
		ContainerName: globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		SensorName:      sensor_name,
		MeasurementName: measurement_name,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
		ChannelNumber:   channel,
		Gain:            gain,
		Rate:            rate,
	}

	return &msg
}

func NewDistanceSensorMessage(sensor_name string, measurement_name string, value float64, units string, direction string, distanceCm float64, distanceIn float64) (pmsg *DistanceSensorMessage) {
	msg := DistanceSensorMessage{
		DeviceId:      globals.MyDevice.DeviceID,
		StationId:     globals.MyStation.StationID,
		SiteId:        globals.MySite.SiteID,
		ContainerName: globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:     globals.Message_type_measurement,
		SensorName:      sensor_name,
		MeasurementName: measurement_name,
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
		DistanceCm:      distanceCm,
		DistanceIn:      distanceIn,
	}

	return &msg
}

func NewTamperSensorMessage(sensor_name string, value float64, units string, direction string, moveX float64, moveY float64, moveZ float64) (pmsg *TamperEventMessage) {
	msg := TamperEventMessage{
		DeviceId:      globals.MyDevice.DeviceID,
		StationId:     globals.MyStation.StationID,
		SiteId:        globals.MySite.SiteID,
		ContainerName: globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:     "event",
		SensorName:      sensor_name,
		MeasurementName: "tamper",
		Value:           value,
		FloatValue:      value,
		Units:           units,
		Direction:       direction,
		XMove:           moveX,
		YMove:           moveY,
		ZMove:           moveZ,
	}

	return &msg
}
