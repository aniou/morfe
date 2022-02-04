
package cpu_65c816

import (
	"fmt"
)

func printCPUFlags(flag byte, name string) (string) {
	if flag > 0 {
		return name
	} else {
		return "-"
	}
}

func (c *CPU) formatInstructionMode(mode byte, w0 byte, w1 byte, w2 byte, w3 byte) string {
	var arg string

	switch mode {
        case m_Absolute:                     // $9876       - p. 288 or 5.2
		arg = fmt.Sprintf("$%02x%02x", w2, w1)
        case m_Absolute_X:                   // $9876, X    - p. 289 or 5.3
		arg = fmt.Sprintf("$%02x%02x, X", w2, w1)
        case m_Absolute_Y:                   // $9876, Y    - p. 290 or 5.3
		arg = fmt.Sprintf("$%02x%02x, Y", w2, w1)
        case m_Accumulator:                  // A           - p. 296 or 5.6
		arg = "A"
        case m_Immediate:                    // #$aa        - p. 306 or 5.14
		if w0 == 0xf4 {
			arg = fmt.Sprintf("#$%02x%02x", w2, w1)
		} else {
			arg = fmt.Sprintf("#$%02x", w1)
		}
        case m_Immediate_flagM:              // #$aaaa/$#aa - p. 306 or 5.14 
		if c.M == 1 {
			arg = fmt.Sprintf("#$%02x", w1)
		} else {
			arg = fmt.Sprintf("#$%02x%02x", w2, w1)
		}
        case m_Immediate_flagX:              // #$aa        - p. 306 or 5.14 // XXX fix it
		if c.X == 1 {
			arg = fmt.Sprintf("#$%02x", w1)
		} else {
			arg = fmt.Sprintf("#$%02x%02x", w2, w1)
		}
        case m_Implied:                      // -           - p. 307 or 5.15
		arg = ""
        case m_DP:                           // $12         - p. 298 or 5.7
		arg = fmt.Sprintf("$%02x", w1)
        case m_DP_X:                         // $12, X      - p. 299 or 5.8
		arg = fmt.Sprintf("$%02x, X", w1)
        case m_DP_Y:                         // $12, Y      - p. 300 or 5.8
		arg = fmt.Sprintf("$%02x, Y", w1)
        case m_DP_X_Indirect:                // ($12, X)    - p. 301 or 5.11
		arg = fmt.Sprintf("($%02x, X)", w1)
        case m_DP_Indirect:                  // ($12)       - p. 302 or 5.9
		arg = fmt.Sprintf("($%02x)", w1)
        case m_DP_Indirect_Long:             // [$12]       - p. 303 or 5.10
		arg = fmt.Sprintf("[$%02x]", w1)
        case m_DP_Indirect_Y:                // ($12), Y    - p. 304 or 5.12
		arg = fmt.Sprintf("($%02x), Y", w1)
        case m_DP_Indirect_Long_Y:           // [$12], Y    - p. 305 or 5.13
		arg = fmt.Sprintf("[$%02x], Y", w1)
        case m_Absolute_X_Indirect:          // ($1234, X)  - p. 291 or 5.5
		arg = fmt.Sprintf("($%02x%02x, X)", w2, w1)
        case m_Absolute_Indirect:            // ($1234)     - p. 292 or 5.4
		arg = fmt.Sprintf("($%02x%02x)", w2, w1)
        case m_Absolute_Indirect_Long:       // [$1234]     - p. 293 or 5.10
		arg = fmt.Sprintf("[$%02x%02x]", w2, w1)
        case m_Absolute_Long:                // $abcdef     - p. 294 or 5.16
		arg = fmt.Sprintf("$%02x%02x%02x", w3, w2, w1)
        case m_Absolute_Long_X:              // $abcdex, X  - p. 295 or 5.17
		arg = fmt.Sprintf("$%02x%02x%02x, X", w3, w2, w1)
        case m_BlockMove:                    // #$12,#$34   - p. 297 or 5.19 (MVN, MVP)
		arg = fmt.Sprintf("#$%02x,#$%02x", w2, w1) // XXX - verify it!
        case m_PC_Relative:                  // rel8        - p. 308 or 5.18 (BRA)
		w216 := uint16(w1)
                if w2 < 0x80 {
                        dest := c.PC + 2 + w216
			arg = fmt.Sprintf("$%02x ($%04x +)", w216, dest)
                } else {
                        dest := c.PC + 2 + w216 - 0x100
			arg = fmt.Sprintf("$%02x ($%04x -)", w216, dest)
                }
        case m_PC_Relative_Long:             // rel16       - p. 309 or 5.18 (BRL)
		arg = "TODO"
        case m_Stack_Relative:               // $32, S      - p. 324 or 5.20
		arg = fmt.Sprintf("$%02x, S",  w1)
        case m_Stack_Relative_Indirect_Y:    // ($32, S), Y - p. 325 or 5.21 (STACK,S),Y
		arg = fmt.Sprintf("($%02x, S), Y",  w1)
	default:
		arg = "! unknown !"
	}
	return arg
}

// XXX - create disassemble line
func (c *CPU) DisassemblePreviousPC() string {
	return c.Disassemble(c.PPC)
}

func (c *CPU) DisassembleCurrentPC() string {
	return c.Disassemble(c.PC)
}


func (c *CPU) Disassemble(myPC uint16) string {
	//var myPC uint16 = c.PC
	var numeric string
	var output string

	//opcode := c.Read(myPC)
	opcode := c.nRead(c.RK, myPC)
	mode := instructions[opcode].mode

	// crude and incosistent size adjust 
	var sizeAdjust byte;
	if mode == m_Immediate_flagM {
		sizeAdjust = c.M
	}
	if mode == m_Immediate_flagX {
		sizeAdjust = c.X
	}

	bytes := instructions[opcode].size - sizeAdjust
	name := instructions[opcode].name
	w0 := c.nRead(c.RK, myPC+0)
	w1 := c.nRead(c.RK, myPC+1)
	w2 := c.nRead(c.RK, myPC+2)
	w3 := c.nRead(c.RK, myPC+3)

	switch bytes {
	case 4:
		numeric = fmt.Sprintf("%02x %02x %02x %02x", w0, w1, w2, w3)
	case 3:
		numeric = fmt.Sprintf("%02x %02x %02x", w0, w1, w2)
	case 2:
		numeric = fmt.Sprintf("%02x %02x", w0, w1)
	case 1:
		numeric = fmt.Sprintf("%02x", w0)
	default:
		numeric = fmt.Sprintf("err: cmd len %d", bytes)
	}

	arg := c.formatInstructionMode(mode, w0, w1, w2, w3)

	if c.Cycles == 0 {
		output = fmt.Sprintf(output, "--:----│           │                 │")

	}
	// XXX - change to different log system
	//if c.Cycles > 9 {
	//	fmt.Fprintf(logView, "warning: instruction cycles > 10\n")
	//}
	//output = fmt.Sprintf("%d\n%02x:%04x│%-11v│%3s %-13v│",
	//				c.Cycles, c.RK, myPC, numeric, name, arg)

	//fmt.Fprintf(v, "%-38v",   "3│00:000c│02 02      │BEQ 02 ($04fa +)")
	output = fmt.Sprintf("%02x %04x: %-11v : %3s %-13v", c.RK, myPC, numeric, name, arg)

	return output
}

