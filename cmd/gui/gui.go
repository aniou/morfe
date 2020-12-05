
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	//"time"
	"github.com/aniou/go65c816/emulator/platform"
	_ "github.com/aniou/go65c816/lib/mylog"
)

// keyboard memory registers
const INT_MASK_REG1    = 0x00_014D
const INT_PENDING_REG1 = 0x00_0141
// 


const FULLSCREEN = false
const CPU_CLOCK = 14318000 // 14.381Mhz
var   CPU_STEP uint64  = 14318

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int32 = 640, 480

// global, for performance reasons
var fb []uint32
var font [256 * 8 * 8]byte // 256 chars * 8 lines * 8 columns

type GUI struct {
	p *platform.Platform
}

type DEBUG struct {
	gui	bool
}
var debug = DEBUG{true}

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

func printCPUFlags(flag byte, name string) string {
        if flag > 0 {
                return name
        } else {
                return "-"
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
			fmt.Printf("		       │		       │")
		}
	}
	fmt.Printf("\n")
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
	fb = make([]uint32, 640*480)


	//pseudoInit()          // fill LUT table
	// pre-defined font at start
	loadFont(&font_st_8x8)

	// platform init
	p := platform.New()
	gui := GUI{p}

	p.InitGUI()
	p.GPU.FB   = &fb
	//p.LoadHex("/home/aniou/c256/go65c816/data/matrix.hex")

	/*
	p.LoadHex("/home/aniou/c256/src/c256-gui-shim/old-kernel.hex")
	//p.LoadHex("/home/aniou/c256/src/c256-gui-shim/c256-gui-shim.hex")
	p.CPU.PC = 0x0000
	p.CPU.RK = 0x03
	*/

	//memoryDump(p, 0x381000)
	//waitForEnter()
	//p.LoadHex("/home/aniou/c256/IDE/bin/Release/roms/kernel.hex")
	//p.LoadHex("/home/aniou/c256/of816/platforms/C256/forth.hex")
	p.LoadHex("/home/aniou/c256/Kernel_FMX.old/kernel.hex")
	p.LoadHex("/home/aniou/c256/src/c256-gui-shim/c256-gui-shim2.hex")
	p.LoadHex("/home/aniou/c256/of816/platforms/C256/forth.hex")
	//p.LoadHex("/home/aniou/c256/FoenixIDE-release-0.4.2.1/bin/Release/roms/kernel.hex")
	p.CPU.PC = 0xff00
	p.CPU.RK = 0x00
	p.CPU.PC = 0x0000
	p.CPU.RK = 0x03
	//memoryDump(p, 0x00fff0)
	//memoryDump(p, 0x00ff00)
	//waitForEnter()

	p.CPU.Bus.EaWrite(0xAF_0005, 0x20) // border B 
	p.CPU.Bus.EaWrite(0xAF_0006, 0x00) // border G
	p.CPU.Bus.EaWrite(0xAF_0007, 0x20) // border R

	p.CPU.Bus.EaWrite(0xAF_0008, 0x20) // border X
	p.CPU.Bus.EaWrite(0xAF_0009, 0x20) // border Y

	p.CPU.Bus.EaWrite(0xAF_0010, 0x03) // VKY_TXT_CURSOR_CTRL_REG
	p.CPU.Bus.EaWrite(0xAF_0012, 0xB1) // VKY_TXT_CURSOR_CHAR_REG
	p.CPU.Bus.EaWrite(0xAF_0013, 0xC4) // VKY_TXT_CURSOR_COLR_REG


	// test text
	for i := range p.GPU.TEXT { // file text memory areas
		p.GPU.FG[i] = 0x0e
		p.GPU.BG[i] = 0x0d
		p.GPU.TEXT[i] = 32
	}

	for c, _ := range p.GPU.TEXT {
		p.GPU.TEXT[c] = 0
	}
	for c, char := range "This is sparta!" {
		p.GPU.TEXT[c] = uint32(char)
	}
	// end of test text

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


	disasm := false

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


	// text render
	//var text_cols, text_rows uint32
	var text_cols uint32 = (640 - (p.GPU.Border_x_size * 2)) / 8 // xxx - parametrize screen width
	var text_rows uint32 = (480 - (p.GPU.Border_y_size * 2)) / 8 // xxx - parametrize screen height
	if debug.gui {
		fmt.Printf("text_rows: %d\n", text_rows)
		fmt.Printf("text_cols: %d\n", text_cols)
	}

	var cursor_x, cursor_y uint32 // row and column of cursor
	var cursor_char uint32    // cursor character
	var text_x, text_y uint32 // row and column of text
	var text_row_pos uint32   // beginning of current text row in text memory
	var fb_row_pos uint32     // beginning of current FB   row in memory
	var font_pos uint32       // position in font array (char * 64 + char_line * 8)
	var font_line uint32      // line in current font
	var font_row_pos uint32   // position of line in current font (=font_line*8 because every line has 8 bytes)
	var i uint32		  // counter

	// placeholders recalculated per row of text, holds values for text_cols loop --
	var fnttmp [128]uint32 	  // position in font array, from char value
	var fgctmp [128]uint32 	  // foreground color cache (rgba) for one line
	var bgctmp [128]uint32 	  // background color cache (rgba) for one line
	var dsttmp [128]uint32 	  // position in destination memory array

	// -----------------------------------------------------------------------------
	sdl.SetHint("SDL_HINT_RENDER_BATCHING", "1")
	sdl.StartTextInput()

	// variables for performance calculation ---------------------------------------
	var prev_ticks uint32 = sdl.GetTicks()
	var mult       uint32 = prev_ticks
	var ticks_now, frames uint32
	var stepCycles, prevCycles uint64 = 0, 0
	
	// main loop -------------------------------------------------------------------
	starting_fb_row_pos := 640*p.GPU.Border_y_size + (p.GPU.Border_x_size)
	running = true
	for running {
		cursor_char = uint32(p.CPU.Bus.EaRead(0xAF_0012))
		cursor_x    = uint32(p.CPU.Bus.EaRead(0xAF_0014))
		cursor_y    = uint32(p.CPU.Bus.EaRead(0xAF_0016))

		// render text - start
		fb_row_pos = starting_fb_row_pos
		for text_y = 0; text_y < text_rows; text_y += 1 { // over lines of text
			text_row_pos = text_y * 128
			for text_x = 0; text_x < text_cols; text_x += 1 { // pre-calculate data for x-axis
				fnttmp[text_x] = p.GPU.TEXT[text_row_pos+text_x] * 64 // position in font array
				dsttmp[text_x] = text_x * 8                     // position of char in dest FB

				f := p.GPU.FG[text_row_pos+text_x] // fg and bg colors
				b := p.GPU.BG[text_row_pos+text_x]

				if (cursor_y == text_y) && (cursor_x == text_x) {
					f = uint32((p.CPU.Bus.EaRead(0xAF_0013) & 0xf0) >> 4)
					b = uint32((p.CPU.Bus.EaRead(0xAF_0013) & 0x0f))
					fnttmp[text_x] = cursor_char * 64
				}

				fgctmp[text_x] = binary.LittleEndian.Uint32(p.GPU.FG_lut[f][:])
				bgctmp[text_x] = binary.LittleEndian.Uint32(p.GPU.BG_lut[b][:])
			}

			for font_line = 0; font_line < 8; font_line += 1 { // for every line of text - over 8 lines of font
				font_row_pos = font_line * 8
				for text_x = 0; text_x < text_cols; text_x += 1 { // for each line iterate over columns of text
					font_pos = fnttmp[text_x] + font_row_pos
					for i = 0; i < 8; i += 1 { // for every font iterate over 8 pixels of font
						if font[font_pos+i] == 0 {
							fb[fb_row_pos+dsttmp[text_x]+i] = bgctmp[text_x]
						} else {
							fb[fb_row_pos+dsttmp[text_x]+i] = fgctmp[text_x]
						}
					}
				}
				fb_row_pos += 640
			}
		}
		// render text - end


		// update screen - start
		texture.UpdateRGBA(nil, fb, 640)
		//renderer.SetDrawColor(0x40, 0x00, 0x40, 255)
		//renderer.Clear()
		renderer.Copy(texture2, nil, nil)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		// update screen - end

		// calculate speed
		frames++
		if sdl.GetTicks() > ticks_now {
			mult = sdl.GetTicks() - ticks_now
			ticks_now = sdl.GetTicks()
			stepCycles = p.CPU.AllCycles

			// cpu step ---------------------------------------------------------
			for {
				if (p.CPU.AllCycles - stepCycles) > CPU_STEP * uint64(mult) {
					break
				}
				_, stopped := p.CPU.Step()

				//if p.CPU.PC == 0x4c33 && p.CPU.RK == 0x38 {
				//	disasm=true
				//}
				//if p.CPU.PC == 0x4d93 && p.CPU.RK == 0x38 {
				//	disasm=false
				//}

				if disasm {
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.N, "n"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.V, "v"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.M, "m"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.X, "x"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.D, "d"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.I, "i"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.Z, "z"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.C, "c"))
					fmt.Fprintf(os.Stdout, " ")
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.B, "B"))
					fmt.Fprintf(os.Stdout, printCPUFlags(p.CPU.E, "E"))
					if p.CPU.M == 0 {
						fmt.Printf(" A  %04x (%7d) │",          p.CPU.RA, p.CPU.RA)
					} else {
						fmt.Printf(" A %02x %02x (%3d %3d) │", p.CPU.RAh, p.CPU.RAl, p.CPU.RAh, p.CPU.RAl)
					}
					fmt.Printf(" %4x ", p.CPU.RX)
					fmt.Printf("%s", p.CPU.DisassembleCurrentPC())
				}

				if stopped {
					running = false
					break
				}
			}
		}




		if (ticks_now - prev_ticks) >= 1000 {
			cyc, unit := showCPUSpeed(p.CPU.AllCycles - prevCycles)
			prevCycles = p.CPU.AllCycles
			fmt.Fprintf(os.Stdout, "frames: %4d ticks %d cpu cycles %10d speed %2d %s cpu.K:PC %02x:%04x\n", frames, (ticks_now - prev_ticks), p.CPU.AllCycles, cyc, unit, p.CPU.RK, p.CPU.PC)
			prev_ticks = ticks_now
			frames = 0
		}




		// keyboard ----------------------------------------------------------
		// https://github.com/veandco/go-sdl2-examples/blob/master/examples/keyboard-input/keyboard-input.go
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			/*
			case *sdl.TextInputEvent:
				fmt.Printf("TextInputEvent\n")
				for _, val := range t.Text {
					if val == 0 {
						break
					}
					p.GABE.InBuf.Enqueue(val)
				}
			*/

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
					case sdl.K_F9:
						if disasm {
							disasm = false 
						} else {
							disasm = true
						}
					default:
						gui.sendKey(t.Keysym.Scancode, t.State)
					}
				}

				if t.State == sdl.RELEASED {
					gui.sendKey(t.Keysym.Scancode, t.State)
				}


			}
		}


	}

	// return from FULLSCREEN
	if FULLSCREEN {
		window.SetDisplayMode(&current_mode)
	}

	memoryDump(p, 0xaf_0000)
	//renderer.Destroy()
	//window.Destroy()
	//sdl.Quit()
}
