package vicky2

// re-implementing Vicky2 at new bus

import (
        "encoding/binary"
        "fmt"
        "log"
        _ "sync"
        _ "github.com/aniou/morfe/lib/mylog"
        "github.com/aniou/morfe/emulator"
        _ "github.com/aniou/morfe/emulator/ram"
)

const F_MAIN   = 0
const F_TEXT   = 1
const F_TEXT_C = 2
const F_VRAM   = 3

const MASTER_CTRL_REG_L  = 0x0000
const MASTER_CTRL_REG_H  = 0x0001
const GAMMA_CTRL_REG     = 0x0002
//                         0x0003 - reserved
const BORDER_CTRL_REG    = 0x0004
const BORDER_COLOR_B     = 0x0005
const BORDER_COLOR_G     = 0x0006
const BORDER_COLOR_R     = 0x0007
const BORDER_X_SIZE      = 0x0008
const BORDER_Y_SIZE      = 0x0009
//                         0x000a
//                         0x000b
//                         0x000c
const BACKGROUND_COLOR_B = 0x000d
const BACKGROUND_COLOR_G = 0x000e
const BACKGROUND_COLOR_R = 0x000f

const VKY_TXT_CURSOR_CTRL_REG = 0x0010
const VKY_TXT_CURSOR_CHAR_REG = 0x0012
const VKY_TXT_CURSOR_COLR_REG = 0x0013
const VKY_TXT_CURSOR_X_REG_L  = 0x0014
const VKY_TXT_CURSOR_X_REG_H  = 0x0015
const VKY_TXT_CURSOR_Y_REG_L  = 0x0016
const VKY_TXT_CURSOR_Y_REG_H  = 0x0017

const BM0_CONTROL_REG         = 0x0100
const BM0_START_ADDY_L        = 0x0101
const BM0_START_ADDY_M        = 0x0102
const BM0_START_ADDY_H        = 0x0103

const BM1_CONTROL_REG         = 0x0108
const BM1_START_ADDY_L        = 0x0109
const BM1_START_ADDY_M        = 0x010a
const BM1_START_ADDY_H        = 0x010b

const FG_CHAR_LUT_PTR         = 0x1f40
const BG_CHAR_LUT_PTR         = 0x1f40

const GRPH_LUT0_PTR           = 0x2000
const GRPH_LUT7_PTR           = 0x3c00
const FONT_MEMORY_BANK0       = 0x8000
const CS_TEXT_MEM_PTR         = 0xa000
const CS_COLOR_MEM_PTR        = 0xc000

const VRAM_START              = 0x01_0000


type Vicky struct {
        name    string         // id of instance
        Mem     []byte         // general Vicky memory

        text    []uint32       // text memory
        vram    []byte         // VRAM
        tc      []byte         // text color memory
        blut    []uint32       // bitmap LUT cache : 256 colors * 8 banks (lut0 to lut7)
        fg      []uint32       // text foreground LUT cache
        bg      []uint32       // text background LUT cache
        font    []byte         // font cache       : 256 chars  * 8 lines * 8 columns

	c	emu.GPU_common	// common GPU-exported properties and framebuffers

	/*
        TFB     []uint32       // text   framebuffer
        BM0FB   []uint32       // bitmap0 framebuffer
        BM1FB   []uint32       // bitmap1 framebuffer

        // some convinient registers that should be converted
        // into some kind of memory indexes...
        Master_L        byte    // MASTER_CTRL_REG_L
        Master_H        byte    // MASTER_CTRL_REG_H
        Cursor_visible  bool
        Border_visible  bool
        BM0_visible     bool
        BM1_visible     bool

        Border_color_b  byte
        Border_color_g  byte
        Border_color_r  byte
        Border_x_size   int32
        Border_y_size   int32
        Background      [3]byte         // r, g, b
	*/

        starting_fb_row_pos uint32
        text_cols       uint32
        text_rows       uint32

        bm0_blut_pos    uint32
        bm1_blut_pos    uint32
        bm0_start_addr  uint32
        bm1_start_addr  uint32

        x_res           uint32          // preliminary, not fully supported
        y_res           uint32
        pixel_size      uint32          // 1 for normal, 2 for double

}

