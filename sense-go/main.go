package main

import (
	"bubblesnet/edge-device/sense-go/adc"
	"io/ioutil"
	"net/http"
	"os"

	//	grpc "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/powerstrip"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	hc "github.com/jdevelop/golang-rpi-extras/sensor_hcsr04"
	"github.com/stianeikeland/go-rpio"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"math"
	"sync"
	"time"
)



func runTamperDetector() {
	log.Info("runTamperDetector")
	adxl345Adaptor := raspi.NewAdaptor()
	adxl345 := i2c.NewADXL345Driver(adxl345Adaptor)
	lastx := 0.0
	lasty := 0.0
	lastz := 0.0

	xmove := 0.0
	ymove := 0.0
	zmove := 0.0

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			x, y, z, _ := adxl345.XYZ()
//			log.Debug(fmt.Sprintf("x: %.7f | y: %.7f | z: %.7f \n", x, y, z))
			if lastx == 0.0 {
			} else {
				xmove = math.Abs(lastx - x)
				ymove = math.Abs(lasty - y)
				zmove = math.Abs(lastz - z)
				if xmove > .03 || ymove > .03 || zmove > .035 {
					log.Info(fmt.Sprintf("TAMPER!! x: %.3f | y: %.3f | z: %.3f ", xmove, ymove, zmove))
					var tamperMessage globals.TamperMessage
					tamperMessage.SampleTimestamp = getNowMillis()
					tamperMessage.XMove = xmove
					tamperMessage.YMove = ymove
					tamperMessage.ZMove = zmove
//					bytearray, err := json.Marshal(tamperMessage)
//					if err != nil {
//						fmt.Println(err)
//						return
//					}

//					msg := bubblesgrpc.SensorRequest{}
//					_, err = bubblesgrpc.SensorStoreAndForwardClient.StoreAndForward(ctx,msg)
//					if err != nil {
//						log.Error(fmt.Sprintf("runTamperDetector ERROR %v", err))
//					}

				} else {
//					log.Debug(fmt.Sprintf("x: %.3f | y: %.3f | z: %.3f \n", xmove, ymove, zmove))
				}
			}
			lastx = x
			lasty = y
			lastz = z
		})
	}

	robot := gobot.NewRobot("adxl345Bot",
		[]gobot.Connection{adxl345Adaptor},
		[]gobot.Device{adxl345},
		work,
	)

	err := robot.Start()
	if err != nil {
		globals.ReportDeviceFailed("adxl345")
		log.Error(fmt.Sprintf("adxl345 robot start error %v", err))
	}
}

func runDistanceWatcher() {
	log.Info("runDistanceWatcher")
	// Use BCM pin numbering
	// Echo pin
	// Trigger pin
	h := hc.NewHCSR04(20, 21)

	for true {
		distance := h.MeasureDistance()
		nanos := distance * 58000.00
		seconds := nanos/1000000000.0
		mydistance := (float64)(17150.00*seconds)
//		log.Debug(fmt.Sprintf("%.2f inches %.2f distance %.2f nanos %.2f cm\n", distance/2.54, distance, nanos, mydistance))
		dm := globals.DistanceMessage{
			SampleTimestamp: getNowMillis(),
			DistanceCm: mydistance,
			DistanceIn: mydistance/2.54}
		bytearray, err := json.Marshal(dm)
		if err == nil {
			log.Debug(fmt.Sprintf("sending distance msg %s?", string(bytearray)))
//			err = grpc.SendStoreAndForwardMessageWithRetries(grpc.GetSequenceNumber(), string(bytearray[:]), 3)
//			if err != nil {
//				log.Error(fmt.Sprintf("runDistanceWatcher ERROR %v", err))
//			}
		} else {
			globals.ReportDeviceFailed("hcsr04")
			log.Error(fmt.Sprintf("rundistancewatcher error = %v", err ))
			break
		}

		time.Sleep(15 * time.Second)
	}
}

func runLocalStateWatcher() {
	log.Info("runLocalStateWatcher")
	for true {
		bytearray, err := json.Marshal(globals.LocalCurrentState)
		if err == nil {
			log.Debug(fmt.Sprintf("sending local current state msg %s?", string(bytearray)))
//			err = grpc.SendStoreAndForwardMessageWithRetries(grpc.GetSequenceNumber(), string(bytearray[:]), 3)
//			if err != nil {
//				log.Error(fmt.Sprintf("runLocalStateWatcher ERROR %v", err))
//			}
		} else {
			log.Debug(fmt.Sprintf("runLocalStateWatcher error = %v", err ))
			break
		}

		time.Sleep(15 * time.Second)
	}
}

