
package cpu_dummy

import (
	"fmt"
	"log"

        "github.com/aniou/morfe/emulator"
)

type CPU struct {
	Bus	emu.Bus
	name	string
	enabled bool
}

func New(bus emu.Bus, name string) *CPU {
        cpu    := CPU{Bus: bus, name: name, enabled: true}
        return &cpu
}

func (cpu *CPU) Write_8(addr uint32, val byte) {
	return
}

func (cpu *CPU) Read_8(addr uint32) byte {
        return 0x00
}

func (cpu *CPU) SetPC(addr uint32) {
        return
}

func (cpu *CPU) Step() (uint32) {
        return 0x00
}

func (cpu *CPU) TriggerIRQ() {
	return
}

func (cpu *CPU) Reset() {
	return
}

func (cpu *CPU) StatusString() string {
        return "????????"
}

func (cpu *CPU) ResetCycles() {
	return
}

func (cpu *CPU) Execute() (uint32) {
	return 0x00
}

func (cpu *CPU) GetAllCycles() uint64 {
        return 0x00
}

func (cpu *CPU) GetCycles() uint32 {
        return 0x00
}

func (cpu *CPU) IsEnabled() bool {
        return cpu.enabled
}

func (cpu *CPU) GetName() string {
        return cpu.name
}

func (cpu *CPU) Enable(state bool) {
        cpu.enabled = state
}

func (cpu *CPU) Dissasm() string {
        log.Panic("GetType in dummy is not implemented yet!")
        return ""
}

func (cpu *CPU) GetType() uint {
        log.Panic("GetType in dummy is not implemented yet!")
        return 0
}

func (c *CPU) SetRegister(reg string, val uint32) error {
        return fmt.Errorf("SetRegister in dummy is not implemented yet")
}

func (c *CPU) GetRegisters() map[string]uint32 {
        var register = map[string]uint32{}
        log.Panic("GetRegisters in dummy is not implemented yet!")
        return register
}

func (c *CPU) DisassembleCurrentPC() string {
        return "m68k: DisassembleCurrentPC() is not implemented yet\n"
}