func New(name string, size int) *Vicky {
        v       := Vicky{name: name}

	v.Mem    = make([]byte,   size)		 // main memory - TODO: shrink it!

        v.text   = make([]uint32,    0x2000)        // text memory  - 0x4000 in GenX
	v.vram   = make([]byte  , 0x40_0000)        // 4MB - TODO: settable
	v.tc     = make([]byte,      0x2000)	    // text color memory - 0x4000 in GenX
        v.fg     = make([]uint32,    0x2000)        // foreground cache -  0x4000 in GenX
        v.bg     = make([]uint32,    0x2000)        // background color cache - 0x4000 in GenX


        v.blut   = make([]uint32, 0x0800)        // bitmap LUT cache : 256 colors * 8 banks (lut0 to lut7)
        v.font   = make([]byte  , 0x100 * 8 * 8) // font cache 256 chars * 8 lines * 8 columns
        //v.c.TFB    = make([]uint32, 480000)        // for max 800x600
        //v.c.BM0FB  = make([]uint32,  0x40_0000)    // max bitmap area - XXX - too large, we always write from 0x00
        //v.c.BM1FB  = make([]uint32,  0x40_0000)    // max bitmap area - XXX - too large, we always write from 0x00
        v.c.TFB    = make([]uint32,  0x0c_0000)    // text framebuffer - 1024x768 max
        v.c.BM0FB  = make([]uint32,  0x0c_0000)    // bm0  framebuffer - 1024x768 max
        v.c.BM1FB  = make([]uint32,  0x0c_0000)    // bm1  framebuffer - 1024x768 max
        

        v.Mem[ BORDER_CTRL_REG ] = 0x01

	// some tests
	//fmt.Printf("vicky2: v.TFB    %p\n", &v.TFB  )
	//fmt.Printf("vicky2: v.BM0FB  %p\n", &v.BM0FB)
	//fmt.Printf("vicky2: v.BM1FB  %p\n", &v.BM1FB)

	v.c.Master_H	   = 0x00 // 640x480 no pixel doubling
        v.c.Cursor_visible = true
        v.c.BM0_visible    = true
        v.c.BM1_visible    = true
        v.c.Border_color_b = 0x20
        v.c.Border_color_g = 0x00
        v.c.Border_color_r = 0x20
        v.c.Border_x_size  = 0x20
        v.c.Border_y_size  = 0x20
        v.starting_fb_row_pos =  0x00
        v.text_cols = 0x00
        v.text_rows = 0x00
        v.bm0_blut_pos = 0x00
        v.bm1_blut_pos = 0x00
        v.bm0_start_addr = 0x00	// relative from beginning of VRAM
        v.bm1_start_addr = 0x00	// relative from beginning of VRAM

        v.x_res         = 640
        v.y_res         = 480
        v.pixel_size    = 1

        // XXX - just for test
        /*
        for i := range text { // file text memory areas
              fg[i] = 0x00
              bg[i] = 0x0d
              text[i] = 0x20
        } 
        */

        v.recalculateScreen()

        return &v
}

func (v *Vicky) GetCommon() *emu.GPU_common {
	return &v.c
}

// GUI-specific
// updates font cache by converting bits to bytes
// position - position of indyvidual byte in font bank
// val      - particular value
func (v *Vicky) updateFontCache(pos uint32, val byte) {
        pos = pos * 8
        for j := uint32(8); j > 0; j = j - 1 {          // counting down spares from shifting val left
                if (val & 1) == 1 {
                        v.font[pos + j - 1] = 1
                } else {
                        v.font[pos + j - 1] = 0
                }
                val = val >> 1
        }
}

