package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"github.com/go-playground/log"
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

func ControlOxygenation(force bool, DeviceID int64, CurrentStage string, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {
	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlOxygenation - no outlets attached")
		return somethingChanged
	}
	switch CurrentStage {
	case globals.IDLE:
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.AIRPUMP, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlOxygenation", globals.AIRPUMP, true, false)
		}
		break
	case globals.GERMINATION:
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.AIRPUMP, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlOxygenation", globals.AIRPUMP, true, false)
		}
		break
	default:
		if somethingChanged = Powerstrip.TurnOnOutletByName(globals.AIRPUMP, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlOxygenation", globals.AIRPUMP, false, true)
		}
		break
	}
	return somethingChanged
}

func ControlRootWater(force bool, DeviceID int64, CurrentStage string, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {
	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlRootWater - no outlets attached")
		return somethingChanged
	}

	if CurrentStage == globals.IDLE {
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.WATERPUMP, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlRootWater", globals.WATERPUMP, true, false)
		}
		return somethingChanged
	}
	switch CurrentStage {
	case globals.IDLE:
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.WATERPUMP, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlRootWater", globals.WATERPUMP, true, false)
		}
		break
	case globals.GERMINATION:
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.WATERPUMP, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlRootWater", globals.WATERPUMP, true, false)
		}
		break
	default:
		if somethingChanged = Powerstrip.TurnOnOutletByName(globals.WATERPUMP, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlRootWater", globals.WATERPUMP, false, true)
		}
		break
	}
	return somethingChanged
}

func ControlAirflow(force bool, DeviceID int64, CurrentStage string, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {
	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlAirflow - no outlets attached")
		return
	}

	switch CurrentStage {
	case globals.GERMINATION:
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.OUTLETFAN, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlAirflow", globals.OUTLETFAN, true, false)
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.INLETFAN, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlAirflow", globals.INLETFAN, true, false)
		}
		break
	case globals.IDLE:
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.OUTLETFAN, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlAirflow", globals.OUTLETFAN, true, false)
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.INLETFAN, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlAirflow", globals.INLETFAN, true, false)
		}
		break
	default:
		if somethingChanged = Powerstrip.TurnOnOutletByName(globals.OUTLETFAN, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlAirflow", globals.INLETFAN, false, true)
		}
		if somethingChanged = Powerstrip.TurnOnOutletByName(globals.INLETFAN, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlAirflow", globals.INLETFAN, false, true)
		}
		break
	}
	return somethingChanged
}

func ReportSwitchStateChanged(functionName string, switchName string, originalState bool, newState bool) {
	log.Infof("StateChange: switch %s from %v to %v via %s", switchName, originalState, newState, functionName)
}

