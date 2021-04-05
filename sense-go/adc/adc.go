// +build linux,arm

package adc

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"golang.org/x/net/context"
	//	"google.golang.org/grpc"
	"time"
)


/*
type ADCSensorMessage struct {
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	MessageType string `json:"message_type"`
	ModuleName string `json:"sensor_name"`
	ChannelNumber int `json:"channel_number,omitempty"`
	Voltage float64 `json:"value,omitempty"`
	Units string	`json:"units"`
	Gain    int	`json:"gain,omitempty"`
	Rate    int	`json:"rate,omitempty"`
}
*/

func ReadAllChannels(index int, adcMessage *ADCMessage) (err error ) {
	return readAllChannels(ads1115s[index],daps[index], adcMessage)
}

func readAllChannels(ads1115 *i2c.ADS1x15Driver, config AdapterConfig, adcMessage *ADCMessage) ( err error ) {
	log.Debugf("readAllChannels on address 0x%x", config.address)
	var err1 error
	adcMessage.Address = config.address
	adcMessage.BusId = config.bus_id
	for channel := 0; channel < 4; channel++ {
		value, err := ads1115.Read(channel, config.channelConfig[channel].gain, config.channelConfig[channel].rate)
		if err == nil {
			//			log.Debugf("Read  value %.2fV channel %d, gain %d, rate %d\n", value, channel, config.channelConfig[channel].gain, config.channelConfig[channel].rate))
			(*adcMessage).ChannelValues[channel].ChannelNumber = channel
			(*adcMessage).ChannelValues[channel].Voltage = value
			(*adcMessage).ChannelValues[channel].Gain = config.channelConfig[channel].gain
			(*adcMessage).ChannelValues[channel].Rate = config.channelConfig[channel].rate
		} else {
			log.Errorf("readAllChannels Read failed %v", err )
			globals.ReportDeviceFailed("ads1115")
			err1 = err
			break
		}
	}
	return err1
}
var last0 = []float64 {
	0.0,0.0,0.0,0.0,
}
var last1 = []float64 {
	0.0,0.0,0.0,0.0,
}

var ads1115s [2]*i2c.ADS1x15Driver

func RunADCPoller(onceOnly bool) (err error) {

	adcAdaptor := raspi.NewAdaptor() // optional bus/address

	ads1115s[0] = i2c.NewADS1115Driver(adcAdaptor,
		i2c.WithBus(a0.bus_id),
		i2c.WithAddress(a0.address))
	err = ads1115s[0].Start()
	if err != nil {
		log.Errorf("error starting interface %v", err )
		return err
	}

	ads1115s[1] = i2c.NewADS1115Driver(adcAdaptor,
		i2c.WithBus(a1.bus_id),
		i2c.WithAddress(a1.address))
	err = ads1115s[1].Start()
	if err != nil {
		log.Errorf("error starting interface %v", err )
		return err
	}

	for {
		adcMessage := new(ADCMessage)
//		err := readAllChannels(ads1115s[0], a0, adcMessage)
		err := ReadAllChannels(0, adcMessage)
		if err != nil {
			log.Errorf("loopforever error %v", err)
			break
		} else {
			for i := 0; i < len(adcMessage.ChannelValues); i++ {
				direction := ""
				if adcMessage.ChannelValues[i].Voltage > last0[i] {
					direction = "up"
				} else if adcMessage.ChannelValues[i].Voltage < last0[i] {
					direction = "down"
				}
				last0[i] = adcMessage.ChannelValues[i].Voltage

				sensor_name := fmt.Sprintf("adc_%d_%d_%d_%d", adcMessage.BusId, adcMessage.ChannelValues[i].ChannelNumber, adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				ads := messaging.NewADCSensorMessage(sensor_name, sensor_name,
					adcMessage.ChannelValues[i].Voltage,"Volts",
					direction,
					adcMessage.ChannelValues[i].ChannelNumber,adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				bytearray, err := json.Marshal(ads)
				if err != nil {
					log.Errorf("loopforever error %v", err)
					break
				}
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
				_, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("RunADCPoller ERROR %v", err)
				} else {
//					log.Infof("sensor_reply %v", sensor_reply)
				}
			}
		}

		adcMessage = new(ADCMessage)
//		err = readAllChannels(ads1115s[1], a1, adcMessage)
		err = ReadAllChannels(1, adcMessage)
		if err != nil {
			log.Errorf("loopforever error %v", err)
			break
		} else {
			//			bytearray, err := json.Marshal(adcMessage)
			for i := 0; i < len(adcMessage.ChannelValues); i++ {
				direction := ""
				if adcMessage.ChannelValues[i].Voltage > last0[i] {
					direction = "up"
				} else if adcMessage.ChannelValues[i].Voltage < last0[i] {
					direction = "down"
				}
				last0[i] = adcMessage.ChannelValues[i].Voltage
				sensor_name := fmt.Sprintf("adc_%d_%d_%d_%d", adcMessage.BusId, adcMessage.ChannelValues[i].ChannelNumber, adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				ads := messaging.NewADCSensorMessage(sensor_name, sensor_name,
					adcMessage.ChannelValues[i].Voltage,"Volts", direction,
					adcMessage.ChannelValues[i].ChannelNumber,adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				bytearray, err := json.Marshal(ads)
				if err != nil {
					log.Errorf("loopforever error %v", err)
					break
				}
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
				_, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("RunADCPoller ERROR %v", err)
				} else {
//					log.Infof("sensor_reply %v", sensor_reply)
				}
			}
		}
		if onceOnly {
			return nil
		}
		//		readAllChannels(ads1115s[1],a1)
		time.Sleep(time.Duration(globals.MyDevice.TimeBetweenSensorPollingInSeconds) * time.Second)
	}
	log.Errorf("loopforever returning err = %v", err)
	return nil
}


