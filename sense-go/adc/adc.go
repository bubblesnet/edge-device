package adc

import (
	"bubblesnet/edge-device/sense-go/globals"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	//	"google.golang.org/grpc"
	"time"
	"golang.org/x/net/context"
	"bubblesnet/edge-device/sense-go/messaging"
)

type ChannelConfig struct {
	gain int
	rate int
}

type AdapterConfig struct {
	bus_id int
	address int
	channelConfig[4]ChannelConfig
}

type ChannelValue struct {
	ChannelNumber int `json:"channel_number,omitempty"`
	Voltage float64 `json:"voltage,omitempty"`
	Gain    int	`json:"gain,omitempty"`
	Rate    int	`json:"rate,omitempty"`
}

/*
type ADCSensorMessage struct {
	ContainerName string `json:"container_name"`
	ExecutableVersion string `json:"executable_version"`
	MessageType string `json:"message_type"`
	SensorName string `json:"sensor_name"`
	ChannelNumber int `json:"channel_number,omitempty"`
	Voltage float64 `json:"value,omitempty"`
	Units string	`json:"units"`
	Gain    int	`json:"gain,omitempty"`
	Rate    int	`json:"rate,omitempty"`
}
*/

type Channels [4]ChannelValue

var	a0 = AdapterConfig{
	bus_id:  1,
	address: 0x48,
	channelConfig: [4]ChannelConfig{
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
	},
}

var a1 = AdapterConfig{
	bus_id: 1,
	address:    0x49,
	channelConfig: [4]ChannelConfig{
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
		{gain: 1,
			rate: 8},
	},
}

type ADCMessage struct {
	BusId         int      `json:"bus_id"`
	Address       int      `json:"address"`
	ChannelValues Channels `json:"channel_values"`
}

func readAllChannels(ads1115 *i2c.ADS1x15Driver, config AdapterConfig, adcMessage *ADCMessage) ( error ) {
	log.Debug(fmt.Sprintf("readAllChannels on address 0x%x", config.address))
	var err1 error
	adcMessage.Address = config.address
	adcMessage.BusId = config.bus_id
	for channel := 0; channel < 4; channel++ {
		value, err := ads1115.Read(channel, config.channelConfig[channel].gain, config.channelConfig[channel].rate)
		if err == nil {
			//			log.Debug(fmt.Sprintf("Read  value %.2fV channel %d, gain %d, rate %d\n", value, channel, config.channelConfig[channel].gain, config.channelConfig[channel].rate))
			(*adcMessage).ChannelValues[channel].ChannelNumber = channel
			(*adcMessage).ChannelValues[channel].Voltage = value
			(*adcMessage).ChannelValues[channel].Gain = config.channelConfig[channel].gain
			(*adcMessage).ChannelValues[channel].Rate = config.channelConfig[channel].rate
		} else {
			log.Error(fmt.Sprintf( "readAllChannels Read failed %v", err ))
			globals.ReportDeviceFailed("ads1115")
			err1 = err
			break
		}
	}
	return err1
}

func RunADCPoller() (error) {
	var ads1115s [2]*i2c.ADS1x15Driver

	adcAdaptor := raspi.NewAdaptor() // optional bus/address

	ads1115s[0] = i2c.NewADS1115Driver(adcAdaptor,
		i2c.WithBus(a0.bus_id),
		i2c.WithAddress(a0.address))
	err := ads1115s[0].Start()
	if err != nil {
		log.Error(fmt.Sprintf("error starting interface %v", err ))
		return err
	}

	ads1115s[1] = i2c.NewADS1115Driver(adcAdaptor,
		i2c.WithBus(a1.bus_id),
		i2c.WithAddress(a1.address))
	err = ads1115s[1].Start()
	if err != nil {
		log.Error(fmt.Sprintf("error starting interface %v", err ))
		return err
	}

	for {
		adcMessage := new(ADCMessage)
		err := readAllChannels(ads1115s[0], a0, adcMessage)
		if err != nil {
			log.Error(fmt.Sprintf("loopforever error %v", err))
			break
		} else {
			for i := 0; i < len(adcMessage.ChannelValues); i++ {
				sensor_name := fmt.Sprintf("adc_%d_%d_%d_%d", adcMessage.BusId, adcMessage.ChannelValues[i].ChannelNumber, adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				ads := messaging.NewADCSensorMessage(sensor_name,
					adcMessage.ChannelValues[i].Voltage,"Volts",
					adcMessage.ChannelValues[i].ChannelNumber,adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				bytearray, err := json.Marshal(ads)
				if err != nil {
					log.Errorf("loopforever error %v", err)
					break
				}
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
				sensor_reply, err := globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Error(fmt.Sprintf("RunADCPoller ERROR %v", err))
				} else {
					log.Infof("sensor_reply %v", sensor_reply)
				}
			}
		}

		adcMessage = new(ADCMessage)
		err = readAllChannels(ads1115s[1], a1, adcMessage)
		if err != nil {
			log.Error(fmt.Sprintf("loopforever error %v", err))
			break
		} else {
			//			bytearray, err := json.Marshal(adcMessage)
			for i := 0; i < len(adcMessage.ChannelValues); i++ {
				sensor_name := fmt.Sprintf("adc_%d_%d_%d_%d", adcMessage.BusId, adcMessage.ChannelValues[i].ChannelNumber, adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				ads := messaging.NewADCSensorMessage(sensor_name,
					adcMessage.ChannelValues[i].Voltage,"Volts",
					adcMessage.ChannelValues[i].ChannelNumber,adcMessage.ChannelValues[i].Gain, adcMessage.ChannelValues[i].Rate)
				bytearray, err := json.Marshal(ads)
				if err != nil {
					log.Errorf("loopforever error %v", err)
					break
				}
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
				sensor_reply, err := globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Error(fmt.Sprintf("RunADCPoller ERROR %v", err))
				} else {
					log.Infof("sensor_reply %v", sensor_reply)
				}
			}
		}
		//		readAllChannels(ads1115s[1],a1)
		time.Sleep(time.Second * 15)
	}
	log.Error(fmt.Sprintf("loopforever returning err = %v", err))
	return nil
}