/*
func readConfig() error {
	log.Debug(fmt.Sprintf("readglobals.Configuration"))
	file, _ := ioutil.ReadFile("/globals.Configuration/globals.Configuration.json")

	_ = json.Unmarshal([]byte(file), &globals.Config)

	log.Debug(fmt.Sprintf("data = %v", globals.Config ))

	for i := 0; i < len(globals.Config.StageSchedules); i++ {
		if globals.Config.StageSchedules[i].Name == globals.Config.Stage {
			globals.CurrentStageSchedule = globals.Config.StageSchedules[i]
			log.Info(fmt.Sprintf("Current stage is %s - schedule is %v", globals.Config.Stage, globals.CurrentStageSchedule))
			return nil
		}
	}
	log.Error(fmt.Sprintf("ERROR: No schedule for stage %s", globals.Config.Stage))
	return errors.New("No sc:hedule for stage")
}


 */

func makeControlDecisions() {
	log.Info("makeControlDecisions")
	i := 0

	for {
//		gsm := bubblesgrpc.GetStateRequest{}
//		grpc. (gsm)
		if i % 60 == 0 {
			log.Debug(fmt.Sprintf( "LocalCurrentState = %v", globals.LocalCurrentState))
			log.Debug(fmt.Sprintf( "globals.Configuration = %v", globals.Config ))
		}
		ControlLight()
//		turnOnOutletByName(globals.GROWLIGHTVEG)
		ControlHeat()
		ControlHumidity()
//		turnOnOutletByName("Heat lamp")
		time.Sleep(time.Second)
		i++
		if i == 60 {
			i = 0
		}
	}
}

func main() {
	fmt.Printf("sense-go")
	log.Info(fmt.Sprintf("sense-go"))

	err := globals.ReadFromPersistentStore("/go", "", "config.json", &globals.Config, &globals.CurrentStageSchedule)
	globals.ConfigureLogging(globals.Config,"sense-go")
	err = getConfigFromServer()

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")


	log.Info(fmt.Sprintf("globals.Configuration = %v", globals.Config))
	log.Info(fmt.Sprintf("stageSchedule = %v", globals.CurrentStageSchedule))

//	err := readglobals.Configuration()
	if err != nil {
		return
	}

	_ = rpio.Open()
	defer func(){
		err := rpio.Close()
		if err != nil {
			log.Errorf("rpio.close %+v", err)
		}
	}()
	deviceid := 70000007
	if isRelayAttached( deviceid ) {
		log.Infof("Relay is attached to device %d", deviceid)
		powerstrip.InitRpioPins()
		powerstrip.TurnAllOff(1)
	} else {
		log.Infof("There is no relay attached to device %d", deviceid)
	}

	if deviceShouldBeHere(globals.ContainerName,deviceid, globals.Config.DeviceSettings.RootPhSensor,"ezoph") {
		log.Info("Starting Atlas EZO driver")
		ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
		err = ezoDriver.Start()
		if err != nil {
			globals.ReportDeviceFailed("ezoph")
			log.Error(fmt.Sprintf("ezo start error %v", err))
		}
	}

	numGoroutines := 6
	if !globals.Config.DeviceSettings.MovementSensor {
		numGoroutines--
	}
	if !globals.Config.DeviceSettings.WaterLevelSensor {
		numGoroutines--
	}
	if !globals.Config.DeviceSettings.RootPhSensor {
		numGoroutines--
	}
	if !globals.Config.DeviceSettings.HeightSensor {
		numGoroutines--
	}
	if !globals.Config.DeviceSettings.Relay {
		numGoroutines--
	}
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	if deviceShouldBeHere(globals.ContainerName,deviceid, globals.Config.DeviceSettings.MovementSensor, "adxl345") {
		log.Info("MovementSensor should be connected to this device, starting")
		go runTamperDetector()
	} else {
		log.Warn(fmt.Sprint("No adxl345 Configured - skipping tamper detection"))
	}
	if  deviceShouldBeHere(globals.ContainerName,deviceid, globals.Config.DeviceSettings.WaterLevelSensor, "ads1115" ) {
		log.Info("WaterlevelSensor should be connected to this device, starting ADC")
		go func() {
			err :=  adc.RunADCPoller()
			if err != nil {
				log.Errorf("rpio.close %+v", err)
			}
		}()
	} else {
		log.Warn(fmt.Sprint("No ads1115s configured - skipping A to D conversion"))
	}
	if  deviceShouldBeHere(globals.ContainerName,deviceid, globals.Config.DeviceSettings.RootPhSensor, "ezoph" ) {
		log.Info("RootPhSensor should be connected to this device, starting EZO reader")
		go func() {
			err = readPh()
			if err != nil {
				log.Errorf("readPh %+v", err)
			}
		}()
	} else {
		log.Warn(fmt.Sprint("No ezoph configured - skipping pH monitoring"))
	}
	if deviceShouldBeHere(globals.ContainerName, deviceid, globals.Config.DeviceSettings.HeightSensor, "hcsr04" ) {
		log.Info("HeightSensor should be connected to this device, starting HSCR04")
		go runDistanceWatcher()
	} else {
		log.Warn(fmt.Sprint("No hcsr04 Configured - skipping A to D conversion"))
	}
	if isRelayAttached( deviceid ) {
		log.Info("Relay configured")
//		go runPinToggler()
	} else {
		log.Warn(fmt.Sprint("No relay Ccnfigured - skipping GPIO relay control"))
	}

	if len(globals.DevicesFailed) > 0 {
		log.Errorf("Exiting because of device failure %v", globals.DevicesFailed)
		os.Exit(1)
	}
	go runLocalStateWatcher()
	go makeControlDecisions()

	log.Info(fmt.Sprintf("all go routines started, waiting for waitgroup to finish" ))
	wg.Wait()
	log.Info(fmt.Sprintf("exiting main - because waitgroup finished" ))
}

