package ps2

// https://wiki.osdev.org/%228042%22_PS/2_Controller

import (
        "fmt"
)

const (
         PS2_STAT_OBF    = byte(0x01)
         PS2_STAT_IBF    = byte(0x02)
         PS2_STAT_SYS    = byte(0x04)
         PS2_STAT_CMD    = byte(0x08)
         PS2_STAT_INH    = byte(0x10)
         PS2_STAT_TTO    = byte(0x20)
         PS2_STAT_RTO    = byte(0x40)
         PS2_STAT_PE     = byte(0x80)
)

const KBD_DATA     = 0x00 // 0x60 for reading and writing
const KBD_COMMAND  = 0x04 // 0x64 for writing
const KBD_STATUS   = 0x04 // 0x64 for reading

type PS2 struct {
        name            string
        mem             []byte  // to conform with RAM interface

        data            byte    // data (usually keycode)
        status          byte    // controller status
        CCB             byte    // controller configuration byte
        ccb_write_mode  bool    // denotes that next write should go to CCB

        First_enabled   bool
        Second_enabled  bool

        debug_status    bool    // temporary
}

func New(name string, size int) *PS2 {
        s := PS2{
                status: 0, 
          debug_status: true,
                   CCB: 0,
                  name: name,
                   mem:  make([]byte, size)}

        return &s
}

func (s *PS2) Name(fn byte) string {
        return s.name
}

func (s *PS2) Size(fn byte) (uint32, uint32) {
        return 0x00, uint32(len(s.mem))
}

func (s *PS2) Clear() { 
}

func (s *PS2) Read(fn byte, addr uint32) (byte, error) {
        switch addr {
        case KBD_DATA:          // 0x60
                fmt.Printf("ps2: %6s read     KBD_DATA: val %02x\n", s.name, s.data)
                s.status = s.status & ^PS2_STAT_OBF
                return s.data, nil

        case KBD_STATUS:        // 0x64 
                if s.debug_status {
                        fmt.Printf("ps2: %6s read   KBD_STATUS: val %02x\n", s.name, s.status)
                }
                return s.status, nil
        default:
                return 0, fmt.Errorf("ps2: %6s Read  addr %6x is not implemented, 0 returned", s.name, addr)
        }
}

func (s *PS2) Write(fn byte, addr uint32, val byte) error {
        switch addr {
        case KBD_DATA: // 0x60
                if s.ccb_write_mode {
                        fmt.Printf("ps2: %6s write    KBD_DATA: val %02x -> CCB\n", s.name, val)

                        s.ccb_write_mode  = false
                        s.CCB             = val
                        return nil
                }

                switch val {
                case 0xf4: // mouse/keyboard enable
                        s.status = s.status | PS2_STAT_OBF 
                        s.debug_status = false	// to get rid console messages in case of pooling
		case 0xf5: // mouse/keyboard disable
                        s.status = s.status | PS2_STAT_OBF 
		case 0xf6: // mouse - reset without self-test
                case 0xff: // mouse/keyboard reset
                        s.status = s.status | PS2_STAT_OBF 

                default:
                        fmt.Printf("ps2: %6s write    KBD_DATA: val %02x - data UNKNOWN\n", s.name, val)
			return nil
                }

                fmt.Printf("ps2: %6s write    KBD_DATA: val %02x\n", s.name, val)


                /*
                if val == 0x69 {                // 
                        s.command = 1           // 
                }
                if val == 0xEE {                // echo
                        s.command = 1           // 
                }
                if val == 0xF4 {                // kbd reset
                          // self-test result
                        s.data = 0xFA
                        s.command = 1
                }
                if val == 0xF6 {                // 
                        s.command = 1
                }
                */

        case KBD_COMMAND: // 0x64 - command when write
                fmt.Printf("ps2: %6s write KBD_COMMAND: val %02x\n", s.name, val)
                switch val {
                case 0x60:
                        s.ccb_write_mode    = true
		case 0xd4: // write next byte to second PS/2 port
                        s.status = s.status | PS2_STAT_OBF 
                case 0xa7: // disable second PS/2 port
                        s.status = s.status | PS2_STAT_OBF 
                        s.Second_enabled = false
                case 0xa8: // enable second PS/2 port
                        s.Second_enabled = true
		case 0xa9: // test second PS/2 port
			s.data = 0x00
                        s.status = s.status | PS2_STAT_OBF 
		case 0xaa: // test PS/2 controller
                        s.data = 0x55
                        s.status = s.status | PS2_STAT_OBF 
		case 0xab: // test first PS/2 port
                        s.data = 0x00
                        s.status = s.status | PS2_STAT_OBF 
                case 0xad: // disable first PS/2 port
                        s.status = s.status | PS2_STAT_OBF 
                        s.First_enabled = false
                case 0xae: // enable first PS/2 port
                        s.First_enabled = true

                default:
                        fmt.Printf("ps2: %6s write KBD_COMMAND: val %02x - command UNKNOWN\n", s.name, val)
                }
                /*
                if val == 0x20 {                // 
                        s.command = 1
                }
                if val == 0x60 {                // Write next byte to "byte 0" of internal RAM (Controller Configuration Byte)
                        s.ccb_write_mode    = true
                }
                if val == 0xA8 {                // 
                        s.command = 1
                }
                if val == 0xA9 {                // 
                        s.data = 0x00
                        s.command = 0x01
                }
                }
                if val == 0xD4 { 
                        s.command = 0x01
                }
                if val == 0x00 {
                        s.command = 0x00
                }
                */
        default:
                return fmt.Errorf("ps2: %6s Write addr %6x val %2x is not implemented", s.name, addr, val)
        }
        return nil
}

func (s *PS2) AddKeyCode(val byte) {
        fmt.Printf("ps2: AddKeyCode: %02x\n", val)
        s.data = val
        s.status = s.status | PS2_STAT_OBF 
}
