package main

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/a2dconverter"
	"bubblesnet/edge-device/sense-go/modules/accelerometer"
	"bubblesnet/edge-device/sense-go/modules/camera"
	"bubblesnet/edge-device/sense-go/modules/distancesensor"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"bubblesnet/edge-device/sense-go/modules/phsensor"
	"bubblesnet/edge-device/sense-go/modules/rpio"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/go-stomp/stomp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"strings"
	"sync"
	"time"
)

var BubblesnetVersionMajorString string
var BubblesnetVersionMinorString = ""
var BubblesnetVersionPatchString = ""
var BubblesnetBuildNumberString = ""
var BubblesnetBuildTimestamp = ""
var BubblesnetGitHash = ""

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
		delay := globals.MyDevice.TimeBetweenSensorPollingInSeconds
		if globals.RunningOnUnsupportedHardware() {
			delay = 1
		}
		time.Sleep(time.Duration(delay) * time.Second)
		if globals.RunningOnUnsupportedHardware() {
			break
		}
	}
}

func getNowMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func countACOutlets() int {
	return len(globals.MyDevice.ACOutlets)
}

func processCommand(msg *stomp.Message) (resub bool, err error) {
	if msg == nil || msg.Err != nil {
		if msg != nil && msg.Err != nil {
			if strings.Contains(fmt.Sprintf("%v", msg.Err), "timeout") {
				log.Debugf("queue read timed out - resubscribing %v", msg.Err)
			} else {
				log.Errorf("listenForCommands read topic error %v", msg.Err)
			}
			time.Sleep(2 * time.Second)
			return true, msg.Err
		} else {
			//				log.Errorf("listenForCommands read topic error %v", msg)
		}
		time.Sleep(2 * time.Second)
		return false, nil
	}
	type MessageHeader struct {
		Command string `json:"command"`
	}
	type SwitchMessage struct {
		Command    string `json:"command"`
		SwitchName string `json:"switch_name"`
		On         bool   `json:"on"`
	}
	type StageMessage struct {
		Command   string `json:"command"`
		StageName string `json:"stage_name"`
	}
	header := MessageHeader{}
	err = json.Unmarshal(msg.Body, &header)
	if err != nil {
		log.Errorf("listenForCommands marshal error %v", err)
		return false, err
	}
	log.Infof("listenForCommands parsed body into %v", header)
	log.Infof("header.Command === %s", header.Command)
	switch header.Command {
	case "stage":
		log.Infof("Changing stage via message %s", msg.Body)
		stageMessage := StageMessage{}
		if err := json.Unmarshal(msg.Body, &stageMessage); err != nil {
			log.Errorf("couldn't parse stage message %s, %v", msg.Body, err)
			break
		}
		log.Infof("listenForCommands parsed body into SwitchMessage %v", stageMessage)
		globals.MyStation.CurrentStage = stageMessage.StageName
		break
	case "picture":
		if globals.MyDevice.Camera.PiCamera == false {
			log.Infof("No camera configured, skipping picture")
		} else {
			log.Infof("switch calling takeAPicture")
			camera.TakeAPicture()
		}
		break
	case "status":
		fmt.Printf("\n\nReceived status message\n\n")
		gpiorelay.PowerstripSvc.ReportAll(200 * time.Millisecond)
		gpiorelay.PowerstripSvc.SendSwitchStatusChangeEvent("automaticControl", globals.MyStation.AutomaticControl)
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
			if !gpiorelay.PowerstripSvc.IsMySwitch(switchMessage.SwitchName) {
				log.Infof("Not my switch %s", switchMessage.SwitchName)
				break
			}
			if switchMessage.SwitchName == "automaticControl" {
				log.Infof("listenForCommands setting %s to %v", switchMessage.SwitchName, switchMessage.On)
				globals.MyStation.AutomaticControl = switchMessage.On
				if globals.MyStation.AutomaticControl {
					initializeOutletsForAutomation() // Make sure the switches conform to currently configured automation
				}
			} else if switchMessage.On == true {
				log.Infof("listenForCommands turning on %s", switchMessage.SwitchName)
				gpiorelay.GetPowerstripService().TurnOnOutletByName(switchMessage.SwitchName, true)
			} else {
				log.Infof("listenForCommands turning off %s", switchMessage.SwitchName)
				gpiorelay.GetPowerstripService().TurnOffOutletByName(switchMessage.SwitchName, true)
			}
			break
		}
	default:
		{
			break
		}
	}
	log.Infof("listenForCommands received message %s", string(msg.Body))
	return false, nil
}

