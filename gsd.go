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
	// "gopkg.in/yaml.v2"
	"sync"
	"time"
)

var statusLock sync.Mutex
var statusbar []string

func updateStatusbar(pos int, text string) {
	statusLock.Lock()
	defer statusLock.Unlock()

	statusbar[pos] = text
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
	C.xname(C.CString("placeholder"))
}
