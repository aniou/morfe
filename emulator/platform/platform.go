// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	//"log"

        "github.com/aniou/go65c816/emulator"
        "github.com/aniou/go65c816/emulator/mathi"
        "github.com/aniou/go65c816/emulator/superio"
        "github.com/aniou/go65c816/emulator/vicky2"
        "github.com/aniou/go65c816/emulator/vicky3"
)

/*

there is a problem here - a definition like this is not flexible
and force me to implement one-most-compatible-graphics-card of all

	GPU      *vicky2.Vicky
	GPU0     *vicky3.Vicky
	GPU1     *vicky3.Vicky

but...

	GPU      emu.Memory
	GPU0     emu.Memory
	GPU1     emu.Memory

does not allow me to access directly into exported variables
(memory areas). It is a big problem here.

*/

type Platform struct {
        CPU0     emu.Processor           // on-board one (65c816 by default)
        CPU1     emu.Processor           // add-on
	GPU      *vicky2.Vicky
	GPU0     *vicky3.Vicky
	GPU1     *vicky3.Vicky
        SIO      *superio.SIO
	MATHI    *mathi.MathInt
	Init	 func()
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

