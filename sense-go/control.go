package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/modules/gpiorelay"
	"github.com/go-playground/log"
	"time"
)

func setEnvironmentalControlString() {
	globals.LocalCurrentState.EnvironmentalControl = ""
	if globals.LocalCurrentState.Heater == true {
		globals.LocalCurrentState.EnvironmentalControl += "HEATING "
	} else {
		globals.LocalCurrentState.EnvironmentalControl += "COOLING "
	}
	if globals.LocalCurrentState.Humidifier == true {
		globals.LocalCurrentState.EnvironmentalControl += "HUMIDIFYING "
	} else {
		globals.LocalCurrentState.EnvironmentalControl += "DRYING "
	}
}

func ControlOxygenation(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: ControlOxygenation - no outlets attached")
		return
	}
	switch globals.MyStation.CurrentStage {
	case globals.IDLE:
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.AIRPUMP, false)
		break
	case globals.GERMINATION:
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.AIRPUMP, false)
		break
	default:
		gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.AIRPUMP, force)
		break
	}

}

func ControlRootWater(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: ControlRootWater - no outlets attached")
		return
	}

	if globals.MyStation.CurrentStage == globals.IDLE {
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.WATERPUMP, false)
		return
	}
	switch globals.MyStation.CurrentStage {
	case globals.IDLE:
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.WATERPUMP, false)
		break
	case globals.GERMINATION:
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.WATERPUMP, false)
		break
	default:
		gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.WATERPUMP, force)
		break
	}

}

func ControlAirflow(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: ControlAirflow - no outlets attached")
		return
	}

	switch globals.MyStation.CurrentStage {
	case globals.GERMINATION:
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.OUTLETFAN, false)
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.INLETFAN, false)
		break
	case globals.IDLE:
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.OUTLETFAN, false)
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.INLETFAN, false)
		break
	default:
		gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.OUTLETFAN, force)
		gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.INLETFAN, force)
		break
	}

}

func ControlLight(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: ControlLight - no outlets attached")
		return
	}
	if globals.MyStation.CurrentStage == globals.IDLE {
		//		log.Debugf("automation: ControlList - stage is idle, turning off")
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.GROWLIGHTVEG, false)
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.GROWLIGHTBLOOM, false)
		return
	}
	localTimeHours := time.Now().Hour()
	offsetHours := 5
	if localTimeHours-offsetHours < 0 {
		localTimeHours = 24 + (localTimeHours - offsetHours)
	} else {
		localTimeHours = localTimeHours - offsetHours
	}
	bloomlight := false

	if globals.MyStation.CurrentStage == "germination" || globals.MyStation.CurrentStage == "seedling" || globals.MyStation.CurrentStage == "vegetative" || globals.MyStation.CurrentStage == "blooming" {
		// If it's time for grow light veg to be on
		if inRange(globals.MyStation.LightOnHour, globals.CurrentStageSchedule.HoursOfLight, localTimeHours) {
			log.Infof("automation: ControlLight turning on %s because local hour %d is within %d hours of %d", globals.GROWLIGHTBLOOM, localTimeHours, globals.CurrentStageSchedule.HoursOfLight, globals.MyStation.LightOnHour)
			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.GROWLIGHTBLOOM, force)
			bloomlight = true
		} else {
			// If it's time for grow light veg to be off
			if globals.LocalCurrentState.GrowLightBloom == true {
				log.Infof("automation: ControlLight turning off %s because local hour %d is outside %d hours of %d", globals.GROWLIGHTBLOOM, localTimeHours, globals.CurrentStageSchedule.HoursOfLight, globals.MyStation.LightOnHour)
			}
			gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.GROWLIGHTBLOOM, force)
			bloomlight = false
		}
	} else {
	}
	if bloomlight && !globals.LocalCurrentState.GrowLightBloom {
		log.Infof("automation: ControlLight Turned veg light ON")
	} else if !bloomlight && globals.LocalCurrentState.GrowLightBloom {
		log.Infof("automation: ControlLight Turned veg light OFF")
	}
	globals.LocalCurrentState.GrowLightBloom = bloomlight
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

func ControlWaterTemp(force bool, DeviceID int64, Station globals.Station, Stage globals.StageSchedule,
	CurrentStage string, ExternalCurrentState globals.ExternalState) {
	if !isRelayAttached(DeviceID) {
		log.Infof("automation: ControlWaterTemp - no outlets attached")
		return
	}
	if CurrentStage == globals.IDLE {
		log.Infof("automation: ControlWaterTemp - stage is idle, turning off")
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.WATERHEATER, false) // MAKE SURE HEAT IS OFF
		return
	}
	if globals.ExternalCurrentState.TempF == globals.TEMPNOTSET {
		log.Infof("automation: ControlWaterTemp TEMPNOTSET ExternalCurrentState.WaterTempF %f - ignoring", globals.ExternalCurrentState.WaterTempF)
		return
	}
	// Go from 62 to 68
	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.WaterTemperature + 3.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.WaterTemperature - 3.0
	if globals.ExternalCurrentState.WaterTempF > highLimit {
		if globals.LastWaterTemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("automation: ControlWaterTemp turning off %s because WaterTemp (%f) just rolled over %f on way up", globals.WATERHEATER, globals.ExternalCurrentState.WaterTempF, highLimit)
			force = true
		}
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.WATERHEATER, force) // MAKE SURE HEAT IS OFF

		globals.LocalCurrentState.WaterHeater = false
		setEnvironmentalControlString()

	} else {                                                    // NOT TOO HOT
		if globals.ExternalCurrentState.WaterTempF < lowLimit { // TOO COLD
			if globals.LastWaterTemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlWaterTemp turning on %s because WaterTemp just fell below %f on way down - %f", globals.WATERHEATER, lowLimit, globals.ExternalCurrentState.WaterTempF)
				force = true
			}
			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.WATERHEATER, false) // MAKE SURE HEAT IS ON

			globals.LocalCurrentState.WaterHeater = true
		} else { // JUST RIGHT
			if globals.LastWaterTemp < lowLimit {
				log.Infof("automation: ControlWaterTemp Water Temp just entered sweet spot on way up - %f", globals.ExternalCurrentState.WaterTempF)
			} else {
				if globals.LastWaterTemp > highLimit {
					log.Infof("automation: ControlWaterTemp Water Temp just entered sweet spot on way down - %f", globals.ExternalCurrentState.WaterTempF)
				} else {
				}
			}
		}
	}
}

