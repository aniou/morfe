
// 65c816 simulator, based on 65c02-one written for
// NES emulator https://github.com/fogleman/nes
// Copyright  2015-2019  Michael Fogleman, released under MIT license
// Copyright  2019       Piotr Meyer

package cpu65c816

import (
	"log"
	_ "bytes"
	"fmt"

	"github.com/aniou/go65c816/emulator/bus"
	"github.com/aniou/go65c816/lib/mylog"
)

type instructionType struct {
	opcode	byte
	name    string
	mode	byte
	size	byte
	cycles  byte
	proc	func(*stepInfo)
}

var (
	instructions [256]instructionType
)

const CPUFrequency = 14000000 // 14MHz. XXX - fix it and move to platform?

// interrupt types
const (
	_ = iota
	interruptNone
	interruptNMI
	interruptIRQ
)

// addressing modes
// reference: 
// 1 - "Programming the 65816" / WDC 2007           
// 2 - http://6502.org/tutorials/65c816opcodes.html 
// 3- http://datasheets.chipdb.org/Western%20Design/w65c816s.pdf 

const (
	_ = iota
	m_Absolute			// $9876          - p. 288 or 5.2
	m_Absolute_X			// $9876, X       - p. 289 or 5.3
	m_Absolute_Y			// $9876, Y       - p. 290 or 5.3
	m_Accumulator			// A              - p. 296 or 5.6
	m_Immediate			// #$aa           - p. 306 or 5.14
	m_Immediate_flagM		// #$aa or #$aabb - p. 306 or 5.14, flag M dependent size
	m_Immediate_flagX		// #$aa or #$aabb - p. 306 or 5.14, flag X dependent size
	m_Implied			// -              - p. 307 or 5.15
	m_DP				// $12            - p. 298 or 5.7
	m_DP_X				// $12, X         - p. 299 or 5.8
	m_DP_Y				// $12, Y         - p. 300 or 5.8
	m_DP_X_Indirect			// ($12, X)       - p. 301 or 5.11
	m_DP_Indirect			// ($12)          - p. 302 or 5.9
	m_DP_Indirect_Long		// [$12]          - p. 303 or 5.10
	m_DP_Indirect_Y			// ($12), Y       - p. 304 or 5.12
	m_DP_Indirect_Long_Y		// [$12], Y       - p. 305 or 5.13
	m_Absolute_X_Indirect		// ($1234, X)     - p. 291 or 5.5
	m_Absolute_Indirect		// ($1234)        - p. 292 or 5.4
	m_Absolute_Indirect_Long	// [$1234]        - p. 293 or 5.10
	m_Absolute_Long			// $abcdef        - p. 294 or 5.16
	m_Absolute_Long_X		// $abcdex, X     - p. 295 or 5.17
	m_BlockMove			// #$12,#$34      - p. 297 or 5.19  (MVN, MVP)
	m_PC_Relative			// rel8           - p. 308 or 5.18  (BRA)
	m_PC_Relative_Long		// rel16          - p. 309 or 5.18  (BRL)
	//m_Stack_Implied		// #$1234         - p. 310 or 6.8.1 (PEA) -> Immediate
	//m_Stack_DP_Indirect		// $12            - p. 312 or 6.8.1 (PEI) -> DP
	//m_Stack_PC_Relative		// rel16          - p. 316 or 6.8.1 (PER) -> PC_Relative_Long 
	m_Stack_Relative		// $32, S         - p. 324 or 5.20
	m_Stack_Relative_Indirect_Y	// ($32, S), Y    - p. 325 or 5.21  (STACK,S),Y
)

