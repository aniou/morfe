
package ram

import (
	_ "fmt"
)

type Ram struct {
	data	[][]byte
	name    string
	bank	int	// current bank
}

func New(name string, banks int, size uint32) *Ram {
	data := make([][]byte, banks)

	for i:=0; i<banks; i++ {
		data[i] = make([]byte, size)
	}
	mem := Ram{data: data, name: name, bank: 0}
	return &mem
}

func (mem *Ram) Write(addr uint32, val byte) {
	mem.data[mem.bank][addr] = val
}

func (mem *Ram) Read(addr uint32) byte {
	return mem.data[mem.bank][addr]
}

func (mem *Ram) Name() string {
	return mem.name
}

func (mem *Ram) Size() (uint32, uint32) {
	return uint32(len(mem.data)), 
	       uint32(len(mem.data[0]))
}

// XXX - to review
/*
func (mem *Ram) Clear() {
	for i := range mem.data {
		mem.data[i] = 0
	}
}

func (mem *Ram) Dump(start uint32) []byte {
	addr := start - mem.offset
	//fmt.Printf(" %06X - %06X - %06X \n", mem.offset, start, addr)
	//return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	return mem.data[addr:addr+0x10]		// XXX: configurable?
}

func (mem *Ram) Shutdown() {
}

*/