func ControlHeat(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: ControlHeat - no outlets attached")
		return
	}
	//	log.Infof("automation: ControlHeat - current stage is %s", globals.MyStation.CurrentStage)
	if globals.MyStation.CurrentStage == globals.IDLE {
		log.Debugf("automation: ControlHeat - stage is idle, turning off")
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATER, false) // MAKE SURE HEAT IS OFF
		return
	}

	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature + 2.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature - 2.0

	//	log.Infof("automation: checking temp %f for stage %s with highLimit %f, lowLimit %f", globals.ExternalCurrentState.TempF, globals.MyStation.CurrentStage, highLimit,lowLimit)
	if globals.ExternalCurrentState.TempF == globals.TEMPNOTSET {
		log.Debugf("automation: ControlHeat TEMPNOTSET ExternalCurrentState.TempF %f - ignoring", globals.ExternalCurrentState.TempF)
		return
	}
	if globals.ExternalCurrentState.TempF > highLimit { // TOO HOT
		if globals.LastTemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("automation: ControlHeat turning off %s because internal temp %f just exceeded high limit %f on way up", globals.HEATER, globals.ExternalCurrentState.TempF, highLimit)
			force = true
		}
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATER, force) // MAKE SURE HEAT IS OFF

		globals.LocalCurrentState.Heater = false
		globals.LocalCurrentState.HeaterPad = false
		setEnvironmentalControlString()
	} else {                                               // NOT TOO HOT
		if globals.ExternalCurrentState.TempF < lowLimit { // TOO COLD
			if globals.LastTemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlHeat turning on %s because internal temp %f just fell below low limit %f on way down", globals.HEATER, globals.ExternalCurrentState.TempF, lowLimit)
				force = true
			}
			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.HEATER, false) // MAKE SURE HEAT IS ON

			globals.LocalCurrentState.Heater = true
			globals.LocalCurrentState.HeaterPad = true
		} else { // JUST RIGHT
			if globals.LastTemp < lowLimit {
				log.Infof("automation: ControlHeat Temp just entered sweet spot on way up - %f", globals.ExternalCurrentState.TempF)
			} else {
				if globals.LastTemp > highLimit {
					log.Infof("automation: ControlHeat Temp just entered sweet spot on way down - %f", globals.ExternalCurrentState.TempF)
				} else {
				}
			}
		}
	}
	setEnvironmentalControlString()
	globals.LastTemp = globals.ExternalCurrentState.TempF
}

func ControlHumidity(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("automation: ControlHumidity - no outlets attached")
		return
	}
	if globals.MyStation.CurrentStage == globals.IDLE {
		//		log.Debugf("automation: ControlHumidity - stage is idle, turning off")
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HUMIDIFIER, false) // MAKE SURE HUMIDIFIER IS OFF
		return
	}
	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity + 5.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity - 5.0

	if globals.ExternalCurrentState.Humidity == globals.HUMIDITYNOTSET {
		//		log.Debugf("automation: HUMIDITYNOTSET ExternalCurrentState.Humidity %f - ignoring", globals.ExternalCurrentState.Humidity))
		return
	}
	if globals.ExternalCurrentState.Humidity > highLimit { // TOO HUMID
		if globals.LastHumidity < highLimit { // JUST BECAME TOO HUMID
			log.Infof("automation: ControlHumidity turning off %s because Humidity %f just exceeded high limit %f on way up", globals.HUMIDIFIER, globals.ExternalCurrentState.Humidity, highLimit)
			force = true
		}
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HUMIDIFIER, force) // MAKE SURE HUMIDIFIER IS OFF
		globals.LocalCurrentState.Humidifier = false
	} else {                                                  // NOT TOO HOT
		if globals.ExternalCurrentState.Humidity < lowLimit { // TOO COLD
			if globals.LastHumidity > lowLimit { // JUST BECAME TOO COLD
				log.Infof("automation: ControlHumidity turning on %s because Humidity %f just fell below low limit %f on way down", globals.HUMIDIFIER, globals.ExternalCurrentState.Humidity, highLimit)
				force = true
			}
			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.HUMIDIFIER, force) // MAKE SURE HUMIDIFIER IS ON

			globals.LocalCurrentState.Humidifier = true
		} else { // JUST RIGHT
			if globals.LastHumidity < lowLimit {
				log.Infof("automation: ControlHumidity Humidity %f just entered sweet spot on way up", globals.ExternalCurrentState.Humidity)
			} else {
				if globals.LastHumidity > highLimit {
					log.Infof("automation: ControlHumidity Humidity %f just entered sweet spot on way down", globals.ExternalCurrentState.Humidity)
				} else {
				}
			}
		}
	}
	setEnvironmentalControlString()
	globals.LastHumidity = globals.ExternalCurrentState.Humidity
}
