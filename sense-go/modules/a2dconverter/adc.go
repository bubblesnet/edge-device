//go:build (linux && arm) || arm64
// +build linux,arm arm64

package a2dconverter

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

func ReadAllChannels(index int, adcMessage *ADCMessage) (err error) {
	return readAllChannels(index, ads1115s[index], daps[index], adcMessage)
}

func readAllChannels(moduleIndex int, ads1115 *i2c.ADS1x15Driver, config AdapterConfig, adcMessage *ADCMessage) (err error) {

	log.Debugf("readAllChannels on address 0x%x", config.address)
	var err1 error
	adcMessage.Address = config.address
	adcMessage.BusId = config.bus_id
	for channel := 0; channel < 4; channel++ {
		value, err := ads1115.Read(channel, config.channelConfig[channel].gain, config.channelConfig[channel].rate)

		//		value, err := ads1115.Read(channel, config.channelConfig[channel].gain, 860)
		if err == nil {
			log.Debugf("Read value %.2fV moduleIndex %d addr 0x%x channel %d, gain %d, rate %d", value, moduleIndex, config.address, channel, config.channelConfig[channel].gain, config.channelConfig[channel].rate)

			(*adcMessage).ChannelValues[channel].ChannelNumber = channel
			(*adcMessage).ChannelValues[channel].Voltage = value
			(*adcMessage).ChannelValues[channel].Gain = config.channelConfig[channel].gain
			(*adcMessage).ChannelValues[channel].Rate = config.channelConfig[channel].rate
		} else {

			log.Errorf("readAllChannels Read failed on address 0x%x channel %d %#v", config.address, channel, err)

			globals.ReportDeviceFailed("ads1115")
			err1 = err
			break
		}

		time.Sleep(time.Duration(config.channelWaitMillis) * time.Millisecond)
	}
	return err1
}

var lastValues = [][]float64{
	{0.0, 0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0, 0.0},
}

var ads1115s [2]*i2c.ADS1x15Driver
var i2cAddresses = []int{0x48, 0x49, 0x4a, 0x4d}

func RunADCPoller(onceOnly bool, pollingWaitInSeconds int) (err error) {

	adcAdaptor := raspi.NewAdaptor() // optional bus/address

	for moduleIndex := 0; moduleIndex < 2; moduleIndex++ {
		ads1115s[moduleIndex] = i2c.NewADS1115Driver(adcAdaptor,
			i2c.WithBus(a0.bus_id),
			i2c.WithAddress(i2cAddresses[moduleIndex]))
		err = ads1115s[moduleIndex].Start()
		if err != nil {
			log.Errorf("error starting interface %#v", err)
			return err
		}
	}

	for {
		for moduleIndex := 0; moduleIndex < 2; moduleIndex++ {
			adcMessage := new(ADCMessage)
			err := ReadAllChannels(moduleIndex, adcMessage)
			if err != nil {
				log.Errorf("loopforever error %#v", err)
				break
			} else {
				for channelIndex := 0; channelIndex < len(adcMessage.ChannelValues); channelIndex++ {
					log.Infof("adc message for moduleIndex %d channel %d %#v", moduleIndex, channelIndex, adcMessage)
					direction := ""
					if adcMessage.ChannelValues[channelIndex].Voltage > lastValues[moduleIndex][channelIndex] {
						direction = "up"
					} else if adcMessage.ChannelValues[channelIndex].Voltage < lastValues[moduleIndex][channelIndex] {
						direction = "down"
					}
					lastValues[moduleIndex][channelIndex] = adcMessage.ChannelValues[channelIndex].Voltage
					err, message := getADCSensorMessageForChannel(adcMessage, moduleIndex, channelIndex, direction)
					log.Infof("adc %#v", adcMessage)
					sensorReply, err := globals.Client.StoreAndForward(context.Background(), &message)
					if err != nil {
						log.Errorf("RunADCPoller ERROR %#v", err)
					} else {
						log.Infof("adc message reply for moduleIndex %d channel %d %#v", moduleIndex, channelIndex, sensorReply)
					}
					_ = sendTranslatedADCSensorMessages(moduleIndex, channelIndex, adcMessage)
				}
			}
			if onceOnly {
				return nil
			}
			//		readAllChannels(ads1115s[1],a1)
			time.Sleep(time.Duration(pollingWaitInSeconds) * time.Second)
		}
	}
	log.Errorf("loopforever returning err = %#v", err)
	return nil
}

