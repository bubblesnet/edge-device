package globals

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)


/**
Top-level object in the data hierarchy.  A site is identified by the user/owner
and contains multiple stations.
 */
type Site struct {
	SiteID             int64     `json:"siteid"`
	UserID             int64     `json:"userid"`
	ControllerHostName string    `json:"controller_hostname"`
	ControllerAPIPort  int       `json:"controller_api_port"`
	LogLevel           string    `json:"log_level,omitempty"`
	AutomaticControl   bool      `json:"automatic_control"`
	Stations           []Station `json:"stations,omitempty"`
}

/**
A station is a grow-unit, typically either a cabinet or a tent.  A station
contains multiple edge devices, typically Raspberry Pi.
 */
type Station struct {
	StationID              int64 `json:"stationid"`
	HeightSensor           bool  `json:"height_sensor,omitempty"`
	Humidifier             bool  `json:"humidifier,omitempty"`
	HumiditySensor         bool  `json:"humidity_sensor_internal,omitempty"`
	ExternalHumiditySensor bool   `json:"humidity_sensor_external,omitempty"`
	Heater                 bool   `json:"heater,omitempty"`
	ThermometerTop         bool   `json:"thermometer_top,omitempty"`
	ThermometerMiddle      bool   `json:"thermometer_middle,omitempty"`
	ThermometerBottom      bool   `json:"thermometer_bottom,omitempty"`
	ThermometerExternal    bool   `json:"thermometer_external,omitempty"`
	ThermometerWater       bool   `json:"thermometer_water,omitempty"`
	WaterPump              bool   `json:"waterPump,omitempty"`
	AirPump                bool   `json:"airPump,omitempty"`
	LightSensor            bool   `json:"light_sensor_internal,omitempty"`
	StationDoorSensor      bool   `json:"station_door_sensor,omitempty"`
	OuterDoorSensor        bool   `json:"outer_door_sensor,omitempty"`
	MovementSensor         bool   `json:"movement_sensor,omitempty"`
	PressureSensor         bool   `json:"pressure_sensors,omitempty"`
	RootPhSensor           bool   `json:"root_ph_sensor,omitempty"`
	EnclosureType          string `json:"enclosure_type,omitempty"`
	WaterLevelSensor       bool   `json:"water_level_sensor,omitempty"`
	IntakeFan              bool   `json:"intakeFan,omitempty"`
	ExhaustFan             bool   `json:"exhaustFan,omitempty"`
	HeatLamp               bool   `json:"heatLamp,omitempty"`
	HeatingPad             bool   `json:"heatingPad,omitempty"`
	LightBloom 		bool `json:"lightBloom,omitempty"`
	LightVegetative bool `json:"lightVegetative,omitempty"`
	LightGerminate 	bool `json:"lightGerminate,omitempty"`
	Relay          	bool `json:"relay,omitempty,omitempty"`
	EdgeDevices		[]EdgeDevice `json:"edge_devices,omitempty"`
	StageSchedules  []StageSchedule  `json:"stage_schedules,omitempty"`
	CurrentStage	string `json:"current_stage,omitempty"`
	LightOnHour     int  `json:"light_on_hour,omitempty"`
	TamperSpec		Tamper           `json:"tamper,omitempty"`
}

/**
An EdgeDevice is a single-board-computer that, with the other
AttachedDevices in the Station, implements the intelligence of the Station
such that the ideal grow-conditions for the plants inside the Station are
always maintained, and a stream of event and environmental sensor messages are sent to
the time-series database.
*/
type EdgeDevice struct {
	DeviceID	int64  `json:"deviceid"`
	DeviceType	string	`json:"devicetypename,omitempty"`
	ExternalID 	string `json:"externalid,omitempty"`
	IPAddress 	string `json:"ipaddress,omitempty"`
	MacAddress 	string `json:"macaddress,omitempty"`
	DeviceModules []DeviceModule `json:"modules,omitempty"`
	Camera 		PiCam	`json:"camera,omitempty"`
	TimeBetweenSensorPollingInSeconds int64 `json:"time_between_sensor_polling_in_seconds,omitempty"`
	ACOutlets  []ACOutlet      `json:"ac_outlets,omitempty"`
}

