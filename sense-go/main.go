package main

import (
	"bubblesnet/edge-device/sense-go/adc"
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"bubblesnet/edge-device/sense-go/powerstrip"
	"bubblesnet/edge-device/sense-go/rpio"
	"bubblesnet/edge-device/sense-go/video"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/go-stomp/stomp"
	hc "github.com/jdevelop/golang-rpi-extras/sensor_hcsr04"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math"
	"os"
	"sync"
	"time"
)

var BubblesnetVersionMajorString string
var BubblesnetVersionMinorString = ""
var BubblesnetVersionPatchString = ""
var BubblesnetBuildNumberString = ""
var BubblesnetBuildTimestamp = ""
var BubblesnetGitHash = ""

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
			//			log.Debugf("x: %.7f | y: %.7f | z: %.7f \n", x, y, z))
			if lastx == 0.0 {
			} else {
				xmove = math.Abs(lastx - x)
				ymove = math.Abs(lasty - y)
				zmove = math.Abs(lastz - z)
				if xmove > globals.MyCabinet.TamperSpec.Xmove ||  ymove > globals.MyCabinet.TamperSpec.Ymove ||  zmove > globals.MyCabinet.TamperSpec.Zmove {
					log.Infof("new tamper message !! x: %.3f | y: %.3f | z: %.3f ", xmove, ymove, zmove)
					var tamperMessage = messaging.NewTamperSensorMessage("tamper_sensor",
						0.0, "", "", xmove, ymove, zmove)
					bytearray, err := json.Marshal(tamperMessage)
					if err != nil {
						fmt.Println(err)
						return
					}
					message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
					_, err = globals.Client.StoreAndForward(context.Background(), &message)
					if err != nil {
						log.Errorf("runTamperDetector ERROR %v", err)
					} else {
//						log.Debugf("%v", sensor_reply)
					}

				} else {
					//					log.Debugf("x: %.3f | y: %.3f | z: %.3f \n", xmove, ymove, zmove))
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
		log.Errorf("adxl345 robot start error %v", err)
	}
}

var lastDistance = float64(0.0)

func runDistanceWatcher() {
	log.Info("runDistanceWatcher")
	if globals.RunningOnUnsupportedHardware() {
		return
	}
	// Use BCM pin numbering
	// Echo pin
	// Trigger pin
	h := hc.NewHCSR04(20, 21)

	for true {
		distance := h.MeasureDistance()
		nanos := distance * 58000.00
		seconds := nanos / 1000000000.0
		mydistance := (float64)(17150.00 * seconds)
		direction := ""
		if mydistance > lastDistance {
			direction = "up"
		} else if mydistance < lastDistance {
			direction = "down"
		}
		lastDistance = mydistance
		//		log.Debugf("%.2f inches %.2f distance %.2f nanos %.2f cm\n", distance/2.54, distance, nanos, mydistance))
		dm := messaging.NewDistanceSensorMessage("height_sensor", "plant_height", mydistance, "cm", direction, mydistance, mydistance/2.54)
		bytearray, err := json.Marshal(dm)
		if err == nil {
			log.Debugf("sending distance msg %s?", string(bytearray))
			message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
			_, err := globals.Client.StoreAndForward(context.Background(), &message)
			if err != nil {
				log.Errorf("runDistanceWatcher ERROR %v", err)
			} else {
//				log.Debugf("%v", sensor_reply)
			}
		} else {
			globals.ReportDeviceFailed("hcsr04")
			log.Errorf("rundistancewatcher error = %v", err)
			break
		}
		if globals.RunningOnUnsupportedHardware() {
			return
		}
		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
	}
}

func runLocalStateWatcher() {
	log.Info("runLocalStateWatcher")
	for true {
		bytearray, err := json.Marshal(globals.LocalCurrentState)
		if err == nil {
			log.Debugf("sending local current state msg %s?", string(bytearray))
			//			err = grpc.SendStoreAndForwardMessageWithRetries(grpc.GetSequenceNumber(), string(bytearray[:]), 3)
			//			if err != nil {
			//				log.Error(fmt.Sprintf("runLocalStateWatcher ERROR %v", err))
			//			}
		} else {
			//			log.Debugf("runLocalStateWatcher error = %v", err ))
			break
		}
		if globals.RunningOnUnsupportedHardware() {
			break
		}
		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
	}
}