//
//
func (c *CPU) createTable() {
	instructions = [256]instructionType{
		{0x00, "brk", m_Implied,                   1, 8, c.op_brk},	// BRK
		{0x01, "ora", m_DP_X_Indirect,             2, 7, c.op_ora},	// ORA ($10,X)
		{0x02, "cop", m_Immediate,                 2, 8, c.op_cop},	// COP #$12
		{0x03, "ora", m_Stack_Relative,            2, 5, c.op_ora},	// ORA $32,S
		{0x04, "tsb", m_DP,                        2, 7, c.op_tsb},	// TSB $10
		{0x05, "ora", m_DP,                        2, 4, c.op_ora},	// ORA $10
		{0x06, "asl", m_DP,                        2, 7, c.op_asl},	// ASL $10
		{0x07, "ora", m_DP_Indirect_Long,          2, 7, c.op_ora},	// ORA [$10]
		{0x08, "php", m_Implied,                   1, 3, c.op_php},	// PHP
		{0x09, "ora", m_Immediate_flagM,           3, 3, c.op_ora},	// ORA #$54
		{0x0a, "asl", m_Accumulator,               1, 2, c.op_asl},	// ASL
		{0x0b, "phd", m_Implied,                   1, 4, c.op_phd},	// PHD
		{0x0c, "tsb", m_Absolute,                  3, 8, c.op_tsb},	// TSB $9876
		{0x0d, "ora", m_Absolute,                  3, 5, c.op_ora},	// ORA $9876
		{0x0e, "asl", m_Absolute,                  3, 8, c.op_asl},	// ASL $9876
		{0x0f, "ora", m_Absolute_Long,             4, 6, c.op_ora},	// ORA $FEDBCA
		{0x10, "bpl", m_PC_Relative,               2, 2, c.op_bpl},	// BPL LABEL
		{0x11, "ora", m_DP_Indirect_Y,             2, 7, c.op_ora},	// ORA ($10),Y
		{0x12, "ora", m_DP_Indirect,               2, 6, c.op_ora},	// ORA ($10)
		{0x13, "ora", m_Stack_Relative_Indirect_Y, 2, 8, c.op_ora},	// ORA ($32,S),Y
		{0x14, "trb", m_DP,                        2, 7, c.op_trb},	// TRB $10
		{0x15, "ora", m_DP_X,                      2, 5, c.op_ora},	// ORA $10,X
		{0x16, "asl", m_DP_X,                      2, 8, c.op_asl},	// ASL $10,X
		{0x17, "ora", m_DP_Indirect_Long_Y,        2, 7, c.op_ora},	// ORA [$10],Y
		{0x18, "clc", m_Implied,                   1, 2, c.op_clc},	// CLC
		{0x19, "ora", m_Absolute_Y,                3, 6, c.op_ora},	// ORA $9876,Y
		{0x1a, "inc", m_Accumulator,               1, 2, c.op_inc},	// INC
		{0x1b, "tcs", m_Implied,                   1, 2, c.op_tcs},	// TCS
		{0x1c, "trb", m_Absolute,                  3, 8, c.op_trb},	// TRB $9876
		{0x1d, "ora", m_Absolute_X,                3, 6, c.op_ora},	// ORA $9876,X
		{0x1e, "asl", m_Absolute_X,                3, 9, c.op_asl},	// ASL $9876,X
		{0x1f, "ora", m_Absolute_Long_X,           4, 6, c.op_ora},	// ORA $FEDCBA,X
		{0x20, "jsr", m_Absolute,                  3, 6, c.op_jsr},	// JSR $1234
		{0x21, "and", m_DP_X_Indirect,             2, 7, c.op_and},	// AND ($10,X)
		{0x22, "jsl", m_Absolute_Long,             4, 8, c.op_jsl},	// JSL $123456
		{0x23, "and", m_Stack_Relative,            2, 5, c.op_and},	// AND $32,S
		{0x24, "bit", m_DP,                        2, 4, c.op_bit},	// BIT $10
		{0x25, "and", m_DP,                        2, 4, c.op_and},	// AND $10
		{0x26, "rol", m_DP,                        2, 7, c.op_rol},	// ROL $10
		{0x27, "and", m_DP_Indirect_Long,          2, 7, c.op_and},	// AND [$10]
		{0x28, "plp", m_Implied,                   1, 4, c.op_plp},	// PLP
		{0x29, "and", m_Immediate_flagM,           3, 3, c.op_and},	// AND #$54
		{0x2a, "rol", m_Accumulator,               1, 2, c.op_rol},	// ROL
		{0x2b, "pld", m_Implied,                   1, 5, c.op_pld},	// PLD
		{0x2c, "bit", m_Absolute,                  3, 5, c.op_bit},	// BIT $9876
		{0x2d, "and", m_Absolute,                  3, 5, c.op_and},	// AND $9876
		{0x2e, "rol", m_Absolute,                  3, 8, c.op_rol},	// ROL $9876
		{0x2f, "and", m_Absolute_Long,             4, 6, c.op_and},	// AND $FEDBCA
		{0x30, "bmi", m_PC_Relative,               2, 2, c.op_bmi},	// BMI LABEL
		{0x31, "and", m_DP_Indirect_Y,             2, 7, c.op_and},	// AND ($10),Y
		{0x32, "and", m_DP_Indirect,               2, 6, c.op_and},	// AND ($10)
		{0x33, "and", m_Stack_Relative_Indirect_Y, 2, 8, c.op_and},	// AND ($32,S),Y
		{0x34, "bit", m_DP_X,                      2, 5, c.op_bit},	// BIT $10,X
		{0x35, "and", m_DP_X,                      2, 5, c.op_and},	// AND $10,X
		{0x36, "rol", m_DP_X,                      2, 8, c.op_rol},	// ROL $10,X
		{0x37, "and", m_DP_Indirect_Long_Y,        2, 7, c.op_and},	// AND [$10],Y
		{0x38, "sec", m_Implied,                   1, 2, c.op_sec},	// SEC
		{0x39, "and", m_Absolute_Y,                3, 6, c.op_and},	// AND $9876,Y
		{0x3a, "dec", m_Accumulator,               1, 2, c.op_dec},	// DEC
		{0x3b, "tsc", m_Implied,                   1, 2, c.op_tsc},	// TSC
		{0x3c, "bit", m_Absolute_X,                3, 6, c.op_bit},	// BIT $9876,X
		{0x3d, "and", m_Absolute_X,                3, 6, c.op_and},	// AND $9876,X
		{0x3e, "rol", m_Absolute_X,                3, 9, c.op_rol},	// ROL $9876,X
		{0x3f, "and", m_Absolute_Long_X,           4, 6, c.op_and},	// AND $FEDCBA,X
		{0x40, "rti", m_Implied,                   1, 7, c.op_rti},	// RTI
		{0x41, "eor", m_DP_X_Indirect,             2, 7, c.op_eor},	// EOR ($10,X)
		{0x42, "wdm", m_Immediate,                 2, 2, c.op_wdm},	// WDM
		{0x43, "eor", m_Stack_Relative,            2, 5, c.op_eor},	// EOR $32,S
		{0x44, "mvp", m_BlockMove,                 3, 7, c.op_mvp},	// MVP #$12,#$34
		{0x45, "eor", m_DP,                        2, 4, c.op_eor},	// EOR $10
		{0x46, "lsr", m_DP,                        2, 7, c.op_lsr},	// LSR $10
		{0x47, "eor", m_DP_Indirect_Long,          2, 7, c.op_eor},	// EOR [$10]
		{0x48, "pha", m_Implied,                   1, 4, c.op_pha},	// PHA
		{0x49, "eor", m_Immediate_flagM,           3, 3, c.op_eor},	// EOR #$54
		{0x4a, "lsr", m_Accumulator,               1, 2, c.op_lsr},	// LSR
		{0x4b, "phk", m_Implied,                   1, 3, c.op_phk},	// PHK
		{0x4c, "jmp", m_Absolute,                  3, 3, c.op_jmp},	// JMP $1234
		{0x4d, "eor", m_Absolute,                  3, 5, c.op_eor},	// EOR $9876
		{0x4e, "lsr", m_Absolute,                  3, 8, c.op_lsr},	// LSR $9876
		{0x4f, "eor", m_Absolute_Long,             4, 6, c.op_eor},	// EOR $FEDBCA
		{0x50, "bvc", m_PC_Relative,               2, 2, c.op_bvc},	// BVC LABEL
		{0x51, "eor", m_DP_Indirect_Y,             2, 7, c.op_eor},	// EOR ($10),Y
		{0x52, "eor", m_DP_Indirect,               2, 6, c.op_eor},	// EOR ($10)
		{0x53, "eor", m_Stack_Relative_Indirect_Y, 2, 8, c.op_eor},	// EOR ($32,S),Y
		{0x54, "mvn", m_BlockMove,                 3, 7, c.op_mvn},	// MVN #$12,#$34
		{0x55, "eor", m_DP_X,                      2, 5, c.op_eor},	// EOR $10,X
		{0x56, "lsr", m_DP_X,                      2, 8, c.op_lsr},	// LSR $10,X
		{0x57, "eor", m_DP_Indirect_Long_Y,        2, 7, c.op_eor},	// EOR [$10],Y
		{0x58, "cli", m_Implied,                   1, 2, c.op_cli},	// CLI
		{0x59, "eor", m_Absolute_Y,                3, 6, c.op_eor},	// EOR $9876,Y
		{0x5a, "phy", m_Implied,                   1, 4, c.op_phy},	// PHY
		{0x5b, "tcd", m_Implied,                   1, 2, c.op_tcd},	// TCD
		{0x5c, "jmp", m_Absolute_Long,             4, 4, c.op_jmp},	// JMP $FEDCBA
		{0x5d, "eor", m_Absolute_X,                3, 6, c.op_eor},	// EOR $9876,X
		{0x5e, "lsr", m_Absolute_X,                3, 9, c.op_lsr},	// LSR $9876,X
		{0x5f, "eor", m_Absolute_Long_X,           4, 6, c.op_eor},	// EOR $FEDCBA,X
		{0x60, "rts", m_Implied,                   1, 6, c.op_rts},	// RTS
		{0x61, "adc", m_DP_X_Indirect,             2, 7, c.op_adc},	// ADC ($10,X)
		{0x62, "per", m_PC_Relative_Long,          3, 6, c.op_per},	// PER LABEL
		{0x63, "adc", m_Stack_Relative,            2, 5, c.op_adc},	// ADC $32,S
		{0x64, "stz", m_DP,                        2, 4, c.op_stz},	// STZ $10
		{0x65, "adc", m_DP,                        2, 4, c.op_adc},	// ADC $10
		{0x66, "ror", m_DP,                        2, 7, c.op_ror},	// ROR $10
		{0x67, "adc", m_DP_Indirect_Long,          2, 7, c.op_adc},	// ADC [$10]
		{0x68, "pla", m_Implied,		   1, 5, c.op_pla},	// PLA
		{0x69, "adc", m_Immediate_flagM,           3, 3, c.op_adc},	// ADC #$54
		{0x6a, "ror", m_Accumulator,               1, 2, c.op_ror},	// ROR
		{0x6b, "rtl", m_Implied,                   1, 6, c.op_rtl},	// RTL
		{0x6c, "jmp", m_Absolute_Indirect,         3, 5, c.op_jmp},	// JMP ($1234)
		{0x6d, "adc", m_Absolute,                  3, 5, c.op_adc},	// ADC $9876
		{0x6e, "ror", m_Absolute,                  3, 8, c.op_ror},	// ROR $9876
		{0x6f, "adc", m_Absolute_Long,             4, 6, c.op_adc},	// ADC $FEDBCA
		{0x70, "bvs", m_PC_Relative,               2, 2, c.op_bvs},	// BVS LABEL
		{0x71, "adc", m_DP_Indirect_Y,             2, 7, c.op_adc},	// ADC ($10),Y
		{0x72, "adc", m_DP_Indirect,               2, 6, c.op_adc},	// ADC ($10)
		{0x73, "adc", m_Stack_Relative_Indirect_Y, 2, 8, c.op_adc},	// ADC ($32,S),Y
		{0x74, "stz", m_DP_X,                      2, 5, c.op_stz},	// STZ $10,X
		{0x75, "adc", m_DP_X,                      2, 5, c.op_adc},	// ADC $10,X
		{0x76, "ror", m_DP_X,                      2, 8, c.op_ror},	// ROR $10,X
		{0x77, "adc", m_DP_Indirect_Long_Y,        2, 7, c.op_adc},	// ADC [$10],Y
		{0x78, "sei", m_Implied,                   1, 2, c.op_sei},	// SEI
		{0x79, "adc", m_Absolute_Y,                3, 6, c.op_adc},	// ADC $9876,Y
		{0x7a, "ply", m_Implied,                   1, 5, c.op_ply},	// PLY
		{0x7b, "tdc", m_Implied,                   1, 2, c.op_tdc},	// TDC
		{0x7c, "jmp", m_Absolute_X_Indirect,       3, 6, c.op_jmp},	// JMP ($1234,X)
		{0x7d, "adc", m_Absolute_X,                3, 6, c.op_adc},	// ADC $9876,X
		{0x7e, "ror", m_Absolute_X,                3, 9, c.op_ror},	// ROR $9876,X
		{0x7f, "adc", m_Absolute_Long_X,           4, 6, c.op_adc},	// ADC $FEDCBA,X
		{0x80, "bra", m_PC_Relative,               2, 3, c.op_bra},	// BRA LABEL
		{0x81, "sta", m_DP_X_Indirect,             2, 7, c.op_sta},	// STA ($10,X)
		{0x82, "brl", m_PC_Relative_Long,          3, 4, c.op_brl},	// BRL LABEL
		{0x83, "sta", m_Stack_Relative,            2, 5, c.op_sta},	// STA $32,S
		{0x84, "sty", m_DP,                        2, 4, c.op_sty},	// STY $10
		{0x85, "sta", m_DP,                        2, 4, c.op_sta},	// STA $10
		{0x86, "stx", m_DP,                        2, 4, c.op_stx},	// STX $10
		{0x87, "sta", m_DP_Indirect_Long,          2, 7, c.op_sta},	// STA [$10]
		{0x88, "dey", m_Implied,                   1, 2, c.op_dey},	// DEY
		{0x89, "bit", m_Immediate_flagM,           3, 3, c.op_bit},	// BIT #$54
		{0x8a, "txa", m_Implied,                   1, 2, c.op_txa},	// TXA
		{0x8b, "phb", m_Implied,                   1, 3, c.op_phb},	// PHB
		{0x8c, "sty", m_Absolute,                  3, 5, c.op_sty},	// STY $9876
		{0x8d, "sta", m_Absolute,                  3, 5, c.op_sta},	// STA $9876
		{0x8e, "stx", m_Absolute,                  3, 5, c.op_stx},	// STX $9876
		{0x8f, "sta", m_Absolute_Long,             4, 6, c.op_sta},	// STA $FEDBCA
		{0x90, "bcc", m_PC_Relative,               2, 2, c.op_bcc},	// BCC LABEL
		{0x91, "sta", m_DP_Indirect_Y,             2, 7, c.op_sta},	// STA ($10),Y
		{0x92, "sta", m_DP_Indirect,               2, 6, c.op_sta},	// STA ($10)
		{0x93, "sta", m_Stack_Relative_Indirect_Y, 2, 8, c.op_sta},	// STA ($32,S),Y
		{0x94, "sty", m_DP_X,                      2, 5, c.op_sty},	// STY $10,X
		{0x95, "sta", m_DP_X,                      2, 5, c.op_sta},	// STA $10,X
		{0x96, "stx", m_DP_Y,                      2, 5, c.op_stx},	// STX $10,Y
		{0x97, "sta", m_DP_Indirect_Long_Y,        2, 7, c.op_sta},	// STA [$10],Y
		{0x98, "tya", m_Implied,                   1, 2, c.op_tya},	// TYA
		{0x99, "sta", m_Absolute_Y,                3, 6, c.op_sta},	// STA $9876,Y
		{0x9a, "txs", m_Implied,                   1, 2, c.op_txs},	// TXS
		{0x9b, "txy", m_Implied,                   1, 2, c.op_txy},	// TXY
		{0x9c, "stz", m_Absolute,                  3, 5, c.op_stz},	// STZ $9876
		{0x9d, "sta", m_Absolute_X,                3, 6, c.op_sta},	// STA $9876,X
		{0x9e, "stz", m_Absolute_X,                3, 6, c.op_stz},	// STZ $9876,X
		{0x9f, "sta", m_Absolute_Long_X,           4, 6, c.op_sta},	// STA $FEDCBA,X
		{0xa0, "ldy", m_Immediate_flagX,           3, 3, c.op_ldy},	// LDY #$54
		{0xa1, "lda", m_DP_X_Indirect,             2, 7, c.op_lda},	// LDA ($10,X)
		{0xa2, "ldx", m_Immediate_flagX,           3, 3, c.op_ldx},	// LDX #$54
		{0xa3, "lda", m_Stack_Relative,            2, 5, c.op_lda},	// LDA $32,S
		{0xa4, "ldy", m_DP,                        2, 4, c.op_ldy},	// LDY $10
		{0xa5, "lda", m_DP,                        2, 4, c.op_lda},	// LDA $10
		{0xa6, "ldx", m_DP,                        2, 4, c.op_ldx},	// LDX $10
		{0xa7, "lda", m_DP_Indirect_Long,          2, 7, c.op_lda},	// LDA [$10]
		{0xa8, "tay", m_Implied,                   1, 2, c.op_tay},	// TAY
		{0xa9, "lda", m_Immediate_flagM,           3, 3, c.op_lda},	// LDA #$54
		{0xaa, "tax", m_Implied,                   1, 2, c.op_tax},	// TAX
		{0xab, "plb", m_Implied,                   1, 4, c.op_plb},	// PLB
		{0xac, "ldy", m_Absolute,                  3, 5, c.op_ldy},	// LDY $9876
		{0xad, "lda", m_Absolute,                  3, 5, c.op_lda},	// LDA $9876
		{0xae, "ldx", m_Absolute,                  3, 5, c.op_ldx},	// LDX $9876
		{0xaf, "lda", m_Absolute_Long,             4, 6, c.op_lda},	// LDA $FEDBCA
		{0xb0, "bcs", m_PC_Relative,               2, 2, c.op_bcs},	// BCS LABEL
		{0xb1, "lda", m_DP_Indirect_Y,             2, 7, c.op_lda},	// LDA ($10),Y
		{0xb2, "lda", m_DP_Indirect,               2, 6, c.op_lda},	// LDA ($10)
		{0xb3, "lda", m_Stack_Relative_Indirect_Y, 2, 8, c.op_lda},	// LDA ($32,S),Y
		{0xb4, "ldy", m_DP_X,                      2, 5, c.op_ldy},	// LDY $10,X
		{0xb5, "lda", m_DP_X,                      2, 5, c.op_lda},	// LDA $10,X
		{0xb6, "ldx", m_DP_Y,                      2, 5, c.op_ldx},	// LDX $10,Y
		{0xb7, "lda", m_DP_Indirect_Long_Y,        2, 7, c.op_lda},	// LDA [$10],Y
		{0xb8, "clv", m_Implied,                   1, 2, c.op_clv},	// CLV
		{0xb9, "lda", m_Absolute_Y,                3, 6, c.op_lda},	// LDA $9876,Y
		{0xba, "tsx", m_Implied,                   1, 2, c.op_tsx},	// TSX
		{0xbb, "tyx", m_Implied,                   1, 2, c.op_tyx},	// TYX
		{0xbc, "ldy", m_Absolute_X,                3, 6, c.op_ldy},	// LDY $9876,X
		{0xbd, "lda", m_Absolute_X,                3, 6, c.op_lda},	// LDA $9876,X
		{0xbe, "ldx", m_Absolute_Y,                3, 6, c.op_ldx},	// LDX $9876,Y
		{0xbf, "lda", m_Absolute_Long_X,           4, 6, c.op_lda},	// LDA $FEDCBA,X
		{0xc0, "cpy", m_Immediate_flagX,           3, 3, c.op_cpy},	// CPY #$54
		{0xc1, "cmp", m_DP_X_Indirect,             2, 7, c.op_cmp},	// CMP ($10,X)
		{0xc2, "rep", m_Immediate,                 2, 3, c.op_rep},	// REP #$12
		{0xc3, "cmp", m_Stack_Relative,            2, 5, c.op_cmp},	// CMP $32,S
		{0xc4, "cpy", m_DP,                        2, 4, c.op_cpy},	// CPY $10
		{0xc5, "cmp", m_DP,                        2, 4, c.op_cmp},	// CMP $10
		{0xc6, "dec", m_DP,                        2, 7, c.op_dec},	// DEC $10
		{0xc7, "cmp", m_DP_Indirect_Long,          2, 7, c.op_cmp},	// CMP [$10]
		{0xc8, "iny", m_Implied,                   1, 2, c.op_iny},	// INY
		{0xc9, "cmp", m_Immediate_flagM,           3, 3, c.op_cmp},	// CMP #$54
		{0xca, "dex", m_Implied,                   1, 2, c.op_dex},	// DEX
		{0xcb, "wai", m_Implied,                   1, 3, c.wai},	// WAI
		{0xcc, "cpy", m_Absolute,                  3, 5, c.op_cpy},	// CPY $9876
		{0xcd, "cmp", m_Absolute,                  3, 5, c.op_cmp},	// CMP $9876
		{0xce, "dec", m_Absolute,                  3, 8, c.op_dec},	// DEC $9876
		{0xcf, "cmp", m_Absolute_Long,             4, 6, c.op_cmp},	// CMP $FEDBCA
		{0xd0, "bne", m_PC_Relative,               2, 2, c.op_bne},	// BNE LABEL
		{0xd1, "cmp", m_DP_Indirect_Y,             2, 7, c.op_cmp},	// CMP ($10),Y
		{0xd2, "cmp", m_DP_Indirect,               2, 6, c.op_cmp},	// CMP ($10)
		{0xd3, "cmp", m_Stack_Relative_Indirect_Y, 2, 8, c.op_cmp},	// CMP ($32,S),Y
		{0xd4, "pei", m_DP,                        2, 6, c.op_pei},	// PEI $12
		{0xd5, "cmp", m_DP_X,                      2, 5, c.op_cmp},	// CMP $10,X
		{0xd6, "dec", m_DP_X,                      2, 8, c.op_dec},	// DEC $10,X
		{0xd7, "cmp", m_DP_Indirect_Long_Y,        2, 7, c.op_cmp},	// CMP [$10],Y
		{0xd8, "cld", m_Implied,                   1, 2, c.op_cld},	// CLD
		{0xd9, "cmp", m_Absolute_Y,                3, 6, c.op_cmp},	// CMP $9876,Y
		{0xda, "phx", m_Implied,                   1, 4, c.op_phx},	// PHX
		{0xdb, "stp", m_Implied,                   1, 3, c.stp},	// STP
		{0xdc, "jmp", m_Absolute_Indirect_Long,    3, 6, c.op_jmp},	// JMP [$1234]
		{0xdd, "cmp", m_Absolute_X,                3, 6, c.op_cmp},	// CMP $9876,X
		{0xde, "dec", m_Absolute_X,                3, 9, c.op_dec},	// DEC $9876,X
		{0xdf, "cmp", m_Absolute_Long_X,           4, 6, c.op_cmp},	// CMP $FEDCBA,X
		{0xe0, "cpx", m_Immediate_flagX,           3, 3, c.op_cpx},	// CPX #$54
		{0xe1, "sbc", m_DP_X_Indirect,             2, 7, c.op_sbc},	// SBC ($10,X)
		{0xe2, "sep", m_Immediate,                 2, 3, c.op_sep},	// SEP #$12
		{0xe3, "sbc", m_Stack_Relative,            2, 5, c.op_sbc},	// SBC $32,S
		{0xe4, "cpx", m_DP,                        2, 4, c.op_cpx},	// CPX $10
		{0xe5, "sbc", m_DP,                        2, 4, c.op_sbc},	// SBC $10
		{0xe6, "inc", m_DP,                        2, 7, c.op_inc},	// INC $10
		{0xe7, "sbc", m_DP_Indirect_Long,          2, 7, c.op_sbc},	// SBC [$10]
		{0xe8, "inx", m_Implied,                   1, 2, c.op_inx},	// INX
		{0xe9, "sbc", m_Immediate_flagM,           3, 3, c.op_sbc},	// SBC #$54
		{0xea, "nop", m_Implied,                   1, 2, c.op_nop},	// NOP
		{0xeb, "xba", m_Implied,                   1, 3, c.op_xba},	// XBA
		{0xec, "cpx", m_Absolute,                  3, 5, c.op_cpx},	// CPX $9876
		{0xed, "sbc", m_Absolute,                  3, 5, c.op_sbc},	// SBC $9876
		{0xee, "inc", m_Absolute,                  3, 8, c.op_inc},	// INC $9876
		{0xef, "sbc", m_Absolute_Long,             4, 6, c.op_sbc},	// SBC $FEDBCA
		{0xf0, "beq", m_PC_Relative,               2, 2, c.op_beq},	// BEQ LABEL
		{0xf1, "sbc", m_DP_Indirect_Y,             2, 7, c.op_sbc},	// SBC ($10),Y
		{0xf2, "sbc", m_DP_Indirect,               2, 6, c.op_sbc},	// SBC ($10)
		{0xf3, "sbc", m_Stack_Relative_Indirect_Y, 2, 8, c.op_sbc},	// SBC ($32,S),Y
		{0xf4, "pea", m_Immediate,                 3, 5, c.op_pea},	// PEA #$1234
		{0xf5, "sbc", m_DP_X,                      2, 5, c.op_sbc},	// SBC $10,X
		{0xf6, "inc", m_DP_X,                      2, 8, c.op_inc},	// INC $10,X
		{0xf7, "sbc", m_DP_Indirect_Long_Y,        2, 7, c.op_sbc},	// SBC [$10],Y
		{0xf8, "sed", m_Implied,                   1, 2, c.op_sed},	// SED
		{0xf9, "sbc", m_Absolute_Y,                3, 6, c.op_sbc},	// SBC $9876,Y
		{0xfa, "plx", m_Implied,                   1, 5, c.op_plx},	// PLX
		{0xfb, "xce", m_Implied,                   1, 2, c.op_xce},	// XCE
		{0xfc, "jsr", m_Absolute_X_Indirect,       3, 8, c.op_jsr},	// JSR ($1234,X)
		{0xfd, "sbc", m_Absolute_X,                3, 6, c.op_sbc},	// SBC $9876,X
		{0xfe, "inc", m_Absolute_X,                3, 9, c.op_inc},	// INC $9876,X
		{0xff, "sbc", m_Absolute_Long_X,           4, 6, c.op_sbc},	// SBC $FEDCBA,X
	}
}

