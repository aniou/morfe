
package main

import (
	_ "encoding/binary"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	//"time"
	"github.com/aniou/go65c816/emulator/platform"
	_ "github.com/aniou/go65c816/lib/mylog"
)

// keyboard memory registers
const INT_MASK_REG1	= 0x00_014D
const INT_PENDING_REG1	= 0x00_0141

// some general consts
const CPU_CLOCK		= 14318000 // 14.381Mhz
const CURSOR_BLINK_RATE = 500

type GUI struct {
	p	   *platform.Platform
	fullscreen bool
}

type DEBUG struct {
	gui	bool
	cpu	bool
}
var debug = DEBUG{true, false}


// some support routines
// xxx - move it
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

/*
func loadFont(p *platform.Platform, fontset *[2048]uint32) {
	for i, v := range fontset {
		for j := 0; j < 8; j = j + 1 {
			v = v << 1
			if (v & 256) == 256 {
				//fmt.Printf("#")
				p.GPU.FONT[i*8+j] = 1
			} else {
				//fmt.Printf(" ")
				p.GPU.FONT[i*8+j] = 0
			}
		}
		//fmt.Printf("\n")
	}
}
*/

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

// xxx - window parametrization
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

/*
   a 'nicer' form, i.e. func (g *GUI) setFullscreen { ... } with 
   orig_mode as field of GUI struct leads to error:
   "panic: runtime error: cgo argument has Go pointer to Go pointer",
   during return to original mode,  so don't improve following in this way
*/

func setFullscreen(window *sdl.Window) sdl.DisplayMode {
	var wanted_mode = sdl.DisplayMode{sdl.PIXELFORMAT_ARGB8888, 640, 480, 60, nil}
	var result_mode sdl.DisplayMode
	display_index, _ := window.GetDisplayIndex()
	orig_mode, _ := sdl.GetCurrentDisplayMode(display_index)
	fmt.Printf("original mode width: %d\n", orig_mode.W)
	fmt.Printf("original mode heigt: %d\n", orig_mode.H)

	_, err := sdl.GetClosestDisplayMode(display_index, &wanted_mode, &result_mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get ClosestMode: %s\n", err)
		os.Exit(2)
	}
	fmt.Printf("wanted mode width: %d\n", result_mode.W)
	fmt.Printf("wanted mode heigt: %d\n", result_mode.H)
	window.SetDisplayMode(&result_mode)
	window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	return orig_mode
}


