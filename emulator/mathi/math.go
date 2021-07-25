
// Math Co-Procesor
// integer operations

package mathi

import (
)


type MathInt struct {
	name	string
	mem	[]byte
}

func New(name string, size int) *MathInt {
        m := MathInt{
                      name: name,
                       mem: make([]byte, size) }
        return &m
}

func (m *MathInt) Read(addr uint32) (byte, error) {
	return m.mem[addr], nil
}

func (m *MathInt) Write(addr uint32, val byte) error {
        switch addr {
        case 0x00, 0x01, 0x02, 0x03:   // UNSIGNED_MULT_A, UNSIGNED_MULT_B
        	m.mem[addr] = val
                op1    := uint16(m.mem[0x00]) + uint16(m.mem[0x01]) << 8
                op2    := uint16(m.mem[0x02]) + uint16(m.mem[0x03]) << 8

                result := uint32(op1 * op2)

                m.mem[0x04] = byte(result       & 0xff)
                m.mem[0x05] = byte(result >> 8  & 0xff)
                m.mem[0x06] = byte(result >> 16 & 0xff)
                m.mem[0x07] = byte(result >> 24 & 0xff)


        case 0x04, 0x05, 0x06, 0x07:   // UNSIGNED_MULT_result
		break


        case 0x08, 0x09, 0x0a, 0x0b:   // SIGNED_MULT_A, SIGNED_MULT_B
                m.mem[addr] = val
                op1    := int16(m.mem[0x08]) + int16(m.mem[0x09]) << 8
                op2    := int16(m.mem[0x0a]) + int16(m.mem[0x0b]) << 8

                result := int32(op1 * op2)

                m.mem[0x0c] = byte(result       & 0xff)
                m.mem[0x0d] = byte(result >> 8  & 0xff)
                m.mem[0x0e] = byte(result >> 16 & 0xff)
                m.mem[0x0f] = byte(result >> 24 & 0xff)


        case 0x0c, 0x0d, 0x0e, 0x0f:   // SIGNED_MULT_result
		break


        case 0x10, 0x11, 0x12, 0x13:   // UNSIGNED_DIV_DEM, UNSIGNED_DIV_NUM
                m.mem[addr] = val
                op1    := uint16(m.mem[0x10]) + uint16(m.mem[0x11]) << 8
                op2    := uint16(m.mem[0x12]) + uint16(m.mem[0x13]) << 8
                        
                var result, remainder uint16
                if (op1 != 0) {
                	result = op2 / op1
                        remainder = op2 % op1
                }

                m.mem[0x14] = byte(result          & 0xff)
                m.mem[0x15] = byte(result    >> 8  & 0xff)
                m.mem[0x16] = byte(remainder       & 0xff)
                m.mem[0x17] = byte(remainder >> 8  & 0xff)


        case 0x14, 0x15, 0x16, 0x17:   // UNSIGNED_DIV_result
		break


        case 0x18, 0x19, 0x1a, 0x1b:   // SIGNED_DIV_DEM, SIGNED_DIV_NUM
                m.mem[addr] = val
                op1    := int16(m.mem[0x18]) + int16(m.mem[0x19]) << 8
                op2    := int16(m.mem[0x1A]) + int16(m.mem[0x1B]) << 8
                        
                var result, remainder int16
                if (op1 != 0) {
                	result = op2 / op1
                	remainder = op2 % op1
                }

                m.mem[0x1C] = byte(result          & 0xff)
                m.mem[0x1D] = byte(result    >> 8  & 0xff)
                m.mem[0x1E] = byte(remainder       & 0xff)
                m.mem[0x1F] = byte(remainder >> 8  & 0xff)


        case 0x1c, 0x1d, 0x1e, 0x1f:   // SIGNED_DIV_result
		break


        case 0x20, 0x21, 0x22, 0x23,  // ADDER32_A
             0x24, 0x25, 0x26, 0x27:  // ADDER32_B

                m.mem[addr] = val
                op1    := int32(m.mem[0x20])       + 
                          int32(m.mem[0x21]) <<  8 + 
                          int32(m.mem[0x22]) << 16 + 
                          int32(m.mem[0x23]) << 24

                op2    := int32(m.mem[0x24])       + 
                          int32(m.mem[0x25]) << 8  + 
                          int32(m.mem[0x26]) << 16 + 
                          int32(m.mem[0x27]) << 24

                result := int32(op1 + op2)

                m.mem[0x28] = byte(result       & 0xff)
                m.mem[0x29] = byte(result >> 8  & 0xff)
                m.mem[0x2a] = byte(result >> 16 & 0xff) 
                m.mem[0x2b] = byte(result >> 24 & 0xff) 

        case 0x28, 0x29, 0x2a, 0x2b:   // ADDER32_result 
		break

        default:
		m.mem[addr] = val
        }

	return nil
}


func (m *MathInt) Name() string {
        return m.name
}

func (m *MathInt) Size() (uint32, uint32) {
        return 0x00, uint32(len(m.mem))
}

