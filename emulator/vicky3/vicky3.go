package vicky3

// a foundation for Vicky3
// a big-endian implementation, suitable for A2560* platforms 
// and - probably - GenX ones

// XXX big warning XXX
// at this moment this is a verbatim copy of little-endian vicky2 code
// i.e. a little-endian one. converting to new register layout is on way

import (
        "encoding/binary"
        "fmt"
        "log"
        _ "sync"
        _ "github.com/aniou/morfe/lib/mylog"
        "github.com/aniou/morfe/emulator"
        _ "github.com/aniou/morfe/emulator/ram"
)

const (
	_ = iota
	F_MAIN
	F_TEXT
	F_TEXT_C
	F_VRAM
	F_CRAM
)

const MasterControlReg_A          = 0x0000              // uint32_t 0x00C40000 

const MasterControlReg_A_b4       = 0x0003
const VKY3_MCR_TEXT_EN            = 0x00_00_00_01       /* Text Mode Enable */
const VKY3_MCR_TEXT_OVRLY         = 0x00_00_00_02       /* Text Mode overlay */
const VKY3_MCR_GRAPH_EN           = 0x00_00_00_04       /* Graphic Mode Enable */

const MasterControlReg_A_b3       = 0x0002
const VKY3_MCR_RESOLUTION_MASK    = 0x00_00_03_00 >> 8  /* Resolution - 00: 640x480, 01:800x600, 10: 1024x768, 11: 640x400 */
const VKY3_MCR_DOUBLE_EN          = 0x00_00_04_00 >> 8  /* Doubling Pixel */

const MasterControlReg_A_b2       = 0x0001
const VKY3_MCR_GAMMA_EN           = 0x00_01_00_00 >> 16 /* GAMMA Enable */
const VKY3_MCR_MANUAL_GAMMA_EN    = 0x00_02_00_00 >> 16 /* Enable Manual GAMMA Enable */
const VKY3_MCR_BLANK_EN           = 0x00_04_00_00 >> 16 /* Turn OFF sync (to monitor in sleep mode) */

const MasterControlReg_A_b1       = 0x0000


const BorderControlReg_L_A        = 0x0004         // uint32_t 0x00C40004
const BORDER_CONTROL		  = 0x0007
const BORDER_X_SIZE               = 0x0006
const BORDER_Y_SIZE               = 0x0005

const VKY3_BRDR_EN                = 0x00_00_00_01       /* Border Enable */
const VKY3_X_SCROLL_MASK          = 0x00_00_00_70       /* X Scroll */
const VKY3_X_SIZE_MASK            = 0x00_00_3f_00 >> 8  /* X Size */
const VKY3_Y_SIZE_MASK            = 0x00_3f_00_00 >> 16 /* Y Size */

const BorderControlReg_H_A        = 0x0008	// uint32_t 0x00C40008
                              //  = 0x0008      // ?? alpha?
const BORDER_COLOR_R              = 0x0009      // guessing
const BORDER_COLOR_G              = 0x000a      // guessing
const BORDER_COLOR_B              = 0x000b      // guessing

const BackGroundControlReg_A      = 0x000c	// uint32_t 0x00C4000C

const CursorControlReg_L_A        = 0x0010	// uint32_t 0x00C40010  - cursor settings : co_ca_ra_en
const CURSOR_COLOR                = 0x0010      // co - color            0-255
const CURSOR_CHARACTER            = 0x0011      // ca - cursor character 0-255
const CURSOR_RATE                 = 0x0012      // ra - rate   & 0x02    (0 = 1 per sec, 1 = 2 per sec, 2 = 4 per sec, 3 = 5 per sec)
const CURSOR_ENABLE               = 0x0013      // en - enable & 0x01

const CursorControlReg_H_A        = 0x0014	// uint32_t 0x00C40014  - cursor position
const CURSOR_Y_H                  = 0x0014
const CURSOR_Y_L                  = 0x0015
const CURSOR_X_H                  = 0x0016
const CURSOR_X_L                  = 0x0017

