package cpu68xxx

// #cgo CFLAGS: -I../../../Musashi  -I../../lib/musashi-c-wrapper
// #cgo LDFLAGS: ../../lib/musashi-c-wrapper/shim.o ../../../Musashi/m68kcpu.o ../../../Musashi/m68kdasm.o ../../../Musashi/m68kops.o  ../../../Musashi/softfloat/softfloat.o
// #include "shim.h"
// #include "m68k.h"
import "C"

import (
        _ "encoding/binary"
        "fmt"
        "github.com/aniou/go65c816/emulator"
)


var bus emu.Bus
var (
        regMapping = map[string]C.m68k_register_t{
                "D0"   : C.M68K_REG_D0,            /* Data registers */
                "D1"   : C.M68K_REG_D1,
                "D2"   : C.M68K_REG_D2,
                "D3"   : C.M68K_REG_D3,
                "D4"   : C.M68K_REG_D4,
                "D5"   : C.M68K_REG_D5,
                "D6"   : C.M68K_REG_D6,
                "D7"   : C.M68K_REG_D7,
                "A0"   : C.M68K_REG_A0,            /* Address registers */
                "A1"   : C.M68K_REG_A1,
                "A2"   : C.M68K_REG_A2,
                "A3"   : C.M68K_REG_A3,
                "A4"   : C.M68K_REG_A4,
                "A5"   : C.M68K_REG_A5,
                "A6"   : C.M68K_REG_A6,
                "A7"   : C.M68K_REG_A7,
                "PC"   : C.M68K_REG_PC,            /* Program Counter */
                "SR"   : C.M68K_REG_SR,            /* Status Register */
                "SP"   : C.M68K_REG_SP,            /* The current Stack Pointer (located in A7) */
                "USP"  : C.M68K_REG_USP,           /* User Stack Pointer */
                "ISP"  : C.M68K_REG_ISP,           /* Interrupt Stack Pointer */
                "MSP"  : C.M68K_REG_MSP,           /* Master Stack Pointer */
                "SFC"  : C.M68K_REG_SFC,           /* Source Function Code */
                "DFC"  : C.M68K_REG_DFC,           /* Destination Function Code */
                "VBR"  : C.M68K_REG_VBR,           /* Vector Base Register */
                "CACR" : C.M68K_REG_CACR,          /* Cache Control Register */
                "CAAR" : C.M68K_REG_CAAR,          /* Cache Address Register */
                "PPC"  : C.M68K_REG_PPC,           /* Previous value in the program counter */
        }

)


type CPU struct {
        Speed       uint32      // in milliseconds
        Enabled     bool
        Type        uint
        name        string      // cpu0, cpu1, etc.

        Cycles      uint32      // number of cycles used by last step
        AllCycles   uint64      // cumulative number of cycles of CPU instance
}

//export go_m68k_read_memory_8
func go_m68k_read_memory_8(addr C.uint) C.uint {
        //fmt.Printf("m68k read8  %8x", addr)

        a   := uint32(addr)
        val := bus.Read_8(a)

        //fmt.Printf(" val  %8x %d\n", val, val)
        return C.uint(val)
}

//export go_m68k_read_memory_16
func go_m68k_read_memory_16(addr C.uint) C.uint {
        //fmt.Printf("m68k read16  %8x", addr)

        a   := uint32(addr)
        val := ( uint32(bus.Read_8(a))   << 8 ) |
                 uint32(bus.Read_8(a+1))

        //fmt.Printf(" val  %8x %d\n", val, val)
        return C.uint(val)
}

//export go_m68k_read_memory_32
func go_m68k_read_memory_32(addr C.uint) C.uint {
        //fmt.Printf("m68k read32  %8x", addr)

        a   := uint32(addr)
        val := ( uint32(bus.Read_8(a))   <<  24 ) |
               ( uint32(bus.Read_8(a+1)) <<  16 ) |
               ( uint32(bus.Read_8(a+2)) <<   8 ) |
                 uint32(bus.Read_8(a+3))

        //fmt.Printf(" val  %8x %d\n", val, val)
        return C.uint(val)
}

//export go_m68k_write_memory_8
func go_m68k_write_memory_8(addr, val C.uint) {
        //fmt.Printf("m68k write8  %8x val  %8x %d\n", addr, val, val)

        a   := uint32(addr)
        bus.Write_8(a, byte(val))
        return
}

//export go_m68k_write_memory_16
func go_m68k_write_memory_16(addr, val C.uint) {
        //fmt.Printf("m68k write16 %8x val  %8x %d\n", addr, val, val)

        a   := uint32(addr)
        bus.Write_8(a,   byte((val >> 8) & 0xff))
        bus.Write_8(a+1, byte( val       & 0xff))
        return
}

//export go_m68k_write_memory_32
func go_m68k_write_memory_32(addr, val C.uint) {
        //fmt.Printf("m68k write32 %8x val  %8x %d\n", addr, val, val)

        a   := uint32(addr)
        bus.Write_8(a,   byte((val >> 24) & 0xff))
        bus.Write_8(a+1, byte((val >> 16) & 0xff))
        bus.Write_8(a+2, byte((val >>  8) & 0xff))
        bus.Write_8(a+3, byte( val        & 0xff))
        return
}


