package lawg

import (
	"bytes"
	"fmt"
	llog "github.com/go-playground/log"
	//	"github.com/go-playground/log"
	"os"
	"strings"
	"testing"
)

var t *testing.T
var Configured = false
var LogF *os.File
var logLevel = "warn,error,fatal,debug,info"

func ConfigureTestLogging( LogLevel string, _ string, tester *testing.T) {
	t = tester
	ConfigureLogging(LogLevel,".")
}

// Errorf logs an error log entry with formatting
func Errorf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "error") {
		ss := "ERROR: " + s
		fmt.Printf(ss, v ...)
		fmt.Printf("\n")
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss, v ...)
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

func Warn(s string) {
	if !Configured || strings.Contains(logLevel, "warn") {
		ss := "WARN: " + s
		fmt.Println(ss)
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss)
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

func Warnf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "warn"){
		ss := "WARN: " + s
		fmt.Printf(ss, v ...)
		fmt.Printf("\n")
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss, v ... )
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Debug(s string) {
	if !Configured || strings.Contains(logLevel, "debug") {
		ss := "DEBUG: " + s
		fmt.Println(ss)
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss )
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Debugf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "debug") {
		ss := "DEBUG: " + s
		fmt.Printf(ss, v ...)
		fmt.Printf("\n")
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss, v ... )
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}
// Errorf logs an error log entry with formatting
func Infof(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "info") {
		ss := "INFO: " + s
		fmt.Printf(ss, v ...)
		fmt.Printf("\n")
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss, v ... )
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Fatalf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "fatal") {
		llog.Fatalf(s, v...)
	}
}
// Errorf logs an error log entry with formatting
func Info(s string) {
	if !Configured || strings.Contains(logLevel, "info") {
		ss := "INFO: " + s
		fmt.Print(ss)
		fmt.Printf("\n")
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, ss )
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Error(e error) {
	if !Configured || strings.Contains(logLevel, "error") {
		fmt.Printf("%+v", e)
		fmt.Printf("\n")
		if LogF != nil {
			_,_ = fmt.Fprintf(LogF, "ERROR: %+v",e )
			_,_ = fmt.Fprintf(LogF, "\n")
		}
	}
}

func ConfigureLogging(LogLevel string, PersistentOutputDirectory string) {
	logLevel = LogLevel
	Configured = true

	if PersistentOutputDirectory != "" {
		var err error
		LogF, err = os.OpenFile(PersistentOutputDirectory+"/icebreaker_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			llog.Fatalf("error opening file: %v", err)
		}
	}

	Debug("debug")
	Info("info")
//	Notice("notice")
	Warn("warn")
	Errorf("error")
	// log.Panic("panic") // this will panic
//	Alert("alert")
}

func TestLog(e llog.Entry) {

	// below prints to os.Stderr but could marshal to JSON
	// and send to central logging server
	//																						       ---------
	// 				                                                                 |----------> | console |
	//                                                                               |             ---------
	// i.e. -----------------               -----------------     Unmarshal    -------------       --------
	//     | app log handler | -- json --> | central log app | --    to    -> | log handler | --> | syslog |
	//      -----------------               -----------------       Entry      -------------       --------
	//      																         |             ---------
	//                                  									         |----------> | DataDog |
	//
	//         																	        	   ---------
	b := new(bytes.Buffer)
	b.Reset()
	b.WriteString(e.Message)

	for _, f := range e.Fields {
		_, err := fmt.Fprintf(b, " %s=%v", f.Key, f.Value)
		fmt.Printf("logging error %+v\n", err)
	}
	t.Logf("%s\n", b.String())
}

