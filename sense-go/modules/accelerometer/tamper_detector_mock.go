//go:build darwin || (windows && amd64) || (linux && amd64)
// +build darwin windows,amd64 linux,amd64

package accelerometer

var singletonTamperDetectorService = MockTamperDetector{Real: false}

type MockTamperDetector struct {
	Real bool
}

func GetTamperDetectorService() TamperDetectorService {
	return &singletonTamperDetectorService
}

func (r *MockTamperDetector) RunTamperDetector(onceOnly bool) {
}