type CPU struct {
	Bus	*bus.Bus

	// additional emulator variables
	AllCycles uint64 // total number of cycles of CPU instance
	Cycles    byte   // number of cycles for this step
	stepPC    uint16 // how many bytes should PC be increased in this step?
	abort     bool   // temporary flag to determine that cpu should stop

	// previous register's value exists for debugging purposes
	PRK	  byte   // previous value of program banK register
	PPC       uint16 // previous value Program Counter 

	// 65c816 registers
	PC        uint16 // Program Counter
	SP        uint16 // Stack Pointer

	RA        uint16 // Accumulator
	RX        uint16 // X register
	RY        uint16 // Y register

	RAh	  byte   // Accumulator - upper part
	RAl	  byte   // Accumulator - lower part
	RXl       byte   // X register - lower part
	RYl       byte   // Y register - lower part

	RDBR	  byte   // Data Bank Register
	RD	  uint16 // Direct register
	RK	  byte   // program banK register

	N         byte   // Negative flag
	V         byte   // oVerflow flag
	M         byte   // accumulator and Memory width flag
	X         byte   // indeX register width flag
	D         byte   // Decimal mode flag
	I         byte   // Interrupt disable flag
	Z         byte   // Zero flag
	C         byte   // Carry flag

	B         byte   // Break flag
	E         byte   // Emulation mode flag

	interrupt byte   // interrupt type to perform
	stall     int    // number of cycles to stall
	table     [256]func(*stepInfo)
}



func New(bus *bus.Bus) (*CPU, error) {
	cpu    := CPU{Bus: bus}
	//cpu.LogBuf = bytes.Buffer{}
	cpu.createTable()
	//cpu.Reset()
	mylog.Logger.Log("cpu: 65c816 initialized")
	return &cpu, nil
}

// Reset resets the CPU to its initial powerup state
func (cpu *CPU) Reset() {
	//cpu.E    = 1     - should be
	cpu.SP   = 0x01ff
	cpu.RD   = 0x0000
	cpu.RK   = 0x0000
	cpu.RDBR = 0x0000
	//cpu.PC   = cpu.Read16(0xFFFC)
	cpu.PC   = cpu.nRead16_cross(0x00, 0xFFFC)
	cpu.SetFlags(0x34)
    cpu.abort = false
}

/* ====================================================================
 *
 *
 *
 * ====================================================================
 */