const LineInterrupt0_A            = 0x0018      // uint16_t 0x00C40018
const LineInterrupt1_A            = 0x001a      // uint16_t 0x00C4001A
const LineInterrupt2_A            = 0x001c      // uint16_t 0x00C4001C
const LineInterrupt3_A            = 0x001e      // uint16_t 0x00C4001E

const MousePointer_Mem_A          = 0x0400      // uint16_t 0x00C40400
const MousePtr_A_CTRL_Reg         = 0x0c00      // uint16_t 0x00C40C00

const MousePtr_A_X_Pos            = 0x0c02      // uint16_t 0x00C40C02
const MousePtr_A_Y_Pos            = 0x0c04      // uint16_t 0x00C40C04
const MousePtr_A_Mouse0           = 0x0c0a      // uint16_t 0x00C40C0A
const MousePtr_A_Mouse1           = 0x0c0c      // uint16_t 0x00C40C0C
const MousePtr_A_Mouse2           = 0x0c0e      // uint16_t 0x00C40C0E

const FONT_MEMORY_BANK0           = 0x8000

/* implemented as separate "functions" in Vicky3 module, see F_* const */

//const ScreenText_A				  // char 0x00C60000       /* Text matrix */
//const ColorText_A                               // uint8_t 0x00C68000    /* Color matrix */
//const FG_CLUT_A                                 // uint16_t 0x00C6C400   /* Foreground LUT */
//const BG_CLUT_A                                 // uint16_t 0x00C6C440   /* Background LUT */

type Vicky struct {
        name    string         // id of instance

        mem     []byte         // general Vicky memory
        text    []uint32       // text memory
        vram    []byte         // VRAM
        tc      []byte         // text color memory
        blut    []uint32       // bitmap LUT cache : 256 colors * 8 banks (lut0 to lut7)
        fg      []uint32       // text foreground LUT cache
        bg      []uint32       // text background LUT cache
        font    []byte         // font cache       : 256 chars  * 8 lines * 8 columns
	cram    []byte	       // XXX - temporary ram for FG clut/BG clut and others

	fg_clut [16]uint32     // 16 pre-calculated RGBA colors for text fore- 
	bg_clut [16]uint32     // ...and background

	overlay_enabled bool

        c       emu.GPU_common  // common GPU-exported properties and framebuffers

        starting_fb_row_pos uint32
        text_cols       uint32
        text_rows       uint32

        bm0_blut_pos    uint32
        bm1_blut_pos    uint32
        bm0_start_addr  uint32
        bm1_start_addr  uint32

	pixel_size      uint32          // 1 for normal, 2 for double - XXX: not used yet
	resolution      byte            // for tracking resolution changes

}

