package main

const TEMPNOTSET float64 = -100.5
const HUMIDITYNOTSET float64 = -100.5
const PRESSURENOTSET float64 = -100.5
const LIGHTNOTSET float64 = -100.5

type externalState struct {
	WaterTempF       float64 `json:"tempF,omitemp"`
	TempF            float64 `json:"tempF,omitempty"`
	Humidity         float64 `json:"humidity,omitempty"`
	PressureInternal float64 `json:"pressure_internal,omitempty"`
	LightInternal    float64 `json:"pressure_internal,omitempty"`
}

var ExternalCurrentState = externalState{
	WaterTempF:       TEMPNOTSET,
	TempF:            TEMPNOTSET,
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
