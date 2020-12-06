package vicky

import (
	"fmt"
	"encoding/binary"
	"github.com/aniou/go65c816/lib/mylog"
)

var bfb  []uint32
var tfb  []uint32
var text []uint32
var fg   []uint32
var bg   []uint32
var mem  []byte
var font []byte // 256 chars * 8 lines * 8 columns

type Vicky struct {
	TFB    []uint32		// text   framebuffer
	BFB    []uint32		// bitmap framebuffer
	TEXT   []uint32
	FG     []uint32
	BG     []uint32
	FG_lut *[16][4]byte
	BG_lut *[16][4]byte
	FONT   []byte
	mem    []byte

	Cursor_visible  bool

        border_ctrl_reg byte
        border_color_b  byte
        border_color_g  byte
        border_color_r  byte
        Border_x_size   uint32
        Border_y_size   uint32

	starting_fb_row_pos uint32
	text_cols	uint32
	text_rows	uint32
}

func init() {
	text = make([]uint32,  8192)
	fg   = make([]uint32,  8192)
	bg   = make([]uint32,  8192)
	mem  = make([]byte  , 65536 * 7) // bank $A0 to F0
	tfb  = make([]uint32,    480000)     // for max 800x600
	bfb  = make([]uint32, 65536 * 6) // whole bank $B0 for bitmap start from BM_START_ADDY
	font = make([]byte, 256 * 8 * 8)
	fmt.Println("vicky areas are initialized")
}


func New() (*Vicky, error) {
	//vicky := Vicky{nil, nil, nil, nil, nil}
	vicky := Vicky{tfb, bfb, text, fg, bg, &f_color_lut, &b_color_lut, font, mem, true, 0x1, 0x20, 0x00, 0x20, 0x20, 0x20, 0x00, 0x00, 0x00}
	return &vicky, nil
}

// GUI-specific
func (v *Vicky) FillByBorderColor() {
        val := binary.LittleEndian.Uint32([]byte{v.border_color_r, v.border_color_g, v.border_color_b, 0xff})                                             
        tfb[0] = val
        for bp := 1; bp < len(tfb); bp *= 2 {
                copy(tfb[bp:], tfb[:bp])
        }

	v.starting_fb_row_pos = 640*v.Border_y_size + (v.Border_x_size)
        v.text_cols = (640 - (v.Border_x_size * 2)) / 8 // xxx - parametrize screen width
        v.text_rows = (480 - (v.Border_y_size * 2)) / 8 // xxx - parametrize screen height
        //if debug.gui {
                fmt.Printf("text_rows: %d\n", v.text_rows)
                fmt.Printf("text_cols: %d\n", v.text_cols)
        //}

}

func (v *Vicky) RederBitmapText() {
        var cursor_x, cursor_y uint32 // row and column of cursor
        var cursor_char uint32    // cursor character
        var cursor_state byte     // cursor register, various states
        var text_x, text_y uint32 // row and column of text
        var text_row_pos uint32   // beginning of current text row in text memory
        var fb_row_pos uint32     // beginning of current FB   row in memory
        var font_pos uint32       // position in font array (char * 64 + char_line * 8)
        var font_line uint32      // line in current font
        var font_row_pos uint32   // position of line in current font (=font_line*8 because every line has 8 bytes)
        var i uint32              // counter

        // placeholders recalculated per row of text, holds values for text_cols loop --
        var fnttmp [128]uint32    // position in font array, from char value
        var fgctmp [128]uint32    // foreground color cache (rgba) for one line
        var bgctmp [128]uint32    // background color cache (rgba) for one line
        var dsttmp [128]uint32    // position in destination memory array

	cursor_state =        mem[0x00_0010]
	cursor_char  = uint32(mem[0x00_0012])
	cursor_x     = uint32(mem[0x00_0014])
	cursor_y     = uint32(mem[0x00_0016])

	// render text - start
	fb_row_pos = v.starting_fb_row_pos
	for text_y = 0; text_y < v.text_rows; text_y += 1 { // over lines of text
		text_row_pos = text_y * 128
		for text_x = 0; text_x < v.text_cols; text_x += 1 { // pre-calculate data for x-axis
			fnttmp[text_x] = text[text_row_pos+text_x] * 64 // position in font array
			dsttmp[text_x] = text_x * 8                     // position of char in dest FB

			f := fg[text_row_pos+text_x] // fg and bg colors
			b := bg[text_row_pos+text_x]

			if v.Cursor_visible && (cursor_y == text_y) && (cursor_x == text_x) && (cursor_state & 0x01 == 1) {
				f = uint32((mem[0x00_0013] & 0xf0) >> 4)
				b = uint32(mem[0x00_0013] & 0x0f)
				fnttmp[text_x] = cursor_char * 64
			}

			fgctmp[text_x] = binary.LittleEndian.Uint32(v.FG_lut[f][:])
			bgctmp[text_x] = binary.LittleEndian.Uint32(v.BG_lut[b][:])
		}

		for font_line = 0; font_line < 8; font_line += 1 { // for every line of text - over 8 lines of font
			font_row_pos = font_line * 8
			for text_x = 0; text_x < v.text_cols; text_x += 1 { // for each line iterate over columns of text
				font_pos = fnttmp[text_x] + font_row_pos
				for i = 0; i < 8; i += 1 { // for every font iterate over 8 pixels of font
					if font[font_pos+i] == 0 {
						tfb[fb_row_pos+dsttmp[text_x]+i] = bgctmp[text_x]
					} else {
						tfb[fb_row_pos+dsttmp[text_x]+i] = fgctmp[text_x]
					}
				}
			}
			fb_row_pos += 640
		}
	}
	// render text - end
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
	
	case address >= 0xB0_0000 && address <= 0xBF_FFFF:
		return mem[a]
	
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

	case address >= 0xB0_0000 && address <= 0xBF_FFFF:
		mem[a] = val
		if val == 2 {
			bfb[a-0x01_0000] = (uint32(val) << 8) | 0x000000FF	// just for test, no LUT at this moment
		} else {
			fmt.Printf("> %d %d\n", a, val)
			bfb[a-0x01_0000] = 0xFFFF00FF	// just for test, no LUT at this moment
		}

		//bfb[a-0x01_0000] = 0x13458BFF 	// just for test, no LUT at this moment
	
	default:
		mylog.Logger.Log(fmt.Sprintf("vicky: write for addr %6X is not implemented", address))
	}
}

