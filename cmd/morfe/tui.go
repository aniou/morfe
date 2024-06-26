package main

import (
	"container/list"
        "fmt"
        "log"
        "sort"
        "strconv"
        "strings"
        _ "time"
        "github.com/aniou/morfe/lib/mylog"
        "github.com/aniou/morfe/emulator/emu"
        "github.com/aniou/morfe/emulator/platform"

        "github.com/awesome-gocui/gocui"
)

const ansi_red   = "\x1b[0;31m"
const ansi_green = "\x1b[0;32m"
const ansi_reset = "\x1b[0m"

const history_size = 30

type Ui struct {
        viewArr     []string           // names of views to cycle
        active      int                // index in viewArr
        lastView    string             // last view.Name when cmd is called

        memPosition uint32             // displayed memory region
        memPrevious *[128]byte         // previous memory region

        logView     *gocui.View        // shortcut for current UI
        ch          chan string
        cpu         emu.Processor      // for what CPU is that GUI?     
        preg        map[string]uint32  // previous values of registers
        watch       map[int][4]byte    // list of addresses and previous values

	h_queue	    *list.List         // cmd history
	h_element   *list.Element      // actual position in history
}

func NewTUI(ch chan string, cpu emu.Processor) *Ui {
	ui      := Ui{ch: ch, cpu: cpu}
        return &ui
}

