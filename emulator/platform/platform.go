// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/bus"
	"github.com/aniou/go65c816/emulator/cpu"
	"github.com/aniou/go65c816/emulator/cpu65c816"
	"github.com/aniou/go65c816/emulator/memory"
	"github.com/aniou/go65c816/emulator/netconsole"
	"github.com/aniou/go65c816/emulator/vicky"
	"github.com/aniou/go65c816/emulator/gabe"
)

type Platform struct {
	CPU     *cpu65c816.CPU		// legacy, to be converted
	CPU0    *cpu.Processor		// on-board one (65c816 by default)
	CPU1    *cpu.Processor		// add-on
	GPU     *vicky.Vicky
	GABE    *gabe.Gabe
	Bus     *bus.Bus
	Console *netconsole.Console
}

func New() (*Platform) {
	p            := Platform{nil, nil, nil, nil, nil, nil, nil}
	return &p
}

// platform with Vicky I
func (platform *Platform) InitGUI() {
	platform.Bus, _	  = bus.New()
	platform.CPU, _   = cpu65c816.New(platform.Bus.EaRead, platform.Bus.EaWrite)	// to be removed
	ram, _	         := memory.New(0x400000, 0x000000)
	platform.GPU, _	  = vicky.New()
	platform.GABE, _  = gabe.New()
	

	platform.Bus.Attach(ram,           "ram",       0x00_0000, 0x3F_FFFF)  // xxx - 1: ram.offset, ram.size 2: get rid that?
	platform.Bus.Attach(platform.GPU,  "vicky",     0xAF_0000, 0xEF_FFFF)
	platform.Bus.Attach(platform.GABE, "gabe",      0xAF_1000, 0xAF_13FF)  // probably should be splitted
	platform.Bus.Attach(platform.GABE, "gabe-math", 0x00_0100, 0x00_012F)  // XXX error GABE coop is 0x2c bytes but we need mult of 16

	platform.Bus.EaWrite(0xAF070B, 0x01)			   // fake platform version, my HW ha 43 here IDE has 00
        platform.Bus.EaWrite(  0xFFFC, 0x00)			   // boot vector for 65c816
        platform.Bus.EaWrite(  0xFFFD, 0x10)
	platform.CPU.Reset()				           // XXX - move it to main binary?

	mylog.Logger.Log("platform: initialized")
}

/*
it does not wor at this moment

// simple platform with Text User Interface only
func (platform *Platform) InitTUI() {
	bus, _		:= bus.New()
	platform.CPU, _  = cpu65c816.New(bus)
	console, _	:= netconsole.NewNetConsole()
	platform.Console = console
	ram, _	        := memory.New(0x400000, 0x000000)
	platform.GPU, _	 = vicky.New()
	vram, _		:= memory.New(0x400000, 0xb00000)		   // XXX - placeholder
	

	platform.CPU.Bus.Attach(ram,            "ram", 0x000000, 0x3FFFFF) // xxx - 1: ram.offset, ram.size 2: get rid that?
	platform.CPU.Bus.Attach(console, "netconsole", 0x000EF0, 0x000FFF)
	platform.CPU.Bus.Attach(platform.GPU, "vicky", 0xAF0000, 0xAFFFFF)
	platform.CPU.Bus.Attach(vram,          "vram", 0xB00000, 0xEFFFFF)

	platform.CPU.Bus.EaWrite(0xAF070B, 0x01)			   // fake platform version, my HW ha 43 here IDE has 00

        platform.CPU.Bus.EaWrite(0xFFFC, 0x00)
        platform.CPU.Bus.EaWrite(0xFFFD, 0x10)
	platform.CPU.Reset()

	mylog.Logger.Log("platform: initialized")
}
*/