// pagesDiffer returns true if the two addresses reference different pages
func pagesDiffer(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// addBranchCycles adds a cycle for taking a branch and adds another cycle
// if the branch jumps to a new page
// note 5 and 6 in [3]
func (cpu *CPU) addBranchCycles(info *stepInfo) {
	cpu.Cycles++
	if pagesDiffer(cpu.PC+2, info.addr) {  // at this moment PC points to jump and addr contains new addr
		cpu.Cycles++
	}
}

func (cpu *CPU) compare8(a, b byte) {
	cpu.setZN8(a - b)
	if a >= b {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}
func (cpu *CPU) compare16(a, b uint16) {
	cpu.setZN16(a - b)
	if a >= b {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// ----------------------------------------------------------------

func (cpu *CPU) nWrite(bank byte, addr uint16, value byte) {
	cpu.Bus.EaWrite(uint32(bank) << 16 | uint32(addr), value)
}

// probably not needed...
func (cpu *CPU) nWrite16_wrap(bank byte, addr uint16, value uint16) {
	bank32 := uint32(bank) << 16
	ll     := byte(value)
	hh     := byte(value >> 8)
	cpu.Bus.EaWrite(bank32 | uint32(addr)  , ll)
	cpu.Bus.EaWrite(bank32 | uint32(addr+1), hh)
}

func (cpu *CPU) nWrite16_cross(bank byte, addr uint16, value uint16) {
	ea     := uint32(bank) << 16 | uint32(addr)
	ll     := byte(value)
	hh     := byte(value >> 8)
	cpu.Bus.EaWrite(ea,   ll)
	cpu.Bus.EaWrite(ea+1, hh)
}

func (cpu *CPU) nRead(bank byte, addr uint16) byte {
	return cpu.Bus.EaRead(uint32(bank) << 16 | uint32(addr))
}

func (cpu *CPU) nRead16_wrap(bank byte, addr uint16) uint16 {
	bank32 := uint32(bank) << 16
	ll     := cpu.Bus.EaRead(bank32 | uint32(addr))
	hh     := cpu.Bus.EaRead(bank32 | uint32(addr+1))
	return uint16(hh) << 8 | uint16(ll)
}

func (cpu *CPU) nRead24_wrap(bank byte, addr uint16) uint32 {
	bank32 := uint32(bank) << 16
	ll     := cpu.Bus.EaRead(bank32 | uint32(addr))
	mm     := cpu.Bus.EaRead(bank32 | uint32(addr+1))
	hh     := cpu.Bus.EaRead(bank32 | uint32(addr+2))
	return uint32(hh) << 16 | uint32(mm) << 8 | uint32(ll)
}

func (cpu *CPU) nRead16_cross(bank byte, addr uint16) uint16 {
	ea     := uint32(bank) << 16 | uint32(addr)
	ll     := cpu.Bus.EaRead(ea)
	hh     := cpu.Bus.EaRead((ea+1) & 0x00ffffff)        // wrap on 24bits
	return uint16(hh) << 8 | uint16(ll)
}




func (cpu *CPU) cmdRead(info *stepInfo) byte {
	switch info.mode {

	case m_Accumulator:
		return cpu.RAl

	case m_Immediate, m_Immediate_flagM, m_Immediate_flagX:
		//return cpu.Bus.EaRead(uint32(cpu.RK) << 16 | uint32(info.addr)) // addr=PC+1
		return cpu.nRead(cpu.RK, info.addr)

	case m_DP, m_DP_X, m_DP_Y, m_Stack_Relative:
		return cpu.Bus.EaRead(uint32(info.addr))  // info.addr is uint16

        case m_DP_Indirect_Long,
             m_DP_Indirect_Long_Y,
             m_Absolute_Long,
             m_Absolute_Long_X,
             m_Absolute_X,
             m_Absolute_Y,
             m_Stack_Relative_Indirect_Y:
                return cpu.Bus.EaRead(info.ea)

        case m_Absolute,
             m_DP_X_Indirect,
             m_DP_Indirect,
             m_DP_Indirect_Y:
		return cpu.nRead(cpu.RDBR, info.addr)


	default:
		//fmt.Fprintf(&cpu.LogBuf, "cmdRead8: unknown mode %v\n", info.mode)
		mylog.Logger.Log(fmt.Sprintf("cmdRead8: unknown mode %v", info.mode))

		return 0
	}

}

func (cpu *CPU) cmdRead16(info *stepInfo) uint16 {
	switch info.mode {

	case m_Accumulator:
		return cpu.RA

	case m_Immediate, m_Immediate_flagM, m_Immediate_flagX:
		return cpu.nRead16_wrap(cpu.RK, info.addr)

	case m_DP, m_DP_X, m_DP_Y, m_Stack_Relative:
		return cpu.nRead16_wrap(0x00, info.addr)

	case m_DP_Indirect_Long,
             m_DP_Indirect_Long_Y,
             m_Absolute_Long,
             m_Absolute_Long_X,
             m_Absolute_X,
             m_Absolute_Y,
             m_Absolute_X_Indirect,
             m_Stack_Relative_Indirect_Y:
                ll := cpu.Bus.EaRead(info.ea)                  // todo - zastapic to jakos?
		hh := cpu.Bus.EaRead(info.ea+1)
		return uint16(hh) << 8 | uint16(ll)

	case m_Absolute,
             m_DP_X_Indirect,
             m_DP_Indirect,
             m_DP_Indirect_Y:
                     return cpu.nRead16_cross(cpu.RDBR, info.addr)

	default:
		//fmt.Fprintf(&cpu.LogBuf, "cmdRead16: unknown mode %v\n", info.mode)
		mylog.Logger.Log(fmt.Sprintf("cmdRead16: unknown mode %v", info.mode))
		return 0
	}

}


func (cpu *CPU) cmdWrite(info *stepInfo, value byte) {
	switch info.mode {

	case m_DP, m_DP_X, m_DP_Y, m_Stack_Relative:
		cpu.Bus.EaWrite(uint32(info.addr), value)  // info.addr is uint16

        case m_DP_Indirect_Long,
             m_DP_Indirect_Long_Y,
             m_Absolute_Long,
             m_Absolute_Long_X,
             m_Absolute_X,
             m_Absolute_Y,
             m_Stack_Relative_Indirect_Y:
                cpu.Bus.EaWrite(info.ea, value)

        case m_Absolute,
             m_DP_X_Indirect,
             m_DP_Indirect,
             m_DP_Indirect_Y:
		cpu.nWrite(cpu.RDBR, info.addr, value)


	default:
		mylog.Logger.Log(fmt.Sprintf("cmdWrite: unknown mode %v", info.mode))
	}

}


func (cpu *CPU) cmdWrite16(info *stepInfo, value uint16) {
	switch info.mode {

	// I'm not so sure though that m_DP modes wraps at bank0
	case m_DP, m_DP_X, m_DP_Y, m_Stack_Relative:
		cpu.nWrite16_wrap(0x00, info.addr, value)

	case m_DP_Indirect_Long,
             m_DP_Indirect_Long_Y,
             m_Absolute_Long,
             m_Absolute_Long_X,
             m_Absolute_X,
             m_Absolute_Y,
             m_Stack_Relative_Indirect_Y:
		ll     := byte(value)
		hh     := byte(value >> 8)
		cpu.Bus.EaWrite(info.ea,   ll)
		cpu.Bus.EaWrite(info.ea+1, hh)

	case m_Absolute,
             m_DP_X_Indirect,
             m_DP_Indirect,
             m_DP_Indirect_Y:
                cpu.nWrite16_cross(cpu.RDBR, info.addr, value)

	default:
		mylog.Logger.Log(fmt.Sprintf("cmdWrite16: unknown mode %v", info.mode))
	}
}


// push pushes a byte onto the stack
func (cpu *CPU) push(value byte) {
	//cpu.Write(0x100|uint16(cpu.SP), value)
	//cpu.Write(cpu.SP, value)
	cpu.nWrite(0x00, cpu.SP, value)
	cpu.SP--

	if cpu.E == 1 {
		cpu.SP = cpu.SP & 0x00FF
		cpu.SP = cpu.SP | 0x1000
	}
}

// pull pops a byte from the stack
func (cpu *CPU) pull() byte {
	cpu.SP++

	if cpu.E == 1 {
		cpu.SP = cpu.SP & 0x00FF
		cpu.SP = cpu.SP | 0x1000
	}
	//return cpu.Read(0x100 | uint16(cpu.SP))
	//return cpu.Read(cpu.SP)
	return cpu.nRead(0x00, cpu.SP)
}

// push16 pushes two bytes onto the stack
func (cpu *CPU) push16(value uint16) {
	hi := byte(value >> 8)
	lo := byte(value & 0xFF)
	cpu.push(hi)
	cpu.push(lo)
}

// pull16 pops two bytes from the stack
func (cpu *CPU) pull16() uint16 {
	lo := uint16(cpu.pull())
	hi := uint16(cpu.pull())
	return hi<<8 | lo
}

// Flags returns the processor status flags
func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.C << 0
	flags |= cpu.Z << 1
	flags |= cpu.I << 2
	flags |= cpu.D << 3
	flags |= cpu.X << 4
	flags |= cpu.M << 5
	flags |= cpu.V << 6
	flags |= cpu.N << 7
	return flags
}

// SetFlags sets the processor status flags
//
// see page 51 for rule for register translation (16/8)
func (cpu *CPU) SetFlags(flags byte) {
	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1

	if cpu.E == 1 {
		cpu.X = 1 // maybe it should be changed to 
		cpu.M = 1 // "no change if E==1" - X/M are forced when E 0->1
	} else {
		oldX := cpu.X
		oldM := cpu.M
		cpu.X = (flags >> 4) & 1
		cpu.M = (flags >> 5) & 1

		if oldX != cpu.X {
			cpu.ChangeRegisterSizes_X()
		}

		if oldM != cpu.M {
			cpu.ChangeRegisterSizes_M()
		}
	}

	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
}

func (cpu *CPU) ChangeRegisterSizes_X() {
	if cpu.X == 1 {
		cpu.RX  = cpu.RX & 0x00ff
		cpu.RXl = uint8(cpu.RX)

		cpu.RY  = cpu.RY & 0x00ff
		cpu.RYl = uint8(cpu.RY)
	} else {
		cpu.RX  = uint16(cpu.RXl)
		cpu.RY  = uint16(cpu.RYl)
	}
}

func (cpu *CPU) ChangeRegisterSizes_M() {
	if cpu.M == 1 {
		cpu.RAl = byte(cpu.RA)
		cpu.RAh = byte(cpu.RA >> 8)
		//fmt.Fprintf(&cpu.LogBuf, "M 0->1 ")
		//fmt.Fprintf(&cpu.LogBuf, "M 0->1 ")
	} else {
		cpu.RA  = uint16(cpu.RAh) << 8 | uint16(cpu.RAl)
		//fmt.Fprintf(&cpu.LogBuf, "M 1->0 ")
	}
	//fmt.Fprintf(&cpu.LogBuf, " $%04x $%02x $%02x\n", cpu.RA, cpu.RAh, cpu.RAl)
}


// setZ sets the zero flag if the argument (8bit) is zero
func (cpu *CPU) setZ8(value byte) {
	if value == 0 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}

// setZ sets the zero flag if the argument (16bit) is zero
func (cpu *CPU) setZ16(value uint16) {
	if value == 0 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}


// setN sets the negative flag if the argument (8bit) is negative (high bit is set)
func (cpu *CPU) setN8(value byte) {
	if value&0x80 != 0 {
		cpu.N = 1
	} else {
		cpu.N = 0
	}
}
// setN sets the negative flag if the argument (16bit) is negative (high bit is set)
func (cpu *CPU) setN16(value uint16) {
	if value&0x8000 != 0 {
		cpu.N = 1
	} else {
		cpu.N = 0
	}
}

// setZN sets the zero flag and the negative flag
func (cpu *CPU) setZN8(value byte) {
	cpu.setZ8(value)
	cpu.setN8(value)
}

// setZN sets the zero flag and the negative flag
func (cpu *CPU) setZN16(value uint16) {
	cpu.setZ16(value)
	cpu.setN16(value)
}

// triggerNMI causes a non-maskable interrupt to occur on the next cycle
func (cpu *CPU) triggerNMI() {
	cpu.interrupt = interruptNMI
}

// triggerIRQ causes an IRQ interrupt to occur on the next cycle
func (cpu *CPU) TriggerIRQ() {
	if cpu.I == 0 {
		cpu.interrupt = interruptIRQ
	}
}

/* ====================================================================
 *
 *
 *
 * ====================================================================
 */

// stepInfo contains information that the instruction functions use
type stepInfo struct {
	ea      uint32  // effective addres    - 24 bit in 65c816
	addr	uint16  // address within bank - used in place of ea in some modes
	pc      uint16
	mode    byte
}


// Step executes a single CPU instruction
func (cpu *CPU) Step() (int, bool) {

	//if cpu.stall > 0 {
	//	cpu.stall--
	//	return 1
	//}

	//cycles := cpu.Cycles

	switch cpu.interrupt {
	case interruptNMI:
		cpu.nmi()
	case interruptIRQ:
		cpu.irq()
	}
	cpu.interrupt = interruptNone

	cpu.PPC    = cpu.PC
	cpu.PRK    = cpu.RK
	opcode    := cpu.nRead(cpu.RK, cpu.PC)
	mode      := instructions[opcode].mode
        cpu.stepPC = uint16(instructions[opcode].size)
        cpu.Cycles = instructions[opcode].cycles



	// addressing mode calculation
	//
	var pageCrossed bool
	var arg8        byte     // temporary arg
	var arg16       uint16   // temporary arg
	var addr        uint16   // final, 16-bit addre
	var ea          uint32   // final, effective addres

	switch mode {

	// $9876          - p. 288 or 5.2
	case m_Absolute:
		addr  = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)


	// $9876, X       - p. 289 or 5.3
	case m_Absolute_X:
		arg16 = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			ea    = (uint32(cpu.RDBR) << 16 | uint32(arg16)) + uint32(cpu.RXl)
			pageCrossed = pagesDiffer(arg16, arg16 + uint16(cpu.RXl))
		} else {
			ea    = (uint32(cpu.RDBR) << 16 | uint32(arg16)) + uint32(cpu.RX)
			pageCrossed = pagesDiffer(arg16, arg16 + cpu.RX)
		}
		//fmt.Fprintf(&cpu.LogBuf, "m_Absolute_X: arg16 %04x ea $%06x\n", arg16, ea)


	// $9876, Y       - p. 290 or 5.3
	case m_Absolute_Y:
		arg16 = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			ea    = (uint32(cpu.RDBR) << 16 | uint32(arg16)) + uint32(cpu.RYl)
			pageCrossed = pagesDiffer(arg16, arg16 + uint16(cpu.RYl))
		} else {
			ea    = (uint32(cpu.RDBR) << 16 | uint32(arg16)) + uint32(cpu.RY)
			pageCrossed = pagesDiffer(arg16, arg16 + cpu.RY)
		}
		//fmt.Fprintf(&cpu.LogBuf, "m_Absolute_Y: arg16 %04x ea $%06x\n", arg16, ea)


	// A              - p. 296 or 5.6
	case m_Accumulator:
		ea  = 0  // effective no-op


	// #$aa           - p. 306 or 5.14
	case m_Immediate:
		addr = cpu.PC + 1
		ea  = uint32(cpu.RK) << 16 | uint32(cpu.PC + 1)  // XXX - required for old SEP implementation
		//fmt.Fprintf(&cpu.LogBuf, "m_Immediate: ea $%06x\n", ea)


	// #$aa or #$aabb - p. 306 or 5.14, flag M size dependent
	case m_Immediate_flagM:
		addr = cpu.PC + 1
		cpu.stepPC -= uint16(cpu.M)   // M = 0 for 16 bit or 1 for 8


	// #$aa or #$aabb - p. 306 or 5.14, flag X size dependent
	case m_Immediate_flagX:
		addr = cpu.PC + 1
		cpu.stepPC -= uint16(cpu.X)   // X = 0 for 16 bit or 1 for 8


	// -              - p. 307 or 5.15
	case m_Implied:
		ea  = 0 // effective no-op


	// $12            - p. 298 or 5.7
	case m_DP:
		arg8    = cpu.nRead(cpu.RK, cpu.PC + 1)
		addr    = uint16(arg8) + cpu.RD

	// $12, X         - p. 299 or 5.8
	case m_DP_X:
		arg8    = cpu.nRead(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			addr = uint16(arg8) + uint16(cpu.RXl) + cpu.RD
		} else {
			addr = uint16(arg8) +        cpu.RX   + cpu.RD
		}


	// $12, Y         - p. 300 or 5.8
	case m_DP_Y:
		arg8    = cpu.nRead(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			addr = uint16(arg8) + uint16(cpu.RYl) + cpu.RD
		} else {
			addr = uint16(arg8) +        cpu.RY   + cpu.RD
		}


	// ($12, X)       - p. 301 or 5.11
	case m_DP_X_Indirect:
		arg8     = cpu.nRead(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			addr = cpu.nRead16_wrap(0x00, uint16(arg8) + uint16(cpu.RXl) + cpu.RD)
		} else {
			addr = cpu.nRead16_wrap(0x00, uint16(arg8) +        cpu.RX   + cpu.RD)
		}
		//fmt.Fprintf(&cpu.LogBuf, "m_DP_X_Indirect: addr $%04x\n", addr)


	// ($12)          - p. 302 or 5.9
	case m_DP_Indirect:
		arg8   = cpu.nRead(cpu.RK, cpu.PC + 1)
		addr   = cpu.nRead16_wrap(0x00, uint16(arg8) + cpu.RD)


	// [$12]          - p. 303 or 5.10
	case m_DP_Indirect_Long:
		// address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1) + cpu.RX))
		arg8   = cpu.nRead(cpu.RK, cpu.PC + 1)
		ea     = cpu.nRead24_wrap(0x00, uint16(arg8) + cpu.RD)


	// ($12), Y       - p. 304 or 5.12
	case m_DP_Indirect_Y:
		arg8   = cpu.nRead(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			addr = cpu.nRead16_wrap(0, uint16(arg8) + cpu.RD) + uint16(cpu.RYl)
			pageCrossed = pagesDiffer(addr-uint16(cpu.RYl), addr)
		} else {
			addr = cpu.nRead16_wrap(0, uint16(arg8) + cpu.RD) +        cpu.RY
			pageCrossed = pagesDiffer(addr-cpu.RY, addr)
		}


	// [$12], Y       - p. 305 or 5.13
	case m_DP_Indirect_Long_Y:
		arg8   = cpu.nRead(cpu.RK, cpu.PC + 1)
		ea     = cpu.nRead24_wrap(0x00, uint16(arg8) + cpu.RD)
		if cpu.X == 1 {
			ea = ea + uint32(cpu.RYl)
		} else {
			ea = ea + uint32(cpu.RY)
		}


	// ($1234, X)     - p. 291 or 5.5
	case m_Absolute_X_Indirect:
		arg16  = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			arg16 = arg16 + uint16(cpu.RXl)
		} else {
			arg16 = arg16 + cpu.RX
		}
		addr = cpu.nRead16_wrap(cpu.RK, arg16)
		ea   = uint32(cpu.RK) << 16 | uint32(arg16)
		//fmt.Fprintf(&cpu.LogBuf, "m_Absolute_X_Indirect: arg16 $%04x addr $%04x ea $%06x\n", arg16, addr, ea)


	// ($1234)        - p. 292 or 5.4
	case m_Absolute_Indirect:
		addr  = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)
		//addr   = cpu.nRead16_wrap(0x00,   arg16)


	// [$1234]        - p. 293 or 5.10
	case m_Absolute_Indirect_Long:
		addr   = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)
		//ea     = cpu.nRead24_wrap(0x00, arg16)
		//fmt.Fprintf(&cpu.LogBuf, "m_Absolute_Indirect_Long: arg16 $%04x ea $%06x\n", arg16, ea)


	// $abcdef        - p. 294 or 5.16
	case m_Absolute_Long:
		ea  = cpu.nRead24_wrap(cpu.RK, cpu.PC + 1)
		//fmt.Fprintf(&cpu.LogBuf, "m_Absolute_Long: ea $%06x\n", ea)


	// $abcdex, X     - p. 295 or 5.17
	case m_Absolute_Long_X:
		ea  = cpu.nRead24_wrap(cpu.RK, cpu.PC + 1)
		if cpu.X == 1 {
			ea = ea + uint32(cpu.RXl)
		} else {
			ea = ea + uint32(cpu.RX)
		}
		//fmt.Fprintf(&cpu.LogBuf, "m_Absolute_Long: ea $%06x\n", ea)


	// rel8           - p. 308 or 5.18 (BRA)
	case m_PC_Relative:
		arg16 = uint16(cpu.nRead(cpu.RK, cpu.PC + 1))
		if arg16 < 0x80 {
			addr = cpu.PC + 2 + arg16
		} else {
			addr = cpu.PC + 2 + arg16 - 0x100
		}


	// rel16          - p. 309 or 5.18 (BRL)
	case m_PC_Relative_Long:
		arg16 = cpu.nRead16_wrap(cpu.RK, cpu.PC + 1)
		if arg16 < 0x8000 {
			addr = cpu.PC + 3 + arg16
		} else {
			addr = cpu.PC + 1 + arg16 - 0xffff // zamiast 0x10000 bo overflow na 16 bit
		}


	// $32, S         - p. 324 or 5.20
	case m_Stack_Relative:
		arg8   = cpu.nRead(cpu.RK, cpu.PC + 1)
		addr   = uint16(arg8) + cpu.SP


	// ($32, S), Y    - p. 325 or 5.21 (STACK,S),Y
	case m_Stack_Relative_Indirect_Y:
		arg8   = cpu.nRead(cpu.RK, cpu.PC + 1)
		arg16  = cpu.nRead16_wrap(0x00, uint16(arg8) + cpu.SP)
		//fmt.Fprintf(&cpu.LogBuf, "m_Stack_Relative_Indirect_Y: arg16 $%04x ", arg16)
		if cpu.X == 1 {
			ea     = (uint32(cpu.RDBR) << 16 | uint32(arg16)) + uint32(cpu.RYl)
		} else {
			ea     = (uint32(cpu.RDBR) << 16 | uint32(arg16)) + uint32(cpu.RY)
		}
		//fmt.Fprintf(&cpu.LogBuf, "ea $%06x ", ea)

	case m_BlockMove:
		addr = cpu.PC + 1


	default:
		cpu.stepPC = 0
		mylog.Logger.Log(fmt.Sprintf("unknown addressing mode PC $%02x:%04x", cpu.RK, cpu.PC))
	}


	// cycles adjust calculation
	// M,X and DL here      - here
	// X and Page crossing  - here
	// E, EP for commands   - op_*
	// branching            - op_*
	if cpu.M == 1 {
		cpu.Cycles -= decCycles_flagM[opcode]
	}

	if cpu.X == 1 {
		cpu.Cycles -= decCycles_flagX[opcode]
		if pageCrossed {
			cpu.Cycles += incCycles_PageCross[opcode]
		}
	}

	if cpu.RD&0x00ff != 0x00 {
		cpu.Cycles += incCycles_regDL_not00[opcode]
	}


	// instruction execution
	info := &stepInfo{ea, addr, cpu.PC, mode}
	instructions[opcode].proc(info)

	// counter and PC update
	cpu.AllCycles += uint64(cpu.Cycles)
	cpu.PC += cpu.stepPC
    if cpu.abort {
        cpu.abort = false
        return int(cpu.Cycles), true
    } else {
	    return int(cpu.Cycles), false
    }
}

