
// +build m68k

package platform

import (
	"log"

	"github.com/aniou/morfe/lib/mylog"

        "github.com/aniou/morfe/emulator"
        "github.com/aniou/morfe/emulator/bus"
        "github.com/aniou/morfe/emulator/cpu_65c816"
        "github.com/aniou/morfe/emulator/cpu_68xxx"
        "github.com/aniou/morfe/emulator/cpu_dummy"
        "github.com/aniou/morfe/emulator/vicky2"
        "github.com/aniou/morfe/emulator/vicky3"
        "github.com/aniou/morfe/emulator/superio"
        "github.com/aniou/morfe/emulator/ram"
        "github.com/aniou/morfe/emulator/mathi"
)



// Tentative Memory Map for the A2560U - PRELIMINARY
// MC68SEC000 Memory Map Model
//1Mx16   (2x 1Mx8) <- $0000_0000 - $001F_FFFF - RAM (U Model)
//2Mx16   (2x 2Mx8) <- $0000_0000 - $003F_FFFF - RAM (U + Model)
//                                          $0040_0000 - $00AF_FFFF - FREE SPACE (Future SDRAM Expansion Card?)

//                     $00B0_0000 - $00B1_FFFF - GABE Registers (SuperIO/Math Block/SDCard/IDE/Ethernet/SDMA)
//                     $00B2_0000 - $00B3_FFFF - BEATRIX Registers (CODEC/ADC/DAC0/DAC1/PSG/SID)
//                             $00B4_0000 - $00B5_FFFF - VICKY Registers
//                             $00B6_0000 - $00B6_3FFF - TEXT Memory
//                             $00B6_4000 - $00B6_7FFF - Color Memory

//                            $00BF_0000 - $00BF_FFFF - EXPANSION Chip Select
// 512Kx32 (2Mx8) <- $00C0_0000 - $00DF_FFFF - VRAM MEMORY
// 1Mx16   (2Mx8) <- $00E0_0000 - $00FF_FFFF - FLASH0

// TODO - parametrize this, U has 2megs and U+ - 4
func (p *Platform) SetA2560U() {

	log.Panic("work in progress, see a2560u branch for latest code for this platform!")

	p.Init     = p.InitA2560U

        bus0       := bus.New("bus0")
        bus1       := bus.New("bus1") // dummy one

        //p.MATHI     =   mathi.New("mathi",       0x100)  // not implemented yet
        //p.SIO       = superio.New("sio",         0x400)  // not implemented yet
        p.GPU       =  vicky3.New("gpu0",       0x2000)    // vram = 0x20_0000, text = 0x4000


        // m68k has RAM attached directly
        //bus0.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky3.F_MAIN, 0xB4_0000, 0xB5_FFFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky3.F_TEXT, 0xB6_0000, 0xB6_3FFF)
        bus0.Attach(emu.M_USER, p.GPU,  vicky3.F_TEXT_C, 0xB6_4000, 0xB6_7FFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky3.F_VRAM, 0xC0_0000, 0xDF_FFFF)

        //bus0.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(emu.M_SV,   p.GPU,    vicky3.F_MAIN, 0xB4_0000, 0xB5_FFFF)
        bus0.Attach(emu.M_SV,   p.GPU,    vicky3.F_TEXT, 0xB6_0000, 0xB6_3FFF)
        bus0.Attach(emu.M_SV,   p.GPU,  vicky3.F_TEXT_C, 0xB6_4000, 0xB6_7FFF)
        bus0.Attach(emu.M_SV,   p.GPU,    vicky3.F_VRAM, 0xC0_0000, 0xDF_FFFF)

	p.CPU0     = cpu_68xxx.New(bus0,  "cpu0") // TODO - add type? Or another routine for type? And pass RAM size
        p.CPU1     = cpu_dummy.New(bus1,  "cpu1")

        mylog.Logger.Log("platform: A2560-like created")

}

func (p *Platform) InitA2560U() {
        p.CPU0.Reset()
}


