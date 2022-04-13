//go:build (linux && arm) || arm64
// +build linux,arm arm64

package camera

// copyright and license inspection - no issues 4/13/22

import (
	"fmt"
	"github.com/dhowden/raspicam"
	"github.com/go-playground/log"
	"os"
	"time"
)

func TakeAPicture() {

	//	log.Infof("takeAPicture()")
	t := time.Now()
	filename := fmt.Sprintf("%4.4d%2.2d%2.2d_%2.2d%2.2d_%2.2d.jpg", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	//	log.Debugf("Creating file %s", filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file: %#v", err)

		return
	}
	defer f.Close()

	//	log.Debugf("NewStill")
	s := raspicam.NewStill()
	errCh := make(chan error)
	go func() {
		for x := range errCh {
			log.Debugf("CAPTURE ERROR %#v", x)
		}
	}()
	//	log.Debugf("Capturing image...")
	raspicam.Capture(s, f, errCh)
	log.Debugf("Uploading picture %s", f.Name())
	err = uploadFile(f.Name())
	if err != nil {
		log.Errorf("os.Upload failed for %s", f.Name())
	}
	err = os.Remove(f.Name())
	if err != nil {
		log.Errorf("os.Remove failed for %s", f.Name())
	}
	SendPictureTakenEvent()

}
