package globals

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/log"
	"io/ioutil"
	"net/http"
	"strings"
)


/**
Top-level object in the data hierarchy.  A farm is identified by the user/owner
and contains multiple cabinets.
 */
type Farm struct {
	FarmID int64 `json:"farmid"`
	UserID int64 `json:"userid"`
	ControllerHostName	string		`json:"controller_hostname"`
	ControllerAPIPort	int			`json:"controller_api_port"`
	LogLevel         string           `json:"log_level,omitempty"`
	AutomaticControl bool 			`json:"automatic_control"`
	Cabinets []Cabinet `json:"cabinets"`
}

/**
A cabinet is a grow-unit, typically either a cabinet or a tent.  A cabinet
contains multiple edge devices, typically Raspberry Pi.
 */
type Cabinet struct {
	CabinetID	int64`json:"cabinetid"`
	HeightSensor bool `json:"height_sensor"`
	Humidifier bool `json:"humidifier"`
	HumiditySensor bool `json:"humidity_sensor_internal"`
	ExternalHumiditySensor bool `json:"humidity_sensor_external"`
	Heater bool `json:"heater"`
	ThermometerTop bool `json:"thermometer_top"`
	ThermometerMiddle bool `json:"thermometer_middle"`
	ThermometerBottom bool `json:"thermometer_bottom"`
	ThermometerExternal bool `json:"thermometer_external"`
	ThermometerWater bool `json:"thermometer_water"`
	WaterPump bool `json:"waterPump"`
	AirPump bool `json:"airPump"`
	LightSensor bool `json:"light_sensor_internal"`
	CabinetDoorSensor bool `json:"cabinet_door_sensor"`
	OuterDoorSensor bool `json:"outer_door_sensor"`
	MovementSensor bool `json:"movement_sensor"`
	PressureSensor bool `json:"pressure_sensors"`
	RootPhSensor bool `json:"root_ph_sensor"`
	EnclosureType string `json:"enclosure_type"`
	WaterLevelSensor bool `json:"water_level_sensor"`
	IntakeFan bool `json:"intakeFan"`
	ExhaustFan bool `json:"exhaustFan"`
	HeatLamp bool `json:"heatLamp"`
	HeatingPad bool `json:"heatingPad"`
	LightBloom bool `json:"lightBloom"`
	LightVegetative bool `json:"lightVegetative"`
	LightGerminate bool `json:"lightGerminate"`
	Relay          bool `json:"relay,omitempty"`
	AttachedDevices	[]AttachedDevice `json:"attached_devices"`
	StageSchedules  []StageSchedule  `json:"stage_schedules,omitempty"`
	CurrentStage	string `json:current_stage`
	LightOnHour     int  `json:"light_on_hour,omitempty"`
	TamperSpec                        Tamper           `json:"tamper"`
}

/**
A DeviceModule is typically an add-on board attached to the edge device that
generates one or more types of measurements.  An AttachedDevice can have multiple
DeviceModules.
*/
type DeviceModule struct {
	ModuleID	int64	`json:"moduleid"`
	ModuleName      string `json:"module_name"`
	InternalAddress string `json:"internal_address"`
	IncludedSensors []Sensor `json:"included_sensors"`
}

type Sensor struct {
	SensorID int64 `json:"sensorid"`
	SensorName	string	`json:"sensor_name"`
	MeasurementName	string `json:"measurement_name"`
}
/**
An AttachedDevice is a single-board-computer that, with the other
AttachedDevices in the Cabinet, implements the intelligence of the Cabinet
such that the ideal grow-conditions for the plants inside the Cabinet are
always maintained, and a stream of event and environmental sensor messages are sent to
the time-series database.
 */
type AttachedDevice struct {
	ContainerName	string        `json:"container_name"`
	DeviceID		int64          `json:"deviceid"`
	DeviceType	string           `json:"device_type"`
	Protocol	string             `json:"protocol"`
	Address	string              `json:"address"`
	DeviceModules []DeviceModule `json:"included_modules"`
	ACOutlets  [8]ACOutlet      `json:"ac_outlets,omitempty"`
	Camera PiCam	`json:"camera,omitempty"`
	TimeBetweenSensorPollingInSeconds int64            `json:"time_between_sensor_polling_in_seconds"`
}

type EnvironmentalTarget struct {
	Temperature float32 `json:"temperature,omitempty"`
	Humidity float32 `json:"humidity,omitempty"`
}

type StageSchedule struct {
	Name                 string              `json:"name,omitempty"`
	HoursOfLight         int                 `json:"hours_of_light,omitempty"`
	EnvironmentalTargets EnvironmentalTarget `json:"environmental_targets,omitempty"`
}

type ControlState struct {
}

type PiCam struct {
	PiCamera	bool		`json:"picamera"`
	ResolutionX	int		`json:"resolutionX"`
	ResolutionY	int		`json:"resolutionY"`
}

type Tamper struct {
	Xmove float64		`json:"xmove"`
	Ymove float64			`json:"ymove"`
	Zmove float64			`json:"zmove"`
}