func ControlLight(force bool, DeviceID int64, CurrentStage string,
	MyStation globals.Station, CurrentStageSchedule globals.StageSchedule,
	LocalCurrentState *globals.LocalState, currentTime time.Time, Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {

	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlLight - no outlets attached")
		return somethingChanged
	}
	if CurrentStage == globals.IDLE {
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.GROWLIGHTVEG, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlLight", globals.GROWLIGHTVEG, true, false)
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.GROWLIGHTBLOOM, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlLight", globals.GROWLIGHTBLOOM, true, false)
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

	if CurrentStage == "germination" || CurrentStage == "seedling" || CurrentStage == "vegetative" || CurrentStage == "blooming" {
		// If it's time for grow light veg to be on
		if inRange(MyStation.LightOnHour, CurrentStageSchedule.HoursOfLight, localTimeHours) {
			log.Infof("automation: ControlLight turning on %s because local hour %d is within %d hours of %d", globals.GROWLIGHTBLOOM, localTimeHours, CurrentStageSchedule.HoursOfLight, MyStation.LightOnHour)
			if somethingChanged = Powerstrip.TurnOnOutletByName(globals.GROWLIGHTBLOOM, force); somethingChanged == true {
				ReportSwitchStateChanged("ControlLight", globals.GROWLIGHTBLOOM, false, true)
			}
			bloomlight = true
		} else {
			// If it's time for grow light veg to be off
			if LocalCurrentState.GrowLightBloom == true {
				log.Infof("automation: ControlLight turning off %s because local hour %d is outside %d hours of %d", globals.GROWLIGHTBLOOM, localTimeHours, CurrentStageSchedule.HoursOfLight, MyStation.LightOnHour)
			}
			if somethingChanged = Powerstrip.TurnOffOutletByName(globals.GROWLIGHTBLOOM, force); somethingChanged == true {
				ReportSwitchStateChanged("ControlLight", globals.GROWLIGHTBLOOM, true, false)
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
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.WATERHEATER, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlWaterTemp", globals.WATERHEATER, true, false)
		} // MAKE SURE HEAT IS OFF
		return somethingChanged
	}
	if ExternalCurrentState.TempF == globals.TEMPNOTSET {
		log.Infof("automation: ControlWaterTemp TEMPNOTSET ExternalCurrentState.WaterTempF %.3f - ignoring", ExternalCurrentState.WaterTempF)
		return somethingChanged
	}
	// Go from 62 to 68
	highLimit := StageSchedule.EnvironmentalTargets.WaterTemperature + 3.0
	lowLimit := StageSchedule.EnvironmentalTargets.WaterTemperature - 3.0
	if ExternalCurrentState.WaterTempF > highLimit {
		if *LastWaterTemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("automation: ControlWaterTemp turning off %s because WaterTemp (%.3f) just rolled over (%.2f/%.3f/%.2f) on way up", globals.WATERHEATER, ExternalCurrentState.WaterTempF, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit)
			force = true
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.WATERHEATER, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlWaterTemp", globals.WATERHEATER, true, false)
		}

		(*LocalCurrentState).WaterHeater = false
		(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)

	} else { // NOT TOO HOT
		if ExternalCurrentState.WaterTempF < lowLimit { // TOO COLD
			if *LastWaterTemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlWaterTemp turning on %s because Water Temp (%.3f) just fell below (%.2f/%.1f/%.2f) on way up from %.2f", globals.WATERHEATER, ExternalCurrentState.WaterTempF, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
				force = true
			}
			if somethingChanged = Powerstrip.TurnOnOutletByName(globals.WATERHEATER, force); somethingChanged == true {
				ReportSwitchStateChanged("ControlWaterTemp", globals.WATERHEATER, false, true)
			} // MAKE SURE HEAT IS ON

			(*LocalCurrentState).WaterHeater = true
		} else { // JUST RIGHT
			if *LastWaterTemp < lowLimit {
				log.Infof("automation: ControlWaterTemp Water Temp (%.3f) just entered sweet spot (%.2f/%.1f/%.2f) on way up from %.2f", ExternalCurrentState.WaterTempF, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
			} else {
				if *LastWaterTemp > highLimit {
					log.Infof("automation: ControlWaterTemp Water Temp (%.3f) just entered sweet spot (%.2f/%.1f/%.2f) on way down from %.3f", ExternalCurrentState.WaterTempF, lowLimit, StageSchedule.EnvironmentalTargets.WaterTemperature, highLimit, *LastWaterTemp)
				} else {
				}
			}
		}
	}
	*LastWaterTemp = ExternalCurrentState.WaterTempF
	return somethingChanged
}

func ControlHeat(force bool,
	DeviceID int64,
	CurrentStage string,
	CurrentStageSchedule globals.StageSchedule,
	ExternalCurrentState globals.ExternalState,
	LocalCurrentState *globals.LocalState,
	LastTemp *float32,
	Powerstrip gpiorelay.PowerstripService) (somethingChanged bool) {

	somethingChanged = false
	if !isPowerstripAttached(DeviceID) {
		log.Debugf("automation: ControlHeat - no outlets attached")
		return somethingChanged
	}
	//	log.Infof("automation: ControlHeat - current stage is %s", globals.MyStation.CurrentStage)
	if CurrentStage == globals.IDLE {
		log.Debugf("automation: ControlHeat - stage is idle, turning off")
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.HEATER, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlHeat", globals.HEATER, true, false)
		} // MAKE SURE HEAT IS OFF
		return somethingChanged
	}

	highLimit := CurrentStageSchedule.EnvironmentalTargets.Temperature + 2.0
	lowLimit := CurrentStageSchedule.EnvironmentalTargets.Temperature - 2.0

	//	log.Infof("automation: checking temp %.3f for stage %s with highLimit %.3f, lowLimit %.3f", globals.ExternalCurrentState.TempF, globals.MyStation.CurrentStage, highLimit,lowLimit)
	if ExternalCurrentState.TempF == globals.TEMPNOTSET {
		log.Debugf("automation: ControlHeat TEMPNOTSET ExternalCurrentState.TempF %.3f - ignoring", ExternalCurrentState.TempF)
		return somethingChanged
	}
	if ExternalCurrentState.TempF > highLimit { // TOO HOT
		log.Infof("automation: ControlHeat turning off %s because internal temp %.3f is over high limit %.3f on way up", globals.HEATER, ExternalCurrentState.TempF, highLimit)
		if globals.LastTemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("automation: ControlHeat turning off %s because internal temp %.3f just exceeded high limit %.3f on way up", globals.HEATER, ExternalCurrentState.TempF, highLimit)
			force = true
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.HEATER, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlHeat", globals.HEATER, true, false)
		} // MAKE SURE HEAT IS OFF

		(*LocalCurrentState).Heater = false
		(*LocalCurrentState).HeaterPad = false
		(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
	} else { // NOT TOO HOT
		if ExternalCurrentState.TempF < lowLimit { // TOO COLD
			//			log.Infof("automation: ControlHeat TOO COLD %.3f < lowLimit %.2f", ExternalCurrentState.TempF, lowLimit)
			if *LastTemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlHeat turning on %s because internal temp %.3f just fell below low limit %.3f on way down", globals.HEATER, ExternalCurrentState.TempF, lowLimit)
				force = true
			}
			if somethingChanged = Powerstrip.TurnOnOutletByName(globals.HEATER, false); somethingChanged == true {
				ReportSwitchStateChanged("ControlHeat", globals.HEATER, false, true)
			} // MAKE SURE HEAT IS ON

			(*LocalCurrentState).Heater = true
			(*LocalCurrentState).HeaterPad = true
		} else { // JUST RIGHT
			log.Infof("automation: ControlHeat JUST RIGHT %.3f", ExternalCurrentState.TempF)
			if *LastTemp < lowLimit {
				log.Infof("automation: ControlHeat Temp just entered sweet spot on way up - %.3f", ExternalCurrentState.TempF)
			} else {
				if *LastTemp > highLimit {
					log.Infof("automation: ControlHeat Temp just entered sweet spot on way down - %.3f", ExternalCurrentState.TempF)
				} else {
				}
			}
		}
	}

	(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
	*LastTemp = ExternalCurrentState.TempF
	return somethingChanged
}

