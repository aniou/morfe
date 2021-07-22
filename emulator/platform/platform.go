// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	"log"

        "github.com/aniou/go65c816/lib/mylog"

        "github.com/aniou/go65c816/emulator"
        "github.com/aniou/go65c816/emulator/bus_fmx"
        "github.com/aniou/go65c816/emulator/bus_genx"
        "github.com/aniou/go65c816/emulator/cpu65c816"
        "github.com/aniou/go65c816/emulator/cpu68xxx"
        "github.com/aniou/go65c816/emulator/memory"
        "github.com/aniou/go65c816/emulator/vicky"
        "github.com/aniou/go65c816/emulator/gabe"
)

type Platform struct {
        CPU0    emu.Processor           // on-board one (65c816 by default)
        CPU1    emu.Processor           // add-on
        GPU     *vicky.Vicky
        GABE    *gabe.Gabe
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

// something like GenX
func (p *Platform) InitGUI() {

	bus2       := bus_genx.New()
        bus2.Attach(nil,   "ram0",       0, 0x00_0000, 0x3F_FFFF)
        bus2.Attach(nil,   "gpu0-vram",  0, 0x40_0000, 0x7F_FFFF) //  2 pages
        bus2.Attach(nil,   "gpu1-vram",  0, 0x80_0000, 0xBF_FFFF) //  2 pages
        bus2.Attach(nil,   "gabe",       0, 0xC0_0000, 0xC1_FFFF)
        bus2.Attach(nil,   "beatrix",    0, 0xC2_0000, 0xC3_FFFF)
        bus2.Attach(nil,   "gpu0-reg",   0, 0xC4_0000, 0xC5_FFFF)
        bus2.Attach(nil,   "gpu0-text",  0, 0xC6_0000, 0xC6_3FFF)
        bus2.Attach(nil,   "gpu0-color", 0, 0xC6_4000, 0xC6_7FFF)
        bus2.Attach(nil,   "reserved0",  0, 0xC6_8000, 0xC7_FFFF) // todo put placeholder for restricted access
        bus2.Attach(nil,   "gpu1-reg",   0, 0xC8_0000, 0xC9_FFFF)
        bus2.Attach(nil,   "gpu1-text",  0, 0xCA_0000, 0xCA_3FFF)
        bus2.Attach(nil,   "gpu1-color", 0, 0xCA_4000, 0xCA_7FFF)
        bus2.Attach(nil,   "reserved1",  0, 0xCA_8000, 0xCF_FFFF) // todo put placeholder for restricted access
        bus2.Attach(nil,   "reserved2",  0, 0xD0_0000, 0xDF_FFFF) // todo put placeholder for restricted access
        bus2.Attach(nil,   "dram0",      0, 0xE0_0000, 0xFF_FFFF) // 32 pages

	log.Panicln("it is ok to halt here")


	bus0,_	   := bus_fmx.New()
	bus1,_     := bus_fmx.New()
        ram0, _    := memory.New(0x400000, 0x000000)
        p.GPU, _    = vicky.New()
        p.GABE, _   = gabe.New()

        bus0.Attach(ram0,   "ram0",      0x00_0000, 0x3F_FFFF)  // xxx - 1: ram.offset, ram.size 2: get rid that?
        bus0.Attach(p.GPU,  "vicky",     0xAF_0000, 0xEF_FFFF)
        bus0.Attach(p.GABE, "gabe",      0xAF_1000, 0xAF_13FF)  // probably should be splitted
        bus0.Attach(p.GABE, "gabe-math", 0x00_0100, 0x00_012F)  // XXX error GABE coop is 0x2c bytes but we need mult of 16

        bus1.Attach(p.GPU,  "vicky",     0xAF_0000, 0xEF_FFFF)

        p.CPU0     = cpu65c816.New(bus0, "cpu0")
        p.CPU1     = cpu68xxx.New(bus1,  "cpu1")   // 20Mhz - not used yet
        
        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()                                      // XXX - move it to main binary?

        p.CPU1.Reset()

        mylog.Logger.Log("platform initialized")
}