// NMI - Non-Maskable Interrupt
// XXXX - valid for 6502
func (cpu *CPU) nmi() {
	cpu.push16(cpu.PC)
	cpu.op_php(nil)
	//cpu.PC = cpu.Read16(0xFFFA)
	cpu.PC = cpu.nRead16_cross(0x00, 0xFFEA)
	cpu.I = 1
	cpu.Cycles += 7
}

// IRQ - IRQ Interrupt
// XXXX - valid for 6502
func (cpu *CPU) irq() {

	cpu.push(cpu.RK)
	cpu.push16(cpu.PC)		// "next instruction to be executed" means current pc because irq is processed at start of new comm.
	cpu.push(cpu.Flags())

	cpu.I  = 1
	cpu.D  = 0
	cpu.RK = 0
	cpu.PC   = cpu.nRead16_cross(0x00, 0xFFEE)
	//mylog.Logger.Log(fmt.Sprintf("\ncpu: irq triggered, PC %4X", cpu.PC))
}


/* ====================================================================
 *
 *
 *
 * ====================================================================
 */

// ADC - Add with Carry
// I'm not sure what I'm doing ;)
func (cpu *CPU) op_adc(info *stepInfo) {
	if cpu.M == 1 {
		a := uint16(cpu.RAl)
		d := uint16(cpu.cmdRead(info))
		c := uint16(cpu.C)
		sum := a + d + c

		if cpu.D == 1 {
			if (sum & 0x0F) > 0x09 {
				sum = sum + 0x06
			}
			if (sum & 0xF0) > 0x90 {
				sum = sum + 0x60
			}
		}

		if sum > 0xFF {
			cpu.C = 1
		} else {
			cpu.C = 0
		}

		// overflow = ~(a ^ arg) & (a ^ sum) & 0x80;
		if (a^d)&0x80 == 0 && (a^sum)&0x80 != 0 {
			cpu.V = 1
		} else {
			cpu.V = 0
		}

		cpu.RAl = byte(sum)
		cpu.setZN8(cpu.RAl)
	} else {
		a := uint32(cpu.RA)
		d := uint32(cpu.cmdRead16(info))
		c := uint32(cpu.C)
		sum := a + d + c

		if cpu.D == 1 {
			if (sum & 0x000F) > 0x0009 {
				sum = sum + 0x0006
			}
			if (sum & 0x00F0) > 0x0090 {
				sum = sum + 0x0060
			}
			if (sum & 0x0F00) > 0x0900 {
				sum = sum + 0x0600
			}
			if (sum & 0xF000) > 0x9000 {
				sum = sum + 0x6000
			}
		}

		if sum > 0xFFFF {
			cpu.C = 1
		} else {
			cpu.C = 0
		}

		if (a^d)&0x8000 == 0 && (a^sum)&0x8000 != 0 {
			cpu.V = 1
		} else {
			cpu.V = 0
		}

		cpu.RA = uint16(sum)
		cpu.setZN16(cpu.RA)
	}
}



