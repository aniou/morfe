
package bus_genx

import (
	"fmt"
	"log"

	_ "github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator"
)

const MAX_MEM_SIZE = 0xff_ffff + 1
const PAGE_BITS    = 8   // 8 - 0x100 bytes, 10 - 0x400 bytes, 14 - 0x4000
const PAGE_SIZE    = 1 << PAGE_BITS
const PAGE_MASK    = PAGE_SIZE-1
const SEGMENTS     = MAX_MEM_SIZE >> PAGE_BITS

type busEntry struct {
        mem    emu.Memory
        offset uint32
}

type Bus struct {
        EA        uint32                  // last memory access - r/w
        Write     bool                    // is write op?
	name	  string
        segment   [2][SEGMENTS]busEntry
}

func New(name string) *Bus {
	b := Bus{name: name}
	fmt.Printf("bus_genx: %s max addr: %06X bits: %d page size: %04X page mask: %04x segments: %d\n",
			name, MAX_MEM_SIZE - 1, PAGE_BITS, PAGE_SIZE, PAGE_MASK, SEGMENTS)
	return &b
}

func (b *Bus) Attach(mem emu.Memory, mode int, start uint32, end uint32) {

	fmt.Printf("bus_genx: attaching mode %d start %06x end %06x name %s\n", mode, start, end, mem.Name())

        if (start & PAGE_MASK) != 0 {
                log.Panicf("bus_genx: start are not properly aligned: %06X", start)
        }

        if ((end+1) & PAGE_MASK) != 0 {
                log.Panicf("bus_genx:   end are not properly aligned: %06X", end)
        }

	region_size := end - start + 1
        if (region_size % PAGE_SIZE) != 0 {
                log.Panicf("bus_genx:  size %06X is not multiplication of %04X", region_size, PAGE_SIZE)
        }

	_, ram_size := mem.Size()
	if (region_size != ram_size) {
                log.Panicf("bus_genx:  region_size %06X does not match ram size %06X", region_size, ram_size)
	}


        for x:=(start >> PAGE_BITS); x<=(end >> PAGE_BITS) ; x++ {
                fmt.Printf("bus_genx: %06x %06x - %s\n", start, x, mem.Name())
		b.segment[mode][x] = busEntry{mem: mem, offset: start}
        }

        return
}

func (b *Bus) Write_8(mode byte, addr uint32, val byte) {
	//fmt.Printf("bus_genx: %s Write_8 mode %d addr %06x val %02x\n", b.name, mode, addr, val)
	s      := addr >> PAGE_BITS
	//fmt.Printf("bus_genx: %s Write_8 segment %06x\n", b.name, s)
	offset := b.segment[mode][s].offset
	//fmt.Printf("bus_genx: %s Write_8 offset  %06x\n", b.name, offset)
	//fmt.Printf("bus_genx: %s Write_8 name    %s\n", b.name, b.segment[mode][s].mem.Name() )
	b.segment[mode][s].mem.Write(addr - offset, val)
}

func (b *Bus) Read_8(mode byte, addr uint32) byte {
	//fmt.Printf("bus_genx: %s Read_8 mode %d addr %06x\n", b.name, mode, addr)
	s      := addr >> PAGE_BITS
	offset := b.segment[mode][s].offset
	return b.segment[mode][s].mem.Read(addr - offset)
}
