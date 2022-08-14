package co2vocmeter

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
 */

/**
Props to https://github.com/JohannesEH/climate-monitor
*/

/* Error register
Bit ERROR_CODE Description
0 WRITE_REG_INVALID The CCS811 received an I²C write request addressed to this station but with
invalid register address ID
1 READ_REG_INVALID The CCS811 received an I²C read request to a mailbox ID that is invalid
2 MEASMODE_INVALID The CCS811 received an I²C request to write an unsupported mode to
MEAS_MODE
3 MAX_RESISTANCE The sensor resistance measurement has reached or exceeded the maximum
range
4 HEATER_FAULT The Heater current in the CCS811 is not in range
5 HEATER_SUPPLY The Heater voltage is not being applied correctly
6 - Reserved for Future Use
7 - Reserved for Future Use
*/
/* Status Register
7 6 5 4 3 2 1 0
FW_MODE APP_ERASE APP_VERIFY APP_VALID DATA_READY - ERROR

Bit(s) Field Description
7 FW_MODE
0: Firmware is in boot mode, this allows new firmware to be loaded
1: Firmware is in application mode. CCS811 is ready to take ADC measurements

6 APP_ERASE
Boot Mode only.
0: No erase completed
1: Application erase operation completed successfully (flag is cleared by APP_DATA
and also by APP_START, SW_RESET, nRESET and APP_VERIFY)
After issuing the ERASE command the application software must wait 500ms
before issuing any transactions to the CCS811 over the I2C interface.

5 APP_VERIFY
Boot Mode only.
0: No verify completed
1: Application verify operation completed successfully (flag is cleared by
APP_START, SW_RESET and nRESET)
After issuing a VERIFY command the application software must wait 70ms before
issuing any transactions to CCS811 over the I²C interface

4 APP_VALID
0: No application firmware loaded
1: Valid application firmware loaded

3 DATA_READY
0: No new data samples are ready
1: A new data sample isready in ALG_RESULT_DATA,this bitis cleared when
ALG_RESULT_DATA is read on the I²C interface

2:1 - Reserved

0 ERROR
Thisbitis clearedby readingERROR_ID(itisnotsufficienttoreadthe ERRORfieldof
ALG_RESULT_DATA and STATUS )
0: No error has occurred
1: There is an error on the I²C or sensor, the ERROR_ID register (0xE0) contains the
error source

*/

import (
	pb "bubblesnet/edge-device/sense-go/bubblesgrpc"
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/messaging"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"golang.org/x/net/context"
	"os"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/experimental/devices/ccs811"
	"periph.io/x/periph/host"
	"time"
)

const (
	ERROR_MASK_WRITE_REG_INVALID byte = 0x01
	ERROR_MASK_READ_REG_INVALID  byte = 0x02
	ERROR_MASK_MEASMODE_INVALID  byte = 0x04
	ERROR_MASK_MAX_RESISTANCE    byte = 0x08
	ERROR_MASK_HEATER_FAULT      byte = 0x10
	ERROR_MASK_HEATER_SUPPLY     byte = 0x20
	ERROR_MASK_RESERVED1         byte = 0x40
	ERROR_MASK_RESERVED2         byte = 0x80

	STATUS_MASK_FIRMWAREMODE byte = 0x80
	STATUS_MASK_APPERASE     byte = 0x40
	STATUS_MASK_APPVERIFY    byte = 0x20
	STATUS_MASK_APPVALID     byte = 0x10
	STATUS_MASK_DATAREADY    byte = 0x08
	STATUS_MASK_RESERVED1    byte = 0x04
	STATUS_MASK_RESERVED2    byte = 0x02
	STATUS_MASK_ERROR        byte = 0x01

	baselineFile = "BASELINE"
)

var baseline = []byte{253, 184}

