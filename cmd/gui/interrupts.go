package main

// copy-pasted Devices/Interrupts.cs

const (
	r0_FNX0_INT00_SOF   = 1    // Start of Frame @ 60FPS
	r0_FNX0_INT01_SOL   = 2    // Start of Line (Programmable)
	r0_FNX0_INT02_TMR0  = 4    // Timer 0 Interrupt
	r0_FNX0_INT03_TMR1  = 8    // Timer 1 Interrupt
	r0_FNX0_INT04_TMR2  = 0x10 // Timer 2 Interrupt
	r0_FNX0_INT05_RTC   = 0x20 // Real-Time Clock Interrupt
	r0_FNX0_INT06_FDC   = 0x40 // Floppy Disk Controller
	r0_FNX0_INT07_MOUSE = 0x80 // Mouse Interrupt (INT12 in SuperIO IOspace)
)

const (
	r1_FNX1_INT00_KBD    = 1    // Keyboard Interrupt
	r1_FNX1_INT01_SC0    = 2    // Sprite 2 Sprite Collision
	r1_FNX1_INT02_SC1    = 4    // Sprite 2 Tiles Collision
	r1_FNX1_INT03_COM2   = 8    // Serial Port 2
	r1_FNX1_INT04_COM1   = 0x10 // Serial Port 1
	r1_FNX1_INT05_MPU401 = 0x20 // Midi Controller Interrupt
	r1_FNX1_INT06_LPT    = 0x40 // Parallel Port
	r1_FNX1_INT07_SDCARD = 0x80 // SD Card Controller Interrupt
)

const (
	r2_FNX2_INT00_OPL2R   = 1    // OPL2 Right Channel
	r2_FNX2_INT01_OPL2L   = 2    // OPL2 Left Channel
	r2_FNX2_INT02_BTX_INT = 4    // Beatrix Interrupt (TBD)
	r2_FNX2_INT03_SDMA    = 8    // System DMA
	r2_FNX2_INT04_VDMA    = 0x10 // Video DMA
	r2_FNX2_INT05_DACHP   = 0x20 // DAC Hot Plug
	r2_FNX2_INT06_EXT     = 0x40 // External Expansion
	r2_FNX2_INT07_ALLONE  = 0x80 // ??
)

const (
	r2fmx_FNX2_INT00_OPL3       = 1    // OPL3
	r2fmx_FNX2_INT01_GABE_INT0  = 2    // GABE (INT0) - TBD
	r2fmx_FNX2_INT02_GABE_INT1  = 4    // GABE (INT1) - TBD
	r2fmx_FNX2_INT03_SDMA       = 8    // VICKY_II (INT4)
	r2fmx_FNX2_INT04_VDMA       = 0x10 // VICKY_II (INT5)
	r2fmx_FNX2_INT05_GABE_INT2  = 0x20 // GABE (INT2) - TBD
	r2fmx_FNX2_INT06_EXT        = 0x40 // External Expansion
	r2fmx_FNX2_INT07_SDCARD_INS = 0x80 // SDCARD Insertion
)

const (
	r3fmx_FNX3_INT00_OPN2 = 1    // OPN2
	r3fmx_FNX3_INT01_OPM  = 2    // OPM
	r3fmx_FNX3_INT02_IDE  = 4    // HDD IDE INTERRUPT
	r3fmx_FNX3_INT03_TBD  = 8    // TBD
	r3fmx_FNX3_INT04_TBD  = 0x10 // TBD
	r3fmx_FNX3_INT05_TBD  = 0x20 // GABE (INT2) - TBD
	r3fmx_FNX3_INT06_TBD  = 0x40 // External Expansion
	r3fmx_FNX3_INT07_TBD  = 0x80 // SDCARD Insertion
)