func (v *Vicky) recalculateScreen() {
        v.starting_fb_row_pos = v.x_res * uint32(v.c.Border_y_size) + uint32(v.c.Border_x_size)

        //v.text_cols = (640 - (uint32(v.Border_x_size) * 2)) / 8 // xxx - parametrize screen width
        //v.text_rows = (480 - (uint32(v.Border_y_size) * 2)) / 8 // xxx - parametrize screen height
        //v.text_cols = (v.x_res - (uint32(v.Border_x_size) * 2)) / (v.pixel_size * 8)
        //v.text_rows = (v.y_res - (uint32(v.Border_y_size) * 2)) / (v.pixel_size * 8)
        v.text_cols = uint32(v.x_res / 8)
        v.text_rows = uint32(v.y_res / 8)

        fmt.Printf("text_rows: %d\n", v.text_rows)
        fmt.Printf("text_cols: %d\n", v.text_cols)

}

func (v *Vicky) RenderBitmapText() {
        var cursor_x, cursor_y uint32 // row and column of cursor
        var cursor_state byte     // cursor register, various states
        var text_x, text_y uint32 // row and column of text
        var text_row_pos uint32   // beginning of current text row in text memory
        var fb_row_pos uint32     // beginning of current FB   row in memory
        var font_pos uint32       // position in font array (char * 64 + char_line * 8)
        var   fb_pos uint32       // position in destination framebuffer
        var font_line uint32      // line in current font
        var font_row_pos uint32   // position of line in current font (=font_line*8 because every line has 8 bytes)
        var i uint32              // counter
        var is_overlay bool       // is overlay text over bm0 enabled?

        // placeholders recalculated per row of text, holds values for text_cols loop
        // current max size is 100 columns (800/8)
        var fnttmp [128]uint32    // position in font array, from char value
        var fgctmp [128]uint32    // foreground color cache (rgba) for one line
        var bgctmp [128]uint32    // background color cache (rgba) for one line
        var dsttmp [128]uint32    // position in destination memory array

        cursor_state =        v.Mem[ VKY_TXT_CURSOR_CTRL_REG ]
        cursor_x     = uint32(v.Mem[ VKY_TXT_CURSOR_X_REG_L  ])
        cursor_y     = uint32(v.Mem[ VKY_TXT_CURSOR_Y_REG_L  ])
        
        if (v.Mem[ MASTER_CTRL_REG_L ] & 0x02) == 0x02 {
                is_overlay = true
        } else {
                is_overlay = false
        }

        // render text - start
        // I prefer to keep it because it allow to simply re-drawing single line in future,
        // by manupipulating starting point (now 0) and end clause (now <v.text_rows)
        fb_row_pos = v.starting_fb_row_pos
        for text_y = 0; text_y < v.text_rows; text_y += 1 { // over lines of text
                text_row_pos = text_y * v.text_cols
                for text_x = 0; text_x < v.text_cols; text_x += 1 { // pre-calculate data for x-axis
                        fnttmp[text_x] = v.text[text_row_pos+text_x] * 64 // position in font array
                        dsttmp[text_x] = text_x * 8                     // position of char in dest FB

                        f := v.fg[text_row_pos+text_x] // fg and bg colors
                        b := v.bg[text_row_pos+text_x]

                        if v.c.Cursor_visible && (cursor_y == text_y) && (cursor_x == text_x) && (cursor_state & 0x01 == 1) {
                                f = uint32((v.Mem[ VKY_TXT_CURSOR_COLR_REG ] & 0xf0) >> 4)
                                b = uint32( v.Mem[ VKY_TXT_CURSOR_COLR_REG ] & 0x0f)
                                fnttmp[text_x] = uint32(v.Mem[ VKY_TXT_CURSOR_CHAR_REG ]) * 64
                        }

                        fgctmp[text_x] = binary.LittleEndian.Uint32(f_color_lut[f][:]) // text LUT - xxx: change name
                        if is_overlay == false {
                                bgctmp[text_x] = binary.LittleEndian.Uint32(b_color_lut[b][:]) // text LUT
                        } else {
                                bgctmp[text_x] = 0x00FFFFFF                     // full alpha
                        }
                }

                for font_line = 0; font_line < 8; font_line += 1 { // for every line of text - over 8 lines of font
                        font_row_pos = font_line * 8
                        for text_x = 0; text_x < v.text_cols; text_x += 1 { // for each line iterate over columns of text
                                font_pos = fnttmp[text_x] + font_row_pos
	                        fb_pos   = dsttmp[text_x] + fb_row_pos
                                for i = 0; i < 8; i += 1 { // for every font iterate over 8 pixels of font
                                        if v.font[font_pos+i] == 0 {
                                                v.c.TFB[fb_pos+i] = bgctmp[text_x]
                                        } else {
                                                v.c.TFB[fb_pos+i] = fgctmp[text_x]
                                        }
                                }
                        }
                        fb_row_pos += v.x_res
                }
        }
        // render text - end
}

