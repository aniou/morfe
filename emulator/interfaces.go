package emu

import (
	"github.com/aniou/morfe/emulator/vram"
)

type Processor interface {
        Reset()
        Execute() uint32	           // execute one or more steps and returns used cycles
        Step() uint32		           // execute single step and returns used cycles
	GetRegisters() map[string]uint32   // returns all registers of CPU
	GetType() uint			   // returns CPU id (when many types are available)
	IsEnabled() bool                   // as name suggests
	Enable(bool)			   // enables/disables CPU
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
	Write_8(byte, uint32, byte)
	Read_8 (byte, uint32) byte
}

type Memory interface {
        Write(addr uint32, value byte)  error
        Read (addr uint32)             (byte, error)
	Name() string
        Size() (uint32, uint32)

        //Shutdown()
        //Clear()
        //Dump(address uint32) []byte
}

type GPU interface {
        Write(addr uint32, value byte)  error
        Read (addr uint32)             (byte, error)
	Name() string
        Size() (uint32, uint32)

	GetCommon() *GPU_common
	RenderBitmapText()
}

const (
        CPU_65c816 = 0
        CPU_68000  = 1
        CPU_68030  = 2
)

const (
	M_USER  = 0
	M_SV    = 1
)

// a 'common' set of Vicky's data
type GPU_common struct {
	Text    *vram.Vram	// text memory attached at platform level

        TFB     []uint32       // text   framebuffer
        BM0FB   []uint32       // bitmap0 framebuffer
        BM1FB   []uint32       // bitmap1 framebuffer

        // some convinient registers that should be converted
        // into some kind of memory indexes...
        Master_L        byte    // MASTER_CTRL_REG_L
        Master_H        byte    // MASTER_CTRL_REG_H
        Cursor_visible  bool
        Border_visible  bool
        BM0_visible     bool
        BM1_visible     bool

        Border_color_b  byte
        Border_color_g  byte
        Border_color_r  byte
        Border_x_size   int32
        Border_y_size   int32
        Background      [3]byte         // r, g, b
}

