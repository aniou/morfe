// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"github.com/jroimartin/gocui"
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/platform"
	"github.com/aniou/go65c816/emulator/tui"
)


func main() {
        g, err := gocui.NewGui(gocui.Output256)
        if err != nil {
                log.Panicln(err)
        }
        defer g.Close()

	logger := mylog.New()
	ui     := tui.New()
	p      := platform.New()

	ui.Init(g, logger, p)
	p.Init(logger)

	ui.Run(g)
}
