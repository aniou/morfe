package emu

type Processor interface {
        Reset()
        Step() uint32
	GetCycles() uint64
	ResetCycles()
	TriggerIRQ()
	SetPC(uint32)

	Write_8(uint32, byte)
	Read_8(uint32) byte

	GetName() string

        // at leas two attributes should be available
        // Cycles
        // Enabled
        // Type
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

