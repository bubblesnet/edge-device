// +build linux,arm

package main

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"math"
	"time"
)

func RunTamperDetector() {
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
				if xmove > globals.MyStation.TamperSpec.Xmove ||  ymove > globals.MyStation.TamperSpec.Ymove ||  zmove > globals.MyStation.TamperSpec.Zmove {
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

