package globals

const TEMPNOTSET float32 = -100.5
const HUMIDITYNOTSET float32 = -100.5

type DistanceMessage struct {
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	DistanceCm float64 `json:"distance_cm,omitempty"`
	DistanceIn float64 `json:"distance_in,omitempty"`
}

type PhMessage struct {
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	Ph float64 `json:"ph,omitempty"`
}

type TamperMessage struct {
	SampleTimestamp int64 `json:"sample_timestamp,omitempty"`
	XMove float64	`json:"xmove,omitempty"`
	YMove float64	`json:"ymove,omitempty"`
	ZMove float64	`json:"zmove,omitempty"`
}

type externalState struct {
	TempF                float32 `json:"tempF,omitempty"`
	Humidity                float32 `json:"humidity,omitempty"`
}

var ExternalCurrentState = externalState{
	TempF: TEMPNOTSET,
	Humidity: HUMIDITYNOTSET,
}

