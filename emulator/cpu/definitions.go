package cpu

type Processor interface {
        Init()
        Reset()
        Step(uint32)

        // at leas two attributes should be available
        // Cycles
        // Enabled
        // Type
}

const (
        CPU_65c816 = 0
        CPU_68000  = 1
        CPU_68030  = 2
)