func New(name string, size int) *Vicky {
        v       := Vicky{name: name}

        v.mem    = make([]byte,   size)             // main memory - TODO: shrink it!

        v.text   = make([]uint32,    0x4000)        // text memory  - 0x4000 in GenX
        v.vram   = make([]byte  , 0x40_0000)        // 4MB - TODO: settable
        v.tc     = make([]byte,      0x4000)        // text color memory - 0x4000 in GenX
        v.fg     = make([]uint32,    0x4000)        // foreground cache -  0x4000 in GenX
        v.bg     = make([]uint32,    0x4000)        // background color cache - 0x4000 in GenX
	v.cram   = make([]byte,        0x100)       // 'misc' color ram

        v.blut   = make([]uint32, 0x0800)           // bitmap LUT cache : 256 colors * 8 banks (lut0 to lut7)
        v.font   = make([]byte  , 0x100 * 8 * 8)    // font cache 256 chars * 8 lines * 8 columns
        v.c.TFB    = make([]uint32,  0x0c_0000)     // text framebuffer - 1024x768 max
        v.c.BM0FB  = make([]uint32,  0x0c_0000)     // bm0  framebuffer - 1024x768 max
        v.c.BM1FB  = make([]uint32,  0x0c_0000)     // bm1  framebuffer - 1024x768 max
        
	//
	v.c.Screen_x_size  = 640
	v.c.Screen_y_size  = 480
	v.c.Screen_resized = false

	v.resolution    = 0
        v.pixel_size    = 1

	// 'verified' knobs
	v.c.Bitmap_enabled = true // XXX: there is no way to change it in vicky3?

	// XXX: test
	v.c.Text_enabled	= true

        //v.mem[ BORDER_CTRL_REG ] = 0x01 - XXX - initial state?

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
        v.bm0_start_addr = 0x00 // relative from beginning of VRAM
        v.bm1_start_addr = 0x00 // relative from beginning of VRAM


        // XXX - just for test
        for i := range v.text { // file text memory areas
              v.fg[i] = 0x0d
              v.bg[i] = 0x02
              v.text[i] = 0x20
        } 

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
	if v.c.Border_enabled {
		v.starting_fb_row_pos = uint32(v.c.Screen_x_size) * uint32(v.c.Border_y_size) + uint32(v.c.Border_x_size)
	} else {
		v.starting_fb_row_pos = 0
	}

        //v.text_cols = (640 - (uint32(v.Border_x_size) * 2)) / 8 // xxx - parametrize screen width
        //v.text_rows = (480 - (uint32(v.Border_y_size) * 2)) / 8 // xxx - parametrize screen height
        //v.text_cols = (v.c.Screen_x_size - (uint32(v.Border_x_size) * 2)) / (v.pixel_size * 8)
        //v.text_rows = (v.c.Screen_y_size - (uint32(v.Border_y_size) * 2)) / (v.pixel_size * 8)
        v.text_cols = uint32(v.c.Screen_x_size / 8)
        v.text_rows = uint32(v.c.Screen_y_size / 8)

	fmt.Printf("vicky3: text_rows: %d\n", v.text_rows)
	fmt.Printf("vicky3: text_cols: %d\n", v.text_cols)

}

// RAM-interface specific
func (v *Vicky) Dump(address uint32) []byte {
        log.Panicf("vicky3 Dump is not implemented yet")
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
        case F_CRAM:
                return v.name + "-cram"
        }
        return v.name + "-UNKNOWN"
}

func (v *Vicky) Clear() { 
        log.Panicf("vicky3 Clear is not implemented yet")
}

func (v *Vicky) Size(fn byte) (uint32, uint32) {
        switch fn {
        case F_MAIN:
                return uint32(1), uint32(len(v.mem))
        case F_TEXT:
                return uint32(1), uint32(len(v.text))
        case F_TEXT_C:
                return uint32(1), uint32(len(v.tc))
        case F_VRAM:
                return uint32(1), uint32(len(v.vram))
        case F_CRAM:
                return uint32(1), uint32(len(v.cram))
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
        case F_CRAM:
                return v.cram[addr], nil
        }
        return 0, fmt.Errorf(" vicky3: %s Read addr %6X fn %d is not implemented", v.name, addr, fn)
}

