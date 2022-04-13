package globals

// copyright and license inspection - no issues 4/13/22

const TEMPNOTSET float32 = -100.5
const HUMIDITYNOTSET float32 = -100.5

type ExternalState struct {
	WaterTempF    float32 `json:"waterTempF"`
	TempF         float32 `json:"tempF"`
	Humidity      float32 `json:"humidity"`
	ExternalTempF float32 `json:"externalTempF"`
}

var ExternalCurrentState = ExternalState{
	WaterTempF:    TEMPNOTSET,
	TempF:         TEMPNOTSET,
	Humidity:      HUMIDITYNOTSET,
	ExternalTempF: TEMPNOTSET,
}