func (ui *Ui) Init(g *gocui.Gui) {
        ui.viewArr = []string{"code", "cmd"}
        ui.active = 1 // "cmd"
        //ui.logger = logger
        ui.watch = make(map[int][4]byte)
        ui.preg  = ui.cpu.GetRegisters()
        ui.memPrevious = &[128]byte{}
	ui.h_queue = list.New()
	ui.h_queue.PushFront("")
	ui.h_element = ui.h_queue.Front()

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

// view support functions
func (ui *Ui) test_step(g *gocui.Gui, v *gocui.View) error {
        ui.ch<-"step"
        _ = <-ui.ch     // wait for ack from goroutine
        ui.updateStatusView(g)
        ui.updateCodeView(g)
        ui.updateWatchView(g)
        ui.updateMemoryView(g)
        ui.updateStackView(g)
        return nil
}

// XXX - kind of PoC, there may be a generic function for printing addresses at specific
//                    positions in View

func makeHex32(val uint32, separator string) string {
        var high string
        var low  string

        if (val >> 16) == 0 {
                high = "    "
                low  = fmt.Sprintf("%4x", val & 0x0000_FFFF)
        } else {
                high = fmt.Sprintf("%4x",  val >> 16)
                low  = fmt.Sprintf("%04x", val & 0x0000_FFFF)
        }

        hex := fmt.Sprintf("%s%s%s",
                                high,
                                separator,
                                low)
        return hex
}

func hex2uint32(hexStr string) (uint32, error) {
        // remove 0x suffix, $ and : characters
        cleaned := strings.Replace(hexStr,  "0x", "",  1)
        cleaned  = strings.Replace(cleaned,  "$", "",  1)
        cleaned  = strings.Replace(cleaned,  ":", "", -1)
        cleaned  = strings.Replace(cleaned,  "_", "", -1)

        result, err := strconv.ParseUint(cleaned, 16, 32)
        if err != nil {
                err = fmt.Errorf("hex2uint32: %w", err)
        }
        return uint32(result), err
}


func makeSafeAscii(val byte) string {
        if val >= 33 && val < 127 {
                return string(val)
        } else {
                return "."
        }
}

func (ui *Ui) cmd_enable(g *gocui.Gui, tokens []string) {
	switch tokens[1] {
	case "cpu0":
		ui.ch<-"enable_cpu0"
	case "cpu1":
		ui.ch<-"enable_cpu1"
        default:
                fmt.Fprintf(ui.logView, "enable: unknown parameter '%s'\n", tokens[1:])
	}
}

func (ui *Ui) cmd_disable(g *gocui.Gui, tokens []string) {
	switch tokens[1] {
	case "cpu0":
		ui.ch<-"disable_cpu0"
	case "cpu1":
		ui.ch<-"disable_cpu1"
        default:
                fmt.Fprintf(ui.logView, "disable: unknown parameter '%s'\n", tokens[1:])
	}
}

func (ui *Ui) cmd_watch(g *gocui.Gui, tokens []string) {
        switch tokens[1] {
        case "add":
                if addr, err := hex2uint32(tokens[2]); err == nil {
                        ui.addWatch(int(addr))
                        ui.updateWatchView(g)
                } else {
                        fmt.Fprintf(ui.logView, "watch add: error: %s\n", err)
                }
        case "del":
                if addr, err := hex2uint32(tokens[2]); err == nil {
                        ui.delWatch(int(addr))
                        ui.updateWatchView(g)
                } else {
                        fmt.Fprintf(ui.logView, "watch add: error: %s\n", err)
                }
        default:
                fmt.Fprintf(ui.logView, "watch: unknown parameter '%s'\n", tokens[1:])
        }
}

func (ui *Ui) cmd_set(g *gocui.Gui, tokens []string) {
        var err error
        var val uint32

        switch tokens[1] {
        case "mem":
                if len(tokens) < 3 {
                        err = fmt.Errorf("not enough parameters")
                        break
                }
                if ui.memPosition, err = hex2uint32(tokens[2]); err == nil {
                        mem := [128]byte{}
                        for x:=uint32(0) ; x<128; x=x+1 {
                                mem[x] = ui.cpu.Read_8(ui.memPosition + x)
                        }
                        ui.memPrevious = &mem
                        ui.updateMemoryView(g)
                }
        case "reg": // set reg <name> <value>
                if len(tokens) < 4 {
                        err = fmt.Errorf("not enough parameters")
                        break
                }
                if val, err = hex2uint32(tokens[3]); err == nil {
                        err = ui.cpu.SetRegister(strings.ToUpper(tokens[2]), val)
                        ui.updateStatusView(g)
                }
        case "pc": // set pc <value>
                if len(tokens) < 3 {
                        err = fmt.Errorf("not enough parameters")
                        break
                }
                if val, err = hex2uint32(tokens[2]); err == nil {
                        ui.cpu.SetPC(val)
                        ui.updateStatusView(g)
                        ui.updateCodeView(g)
                }
        default:
                err = fmt.Errorf("unknown parameter %v", tokens[2])
        }

        if err != nil {
                fmt.Fprintf(ui.logView, "ERROR cmd_set %s\n", err)

        }
}

func (ui *Ui) cmd_load(g *gocui.Gui, tokens []string) {
        var err error = nil

        switch tokens[1] {
        case "hex":
                platform.LoadHex(ui.cpu, tokens[2])     // TODO - convert to error support
                ui.updateStatusView(g)
                ui.updateCodeView(g)
                ui.updateWatchView(g)
                ui.updateMemoryView(g)
                ui.updateStackView(g)
        default:
                err = fmt.Errorf("unknown parameter %v", tokens[1])
        }

        if err != nil {
                fmt.Fprintf(ui.logView, "ERROR cmd_load: %s\n", err)
        }
}

// executes multiple commands, separated by ';'
func (ui *Ui) executeCommands(g *gocui.Gui, v *gocui.View) error {
        line      := strings.TrimSpace(v.Buffer())

	if line != "history" && line != "hi" {
		ui.h_element = ui.h_queue.Front()
		ui.h_queue.InsertAfter(line, ui.h_element)	// remember history
	}

        commands  := strings.Split(line, ";")
	for _, c  := range commands {
		if err := ui.executeSingleCommand(g, v, c); err != nil {
			return err
		}

	}

        v.Clear()
        v.SetCursor(0, 0)

        return nil
}

// executes single command
func (ui *Ui) executeSingleCommand(g *gocui.Gui, v *gocui.View, command string) error {
        line   := strings.TrimSpace(command)
	if len(line) == 0 {
		return nil
	}
        tokens := strings.Split(line, " ")
        switch tokens[0] {
        case "di", "disable":
                ui.cmd_disable(g, tokens)
        case "en", "enable":
                ui.cmd_enable(g, tokens)
        case "wa", "watch":
                ui.cmd_watch(g, tokens)
        case "se", "set":
                ui.cmd_set(g, tokens)
        case "reset":
                ui.cpu.Reset()
                ui.cpu.ResetCycles()
                v, err := g.View("code")
                if err != nil {
                        return err
                } else {
                        v.Clear()
                }
                ui.updateStatusView(g)
                ui.updateCodeView(g)
                ui.updateWatchView(g)
        case "lo", "load":
                ui.cmd_load(g, tokens)
	case "hi", "history":
		count := 0
		for e := ui.h_queue.Back(); e != nil; e = e.Prev() {
			fmt.Fprintf(ui.logView, "%d ", count)
			fmt.Fprintf(ui.logView, e.Value.(string))
			fmt.Fprintf(ui.logView, "\n")
			count=count+1
		}
        case "quit":
                return ui.quit(g, v)
        /*
        case "run":
                ui.runCPU(g, v)
        case "peek", "peek8", "peek16", "peek24":
                ui.peek(g, tokens)
        */
        default:
                fmt.Fprintf(ui.logView, "executeCommand: unknown command: %s\n", command)
        }

	return nil
}



func (ui *Ui) updateStatusView(g *gocui.Gui) error {
        v, err := g.View("status")
        if err != nil {
                //fmt.Fprintf(ui.logView, "%s\n", err)  - XXX - doesn't work yet
                return err
        }

        v.Clear()

        reg := ui.cpu.GetRegisters()

        order := [][]string {   {"D0", "PC"},
                                {"D1", "PPC"},
                                {"D2", "SR"},
                                {"D3", "IR"},
                                {"D4", "SP"},
                                {"D5", "USP"},
                                {"D6", "ISP"},
                                {"D7", "MSP"},
                                {"",   "SFC"},
                                {"A0", "DFC"},
                                {"A1", "VBR"},
                                {"A2", "CACR"},
                                {"A3", "CAAR"},
                                {"A4", ""},
                                {"A5", ""},
                                {"A6", ""},
                                {"A7", ""},
                            }

        for _, val := range order {
                if val[0] == "" {
                        fmt.Fprintf(v, "              ")
                } else {
                        hex       := makeHex32(reg[val[0]], " ")
                        different := ui.preg[val[0]] != reg[val[0]]

                        fmt.Fprintf(v, "%2s ", val[0])
                        printColored(v, different, hex)
                        fmt.Fprintf(v, "  ")
                }


                if val[1] == "" {
                        fmt.Fprintf(v, "\n")
                } else {
                        hex       := makeHex32(reg[val[1]], " ")
                        different := ui.preg[val[1]] != reg[val[1]]
                                   
                        fmt.Fprintf(v, "%4s ", val[1])
                        printColored(v, different, hex)
                        fmt.Fprintf(v, "\n")
                }

        }

        fmt.Fprintf(v, "\n%s", ui.cpu.StatusString())

        ui.preg = reg
        return nil
}

func (ui *Ui) updateLogView(g *gocui.Gui) error {
        v, err := g.View("log")
        if err != nil {
                return err
        }

        v.Clear()
        fmt.Fprintf(v, "Preliminary debug interface\npress keys:\n")
        fmt.Fprintf(v, "F5     to execute single step\n")
        fmt.Fprintf(v, "CTRL+Q or 'quit' cmd to exit debugger\n")
        fmt.Fprintf(v, "\n")
        fmt.Fprintf(v, "Commands (can be separated by ';' in single line):\n")
        fmt.Fprintf(v, "watch {add|del} <value>    - manage watch list vals\n")
        fmt.Fprintf(v, "set reg <regname> <value>  - set value of register\n")
        fmt.Fprintf(v, "set pc <addr>              - set Program Counter of cpu\n")
        fmt.Fprintf(v, "reset                      - perform cpu reset \n")
        fmt.Fprintf(v, "load hex <filename>        - load IntelHex format into memory\n")
        fmt.Fprintf(v, "set mem <addr>             - set memory view address\n")
        fmt.Fprintf(v, "history                    - display history \n")
        fmt.Fprintf(v, "enable|disable {cpu0|cpu1} - enables/disables cpu\n")


        return nil
}

func (ui *Ui) updateCodeView(g *gocui.Gui) error {
        v, err := g.View("code")
        if err != nil {
                return err
        }
        //v.Clear()

        line   := ui.cpu.Dissasm()
        cycles := ui.cpu.GetCycles()
        if cycles > 0 {
                fmt.Fprintf(v, " : %d\n", ui.cpu.GetCycles())
        } else {
                fmt.Fprintf(v, "\n")
        }
        fmt.Fprintf(v, "%-58s", line)

        return nil
}

func (ui *Ui) updateMemoryView(g *gocui.Gui) error {
        v, err := g.View("memory")
        if err != nil {
                return err
        }
        v.Clear()

        var a,x uint32
        var val byte
        var mem [128]byte

        for a=0; a<0x80; a=a+16 {
                start := ui.memPosition + a

                // addr
                high := fmt.Sprintf("%04x", start >> 16)
                low  := fmt.Sprintf("%04x", start & 0x0000_FFFF)
                fmt.Fprintf(v, "%s %s: ", high, low)

                // hex format
                for x = 0; x<16;  x=x+1 {
                        val      = ui.cpu.Read_8(start+x)
                        mem[a+x] = val
                        changed := ui.memPrevious[a+x] != val
                        printColored(v, changed, fmt.Sprintf("%02x ", val))
                        if x == 7 {
                                fmt.Fprintf(v, "| ")
                        }
                }

                // ascii format
                for x = 0; x<16;  x=x+1 {
                        val      = ui.cpu.Read_8(start+x)
                        changed := ui.memPrevious[a+x] != val
                        printColored(v, changed, makeSafeAscii(val))
                        if x == 7 {
                                fmt.Fprintf(v, " ")
                        }
                }
                fmt.Fprintf(v, "\n")
        }
        ui.memPrevious = &mem
        return nil
}

func (ui *Ui) updateStackView(g *gocui.Gui) error {
        v, err := g.View("stack")
        if err != nil {
                return err
        }
        v.Clear()

        reg := ui.cpu.GetRegisters()
	sp  := reg["SP"]

	switch (reg["SR"] >> 12) & 3 {
	case 0:
	case 1:
		v.Title = "Stack (USP)"	
	case 2:
		v.Title = "Stack (ISP)"	
	case 3:
		v.Title = "Stack (MSP)"
	}

	_, sizeY := v.Size() 

	offset := uint32((sizeY / 2) * 4)	// half of window and * 4 bytes

	var start uint32
	if (sp <= offset) {
		start = 0
	} else {
		start = sp - offset
	}

	for x:=0; x<sizeY; x=x+1 {
		addr := start+uint32(x*4)
		if addr == sp {
			fmt.Fprintf(v, ansi_green)
		}

		fmt.Fprintf(v, makeHex32(addr, " "))
		fmt.Fprintf(v, ": ")
		fmt.Fprintf(v, "%02x", ui.cpu.Read_8(addr  ))
		fmt.Fprintf(v, "%02x", ui.cpu.Read_8(addr+1))
		fmt.Fprintf(v, " ")
		fmt.Fprintf(v, "%02x", ui.cpu.Read_8(addr+2))
		fmt.Fprintf(v, "%02x", ui.cpu.Read_8(addr+3))

		if addr == sp {
			fmt.Fprintf(v, ansi_reset)
		} 

		fmt.Fprintf(v, "\n")
	}

	return nil
}





func printColored(v *gocui.View, highlighted bool, s string) {
        if highlighted {
                fmt.Fprintf(v, ansi_red)
                fmt.Fprintf(v, s)
                fmt.Fprintf(v, ansi_reset)
        } else {
                fmt.Fprintf(v, s)
        }
}


func (ui *Ui) addWatch(addr int) {
        // only distinct values
        if _, exists := ui.watch[addr]; exists {
                return
        }

        // initialize values
        vals := [4]byte{}
        for i:=0; i<4; i++ {
                b := ui.cpu.Read_8(uint32(addr + i))
                vals[i] = b
        }
        ui.watch[addr] = vals
}

func (ui *Ui) delWatch(addr int) {
        // only distinct values
        if _, exists := ui.watch[addr]; exists {
                delete(ui.watch, addr)
        }
}

func (ui *Ui) updateWatchView(g *gocui.Gui) error {
        v, err := g.View("watch")
        if err != nil {
                return err
        }
        v.Clear()

        // create sorted list of addresses...
        keys := []int{}
        for k := range ui.watch {
            keys = append(keys, int(k))
        }
        sort.Ints(keys)
        
        // ...and use them to iterate over ui.watch
        for _, key := range keys {
                fmt.Fprintf(v, "%s: ", makeHex32(uint32(key), " "))
                tmparr := [4]byte{}

                // print HEX represenation of four consecutive bytes
                for i:=0; i < 4; i++ {
                        b := ui.cpu.Read_8(uint32(key+i))
                        printColored(v, 
                                     b!=ui.watch[key][i], 
                                     fmt.Sprintf("%02x ", b))
                        tmparr[i] = b
                }

                // print ASCII represenation of byte
                fmt.Fprintf(v, " ")
                for i:=0; i < 4; i++ {
                        b := ui.cpu.Read_8(uint32(key+i))
                        printColored(v,
                                     b!=ui.watch[key][i],
                                     makeSafeAscii(b))
                }
                fmt.Fprintf(v, "\n")

                // update old values
                ui.watch[key] = tmparr
        }

        return nil
}

func (ui *Ui) historyPrev(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	if ui.h_element == nil {
		return nil
	}

	if e := ui.h_element.Next(); e != nil {
		ui.h_element = e
	}

	// https://tour.golang.org/methods/15
	t := ui.h_element.Value.(string)
	fmt.Fprintf(v, t)
	return nil
}

func (ui *Ui) historyNext(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	if ui.h_element == nil {
		return nil
	}

	if e := ui.h_element.Prev(); e != nil {
		ui.h_element = e
	}

	// https://tour.golang.org/methods/15
	t := ui.h_element.Value.(string)
	fmt.Fprintf(v, t)
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
        if err := g.SetKeybinding("cmd", gocui.KeyEnter, gocui.ModNone, ui.executeCommands); err != nil {
                return err
        }
        if err := g.SetKeybinding("cmd", gocui.KeyArrowUp, gocui.ModNone, ui.historyPrev); err != nil {
                return err
        }
        if err := g.SetKeybinding("cmd", gocui.KeyArrowDown, gocui.ModNone, ui.historyNext); err != nil {
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

        maxX, maxY := g.Size()

        const v_stat_x1    = 0
        const v_stat_y1    = 0
        const v_stat_x2    = 29
        const v_stat_y2    = 21

        const v_watch_x1   = 0
        const v_watch_y1   = v_stat_y2 + 1
        const v_watch_x2   = v_stat_x2
              v_watch_y2  := maxY - 1

        const v_stack_x1   = v_stat_x2 + 1
        const v_stack_y1   = 0
        const v_stack_x2   = v_stat_x2 + 22
              v_stack_y2  := maxY - 1

        const v_code_x1    = v_stack_x2 + 1
        const v_code_y1    = 0
              v_code_x2   := maxX - 1
              v_code_y2   := maxY / 3
              
        const v_memory_x1  = v_stack_x2 + 1
              v_memory_y1 := v_code_y2  + 1
              v_memory_x2 := maxX - 1
              v_memory_y2 := v_memory_y1  + 9

        const v_log_x1     = v_stack_x2 + 1
              v_log_y1    := v_memory_y2  + 1
              v_log_x2    := maxX - 1
              v_log_y2    := maxY - 4
              
        const v_cmd_x1     = v_stack_x2 + 1
              v_cmd_y1    := v_log_y2   + 1
              v_cmd_x2    := maxX - 1
              v_cmd_y2    := maxY - 1

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

                ui.updateStatusView(g)
        }

        // watch mem view window
        if v, err := g.SetView("watch", v_watch_x1, v_watch_y1, v_watch_x2, v_watch_y2, 0); err != nil {
                if err != gocui.ErrUnknownView {
                        return err
                }

                v.Editable = false
                v.Wrap = false
                v.Frame = true
                v.Highlight = false
                v.Autoscroll = true
                v.Title = "Watch"

                ui.updateWatchView(g)
        }

        if v, err := g.SetView("stack", v_stack_x1, v_stack_y1, v_stack_x2, v_stack_y2, 0); err != nil {
                if err != gocui.ErrUnknownView {
                        return err
                }

                v.Editable = false
                v.Wrap = false
                v.Frame = true
                v.Highlight = false
                v.Autoscroll = false

                ui.updateStackView(g)
        }

        if v, err := g.SetView("code", v_code_x1, v_code_y1, v_code_x2, v_code_y2, 0); err != nil {
                if err != gocui.ErrUnknownView {
                        return err
                }

                v.Editable = false
                v.Wrap = false
                v.Frame = true
                v.Highlight = false
                v.Autoscroll = true
		v.Title = "code"

                ui.updateCodeView(g)
        }

        if v, err := g.SetView("memory", v_memory_x1, v_memory_y1, v_memory_x2, v_memory_y2, 0); err != nil {
                if err != gocui.ErrUnknownView {
                        return err
                }

                v.Editable = false
                v.Wrap = false
                v.Frame = true
                v.Highlight = false
                v.Autoscroll = false
                v.Title = "Memory"

                ui.updateMemoryView(g)
        }

        if v, err := g.SetView("log", v_log_x1, v_log_y1, v_log_x2, v_log_y2, 0); err != nil {
                if err != gocui.ErrUnknownView {
                        return err
                }

                v.Editable = false
                v.Wrap = false
                v.Frame = true
                v.Highlight = false
                v.Autoscroll = true
                v.Title = "Log"

                ui.logView = v
                mylog.Logger.LogOutput = v      // not nice
                ui.updateLogView(g)
        }

        if v, err := g.SetView("cmd", v_cmd_x1, v_cmd_y1, v_cmd_x2, v_cmd_y2, 0); err != nil {
                if err != gocui.ErrUnknownView {
                        return err
                }

                v.Editable = true
                v.Wrap = false
                v.Frame = true
                v.Highlight = false
                v.Autoscroll = false
                v.Title = "Cmd"

                if _, err := g.SetCurrentView("cmd"); err != nil {
                        return err
                }

                v.Clear()

                //ui.updateStatusView(g)
        }

        return nil
}

func mainTUI(ch chan string, cpu emu.Processor) {
        g, err := gocui.NewGui(gocui.Output256, false)
        if err != nil {
                log.Panicln(err)
        }

        ui     := NewTUI(ch, cpu)
        ui.Init(g)

        //ui.addWatch(0x10_0000)
        //ui.addWatch(0xaf_a000)

        ui.Run(g)
        g.Close()
        fmt.Print("sending quitting...")
        ch<-"run"
        fmt.Print("quitting...")
}
