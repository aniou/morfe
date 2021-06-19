package vicky

import (
	"fmt"
	"encoding/binary"
	_ "sync"
	"github.com/aniou/go65c816/lib/mylog"
)

var bm0fb  []uint32	// bitmap0 framebuffer
var bm1fb  []uint32	// bitmap1 framebuffer
var tfb    []uint32	// text    framebuffer
var text   []uint32	// text memory cache
var fg     []uint32	// text foreground LUT cache
var bg     []uint32	// text background LUT cache
var mem    []byte	// main vicky memory area
var font   []byte	// font cache       : 256 chars  * 8 lines * 8 columns
var blut   []uint32	// bitmap LUT cache : 256 colors * 8 banks (lut0 to lut7)

type Vicky struct {
	TFB    []uint32		// text   framebuffer
	BM0FB  []uint32		// bitmap0 framebuffer
	BM1FB  []uint32		// bitmap1 framebuffer
	TEXT   []uint32
	MEM    []byte

	Master_L	byte	// MASTER_CTRL_REG_L
	Master_H	byte	// MASTER_CTRL_REG_H
	Cursor_visible  bool
	Border_visible  bool
	BM0_visible     bool
	BM1_visible     bool

        border_ctrl_reg byte
        Border_color_b  byte
        Border_color_g  byte
        Border_color_r  byte
        Border_x_size   int32
        Border_y_size   int32
	Background      [3]byte		// r, g, b

	starting_fb_row_pos uint32
	text_cols	uint32
	text_rows	uint32

	bm0_blut_pos    uint32
	bm1_blut_pos	uint32
	bm0_start_addr  uint32
	bm1_start_addr  uint32

	x_res		uint32		// preliminary, not fully supported
	y_res		uint32
	pixel_size	uint32		// 1 for normal, 2 for double

	//Mu_tfb		sync.Mutex
}

func init() {
	text = make([]uint32,  8192)
	fg   = make([]uint32,  8192)
	bg   = make([]uint32,  8192)
	mem  = make([]byte  ,  0x10_0000 + 0x40_0000)	// vicky and bitmap area
	tfb  = make([]uint32,    480000)		// for max 800x600
	bm0fb  = make([]uint32,  0x40_0000)		// max bitmap area - XXX - too large, we always write from 0x00
	bm1fb  = make([]uint32,  0x40_0000)		// max bitmap area - XXX - too large, we always write from 0x00
	font = make([]byte, 256 * 8 * 8)
	blut = make([]uint32, 256*8)
	fmt.Println("vicky areas are initialized")
}


func New() (*Vicky, error) {
	v := new(Vicky)
	v.TFB = tfb
	v.BM0FB = bm0fb
	v.BM1FB = bm1fb
	v.TEXT = text
	v.Cursor_visible = true
	v.BM0_visible    = true
	v.BM1_visible    = true
	v.border_ctrl_reg = 0x01
	v.Border_color_b = 0x20
	v.Border_color_g = 0x00
	v.Border_color_r = 0x20
	v.Border_x_size = 0x20
	v.Border_y_size = 0x20
	v.starting_fb_row_pos =  0x00
	v.text_cols = 0x00
	v.text_rows = 0x00
	v.bm0_blut_pos = 0x00
	v.bm1_blut_pos = 0x00
	v.bm0_start_addr = 0xB0_0000
	v.bm1_start_addr = 0xB0_0000

	v.x_res		= 640
	v.y_res		= 480
	v.pixel_size	= 1

	// XXX - just for test
	/*
        for i := range text { // file text memory areas
              fg[i] = 0x00
              bg[i] = 0x0d
              text[i] = 0x20
        } 
	*/

	v.recalculateScreen()

	return v, nil
}

// GUI-specific
// updates font cache by converting bits to bytes
// position - position of indyvidual byte in font bank
// val      - particular value
func updateFontCache(pos uint32, val byte) {
	pos = pos * 8
	for j := uint32(8); j > 0; j = j - 1 {		// counting down spares from shifting val left
		if (val & 1) == 1 {
			font[pos + j - 1] = 1
		} else {
			font[pos + j - 1] = 0
		}
		val = val >> 1
	}
}

