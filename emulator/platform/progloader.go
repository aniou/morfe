
package platform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcinbor85/gohex"
)

func (p *Platform) LoadHex(filename string) {
	path := filepath.Join(filename)
	file, err := os.Open(path)
	if err != nil {
		p.Logger.Log(fmt.Sprintf("LoadHex failed: %s", err))
		return
	}
	defer file.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		panic(err)
	}

	p.Logger.Log(fmt.Sprintf("LoadHex loading: %s", path))
	for idx, segment := range mem.GetDataSegments() {
		p.Logger.Log(fmt.Sprintf("%d addr %06x length %6x (%d)",
					idx, segment.Address, len(segment.Data), len(segment.Data)))
                for i := range segment.Data {
                        p.CPU.Bus.EaWrite(segment.Address + uint32(i), segment.Data[i])
                }
	}
	p.Logger.Log(fmt.Sprintf("LoadHex done"))
}
