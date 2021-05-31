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

var PersistentStoreMountPoint = "/config"	// Can be changed for units

var DevicesFailed []string

var MySite = Site{}
var MyDevice *EdgeDevice
var MyStation *Station

var MyDeviceID = int64(0)

// var DeviceId = int64(0)
// var UserId = int64(0)

const (
	ForwardingAddress = "store-and-forward:50051"
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

const INLETFAN string = "inletFan"
const WATERPUMP string = "waterPump"
const GROWLIGHTVEG string = "lightVegetative"
const GROWLIGHTBLOOM string = "lightBloom"
const HEATPAD string = "lightBloom"
const HEATLAMP string = "heatLamp"
const AIRPUMP string = "airPump"
const OUTLETFAN string = "exhaustFan"
const HUMIDIFIER string = "humidifier"
const HEATER string = "heater"

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