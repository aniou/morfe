package vicky

import (
	"fmt"
	"encoding/binary"
	"github.com/aniou/go65c816/lib/mylog"
)

var text []uint32
var fg   []uint32
var bg   []uint32
var mem  []byte

type Vicky struct {
	FB     *[]uint32
	TEXT   []uint32
	FG     []uint32
	BG     []uint32
	FG_lut *[16][4]byte
	BG_lut *[16][4]byte
	mem    []byte

        border_ctrl_reg byte
        border_color_b  byte
        border_color_g  byte
        border_color_r  byte
        Border_x_size   uint32
        Border_y_size   uint32
}

func init() {
	text = make([]uint32,  8192)
	fg   = make([]uint32,  8192)
	bg   = make([]uint32,  8192)
	mem  = make([]byte  , 65536)
	fmt.Println("vicky areas are initialized")
}


func New() (*Vicky, error) {
	//vicky := Vicky{nil, nil, nil, nil, nil}
	vicky := Vicky{nil, text, fg, bg, &f_color_lut, &b_color_lut, mem, 0x1, 0x20, 0x00, 0x20, 0x20, 0x20}
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
        addr := address - 0xAF_0000
        //fmt.Printf(" %06X - %06X - %06X \n", mem.offset, start, addr)
        //return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
        return v.mem[addr:addr+0x10]         // XXX: configurable?
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
	a := address - 0xAF_0000

	switch {
	case address == 0xAF_0001:
		return 0                // XXX should be resolution

	case address == 0xAF_0004:
		return v.border_ctrl_reg

	case address == 0xAF_0008:
		return byte(v.Border_x_size)

	case address == 0xAF_0009:
		return byte(v.Border_y_size)

	case address >= 0xAF_0010 && address<=0xAF_0017:	// cursor registers
		return mem[a]

	case address == 0xAF_070B:
		return byte(0)

	case address == 0xAF_070C:
		return byte(0)

	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		return byte(text[address-0xAF_A000])

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		fgc := byte(fg[addr]) << 4
		bgc := byte(bg[addr])
		return byte(fgc|bgc)

	case address == 0xAF_E80E:				// this is Trinity, not Vicky, XXX
		return 0x03					// BASIC
	
	default:
		mylog.Logger.Log(fmt.Sprintf("vicky: read from addr %6X is not implemented, 0 returned", address))
		return 0
	}
}

func (v *Vicky) Write(address uint32, val byte) {
	a := address - 0xAF_0000

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

	case address >= 0xAF_0010 && address<=0xAF_0017:	// cursor registers
		mem[a] = val

	case address >= 0xAF_1F40 && address<=0xAF_1F7F:
		a := address-0xAF_1F40
		byte_in_lut := byte(a & 0x03)
		num := byte(a >> 2)
		f_color_lut[num][byte_in_lut] = val

	case address >= 0xAF_1F80 && address<=0xAF_1FFF:
		a := address-0xAF_1F80
		byte_in_lut := byte(a & 0x03)
		num := byte(a >> 2)
		b_color_lut[num][byte_in_lut] = val

	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		text[address-0xAF_A000] = uint32(val)

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		bgc := uint32( val & 0x0F)
		fgc := uint32((val & 0xF0)>> 4)
		fg[addr] = fgc
		bg[addr] = bgc
	
	default:
		mylog.Logger.Log(fmt.Sprintf("vicky: write for addr %6X is not implemented", address))
	}
}

