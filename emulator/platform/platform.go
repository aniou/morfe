// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/bus"
	"github.com/aniou/go65c816/emulator/cpu65c816"
	"github.com/aniou/go65c816/emulator/memory"
	"github.com/aniou/go65c816/emulator/netconsole"
	"github.com/aniou/go65c816/emulator/vicky"
	"github.com/aniou/go65c816/emulator/gabe"
)

type Platform struct {
	CPU    *cpu65c816.CPU
	GPU    *vicky.Vicky
	GABE   *gabe.Gabe
	Console *netconsole.Console
}

func New() (*Platform) {
	p            := Platform{nil, nil, nil, nil}
	return &p
}

// platform with Vicky I
func (platform *Platform) InitGUI() {
	bus, _		 := bus.New()
	platform.CPU, _   = cpu65c816.New(bus)
	ram, _	         := memory.New(0x400000, 0x000000)
	platform.GPU, _	  = vicky.New()
	//vram, _		 := memory.New(0x400000, 0xb00000)		   // XXX - placeholder
	platform.GABE, _  = gabe.New()
	

	platform.CPU.Bus.Attach(ram,            "ram", 0x000000, 0x3FFFFF) // xxx - 1: ram.offset, ram.size 2: get rid that?
	platform.CPU.Bus.Attach(platform.GPU, "vicky", 0xAF0000, 0xEFFFFF)
	platform.CPU.Bus.Attach(platform.GABE, "gabe", 0xAF1000, 0xAF13FF) // probably should be splitted

	platform.CPU.Bus.EaWrite(0xAF070B, 0x01)			   // fake platform version, my HW ha 43 here IDE has 00

        platform.CPU.Bus.EaWrite(0xFFFC, 0x00)				   // boot vector
        platform.CPU.Bus.EaWrite(0xFFFD, 0x10)
	platform.CPU.Reset()

	mylog.Logger.Log("platform: initialized")
}

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