func (v *Vicky) Write(fn byte, addr uint32, val byte) (error) {
        //fmt.Printf("vicky3: %s Write func %02x addr %06x val %02x\n", v.name, fn, addr, val)
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
	case F_CRAM:
		v.cram[addr] = val

		// it is something strange here, like a mix of little and big endian:
		//  0xHHLL, 0xHHLL
		//  0xGGBB, 0xAARR
		//  const unsigned short fg_color_lut [32] = {
		//        0x0000, 0xFF00, // Black (transparent)
		//        0x0000, 0xFF80, // Mid-Tone Red
		//  ...
		//  - it looks like "middle-endian"

		switch {
		case addr >= 0x00 && addr < 0x40:
			color   := addr >> 2
			a       := addr & 0b_1111_1100

			tmp := append(v.cram[a+2:a+4], v.cram[a:a+2]...)
			v.fg_clut[color] = binary.BigEndian.Uint32( tmp )
			//fmt.Printf(" vicky3: FG color %2d %8x %v\n", color, v.fg_clut[color], v.cram[a:a+4] )

		case addr >= 0x40 && addr < 0x80:
			color   := (addr - 0x40) >> 2
			a       := addr & 0b_1111_1100

			tmp := append(v.cram[a+2:a+4], v.cram[a:a+2]...)
			v.bg_clut[color] = binary.BigEndian.Uint32( tmp )
		}

        default:
                return fmt.Errorf(" vicky3: %s Write addr %6X val %2X fn %d is not implemented", v.name, addr, val)
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

func (v *Vicky) ReadReg(addr uint32) (byte, error) {
        //fmt.Printf("vicky3: %s Read addr %06x\n", v.name, addr)
        switch addr {
        //case 0x0001:
        //        return 0x00, nil        // 640x480, no pixel doubling
	/*
        case 0x0002:
                if emu.DIP[6] {
                        return 0x10, nil        // 1 = Hi-Res on BOOT OFF
                } else {
                        return 0x00, nil        // 0 = Hi-Res on BOOT ON
                }

        case 0x070B:                    // model major
                return 0x00, nil

        case 0x070C:                    // model minor
                return 0x00, nil

        case 0xe902:                    // CODEC_WR_CTRL - dummy value
                return 0x00, nil

        case 0xe80e:                    // DIP_BOOTMODE - dummy value, XXX
                return 0x03, nil        // boot to BASIC
        */
        default:
                return v.mem[addr], nil
        }

}

func (v *Vicky) WriteReg(addr uint32, val byte) error {
        //fmt.Printf("vicky3: %s Write addr %06x val %02x\n", v.name, addr, val)
        v.mem[addr] = val

        switch {
        case addr == MasterControlReg_A_b1:
		return nil			// no functions, so do nothing

        //case addr == MasterControlReg_A_b2:
	//	return nil			// XXX: GAMMA and sync OFF not supported yet

        case addr == MasterControlReg_A_b3:		  // XXX: screen resolution and pixel doubling
                switch (val & VKY3_MCR_RESOLUTION_MASK) {
		case 0b_00:
			v.c.Screen_x_size  = 640
			v.c.Screen_y_size  = 480
		case 0b_01:
			v.c.Screen_x_size  = 800
			v.c.Screen_y_size  = 600
		case 0b_10:
			v.c.Screen_x_size  = 1024
			v.c.Screen_y_size  = 768
		case 0b_11:
			v.c.Screen_x_size  = 640
			v.c.Screen_y_size  = 400
		}

		if (val & VKY3_MCR_DOUBLE_EN) != 0 {
			v.pixel_size = 2
		} else {
			v.pixel_size = 1
		}

		if v.resolution != val {
			v.resolution = val & VKY3_MCR_RESOLUTION_MASK
			v.c.Screen_resized = true
		}

                v.recalculateScreen()
		return nil

	case addr == MasterControlReg_A_b4:		  // text mode, overlay and graphic mode 
		v.c.Text_enabled    = (val & 0x01) != 0
		v.overlay_enabled   = (val & 0x02) != 0
		v.c.Graphic_enabled = (val & 0x04) != 0
		return nil

	case addr == BorderControlReg_H_A:		// probably no meaning, do nothing
		return nil

        case addr == BORDER_CONTROL:
                v.c.Border_enabled  = (val & 0x01) != 0
                v.recalculateScreen()

		if (val & 0x70) != 0 {
			return fmt.Errorf(" vicky3: %s Write addr %6X val %2X is not implemented", v.name, addr, val)
		}

        case addr == BORDER_X_SIZE:
                v.c.Border_x_size  = int32( val & 0x3f )
		if v.c.Border_enabled {
                        v.recalculateScreen()
                }

        case addr == BORDER_Y_SIZE:
                v.c.Border_y_size  = int32( val & 0x3f )
		if v.c.Border_enabled {
                        v.recalculateScreen()
                }

        case addr == BORDER_COLOR_R:
                v.c.Border_color_r = val

        case addr == BORDER_COLOR_G:
                v.c.Border_color_g = val

        case addr == BORDER_COLOR_B:
                v.c.Border_color_b = val

        case addr >= FONT_MEMORY_BANK0 && addr < FONT_MEMORY_BANK0 + 0x800:
                v.updateFontCache(addr - FONT_MEMORY_BANK0, val)  // every bit in font cache is mapped to byte

	/*
        case addr == MASTER_CTRL_REG_H:
                v.c.Master_H = val
                if val & 0x01 == 0 {
                        v.c.Screen_x_size = 640
                        v.c.Screen_y_size = 480
                } else {
                        v.c.Screen_x_size = 800
                        v.c.Screen_y_size = 600
                }
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
                break                                   // just write to mem

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
                v.bm0_start_addr = (uint32(v.mem[ BM0_START_ADDY_H ]) << 16) +
                                   (uint32(v.mem[ BM0_START_ADDY_M ]) << 8 ) +
                                   (uint32(v.mem[ BM0_START_ADDY_L ])      )

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
                v.bm0_start_addr = (uint32(v.mem[ BM1_START_ADDY_H ]) << 16) +
                                   (uint32(v.mem[ BM1_START_ADDY_M ]) << 8 ) +
                                   (uint32(v.mem[ BM1_START_ADDY_L ])      )


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
                                        []byte{v.mem[src], 
                                               v.mem[src+1], 
                                               v.mem[src+2], 
                                               v.mem[src+3],
                                        })
                }
                //fmt.Printf("addr: %6x val %2x src %4x dst: %4d pix: %08x ram: %v\n", addr, val, src, dst, v.blut[dst], v.mem[src:src+4])

        case addr >= FONT_MEMORY_BANK0 && addr < FONT_MEMORY_BANK0 + 0x800:
                v.updateFontCache(addr - FONT_MEMORY_BANK0, val)  // every bit in font cache is mapped to byte

	*/
        default:
                return fmt.Errorf(" vicky3: %s Write addr %6X val %2X is not implemented", v.name, addr, val)
        }
        return nil
}

func (v *Vicky) RenderBitmapText() {
        var cursor_x, cursor_y uint32 // row and column of cursor
        var cursor_enabled bool     // cursor register, various states
        var text_x, text_y uint32 // row and column of text
        var text_row_pos uint32   // beginning of current text row in text memory
        var fb_row_pos uint32     // beginning of current FB   row in memory
        var font_pos uint32       // position in font array (char * 64 + char_line * 8)
        var   fb_pos uint32       // position in destination framebuffer
        var font_line uint32      // line in current font
        var font_row_pos uint32   // position of line in current font (=font_line*8 because every line has 8 bytes)
        var i uint32              // counter
        //var is_overlay bool       // is overlay text over bm0 enabled?

        // placeholders recalculated per row of text, holds values for text_cols loop
        // current max size is 100 columns (800/8)
        var fnttmp [128]uint32    // position in font array, from char value
        var fgctmp [128]uint32    // foreground color cache (rgba) for one line
        var bgctmp [128]uint32    // background color cache (rgba) for one line
        var dsttmp [128]uint32    // position in destination memory array

	// XXX: it should be rather updated on register write
        cursor_enabled =       (v.mem[ CURSOR_ENABLE ] & 0x01) == 0x01
        cursor_x       = uint32(v.mem[ CURSOR_X_H ]) << 16 | uint32(v.mem[ CURSOR_X_L ])
        cursor_y       = uint32(v.mem[ CURSOR_Y_H ]) << 16 | uint32(v.mem[ CURSOR_Y_L ])
        
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

                        if v.c.Cursor_visible && cursor_enabled && (cursor_y == text_y) && (cursor_x == text_x) {
                                f = uint32((v.mem[ CURSOR_COLOR ] & 0xf0) >> 4)
                                b = uint32( v.mem[ CURSOR_COLOR ] & 0x0f)
                                fnttmp[text_x] = uint32(v.mem[ CURSOR_CHARACTER ]) * 64
                        }

                        fgctmp[text_x] = v.fg_clut[f]
                        if v.overlay_enabled == false {
                                bgctmp[text_x] = v.bg_clut[b]
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
                        fb_row_pos += uint32(v.c.Screen_x_size)
                }
        }
        // render text - end
}

// eof
