
package tui

import (
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"

	"github.com/jroimartin/gocui"
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/platform"
)

type Ui struct {
	viewArr		 []string	// names of views to cycle
	active		 int		// index in viewArr

	memPosition	 uint32 	// displayed memory region

	isPageMapVisible bool
	isCPUrunning     bool		// XXX move it to CPU
	lastView	 string		// last view.Name when cmd is called
	cpuSpeed	 uint64


	stepControl	 chan byte
	logQuit		 chan bool	// quit signal channel

	logger		 *mylog.MyLog
	p		 *platform.Platform	// pointer to Platform
	logView		 *gocui.View	// shortcut for current UI
}

func New() *Ui {
	ui            := Ui{}
	return &ui
}

func (ui *Ui) Init(g *gocui.Gui, logger *mylog.MyLog, platform *platform.Platform) {
        ui.active      = 0
	ui.viewArr     = []string{"code", "cmd"}
	ui.stepControl = make(chan byte)
	ui.logger      = logger
        ui.logQuit     = make(chan bool)
	ui.p           = platform


	g.Cursor     = true
	g.SelFgColor = gocui.ColorGreen
	g.Highlight  = true

	g.SetManagerFunc(ui.Layout)
	if err := ui.keybindings(g); err != nil {
		log.Panicln(err)
	}

	go ui.Logger(g)
	logger.Log("TUI init complete")
}



func (ui *Ui) Run(g *gocui.Gui) {
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}


func (ui *Ui) Logger(g *gocui.Gui) {
	t := time.NewTimer(time.Millisecond * 250)

	for {
		select {
		case <-t.C:
			if ui.logger.Len() > 0 {
				g.Update(func(g *gocui.Gui) error {

					v, err := g.View("log")
					// XXX - change to something like logger.Flush()
					// and avoid exposing internals of logger to tui
					if err == nil {
						for ui.logger.Len() > 0 {
							fmt.Fprintf(v, "> %s\n", *ui.logger.Dequeue())
						}
					}
					return err
				})
			}
			t.Reset(time.Millisecond * 250)
		case <-ui.logQuit:
			if !t.Stop() {
				<-t.C
			}
			break
		}
	}
}





//
// misc execution-related functions
//
func (ui *Ui) loadProg(g *gocui.Gui, v *gocui.View) error {
	ui.p.LoadHex("data/program.hex")
	ui.p.CPU.PC = 0x1000
	ui.p.CPU.E  = 0
	ui.p.CPU.M  = 0
	ui.p.CPU.X  = 0
	return nil
}

func (ui *Ui) runCPU(g *gocui.Gui, v *gocui.View) error {
	if ui.isCPUrunning {
		ui.stepControl<-1
		return nil
	}
	ui.isCPUrunning = true

	go func() {
		t := time.NewTimer(time.Millisecond * 250)
		//e := time.NewTimer(time.Second * 60)

		g.Update(func(g *gocui.Gui) error {
			fmt.Fprintf(ui.logView, "runCPU: start\n")
			return nil
		})

		var a byte = 0
		var b uint64 = 0
		var previousCycles uint64 = ui.p.CPU.AllCycles
		for a == 0 {
			select {
			case <-ui.stepControl:
				a = 1
			case <-t.C:
				ui.cpuSpeed = (ui.p.CPU.AllCycles - previousCycles) * 4 // XXX - parametrize that, 250 ms now
				previousCycles = ui.p.CPU.AllCycles
				g.Update(func(g *gocui.Gui) error {
					ui.updateStatusView(g)
					ui.updateStackView(g)
					ui.updateMemoryView(g)
					if ui.isPageMapVisible {
						ui.updatePageMap(g)
					}
					return nil
				})
				t.Reset(time.Millisecond * 250)
			//case <-e.C:
			//	a = 1
			default:
				_, abort := ui.p.CPU.Step()
				b++
				if abort {
					a = 1
				}
			}

		}

		g.Update(func(g *gocui.Gui) error {
			fmt.Fprintf(ui.logView, "runCPU: stop after %d cycles\n", b)
			return nil
		})

		if !t.Stop() {
			<-t.C
		}

		//if !e.Stop() {
		//	<-e.C
		//}

		ui.isCPUrunning = false
	}()


	return nil
}