// AND - Logical AND
func (cpu *CPU) op_and(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RAl = cpu.RAl & cpu.cmdRead(info)
		cpu.setZN8(cpu.RAl)
	} else {
		cpu.RA  = cpu.RA  & cpu.cmdRead16(info)
		cpu.setZN16(cpu.RA)
	}
}

// ASL - Arithmetic Shift Left
func (cpu *CPU) op_asl(info *stepInfo) {
	if info.mode == m_Accumulator {
		if cpu.M == 1 {
			cpu.C   = (cpu.RAl >> 7) & 1
			cpu.RAl = (cpu.RAl << 1)        // or cpu.RAl <<= 1
			cpu.setZN8(cpu.RAl)
		} else {
			cpu.C   = byte(cpu.RA >> 15) & 1
			cpu.RA  = (cpu.RA << 1)
			cpu.setZN16(cpu.RA)
		}
	} else {
		if cpu.M == 1 {
			value  := cpu.cmdRead(info)
			cpu.C   = (value >> 7) & 1
			value   = (value << 1)
			cpu.cmdWrite(info, value)
			cpu.setZN8(value)
		} else {
			value  := cpu.cmdRead16(info)
			cpu.C   = byte(value >> 15) & 1
			value   = (value << 1)
			cpu.cmdWrite16(info, value)
			cpu.setZN16(value)
		}
	}
}


// BCC - Branch if Carry Clear
func (cpu *CPU) op_bcc(info *stepInfo) {
	if cpu.C == 0 {
		cpu.addBranchCycles(info) // always before PC change!
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BCS - Branch if Carry Set
func (cpu *CPU) op_bcs(info *stepInfo) {
	if cpu.C != 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BEQ - Branch if Equal
func (cpu *CPU) op_beq(info *stepInfo) {
	if cpu.Z != 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BIT - test BITs
// TODO: there is a mistake here, immediate mode should not set N nor V!


/*
OP LEN CYCLES      MODE      nvmxdizc e SYNTAX
-- --- ----------- --------- ---------- ------
24 2   4-m+w       dir       mm....m. . BIT $10
2C 3   5-m         abs       mm....m. . BIT $9876
34 2   5-m+w       dir,X     mm....m. . BIT $10,X
3C 3   6-m-x+x*p   abs,X     mm....m. . BIT $9876,X
89 3-m 3-m         imm       ......m. . BIT #$54

Immediate addressing only affects the z flag (with the result of the bitwise
And), but does not affect the n and v flags. All other addressing modes of BIT
affect the n, v, and z flags. This is the only instruction in the 6502 family
where the flags affected depends on the addressing mode.

-  The n flag reflects the high bit of the data (note: just the data, not the
bitwise And of the accumulator and the data).

-  The v flag reflects the second highest bit of the data (i.e. bit 14 of the
data when the m flag is 0, and bit 6 of the data when the m flag is 1, and
again, just the data, not the bitwise And).

-  The z flag reflects whether the result (of the bitwise And) is zero. 
*/

func (cpu *CPU) op_bit(info *stepInfo) {
	if cpu.M == 1 {
		val := cpu.cmdRead(info)
		cpu.setZ8(cpu.RAl & val)

		if info.mode != m_Immediate {
			cpu.setN8(val)
			if val&0x40 != 0 {
				cpu.V = 1
			} else {
				cpu.V = 0
			}
		}

	} else {
		val := cpu.cmdRead16(info)
		cpu.setZ16(cpu.RA & val)

		if info.mode != m_Immediate {
			cpu.setN16(val)
			if val&0x4000 != 0 {
				cpu.V = 1
			} else {
				cpu.V = 0
			}
		}
	}
}

// BMI - Branch if Minus
func (cpu *CPU) op_bmi(info *stepInfo) {
	if cpu.N != 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BNE - Branch if Not Equal
func (cpu *CPU) op_bne(info *stepInfo) {
	if cpu.Z == 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BPL - Branch if Positive
func (cpu *CPU) op_bpl(info *stepInfo) {
	if cpu.N == 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BRK - Force Interrupt
// XXX - from now duplicate with irq?
func (cpu *CPU) op_brk(info *stepInfo) {
	if cpu.E == 1 {
		cpu.push16(cpu.PC+2)
		cpu.push(cpu.Flags() | 0x10)

		cpu.I  = 1
		cpu.D  = 0
		cpu.RK = 0
		cpu.PC   = cpu.nRead16_cross(0x00, 0xFFFE)
		cpu.Cycles -= 1 // 7 cycles when E=1
	} else {
		cpu.push(cpu.RK)
		cpu.push16(cpu.PC+2)
		cpu.push(cpu.Flags())

		cpu.I  = 1
		cpu.D  = 0
		cpu.RK = 0
		cpu.PC   = cpu.nRead16_cross(0x00, 0xFFE6)
	}
	cpu.stepPC = 0
}

// BVC - Branch if Overflow Clear
func (cpu *CPU) op_bvc(info *stepInfo) {
	if cpu.V == 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// BVS - Branch if Overflow Set
func (cpu *CPU) op_bvs(info *stepInfo) {
	if cpu.V != 0 {
		cpu.addBranchCycles(info)
		cpu.PC = info.addr
		cpu.stepPC = 0
	}
}

// CLC - Clear Carry Flag
func (cpu *CPU) op_clc(info *stepInfo) {
	cpu.C = 0
}

// CLD - Clear Decimal Mode
func (cpu *CPU) op_cld(info *stepInfo) {
	cpu.D = 0
}

// CLI - Clear Interrupt Disable
func (cpu *CPU) op_cli(info *stepInfo) {
	cpu.I = 0
}

// CLV - Clear Overflow Flag
func (cpu *CPU) op_clv(info *stepInfo) {
	cpu.V = 0
}

// CMP - Compare
func (cpu *CPU) op_cmp(info *stepInfo) {
	if cpu.M == 1 {
		value := cpu.cmdRead(info)
		cpu.compare8(cpu.RAl, value)
	} else {
		value := cpu.cmdRead16(info)
		cpu.compare16(cpu.RA, value)
	}
}

// CPX - Compare X Register
func (cpu *CPU) op_cpx(info *stepInfo) {
	if cpu.X == 1 {
		value := cpu.cmdRead(info)
		cpu.compare8(cpu.RXl, value)
	} else {
		value := cpu.cmdRead16(info)
		cpu.compare16(cpu.RX, value)
	}
}

// CPY - Compare Y Register
func (cpu *CPU) op_cpy(info *stepInfo) {
	if cpu.X == 1 {
		value := cpu.cmdRead(info)
		cpu.compare8(cpu.RYl, value)
	} else {
		value := cpu.cmdRead16(info)
		cpu.compare16(cpu.RY, value)
	}
}

// DEC - Decrement Memory
func (cpu *CPU) op_dec(info *stepInfo) {
	if info.mode == m_Accumulator {
		if cpu.M == 1 {
			cpu.RAl--
			cpu.setZN8(cpu.RAl)
		} else {
			cpu.RA--
			cpu.setZN16(cpu.RA)
		}
	} else {
		if cpu.M == 1 {
			value := cpu.cmdRead(info) - 1
			cpu.cmdWrite(info, value)
			cpu.setZN8(value)
		} else {
			value := cpu.cmdRead16(info) - 1
			cpu.cmdWrite16(info, value)
			cpu.setZN16(value)
		}
	}
}

// DEX - Decrement X Register
func (cpu *CPU) op_dex(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RXl--
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX--
		cpu.setZN16(cpu.RX)
	}
}

// DEY - Decrement Y Register
func (cpu *CPU) op_dey(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RYl--
		cpu.setZN8(cpu.RYl)
	} else {
		cpu.RY--
		cpu.setZN16(cpu.RY)
	}
}

// EOR - Exclusive OR
func (cpu *CPU) op_eor(info *stepInfo) {
        if cpu.M == 1 {
                cpu.RAl = cpu.RAl ^ cpu.cmdRead(info)
                cpu.setZN8(cpu.RAl)
        } else {
                cpu.RA  = cpu.RA  ^ cpu.cmdRead16(info)
                cpu.setZN16(cpu.RA)
        }
}

// INC - Increment Memory
func (cpu *CPU) op_inc(info *stepInfo) {
	if info.mode == m_Accumulator {
		if cpu.M == 1 {
			cpu.RAl++
			cpu.setZN8(cpu.RAl)
		} else {
			cpu.RA++
			cpu.setZN16(cpu.RA)
		}
	} else {
		if cpu.M == 1 {
			value := cpu.cmdRead(info) + 1
			cpu.cmdWrite(info, value)
			cpu.setZN8(value)
		} else {
			value := cpu.cmdRead16(info) + 1
			cpu.cmdWrite16(info, value)
			cpu.setZN16(value)
		}
	}
}

// INX - Increment X Register
func (cpu *CPU) op_inx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RXl++
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX++
		cpu.setZN16(cpu.RX)
	}
}

// INY - Increment Y Register
func (cpu *CPU) op_iny(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RYl++
		cpu.setZN8(cpu.RYl)
	} else {
		cpu.RY++
		cpu.setZN16(cpu.RY)
	}
}

// JMP - Jump
// XXX - improve that!
func (cpu *CPU) op_jmp(info *stepInfo) {
	switch info.mode {
	case m_Absolute:
		cpu.PC = info.addr
	case m_Absolute_Indirect:
                cpu.PC = cpu.nRead16_wrap(0x00, info.addr) 
	case m_Absolute_Long:
		cpu.PC = uint16(info.ea)
		cpu.RK = byte(info.ea >> 16)
	case m_Absolute_Indirect_Long:
		cpu.PC = cpu.nRead16_wrap(0x00, info.addr)
		cpu.RK = cpu.nRead(0x00,  info.addr+2)
	default:
		cpu.PC = cpu.cmdRead16(info)
	}
	cpu.stepPC = 0
}

// JSL - Jump to Subroutine Long
func (cpu *CPU) op_jsl(info *stepInfo) {
	cpu.push(cpu.RK)
	cpu.push16(cpu.PC + 3)
	cpu.PC = uint16(info.ea)
	cpu.stepPC = 0
	cpu.RK = byte(info.ea >> 16)
}

// JSR - Jump to Subroutine
func (cpu *CPU) op_jsr(info *stepInfo) {
	cpu.push16(cpu.PC + 2)
	switch info.mode {
	case m_Absolute:
		cpu.PC = info.addr
        default:
                cpu.PC = cpu.cmdRead16(info)
        }
	cpu.stepPC = 0
}

// LDA - Load Accumulator - for Immediate arguments
func (cpu *CPU) op_lda(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RAl = cpu.cmdRead(info)
		cpu.setZN8(cpu.RAl)
	} else {
		cpu.RA = cpu.cmdRead16(info)
		cpu.setZN16(cpu.RA)
	}
}

// LDX - Load X Register
func (cpu *CPU) op_ldx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RXl = cpu.cmdRead(info)
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX = cpu.cmdRead16(info)
		cpu.setZN16(cpu.RX)
	}
}

