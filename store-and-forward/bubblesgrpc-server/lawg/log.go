/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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

func ConfigureTestLogging(LogLevel string, _ string, tester *testing.T) {
	t = tester
	ConfigureLogging(LogLevel, "./testdata")
}

// Errorf logs an error log entry with formatting
func Errorf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "error") {
		ss := "ERROR: " + s
		fmt.Printf(ss, v...)
		fmt.Printf("\n")
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, ss, v...)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

func Warn(s string) {
	if !Configured || strings.Contains(logLevel, "warn") {
		ss := "WARN: " + s
		fmt.Println(ss)
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, ss)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

func Warnf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "warn") {
		ss := "WARN: " + s
		fmt.Printf(ss, v...)
		fmt.Printf("\n")
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, ss, v...)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Debug(s string) {
	if !Configured || strings.Contains(logLevel, "debug") {
		ss := "DEBUG: " + s
		fmt.Println(ss)
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, ss)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Debugf(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "debug") {
		ss := "DEBUG: " + s
		fmt.Printf(ss, v...)
		fmt.Printf("\n")
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, ss, v...)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Infof(s string, v ...interface{}) {
	if !Configured || strings.Contains(logLevel, "info") {
		ss := "INFO: " + s
		fmt.Printf(ss, v...)
		fmt.Printf("\n")
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, ss, v...)
			_, _ = fmt.Fprintf(LogF, "\n")
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
			_, _ = fmt.Fprintf(LogF, ss)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

// Errorf logs an error log entry with formatting
func Error(e error) {
	if !Configured || strings.Contains(logLevel, "error") {
		fmt.Printf("%+v", e)
		fmt.Printf("\n")
		if LogF != nil {
			_, _ = fmt.Fprintf(LogF, "ERROR: %+v", e)
			_, _ = fmt.Fprintf(LogF, "\n")
		}
	}
}

func ConfigureLogging(LogLevel string, PersistentOutputDirectory string) {
	logLevel = LogLevel
	Configured = true

	if PersistentOutputDirectory != "" {
		var err error
		LogF, err = os.OpenFile(PersistentOutputDirectory+"/bubblesnet_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
