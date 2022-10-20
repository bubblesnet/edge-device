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

package main

// copyright and license inspection - no issues 4/13/22

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"github.com/go-playground/log"
	"os"
	"time"
)

func isPowerstripAttached(deviceid int64) (relayIsAttached bool) {
	if globals.MyDeviceID == deviceid && len(globals.MyDevice.ACOutlets) > 0 {
		return true
	}
	return false
}

func setEnvironmentalControlString(LocalCurrentState *globals.LocalState) (environmentalControlString string) {
	EnvironmentalControl := ""
	if (*LocalCurrentState).Heater == true {
		EnvironmentalControl += "HEATING "
	} else {
		(*LocalCurrentState).EnvironmentalControl += "COOLING "
	}
	if (*LocalCurrentState).Humidifier == true {
		EnvironmentalControl += "HUMIDIFYING "
	} else {
		EnvironmentalControl += "DRYING "
	}
	return EnvironmentalControl
}

func ControlOxygenation(force bool, DeviceID int64, MyDevice *globals.EdgeDevice, CurrentStage string, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {
	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlOxygenation - no outlets attached")
		return somethingChanged
	}
	switch CurrentStage {
	case globals.IDLE:
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.AIRPUMP, false); somethingChanged == true {
			LogSwitchStateChanged("ControlOxygenation", globals.AIRPUMP, true, false)
		}
		break
	case globals.GERMINATION:
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.AIRPUMP, false); somethingChanged == true {
			LogSwitchStateChanged("ControlOxygenation", globals.AIRPUMP, true, false)
		}
		break
	default:
		if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.AIRPUMP, force); somethingChanged == true {
			LogSwitchStateChanged("ControlOxygenation", globals.AIRPUMP, false, true)
		}
		break
	}
	return somethingChanged
}

func ControlRootWater(force bool, DeviceID int64, MyDevice *globals.EdgeDevice, CurrentStage string, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {
	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlRootWater - no outlets attached")
		return somethingChanged
	}

	switch CurrentStage {
	case globals.IDLE:
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.WATERPUMP, false); somethingChanged == true {
			LogSwitchStateChanged("ControlRootWater", globals.WATERPUMP, true, false)
		}
		break
	case globals.GERMINATION:
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.WATERPUMP, false); somethingChanged == true {
			LogSwitchStateChanged("ControlRootWater", globals.WATERPUMP, true, false)
		}
		break
	default:
		if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.WATERPUMP, force); somethingChanged == true {
			LogSwitchStateChanged("ControlRootWater", globals.WATERPUMP, false, true)
		}
		break
	}
	return somethingChanged
}

func ControlAirflow(force bool, DeviceID int64, MyDevice *globals.EdgeDevice, CurrentStage string, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {
	log.Infof("automation: ControlAirflow force=%v", force)
	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Infof("automation: ControlAirflow - no outlets attached")
		return
	}

	TurnFansOn := false
	TurnFansOff := false
	switch CurrentStage {
	case globals.GERMINATION:
		log.Infof("automation: ControlAirflow - turning fans OFF because stage is GERMINATION")
		TurnFansOff = true
		break
	case globals.IDLE:
		log.Infof("automation: ControlAirflow - turning fans OFF because stage is IDLE")
		TurnFansOff = true
		break
	default:
		log.Infof("automation: ControlAirflow - turning fans ON because stage is not IDLE or GERMINATION")
		TurnFansOn = true
		break
	}

	if Powerstrip.IsOutletOn(MyDevice, globals.HEATER) && os.Getenv("NO_FAN_WITH_HEATER") == "true" {
		log.Infof("automation: ControlAirflow - turning fans OFF because heater is ON and NO_FAN_WITH_HEATER == ", os.Getenv("NO_FAN_WITH_HEATER"))
		TurnFansOff = true
		TurnFansOn = false
	}

	if TurnFansOn == true {
		log.Infof("automation: ControlAirflow - turning fans ON for some reason")
		if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.OUTLETFAN, force); somethingChanged == true {
			LogSwitchStateChanged("ControlAirflow", globals.OUTLETFAN, false, true)
		}
		if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.INLETFAN, force); somethingChanged == true {
			LogSwitchStateChanged("ControlAirflow", globals.INLETFAN, false, true)
		}
	} else {
		if TurnFansOff == true {
			log.Infof("automation: ControlAirflow - turning fans OFF for some reason")
			if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.OUTLETFAN, false); somethingChanged == true {
				LogSwitchStateChanged("ControlAirflow", globals.OUTLETFAN, true, false)
			}
			if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.INLETFAN, false); somethingChanged == true {
				LogSwitchStateChanged("ControlAirflow", globals.INLETFAN, true, false)
			}
		} else {
			log.Infof("automation: ControlAirflow - NOT turning fans ON or OFF for some reason")
		}
	}

	return somethingChanged
}

