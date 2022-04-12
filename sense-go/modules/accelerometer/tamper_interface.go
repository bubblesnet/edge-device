package accelerometer

type TamperDetectorService interface {
	RunTamperDetector(onceOnly bool)
}
