
// +build !m68k

package platform

import (
	"log"
)

func (p *Platform) SetGenX() {
	log.Panic("Emulator was built without m68k support")
}

func (p *Platform) SetFranken() {
	log.Panic("Emulator was built without m68k support")
}