/*
func readConfig() error {
	log.Debugf("readglobals.Configuration"))
	file, _ := ioutil.ReadFile("/globals.Configuration/globals.Configuration.json")

	_ = json.Unmarshal([]byte(file), &globals.MyFarm)

	log.Debugf("data = %v", globals.MyFarm ))

	for i := 0; i < len(globals.MyFarm.StageSchedules); i++ {
		if globals.MyFarm.StageSchedules[i].Name == globals.MyFarm.Stage {
			globals.CurrentStageSchedule = globals.MyFarm.StageSchedules[i]
			log.Infof("Current stage is %s - schedule is %v", globals.MyFarm.Stage, globals.CurrentStageSchedule))
			return nil
		}
	}
	log.Error(fmt.Sprintf("ERROR: No schedule for stage %s", globals.MyFarm.Stage))
	return errors.New("No schedule for stage")
}


*/
func getNowMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func countACOutlets() int {
	count := 0
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].DeviceID == globals.MyDevice.DeviceID {
			count++
		}
	}
	return count
}

func isMySwitch(switchName string) bool {
	if switchName == "automaticControl" {
		return true
	}

	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].DeviceID == globals.MyDevice.DeviceID {
			if globals.MyDevice.ACOutlets[i].Name == switchName {
				return true
			}
		}
	}
	return false
}


func listenForCommands() (err error) {
	log.Infof("listenForCommands dial")

	var options func(*stomp.Conn) error = func(*stomp.Conn) error{
		stomp.ConnOpt.Login("userid", "userpassword")
		stomp.ConnOpt.Host(globals.MyFarm.ControllerHostName)
		stomp.ConnOpt.RcvReceiptTimeout(30*time.Second)
		stomp.ConnOpt.HeartBeat(30*time.Second, 30*time.Second) // I put this but seems no impact
		return nil
	}

	for j := 0; ; j++ {
		log.Infof("readtimeout dial at %d", getNowMillis())
		host_port := fmt.Sprintf("%s:%d", globals.MyFarm.ControllerHostName, 61613)
		stompConn, err := stomp.Dial("tcp", host_port, options)
		if err != nil {
			log.Errorf("listenForCommands dial error %v", err)
			return err
		}
		log.Infof("listenForCommands connect")
		defer stompConn.Disconnect()

		topicName := fmt.Sprintf("/topic/%8.8d/%8.8d", globals.MyFarm.UserID, globals.MyDevice.DeviceID)
		log.Infof("listenForCommands subscribe to topic %s", topicName)

		sub, err := stompConn.Subscribe(topicName, stomp.AckClient)
		if err != nil {
			log.Infof("readtimeout error at %d", getNowMillis())
			log.Errorf("listenForCommands subscribe error %v", err)
			return err
		}
		//
		for i := 0; ; i++ {
			//		log.Infof("listenForCommands read %d", i)
			msg := <-sub.C
			if msg == nil || msg.Err != nil {
				if msg != nil && msg.Err != nil {
					log.Errorf("listenForCommands read topic error %v", msg.Err)
					time.Sleep(2 * time.Second)
					break
				} else {
					//				log.Errorf("listenForCommands read topic error %v", msg)
				}
				time.Sleep(2 * time.Second)
				continue
			}
			type MessageHeader struct {
				Command string `json:"command"`
			}
			type SwitchMessage struct {
				Command    string `json:"command"`
				SwitchName string `json:"switch_name"`
				On         bool   `json:"on"`
			}

			header := MessageHeader{}
			err = json.Unmarshal(msg.Body, &header)
			if err != nil {
				log.Errorf("listenForCommands marshal error %v", err)
				continue
			}
			log.Infof("listenForCommands parsed body into %v", header)
			log.Infof("header.Command === %s", header.Command)
			switch header.Command {
			case "picture":
				if globals.MyDevice.Camera.PiCamera == false{
					log.Infof("No camera configured, skipping picture")
				} else {
					log.Infof("switch calling takeAPicture")
					video.TakeAPicture()
				}
				break
			case "switch":
				{
				if countACOutlets() == 0 {
					log.Infof("No ac outlets configured on this device")
					break
				}
					switchMessage := SwitchMessage{}
					err := json.Unmarshal(msg.Body, &switchMessage)
					log.Infof("listenForCommands parsed body into SwitchMessage %v", switchMessage)
					if err != nil {
						log.Errorf("listenForCommands switch error %v", err)
						break
					}
					if !isMySwitch( switchMessage.SwitchName) {
						log.Infof("Not my switch %s", switchMessage.SwitchName)
						break;
					}
					if switchMessage.SwitchName == "automaticControl" {
						log.Infof("listenForCommands setting %s to %v", switchMessage.SwitchName, switchMessage.On)
						globals.MyFarm.AutomaticControl = switchMessage.On
						if globals.MyFarm.AutomaticControl {
							initializeOutletsForAutomation()	// Make sure the switches conform to currently configured automation
						}
					} else if switchMessage.On == true {
						log.Infof("listenForCommands turning on %s", switchMessage.SwitchName)
						powerstrip.TurnOnOutletByName(switchMessage.SwitchName, true)
					} else {
						log.Infof("listenForCommands turning off %s", switchMessage.SwitchName)
						powerstrip.TurnOffOutletByName(switchMessage.SwitchName, true)
					}
					break
				}
			default:
				{
					break
				}
			}

			log.Infof("listenForCommands received message %s", string(msg.Body))
		}
	}
	log.Infof("listenForCommands returning")
	return nil
}

