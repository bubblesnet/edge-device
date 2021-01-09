package globals

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/log"
	"io/ioutil"
	"strings"
)

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

type DeviceSettings struct {
	HeightSensor bool `json:"height_sensor"`
	Humidifier bool `json:"humidifier"`
	HumiditySensor bool `json:"humidity_sensor"`
	ExternalHumiditySensor bool `json:"external_humidity_sensor"`
	Heater bool `json:"heater"`
	ThermometerTop bool `json:"thermometer_top"`
	ThermometerMiddle bool `json:"thermometer_middle"`
	ThermometerBottom bool `json:"thermometer_bottom"`
	ThermometerExternal bool `json:"external_thermometer"`
	ThermometerWater bool `json:"thermometer_water"`
	WaterPump bool `json:"water_pump"`
	AirPump bool `json:"air_pump"`
	LightSensor bool `json:"light_sensor"`
	CabinetDoorSensor bool `json:"cabinet_door_sensor"`
	OuterDoorSensor bool `json:"outer_door_sensor"`
	MovementSensor bool `json:"movement_sensor"`
	PressureSensor bool `json:"pressure_sensors"`
	RootPhSensor bool `json:"root_ph_sensor"`
	EnclosureType string `json:"enclosure_type"`
	WaterLevelSensor bool `json:"water_level_sensor"`
	IntakeFan bool `json:"intake_fan"`
	ExhaustFan bool `json:"exhaust_fan"`
	HeatLamp bool `json:"heat_lamp"`
	HeatingPad bool `json:"heating_pad"`
	LightBloom bool `json:"light_bloom"`
	LightVegetative bool `json:"light_vegetative"`
	LightGerminate bool `json:"light_germinate"`
	Relay          bool `json:"relay,omitempty"`
}

type Configuration struct {
	ControllerHostName	string		`json:"controller_hostname"`
	ControllerAPIPort	int			`json:"controller_api_port"`
	UserID			int			`json:"userid"`
	DeviceID		int			`json:"deviceid"`
	Stage          string          `json:"stage,omitempty"`
	LightOnHour    int             `json:"light_on_hour,omitempty"`
	StageSchedules []StageSchedule `json:"stage_schedules,omitempty"`
	ACOutlets      [8]ACOutlet     `json:"ac_outlets,omitempty"`
	DeviceSettings	DeviceSettings	`json:"device_settings"`
	LogLevel       string          `json:"log_level,omitempty"`
}

type ACOutlet struct {
	Name string `json:"name,omitempty"`
	Index int `json:"index,omitempty"`
	PowerOn bool `json:"power_on,omitempty"`
	BCMPinNumber int `json:"bcm_pin_number,omitempty"`
}

func ReadFromPersistentStore(storeMountPoint string, relativePath string, fileName string, config *Configuration, currentStageSchedule *StageSchedule) error {
	log.Debug(fmt.Sprintf("readConfig"))
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	fmt.Printf("readConfig from %s", fullpath)
	file, _ := ioutil.ReadFile(fullpath)

	_ = json.Unmarshal([]byte(file), config)

	log.Debug(fmt.Sprintf("data = %v", *config))
	for i := 0; i < len(config.StageSchedules); i++ {
		if config.StageSchedules[i].Name == config.Stage {
			*currentStageSchedule = config.StageSchedules[i]
			log.Info(fmt.Sprintf("Current stage is %s - schedule is %v", config.Stage, currentStageSchedule))
			return nil
		}
	}
	log.Error(fmt.Sprintf("ERROR: No schedule for stage %s", config.Stage))
	return errors.New("No sc:hedule for stage")
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

func ConfigureLogging( config Configuration, containerName string) {
	cLog := new(CustomHandler)

	if strings.Contains(config.LogLevel,"error") {
		log.AddHandler(cLog, log.ErrorLevel)
	}
	if strings.Contains(config.LogLevel,"warn") {
		log.AddHandler(cLog, log.WarnLevel)
	}
	if strings.Contains(config.LogLevel,"debug") {
		log.AddHandler(cLog, log.DebugLevel)
	}
	if strings.Contains(config.LogLevel,"info") {
		log.AddHandler(cLog, log.InfoLevel)
	}
	if strings.Contains(config.LogLevel,"notice") {
		log.AddHandler(cLog, log.NoticeLevel)
	}
	if strings.Contains(config.LogLevel,"panic") {
		log.AddHandler(cLog, log.PanicLevel)
	}

}

