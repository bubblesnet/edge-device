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

const WaterLevelChannelIndex = 1

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

const Etape_slope = 11.37795
const Etape_y_intercept = -17.28562

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
					if channelIndex == 1 {
						adcMessage.ChannelValues[channelIndex].SensorName = "water_level"
						adcMessage.ChannelValues[channelIndex].MeasurementName = "water_level"
						adcMessage.ChannelValues[channelIndex].MeasurementUnits = "gallons"
						adcMessage.ChannelValues[channelIndex].Slope = Etape_slope
						adcMessage.ChannelValues[channelIndex].Yintercept = Etape_y_intercept
					}
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
	if channelIndex == WaterLevelChannelIndex {
		direction := "up"
		err, message := getTranslatedADCSensorMessageForChannel(adcMessage, moduleIndex, channelIndex, direction, adcMessage.ChannelValues[channelIndex].SensorName, adcMessage.ChannelValues[channelIndex].MeasurementName,
			adcMessage.ChannelValues[channelIndex].Slope, adcMessage.ChannelValues[channelIndex].Yintercept, adcMessage.ChannelValues[channelIndex].MeasurementUnits)
		//		log.Infof("sendTranslatedADCSensorMessages adc %#v", adcMessage)
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
	//	log.Infof("sensorName %s", sensorName)

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

func getTranslatedADCSensorMessageForChannel(adcMessage *ADCMessage, moduleIndex int, channelIndex int, direction string,
	sensorName string, measurementName string, slope float64, yintercept float64, measurementUnits string) (err error, message pb.SensorRequest) {

	log.Infof("sendTranslatedADCSensorMessages sensorName %s %v", sensorName, adcMessage.ChannelValues[channelIndex].Voltage)
	typeId := "sensor"

	if moduleIndex == 0 {
		measurementValue := 0.0
		if channelIndex == WaterLevelChannelIndex {
			inches := 0.0
			if adcMessage.ChannelValues[channelIndex].Voltage < MinVoltage {
				inches = 0.0
			} else {
				inches = etapeInchesFromVolts(adcMessage.ChannelValues[channelIndex].Voltage, slope, yintercept)
			}
			measurementValue = etapeInchesToGallons(12.5, 18.0, inches)
			log.Infof("sendTranslatedADCSensorMessages raw %f Volts, %s %f inches, %f %s", adcMessage.ChannelValues[channelIndex].Voltage, measurementName, inches, measurementValue, measurementUnits)

			ads := messaging.NewGenericSensorMessage(sensorName, measurementName,
				measurementValue, measurementUnits, direction)
			bytearray, err := json.Marshal(ads)
			if err != nil {
				log.Errorf("loopforever error %#v", err)
				return err, message
			}
			message = pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: typeId, Data: string(bytearray)}
			//			log.Infof("sendTranslatedADCSensorMessages message = %#v", message)
		}
	}
	return nil, message
}
