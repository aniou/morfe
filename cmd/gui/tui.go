package main

import (
	"fmt"
	"log"
	_ "strconv"
	_ "strings"
	_ "time"

	_ "github.com/aniou/go65c816/lib/mylog"
	"github.com/awesome-gocui/gocui"
)

type Ui struct {
	viewArr []string // names of views to cycle
	active  int      // index in viewArr

	memPosition uint32 // displayed memory region

	lastView         string // last view.Name when cmd is called
	cpuSpeed         uint64

	logView *gocui.View        // shortcut for current UI

	ch	chan string
}

func NewTUI(ch chan string) *Ui {
	ui := Ui{ch: ch}
	return &ui
}

func (ui *Ui) Init(g *gocui.Gui) {
	ui.viewArr = []string{"code", "cmd"}
	ui.active = 1 // "cmd"
	//ui.logger = logger

	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.Highlight = true

	g.SetManagerFunc(ui.Layout)
	if err := ui.keybindings(g); err != nil {
		log.Panicln(err)
	}

	//go ui.Logger(g)
	//logger.Log("TUI init complete")
}

func (ui *Ui) Run(g *gocui.Gui) {
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
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
	return gocui.ErrQuit
}

func (ui *Ui) test_step(g *gocui.Gui, v *gocui.View) error {
	ui.ch<-"step"
	return nil
}


func (ui *Ui) keybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("", gocui.KeyF5, gocui.ModNone, ui.test_step); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, ui.quit); err != nil {
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

	return nil
}

func (ui *Ui) Layout(g *gocui.Gui) error {

	const codeView_width = 41   // with frames - no resize
	const stackView_width = 8   // with frames - no resize
	const memoryView_width = 74 // with frames - no resize
	const cmdView_height = 3    // with frames, resizeable

	const logView_height = 10    // with frames
	const memoryView_height = 18 // with frames

	maxX, maxY := g.Size()

	const v_stat_x1   = 0
	const v_stat_y1   = 0
	const v_stat_x2   = 30
	const v_stat_y2   = 21

	const v_quick_x1  = 0
	const v_quick_y1  = v_stat_y2 + 1
	const v_quick_x2  = 30
	      v_quick_y2 := maxY - 1

	const v_stack_x1  = v_stat_x2 + 1
	const v_stack_y1  = 0
        const v_stack_x2  = v_stat_x2 + 20
              v_stack_y2 := maxY - 1

	const v_code_x1   = v_stack_x2 + 1
	const v_code_y1   = 0
	      v_code_x2  := maxX - 1
	      v_code_y2  := maxY / 3
	      
	const v_dump_x1   = v_stack_x2 + 1
	      v_dump_y1  := v_code_y2  + 1
	      v_dump_x2  := maxX - 1
	      v_dump_y2  := v_dump_y1  + 10

	const v_log_x1    = v_stack_x2 + 1
	      v_log_y1   := v_dump_y2  + 1
	      v_log_x2   := maxX - 1
	      v_log_y2   := maxY - 4
              
	const v_cmd_x1    = v_stack_x2 + 1
	      v_cmd_y1   := v_log_y2   + 1
	      v_cmd_x2   := maxX - 1
	      v_cmd_y2   := maxY - 1

	// cpu status window
	if v, err := g.SetView("status", v_stat_x1, v_stat_y1, v_stat_x2, v_stat_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = false
		v.Title = "Status"

		//ui.updateStatusView(g)
	}

	// quick mem view window
	if v, err := g.SetView("quick", v_quick_x1, v_quick_y1, v_quick_x2, v_quick_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true

		//ui.updateStatusView(g)
	}

	// quick mem view window
	if v, err := g.SetView("stack", v_stack_x1, v_stack_y1, v_stack_x2, v_stack_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true

		//ui.updateStatusView(g)
	}

	// quick mem view window
	if v, err := g.SetView("code", v_code_x1, v_code_y1, v_code_x2, v_code_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true

		//ui.updateStatusView(g)
	}

	// quick mem view window
	if v, err := g.SetView("dump", v_dump_x1, v_dump_y1, v_dump_x2, v_dump_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true

		//ui.updateStatusView(g)
	}

	// quick mem view window
	if v, err := g.SetView("log", v_log_x1, v_log_y1, v_log_x2, v_log_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true

		//ui.updateStatusView(g)
	}

	// quick mem view window
	if v, err := g.SetView("cmd", v_cmd_x1, v_cmd_y1, v_cmd_x2, v_cmd_y2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = false
		v.Frame = true
		v.Highlight = false
		v.Autoscroll = true

		if _, err := g.SetCurrentView("cmd"); err != nil {
			return err
		}

        	v.Clear()
		fmt.Fprintf(v, "> ")

		//ui.updateStatusView(g)
	}

	return nil
}

func mainTUI(ch chan string) {
	g, err := gocui.NewGui(gocui.Output256, false)
	if err != nil {
		log.Panicln(err)
	}

	ui     := NewTUI(ch)
	ui.Init(g)
	ui.Run(g)
	g.Close()
	fmt.Print("sending quitting...")
	ch<-"run"
	fmt.Print("quitting...")
}
