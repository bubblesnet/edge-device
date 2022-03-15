package globals

const TEMPNOTSET float32 = -100.5
const HUMIDITYNOTSET float32 = -100.5

type externalState struct {
	WaterTempF float32 `json:"watertempF"`
	TempF      float32 `json:"tempF"`
	Humidity   float32 `json:"humidity"`
}

var ExternalCurrentState = externalState{
	WaterTempF: TEMPNOTSET,
	TempF:      TEMPNOTSET,
	Humidity:   HUMIDITYNOTSET,
}
