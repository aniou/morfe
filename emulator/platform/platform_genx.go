
// +build m68k

package platform

import (
	//"log"

	"github.com/aniou/morfe/lib/mylog"

        _ "github.com/aniou/morfe/emulator"
        "github.com/aniou/morfe/emulator/bus"
        "github.com/aniou/morfe/emulator/cpu_65c816"
        "github.com/aniou/morfe/emulator/cpu_68xxx"
        "github.com/aniou/morfe/emulator/vicky2"
        "github.com/aniou/morfe/emulator/superio"
        "github.com/aniou/morfe/emulator/ram"
        "github.com/aniou/morfe/emulator/mathi"
)

// a "frankenmode", not existing machine that starts 65c816
// but has active m68k
func (p *Platform) SetFranken() {
	p.Init  = p.InitFMX

        bus0       := bus.New("bus0")
        bus1       := bus.New("bus1")

        p.MATHI     =   mathi.New("mathi",       0x100)
        p.SIO       = superio.New("sio",         0x400)
        p.GPU       =  vicky2.New("gpu0",    0x01_0000 + 0x40_0000 ) // +bitmap area
        ram0       :=     ram.New("ram0", 1, 0x40_0000)              // single bank

        bus0.Attach(ram0,       0, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(p.MATHI,    0, 0x00_0100, 0x00_01FF)
        bus0.Attach(p.GPU,      0, 0xAF_0000, 0xEF_FFFF)
        bus0.Attach(p.SIO,      0, 0xAF_1000, 0xAF_13FF)

        //bus0.Attach(ram0,       0, 0x00_0000, 0x3F_FFFF)  m68k has RAM attached directly
        bus1.Attach(p.MATHI,    0, 0x00_0100, 0x00_01FF)
        bus1.Attach(p.GPU,      0, 0xAF_0000, 0xEF_FFFF)
        bus1.Attach(p.SIO,      0, 0xAF_1000, 0xAF_13FF)

        //bus0.Attach(ram0,       1, 0x00_0000, 0x3F_FFFF)  m68k has RAM attached directly
        bus1.Attach(p.MATHI,    1, 0x00_0100, 0x00_01FF)
        bus1.Attach(p.GPU,      1, 0xAF_0000, 0xEF_FFFF)
        bus1.Attach(p.SIO,      1, 0xAF_1000, 0xAF_13FF)


        p.CPU0     = cpu_65c816.New(bus0, "cpu0")
	p.CPU1     = cpu_68xxx.New(bus1,  "cpu1")

        mylog.Logger.Log("platform: frankenplatform created")

}


func (p *Platform) SetGenX() {
	/*
        p.Init  = p.InitGenX

        bus0       := bus.New("bus0")
        bus1       := bus.New("bus1")

        p.MATHI     =   mathi.New("mathi",       0x100)
        p.SIO       = superio.New("sio",         0x400)
        p.GPU       =  vicky2.New("gpu0",    0x01_0000 + 0x40_0000 ) // +bitmap area
        ram0       :=     ram.New("ram0", 1, 0x40_0000)              // single bank

        bus0.Attach(ram0,       0, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(p.MATHI,    0, 0x00_0100, 0x00_01FF)
        bus0.Attach(p.GPU,      0, 0xAF_0000, 0xEF_FFFF)
        bus0.Attach(p.SIO,      0, 0xAF_1000, 0xAF_13FF)

        p.CPU0     = cpu_65c816.New(bus0, "cpu0")
        p.CPU1     = cpu_dummy.New(bus1,  "cpu1")
        
        mylog.Logger.Log("platform: fmx-like created")




	bus0       := bus.New("bus0")
	bus1       := bus.New("bus1")

        ram0        :=    ram.New("ram0", 1, 0x400000)
	p.GPU0       =  vicky2.New("gpu0",    0x01_0000 + 0x40_0000 ) // +bitmap area

        bus0.Attach(ram0,       0, 0x00_0000, 0x3F_FFFF)
        //bus0.Attach(gpu0.VRAM,  0, 0x40_0000, 0x7F_FFFF)
        bus0.Attach(p.GPU,       0, 0xC4_0000, 0xC5_FFFF)
        //bus0.Attach(gpu0.TEXT,  0, 0xC6_0000, 0xC6_3FFF)
        //bus0.Attach(gpu0.COLOR, 0, 0xC6_4000, 0xC6_7FFF)

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

        p.CPU0     = cpu_65c816.New(bus0, "cpu0")
        p.CPU1     = cpu_68xxx.New(bus1,  "cpu1")
	*/

}

func (p *Platform) InitGenX() {
}