/**
A DeviceModule is typically an add-on board attached to the edge device that
generates one or more types of measurements.  An AttachedDevice can have multiple
DeviceModules.
*/
type DeviceModule struct {
	ModuleID	int64	`json:"moduleid"`
	ContainerName	string        `json:"container_name,omitempty"`
	ModuleName  string `json:"module_name,omitempty"`
	ModuleType	string `json:"module_type,omitempty"`
	Protocol	string             `json:"protocol,omitempty"`
	Address		string              `json:"address,omitempty"`
	InternalAddress string `json:"internal_address"`
	IncludedSensors []Sensor `json:"included_sensors"`
}

type Sensor struct {
	SensorID 	int64 `json:"sensorid"`
	SensorName	string	`json:"sensor_name"`
	MeasurementName	string `json:"measurement_name"`
}

type AttachedDevice struct {
	DeviceID		int64          `json:"deviceid"`
	DeviceType		string           `json:"device_type,omitempty"`
	DeviceModules 	[]DeviceModule `json:"included_modules,omitempty"`
	ACOutlets  		[]ACOutlet      `json:"ac_outlets,omitempty"`
}

type EnvironmentalTarget struct {
	Temperature float32 `json:"temperature,omitempty"`
	Humidity 	float32 `json:"humidity,omitempty"`
}

type StageSchedule struct {
	Name                 string              `json:"name,omitempty"`
	HoursOfLight         int                 `json:"hours_of_light,omitempty"`
	EnvironmentalTargets EnvironmentalTarget `json:"environmental_targets,omitempty"`
}

type ControlState struct {
}

type PiCam struct {
	PiCamera	bool	`json:"picamera"`
	ResolutionX	int		`json:"resolutionX"`
	ResolutionY	int		`json:"resolutionY"`
}

type Tamper struct {
	Xmove float64	`json:"xmove"`
	Ymove float64	`json:"ymove"`
	Zmove float64	`json:"zmove"`
}

/*
type Configuration1 struct {
	ControllerHostName                string           `json:"controller_hostname"`
	ControllerAPIPort                 int              `json:"controller_api_port"`
	UserID                            int64            `json:"userid"`
	DeviceID                          int64            `json:"deviceid"`
	Stage                             string           `json:"stage,omitempty"`
	LightOnHour                       int              `json:"light_on_hour,omitempty"`
	StageSchedules                    []StageSchedule  `json:"stage_schedules,omitempty"`
	Camera                            PiCam            `json:"camera,omitempty"`
	DeviceSettings                    Station          `json:"device_settings"`
	LogLevel                          string           `json:"log_level,omitempty"`
	AttachedDevices                   []AttachedDevice `json:"edge_devices"`
	AutomaticControl                  bool             `json:"automatic_control"`
	TimeBetweenSensorPollingInSeconds int64            `json:"time_between_sensor_polling_in_seconds"`
	TamperSpec                        Tamper           `json:"tamper"`
}

 */

type ACOutlet struct {
	Name 		string `json:"name,omitempty"`
	Index 		int `json:"index,omitempty"`
	PowerOn 	bool `json:"on,omitempty"`
	BCMPinNumber int `json:"bcm_pin_number,omitempty"`
}

func ReadMyDeviceId(storeMountPoint string, relativePath string, fileName string,) (id int64, err error) {
	log.Debug("ReadMyDeviceId")
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	if relativePath == "" {
		fullpath = storeMountPoint + "/" + fileName
	}
	fmt.Printf("readConfig from %s", fullpath)
	file, _ := ioutil.ReadFile(fullpath)
	idstring := strings.TrimSpace(string(file))

	id, err = strconv.ParseInt(idstring,10, 64)
	return id, err
}

