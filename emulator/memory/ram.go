
// memory management routines, adapted from
// NES emulator https://github.com/fogleman/nes
// Copyright (C) 2015 Michael Fogleman
// Copyright (c) 2019 Piotr Meyer

package memory

//import (
//	"fmt"
//)

type Ram struct {
	data	[]byte
	offset  uint32
}

func New(size uint32, offset uint32) (*Ram, error) {
	ram := make([]byte, size)
	mem := Ram{data: ram, offset: offset}
	return &mem, nil
}

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

func (mem *Ram) Read(address uint32) byte {
	addr := address - mem.offset
	return mem.data[addr]
}

func (mem *Ram) Size() uint32 {
	return uint32(len(mem.data))
}

func (mem *Ram) Shutdown() {
}

func (mem *Ram) String() string {
	return "RAM"
}

func (mem *Ram) Write(address uint32, value byte) {
	addr := address - mem.offset
	mem.data[addr] = value
}

