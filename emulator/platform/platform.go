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
        CPU0     emu.Processor           // on-board one (65c816 by default)
        CPU1     emu.Processor           // add-on
	//GPU      *vicky2.Vicky
	//GPU0     *vicky2.Vicky
	//GPU1     *vicky2.Vicky
	GPU      emu.GPU
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