func (ui *Ui) executeStep(g *gocui.Gui, v *gocui.View) error {
	if ui.isCPUrunning {
		return nil
	}

	ui.p.CPU.Step()

	ui.updateStatusView(g)
	ui.updateCodeView(g)
	ui.updateMemoryView(g)
	ui.updateStackView(g)

	return nil
}

//
// command view support
//
func hex2uint24(hexStr string) (uint32, error) {
	// remove 0x suffix, $ and : characters
	cleaned := strings.Replace(hexStr,  "0x", "",  1)
	cleaned  = strings.Replace(cleaned,  "$", "",  1)
	cleaned  = strings.Replace(cleaned,  ":", "", -1)

	result, err := strconv.ParseUint(cleaned, 16, 32)
	return uint32(result) & 0x00ffffff, err
}

// not finished yet
func dec2uint24(in string) (uint32, error) {
	// remove 0x suffix, $ and : characters
	cleaned := strings.Replace(in,      "0d", "",  1)
	cleaned  = strings.Replace(cleaned,  "_", "",  1)

	result, err := strconv.ParseUint(cleaned, 10, 32)
	return uint32(result) & 0x00ffffff, err
}

func (ui *Ui) setParameter(g *gocui.Gui, tokens []string) {
	var err error

	switch tokens[1] {
	case "mem":
		if ui.memPosition, err = hex2uint24(tokens[2]); err == nil {
			ui.updateMemoryView(g)
		} else {
			fmt.Fprintf(ui.logView, "set: error: %s\n", err)
		}
	default:
		fmt.Fprintf(ui.logView, "set: unknown parameter '%s'\n", tokens[2])
	}
}

func (ui *Ui) loadProgram(g *gocui.Gui, tokens []string) {
	switch tokens[1] {
	case "hex":
		ui.p.LoadHex(tokens[2])
		ui.p.CPU.PC = 0x1000
		ui.p.CPU.E  = 0
		ui.p.CPU.M  = 0
		ui.p.CPU.X  = 0

        ui.updateStatusView(g)
        ui.updateCodeView(g)
        ui.updateMemoryView(g)
        ui.updateStackView(g)
	default:
		fmt.Fprintf(ui.logView, "load: unknown parameter '%s'\n", tokens)
	}
}

func (ui *Ui) peek(g *gocui.Gui, tokens []string) {
	var ea uint32
	var err error

	if ea, err = hex2uint24(tokens[1]); err != nil {
		fmt.Fprintf(ui.logView, "set: error: %s\n", err)
		return
	}

	switch tokens[0] {
	case "peek", "peek8":
		val := ui.p.CPU.Bus.EaRead(ea)
		fmt.Fprintf(ui.logView, "peek %06x = %6x\n", ea, val)
	case "peek16":
		ll  := ui.p.CPU.Bus.EaRead(ea)
		hh  := ui.p.CPU.Bus.EaRead(ea+1)
		val := uint16(hh) << 8 | uint16(ll)
		fmt.Fprintf(ui.logView, "peek %06x = %6x\n", ea, val)
	case "peek24":
		ll  := ui.p.CPU.Bus.EaRead(ea)
		mm  := ui.p.CPU.Bus.EaRead(ea+1)
		hh  := ui.p.CPU.Bus.EaRead(ea+2)
		val := uint32(hh) << 16 | uint32(mm) << 8 | uint32(ll)
		fmt.Fprintf(ui.logView, "peek %06x = %6x\n", ea, val)
	}
}

