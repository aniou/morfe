// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	//"log"

        "github.com/aniou/morfe/emulator"
        "github.com/aniou/morfe/emulator/mathi"
        "github.com/aniou/morfe/emulator/superio"
        "github.com/aniou/morfe/emulator/vicky2"
)

type Platform struct {
        CPU0     emu.Processor           // on-board one (65c816 by default)
        CPU1     emu.Processor           // add-on
	GPU      *vicky2.Vicky
	GPU0     *vicky2.Vicky
	GPU1     *vicky2.Vicky
        SIO      *superio.SIO
	MATHI    *mathi.MathInt
	Init	 func()
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