func sendTranslatedADCSensorMessages(moduleIndex int, channelIndex int, adcMessage *ADCMessage) (err error) {

	if moduleIndex != 0 {
		return nil
	}
	if channelIndex == 0 || channelIndex == 1 {
		direction := "up"
		err, message := getTranslatedADCSensorMessageForChannel(adcMessage, moduleIndex, channelIndex, direction)
		log.Infof("sendTranslatedADCSensorMessages adc %#v", adcMessage)
		sensorReply, err := globals.Client.StoreAndForward(context.Background(), &message)
		if err != nil {
			log.Errorf("sendTranslatedADCSensorMessages ERROR %#v", err)
		} else {
			log.Infof("sendTranslatedADCSensorMessages message reply for moduleIndex %d channel %d %#v", moduleIndex, channelIndex, sensorReply)
		}
	}
	return err
}

func getADCSensorMessageForChannel(adcMessage *ADCMessage, moduleIndex int, channelIndex int, direction string) (err error, message pb.SensorRequest) {
	sensorName := fmt.Sprintf("adc_%d_%d_%d_%d", moduleIndex, adcMessage.ChannelValues[channelIndex].ChannelNumber, adcMessage.ChannelValues[channelIndex].Gain, adcMessage.ChannelValues[channelIndex].Rate)
	log.Infof("sensorName %s", sensorName)

	ads := messaging.NewADCSensorMessage(sensorName, sensorName,
		adcMessage.ChannelValues[channelIndex].Voltage, "Volts",
		direction,
		adcMessage.ChannelValues[channelIndex].ChannelNumber, adcMessage.ChannelValues[channelIndex].Gain, adcMessage.ChannelValues[channelIndex].Rate)
	bytearray, err := json.Marshal(ads)
	if err != nil {
		log.Errorf("loopforever error %#v", err)
		return err, message
	}
	message = pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
	return nil, message
}

func getTranslatedADCSensorMessageForChannel(adcMessage *ADCMessage, moduleIndex int, channelIndex int, direction string) (err error, message pb.SensorRequest) {
	sensorName := fmt.Sprintf("adc_%d_%d_%d_%d", moduleIndex, adcMessage.ChannelValues[channelIndex].ChannelNumber, adcMessage.ChannelValues[channelIndex].Gain, adcMessage.ChannelValues[channelIndex].Rate)
	log.Infof("sendTranslatedADCSensorMessages sensorName %s %v", sensorName, adcMessage.ChannelValues[channelIndex].Voltage)
	typeId := "sensor"

	if moduleIndex == 0 {
		measurementName := "na"
		measurementUnits := "na"
		measurementValue := 0.0
		if channelIndex == 0 {
			/// TODO convert values using configuration  data
			sensorName = "water_level"
			measurementName = "water_level"
			//			measurementValue = adcMessage.ChannelValues[channelIndex].Voltage
			//			measurementUnits = "Volts"
			// ohms =ohmsAtMax+(ohmsAtMin-ohmsAtMax)-percent*(ohmsAtMin-ohmsAtMax)
			// volts =voltsMin+(voltsMin*(ohms/(ohms+ohmsAtMin)))
			maxVoltage := 3.325
			minResistance := 400.0
			maxResistance := 2200.0
			//			maxGallons := 10.0
			inches, err := etapeInchesFromVolts(maxVoltage, adcMessage.ChannelValues[channelIndex].Voltage, minResistance, maxResistance, 1.0, 12.78)
			if err != nil {
				log.Errorf("etapeInchesFromVolts returned error %v", err)
				return err, message
			}
			gallons := etapeInchesToGallons(inches)
			measurementValue = gallons
			measurementUnits = "Gallons"
			log.Infof("raw Volts %f, inches %f", adcMessage.ChannelValues[channelIndex].Voltage, inches, gallons)
		} else if channelIndex == 1 {
			sensorName = "temp_water"
			measurementName = "temp_water"
			measurementValue = adcMessage.ChannelValues[channelIndex].Voltage
			measurementUnits = "F"
		}
		if channelIndex == 0 || channelIndex == 1 {
			ads := messaging.NewGenericSensorMessage(sensorName, measurementName,
				measurementValue, measurementUnits, direction)
			bytearray, err := json.Marshal(ads)
			if err != nil {
				log.Errorf("loopforever error %#v", err)
				return err, message
			}
			message = pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: typeId, Data: string(bytearray)}
			log.Infof("sendTranslatedADCSensorMessages message = %#v", message)
		}
	}
	return nil, message
}
