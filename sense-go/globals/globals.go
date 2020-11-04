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

const INLETFAN string = "Inlet fan"
const WATERPUMP string = "Water pump"
const GROWLIGHTVEG string = "Grow light - veg"
const HEATPAD string = "Grow light - bloom"
const HEATLAMP string = "Heat lamp"
const AIRPUMP string = "Air pump"
const OUTLETFAN string = "Outlet fan"
const HUMIDIFIER string = "Humidifier"

var CurrentStageSchedule StageSchedule

var Lasttemp float32
var Lasthumidity float32



