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
)

type Platform struct {
	CPU    *cpu65c816.CPU
	Logger *mylog.MyLog
	GPU    *vicky.Vicky
}

func New() (*Platform) {
	p            := Platform{nil, nil, nil}
	return &p
}

func (platform *Platform) Init(logger *mylog.MyLog) {
	bus, _		:= bus.New(logger)
	platform.CPU, _  = cpu65c816.New(bus)
	console, _	:= netconsole.NewNetConsole(logger)
	ram, _	        := memory.New(0x400000, 0x000000)		// xxx - add logger?
	//vicky, _	:= memory.New(0x010000,	0xaf0000)               // xxx - add logger?
	platform.GPU, _	 = vicky.New(logger)
	vram, _		:= memory.New(0x400000, 0xb00000)
	

	platform.CPU.Bus.Attach(ram,            "ram", 0x000000, 0x3FFFFF) // xxx - 1: ram.offset, ram.size 2: get rid that?
	platform.CPU.Bus.Attach(console, "netconsole", 0x000EF0, 0x000FFF)
	platform.CPU.Bus.Attach(platform.GPU, "vicky", 0xAF0000, 0xAFFFFF)
	platform.CPU.Bus.Attach(vram,          "vram", 0xB00000, 0xEFFFFF)

	platform.CPU.Bus.EaWrite(0xAF070B, 0x01)			   // fake platform version, my HW ha 43 here IDE has 00

        platform.CPU.Bus.EaWrite(0xFFFC, 0x00)
        platform.CPU.Bus.EaWrite(0xFFFD, 0x10)
	platform.CPU.Reset()

	platform.Logger = logger
	platform.Logger.Log("platform: initialized")
}

