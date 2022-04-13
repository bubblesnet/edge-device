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
