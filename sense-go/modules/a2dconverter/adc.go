//go:build (linux && arm) || arm64
// +build linux,arm arm64

/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package a2dconverter

// copyright and license inspection - no issues 4/13/22

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
	"time"
)

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
			//			break
		}

		//		time.Sleep(time.Duration(config.channelWaitMillis) * time.Millisecond)
		time.Sleep(15 * time.Second)
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

	count := 0
	for {
		for moduleIndex := 0; moduleIndex < 2; moduleIndex++ {
			adcMessage := new(ADCMessage)
			err := ReadAllChannels(moduleIndex, adcMessage)
			if err != nil {
				log.Errorf("loopforever error %#v", err)
				break
			} else {
				for channelIndex := 0; channelIndex < len(adcMessage.ChannelValues); channelIndex++ {
					if channelIndex == 1 && moduleIndex == 0 {
						adcMessage.ChannelValues[channelIndex].SensorName = "water_level_sensor"
						adcMessage.ChannelValues[channelIndex].MeasurementName = "water_level"
						adcMessage.ChannelValues[channelIndex].MeasurementUnits = "gallons"
						adcMessage.ChannelValues[channelIndex].Slope = Etape_slope
						adcMessage.ChannelValues[channelIndex].Yintercept = Etape_y_intercept
					}
					count++
					log.Infof("ADCMessageCounter %d", count)
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
					if channelIndex == 1 && moduleIndex == 0 {
						_ = sendTranslatedADCSensorMessages(moduleIndex, channelIndex, adcMessage)
					}
				}
			}
			if onceOnly {
				return nil
			}
			//		readAllChannels(ads1115s[1],a1)
		}
		time.Sleep(15 * time.Second)
	}
	log.Errorf("loopforever returning err = %#v", err)
	return nil
}

func sendTranslatedADCSensorMessages(moduleIndex int, channelIndex int, adcMessage *ADCMessage) (err error) {

	if moduleIndex != 0 {
		return nil
	}
	if channelIndex == WaterLevelChannelIndex {

		err, message := getTranslatedADCSensorMessageForChannel(adcMessage, moduleIndex, channelIndex, "", adcMessage.ChannelValues[channelIndex].SensorName, adcMessage.ChannelValues[channelIndex].MeasurementName,
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
			measurementValue = etapeInchesToGallons(12.5, 13.5, 1.0, inches)
			log.Infof("sendTranslatedADCSensorMessages raw %f Volts, %s %f inches, %f %s", adcMessage.ChannelValues[channelIndex].Voltage, measurementName, inches, measurementValue, measurementUnits)

			direction := ""
			if measurementValue > float64(globals.LastWaterLevel) {
				direction = "up"
			} else if measurementValue < float64(globals.LastWaterLevel) {
				direction = "down"
			}
			globals.LastWaterLevel = float32(measurementValue)

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
