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
	"os"
	"runtime"
)

// Creates the binary in the current directory.  It will overwrite any existing
// binary.
func Deploy(balenafleet string, targetos string, targetarch string) (err error) {
	fmt.Printf("Build(buildos=%s, targetos=%s, targetarch=%s)\n", runtime.GOOS, targetos, targetarch)

	cwd, _ := os.Getwd()
	if err = os.Chdir(cwd + "/sense-go"); err != nil {
		return err
	}
	cwd2, _ := os.Getwd()
	fmt.Printf("Current directory %s\n", cwd2)
	if err := sh.RunV("mage", "build", targetos, targetarch); err != nil {
		return err
	}

	if err = os.Chdir(cwd + "/store-and-forward/bubblesgrpc-server"); err != nil {
		return err
	}
	cwd1, _ := os.Getwd()
	fmt.Printf("Current directory %s\n", cwd1)
	if err := sh.RunV("mage", "build", targetos, targetarch); err != nil {
		return err
	}

	if err = os.Chdir(cwd); err != nil {
		return err
	}
	cwd3, _ := os.Getwd()
	fmt.Printf("Current directory %s\n", cwd3)

	if err := sh.RunV("balena", "push", balenafleet); err != nil {
		return err
	}
	return err
}
