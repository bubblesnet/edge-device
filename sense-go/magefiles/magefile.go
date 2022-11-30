//go:build mage
// +build mage

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
 */

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"runtime"
	"strings"
)

// Creates the binary in the current directory.  It will overwrite any existing
// binary.
func Build(targetos string, targetarch string) error {
	fmt.Printf("Build(buildos=%s, targetos=%s, targetarch=%s)\n", runtime.GOOS, targetos, targetarch)

	envmap := make(map[string]string)
	envmap["GOOS"] = targetos
	envmap["GOARCH"] = targetarch
	envmap["GOARM"] = "7"
	//	envmap["CGO_ENABLED"] = "1"
	githash := ""
	timestamp := ""
	if runtime.GOOS == "windows" {
		githash, _ = sh.Output("magefiles\\getGithash.bat")
		githash = strings.ReplaceAll(githash, "\r", "")
		timestamp, _ = sh.Output("magefiles\\getTimestamp.bat")
		timestamp = strings.ReplaceAll(timestamp, "\r", "")
		timestamp = strings.ReplaceAll(timestamp, "'", "")
	} else if runtime.GOOS == "darwin" {
		githash, _ = sh.Output("magefiles\\getGithash.sh")
		githash = strings.ReplaceAll(githash, "\r", "")
		timestamp, _ = sh.Output("magefiles\\getTimestamp.sh")
		timestamp = strings.ReplaceAll(timestamp, "\r", "")
		timestamp = strings.ReplaceAll(timestamp, "'", "")
	}

	// go mod tidy
	ldf := "-X 'main.BubblesnetBuildNumberString=201' -X 'main.BubblesnetVersionMajorString=2' -X 'main.BubblesnetVersionMinorString=1' -X 'main.BubblesnetVersionPatchString=1'  -X 'main.BubblesnetGitHash=" + githash + "' -X 'main.BubblesnetBuildTimestamp=" + timestamp + "'"

	err := sh.RunWithV(envmap, "go", "build", "-tags", "make", "--ldflags="+ldf, "-o", "build", "./...")
	return err
}

func Coverage(targetos string, targetarch string) error {
	fmt.Printf("Coverage(buildos=%s, targetos=%s, targetarch=%s)\n", runtime.GOOS, targetos, targetarch)

	envmap := make(map[string]string)
	envmap["CODECOV_TOKEN"] = "bd6757f7-5f19-40b6-81f3-68547d5b9177"
	envmap["API_HOST"] = "192.168.23.237"
	envmap["NO_FAN_WITH_HEATER"] = "false"
	envmap["SLEEP_ON_EXIT_FOR_DEBUGGING"] = "60"
	envmap["ACTIVEMQ_HOST"] = "192.168.23.237"
	envmap["ACTIVEMQ_PORT"] = "61611"
	envmap["API_PORT"] = "4001"
	envmap["USERID"] = "90000009"
	envmap["DEVICEID"] = "70000008"
	envmap["GOOS"] = targetos
	envmap["GOARCH"] = targetarch
	envmap["GOARM"] = "7"

	if err := sh.RunWithV(envmap, "go", "test", "-tags", "make", "-coverprofile=coverage.out", "./..."); err != nil {
		return err
	}
	if err := sh.RunWithV(envmap, "go", "tool", "cover", "-tags", "make", "-html=coverage.out"); err != nil {
		return err
	}
	return nil
}
