// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	"container/list"
	//"log"

        "github.com/aniou/morfe/emulator/emu"
        "github.com/aniou/morfe/emulator/mathi"
        "github.com/aniou/morfe/emulator/pata"
        "github.com/aniou/morfe/emulator/ps2"
)

type Platform struct {
        CPU      emu.Processor           // active processor
        CPU0     emu.Processor           // on-board one
        CPU1     emu.Processor           // add-on
	GPU      emu.GPU		 // active head on two-display nodes
	GPU0     emu.GPU
	GPU1     emu.GPU
	PATA0    *pata.PATA
	PS2	 *ps2.PS2
	MATHI    *mathi.MathInt
	
	System   byte			 // system type, const emu.SYS_*

	PS2_queue *list.List             // queue for ps2 scancodes

	Init	 func()
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

