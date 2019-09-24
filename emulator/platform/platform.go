// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/bus"
	"github.com/aniou/go65c816/emulator/cpu65c816"
	"github.com/aniou/go65c816/emulator/memory"
	"github.com/aniou/go65c816/emulator/netconsole"
)

type Platform struct {
	CPU    *cpu65c816.CPU
	Logger *mylog.MyLog
}

func New() (*Platform) {
	p            := Platform{nil, nil}
	return &p
}

func (platform *Platform) Init(logger *mylog.MyLog) {
	bus, _		:= bus.New(logger)
	platform.CPU, _  = cpu65c816.New(bus)
	console, _	:= netconsole.NewNetConsole(logger)
	ram, _	        := memory.New(0x40000)		// xxx - add logger?

	platform.CPU.Bus.Attach(ram,            "ram", 0x000000, 0x03FFFF)
	platform.CPU.Bus.Attach(console, "netconsole", 0x00DF00, 0x00DFFF)

        platform.CPU.Bus.EaWrite(0xFFFC, 0x00)
        platform.CPU.Bus.EaWrite(0xFFFD, 0x10)
	platform.CPU.Reset()

	platform.Logger = logger
	platform.Logger.Log("platform: initialized")
}