// XXX - change it to new model (no recalculate text_cols without color change)
func (v *Vicky) recalculateScreen() {
	v.starting_fb_row_pos = v.x_res * uint32(v.Border_y_size) + uint32(v.Border_x_size)

        //v.text_cols = (640 - (uint32(v.Border_x_size) * 2)) / 8 // xxx - parametrize screen width
        //v.text_rows = (480 - (uint32(v.Border_y_size) * 2)) / 8 // xxx - parametrize screen height
	v.text_cols = (v.x_res - (uint32(v.Border_x_size) * 2)) / (v.pixel_size * 8)
	v.text_rows = (v.y_res - (uint32(v.Border_y_size) * 2)) / (v.pixel_size * 8)


        //if debug.gui {
                fmt.Printf("text_rows: %d\n", v.text_rows)
                fmt.Printf("text_cols: %d\n", v.text_cols)
        //}

}

// XXX - there is a room for improvements due to change from "128 bytes per line" model
//       to more linear one - there is no needed to recalculate position of new lines

func (v *Vicky) RenderBitmapText() {
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
	var is_overlay bool	  // is overlay text over bm0 enabled?

        // placeholders recalculated per row of text, holds values for text_cols loop
	// current max size is 100 columns (800/8)
        var fnttmp [128]uint32    // position in font array, from char value
        var fgctmp [128]uint32    // foreground color cache (rgba) for one line
        var bgctmp [128]uint32    // background color cache (rgba) for one line
        var dsttmp [128]uint32    // position in destination memory array

	cursor_state =        mem[0x00_0010]
	cursor_char  = uint32(mem[0x00_0012])
	cursor_x     = uint32(mem[0x00_0014])
	cursor_y     = uint32(mem[0x00_0016])
	
	if (mem[0x0000] & 0x02) == 0x02 {	// MASTER_CTRL_REG_L
		is_overlay = true
	} else {
		is_overlay = false
	}

	// render text - start
	// I prefer to keep it because it allow to simply re-drawing single line in future,
	// by manupipulating starting point (now 0) and end clause (now <v.text_rows)
	fb_row_pos = v.starting_fb_row_pos
	for text_y = 0; text_y < v.text_rows; text_y += 1 { // over lines of text
		text_row_pos = text_y * 80				// calculate from resolution / 8 !
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

			fgctmp[text_x] = binary.LittleEndian.Uint32(f_color_lut[f][:]) // text LUT - xxx: change name
			if is_overlay == false {
				bgctmp[text_x] = binary.LittleEndian.Uint32(b_color_lut[b][:]) // text LUT
			} else {
				bgctmp[text_x] = 0x00FFFFFF			// full alpha
			}
		}

		for font_line = 0; font_line < 8; font_line += 1 { // for every line of text - over 8 lines of font
			font_row_pos = font_line * 8
			for text_x = 0; text_x < v.text_cols; text_x += 1 { // for each line iterate over columns of text
				font_pos = fnttmp[text_x] + font_row_pos
				for i = 0; i < 8; i += 1 { // for every font iterate over 8 pixels of font
					//v.Mu_tfb.Lock()
					if font[font_pos+i] == 0 {
						tfb[fb_row_pos+dsttmp[text_x]+i] = bgctmp[text_x]
					} else {
						tfb[fb_row_pos+dsttmp[text_x]+i] = fgctmp[text_x]
					}
					//v.Mu_tfb.Unlock()
				}
			}
			fb_row_pos += v.x_res
		}
	}
	// render text - end
}

