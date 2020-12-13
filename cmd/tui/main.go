
package main

import (
	"log"
	_ "github.com/jroimartin/gocui"
	_ "github.com/aniou/go65c816/lib/mylog"
	_ "github.com/aniou/go65c816/emulator/platform"
	_ "github.com/aniou/go65c816/emulator/tui"
)


func main() {
	log.Fatalf("TUI code doesn't work at this moment. " + 
	           "See branch https://github.com/aniou/go65c816/tree/tui for older version")


	/*
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
	*/
}
