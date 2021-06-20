// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	//"os"

	"gopkg.in/ini.v1"
        "github.com/aniou/go65c816/lib/mylog"
        "github.com/aniou/go65c816/emulator/bus"
        "github.com/aniou/go65c816/emulator"
        "github.com/aniou/go65c816/emulator/cpu65c816"
        "github.com/aniou/go65c816/emulator/cpu68xxx"
        "github.com/aniou/go65c816/emulator/memory"
        "github.com/aniou/go65c816/emulator/vicky"
        "github.com/aniou/go65c816/emulator/gabe"
)

type Platform struct {
        CPU0    emu.Processor           // on-board one (65c816 by default)
        CPU1    emu.Processor           // add-on
        GPU     *vicky.Vicky
        GABE    *gabe.Gabe
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

// something like GenX
func (p *Platform) InitGUI() {
	bus0,_	   := bus.New()
	bus1,_     := bus.New()
        ram0, _    := memory.New(0x400000, 0x000000)
        p.GPU, _    = vicky.New()
        p.GABE, _   = gabe.New()

        bus0.Attach(ram0,   "ram0",      0x00_0000, 0x3F_FFFF)  // xxx - 1: ram.offset, ram.size 2: get rid that?
        bus0.Attach(p.GPU,  "vicky",     0xAF_0000, 0xEF_FFFF)
        bus0.Attach(p.GABE, "gabe",      0xAF_1000, 0xAF_13FF)  // probably should be splitted
        bus0.Attach(p.GABE, "gabe-math", 0x00_0100, 0x00_012F)  // XXX error GABE coop is 0x2c bytes but we need mult of 16

        bus1.Attach(p.GPU,  "vicky",     0xAF_0000, 0xEF_FFFF)

        p.CPU0     = cpu65c816.New(bus0)
        p.CPU1     = cpu68xxx.New(bus1)   // 20Mhz - not used yet
        
        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()                                      // XXX - move it to main binary?

        p.CPU1.Reset()

        mylog.Logger.Log("platform initialized")
}


// xxx - duplicate in TUI, go to lib
func hex2uint24(hexStr string) (uint32, error) {
        // remove 0x suffix, $ and : characters
        cleaned := strings.Replace(hexStr, "0x", "", 1)
        cleaned = strings.Replace(cleaned, "$", "", 1)
        cleaned = strings.Replace(cleaned, ":", "", -1)

        result, err := strconv.ParseUint(cleaned, 16, 32)
        return uint32(result) & 0x00ffffff, err
}

func hex2uint16(hexStr string) (uint16, error) {
        // remove 0x suffix, $ and : characters
        cleaned := strings.Replace(hexStr, "0x", "", 1)
        cleaned = strings.Replace(cleaned, "$", "", 1)
        cleaned = strings.Replace(cleaned, ":", "", -1)

        result, err := strconv.ParseUint(cleaned, 16, 16)
        return uint16(result), err
}


func (p *Platform) LoadConfig(filename string) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: false,
	}, filename)
	if err != nil {
        	log.Fatalf("Failed to load from %s - error: %s\n", filename, err)
        }

	// hex load -------------------------------------------------------------
	// load hex files (at this time - only hex files)
	number := len(cfg.Section("load").Keys())
	if cfg.Section("load").HasKey("file") {
		number -= 1
		hexfile := cfg.Section("load").Key("file").String()
		p.LoadHex(p.CPU0, hexfile)
		//fmt.Printf("key file found, file %s\n", hexfile)
	}

	for i := 0; i<1000 && number>0; i += 1 {
		keyname := fmt.Sprintf("file%d", i)
		if cfg.Section("load").HasKey(keyname) {
			hexfile := cfg.Section("load").Key(keyname).String()
			p.LoadHex(p.CPU0, hexfile)
			//fmt.Printf("key %s found\n", keyname)
			number -= 1
		}
	}

	// hex load -------------------------------------------------------------
	// CPU setting XXX - change to cpu0
	if cfg.Section("cpu0").HasKey("start") {
		hex_addr := cfg.Section("cpu0").Key("start").String()
		addr, _  := hex2uint24(hex_addr)
		fmt.Printf("start addr set: %06X\n", addr)
		//g.p.CPU.PC = uint16(addr & 0x0000FFFF)
		//g.p.CPU.RK = uint8(addr >> 16)
		p.CPU0.SetPC(uint32(addr))
	}

	/*
	if cfg.Section("cpu0").HasKey("wdm_mode") {
		wdm_mode := cfg.Section("cpu0").Key("wdm_mode").String()
		switch wdm_mode {
		case "debug":
			debug.cpu = true	// XXX - bad behaviour, globals!
		default:
			debug.cpu = false
		}
	}
	*/
}
