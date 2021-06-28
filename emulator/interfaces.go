package emu

type Processor interface {
        Reset()
        Execute() uint32	// execute one or more steps and returns used cycles
        Step() uint32		// execute single step and returns used cycles
	GetCycles() uint64
	ResetCycles()
	TriggerIRQ()
	SetPC(uint32)

	Write_8(uint32, byte)
	Read_8(uint32) byte

	GetName() string	// get id as "cpu0" / "cpu1" of unit
}

type Bus interface {
	Write_8(uint32, uint8)
	Read_8(uint32) uint8
}

const (
        CPU_65c816 = 0
        CPU_68000  = 1
        CPU_68030  = 2
)