// -----------------------------------------------------------------------------
// MAIN HERE
// -----------------------------------------------------------------------------
func main() {
	var orig_mode	sdl.DisplayMode
	var event	sdl.Event
	var err		error
	var running	bool
	var disasm	bool
	var winWidth    int32 = 640
	var winHeight   int32 = 480
	var CPU_STEP	uint64 = 14318


	// platform init ---------------------------------------------------------------
	p := platform.New()
	gui := new(GUI)
	gui.fullscreen = false
	gui.p = p			// xxx - fix that mess
	p.InitGUI()


	// code load and PC set --------------------------------------------------------
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s filename.ini\n", os.Args[0])
	} else {
		gui.loadConfig(os.Args[1])
	}
 

	// some additional tweaks ------------------------------------------------------
	// XXX - move it somewhere
	p.CPU.Bus.EaWrite(0xAF_0005, 0x20) // border B 
	p.CPU.Bus.EaWrite(0xAF_0006, 0x00) // border G
	p.CPU.Bus.EaWrite(0xAF_0007, 0x20) // border R

	p.CPU.Bus.EaWrite(0xAF_0008, 0x20) // border X
	p.CPU.Bus.EaWrite(0xAF_0009, 0x20) // border Y

	p.CPU.Bus.EaWrite(0xAF_0010, 0x03) // VKY_TXT_CURSOR_CTRL_REG
	p.CPU.Bus.EaWrite(0xAF_0012, 0xB1) // VKY_TXT_CURSOR_CHAR_REG
	p.CPU.Bus.EaWrite(0xAF_0013, 0xED) // VKY_TXT_CURSOR_COLR_REG

	// act as gavin/gabe - copy "flash" area from 38:1000 to 00:1000 (0x200) bytes
	// jump tables
	for j := 0x1000; j < 0x1200; j = j + 1 {
		val := p.CPU.Bus.EaRead(uint32(0x38_0000 + j))
		p.CPU.Bus.EaWrite(uint32(j), val)

	}


	// graphics init ---------------------------------------------------------------
	// step 1: SDL
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Panic(sdl.GetError())
	}
	defer sdl.Quit()

	// step 2: Window
	var window *sdl.Window
	window, err = sdl.CreateWindow(
		"go65c816 / c256 emu",
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

	// step 4: textures 
	texture_txt := newTexture(renderer)
	texture_txt.SetBlendMode(sdl.BLENDMODE_BLEND)

	texture_bm0 := newTexture(renderer)
	texture_bm0.SetBlendMode(sdl.BLENDMODE_BLEND)

	texture_bm1 := newTexture(renderer)
	texture_bm1.SetBlendMode(sdl.BLENDMODE_BLEND)


	// -----------------------------------------------------------------------------
	sdl.SetHint("SDL_HINT_RENDER_BATCHING", "1")
	//sdl.StartTextInput()


	// variables for performance calculation ---------------------------------------
	var prev_ticks uint32 = sdl.GetTicks()    // FPS calculation
	var mult       uint32 = prev_ticks	  // CPU speed calculation
	var ticks_now, frames uint32              // CPU step and FPS calculation
	var stepCycles, prevCycles uint64 = 0, 0  // CPU speed calculation
	var prevCounter, maxCounter uint64 = 0, 0 // measuring speed of custom CPU counter
	var cursor_counter int32                  // how many ticks remains to blink cursor 
	var cyclesCheckpoint uint64 = 0		  // CPU cycles between checkpoints [nga]

	// current draw model ----------------------------------------------------------
	//
	// 1. fill by background color
	// 2. update texture from fb and apply with alpha
	//    - for bm0
	//    - then for bm1
	// 3. draw text 
	// 4. update texture with text and apply texture with alpha
	// 5. draw frames
	// 6. present

	// main loop -------------------------------------------------------------------
	running = true
	disasm  = false
	/* nga debug only - XXX - move to plugin
	var ipptr uint32 = 0
	var tval uint32 = 0
	var nval uint32 = 0
	var tos uint32 = 0
	var nos uint32 = 0
	*/
	for running {
		// step 1
		renderer.SetDrawColor(p.GPU.Background[0], p.GPU.Background[1], p.GPU.Background[2], 255)
		renderer.Clear()

		// step 2 - bm0 and bm1 are updated in vicky, when write is made
		if p.GPU.Master_L & 0x0C == 0x0C {					// todo?
			if p.GPU.BM0_visible {
				texture_bm0.UpdateRGBA(nil, p.GPU.BM0FB, 640)
				renderer.Copy(texture_bm0, nil, nil)
			}

			if p.GPU.BM1_visible  {
				texture_bm1.UpdateRGBA(nil, p.GPU.BM1FB, 640)
				renderer.Copy(texture_bm1, nil, nil)
			}
		}

		// step 3, 4
		if p.GPU.Master_L & 0x01 == 0x01 {					// todo ?
			p.GPU.RenderBitmapText()
			texture_txt.UpdateRGBA(nil, p.GPU.TFB, 640)
			renderer.Copy(texture_txt, nil, nil)
		}	

		// stea 5
		if p.GPU.Border_visible {
			renderer.SetDrawColor(p.GPU.Border_color_r, 
					      p.GPU.Border_color_g, 
					      p.GPU.Border_color_b, 
					      255)
			renderer.FillRects([]sdl.Rect{
				sdl.Rect{0, 0, 640, p.GPU.Border_y_size},
				sdl.Rect{0, 480-p.GPU.Border_y_size, 640, p.GPU.Border_y_size},
				sdl.Rect{0, p.GPU.Border_y_size,  p.GPU.Border_x_size, 480-p.GPU.Border_y_size},
				sdl.Rect{640-p.GPU.Border_x_size, p.GPU.Border_y_size, p.GPU.Border_x_size, 480-p.GPU.Border_y_size},
			})
		}

		// step 6
		renderer.Present()


		// cpu speed calculation --------------------------------------------
		frames++
		if sdl.GetTicks() > ticks_now {
			mult = sdl.GetTicks() - ticks_now
			ticks_now = sdl.GetTicks()
			stepCycles = p.CPU.AllCycles

			// cursor calculation - flip every CURSOR_BLINK_RATE ticks ----------
			cursor_counter = cursor_counter - int32(mult)
			if cursor_counter <= 0 {
				cursor_counter = CURSOR_BLINK_RATE
				p.GPU.Cursor_visible = ! p.GPU.Cursor_visible
			}

			// cpu step ---------------------------------------------------------
			cpu_loop:
			for {
				if (p.CPU.AllCycles - stepCycles) > CPU_STEP * uint64(mult) {
					break
				}
				_, stopped := p.CPU.Step()
				
				// debugging interface, created around WDM opcode
				if debug.cpu && stopped {
					switch p.CPU.WDM {
					case 0x0b:		// count cycles
						fmt.Printf("%%checkpoint: %d cycles from previous\n", 
									p.CPU.AllCycles - cyclesCheckpoint)
						cyclesCheckpoint = p.CPU.AllCycles
					case 0x10:		// enable disasm
						if disasm == false {
							disasm = true
							fmt.Printf("%s", p.CPU.DisassemblePreviousPC())
						}
					case 0x11:		// disable disasm
						if disasm {
							disasm = false
							fmt.Printf("... disassembling stopped\n")
						}
					case 0x20:		// stop emulator
						running = false
						fmt.Printf("...emulator stopped by WDM #$20 at cpu.K:PC %02x:%04x\n", p.CPU.PRK, p.CPU.PPC)
						break cpu_loop
					default:
						fmt.Printf("%WARN: unknown WDM opcode %d\n", p.CPU.WDM)
					}
					p.CPU.WDM = 0
					stopped = false
				}
				
				if disasm {
					// XXX - move it do subroutine
					fmt.Printf(printCPUFlags(p.CPU.N, "n"))
					fmt.Printf(printCPUFlags(p.CPU.V, "v"))
					fmt.Printf(printCPUFlags(p.CPU.M, "m"))
					fmt.Printf(printCPUFlags(p.CPU.X, "x"))
					fmt.Printf(printCPUFlags(p.CPU.D, "d"))
					fmt.Printf(printCPUFlags(p.CPU.I, "i"))
					fmt.Printf(printCPUFlags(p.CPU.Z, "z"))
					fmt.Printf(printCPUFlags(p.CPU.C, "c"))
					fmt.Printf(" ")
					fmt.Printf(printCPUFlags(p.CPU.B, "B"))
					fmt.Printf(printCPUFlags(p.CPU.E, "E"))
					fmt.Printf(" DBR %02x", p.CPU.RDBR)
					if p.CPU.M == 0 {
						fmt.Printf("│A  %04x (%7d)",          p.CPU.RA, p.CPU.RA)
					} else {
						fmt.Printf("│A %02x %02x (%3d %3d)", p.CPU.RAh, p.CPU.RAl, p.CPU.RAh, p.CPU.RAl)
					}
					if p.CPU.X == 0 {
						fmt.Printf("│X  %04x (%7d)",          p.CPU.RX, p.CPU.RX)
						fmt.Printf("│Y  %04x (%7d)│",         p.CPU.RY, p.CPU.RY)
					} else {
						fmt.Printf("│X    %02x (    %3d)",  p.CPU.RXl, p.CPU.RXl)
						fmt.Printf("│Y    %02x (    %3d)│", p.CPU.RYl, p.CPU.RYl)
					}

					/*
					// it is going a little too far, but I want to know... (nga debug)
					tos  = (uint32(p.CPU.Bus.EaRead(0xd7)) << 8) + uint32(p.CPU.Bus.EaRead(0xd6)) + 0x10004
					nos  = tos-4
					tval = (uint32(p.CPU.Bus.EaRead(tos+3)) << 24) +
					       (uint32(p.CPU.Bus.EaRead(tos+2)) << 16) +
					       (uint32(p.CPU.Bus.EaRead(tos+1)) << 8)  +
					        uint32(p.CPU.Bus.EaRead(tos))
					nval = (uint32(p.CPU.Bus.EaRead(nos+3)) << 24) +
					       (uint32(p.CPU.Bus.EaRead(nos+2)) << 16) +
					       (uint32(p.CPU.Bus.EaRead(nos+1)) << 8)  +
					        uint32(p.CPU.Bus.EaRead(nos))
					ipptr = (uint32(p.CPU.Bus.EaRead(0xd5)) << 24) +
					        (uint32(p.CPU.Bus.EaRead(0xd4)) << 16) +
					        (uint32(p.CPU.Bus.EaRead(0xd3)) << 8)  +
					        uint32(p.CPU.Bus.EaRead(0xd2))
					fmt.Printf("IP %02x%02x│", p.CPU.Bus.EaRead(0xd1), p.CPU.Bus.EaRead(0xd0))
					fmt.Printf("SP %02x%02x│", p.CPU.Bus.EaRead(0xd7), p.CPU.Bus.EaRead(0xd6))
					fmt.Printf("RP %02x%02x│", p.CPU.Bus.EaRead(0xd9), p.CPU.Bus.EaRead(0xd8))
					fmt.Printf("NOS %08x│", nval)
					fmt.Printf("TOS %08x│", tval)
					fmt.Printf("PTR %08x│", ipptr)
					*/
					fmt.Printf("%s", p.CPU.DisassembleCurrentPC())
					break
				}
			}
		}

		// performance info --------------------------------------------------
		if (ticks_now - prev_ticks) >= 1000 {
			cyc, unit  := showCPUSpeed(p.CPU.AllCycles - prevCycles)
			prevCycles  = p.CPU.AllCycles

			deltaCounter := p.CPU.Counter - prevCounter
			if deltaCounter > maxCounter {
				maxCounter = deltaCounter
			}
			prevCounter = p.CPU.Counter
			if ! disasm {
				fmt.Fprintf(os.Stdout, "frames: %4d ticks %d cpu cycles %10d speed %2d %s cpu counter %10d max %10d cpu.K:PC %02x:%04x\n", 
				                        frames, (ticks_now - prev_ticks), 
							p.CPU.AllCycles, cyc, unit, 
							deltaCounter, maxCounter,
							p.CPU.RK, p.CPU.PC)
			}
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
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)

				if t.State == sdl.PRESSED {
					if t.Repeat > 0 {
						continue
					}
					switch t.Keysym.Sym {
					case sdl.K_F12:
						running = false
					case sdl.K_F11:
						if gui.fullscreen {
							gui.fullscreen = false
							window.SetDisplayMode(&orig_mode)
							window.SetFullscreen(0)
						} else {
							gui.fullscreen = true
							orig_mode = setFullscreen(window)
						}
					case sdl.K_F10:
						//loadFont(p, &font_c256_8x8)
					case sdl.K_F9:
						if disasm {
							disasm = false 
							fmt.Printf("... disasm stopped\n")
						} else {
							disasm = true
						}
					default:
						gui.sendKey(t.Keysym.Scancode, t.State)
					}
				}

				if t.State == sdl.RELEASED {
					switch t.Keysym.Sym {
					case sdl.K_F12,
					     sdl.K_F11,
					     sdl.K_F10,
					     sdl.K_F9:
					default:
						gui.sendKey(t.Keysym.Scancode, t.State)
					}
				}

			} // SDL event switch/case
		} // SDL event loop
	} // main loop

	// return from FULLSCREEN
	if gui.fullscreen {
		window.SetDisplayMode(&orig_mode)
		window.SetFullscreen(0)
	}

	//memoryDump(p, 0x0)

}
