// +build darwin windows,amd64 linux,amd64

package main

import (
	"fmt"
	"gobot.io/x/gobot/drivers/i2c"
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

func clen(b []byte) (n int) {
	return( 0 )
}
// AtlasEZODriver is a driver for the BMP280 temperature/pressure sensor
type AtlasEZODriver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config

	tpc *bmp280CalibrationCoefficients
	Connection func() (err error)
	Name func() (name string)
	Halt func() (error)
	rawPh func() (float64, error)
	Ph func() (float64, error)
	initialization func() (error)
	read func(byte,int)([]byte, error)
}

func connection() (err error) {
	return nil
}


func NewAtlasEZODriver(c i2c.Connector, options ...func(i2c.Config)) *AtlasEZODriver {
	driver := AtlasEZODriver{
		name: "test",
		read: func( address byte, n int ) ( []byte, error ) {
			return []byte{}, nil
		},
		initialization: func()(error) {
			return nil
		},
		rawPh: func()(float64, error) {
			return 0,nil
		},
		Ph: func()(float64, error) {
			return 0,nil
		},
		Halt: func()(error) {
			fmt.Printf("Halt")
			return nil
		},
		Name: func() (string) {
			return "test"
		},
		Connection: func() (err error) {
			fmt.Printf("Connection")
			return nil},
	}
	return &driver
}

func (d *AtlasEZODriver) Start() (err error) {
	return nil
}

