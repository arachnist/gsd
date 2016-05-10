// Copyright 2015 Robert S. Gerus. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

/*
// taken directly from https://github.com/Igneous/libxname/blob/master/xname.c
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
#include <stdio.h>
#include <stdlib.h>

int xname(const char *msg) {
	Display *dpy;
	Window rootwin;
	int scr;

	if(!(dpy=XOpenDisplay(NULL))) {
		fprintf(stderr, "ERROR: could not open display\n");
		exit(1);
	}

	scr = DefaultScreen(dpy);
	rootwin = RootWindow(dpy, scr);

	XStoreName(dpy, rootwin, msg);
	XCloseDisplay(dpy);

	return 0;
}
*/
import "C"

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	content := strings.TrimSpace(string(buf))

	if args["format"] != "" {
		content = fmt.Sprintf(args["format"], content)
	}

	if args["range_from"] != "" && args["range_to"] != "" && args["separator"] != "" {
		from, _ := strconv.Atoi(args["range_from"])
		to, _ := strconv.Atoi(args["range_to"])
		content = strings.Join(strings.Split(content, args["separator"])[from:to], args["separator"])
	}

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

	for i := 0; ; i++ {
		for pos, item := range config.Items {
			if (i % item.Period) == 0 {
				switch item.Type {
				case "fileReader":
					statusbar[pos] = fileReader(item.Args)
				case "timestamp":
					statusbar[pos] = timestamp(item.Args)
				default:
					log.Fatalln("Unknown item type", item.Type)
				}
			}
		}

		C.xname(C.CString(strings.Join(statusbar, config.Separator)))
		time.Sleep(time.Duration(1) * time.Second)
	}
}