// RAM-interface specific
func (v *Vicky) Dump(address uint32) []byte {
        log.Panicf("vicky2 Dump is not implemented yet")
        return []byte{}
}

func (v *Vicky) Name(fn byte) string {
	switch fn {
	case F_MAIN:
		return v.name
	case F_TEXT:
		return v.name + "-text"
	case F_TEXT_C:
		return v.name + "-text_color"
	case F_VRAM:
		return v.name + "-vram"
	}
	return v.name + "-UNKNOWN"
}

func (v *Vicky) Clear() { 
        log.Panicf("vicky2 Clear is not implemented yet")
}

func (v *Vicky) Size(fn byte) (uint32, uint32) {
	switch fn {
	case F_MAIN:
		return uint32(1), uint32(len(v.Mem))
	case F_TEXT:
		return uint32(1), uint32(len(v.text))
	case F_TEXT_C:
		return uint32(1), uint32(len(v.tc))
	case F_VRAM:
		return uint32(1), uint32(len(v.vram))
	}
	return 0, 0
}

func (v *Vicky) Read(fn byte, addr uint32) (byte, error) {
	switch fn {
	case F_MAIN:
		return v.ReadReg(addr)
	case F_TEXT:
		return byte(v.text[addr]), nil
	case F_TEXT_C:
		return v.tc[addr], nil
	case F_VRAM:
		return v.vram[addr], nil
	}
	return 0, fmt.Errorf(" vicky2: %s Read addr %6X fn %d is not implemented", v.name, addr, fn)
}

func (v *Vicky) Write(fn byte, addr uint32, val byte) (error) {
	switch fn {
	case F_MAIN:
		return v.WriteReg(addr, val)
	case F_TEXT:
		v.text[addr] = uint32(val)
	case F_TEXT_C:
                bgc   := uint32( val & 0x0F)
                fgc   := uint32((val & 0xF0)>> 4)
                v.fg[addr] = fgc
                v.bg[addr] = bgc
		v.tc[addr] = val
	case F_VRAM:
		v.vram[addr] = val
		v.UpdateBitmapFB(addr, val)
	default:
		return fmt.Errorf(" vicky2: %s Write addr %6X val %2X fn %d is not implemented", v.name, addr, val)
	}
	return nil

}

// addr is relative here, ie. $B0:1000 means $1000
func (v *Vicky) UpdateBitmapFB(addr uint32, val byte) {
	if addr >= v.bm0_start_addr && addr - v.bm0_start_addr < uint32(len(v.c.BM0FB)) { 
		dst := addr - v.bm0_start_addr
		//fmt.Printf("bm0fb addr: %6X dst: %6X val %2X blut %4X\n", addr, dst, val, v.blut[v.bm0_blut_pos + uint32(val)])
		v.c.BM0FB[dst] = v.blut[v.bm0_blut_pos + uint32(val)]
	}
	if addr >= v.bm1_start_addr && addr - v.bm1_start_addr < uint32(len(v.c.BM1FB)) { 
		dst := addr - v.bm1_start_addr
		//fmt.Printf("bm1fb addr: %6X dst: %6X val %2X blut %4X\n", addr, dst, val, v.blut[v.bm1_blut_pos + uint32(val)])
		v.c.BM1FB[dst] = v.blut[v.bm1_blut_pos + uint32(val)]
	}
}

