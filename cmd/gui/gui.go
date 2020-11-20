package main

import (
	"encoding/binary"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	//"time"
	"github.com/aniou/go65c816/emulator/platform"
	"github.com/aniou/go65c816/lib/mylog"
)

const FULLSCREEN = false
const CPU_CLOCK = 14318000 // 14.381Mhz

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int32 = 640, 480

// global, for performance reasons
var fb []uint32
var font [256 * 8 * 8]byte // 256 chars * 8 lines * 8 columns

type VICKY struct {
	border_ctrl_reg byte
	border_color_b  byte
	border_color_g  byte
	border_color_r  byte
	border_x_size   uint32
	border_y_size   uint32
}

type DEBUG struct {
	gui	bool
}
var debug = DEBUG{true}

func (v *VICKY) FillByBorderColor() {
	val := binary.LittleEndian.Uint32([]byte{v.border_color_r, v.border_color_g, v.border_color_b, 0xff})
	fb[0] = val
	for bp := 1; bp < len(fb); bp *= 2 {
		copy(fb[bp:], fb[:bp])
	}
}

func (v *VICKY) SetBorderX(size byte) {
	v.border_x_size = uint32(size & 0xF8)
	v.FillByBorderColor()
}

func (v *VICKY) SetBorderY(size byte) {
	v.border_y_size = uint32(size & 0xF8)
	v.FillByBorderColor()
}

func showCPUSpeed(cycles uint64) (uint64, string) {
	switch {
	case cycles > 1000000:
		return cycles / 1000000, "MHz"
	case cycles > 1000:
		return cycles / 100, "kHz"
	default:
		return cycles, "Hz"
	}
}

func memoryDump(p *platform.Platform, address uint32) {
	var x uint16
	var a uint16

	for a = 0; a < 0x100; a = a + 16 {
		start, data := p.CPU.Bus.EaDump(address + uint32(a))
		bank := byte(start >> 16)
		addr := uint16(start)
		fmt.Printf("\n%02x:%04x│", bank, addr)
		if data != nil {
			fmt.Printf("% x│% x│", data[0:8], data[8:16])
			for x = 0; x < 16; x++ {
				if data[x] >= 33 && data[x] < 127 {
					fmt.Printf("%s", data[x:x+1])
				} else {
					fmt.Printf(".")
				}
				if x == 7 {
					fmt.Printf(" ")
				}
			}
		} else {
			fmt.Printf("                       │                       │")
		}
	}
}

func waitForEnter() {
	fmt.Println("\nPress the Enter Key")
	fmt.Scanln() // wait for Enter Key
}

func loadFont(fontset *[2048]uint32) {
	for i, v := range fontset {
		for j := 0; j < 8; j = j + 1 {
			v = v << 1
			if (v & 256) == 256 {
				//fmt.Printf("#")
				font[i*8+j] = 1
			} else {
				//fmt.Printf(" ")
				font[i*8+j] = 0
			}
		}
		//fmt.Printf("\n")
	}
}

// debug routines
func debugPixelFormat(window *sdl.Window) {
	pixelformat, err := window.GetPixelFormat()
	if err != nil {
		log.Fatalf("Failed to get pixel format: %s\n", err)
	}
	fmt.Printf("window pixel format: %s\n", sdl.GetPixelFormatName(uint(pixelformat)))
}

func debugRendererInfo(renderer *sdl.Renderer) {
	r_info, err := renderer.GetInfo()
	if err != nil {
		log.Fatalf("Failed to get renderer info: %s\n", err)
	}
	fmt.Printf("renderer: %s\n", r_info.Name)
	fmt.Printf("MaxTextureWidth: %d\n", r_info.MaxTextureWidth)
	fmt.Printf("MaxTextureHeighh: %d\n", r_info.MaxTextureHeight)
	for _, v := range r_info.TextureFormats {
		fmt.Printf("format: %s\n", sdl.GetPixelFormatName(uint(v)))
	}
	fmt.Printf("\n")
}

// TODO - parametryzacja okna
func newTexture(renderer *sdl.Renderer) *sdl.Texture {
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, 640, 480)
	if err != nil {
		log.Fatalf("Failed to create texture font from surface: %s\n", err)
	}
	
	if debug.gui {
		format, _, w, h, err := texture.Query()
		if err != nil {
			log.Fatalf("Failed to query texture: %s\n", err)
		}
		fmt.Printf("texture format: %s\n", sdl.GetPixelFormatName(uint(format)))
		fmt.Printf("texture width: %d\n", w)
		fmt.Printf("texture heigtt: %d\n", h)
	}

	return texture
}


