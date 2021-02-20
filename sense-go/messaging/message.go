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
	DeviceId int64`json:"deviceid"`
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	MeasurementName string `json:"measurement_name"`
	ChannelNumber int `json:"channel_number,omitempty"`
	Value float64 `json:"value,omitempty"`
	Units string	`json:"units"`
	Direction string `json:"direction"`
	Gain    int	`json:"gain,omitempty"`
	Rate    int	`json:"rate,omitempty"`
}

type GenericSensorMessage struct {
	DeviceId int64`json:"deviceid"`
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	MeasurementName string `json:"measurement_name"`
	Value float64 `json:"value"`
	Units string	`json:"units"`
	Direction string `json:"direction"`
}

type DistanceSensorMessage struct {
	DeviceId int64`json:"deviceid"`
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	MeasurementName string `json:"measurement_name"`
	Value float64 `json:"value"`
	Units string `json:"units"`
	Direction string `json:"direction"`
	DistanceCm float64 `json:"distance_cm,omitempty"`
	DistanceIn float64 `json:"distance_in,omitempty"`
}

type PhMessage struct {
	DeviceId int64`json:"deviceid"`
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	MeasurementName string `json:"measurement_name"`
	Value float64 `json:"value"`
	Units string `json:"units"`
	Direction string `json:"direction"`
}


type TamperSensorMessage struct {
	DeviceId int64`json:"deviceid"`
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	MeasurementName string `json:"measurement_name"`
	Value float64 `json:"value"`
	Units string `json:"units"`
	Direction string `json:"direction"`
	XMove float64	`json:"xmove,omitempty"`
	YMove float64	`json:"ymove,omitempty"`
	ZMove float64	`json:"zmove,omitempty"`
}

type SwitchStatusChangeMessage struct {
	DeviceId int64`json:"deviceid"`
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	EventTimestamp int64 `json:"event_timestamp,omitempty"`
	MessageType string `json:"message_type"`
	SwitchName string `json:"switch_name"`
	On bool `json:"on"`
	}

func NewSwitchStatusChangeMessage( switch_name string, on bool ) (pmsg *SwitchStatusChangeMessage){
	msg := SwitchStatusChangeMessage{
		DeviceId: globals.Config.DeviceID,
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		EventTimestamp: getNowMillis(),
		MessageType: "switch_event",
		SwitchName: switch_name,
		On: on}
	return &msg

}

func NewGenericSensorMessage( sensor_name string, measurement_name string, value float64, units string, direction string ) (pmsg *GenericSensorMessage) {
	msg := GenericSensorMessage{
		DeviceId: globals.Config.DeviceID,
		ContainerName:     globals.ContainerName,
		MeasurementName: measurement_name,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
			SensorName: sensor_name,
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		Value:             value,
		Units: units,
		Direction: direction,
	}

	return &msg
}

func NewADCSensorMessage( sensor_name string, measurement_name string, value float64, units string, direction string, channel int, gain int, rate int ) (pmsg *ADCSensorMessage) {
	msg := ADCSensorMessage{
		DeviceId: globals.Config.DeviceID,
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		SensorName: sensor_name,
		MeasurementName: measurement_name,
		Value:             value,
		Units: units,
		Direction: direction,
		ChannelNumber: channel,
		Gain: gain,
		Rate: rate,
	}

	return &msg
}

func NewDistanceSensorMessage( sensor_name string, measurement_name string, value float64, units string, direction string, distanceCm float64, distanceIn float64 ) (pmsg *DistanceSensorMessage) {
	msg := DistanceSensorMessage{
		DeviceId: globals.Config.DeviceID,
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		SensorName: sensor_name,
		MeasurementName: measurement_name,
		Value:             value,
		Units: units,
		Direction: direction,
		DistanceCm: distanceCm,
		DistanceIn: distanceIn,
	}

	return &msg
}

func NewTamperSensorMessage( sensor_name string, value float64, units string, direction string, moveX float64, moveY float64, moveZ float64 ) (pmsg *TamperSensorMessage) {
	msg := TamperSensorMessage{
		DeviceId: globals.Config.DeviceID,
		ContainerName:     globals.ContainerName,
		ExecutableVersion: fmt.Sprintf("%s.%s.%s %s %s",
			globals.BubblesnetVersionMajorString, globals.BubblesnetVersionMinorString,
			globals.BubblesnetVersionPatchString, globals.BubblesnetBuildTimestamp, globals.BubblesnetGitHash),
		SampleTimestamp: getNowMillis(),
		MessageType:       "measurement",
		SensorName: sensor_name,
		MeasurementName: "movement",
		Value:             value,
		Units: units,
		Direction: direction,
		XMove: moveX,
		YMove: moveY,
		ZMove: moveZ,
	}

	return &msg
}
