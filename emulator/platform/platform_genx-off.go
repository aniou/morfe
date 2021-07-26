
// +build !m68k

package platform

import (
	"log"
)

func (p *Platform) InitGenX() {
	log.Panic("Emulator was built without m68k support")
}