func ReadFromPersistentStore(storeMountPoint string, relativePath string, fileName string, site *Site, currentStageSchedule *StageSchedule) error {
	log.Debug("readConfig")
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	if relativePath == "" {
		fullpath = storeMountPoint + "/" + fileName
	}
	fmt.Printf("readConfig from %s\n", fullpath)
	file, err := ioutil.ReadFile(fullpath)
	if err != nil {
		fmt.Printf("Read config from %s failed %v", fullpath, err )
		return err
	}
	str := string(file)
	log.Infof(str)
	err = json.Unmarshal([]byte(file), site)
	if err != nil {
		fmt.Printf("Error unmarshalling %v\n\n", err )
		fmt.Printf("filestr = %s\n", str )
		return err
	}
	fmt.Printf("data = %v", *site)
	found := false
	for i := 0; i < len(site.Stations) && !found; i++ {
		for j := 0; j < len(site.Stations[i].EdgeDevices) && !found; j++ {
			fmt.Printf("Comparing deviceid %d with %v\n", MyDeviceID, site.Stations[i].EdgeDevices[j])
			if MyDeviceID == site.Stations[i].EdgeDevices[j].DeviceID {
				fmt.Printf("My deviceid %d matches %v\n", MyDeviceID, site.Stations[i].EdgeDevices[j])
				MyStation = &site.Stations[i]
				MyDevice = &site.Stations[i].EdgeDevices[j]
				found = true
			}
		}
	}
	if MyStation == nil {
		log.Fatalf("MyStation not found!!")
	}
	fmt.Printf("MyStation = %v\n", MyStation)
	for i := 0; i < len(MyStation.StageSchedules); i++ {
		if MyStation.StageSchedules[i].Name == MyStation.CurrentStage {
			*currentStageSchedule = MyStation.StageSchedules[i]
			log.Infof("Current stage is %s - schedule is %v", MyStation.CurrentStage, currentStageSchedule)
			return nil
		}
	}
	errstr := fmt.Sprintf("ERROR: No schedule for stage (%s)", MyStation.CurrentStage)
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

func ConfigureLogging( site Site, containerName string) {
	cLog := new(CustomHandler)

	if strings.Contains(site.LogLevel,"error") {
		log.AddHandler(cLog, log.ErrorLevel)
	}
	if strings.Contains(site.LogLevel,"warn") {
		log.AddHandler(cLog, log.WarnLevel)
	}
	if strings.Contains(site.LogLevel,"debug") {
		log.AddHandler(cLog, log.DebugLevel)
	}
	if strings.Contains(site.LogLevel,"info") {
		log.AddHandler(cLog, log.InfoLevel)
	}
	if strings.Contains(site.LogLevel,"notice") {
		log.AddHandler(cLog, log.NoticeLevel)
	}
	if strings.Contains(site.LogLevel,"panic") {
		log.AddHandler(cLog, log.PanicLevel)
	}

}

func validateConfigurable() (err error) {
	if t, ok := interface{}(MySite).(Site); ok == false {
		fmt.Printf(" context %s should be %T, is %T\n", "MySite", t, MySite)
		log.Errorf(" context %s should be %T, is %T", "MySite", t, MySite)
		return errors.New("bad global")
	}
	if t, ok := interface{}(MySite.ControllerHostName).(string); ok == false || len(t) == 0 {
		fmt.Printf(" context %s should be %T, is %T value %s\n", "MySite.ControllerHostName", t, MySite.ControllerHostName, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerHostName", t, MySite.ControllerHostName)
		return errors.New("bad global")
	}
	if t, ok := interface{}(MySite.ControllerAPIPort).(int); ok == false || t <= 0 {
		fmt.Printf(" context %s should be %T, is %T value %d\n", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort)
		return errors.New("bad global")
	}
	if t, ok := interface{}(MyDevice).(*EdgeDevice); ok == false || t == nil {
		fmt.Printf(" context %s should be %T, is %T value %v\n", "MyDevice", t, MyDevice, t)
		log.Errorf(" context %s should be %T, is %T", "MyDevice", t, MyDevice)
		return errors.New("bad global")
	}
	if t, ok := interface{}(MyDevice.DeviceID).(int64); ok == false || t <= 0 {
		fmt.Printf(" context %s should be %T, is %T value %d\n", "MyDevice.DeviceID", t, MyDevice.DeviceID, t)
		log.Errorf(" context %s should be %T, is %T", "MyDevice.DeviceID", t, MyDevice.DeviceID)
		return errors.New("bad global")
	}
	return nil
}
func validateConfigured() (err error) {
	if err := validateConfigurable(); err != nil {
		log.Errorf("validateConfigured error %v", err )
		return err
	}
	if t, ok := interface{}(MySite).(Site); ok == false {
		fmt.Printf(" context %s should be %T, is %T\n", "MySite", t, MySite)
		log.Errorf(" context %s should be %T, is %T", "MySite", t, MySite)
		return errors.New("bad global")
	}
	if t, ok := interface{}(MySite.ControllerHostName).(string); ok == false || len(t) == 0 {
		fmt.Printf(" context %s should be %T, is %T value %s\n", "MySite.ControllerHostName", t, MySite.ControllerHostName, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerHostName", t, MySite.ControllerHostName)
		return errors.New("bad global")
	}
	if t, ok := interface{}(MySite.ControllerAPIPort).(int); ok == false || t <= 0 {
		fmt.Printf(" context %s should be %T, is %T value %d\n", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort)
		return errors.New("bad global")
	}
	return nil
}

func GetConfigFromServer(storeMountPoint string, relativePath string, fileName string) (err error) {
	if err = validateConfigurable(); err != nil {
		log.Errorf("GetConfigFromServer error %v", err )
		return err
	}
	if t, ok := interface{}(storeMountPoint).(string); ok == false {
		log.Errorf(" arg %s should be %v, is %v", "storeMountPoint", t, storeMountPoint)
		return errors.New("bad global")
	}
	if t, ok := interface{}(relativePath).(string); ok == false {
		log.Errorf(" arg %s should be %v, is %v", "relativePath", t, relativePath)
		return errors.New("bad global")
	}
	if t, ok := interface{}(fileName).(string); ok == false {
		log.Errorf(" arg %s should be %v, is %v", "fileName", t, fileName)
		return errors.New("bad global")
	}

	url := fmt.Sprintf("http://%s:%d/api/config/site/%8.8d/%8.8d", MySite.ControllerHostName, MySite.ControllerAPIPort, MySite.UserID, MyDevice.DeviceID)
	fmt.Printf("Sending to %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("post error %v\n", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("readall error %v\n", err)
		return err
	}
	fmt.Printf("response %s\n", string(body))
	newconfig := Site{}
	if err = json.Unmarshal(body, &newconfig); err != nil {
		fmt.Printf("err on site %v\n", err)
		return errors.New("err on site")
	}
	MySite.Stations = newconfig.Stations
	js, _ := json.Marshal(MySite)
	fmt.Printf("\nset site to newconfig \n%s\n", string(js) )
	if err = validateConfigured(); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(MySite, "", "  ")
	filepath := fmt.Sprintf("%s/%s/%s", storeMountPoint,relativePath,fileName)
	if len(relativePath) == 0 {
		filepath = fmt.Sprintf("%s/%s", storeMountPoint,fileName)
	}
	fmt.Printf("writing site to file %s\n\n",filepath)
	err = ioutil.WriteFile(filepath, bytes, 0777)
	if err != nil {
		log.Errorf("error save site file %v", err)
		return(err)
	}

	fmt.Printf("received site\n\n")
	return nil
}
