package accelerometer

// copyright and license inspection - no issues 4/13/22

type TamperDetectorService interface {
	RunTamperDetector(onceOnly bool)
}