func (ui *Ui) executeCmd(g *gocui.Gui, v *gocui.View) error {
        command := strings.TrimSpace(v.Buffer())
	tokens := strings.Split(command, " ")
	switch tokens[0] {
	case "se", "set":
		ui.setParameter(g, tokens)
	case "lo", "load":
		ui.loadProgram(g, tokens)
	case "run":
		ui.runCPU(g, v)
	case "peek", "peek8", "peek16", "peek24":
		ui.peek(g, tokens)
	case "quit":
		ui.quit(g, v)
		return gocui.ErrQuit
	default:
		fmt.Fprintf(ui.logView, "unknown command: %s\n", command)
	}

        v.Clear()
	v.SetCursor(0,0)
        return nil
}

func (ui *Ui) toggleCmdView(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "cmd" {
		if _, err := setCurrentViewOnTop(g, ui.lastView); err != nil {
			return err
		}
	} else {
		ui.lastView = v.Name()
		if _, err := setCurrentViewOnTop(g, "cmd"); err != nil {
			return err
		}
	}

	return nil
}

//
// statusView
//
func printCPUFlags(flag byte, name string) (string) {
	if flag > 0 {
		return name
	} else {
		return "-"
	}
}

func showCPUSpeed(cycles uint64) (uint64, string) {
	switch {
	case cycles > 1000000:
		return cycles/1000000, "MHz"
	case cycles > 1000:
		return cycles/100, "kHz"
	default:
		return cycles, "Hz"
	}
}

func (ui *Ui) updateStatusView(g *gocui.Gui) error {

	v, err := g.View("status")
	if err != nil {
		fmt.Fprintf(ui.logView, "%s\n", err)
		return err
	}

	v.Clear()

	// first line
	if ui.p.CPU.M == 0 {
		fmt.Fprintf(v, " A  %04x (%7d) │ SP   %04x │  PC %02x:%04x │",
			ui.p.CPU.RA, ui.p.CPU.RA, ui.p.CPU.SP, ui.p.CPU.RK, ui.p.CPU.PC)
	} else {
		fmt.Fprintf(v, " A %02x %02x (%3d %3d) │ SP   %04x │  PC %02x:%04x │",
			ui.p.CPU.RAh, ui.p.CPU.RAl, ui.p.CPU.RAh, ui.p.CPU.RAl, ui.p.CPU.SP, ui.p.CPU.RK, ui.p.CPU.PC)

	}
	fmt.Fprintf(v, " cycle %11d │\n", ui.p.CPU.AllCycles)


	speed, suffix := showCPUSpeed(ui.cpuSpeed)

	// second and third line
	if ui.p.CPU.X == 0 {
		fmt.Fprintf(v, " X  %04x (%7d) │ DBR    %02x │ pPC 00:%04x │                   │\n",
			ui.p.CPU.RX, ui.p.CPU.RX, ui.p.CPU.RDBR, ui.p.CPU.PPC)
		fmt.Fprintf(v, " Y  %04x (%7d) │ D    %04x │ ",
			ui.p.CPU.RY, ui.p.CPU.RY, ui.p.CPU.RD)
	} else {
		fmt.Fprintf(v, " X    %02x (    %3d) │ DBR    %02x │ pPC 00:%04x │                   │\n",
			ui.p.CPU.RXl, ui.p.CPU.RXl, ui.p.CPU.RDBR, ui.p.CPU.PPC)
		fmt.Fprintf(v, " Y    %02x (    %3d) │ D    %04x │ ",
			ui.p.CPU.RYl, ui.p.CPU.RYl, ui.p.CPU.RD)
	}

	// third line
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.N, "n"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.V, "v"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.M, "m"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.X, "x"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.D, "d"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.I, "i"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.Z, "z"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.C, "c"))
	fmt.Fprintf(v, " ")
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.B, "B"))
	fmt.Fprintf(v, printCPUFlags(ui.p.CPU.E, "E"))
	fmt.Fprintf(v, " │")
	fmt.Fprintf(v, " speed  %6d %3s │", speed, suffix)



	return nil
}