func isRelayAttached( deviceid int ) (relayIsAttached bool){
	for i := 0; i < len(globals.Config.ACOutlets); i++ {
		if globals.Config.ACOutlets[i].DeviceID == deviceid  {
			return true
		}
	}
	return false
}

func deviceShouldBeHere( containerName string, mydeviceid int, deviceInCabinet bool, deviceType string ) ( shouldBePresent bool ) {
	if !deviceInCabinet {
		return false
	}
	for i := 0; i < len(globals.Config.AttachedDevices); i++ {
		if globals.Config.AttachedDevices[i].ContainerName == containerName && globals.Config.AttachedDevices[i].DeviceID == mydeviceid && globals.Config.AttachedDevices[i].DeviceType == deviceType{
			log.Infof("Device %s should be present at %s", globals.Config.AttachedDevices[i].DeviceType, globals.Config.AttachedDevices[i].Address)
			return true
		}
	}
	return false
}

func readPh() error {
	ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
	err := ezoDriver.Start()
	if err != nil {
		log.Error(fmt.Sprintf("ezoDriver.Start returned ph device error %v", err))
		return err
	}
	var e error = nil

	for {
		ph, err := ezoDriver.Ph()
		if err != nil {
			log.Error(fmt.Sprintf("readPh error %v", err))
			e = err
			break
		} else {
			phm := globals.PhMessage{
				SampleTimestamp: getNowMillis(),
				Ph: ph}
//			bytearray, err := json.Marshal(phm)
			_, err := json.Marshal(phm)
			if err == nil {
//				err = grpc.SendStoreAndForwardMessageWithRetries(grpc.GetSequenceNumber(), string(bytearray[:]), 3)
//				if err != nil {
//					log.Error(fmt.Sprintf("readPh ERROR %v", err))
//				}
			} else {
				log.Error(fmt.Sprintf("readph error = %v", err ))
				e = err
				break
			}
		}
		time.Sleep(15*time.Second)
	}
	log.Debug(fmt.Sprintf("returning %v from readph", e ))
	return e
}

func getNowMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func getConfigFromServer() (err error) {
	url := fmt.Sprintf("http://%s:%d/api/config/%8.8d/%8.8d", globals.Config.ControllerHostName, globals.Config.ControllerAPIPort, globals.Config.UserID, globals.Config.DeviceID)
	log.Debugf("Sending to %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("post error %v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("readall error %v", err)
		return err
	}
	log.Debugf("response %s", string(body))
	config, err := json.Marshal(body)
	log.Debugf("received config %v", config)
	return nil
}