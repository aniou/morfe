
package emu

// copy-pasted Devices/Interrupts.cs

const (
	R0_FNX0_INT00_SOF   = 1    // Start of Frame @ 60FPS
	R0_FNX0_INT01_SOL   = 2    // Start of Line (Programmable)
	R0_FNX0_INT02_TMR0  = 4    // Timer 0 Interrupt
	R0_FNX0_INT03_TMR1  = 8    // Timer 1 Interrupt
	R0_FNX0_INT04_TMR2  = 0x10 // Timer 2 Interrupt
	R0_FNX0_INT05_RTC   = 0x20 // Real-Time Clock Interrupt
	R0_FNX0_INT06_FDC   = 0x40 // Floppy Disk Controller
	R0_FNX0_INT07_MOUSE = 0x80 // Mouse Interrupt (INT12 in SuperIO IOspace)
)

const (
	R1_FNX1_INT00_KBD    = 1    // Keyboard Interrupt
	R1_FNX1_INT01_SC0    = 2    // Sprite 2 Sprite Collision
	R1_FNX1_INT02_SC1    = 4    // Sprite 2 Tiles Collision
	R1_FNX1_INT03_COM2   = 8    // Serial Port 2
	R1_FNX1_INT04_COM1   = 0x10 // Serial Port 1
	R1_FNX1_INT05_MPU401 = 0x20 // Midi Controller Interrupt
	R1_FNX1_INT06_LPT    = 0x40 // Parallel Port
	R1_FNX1_INT07_SDCARD = 0x80 // SD Card Controller Interrupt
)

const (
	R2_FNX2_INT00_OPL2R   = 1    // OPL2 Right Channel
	R2_FNX2_INT01_OPL2L   = 2    // OPL2 Left Channel
	R2_FNX2_INT02_BTX_INT = 4    // Beatrix Interrupt (TBD)
	R2_FNX2_INT03_SDMA    = 8    // System DMA
	R2_FNX2_INT04_VDMA    = 0x10 // Video DMA
	R2_FNX2_INT05_DACHP   = 0x20 // DAC Hot Plug
	R2_FNX2_INT06_EXT     = 0x40 // External Expansion
	R2_FNX2_INT07_ALLONE  = 0x80 // ??
)

const (
	R2FMX_FNX2_INT00_OPL3       = 1    // OPL3
	R2FMX_FNX2_INT01_GABE_INT0  = 2    // GABE (INT0) - TBD
	R2FMX_FNX2_INT02_GABE_INT1  = 4    // GABE (INT1) - TBD
	R2FMX_FNX2_INT03_SDMA       = 8    // VICKY_II (INT4)
	R2FMX_FNX2_INT04_VDMA       = 0x10 // VICKY_II (INT5)
	R2FMX_FNX2_INT05_GABE_INT2  = 0x20 // GABE (INT2) - TBD
	R2FMX_FNX2_INT06_EXT        = 0x40 // External Expansion
	R2FMX_FNX2_INT07_SDCARD_INS = 0x80 // SDCARD Insertion
)

const (
	R3FMX_FNX3_INT00_OPN2 = 1    // OPN2
	R3FMX_FNX3_INT01_OPM  = 2    // OPM
	R3FMX_FNX3_INT02_IDE  = 4    // HDD IDE INTERRUPT
	R3FMX_FNX3_INT03_TBD  = 8    // TBD
	R3FMX_FNX3_INT04_TBD  = 0x10 // TBD
	R3FMX_FNX3_INT05_TBD  = 0x20 // GABE (INT2) - TBD
	R3FMX_FNX3_INT06_TBD  = 0x40 // External Expansion
	R3FMX_FNX3_INT07_TBD  = 0x80 // SDCARD Insertion
)
