// +build linux,arm

package camera

import (
	"fmt"
	"github.com/dhowden/raspicam"
	"github.com/go-playground/log"
	"os"
	"time"
)

func TakeAPicture() {
	log.Infof("takeAPicture()")
	t := time.Now()
	filename := fmt.Sprintf("%4.4d%2.2d%2.2d_%2.2d%2.2d_%2.2d.jpg", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	log.Debugf("Creating file %s", filename )
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file: %v", err)
		return
	}
	defer f.Close()

	log.Debugf("NewStill")
	s := raspicam.NewStill()
	errCh := make(chan error)
	go func() {
		log.Infof("called err channel")
		for x := range errCh {
			log.Debugf( "CAPTURE ERROR %v", x)
		}
	}()
	log.Debugf("Capturing image...")
	raspicam.Capture(s, f, errCh)
	log.Debugf("skipping uploading %s", f.Name())
	uploadFile(f.Name())
	SendPictureTakenEvent()

}