func listenForCommands(isUnitTest bool) (err error) {
	log.Infof("listenForCommands dial")

	var options func(*stomp.Conn) error = func(*stomp.Conn) error {
		stomp.ConnOpt.Login("userid", "userpassword")
		stomp.ConnOpt.Host(globals.MySite.ControllerHostName)
		stomp.ConnOpt.RcvReceiptTimeout(30 * time.Second)
		stomp.ConnOpt.HeartBeat(30*time.Second, 30*time.Second) // I put this but seems no impact
		return nil
	}

	for j := 0; ; j++ {
		log.Debugf("stomp.Dial at %d", getNowMillis())
		host_port := fmt.Sprintf("%s:%d", globals.MySite.ControllerHostName, 61613)
		stompConn, err := stomp.Dial("tcp", host_port, options)
		if err != nil {
			log.Errorf("listenForCommands dial error %v", err)
			return err
		}
		log.Infof("listenForCommands connect")
		defer stompConn.Disconnect()

		topicName := fmt.Sprintf("/topic/%8.8d/%8.8d", globals.MySite.UserID, globals.MyDevice.DeviceID)
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
			reSubscribe, err := processCommand(msg)
			if err != nil {
				log.Errorf("processCommand error %v", err)
			}
			if reSubscribe {
				break
			}

		}
	}
	log.Infof("listenForCommands returning")
	return nil
}

func initializeOutletsFromConfiguration() {
	ps := gpiorelay.GetPowerstripService()
	for i := 0; i < len(globals.MyDevice.ACOutlets); i++ {
		if globals.MyDevice.ACOutlets[i].PowerOn {
			ps.TurnOnOutletByName(globals.MyDevice.ACOutlets[i].Name, true)
		} else {
			ps.TurnOffOutletByName(globals.MyDevice.ACOutlets[i].Name, true)
		}
	}
}

func initializeOutletsForAutomation() {
	ControlLight(true)
	ControlHeat(true)
	ControlHumidity(true)
	ControlOxygenation(true)
	ControlRootWater(true)
	ControlAirflow(true)
}

func makeControlDecisions(once_only bool) {
	log.Info("makeControlDecisions endless loop with once_only set to %v", once_only)
	i := 0

	for {
		gsm := pb.GetStateRequest{}
		gsm.Sequence = globals.GetSequence()
		gr, err := globals.Client.GetState(context.Background(), &gsm)
		if err != nil {
			log.Errorf("getState got error %v", err)
		} else {
			globals.ExternalCurrentState.TempF = gr.TempF
			globals.ExternalCurrentState.Humidity = gr.Humidity
		}
		log.Infof("Got state TempF %f Humidity %f", gr.TempF, gr.Humidity)

		if globals.MyStation.AutomaticControl {
			ControlLight(false)
			ControlHeat(false)
			ControlHumidity(false)
			ControlOxygenation(false)
			ControlRootWater(false)
			ControlAirflow(false)
		}
		time.Sleep(time.Second)
		i++
		if i >= 60 {
			i = 0
		}
		if once_only {
			break
		}
	}
	log.Infof("makeControlDecisions returning")
}

func reportVersion() {
	log.Infof("Version %s.%s.%s-%s timestamp %s githash %s", BubblesnetVersionMajorString, BubblesnetVersionMinorString, BubblesnetVersionPatchString,
		BubblesnetBuildNumberString, BubblesnetBuildTimestamp, BubblesnetGitHash)
}

