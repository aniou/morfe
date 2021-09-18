
package main

import (
        _ "encoding/binary"
        "fmt"
        "github.com/veandco/go-sdl2/sdl"
        "log"
        "os"
	"runtime"
	_ "runtime/pprof"
        _ "time"
        // "github.com/aniou/morfe/emulator/cpu"
        "github.com/aniou/morfe/emulator"

        "github.com/aniou/morfe/emulator/platform"
        "github.com/aniou/morfe/lib/mylog"
)

// some general consts
const CPU_CLOCK         = 14318000 // 14.381Mhz (not used)
const CURSOR_BLINK_RATE = 500      // in ms (milliseconds)

type GUI struct {
        p          *platform.Platform
        fullscreen bool
}

type DEBUG struct {
        gui     bool
        cpu     bool
}
var debug = DEBUG{true, false}

var p = platform.New()          // must be global now

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

func memoryDump(cpu emu.Processor, address uint32) {
        var x uint16
        var a,b uint32
	var data []byte = make([]byte, 16)
	start := address & 0xFFFF_FFF0

        for a = 0; a<0x100; a=a+16 {
		for b = 0; b<16; b=b+1 {
			data[b] = cpu.Read_8(start + a + b)
		}
                bank := byte((start+a) >> 16)
                addr := uint16(start+a)
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
                        fmt.Printf("                   │                       │")
                }
        }
        fmt.Printf("\n")
}

