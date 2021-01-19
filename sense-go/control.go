package main

import (
	"bubblesnet/edge-device/sense-go/globals"
	powerstrip "bubblesnet/edge-device/sense-go/powerstrip"
	"fmt"
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
		log.Info(fmt.Sprintf("Turned veg light ON"))
	} else if !veglight && globals.LocalCurrentState.GrowLightVeg {
		log.Info(fmt.Sprintf("Turned veg light OFF"))
	}
	globals.LocalCurrentState.GrowLightVeg = veglight
}

func inRange( starthour int, numhours int, currenthours int ) bool {
	if( starthour + numhours >= 24 ) {
		if currenthours >= globals.Config.LightOnHour && currenthours < (starthour+numhours) {
			return true
		} else {
			return false
		}
	} else {
		if currenthours >= globals.Config.LightOnHour || currenthours < ((starthour + numhours)-24) {
			return true
		} else {
			return false
		}
	}
}

func ControlHeat() {

	high_limit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature + 2.0
	low_limit := globals.CurrentStageSchedule.EnvironmentalTargets.Temperature - 2.0

	if globals.ExternalCurrentState.TempF == globals.TEMPNOTSET {
//		log.Debug(fmt.Sprintf("TEMPNOTSET ExternalCurrentState.TempF %f - ignoring", globals.ExternalCurrentState.TempF))
		return
	}
	if globals.ExternalCurrentState.TempF > high_limit { // TOO HOT
		if globals.Lasttemp < high_limit { // JUST BECAME TOO HOT
			log.Info(fmt.Sprintf("Temp just rolled over %f on way up %f", high_limit, globals.ExternalCurrentState.TempF))
		}
		powerstrip.TurnOffOutletByName(globals.HEATLAMP) // MAKE SURE HEAT IS OFF
		powerstrip.TurnOffOutletByName(globals.HEATPAD)  // MAKE SURE HEAT IS OFF
		globals.LocalCurrentState.Heater = false
		globals.LocalCurrentState.HeaterPad = false
		setEnvironmentalControlString()
	} else {                                                    // NOT TOO HOT
		if globals.ExternalCurrentState.TempF < low_limit { // TOO COLD
			if globals.Lasttemp > low_limit { // JUST BECAME TOO COLD
				log.Info(fmt.Sprintf("Temp just fell below %f on way down - %f", low_limit, globals.ExternalCurrentState.TempF))
			}
			powerstrip.TurnOnOutletByName(globals.HEATLAMP) // MAKE SURE HEAT IS ON
			powerstrip.TurnOnOutletByName(globals.HEATPAD)  // MAKE SURE HEAT IS ON
			globals.LocalCurrentState.Heater = true
			globals.LocalCurrentState.HeaterPad = true
		} else { // JUST RIGHT
			if globals.Lasttemp < low_limit  {
				log.Info(fmt.Sprintf("Temp just entered sweet spot on way up - %f", globals.ExternalCurrentState.TempF))
			} else {
				if globals.Lasttemp > high_limit {
					log.Info(fmt.Sprintf("Temp just entered sweet spot on way down - %f", globals.ExternalCurrentState.TempF))
				} else {
				}
			}
		}
	}
	setEnvironmentalControlString()
	globals.Lasttemp = globals.ExternalCurrentState.TempF
}

func ControlHumidity() {

	high_limit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity + 5.0
	low_limit := globals.CurrentStageSchedule.EnvironmentalTargets.Humidity - 5.0

	if globals.ExternalCurrentState.Humidity == globals.HUMIDITYNOTSET {
//		log.Debug(fmt.Sprintf("HUMIDITYNOTSET ExternalCurrentState.Humidity %f - ignoring", globals.ExternalCurrentState.Humidity))
		return
	}
	if globals.ExternalCurrentState.Humidity > high_limit { // TOO HUMID
		if globals.Lasthumidity < high_limit { // JUST BECAME TOO HUMID
			log.Info(fmt.Sprintf("Humidity just rolled over %f on way up %f", high_limit, globals.ExternalCurrentState.Humidity))
		}
		powerstrip.TurnOffOutletByName(globals.WATERPUMP) // MAKE SURE HUMIDIFIER IS OFF
		globals.LocalCurrentState.Humidifier = false
	} else {                                                       // NOT TOO HOT
		if globals.ExternalCurrentState.Humidity < low_limit { // TOO COLD
			if globals.Lasthumidity > low_limit { // JUST BECAME TOO COLD
				log.Info(fmt.Sprintf("Humidity just fell below %f on way down - %f", low_limit, globals.ExternalCurrentState.Humidity))
			}
			powerstrip.TurnOnOutletByName(globals.WATERPUMP) // MAKE SURE HUMIDIFIER IS ON
			globals.LocalCurrentState.Humidifier = true
		} else { // JUST RIGHT
			if globals.Lasthumidity < low_limit  {
				log.Info(fmt.Sprintf("Humidity just entered sweet spot on way up - %f", globals.ExternalCurrentState.Humidity))
			} else {
				if globals.Lasthumidity > high_limit {
					log.Info(fmt.Sprintf("Humidity just entered sweet spot on way down - %f", globals.ExternalCurrentState.Humidity))
				} else {
				}
			}
		}
	}
	setEnvironmentalControlString()
	globals.Lasthumidity = globals.ExternalCurrentState.Humidity
}