func initializeOutletsForAutomation() {
	ControlLight(true)
	ControlHeat(true)
	ControlHumidity(true)
	ControlOxygenation(true )
	ControlRootWater(true )
	ControlAirflow(true )
}

func makeControlDecisions() {
	log.Info("makeControlDecisions")
	i := 0

	for {
		//		gsm := bubblesgrpc.GetStateRequest{}
		//		grpc. (gsm)
		if i%60 == 0 {
			log.Debugf("LocalCurrentState = %v", globals.LocalCurrentState)
//			log.Debugf("globals.Configuration = %v", globals.MyFarm)
		}
		if globals.MyFarm.AutomaticControl {
			ControlLight(false)
			ControlHeat(false)
			ControlHumidity(false)
			ControlOxygenation(false )
			ControlRootWater(false )
			ControlAirflow(false )
			if globals.RunningOnUnsupportedHardware() {
				return
			}
		}
		time.Sleep(time.Second)
		i++
		if i == 60 {
			i = 0
		}
	}
}

func reportVersion() {
	log.Infof("Version %s.%s.%s-%s timestamp %s githash %s", BubblesnetVersionMajorString, BubblesnetVersionMinorString, BubblesnetVersionPatchString,
		BubblesnetBuildNumberString, BubblesnetBuildTimestamp, BubblesnetGitHash)
}

func main() {
	fmt.Printf(globals.ContainerName)
	log.Infof(globals.ContainerName)

	globals.BubblesnetVersionMajorString = BubblesnetVersionMajorString
	globals.BubblesnetVersionMinorString = BubblesnetVersionMinorString
	globals.BubblesnetVersionPatchString = BubblesnetVersionPatchString
	globals.BubblesnetBuildNumberString = BubblesnetBuildNumberString
	globals.BubblesnetBuildTimestamp = BubblesnetBuildTimestamp
	globals.BubblesnetGitHash = BubblesnetGitHash

	var err error
	globals.MyDeviceID, err = globals.ReadMyDeviceId("/config","", "deviceid")
	if err != nil {
		fmt.Printf("error read device %v\n", err)
		return
	}
	fmt.Printf("Read deviceid %d\n", globals.MyDeviceID)
	if err := globals.ReadFromPersistentStore("/config", "", "config.json", &globals.MyFarm, &globals.CurrentStageSchedule); err != nil {
		return
	}
	fmt.Printf("\ngetconfigfromserver config = %v\n\n", globals.MyFarm)
	if err := globals.GetConfigFromServer("/config", "", "config.json"); err != nil {
		return
	}
	fmt.Printf("done getting config from server %v\n\n", globals.MyFarm)
	globals.ConfigureLogging(globals.MyFarm, "sense-go")

//	globals.MyFarm.Cabinet.HeightSensor = true
	reportVersion()

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")

	log.Infof("globals.Configuration = %v", globals.MyFarm)
	log.Infof("stageSchedule = %v", globals.CurrentStageSchedule)

	// Set up a connection to the server.
	log.Infof("Dialing GRPC server at %s", globals.ForwardingAddress)
	conn, err := grpc.Dial(globals.ForwardingAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	globals.Client = pb.NewSensorStoreAndForwardClient(conn)

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
/*
	log.Info("Calling rpio.open")
	_ = rpio.Open()
	defer func() {
		err := rpio.Close()
		if err != nil {
			log.Errorf("rpio.close %+v", err)
		}
	}()
 */
	rpio.OpenRpio()
	if isRelayAttached(globals.MyDevice.DeviceID) {
		log.Infof("Relay is attached to device %d", globals.MyDevice.DeviceID)
		powerstrip.InitRpioPins()
		if globals.MyFarm.AutomaticControl {
			powerstrip.TurnAllOff(1)	// turn all OFF first since initalizeOutlets doesnt
			initializeOutletsForAutomation()
		} else {
			powerstrip.TurnAllOff(1)
		}
	} else {
		log.Infof("There is no relay attached to device %d", globals.MyDevice.DeviceID)
	}
	powerstrip.SendSwitchStatusChangeEvent("automaticControl",globals.MyFarm.AutomaticControl)

	log.Info("ezo")
	if deviceShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.RootPhSensor, "ezoph") {
		log.Info("Starting Atlas EZO driver")
		ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
		err = ezoDriver.Start()
		if err != nil {
			globals.ReportDeviceFailed("ezoph")
			log.Errorf("ezo start error %v", err)
		}
	} else {
		log.Infof("No root ph sensor configured")
	}
	log.Info("after ezo")

	numGoroutines := 6
	if !globals.MyCabinet.MovementSensor {
		numGoroutines--
	}
	if !globals.MyCabinet.WaterLevelSensor {
		numGoroutines--
	}
	if !globals.MyCabinet.RootPhSensor {
		numGoroutines--
	}
	if !globals.MyCabinet.HeightSensor {
		numGoroutines--
	}
	if !globals.MyCabinet.Relay {
		numGoroutines--
	}
	log.Infof("Waiting for %d goroutines", numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	log.Info("movement")
	if deviceShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.MovementSensor, "adxl345") {
		log.Info("MovementSensor should be connected to this device, starting")
		go runTamperDetector()
	} else {
		log.Warnf("No adxl345 Configured - skipping tamper detection")
	}
	log.Infof("adc %s %d %v ads1115", globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.WaterLevelSensor)
	if deviceShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.WaterLevelSensor, "ads1115") {
		log.Info("WaterlevelSensor should be connected to this device, starting ADC")
		go func() {
			err := adc.RunADCPoller()
			if err != nil {
				log.Errorf("rpio.close %+v", err)
			}
		}()
	} else {
		log.Warnf("No ads1115s configured - skipping A to D conversion")
	}
	log.Info("root ph")
	if deviceShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.RootPhSensor, "ezoph") {
		log.Info("RootPhSensor should be connected to this device, starting EZO reader")
		go func() {
			err = readPh()
			if err != nil {
				log.Errorf("readPh %+v", err)
			}
		}()
	} else {
		log.Warnf("No ezoph configured - skipping pH monitoring")
	}
	log.Infof("deviceShouldBeHere %s %d %v hcsr04", globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.HeightSensor)
	if deviceShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyCabinet.HeightSensor, "hcsr04") {
		log.Info("HeightSensor should be connected to this device, starting HSCR04")
		go runDistanceWatcher()
	} else {
		log.Warnf("No hcsr04 Configured - skipping distance monitoring")
	}
	if isRelayAttached(globals.MyDevice.DeviceID) {
		log.Info("Relay configured")
		//		go runPinToggler()
	} else {
		log.Warnf("No relay Ccnfigured - skipping GPIO relay control")
	}

	if len(globals.DevicesFailed) > 0 {
		log.Errorf("Exiting because of device failure %v", globals.DevicesFailed)
		os.Exit(1)
	}

	go makeControlDecisions()

	go func() {
		err = listenForCommands()
		if err != nil {
			log.Errorf("listenForCommands %+v", err)
		}
	}()

	log.Infof("all go routines started, waiting for waitgroup to finish")
	wg.Wait()
	log.Infof("exiting main - because waitgroup finished")
}