func ControlHumidity(force bool,
	DeviceID int64,
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
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.HUMIDIFIER, false); somethingChanged == true {
			ReportSwitchStateChanged("ControlHumidity", globals.HUMIDIFIER, true, false)
		} // MAKE SURE HUMIDIFIER IS OFF
		return somethingChanged
	}
	highLimit := StageSchedule.EnvironmentalTargets.Humidity + 5.0
	lowLimit := StageSchedule.EnvironmentalTargets.Humidity - 5.0

	if ExternalCurrentState.Humidity == globals.HUMIDITYNOTSET {
		//		log.Debugf("automation: HUMIDITYNOTSET ExternalCurrentState.Humidity %.3f - ignoring", globals.ExternalCurrentState.Humidity))
		return somethingChanged
	}
	if ExternalCurrentState.Humidity > highLimit { // TOO HUMID
		if *LastHumidity < highLimit { // JUST BECAME TOO HUMID
			log.Infof("automation: ControlHumidity turning off %s because Humidity %.3f just exceeded (%.3f/%.1f/%.3f) on way up from %.3f", globals.HUMIDIFIER, ExternalCurrentState.Humidity, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
			force = true
		}
		if somethingChanged = Powerstrip.TurnOffOutletByName(globals.HUMIDIFIER, force); somethingChanged == true {
			ReportSwitchStateChanged("ControlHeat", globals.HUMIDIFIER, true, false)
		} // MAKE SURE HUMIDIFIER IS OFF
		(*LocalCurrentState).Humidifier = false
	} else { // NOT TOO HOT
		if ExternalCurrentState.Humidity < lowLimit { // TOO COLD
			if *LastHumidity > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlHumidity turning on %s because Humidity %.3f just fell below low (%.3f/%.1f/%.3f) on way down from %.3f", globals.HUMIDIFIER, ExternalCurrentState.Humidity, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
				force = true
			}
			if somethingChanged = Powerstrip.TurnOnOutletByName(globals.HUMIDIFIER, force); somethingChanged == true {
				ReportSwitchStateChanged("ControlHumidity", globals.HUMIDIFIER, false, true)
			} // MAKE SURE HUMIDIFIER IS ON

			(*LocalCurrentState).Humidifier = true
		} else { // JUST RIGHT
			if *LastHumidity < lowLimit {
				log.Infof("automation: ControlHumidity Humidity %.3f just entered sweet spot (%.3f/%.1f/%.3f) on way up from %.3f", ExternalCurrentState.Humidity, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
			} else {
				if *LastHumidity > highLimit {
					log.Infof("automation: ControlHumidity Humidity %.3f just entered sweet spot (%.3f/%.1f/%.3f) on way down from %.3f", ExternalCurrentState.Humidity, lowLimit, StageSchedule.EnvironmentalTargets.Humidity, highLimit, *LastHumidity)
				} else {
				}
			}
		}
	}

	(*LocalCurrentState).EnvironmentalControl = setEnvironmentalControlString(LocalCurrentState)
	*LastHumidity = ExternalCurrentState.Humidity
	return somethingChanged
}
