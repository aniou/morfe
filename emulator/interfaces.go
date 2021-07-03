package emu

type Processor interface {
        Reset()
        Execute() uint32	           // execute one or more steps and returns used cycles
        Step() uint32		           // execute single step and returns used cycles
	GetRegisters() map[string]uint32   // returns all registers of CPU
	GetType() uint			   // returns CPU id (when many types are available)
	Dissasm() string
	GetCycles() uint32		   // number of cycles used by last step
	GetAllCycles() uint64	           // cumulative number of cycles used
	StatusString() string		   // string that represents status flags
	ResetCycles()
	TriggerIRQ()
	SetRegister(string, uint32) error  // set selected register
	SetPC(uint32)			   // redundant to SetRegister but convinient

	Write_8(uint32, byte)		   // write byte to   cpu memory
	Read_8(uint32) byte                // read  byte from cpu memory

	GetName() string	           // get id as "cpu0" / "cpu1" of unit
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

