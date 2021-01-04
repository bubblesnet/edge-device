package main

import (
	log "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/lawg"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

type Config struct {
	Stage          string           `json:"stage,omitempty"`
	LightOnHour    int              `json:"light_on_hour,omitempty"`
	StageSchedules [] StageSchedule `json:"stage_schedules,omitempty"`
	ACOutlets      [8]ACOutlet      `json:"ac_outlets,omitempty"`
	BME280         bool             `json:"bme280,omitempty"`
	BH1750         bool             `json:"bh1750,omitempty"`
	ADS1115_1      bool             `json:"ads1115_1,omitempty"`
	ADS1115_2      bool             `json:"ads1115_2,omitempty"`
	ADXL345        bool             `json:"adxl345,omitempty"`
	EZOPH          bool             `json:"ezoph,omitempty"`
	HCSR04         bool             `json:"hcsr04,omitempty"`
	Relay          bool             `json:"relay,omitempty"`
	LogLevel       string           `json:"log_level,omitempty"`
}

type ACOutlet struct {
	Name string `json:"name,omitempty"`
	Index int `json:"index,omitempty"`
	PowerOn bool `json:"power_on,omitempty"`
	BCMPinNumber int `json:"bcm_pin_number,omitempty"`
}

func ReadFromPersistentStore(storeMountPoint string, relativePath string, fileName string, config *Config, currentStageSchedule *StageSchedule) error {
	log.Debugf("readConfig")
	fullpath := storeMountPoint + "/" + relativePath + "/" + fileName
	file, _ := ioutil.ReadFile(fullpath)

	_ = json.Unmarshal([]byte(file), config)

	log.Debugf("data = %v", *config)
	for i := 0; i < len(config.StageSchedules); i++ {
		if config.StageSchedules[i].Name == config.Stage {
			*currentStageSchedule = config.StageSchedules[i]
			log.Info(fmt.Sprintf("Current stage is %s - schedule is %v", config.Stage, currentStageSchedule))
			return nil
		}
	}
	log.Errorf("ERROR: No schedule for stage %s", config.Stage)
	return errors.New("No sc:hedule for stage")
}

// CustomHandler is your custom handler
type CustomHandler struct {
	// whatever properties you need
}
/*
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
*/