func ReadCO2VOC() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Errorf("ccs811: host.Init failed %+v", err)
		return
	}
	log.Info("ccs811: Inited")

	opts := ccs811.Opts{
		Addr:               0x5a,
		MeasurementMode:    ccs811.MeasurementModeConstant250,
		InterruptWhenReady: false,
		UseThreshold:       false,
	}

	log.Info("ccs811: I2C: Open connection")
	bus, err := i2creg.Open("")
	if err != nil {
		log.Errorf("ccs811: i2creg.Open failed %+v", err)
		return
	}
	log.Info("ccs811: i2creg.Open succeeded")

	for {
		ccs, err := ccs811.New(bus, &opts)
		if err != nil {
			log.Errorf("ccs811: Couldn't get new ccs811 bus %#v, opts %#v, err %#v", bus, opts, err)
			return
		}
		log.Info("ccs811: ccs811.New succeeded")

		mode, err := ccs.GetMeasurementModeRegister()
		if err != nil {
			log.Errorf("ccs811: Couldn't get measurement mode register err %#v", err)
			return
		}
		log.Info("ccs811: ccs.GetMeasurementModeRegister succeeded")

		fwData, err := ccs.GetFirmwareData()
		if err != nil {
			log.Errorf("ccs811: Couldn't get firmware data err %#v", err)
			return
		}
		log.Info("ccs811: ccs.GetFirmwareData succeeded")

		log.Infof("ccs811: ========================================================================")
		log.Infof("ccs811: Device Information:")
		log.Infof("ccs811: ========================================================================")
		log.Infof("ccs811: HW Model:     %s", ccs.String())
		log.Infof("ccs811: HW Identifier: 0x%X", fwData.HWIdentifier)
		log.Infof("ccs811: HW Version:    0x%X", fwData.HWVersion)
		log.Infof("ccs811: Boot Version: %s", fwData.BootVersion)
		log.Infof("ccs811: App Version:  %s", fwData.ApplicationVersion)
		log.Infof("ccs811: Mode:          ")
		switch mode.MeasurementMode {
		case ccs811.MeasurementModeIdle:
			log.Infof("ccs811: Idle, low power mode")
			break
		case ccs811.MeasurementModeConstant1000:
			log.Infof("ccs811: Constant power mode, IAQ measurement every second")
			break
		case ccs811.MeasurementModePulse:
			log.Infof("ccs811: Pulse heating mode IAQ measurement every 10 seconds")
			break
		case ccs811.MeasurementModeLowPower:
			log.Infof("ccs811: Low power pulse heating mode IAQ measurement every 60 seconds")
			break
		case ccs811.MeasurementModeConstant250:
			log.Infof("ccs811: Constant power mode, sensor measurement every 250ms")
			break
		default:
			log.Infof("ccs811: Unknown")
			break
		}

		count := 0
		lowestBaseLine := loadBaseline()
		lowestBaseLineConverted := binary.LittleEndian.Uint16(lowestBaseLine)

		err = ccs.SetBaseline(lowestBaseLine)
		checkErr(err)

		var val = ccs811.SensorValues{}
		err = ccs.Sense(&val)
		checkErr(err)

		err = ccs.SetBaseline(lowestBaseLine)
		checkErr(err)

		for {
			status, err := ccs.ReadStatus()
			if err != nil {
				log.Errorf("ccs811: status failed %#v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			reportStatus(status)

			if status&STATUS_MASK_DATAREADY == STATUS_MASK_DATAREADY {
				var sensorValues = ccs811.SensorValues{}
				err = ccs.Sense(&sensorValues)
				if err != nil {
					log.Errorf("ccs811: Sense failed %+v", err)
					time.Sleep(5 * time.Second)
					continue
				}

				//			if sensorValues.Error != nil {
				//				log.Errorf("ccs811: sensorvalues error %#v", sensorValues.Error)
				//				time.Sleep(5 * time.Second)
				//				continue
				//			}

				baseline, err := ccs.GetBaseline()
				checkErr(err)

				baselineConverted := binary.LittleEndian.Uint16(baseline)

				if baselineConverted < lowestBaseLineConverted {
					lowestBaseLine = baseline
					lowestBaseLineConverted = baselineConverted
					saveBaseline(baseline)
				}

				log.Infof("ccs811: ECO2: %d ppm", sensorValues.ECO2)
				log.Infof("ccs811: VOC: %d ppb", sensorValues.VOC)
				log.Infof("ccs811: Current: %d", sensorValues.RawDataCurrent)
				log.Infof("ccs811: Voltage: %d", sensorValues.RawDataVoltage)
				//			fmt.Println("Baseline: ", baseline, baselineConverted)

				typeId := "sensor"

				co2m, vocm, curm, voltm := getCCCS811SensorMessages(sensorValues)
				log.Infof("ccs811: co2m = %#v", co2m)
				bytearray, err := json.Marshal(co2m)
				if err != nil {
					log.Errorf("ccs811: marshal co2m error %#v", err)
				}
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: typeId, Data: string(bytearray)}
				sensorReply, err := globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("ccs811: ReadCO2VOC ERROR %#v", err)
				} else {
					log.Debugf("ccs811: co2m message reply %#v", sensorReply)
				}

				log.Infof("ccs811: vocm = %#v", vocm)
				bytearray, err = json.Marshal(vocm)
				if err != nil {
					log.Errorf("ccs811: marshal vocm error %#v", err)
				}
				message = pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: typeId, Data: string(bytearray)}
				sensorReply, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("ccs811: ReadCO2VOC ERROR %#v", err)
				} else {
					log.Debugf("ccs811: vocm message reply %#v", sensorReply)
				}

				log.Infof("ccs811: curm = %#v", curm)
				bytearray, err = json.Marshal(curm)
				if err != nil {
					log.Errorf("ccs811: marshal curm error %#v", err)
				}
				message = pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: typeId, Data: string(bytearray)}
				sensorReply, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("ccs811: ReadCO2VOC ERROR %#v", err)
				} else {
					log.Debugf("ccs811: curm message reply %#v", sensorReply)
				}

				log.Infof("ccs811: voltm = %#v", voltm)
				bytearray, err = json.Marshal(voltm)
				if err != nil {
					log.Errorf("ccs811: marshal co2m error %#v", err)
				}
				message = pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: typeId, Data: string(bytearray)}
				sensorReply, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("ccs811: ReadCO2VOC ERROR %#v", err)
				} else {
					log.Debugf("ccs811: voltm message reply %#v", sensorReply)
				}

			} else {
				log.Error("ccs811: bad status no data ready")
				if status&STATUS_MASK_FIRMWAREMODE == STATUS_MASK_FIRMWAREMODE {
					log.Error("ccs811: still in firmware mode - new() failed silently, keep trying!!!!!")
					time.Sleep(5 * time.Second)
					break
				}

				time.Sleep(5 * time.Second)
				continue
			}
			count = count + 1

			if count%300 == 0 {
				fmt.Println("setting baseline", lowestBaseLine, lowestBaseLineConverted)
				err = ccs.SetBaseline(lowestBaseLine)
				checkErr(err)

				count = 0
			}

			time.Sleep(30 * time.Second)
		}
	}
}

