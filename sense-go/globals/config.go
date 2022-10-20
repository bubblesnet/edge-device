/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package globals

// copyright and license inspection - no issues 4/13/22

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/log"
	//	_type "google.golang.org/genproto/googleapis/identity/accesscontextmanager/type"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Site is the Top-level object in the data hierarchy.  A site is identified by the user/owner
// and contains multiple stations.
type Site struct {
	SiteID                     int64     `json:"siteid"`
	UserID                     int64     `json:"userid"`
	ControllerAPIHostName      string    // Removed json because hosts come from env. `json:"controller_hostname"`
	ControllerActiveMQHostName string    // Removed json because hosts come from env. `json:"controller_activemq_hostname"`
	ControllerAPIPort          int       // Removed json because hosts come from env. `json:"controller_api_port"`
	ControllerActiveMQPort     int       // Removed json because hosts come from env. `json:"controller_activemq_port"`
	LogLevel                   string    `json:"log_level,omitempty"`
	Stations                   []Station `json:"stations,omitempty"`
}

// AutomationSettings is the set of automation parameters that belong to this station.
type AutomationSettings struct {
	AutomationSettingsID    int64    `json:"automation_settingsid"`
	StageName               string   `json:"stage_name"`
	LightOnStartHour        int      `json:"light_on_start_hour"`
	StageOptions            []string `json:"stage_options"`
	CurrentLightingSchedule string   `json:"current_lighting_schedule"`
	LightingScheduleOptions []string `json:"lighting_schedule_options"`
	TargetTemperature       float64  `json:"target_temperature"`
	TemperatureMin          float64  `json:"temperature_min"`
	TemperatureMax          float64  `json:"temperature_max"`
	TargetWaterTemperature  float64  `json:"target_water_temperature"`
	WaterTemperatureMin     float64  `json:"water_temperature_min"`
	WaterTemperatureMax     float64  `json:"water_temperature_max"`
	HumidityMin             float64  `json:"humidity_min"`
	HumidityMax             float64  `json:"humidity_max"`
	TargetHumidity          float64  `json:"target_humidity"`
	HumidityTargetRangeLow  float64  `json:"humidity_target_range_low"`
	HumidityTargetRangeHigh float64  `json:"humidity_target_range_high"`
	CurrentLightType        string   `json:"current_light_type"`
	LightTypeOptions        []string `json:"light_type_options"`
}

// Station is a grow-unit, typically either a cabinet or a tent.  A station
// contains multiple edge devices, typically Raspberry Pi.  It's the enclosing physical
// structure for one or more plants.
type Station struct {
	StationID              int64              `json:"stationid"`
	AutomaticControl       bool               `json:"automatic_control"`
	HeightSensor           bool               `json:"height_sensor"`
	Humidifier             bool               `json:"humidifier"`
	HumiditySensor         bool               `json:"humidity_sensor_internal"`
	ExternalHumiditySensor bool               `json:"humidity_sensor_external"`
	Heater                 bool               `json:"heater"`
	WaterHeater            bool               `json:"water_heater"`
	ThermometerTop         bool               `json:"thermometer_top"`
	ThermometerMiddle      bool               `json:"thermometer_middle"`
	ThermometerBottom      bool               `json:"thermometer_bottom"`
	ThermometerExternal    bool               `json:"thermometer_external"`
	ThermometerWater       bool               `json:"thermometer_water"`
	WaterPump              bool               `json:"waterPump"`
	AirPump                bool               `json:"airPump"`
	LightSensorInternal    bool               `json:"light_sensor_internal"`
	LightSensorExternal    bool               `json:"light_sensor_external"`
	StationDoorSensor      bool               `json:"station_door_sensor"`
	OuterDoorSensor        bool               `json:"outer_door_sensor"`
	MovementSensor         bool               `json:"movement_sensor"`
	PressureSensor         bool               `json:"pressure_sensors"`
	RootPhSensor           bool               `json:"root_ph_sensor"`
	EnclosureType          string             `json:"enclosure_type,omitempty"`
	WaterLevelSensor       bool               `json:"water_level_sensor"`
	IntakeFan              bool               `json:"intakeFan"`
	ExhaustFan             bool               `json:"exhaustFan"`
	HeatLamp               bool               `json:"heatLamp"`
	HeatingPad             bool               `json:"heatingPad"`
	LightBloom             bool               `json:"lightBloom"`
	LightVegetative        bool               `json:"lightVegetative"`
	LightGerminate         bool               `json:"lightGerminate"`
	Relay                  bool               `json:"relay"`
	EdgeDevices            []EdgeDevice       `json:"edge_devices,omitempty"`
	Dispensers             []Dispenser        `json:"dispensers,omitempty"`
	StageSchedules         []StageSchedule    `json:"stage_schedules,omitempty"`
	TamperSpec             Tamper             `json:"tamper,omitempty"`
	Automation             AutomationSettings `json:"automation_settings"`
	CurrentStage           string             `json:"current_stage"`
	VOCSensor              bool               `json:"voc_sensor"`
	CO2Sensor              bool               `json:"co2_sensor"`
}

