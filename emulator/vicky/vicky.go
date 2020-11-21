package vicky

import (
	//"fmt"
	//"github.com/aniou/go65c816/lib/mylog"
)

type Vicky struct {
	FB     *[8192]uint32
	FG     *[8192]uint32
	BG     *[8192]uint32
	FG_lut *[16][4]byte;
	BG_lut *[16][4]byte;
}

func New() (*Vicky, error) {
	vicky := Vicky{nil, nil, nil, nil, nil}
	return &vicky, nil
}

func (vicky *Vicky) Dump(address uint32) []byte {
	return nil // XXX - todo
}

func (vicky *Vicky) String() string {
	return "vicky area"
}

func (vicky *Vicky) Shutdown() {
}

func (vicky *Vicky) Clear() { // Maybe Reset?
}

func (vicky *Vicky) Size() uint32 {
	return 0x100 // XXX: something
}

func (vicky *Vicky) Read(address uint32) byte {
	switch {
	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		return byte((*vicky.FB)[address-0xAF_A000])
	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		fgc := byte((*vicky.FG)[addr]) << 4
		bgc := byte((*vicky.BG)[addr])
		return byte(fgc|bgc)
	
	default:
		return 0
	}
}

func (vicky *Vicky) Write(address uint32, val byte) {
	switch {
	case address >= 0xAF_1F40 && address<=0xAF_1F7F:
		a := address-0xAF_1F40
		byte_in_lut := byte(a & 0x03)
		num := byte(a >> 2)
		(*vicky.FG_lut)[num][byte_in_lut] = val

	case address >= 0xAF_1F80 && address<=0xAF_1FFF:
		a := address-0xAF_1F80
		byte_in_lut := byte(a & 0x03)
		num := byte(a >> 2)
		(*vicky.BG_lut)[num][byte_in_lut] = val

	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		(*vicky.FB)[address-0xAF_A000] = uint32(val)

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		bgc := uint32( val & 0x0F)
		fgc := uint32((val & 0xF0)>> 4)
		(*vicky.FG)[addr] = fgc
		(*vicky.BG)[addr] = bgc
	
	default:
	}
}

