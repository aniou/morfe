
package platform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcinbor85/gohex"
	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator"
)

func (p *Platform) LoadHex(cpu emu.Processor, filename string) {
	path := filepath.Join(filename)
	file, err := os.Open(path)
	if err != nil {
		mylog.Logger.Log(fmt.Sprintf("LoadHex failed: %s", err))
		return
	}
	defer file.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		panic(err)
	}

	mylog.Logger.Log(fmt.Sprintf("LoadHex loading: %s", path))
	for idx, segment := range mem.GetDataSegments() {
		mylog.Logger.Log(fmt.Sprintf("%d addr %06x length %6x (%d)",
					idx, segment.Address, len(segment.Data), len(segment.Data)))
                for i := range segment.Data {
                        cpu.Write_8(segment.Address + uint32(i), segment.Data[i])
                }
	}
	mylog.Logger.Log(fmt.Sprintf("LoadHex done"))
}
