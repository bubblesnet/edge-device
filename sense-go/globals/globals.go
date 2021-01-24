package globals

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"github.com/go-playground/log"
	"runtime"
)

// These are shadows of vars in main
var BubblesnetVersionMajorString string
var BubblesnetVersionMinorString=""
var BubblesnetVersionPatchString=""
var BubblesnetBuildNumberString=""
var BubblesnetBuildTimestamp=""
var BubblesnetGitHash=""

var ContainerName = "sense-go"

var DevicesFailed []string

var Config = Configuration{}

var DeviceId = int64(70000007)
var UserId = int64(90000009)


const (
	ForwrdingAddress     = "store-and-forward:50051"
)

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

var Sequence int32
var Client pb.SensorStoreAndForwardClient

func RunningOnUnsupportedHardware() (notSupported bool) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || (runtime.GOARCH != "arm"  && runtime.GOARCH != "arm64") {
		return true
	}
	return false
}

func ReportDeviceFailed(devicename string) {
	for i := 0; i < len(DevicesFailed); i++ {
		if DevicesFailed[i] == devicename {
			return
		}
	}
	log.Errorf("Adding device %s to failed list", devicename)
	DevicesFailed = append(DevicesFailed,devicename)
}

func GetSequence() (int32){
	if Sequence > 200000 {
		Sequence = 100001
	} else {
		Sequence = Sequence + 1
	}
	return Sequence
}