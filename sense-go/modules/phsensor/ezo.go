//go:build (linux && arm) || arm64
// +build linux,arm arm64

package phsensor

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"errors"
	"github.com/go-playground/log"
	"gobot.io/x/gobot/platforms/raspi"
	"golang.org/x/net/context"
	"time"
)

func StartEzoDriver() {
	log.Info("Starting Atlas EZO driver")
	ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
	err := ezoDriver.Start()
	if err != nil {
		globals.ReportDeviceFailed("ezoph")
		log.Errorf("ezo start error %#v", err)

	}
}

func StartEzo(once_only bool) {
	log.Info("RootPhSensor should be connected to this device, starting EZO reader")
	go func() {
		if err := ReadPh(once_only); err != nil {
			log.Errorf("ReadPh %+v", err)
		}
	}()
}

var lastPh = float64(0.0)

func ReadPh(once_only bool) error {
	ezoDriver := NewAtlasEZODriver(raspi.NewAdaptor())
	err := ezoDriver.Start()
	if err != nil {
		log.Errorf("ezoDriver.Start returned ph device error %#v", err)

		return err
	}
	var e error = nil

	for {
		ph, err := ezoDriver.Ph()
		if err != nil {
			log.Errorf("ReadPh error %#v", err)

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
			if globals.Client != nil {
				_, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("RunADCPoller ERROR %#v", err)
				} else {
					//				log.Infof("sensor_reply %#v", sensor_reply)

				}
			} else {
				e = errors.New("GRPC client is not connected!")
			}
		}
		if once_only {
			break
		}
		//		x := globals.MyDevice.TimeBetweenSensorPollingInSeconds

		time.Sleep(30 * time.Second)
	}
	log.Debugf("returning %#v from readph", e)

	return e
}
