
// based on pda/go6502/bus/bus.go
// Copyright 2013–2014 Paul Annesley, released under MIT license
// Copyright 2019      Piotr Meyer

// XXX - should I define eaRead(uint32) and bankRead(uint8, uint16)?
// Or maybe Read(uint32)

package bus

import (
	"fmt"

	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/emulator/memory"
)

type busEntry struct {
	mem   memory.Memory
	name  string
	start uint32
	end   uint32
}


// original algorithm was built around array of busEntry segments.
// iterating over declared segments gives us a smallest memory usage
// at expense of O(n) cost when number of declared segments (n) grows

// there are two acceptable alternatives:
// 1 - arbitraly set an smallest, allowed segment to 16 bytes (4 bits)
//     and it gives us a O(1) computational cost at expense of max
//     2^20 array slots (4 bytes each - pointer to function size?)

// 2 - made the same as above but single byte pointer that selects
//     index from pointers table

// 3 - split above into two tables: [00-ff] for banks. Particular index
//     may be nil (no mapping or no memory present) or contains array
//     of fixed size [4096 pointers for example, when combined with 1]

// Bus is a 24-bit address passed via 32-bit variable, 8-bit data bus,
// which maps reads and writes at different locations to different backend
// Memory. For example the lower 32K could be RAM, the upper 8KB ROM, and
// some I/O in the middle.

// implement variant 1 - static table of 2^20 pointers to 16-bytes segments
// simplest, fastest and with greater memory usage
//
type Bus struct {
	segment [1048576]memory.Memory   // 2^10 because segments are 4bits length
	entries []busEntry
	Logger	*mylog.MyLog
}

func (b *Bus) String() string {
	return fmt.Sprintf("Address bus (TODO: describe)")
}

func (b *Bus) Clear() {
	//b.Mem.Clear()        // XXX - not implemented yet
}

func New(l *mylog.MyLog) (*Bus, error) {
	return &Bus{entries: make([]busEntry, 0), Logger: l}, nil
}


// There are two variants possible:
// handler, "name", start, size
// handler, "name", start, end      <- currently selected
func (b *Bus) Attach(mem memory.Memory, name string, start uint32, end uint32) error {

	if (start & 0xf) != 0 {
		b.Logger.Log(fmt.Sprintf("start are not 4-bit aligned: %06X", start))
		return fmt.Errorf("start are not 4-bit aligned: %06X", start)
	}

	if ((end+1)  & 0xf) != 0 {
		b.Logger.Log(fmt.Sprintf("bus: end are not 4-bit aligned: %06X", end))
		return fmt.Errorf("end are not 4-bit aligned: %06X", end)
	}

	for x:=(start>>4); x<=(end>>4); x++ {
		//fmt.Printf("%v", x)
		b.segment[x] = mem
	}
	//fmt.Printf("0x3ffff: %v\n", b.segment[0x3ffff>>4])


	entry := busEntry{mem: mem, name: name, start: start, end: end}
	b.Logger.Log(fmt.Sprintf("bus attach: %-20v %06x %06x", mem, start, end))
	b.entries = append(b.entries, entry)
	return nil
}



// XXX - crudy hack
func (b *Bus) backendFor(a uint32) (memory.Memory, error) {
	var tmpmem memory.Memory = nil
	tmpmem = b.segment[a>>4]
	if tmpmem == nil {
		//fmt.Printf("%v", b.segment)
		return nil, fmt.Errorf("No backend for address 0x%06X index %06x", a, a>>4)
	} else {
		return tmpmem, nil
	}
}


// Shutdown tells the address bus a shutdown is occurring, and to pass the
// message on to subordinates.
func (b *Bus) Shutdown() {
	for _, be := range b.entries {
		be.mem.Shutdown()
	}
}

// Read returns the byte from memory mapped to the given address.
// e.g. if ROM is mapped to 0xC000, then Read(0xC0FF) returns the byte at
// 0x00FF in that RAM device.
func (b *Bus) EaRead(a uint32) byte {
	mem, err := b.backendFor(a)
	if err != nil {
		panic(err)
	}
	value := mem.Read(a)
	return value
}

// Write the byte to the device mapped to the given address.
func (b *Bus) EaWrite(a uint32, value byte) {
	mem, err := b.backendFor(a)
	if err != nil {
		panic(err)
	}
	mem.Write(a, value)
}


// Dumps 16 bytes of memory, aligned to 16 bytes, used by
// memory viewers
func (b *Bus) EaDump(a uint32) (uint32, []byte) {
	a = a & 0x00fffff0          // round to segment
	start := a
	mem, err := b.backendFor(a)
	if err != nil {
		return start, nil
	}
	return start, mem.Dump(start)
}

