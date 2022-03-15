package gonewire

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-playground/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	dir  = flag.String("path", "/sys/bus/w1/devices/", "-path /sys/bus/w1/devices/")
	addr = flag.String("address", ":7777", "-address :7777")
)

func init() {
	flag.Parse()
}

func ReadOneWire() {

	gw, err := New(*dir)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigs:
				cancel()
				return
			case v := <-gw.Values():
				ftemp, _ := strconv.ParseFloat(v.Value, 64)
				ftemp = ftemp / 1000
				fahrenheit := (ftemp * 1.8000) + 32.00

				direction := ""
				if fahrenheit > float64(globals.LastWaterTemp) {
					direction = "up"
				} else if fahrenheit < float64(globals.LastWaterTemp) {
					direction = "down"
				}
				globals.LastWaterTemp = float32(fahrenheit)

				phm := messaging.NewGenericSensorMessage("thermometer_water", "temp_water", fahrenheit, "F", direction)
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
					_ = errors.New("GRPC client is not connected!")
				}
			}
		}
	}()

	gw.OnReadError(func(e error, s *Sensor) {
		fmt.Printf("onReadError\n")
		log.Errorf("blah")
		log.Errorf("[ERR] %s", s.ID())
		log.Errorf("[ERR] %#v", err)
	})

	gw.Start(ctx, 10*time.Second)
}
