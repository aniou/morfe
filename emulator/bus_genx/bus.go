
package bus_genx

import (
	"fmt"
	"log"

	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/memory"
)

const MAX_MEM_SIZE = 0xff_ffff + 1
const PAGE_BITS    = 14
const PAGE_SIZE    = 1 << PAGE_BITS   // 4000 for 14 bits
const PAGE_MASK    = PAGE_SIZE-1      // 3fff for 4000
const SEGMENTS     = MAX_MEM_SIZE >> PAGE_BITS

type busEntry struct {
        mem   memory.Memory
        name  string
        start uint32
        end   uint32
}

type Bus struct {
        EA        uint32                  // last memory access - r/w
        Write     bool                    // is write op?
        segment   [2][SEGMENTS]memory.Memory
	entries   [2][]busEntry
}

func New() *Bus {
	b := Bus{}
	b.entries[0] = make([]busEntry, 0)
	b.entries[1] = make([]busEntry, 0)
	return &b
}

func (b *Bus) Attach(mem memory.Memory, name string, mode int, start uint32, end uint32) {

        if (start & PAGE_MASK) != 0 {
                log.Panicf("bus_genx: start are not properly aligned: %06X", start)
        }

        if ((end+1) & PAGE_MASK) != 0 {
                log.Panicf("bus_genx:   end are not properly aligned: %06X", end)
        }

        if ((end-start+1) % PAGE_SIZE) != 0 {
                log.Panicf("bus_genx:  size %06X is not multiplication of %04X", (end-start+1), PAGE_SIZE)
        }

        for x:=(start >> PAGE_BITS); x<=(end >> PAGE_BITS) ; x++ {
                //fmt.Printf("bus_genx: %v\n", x)
                b.segment[mode][x] = mem
        }

        entry := busEntry{mem: mem, name: name, start: start, end: end}
        mylog.Logger.Log(fmt.Sprintf("bus attach: %-20v %06x %06x %s", mem, start, end, name))
        b.entries[mode] = append(b.entries[mode], entry)

        return
}