// EdgeDevice is a single-board-computer that, with the other
// AttachedDevices in the Station, implements the intelligence of the Station
// such that the ideal grow-conditions for the plants inside the Station are
// always maintained, and a stream of event and environmental sensor messages are sent to
// the time-series database.
type EdgeDevice struct {
	DeviceID                          int64          `json:"deviceid"`
	DeviceType                        string         `json:"devicetypename,omitempty"`
	ExternalID                        string         `json:"externalid,omitempty"`
	IPAddress                         string         `json:"ipaddress,omitempty"`
	MacAddress                        string         `json:"macaddress,omitempty"`
	DeviceModules                     []DeviceModule `json:"modules,omitempty"`
	Camera                            PiCam          `json:"camera,omitempty"`
	TimeBetweenSensorPollingInSeconds int64          `json:"time_between_sensor_polling_in_seconds,omitempty"`
	TimeBetweenPicturesInSeconds      int64          `json:"time_between_pictures_in_seconds"`
	ACOutlets                         []ACOutlet     `json:"ac_outlets,omitempty"`
}

// Dispenser is a combination of a dispenser bottle, a relay-connected dispenser-valve and an additive
// that is loaded into the dispenser bottle.  Combines the server side tables Dispenser and Additive.
type Dispenser struct {
	DispenserID      int64  `json:"dispenserid"`
	DispenserName    string `json:"dispenser_name"`
	DeviceID         int64  `json:"deviceid_device"`
	StationID        int64  `json:"stationid_station"`
	AdditiveID       int64  `json:"currently_loaded_additiveid"`
	OnOff            bool   `json:"onoff"`
	Name             string `json:"name"`
	Index            int    `json:"index"`
	BCMPinNumber     int    `json:"bcm_pin_number"`
	ManufacturerName string `json:"manufacturer_name"`
	AdditiveName     string `json:"additive_name"`
}

// DeviceModule is typically an add-on board attached to the edge device that
// generates one or more types of measurements.  An AttachedDevice can have multiple
// DeviceModules.
type DeviceModule struct {
	ModuleID        int64    `json:"moduleid"`
	ContainerName   string   `json:"container_name,omitempty"`
	ModuleName      string   `json:"module_name,omitempty"`
	ModuleType      string   `json:"module_type,omitempty"`
	Protocol        string   `json:"protocol,omitempty"`
	Address         string   `json:"address,omitempty"`
	InternalAddress string   `json:"internal_address"`
	IncludedSensors []Sensor `json:"included_sensors"`
}

type Sensor struct {
	SensorID        int64  `json:"sensorid"`
	SensorName      string `json:"sensor_name"`
	MeasurementName string `json:"measurement_name"`
}

type AttachedDevice struct {
	DeviceID      int64          `json:"deviceid"`
	DeviceType    string         `json:"device_type,omitempty"`
	DeviceModules []DeviceModule `json:"included_modules,omitempty"`
	ACOutlets     []ACOutlet     `json:"ac_outlets,omitempty"`
}

type EnvironmentalTarget struct {
	Temperature      float32 `json:"temperature"`
	Humidity         float32 `json:"humidity"`
	WaterTemperature float32 `json:"water_temperature"`
}

type StageSchedule struct {
	Name                 string              `json:"name"`
	LightOnStartHour     int                 `json:"light_on_start_hour"`
	HoursOfLight         int                 `json:"hours_of_light"`
	EnvironmentalTargets EnvironmentalTarget `json:"environmental_targets"`
}

type PiCam struct {
	PiCamera    bool `json:"picamera"`
	ResolutionX int  `json:"resolutionX"`
	ResolutionY int  `json:"resolutionY"`
}

type Tamper struct {
	Xmove float64 `json:"xmove"`
	Ymove float64 `json:"ymove"`
	Zmove float64 `json:"zmove"`
}

