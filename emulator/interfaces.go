package emu

type Processor interface {
        Reset()
        Step() uint32
	GetCycles() uint64
	ResetCycles()
	TriggerIRQ()
	SetPC(uint32)

        // at leas two attributes should be available
        // Cycles
        // Enabled
        // Type
}

type Bus interface {
	EaWrite(uint32, uint8)
	EaRead(uint32) uint8
}

const (
        CPU_65c816 = 0
        CPU_68000  = 1
        CPU_68030  = 2
)

