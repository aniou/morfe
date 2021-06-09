package cpu68xxx

// #cgo CFLAGS: -I../../../Musashi 
// #cgo LDFLAGS: ../../../Musashi/m68kcpu.o  ../../../Musashi/m68kops.o  ../../../Musashi/softfloat/softfloat.o
// #include <m68k.h>
import "C"

import (
        _ "encoding/binary"
        _ "fmt"
)

type mem_read   func(uint32) byte
type mem_write  func(uint32, byte)

var EaRead  	mem_read
var EaWrite	mem_write

type CPU struct {
	Speed   uint32		// in milliseconds
	Enabled bool
	Type	byte

	AllCycles uint64	// cumulative number of cycles of CPU instance
}

// XXX - change EaRead/Write to BigEndian ones!

//export m68k_read_memory_8
func m68k_read_memory_8(addr C.uint) C.uint {
        //fmt.Printf("m68k read8  %8x", addr)

        a   := uint32(addr)
        val := EaRead(a)

        //fmt.Printf(" val  %8x %d\n", val, val)
        return C.uint(val)
}

//export m68k_read_memory_16
func m68k_read_memory_16(addr C.uint) C.uint {
        //fmt.Printf("m68k read16  %8x", addr)

        a   := uint32(addr)
        val := ( uint32(EaRead(a))   << 8 ) |
                 uint32(EaRead(a+1))

        //fmt.Printf(" val  %8x %d\n", val, val)
        return C.uint(val)
}

//export m68k_read_memory_32
func m68k_read_memory_32(addr C.uint) C.uint {
        //fmt.Printf("m68k read32  %8x", addr)

        a   := uint32(addr)
        val := ( uint32(EaRead(a))   <<  24 ) |
               ( uint32(EaRead(a+1)) <<  16 ) |
               ( uint32(EaRead(a+2)) <<   8 ) |
                 uint32(EaRead(a+3))

        //fmt.Printf(" val  %8x %d\n", val, val)
        return C.uint(val)
}

//export m68k_write_memory_8
func m68k_write_memory_8(addr, val C.uint) {
        //fmt.Printf("m68k write8  %8x val  %8x %d\n", addr, val, val)

        a   := uint32(addr)
        EaWrite(a, byte(val))
        return
}

//export m68k_write_memory_16
func m68k_write_memory_16(addr, val C.uint) {
        //fmt.Printf("m68k write16 %8x val  %8x %d\n", addr, val, val)

        a   := uint32(addr)
        EaWrite(a,   byte((val >> 8) & 0xff))
        EaWrite(a+1, byte( val       & 0xff))
        return
}

//export m68k_write_memory_32
func m68k_write_memory_32(addr, val C.uint) {
        //fmt.Printf("m68k write32 %8x val  %8x %d\n", addr, val, val)

        a   := uint32(addr)
        EaWrite(a,   byte((val >> 24) & 0xff))
        EaWrite(a+1, byte((val >> 16) & 0xff))
        EaWrite(a+2, byte((val >>  8) & 0xff))
        EaWrite(a+3, byte( val        & 0xff))
        return
}



func New(speed uint32, r mem_read, w mem_write) *CPU {
	cpu := CPU{Speed: speed}
        EaRead  = r
        EaWrite = w
	C.m68k_init();
        C.m68k_set_cpu_type(C.M68K_CPU_TYPE_68EC030)
	return &cpu
}

// to fulfill interface, that doesn't allow direct acces to fields
func (cpu *CPU) GetCycles() uint64 {
        return cpu.AllCycles
}

// to fulfill interface, that doesn't allow direct acces to fields
func (cpu *CPU) ResetCycles() {
        cpu.AllCycles=0
}

func (c *CPU) Reset() {
	// just for test
	m68k_write_memory_32(0,           0x10_0000)    // stack
        m68k_write_memory_32(4,           0x20_0000)    // instruction pointer
        m68k_write_memory_16(0x20_0000,      0x7041)    // moveq  #41, D0
        m68k_write_memory_16(0x20_0002,      0x13C0)    // move.b D0, $AFA000
        m68k_write_memory_32(0x20_0004, 0x00AF_A000)    // ...
        m68k_write_memory_32(0x20_0008, 0x60F6_4E71)    // bra to 20_0000

	// normal
	C.m68k_pulse_reset()
	return
}

// there is a small problem here - go65c816 uses
// number of instruction here and returns spent
// cycles - Musashi gets amount of cycles and
// returns number of spent cycles, so there are
// difference with calculations

func (c *CPU) Step() uint32 {
	cycles := C.m68k_execute(1000)		// dummy value
	c.AllCycles=c.AllCycles+uint64(cycles)
	return uint32(cycles)
}

func (c *CPU) TriggerIRQ() {
	return
}

func (c *CPU) SetPC(addr uint32) {
	return
}
