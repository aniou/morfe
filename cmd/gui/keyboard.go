// from Basic/Scancodes.cs (FoenixIDE)

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	sc_null            = 0x00
	sc_escape          = 0x01
	sc_1               = 0x02
	sc_2               = 0x03
	sc_3               = 0x04
	sc_4               = 0x05
	sc_5               = 0x06
	sc_6               = 0x07
	sc_7               = 0x08
	sc_8               = 0x09
	sc_9               = 0x0A
	sc_0               = 0x0B
	sc_minus           = 0x0C
	sc_equals          = 0x0D
	sc_backspace       = 0x0E
	sc_tab             = 0x0F
	sc_q               = 0x10
	sc_w               = 0x11
	sc_e               = 0x12
	sc_r               = 0x13
	sc_t               = 0x14
	sc_y               = 0x15
	sc_u               = 0x16
	sc_i               = 0x17
	sc_o               = 0x18
	sc_p               = 0x19
	sc_bracketLeft     = 0x1A
	sc_bracketRight    = 0x1B
	sc_enter           = 0x1C
	sc_controlLeft     = 0x1D
	sc_a               = 0x1E
	sc_s               = 0x1F
	sc_d               = 0x20
	sc_f               = 0x21
	sc_g               = 0x22
	sc_h               = 0x23
	sc_j               = 0x24
	sc_k               = 0x25
	sc_l               = 0x26
	sc_semicolon       = 0x27
	sc_apostrophe      = 0x28
	sc_grave           = 0x29
	sc_shiftLeft       = 0x2A
	sc_backslash       = 0x2B
	sc_z               = 0x2C
	sc_x               = 0x2D
	sc_c               = 0x2E
	sc_v               = 0x2F
	sc_b               = 0x30
	sc_n               = 0x31
	sc_m               = 0x32
	sc_comma           = 0x33
	sc_period          = 0x34
	sc_slash           = 0x35
	sc_shiftRight      = 0x36
	sc_numpad_multiply = 0x37
	sc_altLeft         = 0x38
	sc_space           = 0x39
	sc_capslock        = 0x3A
	sc_F1              = 0x3B
	sc_F2              = 0x3C
	sc_F3              = 0x3D
	sc_F4              = 0x3E
	sc_F5              = 0x3F
	sc_F6              = 0x40
	sc_F7              = 0x41
	sc_F8              = 0x42
	sc_F9              = 0x43
	sc_F10             = 0x44
	sc_F11             = 0x57
	sc_F12             = 0x58
	sc_up_arrow        = 0x48 // also maps to num keypad 8
	sc_left_arrow      = 0x4B // also maps to num keypad 4
	sc_right_arrow     = 0x4D // also maps to num keypad 6
	sc_down_arrow      = 0x50 // also maps to num keypad 2
)

