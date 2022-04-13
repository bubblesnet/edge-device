package globals

// copyright and license inspection - no issues 4/13/22

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"github.com/go-playground/log"
	"runtime"
)

// These are shadows of vars in main
var BubblesnetVersionMajorString string
var BubblesnetVersionMinorString = ""
var BubblesnetVersionPatchString = ""
var BubblesnetBuildNumberString = ""
var BubblesnetBuildTimestamp = ""
var BubblesnetGitHash = ""

var ContainerName = "sense-go"

var PersistentStoreMountPoint = "/config" // Can be changed for units

var DevicesFailed []string

var MySite = Site{}
var MyDevice *EdgeDevice
var MyStation *Station

var MyDeviceID = int64(0)

var PollingWaitInSeconds = 10

// var DeviceId = int64(0)
// var UserId = int64(0)

const (
	ForwardingAddress = "store-and-forward:50051"
)

type LocalState struct {
	EnvironmentalControl string `json:"environmental_control,omitempty"`
	Humidifier           bool   `json:"humidifier"`
	Heater               bool   `json:"heater"`
	WaterHeater          bool   `json:"water_heater"`
	HeaterPad            bool   `json:"heater_pad"`
	GrowLightVeg         bool   `json:"grow_light_veg"`
	GrowLightBloom       bool   `json:"grow_light_bloom"`
}

var LocalCurrentState = LocalState{
	EnvironmentalControl: "",
	GrowLightBloom:       false,
	GrowLightVeg:         false,
	HeaterPad:            false,
	Heater:               false,
	WaterHeater:          false,
	Humidifier:           false,
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
const WATERHEATER string = "waterHeater"

const GERMINATION string = "germinate"
const SEEDLING string = "seedling"
const VEGETATIVE string = "vegetative"
const BLOOMING string = "bloom"
const HARVEST string = "harvest"
const CURING string = "cure"
const DRYING string = "dry"
const IDLE string = "idle"

var CurrentStageSchedule StageSchedule

var LastTemp float32
var LastWaterTemp float32
var LastHumidity float32
var LastWaterLevel float32

var Sequence int32
var Client pb.SensorStoreAndForwardClient

func RunningOnUnsupportedHardware() (notSupported bool) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || (runtime.GOARCH != "arm" && runtime.GOARCH != "arm64") {
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
	DevicesFailed = append(DevicesFailed, devicename)
}

func GetSequence() int32 {
	if Sequence > 200000 {
		Sequence = 100001
	} else {
		Sequence = Sequence + 1
	}
	return Sequence
}
