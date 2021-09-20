
package vram

// a simplest, directly accessible, non-banking RAM

import (
)

type Vram struct {
	Data	[]uint32
	name    string
}

func New(name string, size uint32) *Vram {
	data := make([]uint32, size)
	mem  := Vram{Data: data, name: name}
	return &mem
}

func (mem *Vram) Write(addr uint32, val byte) error {
	mem.Data[addr] = uint32(val)
	return nil
}

func (mem *Vram) Read(addr uint32) (byte, error) {
	return byte(mem.Data[addr]), nil
}

func (mem *Vram) Name() string {
	return mem.name
}

func (mem *Vram) Size() (uint32, uint32) {
	return uint32(1), 
	       uint32(len(mem.Data))
}

