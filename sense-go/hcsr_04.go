// +build linux,arm

package main

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"github.com/go-playground/log"
	hc "github.com/jdevelop/golang-rpi-extras/sensor_hcsr04"
	"golang.org/x/net/context"
	"time"
)

func RunDistanceWatcher() {
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

