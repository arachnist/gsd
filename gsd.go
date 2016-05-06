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
var statusbar []string

func timestamp(args map[string]string) string {
	return time.Now().Format(args["format"])
}

func fileReader(args map[string]string) string {
	buf, err := ioutil.ReadFile(args["path"])
	if err != nil {
		log.Println("Error reading", args["path"], err)
		return ""
	}
	content := string(buf)

	return content
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

	statusbar = make([]string, len(config.Items))

	for {
		for pos, item := range config.Items {
			switch item.Type {
			case "fileReader":
				statusbar[pos] = fileReader(item.Args)
			case "timestamp":
				statusbar[pos] = timestamp(item.Args)
			default:
				log.Fatalln("Unknown item type", item.Type)
			}
		}

		C.xname(C.CString(strings.Join(statusbar, config.Separator)))
		time.Sleep(time.Duration(1) * time.Second)
	}
}