type ACOutlet struct {
	Name         string `json:"name"`
	Index        int    `json:"index"`
	PowerOn      bool   `json:"on"`
	BCMPinNumber int    `json:"bcm_pin_number"`
}

// ReadMyDeviceId reads the deviceid of this device from the config directory
func ReadMyDeviceId() (id int64, err error) {
	retval, _ := strconv.ParseInt(os.Getenv("DEVICEID"), 10, 64)
	log.Infof("ReadMyDeviceId returning %d", retval)
	return retval, nil
}

// ReadMyAPIServerHostname reads the name/ip of the server from the config directory
func ReadMyAPIServerHostname() (serverHostname string, err error) {
	log.Debug("ReadMyAPIServerHostname")
	return os.Getenv("API_HOST"), nil
}

// ReadCompleteSiteFromPersistentStore reads a complete site configuration from the specified mount-point/relativePath/fileName
// and sets the station, and currentStageSchedule from there.
func ReadCompleteSiteFromPersistentStore(storeMountPoint string, relativePath string, fileName string, site *Site, currentStageSchedule *StageSchedule) error {
	log.Debug("ReadCompleteSiteFromPersistentStore")
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	if relativePath == "" {
		fullpath = storeMountPoint + "/" + fileName
	}
	fmt.Printf("readConfig from %s\n", fullpath)
	file, err := ioutil.ReadFile(fullpath)
	if err != nil {
		fmt.Printf("Read config from %s failed %#v", fullpath, err)
		return err
	}
	str := string(file)
	//	log.Infof(str)

	err = json.Unmarshal([]byte(file), site)
	if err != nil {
		fmt.Printf("Error unmarshalling %#v\n\n", err)
		fmt.Printf("filestr = %s\n", str)
		return err
	}
	success := setMyStationAndMyDevice(MySite)
	if !success {
		if len(site.Stations) == 0 {
			fmt.Printf("NO STATIONS IN THIS SITE!! %#v\n", site)
		} else {
			fmt.Printf("MyStation not found???\n")

		}
		return errors.New(fmt.Sprintf("DeviceID %d not found in %#v", MyDeviceID, site.Stations))
	}

	//	fmt.Printf("my station is set to %#v\n", MyStation)
	//	fmt.Printf("MyStation = %#v\n", MyStation)

	for i := 0; i < len(MyStation.StageSchedules); i++ {
		fmt.Printf("StageSchedule[%d] = %#v\n", i, MyStation.StageSchedules[i])
		if MyStation.StageSchedules[i].Name == MyStation.CurrentStage {
			*currentStageSchedule = MyStation.StageSchedules[i]
			fmt.Printf("Current stage is %s - schedule is %#v", currentStageSchedule.Name, currentStageSchedule)
			return nil
		}
	}
	errstr := fmt.Sprintf("ERROR: No schedule for stage (%s)", MyStation.CurrentStage)

	fmt.Printf("%s\n", errstr)
	log.Error(errstr)
	return errors.New(errstr)
}

// setMyStationAndMyDevice sets the globally accessible vars MyStation and MyDevice from a full
// site configuration.  Convenience function to keep from accessing site every time
func setMyStationAndMyDevice(site Site) (success bool) {
	var nilerr error

	MySite.ControllerAPIHostName, _ = os.Getenv("API_HOST"), nilerr
	MySite.ControllerAPIPort, _ = strconv.Atoi(os.Getenv("API_PORT"))
	MySite.ControllerActiveMQHostName, _ = os.Getenv("ACTIVEMQ_HOST"), nilerr
	MySite.ControllerActiveMQPort, _ = strconv.Atoi(os.Getenv("ACTIVEMQ_PORT"))
	//	fmt.Printf("data = %#v\n", *site)
	found := false
	//	fmt.Printf("searching %d stations\n", len(site.Stations))
	for stationIndex := 0; stationIndex < len(site.Stations) && !found; stationIndex++ {
		//		fmt.Printf("searching %d devices\n", len(site.Stations[stationIndex].EdgeDevices))
		for deviceIndex := 0; deviceIndex < len(site.Stations[stationIndex].EdgeDevices) && !found; deviceIndex++ {
			//			fmt.Printf("Comparing deviceid %d with %#v\n", MyDeviceID, site.Stations[stationIndex].EdgeDevices[deviceIndex])
			if MyDeviceID == site.Stations[stationIndex].EdgeDevices[deviceIndex].DeviceID {
				//				fmt.Printf("My deviceid %d matches %#v\n", MyDeviceID, site.Stations[stationIndex].EdgeDevices[deviceIndex])
				MyStation = &site.Stations[stationIndex]
				log.Infof("MyStation is %#v", MyStation)
				MyDevice = &site.Stations[stationIndex].EdgeDevices[deviceIndex]
				found = true
				return true
			}
		}
	}
	fmt.Printf("Could not set MyStation and MyDevice!!\n")
	return false
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
		_, _ = fmt.Fprintf(b, " %s=%#v", f.Key, f.Value)
	}
	fmt.Println(b.String())
}