func LogSwitchStateChanged(functionName string, switchName string, originalState bool, newState bool) {
	log.Infof("StateChange: switch %s from %v to %v via %s", switchName, originalState, newState, functionName)
}

func ControlLight(force bool, DeviceID int64, MyDevice *globals.EdgeDevice, CurrentStage string,
	MyStation globals.Station, CurrentStageSchedule globals.StageSchedule,
	LocalCurrentState *globals.LocalState, currentTime time.Time, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {

	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlLight - no outlets attached")
		return somethingChanged
	}
	if CurrentStage == globals.IDLE {
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.GROWLIGHTVEG, false); somethingChanged == true {
			LogSwitchStateChanged("ControlLight", globals.GROWLIGHTVEG, true, false)
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.GROWLIGHTBLOOM, false); somethingChanged == true {
			LogSwitchStateChanged("ControlLight", globals.GROWLIGHTBLOOM, true, false)
		}
		return
	}
	localTimeHours := currentTime.Hour()
	offsetHours := 5
	if localTimeHours-offsetHours < 0 {
		localTimeHours = 24 + (localTimeHours - offsetHours)
	} else {
		localTimeHours = localTimeHours - offsetHours
	}
	bloomlight := false

	if CurrentStage == globals.GERMINATION || CurrentStage == globals.SEEDLING ||
		CurrentStage == globals.VEGETATIVE || CurrentStage == globals.BLOOMING {
		// If it's time for grow light veg to be on
		if inRange(CurrentStageSchedule.LightOnStartHour, CurrentStageSchedule.HoursOfLight, localTimeHours) {
			log.Infof("automation: ControlLight turning on %s because local hour %d is within %d hours of %d", globals.GROWLIGHTBLOOM,
				localTimeHours, CurrentStageSchedule.HoursOfLight, CurrentStageSchedule.LightOnStartHour)
			if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.GROWLIGHTBLOOM, force); somethingChanged == true {
				LogSwitchStateChanged("ControlLight", globals.GROWLIGHTBLOOM, false, true)
			}
			bloomlight = true
		} else {
			// If it's time for grow light veg to be off
			if LocalCurrentState.GrowLightBloom == true {
				log.Infof("automation: ControlLight turning off %s because local hour %d is outside %d hours of %d", globals.GROWLIGHTBLOOM,
					localTimeHours, CurrentStageSchedule.HoursOfLight, CurrentStageSchedule.LightOnStartHour)
			}
			if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.GROWLIGHTBLOOM, force); somethingChanged == true {
				LogSwitchStateChanged("ControlLight", globals.GROWLIGHTBLOOM, true, false)
			}
			bloomlight = false
		}
	} else {
	}
	if bloomlight && !LocalCurrentState.GrowLightBloom {
		log.Infof("automation: ControlLight Turned veg light ON")
	} else if !bloomlight && LocalCurrentState.GrowLightBloom {
		log.Infof("automation: ControlLight Turned veg light OFF")
	}
	(*LocalCurrentState).GrowLightBloom = bloomlight
	return somethingChanged
}

func inRange(starthour int, numhours int, currenthours int) bool {
	if starthour+numhours >= 24 { // cross days
		if currenthours >= starthour {
			return true
		} else {
			if currenthours < (starthour + numhours - 24) {
				return true
			} else {
				return false
			}
		}
	} else { // within day
		if currenthours >= starthour && currenthours < (starthour+numhours) {
			return true
		} else {
			return false
		}
	}
}

