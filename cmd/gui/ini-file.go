package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	//"os"

	"gopkg.in/ini.v1"
)

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


func (g *GUI) loadConfig(filename string) {
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
		g.p.LoadHex(hexfile)
		//fmt.Printf("key file found, file %s\n", hexfile)
	}

	for i := 0; i<1000 && number>0; i += 1 {
		keyname := fmt.Sprintf("file%d", i)
		if cfg.Section("load").HasKey(keyname) {
			hexfile := cfg.Section("load").Key(keyname).String()
			g.p.LoadHex(hexfile)
			//fmt.Printf("key %s found\n", keyname)
			number -= 1
		}
	}

	// hex load -------------------------------------------------------------
	// CPU setting
	if cfg.Section("cpu").HasKey("start") {
		hex_addr := cfg.Section("cpu").Key("start").String()
		addr, _  := hex2uint24(hex_addr)
		fmt.Printf("start addr set: %06X\n", addr)
		g.p.CPU.PC = uint16(addr & 0x0000FFFF)
		g.p.CPU.RK = uint8(addr >> 16)
	}
}