// ConfigureLogging adds log handlers for each log level enabled in the site configuration
func ConfigureLogging(site Site, containerName string) {
	cLog := new(CustomHandler)

	if strings.Contains(site.LogLevel, "error") {
		log.AddHandler(cLog, log.ErrorLevel)
	}
	if strings.Contains(site.LogLevel, "warn") {
		log.AddHandler(cLog, log.WarnLevel)
	}
	if strings.Contains(site.LogLevel, "debug") {
		log.AddHandler(cLog, log.DebugLevel)
	}
	if strings.Contains(site.LogLevel, "info") {
		log.AddHandler(cLog, log.InfoLevel)
	}
	if strings.Contains(site.LogLevel, "notice") {
		log.AddHandler(cLog, log.NoticeLevel)
	}
	if strings.Contains(site.LogLevel, "panic") {
		log.AddHandler(cLog, log.PanicLevel)
	}

}

// ValidateConfigurable checks MySite and MyDevice configuration to determine that all configuration
// structures are present and implement the appropriate interfaces.
func ValidateConfigurable() (err error) {
	if t, ok := interface{}(MySite).(Site); ok == false {
		fmt.Printf("ValidateConfigurable mysite.site context %s should be %T, is %T\n", "MySite", t, MySite)
		log.Errorf(" context %s should be %T, is %T", "MySite", t, MySite)
		return errors.New("bad global MySite")
	}
	if t, ok := interface{}(MySite.ControllerAPIHostName).(string); ok == false || len(t) == 0 {
		fmt.Printf("ValidateConfigurable MySite.ControllerAPIHostName context %s should be %T, is %T value %s length %d\n",
			"MySite.ControllerAPIHostName", t, MySite.ControllerAPIHostName, t, len(t))
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerAPIHostName", t, MySite.ControllerAPIHostName)
		fmt.Printf("\n\n%#v\n\n", MySite)
		return errors.New("bad global MySite.ControllerAPIHostName")
	}
	if t, ok := interface{}(MySite.ControllerAPIPort).(int); ok == false || t <= 0 {
		fmt.Printf("ValidateConfigurable MySite.ControllerAPIPort context %s should be %T, is %T value %d\n", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort)
		return errors.New("bad global MySite.ControllerAPIPort")
	}
	if t, ok := interface{}(MyDevice).(*EdgeDevice); ok == false || t == nil {
		fmt.Printf("ValidateConfigurable EdgeDevice context %s should be %T, is %T value %#v\n", "MyDevice", t, MyDevice, t)
		log.Errorf(" context %s should be %T, is %T", "MyDevice", t, MyDevice)
		return errors.New("bad global MyDevice")
	}
	if t, ok := interface{}(MyDevice.DeviceID).(int64); ok == false || t <= 0 {
		fmt.Printf("ValidateConfigurable MyDevice.DeviceID context %s should be %T, is %T value %d\n", "MyDevice.DeviceID", t, MyDevice.DeviceID, t)
		log.Errorf(" context %s should be %T, is %T", "MyDevice.DeviceID", t, MyDevice.DeviceID)
		return errors.New("bad global")
	}
	return nil
}