func ControlWaterTemp(force bool,
	DeviceID int64,
	MyDevice *globals.EdgeDevice,
	StageSchedule globals.StageSchedule,
	CurrentStage string,
	ExternalCurrentState globals.ExternalState,
	LocalCurrentState *globals.LocalState,
	LastWaterTemp *float32,
	Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {

	somethingChanged = false

	if !isPowerstripAttached(DeviceID) {
		log.Infof("automation: ControlWaterTemp - no outlets attached")
		return somethingChanged
	}
	if CurrentStage == globals.IDLE {
		log.Infof("automation: ControlWaterTemp - stage is idle, turning off")
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.WATERHEATER, false); somethingChanged == true {
			LogSwitchStateChanged("ControlWaterTemp", globals.WATERHEATER, true, false)
		} // MAKE SURE HEAT IS OFF
		return somethingChanged
	}

	if ExternalCurrentState.TempWater == globals.TEMPNOTSET {
		log.Infof("automation: ControlWaterTemp TEMPNOTSET ExternalCurrentState.TempWater %.3f - ignoring", ExternalCurrentState.TempWater)
		return somethingChanged
	}
	// Go from 62 to 68
	highLimit := StageSchedule.EnvironmentalTargets.WaterTemperature + 3.0
	lowLimit := StageSchedule.EnvironmentalTargets.WaterTemperature - 3.0
	if ExternalCurrentState.TempWater > highLimit {
		if *LastWaterTemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("automation: ControlWaterTemp turning off %s because WaterTemp (%.3f) just rolled over (%.2f/%.3f/%.2f) on way up", globals.WATERHEATER, ExternalCurrentState.TempWater, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit)
			force = true
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.WATERHEATER, force); somethingChanged == true {
			LogSwitchStateChanged("ControlWaterTemp", globals.WATERHEATER, true, false)
		}

		(*LocalCurrentState).WaterHeater = false
		(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
		log.Infof("automation: ControlWaterTemp water temp %f is already too high %f", ExternalCurrentState.TempWater, highLimit)
	} else { // NOT TOO HOT
		if ExternalCurrentState.TempWater < lowLimit { // TOO COLD
			if *LastWaterTemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlWaterTemp turning on %s because Water Temp (%.3f) just fell below (%.2f/%.1f/%.2f) on way up from %.2f", globals.WATERHEATER, ExternalCurrentState.TempWater, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
				force = true
			}
			if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.WATERHEATER, force); somethingChanged == true {
				LogSwitchStateChanged("ControlWaterTemp", globals.WATERHEATER, false, true)
			} // MAKE SURE HEAT IS ON

			(*LocalCurrentState).WaterHeater = true
		} else { // JUST RIGHT
			if *LastWaterTemp < lowLimit {
				log.Infof("automation: ControlWaterTemp Water Temp (%.3f) just entered sweet spot (%.2f/%.1f/%.2f) on way up from %.2f", ExternalCurrentState.TempWater, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
			} else {
				if *LastWaterTemp > highLimit {
					log.Infof("automation: ControlWaterTemp Water Temp (%.3f) just entered sweet spot (%.2f/%.1f/%.2f) on way down from %.3f", ExternalCurrentState.TempWater, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
				} else {
					log.Infof("automation: ControlWaterTemp Water Temp (%.3f) living in the sweet spot (%.2f/%.1f/%.2f) on way down from %.3f", ExternalCurrentState.TempWater, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
				}
			}
		}
	}
	*LastWaterTemp = ExternalCurrentState.TempWater
	return somethingChanged
}

func ControlHeat(force bool,
	DeviceID int64,
	MyDevice *globals.EdgeDevice,
	CurrentStage string,
	CurrentStageSchedule globals.StageSchedule,
	ExternalCurrentState globals.ExternalState,
	LocalCurrentState *globals.LocalState,
	LastTemp *float32,
	Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {

	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlHeat - no outlets attached, exiting")
		return somethingChanged
	}
	//	log.Infof("automation: ControlHeat - current stage is %s", globals.MyStation.CurrentStage)
	if CurrentStage == globals.IDLE {
		log.Debugf("automation: ControlHeat - stage is idle, turning off and exiting")
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.HEATER, false); somethingChanged == true {
			LogSwitchStateChanged("ControlHeat", globals.HEATER, true, false)
		} // MAKE SURE HEAT IS OFF
		return somethingChanged
	}

	highLimit := CurrentStageSchedule.EnvironmentalTargets.Temperature + 2.0
	lowLimit := CurrentStageSchedule.EnvironmentalTargets.Temperature - 2.0

	//	log.Infof("automation: checking temp %.3f for stage %s with highLimit %.3f, lowLimit %.3f", globals.ExternalCurrentState.TempAirMiddle, globals.MyStation.CurrentStage, highLimit,lowLimit)
	if ExternalCurrentState.TempAirMiddle == globals.TEMPNOTSET {
		log.Debugf("automation: ControlHeat TEMPNOTSET ExternalCurrentState.TempAirMiddle %.3f - ignoring and exiting", ExternalCurrentState.TempAirMiddle)
		return somethingChanged
	}
	if ExternalCurrentState.TempAirMiddle > highLimit { // TOO HOT
		log.Infof("automation: ControlHeat turning off %s because internal temp %.3f is over high limit %.3f on way up", globals.HEATER, ExternalCurrentState.TempAirMiddle, highLimit)
		if globals.LastTemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("automation: ControlHeat turning off %s because internal temp %.3f just exceeded high limit %.3f on way up", globals.HEATER, ExternalCurrentState.TempAirMiddle, highLimit)
			force = true
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.HEATER, force); somethingChanged == true {
			LogSwitchStateChanged("ControlHeat", globals.HEATER, true, false)
		} // MAKE SURE HEAT IS OFF

		(*LocalCurrentState).Heater = false
		(*LocalCurrentState).HeaterPad = false
		(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
	} else { // NOT TOO HOT
		if ExternalCurrentState.TempAirMiddle < lowLimit { // TOO COLD
			//			log.Infof("automation: ControlHeat TOO COLD %.3f < lowLimit %.2f", ExternalCurrentState.TempAirMiddle, lowLimit)
			if *LastTemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlHeat turning on %s because internal temp %.3f just fell below low limit %.3f on way down", globals.HEATER, ExternalCurrentState.TempAirMiddle, lowLimit)
				force = true
			}
			if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.HEATER, false); somethingChanged == true {
				LogSwitchStateChanged("ControlHeat", globals.HEATER, false, true)
			} // MAKE SURE HEAT IS ON

			(*LocalCurrentState).Heater = true
			(*LocalCurrentState).HeaterPad = true
		} else { // JUST RIGHT
			log.Infof("automation: ControlHeat JUST RIGHT %.3f", ExternalCurrentState.TempAirMiddle)
			if *LastTemp < lowLimit {
				log.Infof("automation: ControlHeat Temp just entered sweet spot on way up - %.3f", ExternalCurrentState.TempAirMiddle)
			} else {
				if *LastTemp > highLimit {
					log.Infof("automation: ControlHeat Temp just entered sweet spot on way down - %.3f", ExternalCurrentState.TempAirMiddle)
				} else {
				}
			}
		}
	}

	(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
	*LastTemp = ExternalCurrentState.TempAirMiddle
	return somethingChanged
}

