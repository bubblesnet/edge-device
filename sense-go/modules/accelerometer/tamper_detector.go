//go:build (linux && arm) || arm64
// +build linux,arm arm64

package accelerometer

import (
	"bubblesnet/edge-device/sense-go/globals"
	"github.com/go-playground/log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

var lastx = 0.0
var lasty = 0.0
var lastz = 0.0

var xmove = 0.0
var ymove = 0.0
var zmove = 0.0

func checkTamper() {
	x, y, z, _ := adxl345.XYZ()
	//			log.Debugf("x: %.7f | y: %.7f | z: %.7f \n", x, y, z))
	DidWeMove(x, y, z, false)
}

func DidWeMove(x int32, y int32, z int32, isUnitTest bool) {
	if lastx == 0.0 {
	} else {
		xmove = math.Abs(lastx - x)
		ymove = math.Abs(lasty - y)
		zmove = math.Abs(lastz - z)
		if xmove > globals.MyStation.TamperSpec.Xmove || ymove > globals.MyStation.TamperSpec.Ymove || zmove > globals.MyStation.TamperSpec.Zmove {
			log.Infof("new tamper message !! x: %.3f | y: %.3f | z: %.3f ", xmove, ymove, zmove)
			var tamperMessage = messaging.NewTamperSensorMessage("tamper_sensor",
				0.0, "", "", xmove, ymove, zmove)
			bytearray, err := json.Marshal(tamperMessage)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !isUnitTest {
				message := pb.SensorRequest{Sequence: globals.GetSequence(), TypeId: "sensor", Data: string(bytearray)}
				_, err = globals.Client.StoreAndForward(context.Background(), &message)
				if err != nil {
					log.Errorf("runTamperDetector ERROR %#v", err)
				}
			}
		}
	}
	lastx = x
	lasty = y
	lastz = z
}

func RunTamperDetector(onceOnly bool) {
	log.Info("runTamperDetector")
	adxl345Adaptor := raspi.NewAdaptor()
	adxl345 := i2c.NewADXL345Driver(adxl345Adaptor)

	work := func() {
		gobot.Every(100*time.Millisecond, checkTamper)
	}

	robot := gobot.NewRobot("adxl345Bot",
		[]gobot.Connection{adxl345Adaptor},
		[]gobot.Device{adxl345},
		work,
	)

	err := robot.Start()
	if err != nil {
		globals.ReportDeviceFailed("adxl345")
		log.Errorf("adxl345 robot start error %#v", err)
	}

	if onceOnly {
		robot.Stop()
	}

	if onceOnly {
		robot.Stop()
	}
}
