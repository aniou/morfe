package vicky

import (
	"fmt"
	"encoding/binary"
	"github.com/aniou/go65c816/lib/mylog"
)

type Vicky struct {
	FB     *[]uint32
	TEXT   *[8192]uint32
	FG     *[8192]uint32
	BG     *[8192]uint32
	FG_lut *[16][4]byte;
	BG_lut *[16][4]byte;

        border_ctrl_reg byte
        border_color_b  byte
        border_color_g  byte
        border_color_r  byte
        Border_x_size   uint32
        Border_y_size   uint32
}

func New() (*Vicky, error) {
	//vicky := Vicky{nil, nil, nil, nil, nil}
	vicky := Vicky{}
	return &vicky, nil
}

// GUI-specific
func (v *Vicky) FillByBorderColor() {
        val := binary.LittleEndian.Uint32([]byte{v.border_color_r, v.border_color_g, v.border_color_b, 0xff})                                             
        (*v.FB)[0] = val
        for bp := 1; bp < len((*v.FB)); bp *= 2 {
                copy((*v.FB)[bp:], (*v.FB)[:bp])
        }
}

// RAM-interface specific

func (v *Vicky) Dump(address uint32) []byte {
	return nil // XXX - todo
}

func (v *Vicky) String() string {
	return "vicky area"
}

func (v *Vicky) Shutdown() {
}

func (v *Vicky) Clear() { // Maybe Reset?
}

func (v *Vicky) Size() uint32 {
	return 0x100 // XXX: something
}

func (v *Vicky) Read(address uint32) byte {
	switch {
	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		return byte((*v.TEXT)[address-0xAF_A000])

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		fgc := byte((*v.FG)[addr]) << 4
		bgc := byte((*v.BG)[addr])
		return byte(fgc|bgc)
	
	default:
		mylog.Logger.Log(fmt.Sprintf("read from addr %6X is not implemented, 0 returned", address))
		return 0
	}
}

func (v *Vicky) Write(address uint32, val byte) {
	switch {
	case address == 0xAF_0005:
		v.border_color_b = val
		v.FillByBorderColor()

	case address == 0xAF_0006:
		v.border_color_g = val
		v.FillByBorderColor()

	case address == 0xAF_0007:
		v.border_color_r = val
		v.FillByBorderColor()

	case address == 0xAF_0008:
		v.Border_x_size = uint32(val & 0x3F)		// XXX: in spec - 0-32, bitmask allows to 0-63
		v.FillByBorderColor()

	case address == 0xAF_0009:
		v.Border_y_size = uint32(val & 0x3F)		// XXX: in spec - 0-32, bitmask allows to 0-63
		v.FillByBorderColor()

	case address >= 0xAF_1F40 && address<=0xAF_1F7F:
		a := address-0xAF_1F40
		byte_in_lut := byte(a & 0x03)
		num := byte(a >> 2)
		(*v.FG_lut)[num][byte_in_lut] = val

	case address >= 0xAF_1F80 && address<=0xAF_1FFF:
		a := address-0xAF_1F80
		byte_in_lut := byte(a & 0x03)
		num := byte(a >> 2)
		(*v.BG_lut)[num][byte_in_lut] = val

	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		(*v.TEXT)[address-0xAF_A000] = uint32(val)

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		bgc := uint32( val & 0x0F)
		fgc := uint32((val & 0xF0)>> 4)
		(*v.FG)[addr] = fgc
		(*v.BG)[addr] = bgc
	
	default:
		mylog.Logger.Log(fmt.Sprintf("write for addr %6X is not implemented", address))
	}
}

