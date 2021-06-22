
package platform

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
	"github.com/marcinbor85/gohex"
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator"
)

// xxx - duplicate in TUI, go to lib
func hex2uint24(hexStr string) (uint32, error) {
        // remove 0x suffix, $ and : characters
        cleaned := strings.Replace(hexStr, "0x", "", 1)
        cleaned =  strings.Replace(cleaned, "$", "", 1)
        cleaned =  strings.Replace(cleaned, ":", "", -1)

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

// load file= and file....= into memory of desired CPU
func (p *Platform) loadFiles(cfg *ini.File, c emu.Processor) {
	cpu_section := c.GetName()

	// at first read file=
	if cfg.Section(cpu_section).HasKey("file") {
		hexfile := cfg.Section(cpu_section).Key("file").String()
		p.LoadHex(c, hexfile)
	}
	// then file0= to file99=
	for i := 0; i<100; i += 1 {
		keyname := fmt.Sprintf("file%d", i)
		if cfg.Section(cpu_section).HasKey(keyname) {
			hexfile := cfg.Section(cpu_section).Key(keyname).String()
			p.LoadHex(c, hexfile)
		}
	}
}

func (p *Platform) LoadConfig(filename string) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: false,
	}, filename)
	if err != nil {
        	log.Fatalf("Failed to load from %s - error: %s\n", filename, err)
        }

	p.loadFiles(cfg, p.CPU0)
	p.loadFiles(cfg, p.CPU1)

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

func (p *Platform) LoadHex(cpu emu.Processor, filename string) {
	path := filepath.Join(filename)
	file, err := os.Open(path)
	if err != nil {
		mylog.Logger.Log(fmt.Sprintf("LoadHex failed: %s", err))
		return
	}
	defer file.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		panic(err)
	}

	mylog.Logger.Log(fmt.Sprintf("LoadHex for cpu %s - loading file %s", cpu.GetName, path))
	for idx, segment := range mem.GetDataSegments() {
		mylog.Logger.Log(fmt.Sprintf("%d addr %06x length %6x (%d)",
					idx, segment.Address, len(segment.Data), len(segment.Data)))
                for i := range segment.Data {
                        cpu.Write_8(segment.Address + uint32(i), segment.Data[i])
                }
	}
	mylog.Logger.Log(fmt.Sprintf("LoadHex done"))
}
