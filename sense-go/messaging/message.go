package messaging

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
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	ChannelNumber int `json:"channel_number,omitempty"`
	Value float64 `json:"value,omitempty"`
	Units string	`json:"units"`
	Gain    int	`json:"gain,omitempty"`
	Rate    int	`json:"rate,omitempty"`
}

type GenericSensorMessage struct {
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	Value float64 `json:"value,omitempty"`
	Units string	`json:"units"`
}

type DistanceSensorMessage struct {
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	Value float64 `json:"value"`
	Units string `json:"units"`
	DistanceCm float64 `json:"distance_cm,omitempty"`
	DistanceIn float64 `json:"distance_in,omitempty"`
}

type PhMessage struct {
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	Value float64 `json:"value"`
	Units string `json:"units"`
}

type TamperSensorMessage struct {
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	Value float64 `json:"value"`
	Units string `json:"units"`
	XMove float64	`json:"xmove,omitempty"`
	YMove float64	`json:"ymove,omitempty"`
	ZMove float64	`json:"zmove,omitempty"`
}

func NewGenericSensorMessage( sensor_name string, value float64, units string ) (pmsg *GenericSensorMessage) {
	msg := GenericSensorMessage{
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		Value:             value,
		Units: units,
	}

	return &msg
}

func NewADCSensorMessage( sensor_name string, value float64, units string, channel int, gain int, rate int ) (pmsg *ADCSensorMessage) {
	msg := ADCSensorMessage{
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		SensorName: sensor_name,
		Value:             value,
		Units: units,
		ChannelNumber: channel,
		Gain: gain,
		Rate: rate,
	}

	return &msg
}

func NewDistanceSensorMessage( sensor_name string, value float64, units string, distanceCm float64, distanceIn float64 ) (pmsg *DistanceSensorMessage) {
	msg := DistanceSensorMessage{
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		SensorName: sensor_name,
		Value:             value,
		Units: units,
		DistanceCm: distanceCm,
		DistanceIn: distanceIn,
	}

	return &msg
}

func NewTamperSensorMessage( sensor_name string, value float64, units string, moveX float64, moveY float64, moveZ float64 ) (pmsg *TamperSensorMessage) {
	msg := TamperSensorMessage{
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		SensorName: sensor_name,
		Value:             value,
		Units: units,
		XMove: moveX,
		YMove: moveY,
		ZMove: moveZ,
	}

	return &msg
}