func New(b emu.Bus, name string) *CPU {
        cpu := CPU{name: name}
        bus     = b
        C.m68k_init_ram();
        C.m68k_init();
        C.m68k_set_cpu_type(C.M68K_CPU_TYPE_68EC030)
        cpu.Type = uint(C.M68K_CPU_TYPE_68EC030)                // XXX - parametrize it!
        return &cpu
}

func (cpu *CPU) GetType() uint {
        return cpu.Type
}

func (cpu *CPU) GetName() string {
        return cpu.name
}

func (cpu *CPU) GetCycles() uint32 {
        return cpu.Cycles
}

// to fulfill interface, that doesn't allow direct acces to fields
func (cpu *CPU) GetAllCycles() uint64 {
        return cpu.AllCycles
}

// to fulfill interface, that doesn't allow direct acces to fields
func (cpu *CPU) ResetCycles() {
        cpu.AllCycles=0
        cpu.Cycles=0
}

func (c *CPU) Write_8(addr uint32, val byte) {
        C.m68k_write_memory_8(C.uint(addr), C.uint(val))
}

func (c *CPU) Read_8(addr uint32) byte {
        return byte(C.m68k_read_memory_8(C.uint(addr)))
}
func (c *CPU) Reset() {

        // just for test
        C.m68k_write_memory_32(0,           0x10_0000)    // stack
        C.m68k_write_memory_32(4,           0x20_0000)    // instruction pointer
        C.m68k_write_memory_16(0x20_0000,      0x7042)    // moveq  #41, D0
        C.m68k_write_memory_16(0x20_0002,      0x13C0)    // move.b D0, $AFA000
        C.m68k_write_memory_32(0x20_0004, 0x00AF_A000)    // ...
        //C.m68k_write_memory_32(0x20_0004, 0x00A0_A000)    // ...
        C.m68k_write_memory_32(0x20_0008, 0x60F6_4E71)    // bra to 20_0000
        
        // normal
        C.m68k_pulse_reset()
        return
}

// there is a small problem here - go65c816 uses
// number of instruction here and returns spent
// cycles - Musashi gets amount of cycles and
// returns number of spent cycles, so there are
// difference with calculations

func (c *CPU) Execute() uint32 {
        cycles      := C.m68k_execute(1000)          // dummy value
        c.AllCycles  = c.AllCycles+uint64(cycles)
        c.Cycles     = uint32(cycles)
        return uint32(cycles)
}

func (c *CPU) Step() uint32 {
        cycles      := C.m68k_execute_step()         // provided by wrapper
        c.AllCycles  = c.AllCycles+uint64(cycles)
        c.Cycles     = uint32(cycles)
        return uint32(cycles)
}


func (c *CPU) TriggerIRQ() {
        return
}

func (c *CPU) SetPC(val uint32) {
        C.m68k_set_reg(C.M68K_REG_PC, C.uint(val))
}

func (c *CPU) SetRegister(reg string, val uint32) error {
	if c_reg_id, exists := regMapping[reg]; exists {
		C.m68k_set_reg(c_reg_id, C.uint(val))
                return nil
        } else {
		return fmt.Errorf("m68k: unknown register %v", reg)
	}
	
}

func (c *CPU) GetRegisters() map[string]uint32 {
        var register = map[string]uint32{}

        for name, id := range regMapping {
                register[name] = uint32(C.m68k_get_reg(nil, id))
        }

        return register
}

func (c *CPU) Dissasm() string {
        dpc := C.m68k_get_reg(nil, C.M68K_REG_PC)
        b := make([]C.char, 512)

        _ = C.m68k_disassemble_program(&b[0], dpc, C.uint(c.Type))
        return C.GoString(&b[0])
}

func cpuFlag(val uint16, flag string) string {
        if val > 0 {
                return flag
        } else {
                return "-"
        }
}

func (c *CPU) StatusString() string {
        // 0 MS 7 XNZVC
        sr     := uint16(C.m68k_get_reg(nil, C.M68K_REG_SR))

        status := fmt.Sprintf("%d ", (sr >> 14))
        status += cpuFlag(sr & 0b0010_0000_0000_0000, "S")
        status += cpuFlag(sr & 0b0001_0000_0000_0000, "M")
        status += fmt.Sprintf(" %d ", (sr & 0b0000_0111_0000_0000) >> 8)
        status += cpuFlag(sr & 0b0000_0000_0001_0000, "X")
        status += cpuFlag(sr & 0b0000_0000_0000_1000, "N")
        status += cpuFlag(sr & 0b0000_0000_0000_0100, "Z")
        status += cpuFlag(sr & 0b0000_0000_0000_0010, "V")
        status += cpuFlag(sr & 0b0000_0000_0000_0001, "C")
        return status
}