func initGlobals() {
	log.Infof("initGlobals")
	globals.BubblesnetVersionMajorString = BubblesnetVersionMajorString
	globals.BubblesnetVersionMinorString = BubblesnetVersionMinorString
	globals.BubblesnetVersionPatchString = BubblesnetVersionPatchString
	globals.BubblesnetBuildNumberString = BubblesnetBuildNumberString
	globals.BubblesnetBuildTimestamp = BubblesnetBuildTimestamp
	globals.BubblesnetGitHash = BubblesnetGitHash

	var err error
	globals.MyDeviceID, err = globals.ReadMyDeviceId(globals.PersistentStoreMountPoint, "", "deviceid")
	if err != nil {
		fmt.Printf("error read device %v\n", err)
		return
	}
	globals.MySite.ControllerHostName, err = globals.ReadMyServerHostname(globals.PersistentStoreMountPoint, "", "hostname")
	if err != nil {
		fmt.Printf("error read serverHostname %v\n", err)
		return
	}
	// Read the configuration file
	fmt.Printf("Read deviceid %d and server_hostname %s\n", globals.MyDeviceID, globals.MySite.ControllerHostName)
	readConfigFromDisk()
	// Get a NEW config file from server and save to disk
	if err := globals.GetConfigFromServer(globals.PersistentStoreMountPoint, "", "config.json"); err != nil {
		fmt.Printf("Exiting because of bad configuration - sleeping for 60 seconds to allow intervention\n")
		time.Sleep(60 * time.Second)
		os.Exit(1)
	}
	// Reread the configuration file
	readConfigFromDisk()

	globals.MySite.LogLevel = "silly,debug,info,warn,fatal,notice,error,alert"
	//	fmt.Printf("done getting config from server %v\n\n", globals.MySite)
	globals.ConfigureLogging(globals.MySite, "sense-go")

	//	globals.MySite.Station.HeightSensor = true
	reportVersion()

	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")

	//	log.Infof("globals.Configuration = %v", globals.MySite)
	//	log.Infof("stageSchedule = %v", globals.CurrentStageSchedule)
}

func readConfigFromDisk() {
	if err := globals.ReadFromPersistentStore(globals.PersistentStoreMountPoint, "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
		fmt.Printf("ReadFromPersistentStore failed - using default config\n")
		//		globals.MySite.ControllerHostName = serverHostname
		globals.MySite.ControllerAPIPort = 3003
		nodeEnv := os.Getenv("NODE_ENV")
		switch nodeEnv {
		case "PRODUCTION":
			globals.MySite.ControllerAPIPort = 3001
			break
		case "DEV":
			globals.MySite.ControllerAPIPort = 3003
			break
		case "TEST":
			globals.MySite.ControllerAPIPort = 3002
			break
		}
		globals.MySite.UserID = 90000009
		d := globals.EdgeDevice{DeviceID: globals.MyDeviceID}
		globals.MyDevice = &d
		//		fmt.Printf("\ngetconfigfromserver config = %v\n\n", globals.MySite)
	}
}

func setupGPIO() {
	rpio.OpenRpio()
	/*	defer func() {
		rpio.CloseRpio()
	}() */

	if isRelayAttached(globals.MyDevice.DeviceID) {
		log.Infof("Relay is attached to device %d", globals.MyDevice.DeviceID)
		gpiorelay.PowerstripSvc.InitRpioPins()
		gpiorelay.PowerstripSvc.TurnAllOff(1) // turn all OFF first since initalizeOutlets doesnt
		if globals.MyStation.AutomaticControl {
			initializeOutletsForAutomation()
		} else {
			initializeOutletsFromConfiguration()
		}
		gpiorelay.PowerstripSvc.SendSwitchStatusChangeEvent("automaticControl", globals.MyStation.AutomaticControl)
	} else {
		log.Infof("There is no relay attached to device %d", globals.MyDevice.DeviceID)
	}
}

func setupPhMonitor() {
	log.Infof("ezo mydevice %v, mystation %v", globals.MyDevice, globals.MyStation)
	globals.ValidateConfigured("setupPhMonitor")
	if moduleShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.RootPhSensor, "ezoph") {
		phsensor.StartEzoDriver()
	} else {
		log.Infof("No root ph sensor configured")
	}
	log.Info("after ezo")
}

func countGoRoutines() (count int) {
	numGoroutines := 6
	if !globals.MyStation.MovementSensor {
		numGoroutines--
	}
	if !globals.MyStation.WaterLevelSensor {
		numGoroutines--
	}
	if !globals.MyStation.RootPhSensor {
		numGoroutines--
	}
	if !globals.MyStation.HeightSensor {
		numGoroutines--
	}
	if !globals.MyStation.Relay {
		numGoroutines--
	}
	log.Infof("Waiting for %d goroutines", numGoroutines)
	return numGoroutines
}

