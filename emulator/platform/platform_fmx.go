
package platform

import (
	//"log"

        "github.com/aniou/go65c816/lib/mylog"

        "github.com/aniou/go65c816/emulator/bus"
        "github.com/aniou/go65c816/emulator/cpu_65c816"
        "github.com/aniou/go65c816/emulator/cpu_dummy"
        "github.com/aniou/go65c816/emulator/vicky2"
        "github.com/aniou/go65c816/emulator/superio"
        "github.com/aniou/go65c816/emulator/ram"
        "github.com/aniou/go65c816/emulator/mathi"
)

// something like FMX
func (p *Platform) InitFMX() {

	bus0       := bus.New("bus0")
	bus1       := bus.New("bus1")

	p.MATHI     =   mathi.New("mathi",       0x100)
        p.SIO       = superio.New("sio",         0x400)
	//p.GPU       =  vicky2.New("gpu0",    0x01_0000)
	p.GPU       =  vicky2.New("gpu0",    0x01_0000 + 0x40_0000 ) // +bitmap area
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
        //bus0.Attach(p.GPU,      0, 0xAF_0000, 0xAF_FFFF)
        bus0.Attach(p.GPU,      0, 0xAF_0000, 0xEF_FFFF)
        bus0.Attach(p.SIO,      0, 0xAF_1000, 0xAF_13FF)

	//log.Panicln("it is ok to halt here")


        p.CPU0     = cpu_65c816.New(bus0, "cpu0")
        p.CPU1     = cpu_dummy.New(bus1,  "cpu1")
        //p.CPU1     = cpu_68xxx.New(bus1,  "cpu1")   // 20Mhz - not used yet
        
        p.CPU0.Write_8(  0xFFFC, 0x00)                      // boot vector for 65c816
        p.CPU0.Write_8(  0xFFFD, 0x10)
        p.CPU0.Reset()                                      // XXX - move it to main binary?
        //p.CPU1.Reset()

        mylog.Logger.Log("platform initialized")
}
