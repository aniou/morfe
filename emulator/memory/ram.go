
// memory management routines, adapted from
// NES emulator https://github.com/fogleman/nes
// Copyright (C) 2015 Michael Fogleman
// Copyright (c) 2019 Piotr Meyer

package memory

type Ram struct {
	data	[]byte
}

func New(size uint32) (*Ram, error) {
	ram := make([]byte, size)
	mem := Ram{data: ram}
	return &mem, nil
}

func (mem *Ram) Clear() {
	for i := range mem.data {
		mem.data[i] = 0
	}
}

func (mem *Ram) Dump(start uint32) []byte {
	return mem.data[start:start+0x10]		// XXX: configurable?
}

func (mem *Ram) Read(address uint32) byte {
	return mem.data[address]
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
	mem.data[address] = value
}

