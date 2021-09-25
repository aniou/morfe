// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	//"log"

        "github.com/aniou/morfe/emulator"
        "github.com/aniou/morfe/emulator/mathi"
        "github.com/aniou/morfe/emulator/superio"
)

type Platform struct {
        CPU      emu.Processor           // active processor
        CPU0     emu.Processor           // on-board one
        CPU1     emu.Processor           // add-on
	GPU      emu.GPU		 // active head on two-display nodes
	GPU0     emu.GPU
	GPU1     emu.GPU
        SIO      *superio.SIO
	MATHI    *mathi.MathInt
	Init	 func()
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

