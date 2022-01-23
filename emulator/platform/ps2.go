
package platform

import (
	"fmt"
	"github.com/aniou/morfe/emulator/emu"
)

// keyboard handling routine

// keyboard memory registers
const FMX_INT_MASK_REG1     = 0x00_014D
const FMX_INT_PENDING_REG1  = 0x00_0141

// at this moment we have some troubles with speed
// of keyboard routines on emu and sometimes is
// worthwhile to queue codes 
func (p *Platform) SendKeyFromQueue() {
		if p.CPU.Read_8(0xAF_1064) == 0 {
			e := p.PS2_queue.Front()
			p.PS2_queue.Remove(e)
			p.PS2.AddKeyCode(e.Value.(byte))
		} 
		
		irq1 := p.CPU.Read_8(FMX_INT_PENDING_REG1) | emu.R1_FNX1_INT00_KBD
		p.CPU.Write_8(FMX_INT_PENDING_REG1, irq1)
		p.CPU.TriggerIRQ(0)  // 65c816 hasn't levels at all
}

func (p *Platform) SendKey(code byte) {
	switch p.System {
	case emu.SYS_FOENIX_FMX:
		p.PS2_queue.PushBack(code)
		p.SendKeyFromQueue()
		fmt.Printf("ps2 queue len %d\n", p.PS2_queue.Len())
	case emu.SYS_FOENIX_A2560K:
		p.PS2.AddKeyCode(code)
		p.CPU.TriggerIRQ(3)
		// there should be a proper irq handling routine for m68k XXX

	default:
		fmt.Printf("platform: sendKey() called for unknown platform: %v (see emu.SYS_*)\n", p.System)
	}
}