func ControlHumidity(force bool,
	DeviceID int64,
	MyDevice *globals.EdgeDevice,
	StageSchedule globals.StageSchedule,
	CurrentStage string,
	ExternalCurrentState globals.ExternalState,
	LocalCurrentState *globals.LocalState,
	LastHumidity *float32,
	Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {

	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlHumidity - no outlets attached")
		return somethingChanged
	}
	if CurrentStage == globals.IDLE {
		//		log.Debugf("automation: ControlHumidity - stage is idle, turning off")
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.HUMIDIFIER, false); somethingChanged == true {
			LogSwitchStateChanged("ControlHumidity", globals.HUMIDIFIER, true, false)
		} // MAKE SURE HUMIDIFIER IS OFF
		return somethingChanged
	}
	highLimit := StageSchedule.EnvironmentalTargets.Humidity + 5.0
	lowLimit := StageSchedule.EnvironmentalTargets.Humidity - 5.0

	if ExternalCurrentState.HumidityInternal == globals.HUMIDITYNOTSET {
		//		log.Debugf("automation: HUMIDITYNOTSET ExternalCurrentState.HumidityInternal %.3f - ignoring", globals.ExternalCurrentState.HumidityInternal))
		return somethingChanged
	}
	if ExternalCurrentState.HumidityInternal > highLimit { // TOO HUMID
		if *LastHumidity < highLimit { // JUST BECAME TOO HUMID
			log.Infof("automation: ControlHumidity turning off %s because HumidityInternal %.3f just exceeded (%.3f/%.1f/%.3f) on way up from %.3f", globals.HUMIDIFIER, ExternalCurrentState.HumidityInternal, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
			force = true
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(MyDevice, globals.HUMIDIFIER, force); somethingChanged == true {
			LogSwitchStateChanged("ControlHeat", globals.HUMIDIFIER, true, false)
		} // MAKE SURE HUMIDIFIER IS OFF
		(*LocalCurrentState).Humidifier = false
	} else { // NOT TOO HOT
		if ExternalCurrentState.HumidityInternal < lowLimit { // TOO COLD
			if *LastHumidity > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlHumidity turning on %s because HumidityInternal %.3f just fell below low (%.3f/%.1f/%.3f) on way down from %.3f", globals.HUMIDIFIER, ExternalCurrentState.HumidityInternal, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
				force = true
			}
			if somethingChanged = Powerstrip.TurnOnOutletByName(MyDevice, globals.HUMIDIFIER, force); somethingChanged == true {
				LogSwitchStateChanged("ControlHumidity", globals.HUMIDIFIER, false, true)
			} // MAKE SURE HUMIDIFIER IS ON

			(*LocalCurrentState).Humidifier = true
		} else { // JUST RIGHT
			if *LastHumidity < lowLimit {
				log.Infof("automation: ControlHumidity HumidityInternal %.3f just entered sweet spot (%.3f/%.1f/%.3f) on way up from %.3f", ExternalCurrentState.HumidityInternal, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
			} else {
				if *LastHumidity > highLimit {
					log.Infof("automation: ControlHumidity HumidityInternal %.3f just entered sweet spot (%.3f/%.1f/%.3f) on way down from %.3f", ExternalCurrentState.HumidityInternal, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
				} else {
				}
			}
			log.Infof("automation: ControlHumidity HumidityInternal %.3f just right", ExternalCurrentState.HumidityInternal)
		}
	}

	(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
	*LastHumidity = ExternalCurrentState.HumidityInternal
	return somethingChanged
}
