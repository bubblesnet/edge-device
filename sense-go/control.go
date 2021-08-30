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
		log.Debugf("ControlOxygenation - no relay attached")
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
		log.Debugf("ControlRootWater - no relay attached")
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
		log.Debugf("ControlAirflow - no relay attached")
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
		log.Debugf("ControlLight - no relay attached")
		return
	}
	if globals.MyStation.CurrentStage == globals.IDLE {
//		log.Debugf("ControlList - stage is idle, turning off")
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.GROWLIGHTVEG, false)
		gpiorelay.GetPowerstripService().TurnOffOutletByName(globals.GROWLIGHTBLOOM, false)
		return
	}
	localTimeHours := time.Now().Hour()
	offsetHours := 5
	if localTimeHours - offsetHours < 0 {
		localTimeHours = 24 + (localTimeHours - offsetHours)
	} else {
		localTimeHours = localTimeHours - offsetHours
	}
	veglight := false

	if globals.MyStation.CurrentStage == "germination" || globals.MyStation.CurrentStage == "seedling" || globals.MyStation.CurrentStage == "vegetative" {
		// If it's time for grow light veg to be on
		if inRange(globals.MyStation.LightOnHour, globals.CurrentStageSchedule.HoursOfLight, localTimeHours) {
			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.GROWLIGHTVEG, force)
			veglight = true
		} else {
			// If it's time for grow light veg to be off
			gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.GROWLIGHTVEG, force)
			veglight = false
		}
	} else {
	}
	if veglight && !globals.LocalCurrentState.GrowLightVeg {
		log.Infof("Turned veg light ON")
	} else if !veglight && globals.LocalCurrentState.GrowLightVeg {
		log.Infof("Turned veg light OFF")
	}
	globals.LocalCurrentState.GrowLightVeg = veglight
}

func inRange( starthour int, numhours int, currenthours int ) bool {
	if starthour + numhours >= 24  { // cross days
		if currenthours >= starthour {
			return true
		} else {
			if currenthours < (starthour+numhours-24) {
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

func ControlHeat(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("ControlHeat - no relay attached")
		return
	}
	log.Infof("ControlHeat - current stage is %s", globals.MyStation.CurrentStage)
	if globals.MyStation.CurrentStage == globals.IDLE {
		log.Debugf("ControlHeat - stage is idle, turning off")
//		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATLAMP, false) // MAKE SURE HEAT IS OFF
//		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATPAD, false)  // MAKE SURE HEAT IS OFF
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATER, false)   // MAKE SURE HEAT IS OFF
		return
	}

	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature + 2.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature - 2.0

//	log.Infof("checking temp %f for stage %s with highLimit %f, lowLimit %f", globals.ExternalCurrentState.TempF, globals.MyStation.CurrentStage, highLimit,lowLimit)
	if globals.ExternalCurrentState.TempF == globals.TEMPNOTSET {
		log.Debugf("TEMPNOTSET ExternalCurrentState.TempF %f - ignoring", globals.ExternalCurrentState.TempF)
		return
	}
	if globals.ExternalCurrentState.TempF > highLimit { // TOO HOT
		if globals.Lasttemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("Temp just rolled over %f on way up %f", highLimit, globals.ExternalCurrentState.TempF)
			force = true
		}
//		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATLAMP, force) // MAKE SURE HEAT IS OFF
//		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATPAD, force)  // MAKE SURE HEAT IS OFF
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HEATER, force)   // MAKE SURE HEAT IS OFF
		globals.LocalCurrentState.Heater = false
		globals.LocalCurrentState.HeaterPad = false
		setEnvironmentalControlString()
	} else {                                               // NOT TOO HOT
		if globals.ExternalCurrentState.TempF < lowLimit { // TOO COLD
			if globals.Lasttemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("Temp just fell below %f on way down - %f", lowLimit, globals.ExternalCurrentState.TempF)
				force = true
			}
//			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.HEATLAMP, false) // MAKE SURE HEAT IS ON
//			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.HEATPAD, false)  // MAKE SURE HEAT IS ON
			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.HEATER, false)   // MAKE SURE HEAT IS ON
			globals.LocalCurrentState.Heater = true
			globals.LocalCurrentState.HeaterPad = true
		} else { // JUST RIGHT
			if globals.Lasttemp < lowLimit {
				log.Infof("Temp just entered sweet spot on way up - %f", globals.ExternalCurrentState.TempF)
			} else {
				if globals.Lasttemp > highLimit {
					log.Infof("Temp just entered sweet spot on way down - %f", globals.ExternalCurrentState.TempF)
				} else {
				}
			}
		}
	}
	setEnvironmentalControlString()
	globals.Lasttemp = globals.ExternalCurrentState.TempF
}

func ControlHumidity(force bool) {
	if !isRelayAttached(globals.MyDevice.DeviceID) {
		log.Debugf("ControlHumidity - no relay attached")
		return
	}
	if globals.MyStation.CurrentStage == globals.IDLE {
//		log.Debugf("ControlHumidity - stage is idle, turning off")
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HUMIDIFIER, false) // MAKE SURE HUMIDIFIER IS OFF
		return
	}
	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity + 5.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity - 5.0

	if globals.ExternalCurrentState.Humidity == globals.HUMIDITYNOTSET {
//		log.Debugf("HUMIDITYNOTSET ExternalCurrentState.Humidity %f - ignoring", globals.ExternalCurrentState.Humidity))
		return
	}
	if globals.ExternalCurrentState.Humidity > highLimit { // TOO HUMID
		if globals.Lasthumidity < highLimit { // JUST BECAME TOO HUMID
			log.Infof("Humidity just rolled over %f on way up %f", highLimit, globals.ExternalCurrentState.Humidity)
			force = true
		}
		gpiorelay.PowerstripSvc.TurnOffOutletByName(globals.HUMIDIFIER, force) // MAKE SURE HUMIDIFIER IS OFF
		globals.LocalCurrentState.Humidifier = false
	} else {                                                  // NOT TOO HOT
		if globals.ExternalCurrentState.Humidity < lowLimit { // TOO COLD
			if globals.Lasthumidity > lowLimit { // JUST BECAME TOO COLD
				log.Infof("Humidity just fell below %f on way down - %f", lowLimit, globals.ExternalCurrentState.Humidity)
				force = true
			}

			gpiorelay.GetPowerstripService().TurnOnOutletByName(globals.HUMIDIFIER, force) // MAKE SURE HUMIDIFIER IS ON
			globals.LocalCurrentState.Humidifier = true
		} else { // JUST RIGHT
			if globals.Lasthumidity < lowLimit {
				log.Infof("Humidity just entered sweet spot on way up - %f", globals.ExternalCurrentState.Humidity)
			} else {
				if globals.Lasthumidity > highLimit {
					log.Infof("Humidity just entered sweet spot on way down - %f", globals.ExternalCurrentState.Humidity)
				} else {
				}
			}
		}
	}
	setEnvironmentalControlString()
	globals.Lasthumidity = globals.ExternalCurrentState.Humidity
}


