
package platform

import (
	//"log"

        "github.com/aniou/morfe/lib/mylog"

        "github.com/aniou/morfe/emulator"
        "github.com/aniou/morfe/emulator/bus"
        "github.com/aniou/morfe/emulator/cpu_65c816"
        "github.com/aniou/morfe/emulator/cpu_dummy"
        "github.com/aniou/morfe/emulator/vicky2"
        _ "github.com/aniou/morfe/emulator/vicky3"
        "github.com/aniou/morfe/emulator/superio"
        _ "github.com/aniou/morfe/emulator/vram"
        "github.com/aniou/morfe/emulator/ram"
        "github.com/aniou/morfe/emulator/mathi"
)

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

// something like FMX
func (p *Platform) SetFMX() {
	p.Init  = p.InitFMX

	bus0       := bus.New("bus0")
	bus1       := bus.New("bus1")

	p.MATHI     =   mathi.New("mathi",       0x100)
        p.SIO       = superio.New("sio",         0x400)
	p.GPU       =  vicky2.New("gpu0",    0x01_0000 + 0x40_0000 ) // +bitmap area
        ram0       :=     ram.New("ram0", 1, 0x40_0000)              // single bank

	bus0.Attach(emu.M_USER, bus.BE{0x00_0000, 0x3F_FFFF,    ram0.Name,    ram0.Size,    ram0.Read,    ram0.Write})
        bus0.Attach(emu.M_USER, bus.BE{0x00_0100, 0x00_01FF, p.MATHI.Name, p.MATHI.Size, p.MATHI.Read, p.MATHI.Write})
        bus0.Attach(emu.M_USER, bus.BE{0xAF_0000, 0xEF_FFFF,   p.GPU.Name,   p.GPU.Size,   p.GPU.Read,   p.GPU.Write})
        bus0.Attach(emu.M_USER, bus.BE{0xAF_1000, 0xAF_13FF,   p.SIO.Name,   p.SIO.Size,   p.SIO.Read,   p.SIO.Write})

        p.CPU0     = cpu_65c816.New(bus0, "cpu0")
        p.CPU1     = cpu_dummy.New(bus1,  "cpu1")
        
	mylog.Logger.Log("platform: fmx-like created")
}

func (p *Platform) InitFMX() {

        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()

        p.CPU0.Write_8(0xAF_0005, 0x20) // border B                                                                                 
        p.CPU0.Write_8(0xAF_0006, 0x00) // border G
        p.CPU0.Write_8(0xAF_0007, 0x20) // border R

        p.CPU0.Write_8(0xAF_0010, 0x03) // VKY_TXT_CURSOR_CTRL_REG
        p.CPU0.Write_8(0xAF_0012, 0xB1) // VKY_TXT_CURSOR_CHAR_REG
        p.CPU0.Write_8(0xAF_0013, 0xED) // VKY_TXT_CURSOR_COLR_REG

	// On boot, Gavin copies the first 64KB of the content of System Flash 
	// (or User Flash, if present) to Bank $00.  The entire 512KB are copied 
	// to address range $18:0000 to $1F:FFFF (or 38:000 to 3F:FFFF)

        // act ersatz - copy jump table
        for j := 0x1000; j <= 0x1800; j = j + 1 {
                val := p.CPU0.Read_8(uint32(0x38_0000 + j))
                p.CPU0.Write_8(uint32(j), val)

        }

	mylog.Logger.Log("platform: fmx-like initialized")
}