func main() {
	var err error

	vicky := VICKY{}
	fb = make([]uint32, 640*480)

	//vicky.border_color_r = 0x20
	vicky.border_color_r = 0x00
	vicky.border_color_g = 0x00
	//vicky.border_color_b = 0x20
	vicky.border_color_b = 0x00
	vicky.SetBorderX(32)
	vicky.SetBorderY(32)

	var text [8192]uint32 // CS_TEXT_MEM
	var fg [8192]uint32   // foreground attributes
	var bg [8192]uint32   // background attributes

	pseudoInit()          // fill LUT table
	for i := range text { // file text memory areas
		fg[i] = 0x0e
		bg[i] = 0x0d
		text[i] = 32
	}

	// pre-defined font at start
	loadFont(&font_st_8x8)

	// test text
	count := 0
	for _, char := range "This is sparta!" {
		text[count] = uint32(char)
		count += 1
	}

	// platform init
	logger := mylog.New()
	p := platform.New()
	p.Init(logger)
	p.GPU.FB = &text
	p.GPU.FG = &fg
	p.GPU.BG = &bg
	p.GPU.FG_lut = &f_color_lut
	p.GPU.BG_lut = &b_color_lut
	//p.LoadHex("/home/aniou/c256/go65c816/data/matrix.hex")
	p.LoadHex("/home/aniou/c256/src/c256-gui-shim/old-kernel.hex")
	p.LoadHex("/home/aniou/c256/of816/platforms/C256/forth.hex")
	p.LoadHex("/home/aniou/c256/src/c256-gui-shim/c256-gui-shim.hex")
	p.CPU.PC = 0x0000
	p.CPU.RK = 0x03
	//memoryDump(p, 0x381000)
	//waitForEnter()


	// step 1: SDL
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Panic(sdl.GetError())
	}
	defer sdl.Quit()


	// step 2: Window
	var window *sdl.Window
	window, err = sdl.CreateWindow(
		winTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight,
		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL,
	)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	defer window.Destroy()
	debugPixelFormat(window)


	// step 3: Renderer
	var renderer *sdl.Renderer
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	defer renderer.Destroy()
	debugRendererInfo(renderer)

	// TODO - move it
	renderer.SetDrawColor(0, 255, 0, 255)
	renderer.Clear()

	var event sdl.Event
	var running bool
	// end of TODO

	// font texture
	texture := newTexture(renderer)

	// bit texture
	texture2 := newTexture(renderer)


	var prev_ticks uint32 = sdl.GetTicks()
	var ticks_now, frames uint32

	// -----------------------------------------------------------------------------------
	// zmiana trybu
	var current_mode sdl.DisplayMode
	if FULLSCREEN {
		var wanted_mode = sdl.DisplayMode{sdl.PIXELFORMAT_ARGB8888, 640, 480, 60, nil}
		var result_mode sdl.DisplayMode
		display_index, _ := window.GetDisplayIndex()
		current_mode, _ = sdl.GetCurrentDisplayMode(display_index)
		fmt.Printf("current mode width: %d\n", current_mode.W)
		fmt.Printf("current mode heigt: %d\n", current_mode.H)

		_, err = sdl.GetClosestDisplayMode(display_index, &wanted_mode, &result_mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get ClosesMode: %s\n", err)
			os.Exit(2)
		}
		fmt.Printf("wanted mode width: %d\n", result_mode.W)
		fmt.Printf("wanted mode heigtt: %d\n", result_mode.H)
		window.SetDisplayMode(&result_mode)
		window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	}
	// -----------------------------------------------------------------------------------

	var text_cols, text_rows uint32
	// text render
	text_cols = (640 - (vicky.border_x_size * 2)) / 8 // xxx - parametrize screen width
	text_rows = (480 - (vicky.border_y_size * 2)) / 8 // xxx - parametrize screen height
	fmt.Printf("text_rows: %d\n", text_rows)

	sdl.SetHint("SDL_HINT_RENDER_BATCHING", "1")

	var text_x, text_y uint32 // row and column of text
	var text_row_pos uint32   // beginning of current text row in text memory
	var fb_row_pos uint32     // beginning of current FB   row in memory
	var font_pos uint32       // position in font array (char * 64 + char_line * 8)
	var font_line uint32      // line in current font
	var font_row_pos uint32   // position of line in current font (=font_line*8 because every line has 8 bytes)
	var i uint32

	// placeholders recalculated per rows, holds values for text_cols loop
	var txttmp [128]uint32
	var fgtmp [128]uint32 // for rgba
	var bgtmp [128]uint32 // for rgba
	var dsttmp [128]uint32

	var prevCycles uint64 = 0
	var cpuSteps uint64 = 10000 // CPU steps, low initial
	//var l uint64

	// -----------------------------------------------------------------------------
	sdl.StartTextInput()

	starting_fb_row_pos := 640*vicky.border_y_size + (vicky.border_x_size)
	running = true
	for running {
		// render text - start
		fb_row_pos = starting_fb_row_pos
		for text_y = 0; text_y < text_rows; text_y += 1 { // over lines of text
			text_row_pos = text_y * 128
			for text_x = 0; text_x < text_cols; text_x += 1 { // pre-calculate data for x-axis
				txttmp[text_x] = text[text_row_pos+text_x] * 64 // position in font array
				dsttmp[text_x] = text_x * 8                     // position of char in dest FB

				f := fg[text_row_pos+text_x] // fg and bg colors
				b := bg[text_row_pos+text_x]
				fgtmp[text_x] = binary.LittleEndian.Uint32(f_color_lut[f][:])
				bgtmp[text_x] = binary.LittleEndian.Uint32(b_color_lut[b][:])

			}

			for font_line = 0; font_line < 8; font_line += 1 { // for every line of text - over 8 lines of font
				font_row_pos = font_line * 8
				for text_x = 0; text_x < text_cols; text_x += 1 { // for each line iterate over columns of text
					font_pos = txttmp[text_x] + font_row_pos
					for i = 0; i < 8; i += 1 { // for every font iterate over 8 pixels of font
						if font[font_pos+i] == 0 {
							fb[fb_row_pos+dsttmp[text_x]+i] = bgtmp[text_x]
						} else {
							fb[fb_row_pos+dsttmp[text_x]+i] = fgtmp[text_x]
						}
					}
				}
				fb_row_pos += 640
			}
		}
		// render text - end
		texture.UpdateRGBA(nil, fb, 640)
		//renderer.SetDrawColor(0x40, 0x00, 0x40, 255)
		//renderer.Clear()
		renderer.Copy(texture2, nil, nil)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		frames++
		ticks_now = sdl.GetTicks()
		if (ticks_now - prev_ticks) >= 1000 {
			if (p.CPU.AllCycles - prevCycles) < CPU_CLOCK {
				cpuSteps += 100
			}
			if (p.CPU.AllCycles - prevCycles) > CPU_CLOCK+10000 {
				cpuSteps -= 10
			}

			cyc, unit := showCPUSpeed(p.CPU.AllCycles - prevCycles)
			prevCycles = p.CPU.AllCycles
			fmt.Fprintf(os.Stdout, "keyq len: %d frames: %d ticks %d desired cycles %d cpu cycles %d speed %d %s cpu.K %02x cpu.PC %04x\n", p.Console.InBuf.Len(), frames, (ticks_now - prev_ticks), cpuSteps, p.CPU.AllCycles, cyc, unit, p.CPU.RK, p.CPU.PC)
			prev_ticks = ticks_now
			frames = 0
			//memoryDump(p, 0x0)
		}

		// keyboard ----------------------------------------------------------
		// https://github.com/veandco/go-sdl2-examples/blob/master/examples/keyboard-input/keyboard-input.go
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.TextInputEvent:
				fmt.Printf("TextInputEvent\n")
				for _, val := range t.Text {
					if val == 0 {
						break
					}
					p.Console.InBuf.Enqueue(val)
				}

			case *sdl.KeyboardEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)

				if t.State == sdl.PRESSED {
					if t.Repeat > 0 {
						continue
					}
					switch t.Keysym.Sym {
					case sdl.K_F12:
						running = false
					case sdl.K_F11:
						loadFont(&font_st_8x8)
					case sdl.K_F10:
						loadFont(&font_c256_8x8)
					case sdl.K_BACKSPACE,
						sdl.K_RETURN:
						p.Console.InBuf.Enqueue(byte(t.Keysym.Sym)) // XXX horrible, terrible

					default:
					}
				}
			}
		}

		// cpu step ----------------------------------------------------------
		// XXX: change it to regular steps and "stalled" steps in CPU routines
		/*
			for l = 0; l < cpuSteps; l += 1 {
				_, stopped := p.CPU.Step()
				if stopped {
					running = false
					break
				}
			}
		*/
		//cycles, stopped := p.CPU.Step()
		//fmt.Printf("CPU %d cycles and stopped %v\n", cycles, stopped)

	}
	if FULLSCREEN {
		window.SetDisplayMode(&current_mode)
	}

	//renderer.Destroy()
	//window.Destroy()
	//sdl.Quit()
}