func reportError(err byte) {
	if err&ERROR_MASK_WRITE_REG_INVALID == ERROR_MASK_WRITE_REG_INVALID {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_READ_REG_INVALID == ERROR_MASK_READ_REG_INVALID {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_MEASMODE_INVALID == ERROR_MASK_MEASMODE_INVALID {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_MAX_RESISTANCE == ERROR_MASK_MAX_RESISTANCE {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_HEATER_FAULT == ERROR_MASK_HEATER_FAULT {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_HEATER_SUPPLY == ERROR_MASK_HEATER_SUPPLY {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_RESERVED1 == ERROR_MASK_RESERVED1 {
		log.Infof("ccs811: ")
	}
	if err&ERROR_MASK_RESERVED2 == ERROR_MASK_RESERVED2 {
		log.Infof("ccs811: ")
	}
}
func reportStatus(status byte) {
	statstring := "ccs811: status - "

	if status&STATUS_MASK_FIRMWAREMODE == STATUS_MASK_FIRMWAREMODE {
		//		statstring = statstring + " firmware in application mode "
	} else {
		statstring = statstring + " firmware in boot mode! "
	}

	if status&STATUS_MASK_APPERASE == STATUS_MASK_APPERASE {
		statstring = statstring + " apperase erase successful - boot mode only "
	}

	if status&STATUS_MASK_APPVERIFY == STATUS_MASK_APPVERIFY {
		statstring = statstring + " verify successful - boot mode only "
	}

	if status&STATUS_MASK_APPVALID == STATUS_MASK_APPVALID {
		//		statstring = statstring + " Valid application firmware loaded "
	} else {
		statstring = statstring + " No application firmware loaded "
	}

	if status&STATUS_MASK_DATAREADY == STATUS_MASK_DATAREADY {
		statstring = statstring + " Data ready! "
	} else {
		statstring = statstring + " No data ready "
	}

	if status&STATUS_MASK_RESERVED1 == STATUS_MASK_RESERVED1 {
		statstring = statstring + " reserved1 "
	}
	if status&STATUS_MASK_RESERVED2 == STATUS_MASK_RESERVED2 {
		statstring = statstring + " reserved2 "
	}

	if status&STATUS_MASK_ERROR == STATUS_MASK_ERROR {
		statstring = statstring + " ERROR "
	}
	log.Info(statstring)
}

func getCCCS811SensorMessages(sensorValues ccs811.SensorValues) (co2msg *messaging.CO2SensorMessage, vocmsg *messaging.VOCSensorMessage,
	rawcurrentmsg *messaging.CCS811CurrentMessage, rawvoltagemsg *messaging.CCS811VoltageMessage) {

	co2msg = messaging.NewCO2SensorMessage("co2_sensor", "co2", float64(sensorValues.ECO2), "ppm", "")
	vocmsg = messaging.NewVOCSensorMessage("voc_sensor", "voc", float64(sensorValues.VOC), "ppb", "")
	rawcurrentmsg = messaging.NewCCS811CurrentMessage("ccs811_current_sensor", "ccs811_rawcurrent",
		float64(sensorValues.RawDataCurrent), "ua", "")
	rawvoltagemsg = messaging.NewCCS811VoltageMessage("ccs811_voltage_sensor", "ccs811_rawvoltage",
		float64(sensorValues.RawDataVoltage), "uv", "")

	return co2msg, vocmsg, rawcurrentmsg, rawvoltagemsg
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadBaseline() []byte {
	_, err := os.Stat(baselineFile)

	if os.IsNotExist(err) {
		return []byte{0xFF, 0xFF}
	}

	checkErr(err)

	file, err := os.Open(baselineFile)
	checkErr(err)
	defer file.Close()

	stats, err := file.Stat()
	checkErr(err)

	size := stats.Size()
	bytes := make([]byte, size)

	rdr := bufio.NewReader(file)
	_, err = rdr.Read(bytes)

	return bytes

	return baseline
}

func saveBaseline(baseline []byte) {
	file, err := os.Create(baselineFile)
	checkErr(err)
	defer file.Close()

	wrt := bufio.NewWriter(file)
	_, err = wrt.Write(baseline)
	checkErr(err)

	err = wrt.Flush()
	checkErr(err)
}
