// Copyright © 2016 Dennis Chen <barracks510@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bcatlib

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// commandLookup searches for browser in the map of supported commands.
func commandLookup(browser string) (string, error) {
	command, exists := commands[browser]
	if !exists {
		return "", errors.New("browser not supported")
	}
	return command, nil
}

// A Browser holds the environment to manipulate various browsers.
type Browser struct {
	browser string
	command string
}

// New creates a new browser environment
func NewBrowser(browser, command string) (*Browser, error) {
	browser = strings.ToLower(browser)
	if command == "" {
		var err error
		command, err = commandLookup(browser)
		if err != nil {
			return nil, err
		}
	}
	b := &Browser{
		browser: browser,
		command: command,
	}
	return b, nil
}

// Open opens a browser window pointing to url
func (b *Browser) Open(url string) error {
	args := strings.Split(fmt.Sprintf("%s %s", b.command, url), " ")
	binary, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	wd, _ := os.Getwd()
	pa := &syscall.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: nil,
		Sys:   nil,
	}
	if _, err := syscall.ForkExec(binary, args, pa); err != nil {
		return err
	}
	return nil
}