func startGoRoutines(onceOnly bool) {
	log.Info("movement")
	if moduleShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.MovementSensor, "adxl345") {
		log.Info("MovementSensor should be connected to this device, starting")
		go accelerometer.RunTamperDetector(onceOnly)
	} else {
		log.Warnf("No adxl345 Configured - skipping tamper detection")
	}
	log.Infof("adc %s %d %v ads1115", globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.WaterLevelSensor)
	if moduleShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.WaterLevelSensor, "ads1115") {
		log.Info("WaterlevelSensor should be connected to this device, starting ADC")
		go func() {
			err := a2dconverter.RunADCPoller(onceOnly)
			if err != nil {
				log.Errorf("rpio.close %+v", err)
			}
		}()
	} else {
		log.Warnf("No ads1115s configured - skipping A to D conversion because %v", globals.MyStation.WaterLevelSensor)
	}
	log.Info("root ph")
	if moduleShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.RootPhSensor, "ezoph") {
		phsensor.StartEzo(onceOnly)
	} else {
		log.Warnf("No ezoph configured - skipping pH monitoring")
	}
	log.Infof("moduleShouldBeHere %s %d %v hcsr04", globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.HeightSensor)
	if moduleShouldBeHere(globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.HeightSensor, "hcsr04") {
		log.Info("HeightSensor should be connected to this device, starting HSCR04")
		go distancesensor.RunDistanceWatcher(onceOnly)
	} else {
		log.Warnf("No hcsr04 Configured - skipping distance monitoring")
	}
	if (globals.MyDevice.Camera.PiCamera == true) {
		go pictureTaker(onceOnly)
	}
}

func pictureTaker(onceOnly bool) {
	for {
		camera.TakeAPicture()
		time.Sleep(30 * time.Second)
		if (onceOnly) {
			break;
		}
	}

}

func main() {
	testableSubmain(false)
}

func testableSubmain(isUnitTest bool) {
	fmt.Printf(globals.ContainerName)
	log.Infof(globals.ContainerName)

	initGlobals()

	// Set up a connection to the server.
	if !isUnitTest {
		log.Infof("Dialing GRPC server at %s", globals.ForwardingAddress)
		conn, err := grpc.Dial(globals.ForwardingAddress, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		globals.Client = pb.NewSensorStoreAndForwardClient(conn)
	}

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	setupGPIO()

	setupPhMonitor()
	numGoRoutines := countGoRoutines()
	var wg sync.WaitGroup
	wg.Add(numGoRoutines)
	startGoRoutines(isUnitTest)

	if len(globals.DevicesFailed) > 0 {
		log.Errorf("Exiting because of device failure %v", globals.DevicesFailed)
		os.Exit(1)
	}

	go makeControlDecisions(isUnitTest)

	go func() {
		if isUnitTest {
			return
		}
		err := listenForCommands(isUnitTest)
		if err != nil {
			log.Errorf("listenForCommands %+v", err)
		}
	}()

	if !isUnitTest {
		log.Infof("all go routines started, waiting for waitgroup to finish")
		wg.Wait()
		log.Infof("exiting main - because waitgroup finished")
	}
}

func isRelayAttached(deviceid int64) (relayIsAttached bool) {
	if len(globals.MyDevice.ACOutlets) > 0 {
		return true
	}
	return false
}

func moduleShouldBeHere(containerName string, mydeviceid int64, deviceInStation bool, moduleType string) (shouldBePresent bool) {
	if !deviceInStation {
		return false
	}
	for i := 0; i < len(globals.MyStation.EdgeDevices); i++ {
		//		log.Infof("%v", globals.MySite.AttachedDevices[i])
		for j := 0; j < len(globals.MyStation.EdgeDevices[i].DeviceModules); j++ {
			if globals.MyStation.EdgeDevices[i].DeviceModules[j].ContainerName == containerName && globals.MyStation.EdgeDevices[i].DeviceID == mydeviceid && globals.MyStation.EdgeDevices[i].DeviceModules[j].ModuleType == moduleType {
				log.Infof("Device %s should be present at %s", globals.MyStation.EdgeDevices[i].DeviceType, globals.MyStation.EdgeDevices[i].DeviceModules[j].Address)
				return true
			}
		}
	}
	return false
}
