package globals

var Config = Configuration{}

type LocalState struct {
	EnvironmentalControl string  `json:"environmental_control,omitempty"`
	Humidifier bool `json:"humidifier"`
	Heater bool `json:"heater"`
	HeaterPad bool `json:"heater_pad"`
	GrowLightVeg bool `json:"grow_light_veg"`

}

var LocalCurrentState = LocalState {
	EnvironmentalControl: "",
	GrowLightVeg:             false,
	HeaterPad:              false,
	Heater:              false,
	Humidifier: false,
}

const INLETFAN string = "inlet_fan"
const WATERPUMP string = "water_pump"
const GROWLIGHTVEG string = "light_vegetative"
const HEATPAD string = "light_bloom"
const HEATLAMP string = "heat_lamp"
const AIRPUMP string = "air_pump"
const OUTLETFAN string = "exhaust_fan"
const HUMIDIFIER string = "humidifier"

var CurrentStageSchedule StageSchedule

var Lasttemp float32
var Lasthumidity float32