// LDY - Load Y Register
func (cpu *CPU) op_ldy(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RYl = cpu.cmdRead(info)
		cpu.setZN8(cpu.RYl)
	} else {
		cpu.RY = cpu.cmdRead16(info)
		cpu.setZN16(cpu.RY)
	}
}

// LSR - Logical Shift Right
func (cpu *CPU) op_lsr(info *stepInfo) {
	if info.mode == m_Accumulator {
		if cpu.M == 1 {
			cpu.C     = cpu.RAl & 1
			cpu.RAl >>= 1
			cpu.setZN8(cpu.RAl)
		} else {
			cpu.C     = byte(cpu.RA & 1)
			cpu.RA  >>= 1
			cpu.setZN16(cpu.RA)
		}
	} else {
		if cpu.M == 1 {
			value  := cpu.cmdRead(info)
			cpu.C   = value & 1
			value >>= 1
			cpu.cmdWrite(info, value)
			cpu.setZN8(value)
		} else {
			value  := cpu.cmdRead16(info)
			cpu.C   = byte(value & 1)
			value >>= 1
			cpu.cmdWrite16(info, value)
			cpu.setZN16(value)
		}
	}
}

// NOP - No Operation
func (cpu *CPU) op_nop(info *stepInfo) {
}

// ORA - Logical Inclusive OR
func (cpu *CPU) op_ora(info *stepInfo) {
        if cpu.M == 1 {
                cpu.RAl = cpu.RAl | cpu.cmdRead(info)
                cpu.setZN8(cpu.RAl)
        } else {
                cpu.RA  = cpu.RA  | cpu.cmdRead16(info)
                cpu.setZN16(cpu.RA)
        }
}

// PHA - Push Accumulator
func (cpu *CPU) op_pha(info *stepInfo) {
	if cpu.M == 1 {
		cpu.push(cpu.RAl)
	} else {
		cpu.push16(cpu.RA)
	}
}

// PHP - Push Processor Status
func (cpu *CPU) op_php(info *stepInfo) {
	//cpu.push(cpu.Flags() | 0x10)
	cpu.push(cpu.Flags())
}

// PHX - PusH X Register
func (cpu *CPU) op_phx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.push(cpu.RXl)
	} else {
		cpu.push16(cpu.RX)
	}
}

// PHY - PusH Y Register
func (cpu *CPU) op_phy(info *stepInfo) {
	if cpu.X == 1 {
		cpu.push(cpu.RYl)
	} else {
		cpu.push16(cpu.RY)
	}
}

// PLA - Pull Accumulator
func (cpu *CPU) op_pla(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RAl = cpu.pull()
		cpu.setZN8(cpu.RAl)
	} else {
		cpu.RA = cpu.pull16()
		cpu.setZN16(cpu.RA)
	}
}

// PLP - Pull Processor Status
func (cpu *CPU) op_plp(info *stepInfo) {
	//cpu.SetFlags(cpu.pull()&0xEF | 0x20)
	cpu.SetFlags(cpu.pull())
}

// ROL - Rotate Left
func (cpu *CPU) op_rol(info *stepInfo) {
	c := cpu.C
	if info.mode == m_Accumulator {
		if cpu.M == 1 {
			cpu.C   = (cpu.RAl >> 7) & 1
			cpu.RAl = (cpu.RAl << 1) | c
			cpu.setZN8(cpu.RAl)
		} else {
			cpu.C   = byte(cpu.RA >> 15) & 1
			cpu.RA  = (cpu.RA << 1) | uint16(c)
			cpu.setZN16(cpu.RA)
		}
	} else {
		if cpu.M == 1 {
			value  := cpu.cmdRead(info)
			cpu.C   = (value >> 7) & 1
			value   = (value << 1) | c
			cpu.cmdWrite(info, value)
			cpu.setZN8(value)
		} else {
			value  := cpu.cmdRead16(info)
			cpu.C   = byte(value >> 15) & 1
			value   = (value << 1) | uint16(c)
			cpu.cmdWrite16(info, value)
			cpu.setZN16(value)
		}
	}
}

// ROR - Rotate Right
func (cpu *CPU) op_ror(info *stepInfo) {
	c := cpu.C
	if info.mode == m_Accumulator {
		if cpu.M == 1 {
			cpu.C   =  cpu.RAl & 1
			cpu.RAl = (cpu.RAl >> 1) | (c << 7)
			cpu.setZN8(cpu.RAl)
		} else {
			cpu.C    = byte(cpu.RA  & 1)
			cpu.RA   = (cpu.RA  >> 1) | (uint16(c) << 15)
			cpu.setZN16(cpu.RA)
		}
	} else {
		if cpu.M == 1 {
			value  := cpu.cmdRead(info)
			cpu.C   = value & 1
			value   = (value >> 1) | (c << 7)
			cpu.cmdWrite(info, value)
			cpu.setZN8(value)
		} else {
			value  := cpu.cmdRead16(info)
			cpu.C   = byte(value & 1)
			value   = (value >> 1) | (uint16(c) << 15)
			cpu.cmdWrite16(info, value)
			cpu.setZN16(value)
		}
	}
}


// RTI - Return from Interrupt
func (cpu *CPU) op_rti(info *stepInfo) {
	//mylog.Logger.Log("cpu: rti")
	//cpu.SetFlags(cpu.pull()&0xEF | 0x20)
	if cpu.E == 1 {
		cpu.SetFlags(cpu.pull())
		cpu.PC = cpu.pull16()
		cpu.Cycles -= 1 // 6 when E=1
	} else {
		cpu.SetFlags(cpu.pull())
		cpu.PC = cpu.pull16()
		cpu.RK = cpu.pull()
	}
	cpu.I = 0
	cpu.stepPC = 0
}

// RLK - ReTurn from subroutine Long
func (cpu *CPU) op_rtl(info *stepInfo) {
	cpu.PC = cpu.pull16() + 1
	cpu.RK = cpu.pull()
	cpu.stepPC = 0
}

// RTS - Return from Subroutine
func (cpu *CPU) op_rts(info *stepInfo) {
	cpu.PC = cpu.pull16() + 1
	cpu.stepPC = 0
}

// SBC - SuBstract with Carry
// I'm not sure what I'm doing ;)
func (cpu *CPU) op_sbc(info *stepInfo) {
	if cpu.M == 1 {
		a := uint16(cpu.RAl)
		d := uint16(^cpu.cmdRead(info))
		c := uint16(cpu.C)
		sum := a + d + c

		if cpu.D == 1 {
			if (sum & 0x0F) > 0x09 {
				sum = sum + 0x06
			}
			if (sum & 0xF0) > 0x90 {
				sum = sum + 0x60
			}
		}

		if sum > 0xFF {
			cpu.C = 1
		} else {
			cpu.C = 0
		}

		// overflow = ~(a ^ arg) & (a ^ sum) & 0x80;
		if (a^d)&0x80 == 0 && (a^sum)&0x80 != 0 {
			cpu.V = 1
		} else {
			cpu.V = 0
		}

		cpu.RAl = byte(sum)
		cpu.setZN8(cpu.RAl)
	} else {
		a := uint32(cpu.RA)
		d := uint32(^cpu.cmdRead16(info))
		c := uint32(cpu.C)
		sum := a + d + c

		if cpu.D == 1 {
			if (sum & 0x000F) > 0x0009 {
				sum = sum + 0x0006
			}
			if (sum & 0x00F0) > 0x0090 {
				sum = sum + 0x0060
			}
			if (sum & 0x0F00) > 0x0900 {
				sum = sum + 0x0600
			}
			if (sum & 0xF000) > 0x9000 {
				sum = sum + 0x6000
			}
		}

		if sum > 0xFFFF {
			cpu.C = 1
		} else {
			cpu.C = 0
		}

		if (a^d)&0x8000 == 0 && (a^sum)&0x8000 != 0 {
			cpu.V = 1
		} else {
			cpu.V = 0
		}

		cpu.RA = uint16(sum)
		cpu.setZN16(cpu.RA)
	}
}


/*
// SBC - Subtract with Carry
func (cpu *CPU) sbc(info *stepInfo) {
	a := cpu.RA
	b := cpu.Read(info.address)
	c := cpu.C
	cpu.RA = a - b - (1 - c)
	cpu.setZN(cpu.RA)
	if int(a)-int(b)-int(1-c) >= 0 {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^b)&0x80 != 0 && (a^cpu.RA)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}
*/

// SEC - Set Carry Flag
func (cpu *CPU) op_sec(info *stepInfo) {
	cpu.C = 1
}

// SED - Set Decimal Flag
func (cpu *CPU) op_sed(info *stepInfo) {
	cpu.D = 1
}

// SEI - Set Interrupt Disable
func (cpu *CPU) op_sei(info *stepInfo) {
	cpu.I = 1
}

// STA - Store Accumulator
func (cpu *CPU) op_sta(info *stepInfo) {
	if cpu.M == 1 {
		cpu.cmdWrite(info, cpu.RAl)
	} else {
		cpu.cmdWrite16(info, cpu.RA)
	}
}

// STX - Store X Register
func (cpu *CPU) op_stx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.cmdWrite(info, cpu.RXl)
	} else {
		cpu.cmdWrite16(info, cpu.RX)
	}
}

// STY - Store Y Register
func (cpu *CPU) op_sty(info *stepInfo) {
	if cpu.X == 1 {
		cpu.cmdWrite(info, cpu.RYl)
	} else {
		cpu.cmdWrite16(info, cpu.RY)
	}
}

// TAX - Transfer Accumulator to X
func (cpu *CPU) op_tax(info *stepInfo) {
	var src uint16
	if cpu.M == 1 {
		src = uint16(cpu.RAl)
	} else {
		src = cpu.RA
	}

	if cpu.X == 1 {
		cpu.RXl = byte(src)
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX = uint16(src)
		cpu.setZN16(cpu.RX)
	}
}


// TAY - Transfer Accumulator to Y
func (cpu *CPU) op_tay(info *stepInfo) {
	var src uint16
	if cpu.M == 1 {
		src = uint16(cpu.RAl)
	} else {
		src = cpu.RA
	}

	if cpu.X == 1 {
		cpu.RYl = byte(src)
		cpu.setZN8(cpu.RYl)
	} else {
		cpu.RY = uint16(src)
		cpu.setZN16(cpu.RY)
	}
}

// TSX - Transfer Stack Pointer to X
func (cpu *CPU) op_tsx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RXl = byte(cpu.SP & 0x00ff)
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX  = cpu.SP
		cpu.setZN16(cpu.RX)
	}
}