func waitForEnter() {
        fmt.Println("\nPress the Enter Key")
        fmt.Scanln() // wait for Enter Key
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
        var orig_mode   sdl.DisplayMode
        var event       sdl.Event
        var err         error
        var running     bool
        var disasm      bool		// indicator for debug/disasm mode
        var live_disasm bool		// indicator for debug/disasm mode
        var winWidth    int32 = 640
        var winHeight   int32 = 480
        var CPU0_STEP   uint64 = 14318 // 14.318 MHz in milliseconds, apply for 65c816
        var CPU1_STEP   uint64 = 20000 // I'm able to achieve 25Mhz too
	var ch		chan string
	var msg	        string
	var pcfg	*platform.Config
	var gpu		*emu.GPU_common

	// first things first
        if len(os.Args) < 2 {
                log.Fatalf("Usage: %s filename.ini\n", os.Args[0])
        } 

	runtime.LockOSThread()

	/*
        f, err := os.Create("go65c816.profile")
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
	*/

        // platform init ---------------------------------------------------------------
        //p := platform.New()           // must be global now - but it is still true?
	gui := new(GUI)
        gui.fullscreen = false
        gui.p = p                       // xxx - fix that mess

	if pcfg, err = p.LoadPlatformConfig(os.Args[1]); err != nil {
		log.Fatalf("%s", err)
	}

	switch pcfg.Mode {
	case "fmx-like":
		p.SetFMX()
	case "frankenmode":
		p.SetFranken()
	case "genx-like":
		p.SetGenX()
	default:
		log.Fatalf("unknown mode %s", pcfg.Mode)
	}
	gpu = p.GPU.GetCommon()

	// kernel and others files loading also here
        p.LoadCpuConfig(os.Args[1])

	// platform-specific init function 
	p.Init()

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
        //renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
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
        var prev_ticks uint32 = sdl.GetTicks()      // FPS calculation (1/1000 of second)
        var ms_elapsed uint64 = uint64(prev_ticks)  // how many ms elapsed from last check?
        var ticks_now, frames uint32                // CPU step and FPS calculation
        var prevCycles0, prevCycles1 uint64 = 0, 0  // CPU speed calculation
        var cursor_counter int32                    // how many ticks remains to blink cursor 

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
        live_disasm  = false

	desired_cycles0 := uint64(CPU0_STEP)
	desired_cycles1 := uint64(CPU1_STEP)

        for running {
                // step 1
                renderer.SetDrawColor(gpu.Background[0], gpu.Background[1], gpu.Background[2], 255)
                renderer.Clear()

                // step 2 - bm0 and bm1 are updated in vicky, when write is made
                if gpu.Master_L & 0x0C == 0x0C {                                      // todo?
                        if gpu.BM0_visible {
                                texture_bm0.UpdateRGBA(nil, gpu.BM0FB, 640)
                                renderer.Copy(texture_bm0, nil, nil)
                        }

                        if gpu.BM1_visible  {
                                texture_bm1.UpdateRGBA(nil, gpu.BM1FB, 640)
                                renderer.Copy(texture_bm1, nil, nil)
                        }
                }

                // step 3, 4
                if gpu.Master_L & 0x01 == 0x01 {                                      // todo ?
                        p.GPU.RenderBitmapText()
                        texture_txt.UpdateRGBA(nil, gpu.TFB, 640)
                        renderer.Copy(texture_txt, nil, nil)
                }       

                // stea 5
                if gpu.Border_visible {
                        renderer.SetDrawColor(gpu.Border_color_r, 
                                              gpu.Border_color_g, 
                                              gpu.Border_color_b, 
                                              255)
                        renderer.FillRects([]sdl.Rect{
                                sdl.Rect{0, 0, 640, gpu.Border_y_size},
                                sdl.Rect{0, 480-gpu.Border_y_size, 640, gpu.Border_y_size},
                                sdl.Rect{0, gpu.Border_y_size,  gpu.Border_x_size, 480-gpu.Border_y_size},
                                sdl.Rect{640-gpu.Border_x_size, gpu.Border_y_size, gpu.Border_x_size, 480-gpu.Border_y_size},
                        })
                }

                // step 6
                renderer.Present()


                // cpu running routines ---------------------------------------------
		if ! disasm {
			frames++
			if sdl.GetTicks() > ticks_now {
				ms_elapsed = uint64(sdl.GetTicks() - ticks_now)
				ticks_now = sdl.GetTicks()

				// cursor calculation - flip every CURSOR_BLINK_RATE ticks ----------
				cursor_counter = cursor_counter - int32(ms_elapsed)
				if cursor_counter <= 0 {
					cursor_counter = CURSOR_BLINK_RATE
					gpu.Cursor_visible = ! gpu.Cursor_visible
				}

				// WARNING - it has tendency to going in tight loop if
				//           system is too slow to do desired number of
				//           cycles per ms when *ms is used

				if p.CPU0.IsEnabled() {
					for p.CPU0.GetAllCycles() < desired_cycles0 {
						if live_disasm {
							fmt.Printf("%s", p.CPU0.DisassembleCurrentPC())	// XXX - change it for CURRENT CPU
						}

						p.CPU0.Execute()
					}
					desired_cycles0 = desired_cycles0 + CPU0_STEP*ms_elapsed
				}

				if p.CPU1.IsEnabled() {
					for p.CPU1.GetAllCycles() < desired_cycles1 {
						p.CPU1.Execute()
					}
					desired_cycles1 = desired_cycles1 + CPU1_STEP*ms_elapsed
				}
			}

			// performance info --------------------------------------------------
			if (ticks_now - prev_ticks) >= 1000 {   // once per second
				if ! disasm {  // TODO - redundant, but remove after shaping a main routing
					spd0, unit0  := showCPUSpeed(p.CPU0.GetAllCycles() - prevCycles0)
					spd1, unit1  := showCPUSpeed(p.CPU1.GetAllCycles() - prevCycles1)
					prevCycles0  = p.CPU0.GetAllCycles()
					prevCycles1  = p.CPU1.GetAllCycles()
					fmt.Fprintf(os.Stdout, 
						    "frames: %4d ticks %d cpu0 cycles %10d speed %2d %s cpu1 cycles %10d speed %d %s\n", 
							    frames, (ticks_now - prev_ticks), 
							    p.CPU0.GetAllCycles(), spd0, unit0, 
							    p.CPU1.GetAllCycles(), spd1, unit1)
				}
				prev_ticks = ticks_now
				frames = 0
			}
			
		} else {
			msg = <-ch
			switch msg {
			case "step":
				p.CPU1.Step()
				ch<-"done"
			case "disable_cpu0":
				p.CPU0.Enable(false)
			case "disable_cpu1":
				p.CPU1.Enable(false)
			case "enable_cpu0":
				p.CPU0.Enable(true)
			case "enable_cpu1":
				p.CPU1.Enable(true)
			case "run":
				ticks_now = sdl.GetTicks()
				disasm = false
				close(ch)
				mylog.Logger.LogOutput = os.Stdout
			default:
				log.Panicln("channel from tui: received unknown message %s", msg)
			}
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
                                //      t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)

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
						if ! live_disasm {
							live_disasm = true
						} else {
							live_disasm = false
						}
                                        case sdl.K_F9:
                                                if ! disasm {
                                                        disasm = true
							ch = make(chan string, 1)
							p.CPU1.ResetCycles()
							go func() {
								mainTUI(ch, p.CPU1)		// XXX - parametrize that!
							}()
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

        memoryDump(p.CPU0, 0x0c)
        //memoryDump(p.CPU0, 0xAF_8000)
        //memoryDump(p.CPU0, 0xAF_1f40)

}
