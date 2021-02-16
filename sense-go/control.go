package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	"bubblesnet/edge-device/sense-go/powerstrip"
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

func ControlLight() {
	localTimeHours := time.Now().Hour()
	offsetHours := 5
	if localTimeHours - offsetHours < 0 {
		localTimeHours = 24 + (localTimeHours - offsetHours)
	} else {
		localTimeHours = localTimeHours - offsetHours
	}
	veglight := false

	if globals.Config.Stage == "germination" || globals.Config.Stage == "seedling" || globals.Config.Stage == "vegetative" {
		// If it's time for grow light veg to be on
		if inRange(globals.Config.LightOnHour, globals.CurrentStageSchedule.HoursOfLight, localTimeHours) {
			powerstrip.TurnOnOutletByName(globals.GROWLIGHTVEG)
			veglight = true
		} else {
			// If it's time for grow light veg to be off
			powerstrip.TurnOffOutletByName(globals.GROWLIGHTVEG)
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

func ControlHeat() {

	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature + 2.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature - 2.0

	if globals.ExternalCurrentState.TempF == globals.TEMPNOTSET {
//		log.Debugf("TEMPNOTSET ExternalCurrentState.TempF %f - ignoring", globals.ExternalCurrentState.TempF))
		return
	}
	if globals.ExternalCurrentState.TempF > highLimit { // TOO HOT
		if globals.Lasttemp < highLimit { // JUST BECAME TOO HOT
			log.Infof("Temp just rolled over %f on way up %f", highLimit, globals.ExternalCurrentState.TempF)
		}
		powerstrip.TurnOffOutletByName(globals.HEATLAMP) // MAKE SURE HEAT IS OFF
		powerstrip.TurnOffOutletByName(globals.HEATPAD)  // MAKE SURE HEAT IS OFF
		powerstrip.TurnOffOutletByName(globals.HEATER)  // MAKE SURE HEAT IS OFF
		globals.LocalCurrentState.Heater = false
		globals.LocalCurrentState.HeaterPad = false
		setEnvironmentalControlString()
	} else {                                               // NOT TOO HOT
		if globals.ExternalCurrentState.TempF < lowLimit { // TOO COLD
			if globals.Lasttemp > lowLimit { // JUST BECAME TOO COLD
				log.Infof("Temp just fell below %f on way down - %f", lowLimit, globals.ExternalCurrentState.TempF)
			}
			powerstrip.TurnOnOutletByName(globals.HEATLAMP) // MAKE SURE HEAT IS ON
			powerstrip.TurnOnOutletByName(globals.HEATPAD)  // MAKE SURE HEAT IS ON
			powerstrip.TurnOnOutletByName(globals.HEATER)  // MAKE SURE HEAT IS ON
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

func ControlHumidity() {

	highLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity + 5.0
	lowLimit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity - 5.0

	if globals.ExternalCurrentState.Humidity == globals.HUMIDITYNOTSET {
//		log.Debugf("HUMIDITYNOTSET ExternalCurrentState.Humidity %f - ignoring", globals.ExternalCurrentState.Humidity))
		return
	}
	if globals.ExternalCurrentState.Humidity > highLimit { // TOO HUMID
		if globals.Lasthumidity < highLimit { // JUST BECAME TOO HUMID
			log.Infof("Humidity just rolled over %f on way up %f", highLimit, globals.ExternalCurrentState.Humidity)
		}
		powerstrip.TurnOffOutletByName(globals.HUMIDIFIER) // MAKE SURE HUMIDIFIER IS OFF
		globals.LocalCurrentState.Humidifier = false
	} else {                                                  // NOT TOO HOT
		if globals.ExternalCurrentState.Humidity < lowLimit { // TOO COLD
			if globals.Lasthumidity > lowLimit { // JUST BECAME TOO COLD
				log.Infof("Humidity just fell below %f on way down - %f", lowLimit, globals.ExternalCurrentState.Humidity)
			}
			powerstrip.TurnOnOutletByName(globals.HUMIDIFIER) // MAKE SURE HUMIDIFIER IS ON
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


