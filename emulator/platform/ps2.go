
package platform

import (
	"fmt"
	"github.com/aniou/morfe/emulator/emu"
)

// keyboard handling routine

// keyboard memory registers
const FMX_INT_MASK_REG1     = 0x00_014D
const FMX_INT_PENDING_REG1  = 0x00_0141


func (p *Platform) SendKey(code byte) {
	switch p.System {
	case emu.SYS_FOENIX_FMX:
		mask := p.CPU.Read_8(FMX_INT_MASK_REG1)
		if (^mask & byte(emu.R1_FNX1_INT00_KBD)) == byte(emu.R1_FNX1_INT00_KBD) {
			p.SIO.Data = code + 0x80
			p.CPU.Write_8(0xAF_1064, 0)
			irq1 := p.CPU.Read_8(FMX_INT_PENDING_REG1) | emu.R1_FNX1_INT00_KBD
			p.CPU.Write_8(FMX_INT_PENDING_REG1, irq1)
			p.CPU.TriggerIRQ()
		}
	default:
		fmt.Printf("platform: sendKey() called for unknown platform: %v (see emu.SYS_*)\n", p.System)
	}
}
