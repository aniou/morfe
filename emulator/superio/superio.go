package superio

import (
        "fmt"

        "github.com/aniou/morfe/lib/queue"
)

const F_MAIN = 0

// base 0xAF_1000 for FMX
const KBD_DATA_BUF = 0x60       // for writing
const KBD_INPT_BUF = 0x60       // for reading
const PORT_B       = 0x61       // not used yet
const KBD_CMD_BUF  = 0x64       // for writing
const KBD_STATUS   = 0x64       // for reading

//    $AF:1060 - $AF:1064 - LOGIC DEVICE 7 - KEYBOARD
//    $AF:1100 - $AF:117F - LOGIC DEVICE A - PME (Runtime Registers)
//    $AF:1200 - $AF:1200 - LOGIC DEVICE 9 - GAME PORT
//    $AF:12F8 - $AF:12FF - LOGIC DEVICE 5 - SERIAL 2
//    $AF:1330 - $AF:1331 - LOGIC DEVICE B - MPU-401
//    $AF:1378 - $AF:137F - LOGIC DEVICE 3 - PARALLEL PORT
//    $AF:13F0 - $AF:13F7 - LOGIC DEVICE 0 - FLOPPY CONTROLLER
//    $AF:13F8 - $AF:13FF - LOGIC DEVICE 4 - SERIAL 1

type SIO struct {
        //out    chan byte       // for 'display'
        InBuf    queue.QueueByte // for 'keyboard'
        Data     byte
        command  byte
        name     string
        mem      []byte
}

func New(name string, size int) *SIO {
        s := SIO{InBuf: queue.NewQueueByte(200), 
                  Data: 0, 
               command: 0, 
                  name: name,
                   mem:  make([]byte, size)}
        return &s
}

func (s *SIO) Name(fn byte) string {
        return s.name
}

func (s *SIO) Clear() { 
}

func (s *SIO) Size(fn byte) (uint32, uint32) {
        return 0x00, uint32(len(s.mem))
}

func (s *SIO) Read(fn byte, addr uint32) (byte, error) {
        switch addr {
        case KBD_INPT_BUF:
                return s.Data, nil			// XXX something gone wrong
                if s.InBuf.Len() > 0 {
                        return *s.InBuf.Dequeue(), nil
                } else {
                        return 0, nil
                }

        case KBD_STATUS:
                return s.command, nil
        default:
                return 0, fmt.Errorf("superio: %4s Read  addr %6x is not implemented, 0 returned", s.name, addr)
        }
}

// taken from FoenixIDE
func (s *SIO) Write(fn byte, addr uint32, val byte) error {
        switch addr {
        case KBD_DATA_BUF:                      // AF:1060 FMX
                if val == 0x69 {                // 
                        s.command = 1           // 
                }
                if val == 0xEE {                // echo
                        s.command = 1           // 
                }
                if val == 0xF4 {                // kbd reset
                        s.InBuf.Enqueue(0xFA)   // self-test result
                        s.Data = 0xFA
                        s.command = 1
                }
                if val == 0xF6 {                // 
                        s.command = 1
                }
        case KBD_CMD_BUF:                       // 0xAF1064 FMX
                if val == 0x20 {                // 
                        s.command = 1
                }
                if val == 0x60 {                // 
                        s.command = 1
                }
                if val == 0xA8 {                // 
                        s.command = 1
                }
                if val == 0xA9 {                // 
                        s.InBuf.Enqueue(0x00)   // 
                        s.Data = 0
                        s.command = 0x01
                }
                if val == 0xAA {                // self-test
                        s.InBuf.Enqueue(0x55)   // self-test result
                        s.Data = 0x55
                        s.command = 0x01
                }
                if val == 0xAB {                // self-test
                        s.InBuf.Enqueue(0x00)   // 
                        s.Data = 0
                }
                if val == 0xD4 {                // 
                        s.command = 0x01
                }
                if val == 0x00 {
                        s.command = 0x00
                }
        default:
                return fmt.Errorf("superio: %4s Write addr %6x val %2x is not implemented", s.name, addr, val)
        }
        return nil
}

