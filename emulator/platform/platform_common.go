
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
	"github.com/aniou/morfe/lib/mylog"
	"github.com/aniou/morfe/emulator"
)

type Config struct {
	Mode	string
}

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
		if err := LoadHex(c, hexfile); err != nil {
			log.Panicln(err)
		}
	}
	// then file0= to file99=
	for i := 0; i<100; i += 1 {
		keyname := fmt.Sprintf("file%d", i)
		if cfg.Section(cpu_section).HasKey(keyname) {
			hexfile := cfg.Section(cpu_section).Key(keyname).String()
			if err := LoadHex(c, hexfile); err != nil {
				log.Panicln(err)
			}
		}
	}
}

// now, only PC, by keyword 'start' is supported - registers in future
func (p *Platform) setRegisters(cfg *ini.File, c emu.Processor) {
	cpu_section := c.GetName()

	if cfg.Section(cpu_section).HasKey("start") {
		hex_addr := cfg.Section(cpu_section).Key("start").String()
		addr, _  := hex2uint24(hex_addr)
		fmt.Printf("start addr set for cpu %s: %06X\n", cpu_section, addr)
		c.SetPC(uint32(addr))
	}
	if cfg.Section(cpu_section).HasKey("enable") {
		state, _ := cfg.Section(cpu_section).Key("enable").Bool()
		c.Enable(state)
	}
}

// XXX - make it more common and move loadFiles to gui?
func (p *Platform) LoadCpuConfig(filename string) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: false,
	}, filename)
	if err != nil {
        	log.Fatalf("Failed to load from %s - error: %s\n", filename, err)
        }

	p.loadFiles(cfg, p.CPU0)
	p.setRegisters(cfg, p.CPU0)

	p.loadFiles(cfg, p.CPU1)
	p.setRegisters(cfg, p.CPU1)

}

func (p *Platform) LoadPlatformConfig(filename string) (*Config, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: false,
	}, filename)
	if err != nil {
        	log.Fatalf("Failed to load from %s - error: %s\n", filename, err)
        }

	pcfg      := Config{}
	pcfg.Mode  = cfg.Section("platform").Key("mode").In("fmx-like", []string{"fmx-like", "frankenmode", "genx-like"})

	// set DIP-switch config in emu
	for i := 1; i<7; i += 1 {
		keyname := fmt.Sprintf("DIP%d", i)
		if cfg.Section("platform").HasKey(keyname) {
			emu.DIP[i], _ = cfg.Section("platform").Key(keyname).Bool()
		}
	}
	return &pcfg, nil

}

// XXX - move error support and displaying into upper layer
func LoadHex(cpu emu.Processor, filename string) error {
	path := filepath.Join(filename)
	file, err := os.Open(path)
	if err != nil {
		mylog.Logger.Log(fmt.Sprintf("LoadHex failed: %s", err))
		return err
	}
	defer file.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		panic(err)
	}

	mylog.Logger.Log(fmt.Sprintf("LoadHex for cpu %s - loading file %s", cpu.GetName(), path))
	for idx, segment := range mem.GetDataSegments() {
		mylog.Logger.Log(fmt.Sprintf("%d addr %06x length %6x (%d)",
					idx, segment.Address, len(segment.Data), len(segment.Data)))
                for i := range segment.Data {
                        cpu.Write_8(segment.Address + uint32(i), segment.Data[i])
                }
	}
	mylog.Logger.Log(fmt.Sprintf("LoadHex done"))
	return nil
}
