package a2dconverter

const MinInches = 0.0
const MaxInches = 12.5
const MaxVoltage = 2.65
const MinVoltage = 1.65

const ohmRange = 1600.0
const MinOhms = 400.0
const MaxOhms = 2400.0

const Etape_slope = 11.37795
const Etape_y_intercept = -17.28562

func etapeInchesToGallons(MaxInches float64, MaxGallons float64, inches float64) (gallons float64) {
	return inches * (MaxGallons / MaxInches)
}

func etapeInchesFromVolts(voltage float64, slope float64, yintercept float64) (inches float64) {
	// INCHES = VOLTS * -12.5 + 33.125
	inches = voltage*slope + yintercept
	return inches
}