//
// codeView
//
func (ui *Ui) updateCodeView(g *gocui.Gui) error {
	//var numeric string

	v, err := g.View("code")
	if err != nil {
		fmt.Fprintf(ui.logView, "%s\n", err)
		return err
	}

	if ui.p.CPU.Cycles == 0 {
		fmt.Fprintf(v, "--:----│           │                 │")

	}
	if ui.p.CPU.Cycles > 9 {
		fmt.Fprintf(ui.logView, "warning: instruction cycles > 10\n")
	}

	line := ui.p.CPU.DisassembleCurrentPC()
	//fmt.Fprintf(v, "%d\n%02x:%04x│%-11v│%3s %-13v│",
	//				ui.p.CPU.Cycles, ui.p.CPU.RK, myPC, numeric, name, arg)
	//fmt.Fprintf(v, "%-38v",   "3│00:000c│02 02      │BEQ 02 ($04fa +)")
	fmt.Fprintf(v, line)

	return nil
}


//
// memoryView
//
func (ui *Ui) updateMemoryView(g *gocui.Gui) error {

	v, err := g.View("memory")
	if err != nil {
		fmt.Fprintf(ui.logView, "%s\n", err)
		return err
	}

	v.Clear()

	var x uint16
	var a uint16

	for a = 0; a < 0x100; a=a+16 {
		start, data := ui.p.CPU.Bus.EaDump(ui.memPosition+uint32(a))
		bank := byte(start >> 16)
		addr := uint16(start)
		fmt.Fprintf(v, "\n%02x:%04x│", bank, addr)
		if data != nil {
			fmt.Fprintf(v, "% x│% x│", data[0:8], data[8:16])
			for x = 0; x < 16; x++  {
				if data[x] >= 33 && data[x] < 127 {
					fmt.Fprintf(v, "%s", data[x:x+1])
				} else {
					fmt.Fprintf(v, ".")
				}
				if x == 7 {
					fmt.Fprintf(v, " ")
				}
			}
		} else {
			fmt.Fprintf(v, "                       │                       │")
		}
	}

	return nil
}

//
// stackView
//
func (ui *Ui) updateStackView(g *gocui.Gui) error {

	v, err := g.View("stack")
	if err != nil {
		fmt.Fprintf(ui.logView, "%s\n", err)
		return err
	}

	v.Clear()

	_, rows  := v.Size()
	halfrows := uint16(rows) >> 1     // half of rows?

	///fmt.Fprintf(v, "%04x\n", 0x100 + p.CPU.SP )

	for a := ui.p.CPU.SP+halfrows; a > ui.p.CPU.SP-halfrows; a=a-1 {
		val := ui.p.CPU.Bus.EaRead(uint32(a))
		if val == 0 {
			fmt.Fprintf(v, "\n%04x \x1b[38;5;236m%02x\x1b[0m", a, val) 
		} else {
			fmt.Fprintf(v, "\n%04x %02x", a, val)
		}
	}
	v.SetCursor(0, int(halfrows))


	return nil
}

//
// pageMap functions
//
func (ui *Ui) togglePageMap(g *gocui.Gui, v *gocui.View) error {
	if ui.isPageMapVisible {
		ui.isPageMapVisible = false
		ui.closePageMap(g)
	} else {
		ui.isPageMapVisible = true
	}
	return nil
}

func (ui *Ui) showPageMap(g *gocui.Gui) error {
	//maxX, maxY := g.Size()
	if v, err := g.SetView("pagemap", 1, 8, 64+2, 16+1+8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = false
		v.Title = "Page Map"
	}
	g.SetViewOnTop("pagemap")

	return nil
}

func (ui *Ui) closePageMap(g *gocui.Gui) error {
	g.DeleteView("pagemap")
	ui.isPageMapVisible = false
	return nil
}

func (ui *Ui) updatePageMap(g *gocui.Gui) error {
        v, err := g.View("pagemap")
        if err != nil {
                fmt.Fprintf(ui.logView, "%s\n", err)
                return err
        }

	//var out string
	v.Clear()
	/*
	for x:=0; x<1024; x++ {
		switch p.mapPage[x] {
		case 1:
			out = "r"
		case 2:
			out = "W"
		default:
			out = "."
		}
		fmt.Fprintf(v, out)
	}
	*/

	return nil
}


