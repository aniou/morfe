
// +build !musashi

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

func (p *Platform) SetA2560U() {
	log.Panic("Emulator was built without m68k support")
}

func (p *Platform) SetA2560K() {
	log.Panic("Emulator was built without m68k support")
}