// still vicky-i
func (v *Vicky) RenderBitmapTextVickyI() {
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
	var is_overlay bool	  // is overlay text over bm0 enabled?

        // placeholders recalculated per row of text, holds values for text_cols loop --
        var fnttmp [128]uint32    // position in font array, from char value
        var fgctmp [128]uint32    // foreground color cache (rgba) for one line
        var bgctmp [128]uint32    // background color cache (rgba) for one line
        var dsttmp [128]uint32    // position in destination memory array

	cursor_state =        mem[0x00_0010]
	cursor_char  = uint32(mem[0x00_0012])
	cursor_x     = uint32(mem[0x00_0014])
	cursor_y     = uint32(mem[0x00_0016])
	
	if (mem[0x0000] & 0x02) == 0x02 {	// MASTER_CTRL_REG_L
		is_overlay = true
	} else {
		is_overlay = false
	}

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

			fgctmp[text_x] = binary.LittleEndian.Uint32(f_color_lut[f][:]) // text LUT - xxx: change name
			if is_overlay == false {
				bgctmp[text_x] = binary.LittleEndian.Uint32(b_color_lut[b][:]) // text LUT
			} else {
				bgctmp[text_x] = 0x00FFFFFF			// full alpha
			}
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
        return mem[addr:addr+0x10]         // XXX: configurable?
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

// XXX - balance between number of entries and possibility to identify
//       not implemented reads
func (v *Vicky) Read(address uint32) byte {
	a := address - 0xAF_0000

	switch {
	case address == 0xAF_0000:
		return mem[a]

	case address == 0xAF_0001:
		return mem[a]

	case address == 0xAF_0004:
		return mem[a]

	case address == 0xAF_0008:
		return mem[a]

	case address == 0xAF_0009:
		return mem[a]

	case address == 0xAF_000D:	// BACKGROUND_COLOR_B
		return mem[a]

	case address == 0xAF_000E:	// BACKGROUND_COLOR_G
		return mem[a]

	case address == 0xAF_000F:	// BACKGROUND_COLOR_R
		return mem[a]

	case address >= 0xAF_0010 && address<=0xAF_0017:	// cursor registers
		return mem[a]

	case address == 0xAF_0100:				// BM0_CONTROL_REG
		return mem[a]

	case address == 0xAF_0108:				// BM1_CONTROL_REG
		return mem[a]

	case address == 0xAF_070B:
		return byte(0)

	case address == 0xAF_070C:
		return byte(0)

        case address >= 0xAF_1F40 && address<=0xAF_1F7F:                                                                             
                a := address-0xAF_1F40
                byte_in_lut := byte(a & 0x03)
                num := byte(a >> 2)
                return f_color_lut[num][byte_in_lut]


        case address >= 0xAF_1F80 && address<=0xAF_1FFF:
                a := address-0xAF_1F80
                byte_in_lut := byte(a & 0x03)
                num := byte(a >> 2)
                return b_color_lut[num][byte_in_lut]

	case address >= 0xAF_2000 && address <= 0xAF_3CFF:	// GRPH_LUT0_PTR to GRPH_LUT7_PTR
		return mem[a]

	case address >= 0xAF_8000 && address <= 0xAF_87FF:	// FONT_MEMORY_BANK0
		return mem[a]

	case address >= 0xAF_8800 && address <= 0xAF_8FFF:	// FONT_MEMORY_BANK1  - xxx: NOT USED?
		return mem[a]

	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		//return byte(text[address-0xAF_A000])		// there is no need to recalculate
		return mem[a]					// because we have a 'generic' mem

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		//addr := address - 0xAF_C000			// same as above
		//fgc := byte(fg[addr]) << 4
		//bgc := byte(bg[addr])
		//return byte(fgc|bgc)
		return mem[a]

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
	mem[a] = val

	switch {
	case address == 0xAF_0000:				// MASTER_CTRL_REG_L
		v.Master_L = val

	case address == 0xAF_0001:				// MASTER_CTRL_REG_H
		v.Master_H = val

	case address == 0xAF_0004:                              // BORDER_CTRL_REG
		if (val & 0x01) == 1 {
			v.Border_x_size  = int32(mem[0x0008])
			v.Border_y_size  = int32(mem[0x0009])
			v.Border_visible = true
			v.recalculateScreen()
		} else {
			v.Border_x_size  = 0
			v.Border_y_size  = 0
			v.recalculateScreen()
			v.Border_visible = false
		}
		// xxx: v.recalculateScreen()

	case address == 0xAF_0005:                              // BORDER_COLOR_B
		v.Border_color_b = val
		v.recalculateScreen()

	case address == 0xAF_0006:				// BORDER_COLOR_G
		v.Border_color_g = val
		v.recalculateScreen()

	case address == 0xAF_0007:				// BORDER_COLOR_R
		v.Border_color_r = val
		v.recalculateScreen()

	case address == 0xAF_0008:				// BORDER_X_SIZE
		if v.Border_visible {
			v.Border_x_size = int32(val & 0x3F)	// XXX: in spec - 0-32, bitmask allows to 0-63
			v.recalculateScreen()
		}

	case address == 0xAF_0009:				// BORDER_Y_SIZE
		if v.Border_visible {
			v.Border_y_size = int32(val & 0x3F)	// XXX: in spec - 0-32, bitmask allows to 0-63
			v.recalculateScreen()
		}

	case address == 0xAF_000D:	// BACKGROUND_COLOR_B
		v.Background[2] = val

	case address == 0xAF_000E:	// BACKGROUND_COLOR_G
		v.Background[1] = val

	case address == 0xAF_000F:	// BACKGROUND_COLOR_R
		v.Background[0] = val

	case address == 0xAF_0010:	// xxx - implement rest of that
		if (val & 0x01) == 0 {
			v.Cursor_visible = false
		} else {
			v.Cursor_visible = true
		}
		
	case address >= 0xAF_0012 && address<=0xAF_0017:	// cursor registers
		return

	case address == 0xAF_0100:				// BM0_CONTROL_REG
		if (val & 0x01) == 0 {
			v.BM0_visible = false
		} else {
			v.BM0_visible = true
		}
		val = (val & 0x0E) >> 1				// extract LUT number
		v.bm0_blut_pos = uint32(val) * 0x100		// position in Bitmap LUT cache

		// XXX - todo: repopulate LUT if pos changed

	
	case address >= 0xAF_0101 && address<=0xAF_103:
		v.bm0_start_addr = 0xB0_0000 + (uint32(mem[0xAF_0103]) << 16) +
		                               (uint32(mem[0xAF_0102]) << 8 ) +
 				               (uint32(mem[0xAD_0101])      )

		// XXX - todo: recalculate bm0 framebuffer from new slice


	case address == 0xAF_0108:				// BM1_CONTROL_REG
		if (val & 0x01) == 0 {
			v.BM1_visible = false
		} else {
			v.BM1_visible = true
		}
		val = (val & 0x0E) >> 1				// extract LUT number
		v.bm1_blut_pos = uint32(val) * 0x100		// position in Bitmap LUT cache

		// XXX - todo: repopulate LUT if pos changed

	case address >= 0xAF_0109 && address<=0xAF_10B:
		v.bm1_start_addr = 0xB0_0000 + (uint32(mem[0xAF_010B]) << 16) +
		                               (uint32(mem[0xAF_010A]) << 8 ) +
 				               (uint32(mem[0xAD_0109])      )

		// XXX - todo: recalculate bm1 framebuffer from new slice

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

	// XXX - probably this needs correction with different
	//       bitmap format than ARGB
	case address >= 0xAF_2000 && address <= 0xAF_3FFF:	// GRPH_LUT0_PTR to GRPH_LUT7_PTR
		src := a & 0xfffffffc
		dst := ((address - 0xAF_2000) >>2 )		// clear bits 0-1, we need 4 bytes for in mem BGRA
								// in memory representation fo uint32: ARGB
		if (dst & 0xff) == 0 {
			blut[dst] = 0x00FFFFFF                  // LUTx[0] is always transparent, by design

		} else {
			blut[dst] = binary.LittleEndian.Uint32(
					[]byte{mem[src], 
					       mem[src+1], 
					       mem[src+2], 
					       mem[src+3],
				        })
		}
		//fmt.Printf("addr: %6x val %2x mem %4x dst: %4d pix: %08x ram: %v\n", address, val, src, dst, blut[dst], mem[src:src+4])

	case address >= 0xAF_8000 && address <= 0xAF_87FF:	// FONT_MEMORY_BANK0
		updateFontCache(address - 0xAF_8000, val)	// every bit in font cache is mapped to byte

	case address >= 0xAF_8800 && address <= 0xAF_8FFF:	// FONT_MEMORY_BANK1  - xxx: NOT USED?
		return
		// XXX - at this moment we don't use second bank at all
		//updateFontCache(address - 0xAF_8800, val)	// every bit in font cache is mapped to byte

	case address >= 0xAF_A000 && address<=0xAF_BFFF:
		text[address-0xAF_A000] = uint32(val)

	case address >= 0xAF_C000 && address<=0xAF_DFFF:
		addr := address - 0xAF_C000
		bgc := uint32( val & 0x0F)
		fgc := uint32((val & 0xF0)>> 4)
		fg[addr] = fgc
		bg[addr] = bgc

	case address >= 0xB0_0000 && address <= 0xEF_FFFF:                             // 4MB, xxx: parametrize
		if address >= v.bm0_start_addr && address<v.bm0_start_addr + 0x75300 { // max 800x600 bytes
			dst := address - v.bm0_start_addr
			//fmt.Printf("bm0fb addr: %6X dst: %6X val %2X blut %4X\n", address, dst, val, blut[v.bm0_blut_pos + uint32(val)])
			bm0fb[dst] = blut[v.bm0_blut_pos + uint32(val)]
		}
		if address >= v.bm1_start_addr && address<v.bm1_start_addr + 0x75300 {  // max 800x600 bytes
			dst := address - v.bm1_start_addr
			//fmt.Printf("bm1fb addr: %6X dst: %6X val %2X blut %4X\n", address, dst, val, blut[v.bm1_blut_pos + uint32(val)])
			bm1fb[dst] = blut[v.bm1_blut_pos + uint32(val)]
		}
	
	default:
		mylog.Logger.Log(fmt.Sprintf("vicky: write for addr %6X val %2X is not implemented", address, val))
	}
}

