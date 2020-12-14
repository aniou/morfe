package gabe

import (
	"fmt"

	"github.com/aniou/go65c816/lib/mylog"
	"github.com/aniou/go65c816/lib/queue"
)

type Gabe struct {
	//out    chan byte       // for 'display'
	InBuf    queue.QueueByte // for 'keyboard'
	Data     byte
	command  byte
	mem      []byte
}

func New() (*Gabe, error) {
	//g := Gabe{make(chan byte, 200), queue.NewQueueByte(200)}
	gabe := Gabe{queue.NewQueueByte(200), 0, 0, make([]byte, 0x2B)}
	return &gabe, nil
}

func (g *Gabe) Dump(address uint32) []byte {
	return nil // XXX - todo
}

func (g *Gabe) String() string {
	return "Gabe area"
}

func (g *Gabe) Shutdown() {
}

func (g *Gabe) Clear() { // Maybe Reset?
}

func (g *Gabe) Size() uint32 {
	return 0x100 // XXX: something
}

func (g *Gabe) Read(address uint32) byte {
	//mylog.Logger.Log(fmt.Sprintf("."))
	switch {
	case address >= 0x00_0100 && address <= 0x00_012B:		// MATH Coprocessor
		return g.mem[address - 0x00_0100]

	case address >= 0x00_012E && address <= 0x00_012F:		// grey area of errors, we simulate memory here
		return g.mem[address - 0x00_0100]

	case address == 0xAF1060:
		return g.Data
		if g.InBuf.Len() > 0 {
			return *g.InBuf.Dequeue()
		} else {
			return 0
		}
	case address == 0xAF1064:			// we support only bit 0 
		return g.command
		/*
		if g.InBuf.Len() > 0 {
			mylog.Logger.Log(fmt.Sprintf("gabe: read from addr %6X 1 returned", address))
			return 2
		} else {
			return 0
		}
		*/
	default:
		mylog.Logger.Log(fmt.Sprintf("gabe: read from addr %6X is not implemented, 0 returned", address))
		return 0
	}
}