// ValidateConfigured checks that the configuration is present and implements all interfaces and that all checkable
// values are within limits.
func ValidateConfigured(situation string) (err error) {
	var nilerr error

	MySite.ControllerAPIHostName, _ = os.Getenv("API_HOST"), nilerr
	MySite.ControllerAPIPort, _ = strconv.Atoi(os.Getenv("API_PORT"))
	MySite.ControllerActiveMQHostName, _ = os.Getenv("ACTIVEMQ_HOST"), nilerr
	MySite.ControllerActiveMQPort, _ = strconv.Atoi(os.Getenv("ACTIVEMQ_PORT"))

	if err := ValidateConfigurable(); err != nil {
		log.Errorf("ValidateConfigured error %#v", err)
		fmt.Printf("Validate failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return err
	}
	if t, ok := interface{}(MySite).(Site); ok == false {
		fmt.Printf("ValidateConfigured (%s) (MySite).(Site context %s should be %T, is %T\n", situation, "MySite", t, MySite)
		log.Errorf(" context %s should be %T, is %T", "MySite", t, MySite)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("nil or empty MySite")
	}
	if MySite.SiteID < 0 {
		fmt.Printf("<0 bad MySite.SiteID\n")
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("<0 bad MySite.SiteID")
	}
	if MySite.UserID <= 0 {
		fmt.Printf("<0 MySite.UserID\n")
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("<0 MySite.UserID")
	}
	if t, ok := interface{}(MySite.ControllerAPIHostName).(string); ok == false || len(t) == 0 {
		fmt.Printf("ValidateConfigured (%s) MySite.ControllerAPIHostName context %s should be %T, is %T value %s\n", situation, "MySite.ControllerAPIHostName", t, MySite.ControllerAPIHostName, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerAPIHostName", t, MySite.ControllerAPIHostName)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("nil or wrong type MySite.ControllerAPIHostName")
	}
	if MySite.ControllerAPIHostName == "localhost" {
		fmt.Printf("MySite.ControllerAPIHostName cannot be localhost\n")
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("MySite.ControllerAPIHostName cannot be localhost")
	}
	if t, ok := interface{}(MySite.ControllerAPIPort).(int); ok == false || t <= 0 {
		fmt.Printf("ValidateConfigured (%s) MySite.ControllerAPIPort context %s should be %T, is %T value %d\n", situation, "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort, t)
		log.Errorf(" context %s should be %T, is %T", "MySite.ControllerAPIPort", t, MySite.ControllerAPIPort)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("nil or wrong type ")
	}
	if len(MySite.Stations) <= 0 {
		fmt.Printf("0 length MySite.Stations\n")
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("0 length MySite.Stations")
	}
	if t, ok := interface{}(MyStation).(*Station); ok == false {
		fmt.Printf("ValidateConfigured (%s) (MyStation).(Station context %s should be %T, is %T\n", situation, "*globals.Station", t, MyStation)
		log.Errorf(" context %s should be %T, is %T", "*globals.Station", t, MyStation)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("nil or wrong type MyStation")
	}
	if MyStation == nil || MyStation.StationID < 0 {
		if MyStation == nil {
			fmt.Printf("nil MyStation\n")

		} else {
			fmt.Printf("<0 MyStation.StationID\n")
		}
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("bad MyStation.StationID")
	}
	if MyStation.EnclosureType != "CABINET" && MyStation.EnclosureType != "TENT" {
		fmt.Printf("bad enclosuretype %s\n", MyStation.EnclosureType)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New(fmt.Sprintf("bad MyStation.EnclosureType %s", MyStation.EnclosureType))

	}
	if len(MyStation.StageSchedules) <= 0 {
		fmt.Printf("0 length MyStation.StageSchedules\n")
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("0 length MyStation.StageSchedules")

	}
	if t, ok := interface{}(MyStation.CurrentStage).(string); ok == false || len(t) <= 0 {
		fmt.Printf("nil, or empty MyStation.Automation.CurrentStage\n")
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("nil, or empty MyStation.CurrentStage")

	}
	if MyStation.TamperSpec.Xmove <= 0.0 {
		fmt.Printf("bad TamperSpec.Xmove %f\n", MyStation.TamperSpec.Xmove)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New(fmt.Sprintf("bad TamperSpec.Xmove %f", MyStation.TamperSpec.Xmove))
	}
	if MyStation.TamperSpec.Ymove <= 0.0 {
		fmt.Printf("bad TamperSpec.Ymove %f\n", MyStation.TamperSpec.Ymove)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New(fmt.Sprintf("bad TamperSpec.Ymove %f", MyStation.TamperSpec.Ymove))
	}
	if MyStation.TamperSpec.Zmove <= 0.0 {
		fmt.Printf("bad TamperSpec.Zmove %f\n", MyStation.TamperSpec.Zmove)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New(fmt.Sprintf("bad TamperSpec.Zmove %f", MyStation.TamperSpec.Zmove))
	}

	// is there at least ONE device in the station?
	if t, ok := interface{}(MyStation.EdgeDevices).([]EdgeDevice); ok == false {
		fmt.Printf("ValidateConfigured (%s) (MyStation).(EdgeDevices context %s should be %T, is %T\n", situation, "[]EdgeDevice]", t, MyStation)
		log.Errorf(" context %s should be %T, is %T", "[]EdgeDevice", t, MyStation.EdgeDevices)
		fmt.Printf("ValidateConfigured failed at %s. Sleeping for 1 minute to allow devops container intervention before container restart", situation)
		if situation != "test" {
			time.Sleep(60 * time.Second)
		}
		return errors.New("no edge devices configured in mystation")
	}

	return nil
}

// GetConfigFromServer get the site config from the host named in the configuration file "hostname" in the config
// directory.
func GetConfigFromServer(storeMountPoint string, relativePath string, fileName string) (err error) {
	fmt.Printf("\n\nGetConfigFromServer\n")
	var nilerr error
	MySite.ControllerAPIHostName, _ = os.Getenv("API_HOST"), nilerr
	MySite.ControllerAPIPort, _ = strconv.Atoi(os.Getenv("API_PORT"))
	MySite.ControllerActiveMQHostName, _ = os.Getenv("ACTIVEMQ_HOST"), nilerr
	MySite.ControllerActiveMQPort, _ = strconv.Atoi(os.Getenv("ACTIVEMQ_PORT"))
	if err = ValidateConfigurable(); err != nil {
		log.Errorf("GetConfigFromServer error %#v", err)
		return err
	}
	if t, ok := interface{}(storeMountPoint).(string); ok == false {
		log.Errorf(" arg %s should be %#v, is %#v", "storeMountPoint", t, storeMountPoint)
		return errors.New("bad global")
	}
	if t, ok := interface{}(relativePath).(string); ok == false {
		log.Errorf(" arg %s should be %#v, is %#v", "relativePath", t, relativePath)
		return errors.New("bad global")
	}
	if t, ok := interface{}(fileName).(string); ok == false {
		log.Errorf(" arg %s should be %#v, is %#v", "fileName", t, fileName)
		return errors.New("bad global")
	}

	url := fmt.Sprintf("http://%s:%d/api/config/site/%8.8d/%8.8d", MySite.ControllerAPIHostName, MySite.ControllerAPIPort, MySite.UserID, MyDevice.DeviceID)
	fmt.Printf("Sending to %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("post error %#v\n", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("readall error %#v\n", err)
		return err
	}

	fmt.Printf("\n\nconfig response from server %s\n\n\n", string(body))
	// body is good here
	newconfig := Site{}
	if err = json.Unmarshal(body, &newconfig); err != nil {
		fmt.Printf("err on site %#v\n", err)
		return errors.New("err on site")
	}
	fmt.Printf("newconfig = %+v\n\n", newconfig)
	if newconfig.Stations == nil {
		fmt.Printf("No stations\n")
		log.Fatalf("stations is nil!!!")
	}
	MySite.Stations = newconfig.Stations
	fmt.Printf("before setMyStation %+v", MySite)
	success := setMyStationAndMyDevice(MySite)
	if !success {
		fmt.Printf("No station\n")
		return errors.New("NO station!!")
	}
	fmt.Printf("after setMyStation %+v", MySite)
	//	js, _ := json.Marshal(MySite)
	//	fmt.Printf("\nset site to newconfig \n%s\n", string(js) )
	if err = ValidateConfigured("getConfigFromServer"); err != nil {
		return err
	}
	if err = WriteConfig(storeMountPoint, relativePath, fileName); err != nil {
		return err
	}
	return nil
}

func WriteConfig(storeMountPoint string, relativePath string, fileName string) (err error) {
	fmt.Printf("WriteConfig MySite %+v\n\n", MySite)
	siteBytes, err := json.MarshalIndent(MySite, "", "  ")
	fmt.Printf("\n\nWriteConfig sighx siteBytes %s\n\n", string(siteBytes[:]))
	if err != nil {
		log.Errorf("error marshalling MySite %#v", err)
		return err
	}
	filepath := fmt.Sprintf("%s/%s/%s", storeMountPoint, relativePath, fileName)
	if len(relativePath) == 0 {
		filepath = fmt.Sprintf("%s/%s", storeMountPoint, fileName)
	}
	//	fmt.Printf("writing site to file %s\n\n",filepath)
	err = ioutil.WriteFile(filepath, siteBytes, 0777)
	if err != nil {
		log.Errorf("error save site file %#v", err)
		return err
	}

	fmt.Printf("WriteConfig wrote config\n\n")
	return nil
}