// TXA - Transfer X to Accumulator
func (cpu *CPU) op_txa(info *stepInfo) {
	var src uint16
	if cpu.X == 1 {
		src = uint16(cpu.RXl)
	} else {
		src = cpu.RX
	}

	if cpu.M == 1 {
		cpu.RAl = byte(src)
		cpu.setZN8(cpu.RAl)
	} else {
		cpu.RA = uint16(src)
		cpu.setZN16(cpu.RA)
	}
}

// TXS - Transfer X to Stack Pointer
func (cpu *CPU) op_txs(info *stepInfo) {
	var src uint16
	if cpu.X == 1 {
		src = uint16(cpu.RXl)
	} else {
		src = cpu.RX
	}

	if cpu.E == 1 {
		cpu.SP = 0x0100 | (src & 0x00ff)
	} else {
		cpu.SP = src
	}
}

// TYA - Transfer Y to Accumulator
func (cpu *CPU) op_tya(info *stepInfo) {
	var src uint16
	if cpu.X == 1 {
		src = uint16(cpu.RYl)
	} else {
		src = cpu.RY
	}

	if cpu.M == 1 {
		cpu.RAl = byte(src)
		cpu.setZN8(cpu.RAl)
	} else {
		cpu.RA = uint16(src)
		cpu.setZN16(cpu.RA)
	}
}

// BRA - BRanch Always
func (cpu *CPU) op_bra(info *stepInfo) {
	cpu.addBranchCycles(info) // always before PC change!
	cpu.PC = info.addr
	cpu.stepPC = 0
}

// BRL - BRanch Long
func (cpu *CPU) op_brl(info *stepInfo) {
	cpu.PC = info.addr
	cpu.stepPC = 0
}

// COP - COProcessor
func (cpu *CPU) op_cop(info *stepInfo) {
	if cpu.E == 1 {
		cpu.push16(cpu.PC+2)
		cpu.push(cpu.Flags())

		cpu.I  = 1
		cpu.D  = 0
		cpu.RK = 0
		cpu.PC   = cpu.nRead16_cross(0x00, 0xFFF4)
		cpu.Cycles -= 1 // COP for E=1 has 7 cycles
	} else {
		cpu.push(cpu.RK)
		cpu.push16(cpu.PC+2)
		cpu.push(cpu.Flags())

		cpu.I  = 1
		cpu.D  = 0
		cpu.RK = 0
		cpu.PC   = cpu.nRead16_cross(0x00, 0xFFE4)
	}
	cpu.stepPC = 0
}

// MVN - MoVe memory Negative
func (cpu *CPU) op_mvn(info *stepInfo) {
	dst := cpu.nRead(cpu.RK, info.addr)
	src := cpu.nRead(cpu.RK, info.addr+1)

	cpu.RDBR = dst
	if cpu.X == 1 {
		cpu.nWrite(dst, uint16(cpu.RYl), cpu.nRead(src, uint16(cpu.RXl)))
		cpu.RYl++
		cpu.RXl++
	} else {
		cpu.nWrite(dst, cpu.RY, cpu.nRead(src, cpu.RX))
		cpu.RY++
		cpu.RX++
	}


	if cpu.M == 1 {
		cpu.RAl--
		if cpu.RAl != 0xff {
			cpu.stepPC = 0
		}
	} else {
		cpu.RA--
		if cpu.RA != 0xffff {
			cpu.stepPC = 0
		}
	}
}	

// MVP - MoVe memory Positive
func (cpu *CPU) op_mvp(info *stepInfo) {
	dst := cpu.nRead(cpu.RK, info.addr)
	src := cpu.nRead(cpu.RK, info.addr+1)

	cpu.RDBR = dst
	if cpu.X == 1 {
		cpu.nWrite(dst, uint16(cpu.RYl), cpu.nRead(src, uint16(cpu.RXl)))
		cpu.RYl--
		cpu.RXl--
	} else {
		cpu.nWrite(dst, cpu.RY, cpu.nRead(src, cpu.RX))
		cpu.RY--
		cpu.RX--
	}


	if cpu.M == 1 {
		cpu.RAl--
		if cpu.RAl != 0xff {
			cpu.stepPC = 0
		}
	} else {
		cpu.RA--
		if cpu.RA != 0xffff {
			cpu.stepPC = 0
		}
	}
}

// PHB - PusH data Bank register
func (cpu *CPU) op_phb(info *stepInfo) {
	cpu.push(cpu.RDBR)
}

// PHD - PusH Direct register
func (cpu *CPU) op_phd(info *stepInfo) {
	cpu.push16(cpu.RD)
}

// PHK - PusH K register
func (cpu *CPU) op_phk(info *stepInfo) {
	cpu.push(cpu.RK)
}


// PLX - PulL X Register
func (cpu *CPU) op_plx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RXl = cpu.pull()
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX = cpu.pull16()
		cpu.setZN16(cpu.RX)
	}
}

// PLY - PulL Y Register
func (cpu *CPU) op_ply(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RYl = cpu.pull()
		cpu.setZN8(cpu.RYl)
	} else {
		cpu.RY = cpu.pull16()
		cpu.setZN16(cpu.RY)
	}
}

// PEA - Push Effective Address
func (cpu *CPU) op_pea(info *stepInfo) {
	cpu.push16(cpu.cmdRead16(info))
	//val := cpu.cmdRead16(info)
	//fmt.Fprintf(&cpu.LogBuf, "PEA %04x %04x\n", val, info.addr)
}

// PLD - PulL Direct register
func (cpu *CPU) op_pld(info *stepInfo) {
	cpu.RD = cpu.pull16()
	cpu.setZN16(cpu.RD)
}

// PER - Push Effective Relative address
func (cpu *CPU) op_per(info *stepInfo) {
	cpu.push16(info.addr)
}

// PEI - Push Effective Indirect address
func (cpu *CPU) op_pei(info *stepInfo) {
	cpu.push16(cpu.cmdRead16(info))
}

// PLB - PulL data Bank register
func (cpu *CPU) op_plb(info *stepInfo) {
	cpu.RDBR = cpu.pull()
	cpu.setZN8(cpu.RDBR)
}

// REset Processor status bits
func (cpu *CPU) op_rep(info *stepInfo) {
	neg_flags := ^cpu.Bus.EaRead(info.ea)
	tmp_flags := cpu.Flags() & neg_flags
	//fmt.Fprintf(&cpu.LogBuf, "op_rep %08b %08b %08b %08b\n", cpu.Bus.EaRead(info.ea), neg_flags, cpu.Flags(), tmp_flags)
	cpu.SetFlags(tmp_flags)
}

// SEt Processor status bits
func (cpu *CPU) op_sep(info *stepInfo) {
	tmp_flags := cpu.Flags() | cpu.Bus.EaRead(info.ea)
	cpu.SetFlags(tmp_flags)
}

// STP - SToP the clock
func (cpu *CPU) stp(info *stepInfo) {
}

// STZ - STore Zero
func (cpu *CPU) op_stz(info *stepInfo) {
	if cpu.M == 1 {
		cpu.cmdWrite(info, 0x00)
	} else {
		cpu.cmdWrite16(info, 0x0000)
	}
}

// TCD - Transfer C accumulator to Direct register
func (cpu *CPU) op_tcd(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RD = (uint16(cpu.RAh) << 8) | uint16(cpu.RAl)
	} else {
		cpu.RD = cpu.RA
	}
	cpu.setZN16(cpu.RD) // XXX: always 16 bit regardless od M register?
}

// TCS - Transfer C accumulator to Stack pointer
func (cpu *CPU) op_tcs(info *stepInfo) {
	var src uint16
	if cpu.M == 1 {
		src = (uint16(cpu.RAh) << 8) | uint16(cpu.RAl)
	} else {
		src = cpu.RA
	}

	if cpu.E == 1 {
		cpu.SP = 0x0100 | (src & 0x00ff)
	} else {
		cpu.SP = src
	}
}

// TCD - Transfer Direct register to C accumulator
func (cpu *CPU) op_tdc(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RAh = byte(cpu.RD >> 8)
		cpu.RAl = byte(cpu.RD)
	} else {
		cpu.RA = cpu.RD
	}
	cpu.setZN16(cpu.RD) // XXX: always 16 bit regardless od M register?
}

// Test and Reset Bits
func (cpu *CPU) op_trb(info *stepInfo) {
	if cpu.M == 1 {
		value := cpu.cmdRead(info)
		cpu.cmdWrite(info, value & (^cpu.RAl))
		cpu.setZN8(value & cpu.RAl)
	} else {
		value := cpu.cmdRead16(info)
		cpu.cmdWrite16(info, value & (^cpu.RA))
		cpu.setZN16(value & cpu.RA)
	}
}

// Test and Set Bits
func (cpu *CPU) op_tsb(info *stepInfo) {
	if cpu.M == 1 {
		value := cpu.cmdRead(info)
		cpu.cmdWrite(info, value | cpu.RAl)
		cpu.setZN8(value & cpu.RAl)
	} else {
		value := cpu.cmdRead16(info)
		cpu.cmdWrite16(info, value | cpu.RA)
		cpu.setZN16(value & cpu.RA)
	}
}

// TSC - Transfer Stack pointer to C accumulator
func (cpu *CPU) op_tsc(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RAh = byte(cpu.SP >> 8)
		cpu.RAl = byte(cpu.SP)
	} else {
		cpu.RA = cpu.SP
	}
	cpu.setZN16(cpu.SP) // XXX: always 16 bit regardless od M register?
}

// TXY - Transfer X register to Y register
func (cpu *CPU) op_txy(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RYl = cpu.RXl
		cpu.setZN8(cpu.RYl)
	} else {
		cpu.RY = cpu.RX
		cpu.setZN16(cpu.RY)
	}
}

// TYX - Transfer Y register to X register
func (cpu *CPU) op_tyx(info *stepInfo) {
	if cpu.X == 1 {
		cpu.RXl = cpu.RYl
		cpu.setZN8(cpu.RXl)
	} else {
		cpu.RX = cpu.RY
		cpu.setZN16(cpu.RX)
	}
}

// WAI - WAit for Interrupt
func (cpu *CPU) wai(info *stepInfo) {
}

// WDM - William D. Mensch, Jr.
func (cpu *CPU) op_wdm(info *stepInfo) {
    cpu.abort = true
}

// XBA - eXchange B and A accumulator
func (cpu *CPU) op_xba(info *stepInfo) {
	if cpu.M == 1 {
		cpu.RAh, cpu.RAl = cpu.RAl, cpu.RAh
		cpu.setZN8(cpu.RAl)
	} else {
		newl := cpu.RA >> 8
		newh := cpu.RA << 8
		cpu.RA = newh | newl
		cpu.setZN8(byte(newl))
	}
}

// XCE - eXchange Carry and Emulation flags
func (cpu *CPU) op_xce(info *stepInfo) {
	if cpu.C == cpu.E {
		return
	}
	if cpu.C == 1 {
		cpu.C = cpu.E
		cpu.SetFlags(cpu.Flags() | 0x30)
		cpu.E = 1 // after SetFlags() due to conflict 
		cpu.SP = 0x0100  | (cpu.SP & 0x00ff)
	} else {
		cpu.C = cpu.E
		cpu.E = 0
	}
}

// unknown opcode - XXX - todo PANIC?
func (cpu *CPU) op_xxx(info *stepInfo) {
	log.Fatalf("unknown instruction %02x at $%04x", cpu.nRead(cpu.RK, cpu.PC), cpu.PC)
}