// based on MathCoproRegister.cs from FoenixIDE
func (g *Gabe) mathCoop(address uint32, val byte) {
	a := address - 0x00_0100
	switch address {
		case 0x00_0100, 0x00_0101,		// UNSIGNED_MULT_A
		     0x00_0102, 0x00_0103:		// UNSIGNED_MULT_B

			g.mem[a] = val
			op1    := uint16(g.mem[0x00]) + uint16(g.mem[0x01]) << 8
			op2    := uint16(g.mem[0x02]) + uint16(g.mem[0x03]) << 8

			result := uint32(op1 * op2)

			g.mem[0x04] = byte(result       & 0xff)
			g.mem[0x05] = byte(result >> 8  & 0xff)
			g.mem[0x06] = byte(result >> 16 & 0xff)
			g.mem[0x07] = byte(result >> 24 & 0xff)

		case 0x00_0108, 0x00_0109,		// SIGNED_MULT_A
		     0x00_010A, 0x00_010B:		// SIGNED_MULT_B

			g.mem[a] = val
			op1    := int16(g.mem[0x08]) + int16(g.mem[0x09]) << 8
			op2    := int16(g.mem[0x0a]) + int16(g.mem[0x0b]) << 8

			result := int32(op1 * op2)

			g.mem[0x0c] = byte(result       & 0xff)
			g.mem[0x0d] = byte(result >> 8  & 0xff)
			g.mem[0x0e] = byte(result >> 16 & 0xff)
			g.mem[0x0f] = byte(result >> 24 & 0xff)

		case 0x00_0110, 0x00_0111,		// UNSIGNED_DIV_DEM
	 	     0x00_0112, 0x00_0113:		// UNSIGNED_DIV_NUM

			g.mem[a] = val
			op1    := uint16(g.mem[0x10]) + uint16(g.mem[0x11]) << 8
			op2    := uint16(g.mem[0x12]) + uint16(g.mem[0x13]) << 8
			
			var result, remainder uint16
			if (op1 != 0) {
				result = op2 / op1
				remainder = op2 % op1
			}

			g.mem[0x14] = byte(result          & 0xff)
			g.mem[0x15] = byte(result    >> 8  & 0xff)
			g.mem[0x16] = byte(remainder       & 0xff)
			g.mem[0x17] = byte(remainder >> 8  & 0xff)

		case 0x00_0118, 0x00_0119,		// SIGNED_DIV_DEM
	 	     0x00_011A, 0x00_011B:		// SIGNED_DIV_NUM

			g.mem[a] = val
			op1    := int16(g.mem[0x18]) + int16(g.mem[0x19]) << 8
			op2    := int16(g.mem[0x1A]) + int16(g.mem[0x1B]) << 8
			
			var result, remainder int16
			if (op1 != 0) {
				result = op2 / op1
				remainder = op2 % op1
			}

			g.mem[0x1C] = byte(result          & 0xff)
			g.mem[0x1D] = byte(result    >> 8  & 0xff)
			g.mem[0x1E] = byte(remainder       & 0xff)
			g.mem[0x1F] = byte(remainder >> 8  & 0xff)

		case 0x00_0120, 0x00_0121, 0x00_122, 0x00_123,	// ADDER32_A
		     0x00_0124, 0x00_0125, 0x00_126, 0x00_127:	// ADDER32_B

			g.mem[a] = val
			op1    := int32(g.mem[0x20])       + 
			          int32(g.mem[0x21]) <<  8 + 
				  int32(g.mem[0x22]) << 16 + 
				  int32(g.mem[0x23]) << 24

			op2    := int32(g.mem[0x24])       + 
			          int32(g.mem[0x25]) << 8  + 
				  int32(g.mem[0x26]) << 16 + 
				  int32(g.mem[0x27]) << 24

			result := int32(op1 + op2)

			g.mem[0x28] = byte(result       & 0xff)
			g.mem[0x29] = byte(result >> 8  & 0xff)
			g.mem[0x2a] = byte(result >> 16 & 0xff) 
			g.mem[0x2b] = byte(result >> 24 & 0xff) 

		default:
			mylog.Logger.Log(fmt.Sprintf("gabe-math: write to addr %6X val %2X is not implemented", address, val))
	}
}


// taken from FoenixIDE
func (g *Gabe) Write(address uint32, val byte) {
	switch {
	case address >= 0x00_0100 && address <= 0x00_012B:	// MATH Coprocessor
		g.mathCoop(address, val)

	case address >= 0x00_012E && address <= 0x00_012F:	// grey area of errors, we simulate memory here
		g.mem[address - 0x00_0100] = val

	case address == 0xAF1060:
		if val == 0x69 {				// 
			g.command = 1		// 
		}
		if val == 0xEE {				// echo
			g.command = 1		// 
		}
		if val == 0xF4 {				// kbd reset
			g.InBuf.Enqueue(0xFA)		// self-test result
			g.Data = 0xFA
			g.command = 1
		}
		if val == 0xF6 {				// 
			g.command = 1
		}
	case address == 0xAF1064:
		if val == 0x20 {				// 
			g.command = 1
		}
		if val == 0x60 {				// 
			g.command = 1
		}
		if val == 0xA8 {				// 
			g.command = 1
		}
		if val == 0xA9 {				// 
			g.InBuf.Enqueue(0x00)		// 
			g.Data = 0
			g.command = 0x01
		}
		if val == 0xAA {				// self-test
			g.InBuf.Enqueue(0x55)		// self-test result
			g.Data = 0x55
			g.command = 0x01
		}
		if val == 0xAB {				// self-test
			g.InBuf.Enqueue(0x00)		// 
			g.Data = 0
		}
		if val == 0xD4 {				// 
			g.command = 0x01
		}
		if val == 0x00 {
			g.command = 0x00
		}
	default:
		mylog.Logger.Log(fmt.Sprintf("gabe: write to addr %6X val %2X is not implemented, 0 returned", address, val))
	}
	return
}

