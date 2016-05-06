// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

/*
#cgo LDFLAGS: -lxname
#include <xname.h>
*/
import "C"

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type statusbarItem struct {
	Name   string            `yaml:"name"`
	Period int               `yaml:"period"`
	Type   string            `yaml:"type"`
	Args   map[string]string `yaml:"args"`
}

type gsdConfig struct {
	Separator string          `yaml:"separator"`
	Items     []statusbarItem `yaml:"items"`
}

var config gsdConfig
var statusLock sync.Mutex
var statusbar []string

func updateStatusbar(pos int, text string) {
	statusLock.Lock()
	defer statusLock.Unlock()

	statusbar[pos] = text
	C.xname(C.CString(strings.Join(statusbar, config.Separator)))
}

func spawnUpdater(pos, period int, args map[string]string, f func(map[string]string) string) {
	go func() {
		for {
			time.Sleep(time.Duration(period) * time.Second)
			updateStatusbar(pos, f(args))
		}
	}()
}

func timestamp(args map[string]string) string {
	return "placeholder"
}

func fileReader(args map[string]string) string {
	return "placeholder"
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage:", os.Args[0], "<configuration file>")
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Error reading configuration file:", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalln("Error parsing configuration file:", err)
	}

}
