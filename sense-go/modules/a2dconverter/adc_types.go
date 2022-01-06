package a2dconverter

type ChannelConfig struct {
	gain int
	rate int
}

type AdapterConfig struct {
	bus_id        int
	address       int
	channelConfig [4]ChannelConfig
}

type ChannelValue struct {
	ChannelNumber int     `json:"channel_number"`
	Voltage       float64 `json:"voltage,omitempty"`
	Gain          int     `json:"gain"`
	Rate          int     `json:"rate"`
}

type Channels [4]ChannelValue

var a0 = AdapterConfig{
	bus_id:  1,
	address: 0x48,
	channelConfig: [4]ChannelConfig{
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
	},
}

var a1 = AdapterConfig{
	bus_id:  1,
	address: 0x49,
	channelConfig: [4]ChannelConfig{
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
	},
}
var daps = []AdapterConfig{a0, a1}

type ADCMessage struct {
	BusId         int      `json:"bus_id"`
	Address       int      `json:"address"`
	ChannelValues Channels `json:"channel_values"`
}
