package main

// copyright and license inspection - no issues 4/13/22

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/a2dconverter"
	"bubblesnet/edge-device/sense-go/modules/accelerometer"
	"bubblesnet/edge-device/sense-go/modules/camera"
	"bubblesnet/edge-device/sense-go/modules/distancesensor"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	gonewire "bubblesnet/edge-device/sense-go/modules/onewire"
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
	log.Info("automation: runLocalStateWatcher")
	for true {
		bytearray, err := json.Marshal(globals.LocalCurrentState)
		if err == nil {
			log.Debugf("sending local current state msg %s?", string(bytearray))
			//			err = grpc.SendStoreAndForwardMessageWithRetries(grpc.GetSequenceNumber(), string(bytearray[:]), 3)
			//			if err != nil {
			//				log.Error(fmt.Sprintf("runLocalStateWatcher ERROR %#v", err))
			//			}
		} else {
			//			log.Debugf("runLocalStateWatcher error = %#v", err ))
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

func processCommand(msg *stomp.Message, Powerstrip gpiorelay.PowerstripService) (resub bool, err error) {
	if msg == nil || msg.Err != nil {
		if msg != nil && msg.Err != nil {
			if strings.Contains(fmt.Sprintf("%#v", msg.Err), "timeout") {
				//				log.Debugf("queue read timed out - NOT resubscribing %#v", msg.Err)
				return true, nil
			} else {
				log.Errorf("listenForCommands read topic error %#v - resubscribing", msg.Err)
			}
			time.Sleep(2 * time.Second)
			return true, msg.Err
		} else {
			//				log.Errorf("listenForCommands read topic error %#v", msg)
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
		log.Errorf("listenForCommands marshal error %#v", err)
		return false, err
	}
	log.Infof("listenForCommands parsed body into %#v", header)
	log.Infof("header.Command === %s", header.Command)
	switch header.Command {
	case "stage":
		log.Infof("Changing stage via message %s", msg.Body)
		stageMessage := StageMessage{}
		if err := json.Unmarshal(msg.Body, &stageMessage); err != nil {
			log.Errorf("couldn't parse stage message %s, %#v", msg.Body, err)
			break
		}
		log.Infof("listenForCommands parsed body into StageMessage %#v", stageMessage)
		ChangeStageTo(stageMessage.StageName)
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
		Powerstrip.ReportAll(globals.MyDevice, 200*time.Millisecond)
		Powerstrip.SendSwitchStatusChangeEvent("automaticControl", globals.MyStation.AutomaticControl, globals.GetSequence())
		break
	case "switch":
		{
			if countACOutlets() == 0 {
				log.Infof("No ac outlets configured on this device")
				break
			}
			switchMessage := SwitchMessage{}
			err := json.Unmarshal(msg.Body, &switchMessage)
			log.Infof("listenForCommands parsed body into SwitchMessage %#v on-demand", switchMessage)
			if err != nil {
				log.Errorf("listenForCommands switch error %#v", err)
				break
			}
			if !Powerstrip.IsMySwitch(globals.MyDevice, switchMessage.SwitchName) {
				log.Infof("Not my switch %s", switchMessage.SwitchName)
				break
			}
			if switchMessage.SwitchName == "automaticControl" {
				log.Infof("listenForCommands setting %s to %#v on-demand", switchMessage.SwitchName, switchMessage.On)
				originalState := globals.MyStation.AutomaticControl
				globals.MyStation.AutomaticControl = switchMessage.On
				LogSwitchStateChanged("processCommand", switchMessage.SwitchName, originalState, switchMessage.On)
				log.Infof("automaticControl - sending switch changed event to console")
				gpiorelay.GetPowerstripService().SendSwitchStatusChangeEvent(switchMessage.SwitchName, switchMessage.On, 1120)
				if globals.MyStation.AutomaticControl {
					initializeOutletsForAutomation() // Make sure the switches conform to currently configured automation
				}
			} else if switchMessage.On == true {
				log.Infof("listenForCommands turning on %s on-demand", switchMessage.SwitchName)
				if stateChanged := Powerstrip.TurnOnOutletByName(globals.MyDevice, switchMessage.SwitchName, true); stateChanged == true {
					LogSwitchStateChanged("processCommand", switchMessage.SwitchName, false, true)
				}
			} else {
				log.Infof("listenForCommands turning off %s on-demand", switchMessage.SwitchName)
				if stateChanged := gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.MyDevice, switchMessage.SwitchName, true); stateChanged == true {
					LogSwitchStateChanged("processCommand", switchMessage.SwitchName, true, false)
				}
			}
			break
		}
	default:
		{
			break
		}
	}
	log.Infof("listenForCommands successfully processed message %s", string(msg.Body))
	return false, nil
}

func ChangeStageTo(StageName string) {
	log.Infof("ChangeStageTo %s", StageName)
	globals.MyStation.CurrentStage = StageName
	globals.MySite.Stations[0].CurrentStage = StageName
	globals.CurrentStageSchedule = findSchedule(StageName)
	globals.WriteConfig(globals.PersistentStoreMountPoint, "", "config.json")
	if globals.MyStation.AutomaticControl {
		initializeOutletsForAutomation() // Make sure the switches conform to newly configured automation
	}
}

func findSchedule(StageName string) (stageSchedule globals.StageSchedule) {
	for i := 0; i < len(globals.MyStation.StageSchedules); i++ {
		if globals.MyStation.StageSchedules[i].Name == StageName {
			return globals.MyStation.StageSchedules[i]
		}
	}
	return globals.StageSchedule{}
}
func listenForCommands(isUnitTest bool) (err error) {
	topicName := fmt.Sprintf("/topic/%8.8d/%8.8d", globals.MySite.UserID, globals.MyDevice.DeviceID)
	hostPort := fmt.Sprintf("%s:%d", globals.MySite.ControllerHostName, 61613)
	log.Infof("listenForCommands at %s topic %s", hostPort, topicName)

	var options func(*stomp.Conn) error = func(*stomp.Conn) error {
		stomp.ConnOpt.Login("userid", "userpassword")
		stomp.ConnOpt.Host(globals.MySite.ControllerHostName)
		stomp.ConnOpt.RcvReceiptTimeout(30 * time.Second)
		stomp.ConnOpt.HeartBeat(60*time.Second, 60*time.Second) // I put this but seems no impact
		return nil
	}

	for j := 0; ; j++ {
		log.Debugf("stomp.Dial %s at %d - if this is the last message you see, open the firewall port 61613 on ActiveMQ host", hostPort, getNowMillis())
		stompConn, err := stomp.Dial("tcp", hostPort, options)
		if err != nil {
			log.Errorf("listenForCommands dial error %#v", err)
			return err
		}
		log.Infof("listenForCommands connected to %s", hostPort)
		defer stompConn.Disconnect()

		log.Infof("subscribing to topic %s", topicName)
		sub, err := stompConn.Subscribe(topicName, stomp.AckClient)
		if err != nil {
			log.Errorf("listenForCommands subscribe error at %d %#v no retry!", getNowMillis(), err)
			return err
		}
		//
		for i := 0; ; i++ {
			//		log.Infof("listenForCommands read %d", i)
			msg := <-sub.C
			reSubscribe, err := processCommand(msg, gpiorelay.GetPowerstripService())
			if err != nil {
				log.Warnf("processCommand error %#v need to redial/resubscribe", err)
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
			ps.TurnOnOutletByName(globals.MyDevice, globals.MyDevice.ACOutlets[i].Name, true)
			ReportSwitchInitialized("initializeOutletsFromConfiguration", globals.MyDevice.ACOutlets[i].Name, true)
		} else {
			ps.TurnOffOutletByName(globals.MyDevice, globals.MyDevice.ACOutlets[i].Name, true)
			ReportSwitchInitialized("initializeOutletsFromConfiguration", globals.MyDevice.ACOutlets[i].Name, false)
		}
	}
}

func ReportSwitchInitialized(functionName string, switchName string, newState bool) {
	log.Infof("StateChange: switch %s initialized to %v via %s", switchName, newState, functionName)
}

func initializeOutletsForAutomation() {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: initializeOutletsForAutomation - no outlets attached")
		return
	}
	log.Infof("automation: initializeOutletsForAutomation currentStage %s", globals.MyStation.CurrentStage)

	ControlLight(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.MyStation.CurrentStage,
		*globals.MyStation,
		globals.CurrentStageSchedule,
		&globals.LocalCurrentState,
		time.Now(),
		gpiorelay.GetPowerstripService())
	ControlWaterTemp(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.CurrentStageSchedule,
		globals.MyStation.CurrentStage,
		globals.ExternalCurrentState,
		&globals.LocalCurrentState,
		&globals.LastWaterTemp,
		gpiorelay.GetPowerstripService())
	ControlHeat(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.MyStation.CurrentStage,
		globals.CurrentStageSchedule,
		globals.ExternalCurrentState,
		&globals.LocalCurrentState,
		&globals.LastTemp,
		gpiorelay.GetPowerstripService())
	ControlHumidity(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.CurrentStageSchedule,
		globals.MyStation.CurrentStage,
		globals.ExternalCurrentState,
		&globals.LocalCurrentState,
		&globals.LastHumidity,
		gpiorelay.GetPowerstripService())
	ControlOxygenation(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.MyStation.CurrentStage,
		gpiorelay.GetPowerstripService())
	ControlRootWater(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.MyStation.CurrentStage,
		gpiorelay.GetPowerstripService())
	ControlAirflow(true,
		globals.MyDevice.DeviceID,
		globals.MyDevice,
		globals.MyStation.CurrentStage,
		gpiorelay.GetPowerstripService())
}

func makeControlDecisions(once_only bool) {
	log.Infof("makeControlDecisions endless loop with once_only set to %t", once_only)
	i := 0

	for {
		gsm := pb.GetStateRequest{}
		gsm.Sequence = globals.GetSequence()
		gr, err := globals.Client.GetState(context.Background(), &gsm)
		if err != nil {
			log.Errorf("getState got error %#v", err)
		} else {
			globals.ExternalCurrentState.TempF = gr.TempF
			globals.ExternalCurrentState.WaterTempF = gr.WaterTempF
			globals.ExternalCurrentState.Humidity = gr.Humidity
			//			fmt.Printf("automation: gr = %+v", gr)
			//			fmt.Printf("automation: TempF %f, WaterTempF %f, Humidity %f\n", gr.TempF, gr.WaterTempF, gr.Humidity)
		}
		//		log.Infof("Got state TempF %f Humidity %f", gr.TempF, gr.Humidity)

		if globals.MyStation.AutomaticControl {
			if !isRelayAttached(globals.MyDevice.DeviceID) {
				log.Debugf("automation: makeControlDecisions - no outlets attached ")
				return
			} else {
				ControlLight(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.MyStation.CurrentStage,
					*globals.MyStation,
					globals.CurrentStageSchedule,
					&globals.LocalCurrentState,
					time.Now(),
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
				ControlHeat(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.MyStation.CurrentStage,
					globals.CurrentStageSchedule,
					globals.ExternalCurrentState,
					&globals.LocalCurrentState,
					&globals.LastTemp,
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
				ControlWaterTemp(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.CurrentStageSchedule,
					globals.MyStation.CurrentStage,
					globals.ExternalCurrentState,
					&globals.LocalCurrentState,
					&globals.LastWaterTemp,
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
				ControlHumidity(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.CurrentStageSchedule,
					globals.MyStation.CurrentStage,
					globals.ExternalCurrentState,
					&globals.LocalCurrentState,
					&globals.LastHumidity,
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
				ControlOxygenation(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.MyStation.CurrentStage,
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
				ControlRootWater(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.MyStation.CurrentStage,
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
				ControlAirflow(false,
					globals.MyDevice.DeviceID,
					globals.MyDevice,
					globals.MyStation.CurrentStage,
					gpiorelay.GetPowerstripService())
				time.Sleep(time.Second) // Try not to toggle AC mains power too quickly
			}
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

func initGlobals(testing bool) {
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
		fmt.Printf("error read device %#v\n", err)
		return
	}
	globals.MySite.ControllerHostName, err = globals.ReadMyServerHostname(globals.PersistentStoreMountPoint, "", "hostname")
	if err != nil {
		fmt.Printf("error read serverHostname %#v\n", err)
		return
	}
	// Read the configuration file
	fmt.Printf("Read deviceid %d and server_hostname %s\n", globals.MyDeviceID, globals.MySite.ControllerHostName)
	readConfigFromDisk()
	// Get a NEW config file from server and save to disk
	if err := globals.GetConfigFromServer(globals.PersistentStoreMountPoint, "", "config.json"); err != nil {
		if testing {
			fmt.Printf("Returning because of bad configuration\n")
			return
		}
		fmt.Printf("Exiting because of bad configuration - sleeping for 60 seconds to allow intervention\n")
		time.Sleep(60 * time.Second)
		os.Exit(1)
	}
	// Reread the configuration file
	readConfigFromDisk()

	globals.MySite.LogLevel = "silly,debug,info,warn,fatal,notice,error,alert"
	//	fmt.Printf("done getting config from server %#v\n\n", globals.MySite)
	globals.ConfigureLogging(globals.MySite, "sense-go")

	//	globals.MySite.Station.HeightSensor = true
	reportVersion()

	log.Debug("debug")
	log.Info("automation: info")
	log.Notice("notice")
	log.Warn("warn")
	log.Error("error")
	// log.Panic("panic") // this will panic
	log.Alert("alert")

	//	log.Infof("globals.Configuration = %#v", globals.MySite)
	//	log.Infof("stageSchedule = %#v", globals.CurrentStageSchedule)
}

func readConfigFromDisk() {
	if err := globals.ReadCompleteSiteFromPersistentStore(globals.PersistentStoreMountPoint, "", "config.json", &globals.MySite, &globals.CurrentStageSchedule); err != nil {
		fmt.Printf("ReadCompleteSiteFromPersistentStore failed - using default config\n")
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
		//		fmt.Printf("\ngetconfigfromserver config = %#v\n\n", globals.MySite)
	}
}

func setupGPIO(MyStation *globals.Station, MyDevice *globals.EdgeDevice, Powerstrip gpiorelay.PowerstripService) {
	rpio.OpenRpio()
	/*	defer func() {
		rpio.CloseRpio()
	}() */

	if isRelayAttached(MyDevice.DeviceID) {
		log.Infof("Relay is attached to device %d", MyDevice.DeviceID)
		Powerstrip.InitRpioPins(globals.MyDevice, globals.RunningOnUnsupportedHardware())
		Powerstrip.TurnAllOff(globals.MyDevice, 1) // turn all OFF first since initalizeOutlets doesnt
		if globals.MyStation.AutomaticControl {
			initializeOutletsForAutomation()
		} else {
			initializeOutletsFromConfiguration()
		}
		Powerstrip.SendSwitchStatusChangeEvent("automaticControl", MyStation.AutomaticControl, globals.GetSequence())
	} else {
		log.Infof("There is no relay attached to device %d", MyDevice.DeviceID)
	}
}

func setupPhMonitor() {
	log.Infof("setupPhMonitor")
	globals.ValidateConfigured("setupPhMonitor")
	if moduleShouldBeHere(globals.ContainerName, globals.MyStation, globals.MyDevice.DeviceID, globals.MyStation.RootPhSensor, "ezoph") {
		log.Info("automation: RootPhSensor configured for this device, starting")
		phsensor.StartEzoDriver()
		log.Debug("Ezo driver started")
	} else {
		log.Info("automation: RootPhSensor not configured for this device")
	}
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
	log.Info("startGoRoutines")
	if moduleShouldBeHere(globals.ContainerName, globals.MyStation, globals.MyDevice.DeviceID, globals.MyStation.ThermometerWater, "DS18B20") {
		log.Info("automation: Water Temperature configured for this device, starting")

		go gonewire.ReadOneWire()
	} else {
		log.Warnf("Water Temperature (DS18B20) not configured for this device - skipping water level")
	}

	if moduleShouldBeHere(globals.ContainerName, globals.MyStation, globals.MyDevice.DeviceID, globals.MyStation.MovementSensor, "adxl345") {
		log.Info("automation: MovementSensor configured for this device, starting")

		go accelerometer.GetTamperDetectorService().RunTamperDetector(onceOnly)
	} else {
		log.Warnf("MovementSensor (adxl345) not configured for this device - skipping tamper detection")
	}
	log.Infof("adc %s %d %#v ads1115", globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.WaterLevelSensor)
	if moduleShouldBeHere(globals.ContainerName, globals.MyStation, globals.MyDevice.DeviceID, globals.MyStation.WaterLevelSensor, "ads1115") {
		log.Info("automation: WaterlevelSensor configured for this device, starting ADC")
		go func() {
			err := a2dconverter.RunADCPoller(onceOnly, globals.PollingWaitInSeconds)

			if err != nil {
				log.Errorf("rpio.close %+v", err)
			}
		}()
	} else {
		log.Warnf("WaterLevelSensor (ads1115) not configured for this device - skipping A to D conversion because globals.MyStation.WaterLevelSensor == %#v", globals.MyStation.WaterLevelSensor)
	}
	log.Info("root ph")
	if moduleShouldBeHere(globals.ContainerName, globals.MyStation, globals.MyDevice.DeviceID, globals.MyStation.RootPhSensor, "ezoph") {
		log.Info("automation: RootPhSensor configured for this device, starting ezoPh")

		phsensor.StartEzo(onceOnly)
	} else {
		log.Warnf("RootPhSensor (ezoPh) not configured for this device, - skipping pH monitoring")
	}
	log.Infof("moduleShouldBeHere %s %d %#v hcsr04", globals.ContainerName, globals.MyDevice.DeviceID, globals.MyStation.HeightSensor)
	if moduleShouldBeHere(globals.ContainerName, globals.MyStation, globals.MyDevice.DeviceID, globals.MyStation.HeightSensor, "hcsr04") {
		log.Info("automation: HeightSensor configured for this device, starting HSCR04")

		go distancesensor.RunDistanceWatcher(onceOnly, false)
	} else {
		log.Warnf("HeightSensor (hcsr04) not configured for this device - skipping distance monitoring")
	}

	if globals.MyDevice.Camera.PiCamera == true {
		log.Info("automation: Camera configured for this device, starting picture taker")
		go pictureTaker(onceOnly)
	} else {
		log.Warnf("Camera (piCamers) not configured for this device - skipping picture taker")
	}
}

func pictureTaker(onceOnly bool) {
	for {
		camera.TakeAPicture()
		time.Sleep(30 * time.Second)
		if onceOnly {
			break
		}
	}

}

func main() {
	testableSubmain(false)
}

func testableSubmain(isUnitTest bool) {
	fmt.Printf(globals.ContainerName)
	log.Infof(globals.ContainerName)

	initGlobals(isUnitTest)

	// Set up a connection to the server.
	if !isUnitTest {
		log.Infof("Dialing GRPC server at %s", globals.ForwardingAddress)
		conn, err := grpc.Dial(globals.ForwardingAddress, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %#v", err)
		}
		defer conn.Close()
		globals.Client = pb.NewSensorStoreAndForwardClient(conn)
	}

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	setupGPIO(globals.MyStation, globals.MyDevice, gpiorelay.GetPowerstripService())

	setupPhMonitor()
	numGoRoutines := countGoRoutines()
	var wg sync.WaitGroup
	wg.Add(numGoRoutines)
	startGoRoutines(isUnitTest)

	if len(globals.DevicesFailed) > 0 {
		log.Errorf("Exiting because of device failure %#v", globals.DevicesFailed)
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
	if globals.MyDeviceID == deviceid && len(globals.MyDevice.ACOutlets) > 0 {
		return true
	}
	return false
}

func moduleShouldBeHere(containerName string, MyStation *globals.Station, mydeviceid int64, deviceInStation bool, moduleType string) (shouldBePresent bool) {
	if !deviceInStation {
		return false
	}
	for i := 0; i < len(MyStation.EdgeDevices); i++ {
		//		log.Infof("%#v", globals.MySite.AttachedDevices[i])
		for j := 0; j < len(MyStation.EdgeDevices[i].DeviceModules); j++ {
			if MyStation.EdgeDevices[i].DeviceModules[j].ContainerName == containerName && MyStation.EdgeDevices[i].DeviceID == mydeviceid && MyStation.EdgeDevices[i].DeviceModules[j].ModuleType == moduleType {
				log.Infof("Module %s should be present at %s", globals.MyStation.EdgeDevices[i].DeviceModules[j].ModuleType, MyStation.EdgeDevices[i].DeviceModules[j].Address)
				return true
			}
		}
	}
	return false
}