type Configuration1 struct {
	ControllerHostName                string           `json:"controller_hostname"`
	ControllerAPIPort                 int              `json:"controller_api_port"`
	UserID                            int64            `json:"userid"`
	DeviceID                          int64            `json:"deviceid"`
	Stage                             string           `json:"stage,omitempty"`
	LightOnHour                       int              `json:"light_on_hour,omitempty"`
	StageSchedules                    []StageSchedule  `json:"stage_schedules,omitempty"`
	Camera                            PiCam            `json:"camera,omitempty"`
	DeviceSettings                    Cabinet          `json:"device_settings"`
	LogLevel                          string           `json:"log_level,omitempty"`
	AttachedDevices                   []AttachedDevice `json:"attached_devices"`
	AutomaticControl                  bool             `json:"automatic_control"`
	TimeBetweenSensorPollingInSeconds int64            `json:"time_between_sensor_polling_in_seconds"`
	TamperSpec                        Tamper           `json:"tamper"`
}

type ACOutlet struct {
	DeviceID int64 `json:"deviceid"`
	Name string `json:"name,omitempty"`
	Index int `json:"index,omitempty"`
	PowerOn bool `json:"on,omitempty"`
	BCMPinNumber int `json:"bcm_pin_number,omitempty"`
}

func ReadFromPersistentStore(storeMountPoint string, relativePath string, fileName string, farm *Farm, currentStageSchedule *StageSchedule) error {
	log.Debug("readConfig")
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	if relativePath == "" {
		fullpath = storeMountPoint + "/" + fileName
	}
	fmt.Printf("readConfig from %s", fullpath)
	file, _ := ioutil.ReadFile(fullpath)

	_ = json.Unmarshal([]byte(file), farm)

	log.Debugf("data = %v", *farm)

	for i := 0; i < len(MyCabinet.StageSchedules); i++ {
		if MyCabinet.StageSchedules[i].Name == MyCabinet.CurrentStage {
			*currentStageSchedule = MyCabinet.StageSchedules[i]
			log.Infof("Current stage is %s - schedule is %v", MyCabinet.CurrentStage, currentStageSchedule)
			return nil
		}
	}
	errstr := fmt.Sprintf("ERROR: No schedule for stage (%s)", MyCabinet.CurrentStage)
	log.Error(errstr)
	return errors.New(errstr)
}

// CustomHandler is your custom handler
type CustomHandler struct {
	// whatever properties you need
}

// Log accepts log entries to be processed
func (c *CustomHandler) Log(e log.Entry) {

	// below prints to os.Stderr but could marshal to JSON
	// and send to central logging server
	//																						       ---------
	// 				                                                                 |----------> | console |
	//                                                                               |             ---------
	// i.e. -----------------               -----------------     Unmarshal    -------------       --------
	//     | app log handler | -- json --> | central log app | --    to    -> | log handler | --> | syslog |
	//      -----------------               -----------------       Entry      -------------       --------
	//      																         |             ---------
	//                                  									         |----------> | DataDog |
	//
	//         																	        	   ---------
	b := new(bytes.Buffer)
	b.Reset()
	b.WriteString(e.Message)

	for _, f := range e.Fields {
		fmt.Fprintf(b, " %s=%v", f.Key, f.Value)
	}
	fmt.Println(b.String())
}

func ConfigureLogging( farm Farm, containerName string) {
	cLog := new(CustomHandler)

	if strings.Contains(farm.LogLevel,"error") {
		log.AddHandler(cLog, log.ErrorLevel)
	}
	if strings.Contains(farm.LogLevel,"warn") {
		log.AddHandler(cLog, log.WarnLevel)
	}
	if strings.Contains(farm.LogLevel,"debug") {
		log.AddHandler(cLog, log.DebugLevel)
	}
	if strings.Contains(farm.LogLevel,"info") {
		log.AddHandler(cLog, log.InfoLevel)
	}
	if strings.Contains(farm.LogLevel,"notice") {
		log.AddHandler(cLog, log.NoticeLevel)
	}
	if strings.Contains(farm.LogLevel,"panic") {
		log.AddHandler(cLog, log.PanicLevel)
	}

}

func GetConfigFromServer(storeMountPoint string, relativePath string, fileName string) (err error) {
	url := fmt.Sprintf("http://%s:%d/api/farm/%8.8d/%8.8d", MyFarm.ControllerHostName, MyFarm.ControllerAPIPort, MyFarm.UserID, MyDevice.DeviceID)
	fmt.Printf("Sending to %s", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("post error %v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("readall error %v", err)
		return err
	}
	fmt.Printf("response %s", string(body))
	newconfig := Farm{}
	if err = json.Unmarshal(body, &newconfig); err != nil {
		fmt.Printf("err on farm %v\n", err)
		return errors.New("err on farm")
	}
	MyFarm = newconfig

	fmt.Printf("set farm to newconfig %v\n", MyFarm)

	bytes, err := json.MarshalIndent(MyFarm, "", "  ")
	filepath := fmt.Sprintf("%s/%s/%s", storeMountPoint,relativePath,fileName)
	if len(relativePath) == 0 {
		filepath = fmt.Sprintf("%s/%s", storeMountPoint,fileName)
	}
	fmt.Printf("writing farm to file %s\n\n",filepath)
	err = ioutil.WriteFile(filepath, bytes, 0777)
	if err != nil {
		log.Errorf("error save farm file %v", err)
		return(err)
	}

	fmt.Printf("received farm\n\n")
	return nil
}