func PS2ScanCode(code sdl.Scancode) byte {
	switch code {
	case sdl.SCANCODE_0:
		return sc_0
	case sdl.SCANCODE_1,
		sdl.SCANCODE_2,
		sdl.SCANCODE_3,
		sdl.SCANCODE_4,
		sdl.SCANCODE_5,
		sdl.SCANCODE_6,
		sdl.SCANCODE_7,
		sdl.SCANCODE_8,
		sdl.SCANCODE_9:
		return byte(sc_1 + (code - sdl.SCANCODE_1))
	case sdl.SCANCODE_A:
		return sc_a
	case sdl.SCANCODE_B:
		return sc_b
	case sdl.SCANCODE_C:
		return sc_c
	case sdl.SCANCODE_D:
		return sc_d
	case sdl.SCANCODE_E:
		return sc_e
	case sdl.SCANCODE_F:
		return sc_f
	case sdl.SCANCODE_G:
		return sc_g
	case sdl.SCANCODE_H:
		return sc_h
	case sdl.SCANCODE_I:
		return sc_i
	case sdl.SCANCODE_J:
		return sc_j
	case sdl.SCANCODE_K:
		return sc_k
	case sdl.SCANCODE_L:
		return sc_l
	case sdl.SCANCODE_M:
		return sc_m
	case sdl.SCANCODE_N:
		return sc_n
	case sdl.SCANCODE_O:
		return sc_o
	case sdl.SCANCODE_P:
		return sc_p
	case sdl.SCANCODE_Q:
		return sc_q
	case sdl.SCANCODE_R:
		return sc_r
	case sdl.SCANCODE_S:
		return sc_s
	case sdl.SCANCODE_T:
		return sc_t
	case sdl.SCANCODE_U:
		return sc_u
	case sdl.SCANCODE_V:
		return sc_v
	case sdl.SCANCODE_W:
		return sc_w
	case sdl.SCANCODE_X:
		return sc_x
	case sdl.SCANCODE_Y:
		return sc_y
	case sdl.SCANCODE_Z:
		return sc_z
	case sdl.SCANCODE_RETURN:
		return sc_enter
	case sdl.SCANCODE_DELETE, sdl.SCANCODE_BACKSPACE:
		return sc_backspace
	case sdl.SCANCODE_SPACE:
		return sc_space
	case sdl.SCANCODE_COMMA:
		return sc_comma
	case sdl.SCANCODE_PERIOD:
		return sc_period
	case sdl.SCANCODE_SEMICOLON:
		return sc_semicolon
	case sdl.SCANCODE_ESCAPE:
		return sc_escape
	case sdl.SCANCODE_GRAVE:
		return sc_grave
	case sdl.SCANCODE_APOSTROPHE:
		return sc_apostrophe
	case sdl.SCANCODE_LEFTBRACKET:
		return sc_bracketLeft
	case sdl.SCANCODE_RIGHTBRACKET:
		return sc_bracketRight
	case sdl.SCANCODE_MINUS:
		return sc_minus
	case sdl.SCANCODE_EQUALS:
		return sc_equals
	case sdl.SCANCODE_TAB:
		return sc_tab
	case sdl.SCANCODE_SLASH:
		return sc_slash
	case sdl.SCANCODE_BACKSLASH:
		return sc_backslash
	case sdl.SCANCODE_LSHIFT:
		return sc_shiftLeft
	case sdl.SCANCODE_RSHIFT:
		return sc_shiftRight
	case sdl.SCANCODE_LALT:
		return sc_altLeft
	case sdl.SCANCODE_LCTRL:
		return sc_controlLeft
	case sdl.SCANCODE_UP:
		return sc_up_arrow
	case sdl.SCANCODE_DOWN:
		return sc_down_arrow
	case sdl.SCANCODE_LEFT:
		return sc_left_arrow
	case sdl.SCANCODE_RIGHT:
		return sc_right_arrow
	case sdl.SCANCODE_F1,
		sdl.SCANCODE_F2,
		sdl.SCANCODE_F3,
		sdl.SCANCODE_F4,
		sdl.SCANCODE_F5,
		sdl.SCANCODE_F6,
		sdl.SCANCODE_F7,
		sdl.SCANCODE_F8,
		sdl.SCANCODE_F9,
		sdl.SCANCODE_F10:
		return byte(sc_F1 + (code - sdl.SCANCODE_F1))
	case sdl.SCANCODE_F11,
		sdl.SCANCODE_F12:
		return byte(sc_F11 + (code - sdl.SCANCODE_F11))
	default:
		return sc_null
	}
}

func (g *GUI) sendKey(code sdl.Scancode, state byte) {
	mask := g.p.Bus.EaRead(INT_MASK_REG1)
	if (^mask & byte(r1_FNX1_INT00_KBD)) == byte(r1_FNX1_INT00_KBD) {
		code := PS2ScanCode(code)
		//fmt.Printf("\nKEY pressed?:%v, mask %2X %2X %2X\n", state, mask, ^mask, byte(r1_FNX1_INT00_KBD))

		if code == sc_null {
			fmt.Printf("unknown scancode\n")
		} else {
			if state == sdl.RELEASED {
				g.p.GABE.Data = code + 0x80
			} else {
				g.p.GABE.Data = code
			}
			g.p.Bus.EaWrite(0xAF_1064, 0)
			irq1 := g.p.Bus.EaRead(INT_PENDING_REG1) | r1_FNX1_INT00_KBD
			g.p.Bus.EaWrite(INT_PENDING_REG1, irq1)
			g.p.CPU.TriggerIRQ()
		}
	}
}