//
// common ui-support functions
//
func (ui *Ui) nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (ui.active + 1) % len(ui.viewArr)
        name := ui.viewArr[nextIndex]

        if _, err := setCurrentViewOnTop(g, name); err != nil {
                return err
        }

        ui.active = nextIndex
	return nil
}

func (ui *Ui) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ui *Ui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
        if _, err := g.SetCurrentView(name); err != nil {
                return nil, err
        }
        return g.SetViewOnTop(name)
}

func (ui *Ui) quit(g *gocui.Gui, v *gocui.View) error {
	ui.logQuit<-true	// terminate logger
	return gocui.ErrQuit
}

func (ui *Ui) keybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("",  gocui.KeyTab, gocui.ModNone, ui.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, ui.runCPU); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, ui.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, ui.loadProg); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '`', gocui.ModNone, ui.toggleCmdView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, ui.togglePageMap); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyEnter, gocui.ModNone, ui.executeCmd); err != nil {
		return err
	}
	if err := g.SetKeybinding("log", gocui.KeyArrowDown, gocui.ModNone, ui.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("log", gocui.KeyArrowUp, gocui.ModNone, ui.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("code", gocui.KeyArrowDown, gocui.ModNone, ui.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("code", gocui.KeyArrowUp, gocui.ModNone, ui.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("code", gocui.KeySpace, gocui.ModNone, ui.executeStep); err != nil {
		return err
	}

	return nil
}



func (ui *Ui) Layout(g *gocui.Gui) error {

	const codeView_width    = 41 // with frames - no resize
	const stackView_width   =  8 // with frames - no resize
	const statusView_height =  4 // with frames - no resize
	const memoryView_width  = 74 // with frames - no resize
	const cmdView_height    =  3 // with frames, resizeable

	const logView_height    = 10 // with frames
	const memoryView_height = 18 // with frames

	maxX, maxY := g.Size()

	// sample assembly code
	if v, err := g.SetView("code", maxX-codeView_width, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true
		//v.SelBgColor = gocui.ColorGreen
		if _, err := g.SetCurrentView("code"); err != nil {
			return err
		}

		ui.updateCodeView(g)
	}

	// sample stack window
	x1_stack := maxX - (codeView_width + stackView_width + 1)
	x2_stack := maxX - (codeView_width + 1)

	if v, err := g.SetView("stack", x1_stack, 0, x2_stack, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = true
		v.Autoscroll = true
		//v.SelBgColor = gocui.ColorGreen
		ui.updateStackView(g)
	}

	// sample status window
	if v, err := g.SetView("status", 0, 0, x1_stack - 1, statusView_height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = false

		ui.updateStatusView(g)
	}

	// sample memory view window
	if v, err := g.SetView("memory", 0, maxY-memoryView_height, memoryView_width, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = false
		v.Title = "Memory"
		g.SetViewOnBottom("memory")

		ui.updateMemoryView(g)
	}

	// sample memory view window
	// cmdView_width euqals to memoryView_width
	y1_cmd := maxY - (memoryView_height + cmdView_height)
	y2_cmd := y1_cmd + cmdView_height - 1
	if v, err := g.SetView("cmd", 0, y1_cmd, x1_stack - 1, y2_cmd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true
		v.Title = "Command"

	}

	/*
	// sample log view window
	sizes  := (y1_cmd - statusView_height - 3) >> 1
	y1_log := statusView_height + 1
	y2_log := y1_log + sizes
	if v, err := g.SetView("out", 0, y1_log, x1_stack-1, y2_log); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true
		v.Title = "Output"
		v.Editor = gocui.EditorFunc(outViewEditor)
		outView = v
	}
	*/
	// log view window

	if v, err := g.SetView("log", 0, statusView_height + 1, x1_stack - 1, y1_cmd - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = true
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true
		v.Title = "Log"
		ui.logView = v
	}

	if ui.isPageMapVisible {
		ui.showPageMap(g)
	}

	return nil
}

