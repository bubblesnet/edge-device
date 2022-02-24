package a2dconverter

import (
	"errors"
	"fmt"
	"github.com/go-playground/log"
)

func etapeInchesToGallons(inches float64) (gallons float64) {
	return inches * (12.0 / 18.0)
}

func etapeInchesFromVolts(Vdd float64, voltage float64, minResistance float64, maxResistance float64, minInches float64, maxInches float64) (inches float64, err error) {
	ohms, err := etapeVoltageToOhms(Vdd, voltage, minResistance, maxResistance)
	if err != nil {
		return 0.0, err
	}
	inches = etapeOhmsToInches(ohms)
	log.Infof("etapeInchesFromVolts returning %f inches for %f volts", inches, voltage)
	if inches < minInches || inches > maxInches {
		return 0, errors.New(fmt.Sprintf("returned Inches %f out of range %f - %f", inches, minInches, maxInches))
	}
	return inches, nil
}

func etapeVoltageToOhms(Vdd float64, voltage float64, minResistance float64, maxResistance float64) (ohms float64, err error) {
	maxVoltage := Vdd
	minVoltage := Vdd / 2.0

	v := voltage - minVoltage

	// y = mx + b
	x := v
	m := (maxResistance - minResistance) / (minVoltage - maxVoltage)
	b := maxResistance - minResistance
	y := m*x + b
	ohms = minResistance + y
	if ohms > maxResistance {
		return minResistance, errors.New(fmt.Sprintf("resistance %f out of range %f-%f ohms for voltage %f", ohms, minResistance, maxResistance, voltage))

	}
	log.Infof("etapeVoltageToOhms returning %f ohms for %f volts", ohms, voltage)
	return ohms, nil
}
func etapeOhmsToInches(ohms float64) (inches float64) {
	slope := -1800.0 / 11.5
	b := 2400.0
	// y = mx+b
	// mx = y-b
	// x = (y-b)/m
	inches = (ohms - b) / slope
	log.Infof("etapeOhmsToInches returning %f inches for %f ohms", inches, ohms)
	return inches
}