func (v *Vicky) ReadVram(addr uint32) (byte, error) {
	return v.vram[addr], nil
}


func (v *Vicky) ReadReg(addr uint32) (byte, error) {
        //fmt.Printf("vicky2: %s Read addr %06x\n", v.name, addr)
        switch addr {
        //case 0x0001:
        //        return 0x00, nil        // 640x480, no pixel doubling

        case 0x0002:
		if emu.DIP[6] {
			return 0x10, nil        // 1 = Hi-Res on BOOT OFF
		} else {
			return 0x00, nil        // 0 = Hi-Res on BOOT ON
		}

        case 0x070B:			// model major
                return 0x00, nil

        case 0x070C:			// model minor
                return 0x00, nil

        case 0xe902:			// CODEC_WR_CTRL - dummy value
                return 0x00, nil

	case 0xe80e:			// DIP_BOOTMODE - dummy value, XXX
		return 0x03, nil	// boot to BASIC

        default:
                return v.Mem[addr], nil
        }
}

func (v *Vicky) WriteReg(addr uint32, val byte) error {
        //fmt.Printf("vicky2: %s Write addr %06x val %02x\n", v.name, addr, val)
        v.Mem[addr] = val

        switch {
        case addr == MASTER_CTRL_REG_L:
                v.c.Master_L = val

        case addr == MASTER_CTRL_REG_H:
                v.c.Master_H = val
		if val & 0x01 == 0 {
			v.x_res = 640
			v.y_res = 480
		} else {
			v.x_res = 800
			v.y_res = 600
		}
		v.recalculateScreen()

        case addr == BORDER_CTRL_REG:
                if (val & 0x01) == 1 {
                        v.c.Border_x_size  = int32(v.Mem[ BORDER_X_SIZE ])
                        v.c.Border_y_size  = int32(v.Mem[ BORDER_Y_SIZE ])
                        v.c.Border_visible = true
                        v.recalculateScreen()
                } else {
                        v.c.Border_x_size  = 0
                        v.c.Border_y_size  = 0
                        v.recalculateScreen()
                        v.c.Border_visible = false
                }

        case addr == BORDER_COLOR_B:
                v.c.Border_color_b = val
                v.recalculateScreen()

        case addr == BORDER_COLOR_G:
                v.c.Border_color_g = val
                v.recalculateScreen()

        case addr == BORDER_COLOR_R:
                v.c.Border_color_r = val
                v.recalculateScreen()

        case addr == BORDER_X_SIZE:
                v.c.Border_x_size = int32(val & 0x3F)     // XXX: in spec - 0-32, bitmask allows to 0-63
                if v.c.Border_visible {
                        v.recalculateScreen()
                }

        case addr == BORDER_Y_SIZE:
                v.c.Border_y_size = int32(val & 0x3F)     // XXX: in spec - 0-32, bitmask allows to 0-63
                if v.c.Border_visible {
                        v.recalculateScreen()
                }

	case addr == VKY_TXT_CURSOR_X_REG_L:
	case addr == VKY_TXT_CURSOR_X_REG_H:
	case addr == VKY_TXT_CURSOR_Y_REG_L:
	case addr == VKY_TXT_CURSOR_Y_REG_H:
		break					// just write to mem

        case addr == BACKGROUND_COLOR_B:
                v.c.Background[2] = val

        case addr == BACKGROUND_COLOR_G:
                v.c.Background[1] = val

        case addr == BACKGROUND_COLOR_R:
                v.c.Background[0] = val

        case addr == VKY_TXT_CURSOR_CTRL_REG:
                if (val & 0x01) == 0 {
                        v.c.Cursor_visible = false
                } else {
                        v.c.Cursor_visible = true
                }
                
        case addr == BM0_CONTROL_REG:
                if (val & 0x01) == 0 {
                        v.c.BM0_visible = false
                } else {
                        v.c.BM0_visible = true
                }
                val = (val & 0x0E) >> 1                         // extract LUT number
                v.bm0_blut_pos = uint32(val) * 0x100            // position in Bitmap LUT cache
                // XXX - todo: repopulate LUT if pos changed

        
        case addr == BM0_START_ADDY_L:
        case addr == BM0_START_ADDY_M:
        case addr == BM0_START_ADDY_H:
                v.bm0_start_addr = (uint32(v.Mem[ BM0_START_ADDY_H ]) << 16) +
                                   (uint32(v.Mem[ BM0_START_ADDY_M ]) << 8 ) +
                                   (uint32(v.Mem[ BM0_START_ADDY_L ])      )

        case addr == BM1_CONTROL_REG:
                if (val & 0x01) == 0 {
                        v.c.BM1_visible = false
                } else {
                        v.c.BM1_visible = true
                }
                val = (val & 0x0E) >> 1                         // extract LUT number
                v.bm1_blut_pos = uint32(val) * 0x100            // position in Bitmap LUT cache
                // XXX - todo: repopulate LUT if pos changed


        case addr == BM1_START_ADDY_L:
        case addr == BM1_START_ADDY_M:
        case addr == BM1_START_ADDY_H:
                v.bm0_start_addr = (uint32(v.Mem[ BM1_START_ADDY_H ]) << 16) +
                                   (uint32(v.Mem[ BM1_START_ADDY_M ]) << 8 ) +
                                   (uint32(v.Mem[ BM1_START_ADDY_L ])      )

        case addr >= FG_CHAR_LUT_PTR && addr < FG_CHAR_LUT_PTR + 64:
                a := addr - FG_CHAR_LUT_PTR
                byte_in_lut := byte(a & 0x03)
                num := byte(a >> 2)
                f_color_lut[num][byte_in_lut] = val // XXX - global one!


        case addr >= BG_CHAR_LUT_PTR && addr < BG_CHAR_LUT_PTR + 64:
                a := addr - BG_CHAR_LUT_PTR
                byte_in_lut := byte(a & 0x03)
                num := byte(a >> 2)
                b_color_lut[num][byte_in_lut] = val // XXX - global one!

        // XXX - probably this needs correction with different
        //       bitmap format than ARGB
        case addr >= GRPH_LUT0_PTR  && addr < GRPH_LUT7_PTR + 0x400:
                src :=  addr & 0xfffffffc
                dst := (addr - GRPH_LUT0_PTR) >> 2              // clear bits 0-1, we need 4 bytes for in mem BGRA
                                                                // in memory representation fo uint32: ARGB
                if (dst & 0xff) == 0 {
                        v.blut[dst] = 0x00FFFFFF                  // LUTx[0] is always transparent, by design

                } else {
                        v.blut[dst] = binary.LittleEndian.Uint32(
                                        []byte{v.Mem[src], 
                                               v.Mem[src+1], 
                                               v.Mem[src+2], 
                                               v.Mem[src+3],
                                        })
                }
                //fmt.Printf("addr: %6x val %2x src %4x dst: %4d pix: %08x ram: %v\n", addr, val, src, dst, v.blut[dst], v.Mem[src:src+4])

        case addr >= FONT_MEMORY_BANK0 && addr < FONT_MEMORY_BANK0 + 0x800:
                v.updateFontCache(addr - FONT_MEMORY_BANK0, val)  // every bit in font cache is mapped to byte

        default:
                return fmt.Errorf(" vicky2: %s Write addr %6X val %2X is not implemented", v.name, addr, val)
        }
	return nil
}

