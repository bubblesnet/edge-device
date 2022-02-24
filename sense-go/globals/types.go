package globals

const TEMPNOTSET float32 = -100.5
const HUMIDITYNOTSET float32 = -100.5

type externalState struct {
	TempF    float32 `json:"tempF"`
	Humidity float32 `json:"humidity"`
}

var ExternalCurrentState = externalState{
	TempF:    TEMPNOTSET,
	Humidity: HUMIDITYNOTSET,
}