// a "frankenmode", not existing machine that starts 65c816
// but has active m68k
func (p *Platform) SetFranken() {
	p.Init  = p.InitFMX

        bus0       := bus.New("bus0")
        bus1       := bus.New("bus1")

        p.MATHI     =   mathi.New("mathi",       0x100)
        p.SIO       = superio.New("sio",         0x400)
        p.GPU       =  vicky2.New("gpu0",     0x1_0000)       // should be 0xA000 but there is no support for 0xAF:Exxx
        //p.GPU     =  vicky2.New("gpu0",       0xA000)
        ram0       :=     ram.New("ram0", 1, 0x40_0000)       // single bank                                                             

        bus0.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(emu.M_USER, p.MATHI,   mathi.F_MAIN, 0x00_0100, 0x00_01FF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky2.F_MAIN, 0xAF_0000, 0xAF_FFFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky2.F_TEXT, 0xAF_A000, 0xAF_BFFF)
        bus0.Attach(emu.M_USER, p.GPU,  vicky2.F_TEXT_C, 0xAF_C000, 0xAF_DFFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky2.F_VRAM, 0xB0_0000, 0xEF_FFFF)  // TODO - parametrize that
        bus0.Attach(emu.M_USER, p.SIO,   superio.F_MAIN, 0xAF_1000, 0xAF_13FF)


        // m68k has RAM attached directly
        //bus1.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus1.Attach(emu.M_USER, p.MATHI,   mathi.F_MAIN, 0x00_0100, 0x00_01FF)
        bus1.Attach(emu.M_USER, p.GPU,    vicky2.F_MAIN, 0xAF_0000, 0xAF_FFFF)
        bus1.Attach(emu.M_USER, p.GPU,    vicky2.F_TEXT, 0xAF_A000, 0xAF_BFFF)
        bus1.Attach(emu.M_USER, p.GPU,  vicky2.F_TEXT_C, 0xAF_C000, 0xAF_DFFF)
        bus1.Attach(emu.M_USER, p.GPU,    vicky2.F_VRAM, 0xB0_0000, 0xEF_FFFF)  // TODO - parametrize that
        bus1.Attach(emu.M_USER, p.SIO,   superio.F_MAIN, 0xAF_1000, 0xAF_13FF)

        //bus1.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus1.Attach(emu.M_SV  , p.MATHI,   mathi.F_MAIN, 0x00_0100, 0x00_01FF)
        bus1.Attach(emu.M_SV  , p.GPU,    vicky2.F_MAIN, 0xAF_0000, 0xAF_FFFF)
        bus1.Attach(emu.M_SV  , p.GPU,    vicky2.F_TEXT, 0xAF_A000, 0xAF_BFFF)
        bus1.Attach(emu.M_SV  , p.GPU,  vicky2.F_TEXT_C, 0xAF_C000, 0xAF_DFFF)
        bus1.Attach(emu.M_SV  , p.GPU,    vicky2.F_VRAM, 0xB0_0000, 0xEF_FFFF)  // TODO - parametrize that
        bus1.Attach(emu.M_SV  , p.SIO,   superio.F_MAIN, 0xAF_1000, 0xAF_13FF)


        p.CPU0     = cpu_65c816.New(bus0, "cpu0")
	p.CPU1     = cpu_68xxx.New(bus1,  "cpu1")

        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()

        mylog.Logger.Log("platform: frankenplatform created")

}


func (p *Platform) SetGenX() {
}

func (p *Platform) InitGenX() {
}

// an A2560K may have a memory layout similar to GenX-one
func (p *Platform) SetA2560K() {

	p.Init     = p.InitA2560K

        bus0       := bus.New("bus0")
        bus1       := bus.New("bus1") // dummy one

        //p.MATHI     =   mathi.New("mathi",       0x100)  // not implemented yet
        //p.SIO       = superio.New("sio",         0x400)  // not implemented yet
        p.GPU       =  vicky3.New("gpu0",       0x20000)    // vram = 0x20_0000, text = 0x4000
        flash0    :=     ram.New("flash0", 1, 0x10_0000)       // 1MB of "flash"


        // m68k has RAM attached directly
        //bus0.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky3.F_VRAM, 0x40_0000, 0x7F_FFFF)   
        bus0.Attach(emu.M_USER, p.GPU,    vicky3.F_MAIN, 0xC4_0000, 0xC5_FFFF)
        bus0.Attach(emu.M_USER, p.GPU,    vicky3.F_TEXT, 0xC6_0000, 0xC6_3FFF)
	bus0.Attach(emu.M_USER, p.GPU,  vicky3.F_TEXT_C, 0xC6_4000, 0xC6_7FFF)

        //bus0.Attach(emu.M_USER, ram0,        ram.F_MAIN, 0x00_0000, 0x3F_FFFF)
        bus0.Attach(emu.M_SV,   p.GPU,    vicky3.F_VRAM, 0x40_0000, 0x7F_FFFF)   
        bus0.Attach(emu.M_SV,   flash0,  ram.F_MAIN, 0xC0_0000, 0xCF_FFFF)

	p.CPU0     = cpu_68xxx.New(bus0,  "cpu0") // TODO - add type? Or another routine for type? And pass RAM size
        p.CPU1     = cpu_dummy.New(bus1,  "cpu1")

        mylog.Logger.Log("platform: A2560k-like created")

}

func (p *Platform) InitA2560K() {
        p.CPU0.Reset()
}