func isRelayAttached(deviceid int64) (relayIsAttached bool) {
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].DeviceID == deviceid {
			return true
		}
	}
	return false
}

func deviceShouldBeHere(containerName string, mydeviceid int64, deviceInCabinet bool, deviceType string) (shouldBePresent bool) {
	if !deviceInCabinet {
		return false
	}
	for i := 0; i < len(globals.MyCabinet.AttachedDevices); i++ {
		//		log.Infof("%v", globals.MyFarm.AttachedDevices[i])
		if globals.MyCabinet.AttachedDevices[i].ContainerName == containerName && globals.MyCabinet.AttachedDevices[i].DeviceID == mydeviceid && globals.MyCabinet.AttachedDevices[i].DeviceType == deviceType {
			log.Infof("Device %s should be present at %s", globals.MyCabinet.AttachedDevices[i].DeviceType, globals.MyCabinet.AttachedDevices[i].Address)
			return true
		}
	}
	return false
}

var lastPh = float64(0.0)

func readPh() error {
	ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
	err := ezoDriver.Start()
	if err != nil {
		log.Errorf("ezoDriver.Start returned ph device error %v", err)
		return err
	}
	var e error = nil

	for {
		ph, err := ezoDriver.Ph()
		if err != nil {
			log.Errorf("readPh error %v", err)
			e = err
			break
		} else {
			direction := ""
			if ph > lastPh {
				direction = "up"
			} else if ph < lastPh {
				direction = "down"
			}
			lastPh = ph
			phm := messaging.NewGenericSensorMessage("root_ph_sensor", "root_ph", ph, "", direction)
			bytearray, err := json.Marshal(phm)
			message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
			_, err = globals.Client.StoreAndForward(context.Background(), &message)
			if err != nil {
				log.Errorf("RunADCPoller ERROR %v", err)
			} else {
//				log.Infof("sensor_reply %v", sensor_reply)
			}
		}
		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
	}
	log.Debugf("returning %v from readph", e)
	return e
}

