// +build linux,arm arm64

package phsensor

import (
	"bubblesnet/edge-device/sense-go/globals"
	"github.com/go-playground/log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"strconv"
	"time"
)

const (
	atlasEZOAddress            = 0x63
	bmp280RegisterControl      = 0xf4
	bmp280RegisterConfig       = 0xf5
	bmp280RegisterPressureData = 0xf7
	bmp280RegisterTempData     = 0xfa
	bmp280RegisterCalib00      = 0x88
	bmp280SeaLevelPressure     = 1013.25
)

type bmp280CalibrationCoefficients struct {
	t1 uint16
	t2 int16
	t3 int16
	p1 uint16
	p2 int16
	p3 int16
	p4 int16
	p5 int16
	p6 int16
	p7 int16
	p8 int16
	p9 int16
}

// AtlasEZODriver is a driver for the BMP280 temperature/pressure sensor
type AtlasEZODriver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config

	tpc *bmp280CalibrationCoefficients
}

// NewAtlasEZODriver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewAtlasEZODriver(c i2c.Connector, options ...func(i2c.Config)) *AtlasEZODriver {
	b := &AtlasEZODriver{
		name:      gobot.DefaultName("ATLASEZO"),
		connector: c,
		Config:    i2c.NewConfig(),
		tpc:       &bmp280CalibrationCoefficients{},
	}

	for _, option := range options {
		option(b)
	}

	// TODO: expose commands to API
	return b
}

// Name returns the name of the device.
func (d *AtlasEZODriver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *AtlasEZODriver) SetName(n string) {
	d.name = n
}

// Connection returns the connection of the device.
func (d *AtlasEZODriver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the BMP280 and loads the calibration coefficients.
func (d *AtlasEZODriver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(atlasEZOAddress)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		globals.ReportDeviceFailed("ezoph")
		log.Errorf("atlasezo getconnection error %#v", err)
		return err
	}

	if err := d.initialization(); err != nil {
		globals.ReportDeviceFailed("ezoph")
		log.Errorf("atlasezo initialization error %#v", err)
		return err
	}

	return nil
}

// Halt halts the device.
func (d *AtlasEZODriver) Halt() (err error) {
	return nil
}

// Ph returns the current temperature, in celsius degrees.
func (d *AtlasEZODriver) Ph() (pH float64, err error) {
	var rawP float64
	if rawP, err = d.rawPh(); err != nil {
		log.Errorf("Ph read error %#v", err)
		log.Errorf("atlasezo rawPh %#v", err)
		return 0.0, err
	}
	pH = rawP
	return pH, nil
}

// initialization reads the calibration coefficients.
func (d *AtlasEZODriver) initialization() (err error) {
	return nil
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func (d *AtlasEZODriver) rawPh() (pH float64, err error) {
	var data []byte

	if data, err = d.read(0x52, 256); err != nil {
		log.Errorf("atlasezo rawPh err %#v", err)
		return 0, err
	}
	d1 := data[:clen(data)]
	var s = string(d1)
	pH, err = strconv.ParseFloat(s, 64)

	return pH, nil
}

func (d *AtlasEZODriver) read(address byte, n int) ([]byte, error) {
	if _, err := d.connection.Write([]byte{address}); err != nil {
		log.Errorf("atlasezo write err %#v", err)
		return nil, err
	}
	// Documentation says wait 900ms between write and read, but 1000ms doesn't work while 2000ms does
	time.Sleep(1 * time.Second)
	buf := make([]byte, n)
	bytesRead, err := d.connection.Read(buf)
	if bytesRead != n || err != nil {
		log.Errorf("read %d bytes err = %#v", bytesRead, err)
		return nil, err
	}
	buflen := 0
	for i := 1; i < n && buf[i] != 0x0; i++ {
		buflen = buflen + 1
	}
	buf1 := make([]byte, buflen)
	for i := 1; i < buflen && buf[i] != 0x0; i++ {
		buf1[i-1] = buf[i]
	}
	return buf1, nil
}
