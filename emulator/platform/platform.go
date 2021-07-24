// adapted from fogleman/nes/console.go by Michael Fogleman
// <michael.fogleman@gmail.com>

package platform

import (
	//"log"

        "github.com/aniou/go65c816/lib/mylog"

        "github.com/aniou/go65c816/emulator"
        //"github.com/aniou/go65c816/emulator/bus_fmx"
        "github.com/aniou/go65c816/emulator/bus_genx"
        "github.com/aniou/go65c816/emulator/cpu65c816"
        "github.com/aniou/go65c816/emulator/cpu_dummy"
        //"github.com/aniou/go65c816/emulator/cpu68xxx"
        //"github.com/aniou/go65c816/emulator/memory"
        //"github.com/aniou/go65c816/emulator/vicky"
        "github.com/aniou/go65c816/emulator/vicky2"
        "github.com/aniou/go65c816/emulator/superio"
        "github.com/aniou/go65c816/emulator/ram"
        "github.com/aniou/go65c816/emulator/mathi"
)

type Platform struct {
        CPU0     emu.Processor           // on-board one (65c816 by default)
        CPU1     emu.Processor           // add-on
	GPU      *vicky2.Vicky
	GPU0     *vicky2.Vicky
	GPU1     *vicky2.Vicky
        SIO      *superio.SIO
	MATHI    *mathi.MathInt
}

func New() (*Platform) {
        p            := Platform{}
        return &p
}

/*
func (p *Platform) InitGenX() {
	bus0       := bus_genx.New()

        ram0       :=    ram.New("ram0", 1, 0x400000)
	gpu0       := vicky3.New("gpu0") 

        bus0.Attach(ram0,       0, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(gpu0.VRAM,  0, 0x40_0000, 0x7F_FFFF)
        bus0.Attach(gpu0,       0, 0xC4_0000, 0xC5_FFFF)
        bus0.Attach(gpu0.TEXT,  0, 0xC6_0000, 0xC6_3FFF)
        bus0.Attach(gpu0.COLOR, 0, 0xC6_4000, 0xC6_7FFF)

	/*
        bus0.Attach(nil,   "gpu0-vram",  0, 0x40_0000, 0x7F_FFFF) //  2 pages
        bus0.Attach(nil,   "gpu1-vram",  0, 0x80_0000, 0xBF_FFFF) //  2 pages
        bus0.Attach(nil,   "gabe",       0, 0xC0_0000, 0xC1_FFFF)
        bus0.Attach(nil,   "beatrix",    0, 0xC2_0000, 0xC3_FFFF)
        bus0.Attach(nil,   "gpu0-reg",   0, 0xC4_0000, 0xC5_FFFF)
        bus0.Attach(nil,   "gpu0-text",  0, 0xC6_0000, 0xC6_3FFF)
        bus0.Attach(nil,   "gpu0-color", 0, 0xC6_4000, 0xC6_7FFF)
        bus0.Attach(nil,   "reserved0",  0, 0xC6_8000, 0xC7_FFFF) // todo put placeholder for restricted access
        bus0.Attach(nil,   "gpu1-reg",   0, 0xC8_0000, 0xC9_FFFF)
        bus0.Attach(nil,   "gpu1-text",  0, 0xCA_0000, 0xCA_3FFF)
        bus0.Attach(nil,   "gpu1-color", 0, 0xCA_4000, 0xCA_7FFF)
        bus0.Attach(nil,   "reserved1",  0, 0xCA_8000, 0xCF_FFFF) // todo put placeholder for restricted access
        bus0.Attach(nil,   "reserved2",  0, 0xD0_0000, 0xDF_FFFF) // todo put placeholder for restricted access
        bus0.Attach(nil,   "dram0",      0, 0xE0_0000, 0xFF_FFFF) // 32 pages
	log.Panicln("it is ok to halt here")


}
        */

// something like FMX
func (p *Platform) InitFMX() {

	bus0       := bus_genx.New("bus0")
	bus1       := bus_genx.New("bus1")

	p.MATHI     =   mathi.New("mathi",       0x100)
        p.SIO       = superio.New("sio",         0x400)
	p.GPU       =  vicky2.New("gpu0",    0x01_0000)
        ram0       :=     ram.New("ram0", 1, 0x40_0000)  // single bank

        // FMX/U/U+ memory model
	//
	// $00:0000 - $1f:ffff - 2MB RAM
	//   $00:100 - $00:01ff - math core, IRQ CTRL, Timers, SDMA
	// $20:0000 - $3f:ffff - 2MB RAM on FMX revB and U+
	// $40:0000 - $ae:ffff - empty space (for example: extension card)
	// $af:0000 - $af:9fff - IO registers (mostly VICKY)
	//   $af:0800 - $af:080f - RTC
	//   $af:1000 - $af:13ff - GABE 
	// $af:a000 - $af:bfff - VICKY - TEXT  RAM
	// $af:c000 - $af:dfff - VICKY - COLOR RAM
	// $af:e000 - $af:ffff - IO registers (Trinity, Unity, GABE, SDCARD)
	// $b0:0000 - $ef:ffff - VIDEO RAM
	// $f0:0000 - $f7:ffff - 512KB System Flash
	// $f8:0000 - $ff:ffff - 512KB User Flash (if populated)

        bus0.Attach(ram0,       0, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(p.MATHI,    0, 0x00_0100, 0x00_01FF)
        bus0.Attach(p.GPU,      0, 0xAF_0000, 0xAF_FFFF)
        bus0.Attach(p.SIO,      0, 0xAF_1000, 0xAF_13FF)

	//log.Panicln("it is ok to halt here")


        p.CPU0     = cpu65c816.New(bus0, "cpu0")
        p.CPU1     = cpu_dummy.New(bus1, "cpu1")
        //p.CPU1     = cpu68xxx.New(bus1,  "cpu1")   // 20Mhz - not used yet
        
        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()                                      // XXX - move it to main binary?
        //p.CPU1.Reset()

        mylog.Logger.Log("platform initialized")
}
/*
// something like GenX
func (p *Platform) InitGUI() {
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
        //p.CPU1     = cpu68xxx.New(bus1,  "cpu1")   // 20Mhz - not used yet
        
        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()                                      // XXX - move it to main binary?

        p.CPU1.Reset()

        mylog.Logger.Log("platform initialized")
}
*/
