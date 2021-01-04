package main

import (
	"bubblesnet/edge-device/sense-go/adc"
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
		log.Error(fmt.Sprintf("adxl345 robot start error %v", err))
	}
}

func runDistanceWatcher() {
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
			log.Error(fmt.Sprintf("rundistancewatcher error = %v", err ))
			break
		}

		time.Sleep(15 * time.Second)
	}
}

func runLocalStateWatcher() {
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

	powerstrip.InitRpioPins()
	powerstrip.TurnAllOff(1)

	if globals.Config.EZOPH {
		ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
		err = ezoDriver.Start()
		if err != nil {
			log.Error(fmt.Sprintf("ezo start error %v", err))
		}
	}

	numGoroutines := 6
	if !globals.Config.ADXL345 {
		numGoroutines--
	}
	if !globals.Config.ADS1115_1 && !globals.Config.ADS1115_2 {
		numGoroutines--
	}
	if !globals.Config.EZOPH {
		numGoroutines--
	}
	if !globals.Config.HCSR04 {
		numGoroutines--
	}
	if !globals.Config.Relay {
		numGoroutines--
	}
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	if globals.Config.ADXL345 {
		go runTamperDetector()
	} else {
		log.Warn(fmt.Sprint("No adxl345 Configured - skipping tamper detection"))
	}
	if globals.Config.ADS1115_1 || globals.Config.ADS1115_2 {
		go func() {
			err :=  adc.RunADCPoller()
			if err != nil {
				log.Errorf("rpio.close %+v", err)
			}
		}()
	} else {
		log.Warn(fmt.Sprint("No ads1115s configured - skipping A to D conversion"))
	}
	if globals.Config.EZOPH {
		go func() {
			err = readPh()
			if err != nil {
				log.Errorf("readPh %+v", err)
			}
		}()
	} else {
		log.Warn(fmt.Sprint("No ezoph configured - skipping pH monitoring"))
	}
	if globals.Config.HCSR04 {
		go runDistanceWatcher()
	} else {
		log.Warn(fmt.Sprint("No hcsr04 Configured - skipping A to D conversion"))
	}
	if globals.Config.Relay {
//		go runPinToggler()
	} else {
		log.Warn(fmt.Sprint("No relay Ccnfigured - skipping GPIO relay control"))
	}

	go runLocalStateWatcher()
	go makeControlDecisions()

	log.Info(fmt.Sprintf("waiting for waitgroup to finish" ))
	wg.Wait()
	log.Info(fmt.Sprintf("exiting main - because waitgroup finished" ))

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
