package pata

/* WARNING:
   Work-In-Progress code: so far only very selected subset of commands 
   provided, debug enabled and only R/O support
*/

import (
        "fmt"
	"os"
)

const (
	IDE_IDLE = iota
	IDE_CMD
	IDE_DATA_IN
	IDE_DATA_OUT
)

const (
	CMD_READ = 0x20
)

const (
        // controller errors (bits!)
        ERR_NO_ADDR_MARK       = byte(1)
        ERR_NO_TRACK0          = byte(2)
        ERR_COMMAND_ABORTED    = byte(4)
        ERR_MEDIA_CHANGE_REQ   = byte(8)
        ERR_NO_ID_MARK_FOUND   = byte(16)
        ERR_MEDIA_CHANGED      = byte(32)
        ERR_UNCORRECTABLE_DATA = byte(64)
        ERR_BAD_BLOCK          = byte(128)


        DEVH_HEAD = byte(15)
        DEVH_DEV  = byte(16)
        DEVH_LBA  = byte(64)

        BIT0       = byte(1)
        BIT1       = byte(2)
        BIT2       = byte(4)
        BIT3       = byte(8)
        BIT4       = byte(16)
        BIT5       = byte(32)
        BIT6       = byte(64)
        BIT7       = byte(128)

        // controller statusa (bits)
        ST_ERR       = byte(1)
        ST_IDX       = byte(2)
        ST_CORR      = byte(4)
        ST_DRQ       = byte(8)
        ST_DSC       = byte(16)		// drive seek complete
        ST_DF        = byte(32)
        ST_DRDY      = byte(64)
        ST_BSY       = byte(128)

        // addresses (address offset)
        REG_PATA_DATA       = uint32(0)  // data8 and data16
        REG_PATA_ERROR      = uint32(2)  // error on read, feature on write
        REG_PATA_SECT_CNT   = uint32(4)
        REG_PATA_SECT_SRT   = uint32(6)  // 06: LBA0: low
        REG_PATA_CLDR_LO    = uint32(8)  // 08: LBA1: med
        REG_PATA_CLDR_HI    = uint32(10) // 0a: LBA2: hi
        REG_PATA_DEVH       = uint32(12) // 0c: LBA3: top - bit 24 to 27 
        REG_PATA_CMD_STAT   = uint32(14) // 0e: command or status (write or read)
)

/*
Drive / Head Register 
Bit     Abbrev  Function
0 - 3           In CHS addressing, bits 0 to 3 of the head. 
                In LBA addressing, bits 24 to 27 of the block number.
4       DRV     Selects the drive number.
5       1       Always set.
6       LBA     Uses CHS addressing if clear or LBA addressing if set.
7       1       Always set.
*/


type DRIVE struct {
        lba_mode        bool   // LBA or CHS (no implemented)
        lba0            byte   // SECTOR      or  0:7   of LBA
        lba1            byte   // CYLINDER lo or  8:15  of LBA
        lba2            byte   // CYLINDER hi or 16:23  of LBA
        lba3            byte   // 0:3 of HEAD or 24:27  of LBA
        sector_count    byte   // parameter for operations

	command		byte
	status		byte
	err		byte
	state		byte

	fd		*os.File // file descriptor for image
	offset		uint32   // current file position
	data		[512*256]byte
	data_amount	int
	data_pointer    int
}

type PATA struct {
        name            string
        mem             []byte  // to conform with RAM interface

        selected        byte    // selected drive (0, 1)
        drive           [2]DRIVE

        log_level       uint
}

// for debug purposes
var REG = []string{
        "PATA_DATA lo",
        "PATA_DATA hi",
        "PATA_ERROR",
        "offset 0x03",
        "PATA_SECT_CNT",
        "offset 0x05",
        "PATA_SECT_SRT / LBA0",
        "offset 0x07",
        "PATA_CLDR_LO  / LBA1",
        "offset 0x09",
        "PATA_CLDR_HI  / LBA2",
        "offset 0x0b",
        "PATA_DEVH     / LBA3",
        "offset 0x0d",
        "PATA_CMD_STAT",
        "offset 0x0f",
}

const (
        LOG_FATAL = iota 
        LOG_ERROR   
        LOG_INFO    
        LOG_DEBUG   
        LOG_TRACE   
        LOG_VERBOSE 
)

var logLevelLookup = map[uint]string {
        LOG_FATAL   : "FATAL",
        LOG_ERROR   : "ERROR",
        LOG_INFO    : "INFO",
        LOG_DEBUG   : "DEBUG",
        LOG_TRACE   : "TRACE",
        LOG_VERBOSE : "VERBOSE",
}

func New(name string, size int) *PATA {
        s := PATA{
             log_level: LOG_FATAL,
                  name: name,
                   mem: make([]byte, size),
        }
	s.drive[0] = DRIVE { status: ST_DRDY, fd: nil }
	s.drive[1] = DRIVE { status: ST_DRDY, fd: nil }

        return &s
}

func (s *PATA) AttachDisk(number byte, path string) error {
	if path == "" {
		s.debug(LOG_DEBUG, "pata: %6s drive %d no filename, aborting attach disk\n", s.name, number)
		return nil
	}

	file, err := os.Open(path)	
	if err != nil {
		return err
	}
        s.debug(LOG_DEBUG, "pata: %6s drive %d file was succesfully opened %s\n", s.name, number, path)
	s.drive[number].fd = file
	return nil
}

func (s *PATA) DetachDisk(number byte) {
	if s.drive[number].fd != nil {
		s.drive[number].fd.Close()
		s.debug(LOG_DEBUG, "pata: %6s drive %d was closed\n", s.name, number)
	}
}

func (s *PATA) calculate_block() (int64, error) {
        var block_number int64

        drive := s.drive[s.selected]
        if (!drive.lba_mode) {
                return 0, fmt.Errorf("CHS not supported yet")
        }
        
        block_number = int64(drive.lba3) << 24 |
                       int64(drive.lba2) << 16 |
                       int64(drive.lba1) <<  8 |
                       int64(drive.lba0)

        return block_number, nil
}

func (s *PATA) cmd_read_sectors() {
	offset, err := s.calculate_block()

	drive := &s.drive[s.selected]
	if err != nil {
                s.debug(LOG_TRACE, "pata: %6s drive %d error %s\n", s.name, s.selected, err )
		drive.status  |= ST_ERR
		drive.status &^= ST_DSC
		drive.err     |= ERR_NO_ID_MARK_FOUND

		drive.status &^= (ST_BSY|ST_DRQ)
		drive.status  |= ST_DRDY
		drive.state    = IDE_IDLE
		return
	}

	_, err = drive.fd.Seek(offset * 512, 0)		// XXX - block size always as 512?
	if err != nil {
                s.debug(LOG_TRACE, "pata: %6s drive %d seek error %s\n", s.name, s.selected, err )
		drive.status  |= ST_ERR
		drive.status &^= ST_DSC
		drive.err     |= ERR_NO_ID_MARK_FOUND

		drive.status &^= (ST_BSY|ST_DRQ)
		drive.status  |= ST_DRDY
		drive.state    = IDE_IDLE
		return
	}

	drive.status  |= (ST_DRQ | ST_DSC | ST_DRDY)    // FoenixMCP required DATA_READY
	drive.status &^= ST_BSY

	data_to_read  := int(drive.sector_count) * 512
	if data_to_read == 0 {				// 0 means '256'
		data_to_read = 256 * 512    
	}
        _, err = drive.fd.Read(drive.data[0 : data_to_read])
	if err != nil {
                s.debug(LOG_ERROR, "pata: %6s drive %d error %s\n", s.name, s.selected, err )
		drive.status |= ST_ERR
		drive.err     = ERR_UNCORRECTABLE_DATA
		return
	}

	// there are data in buffer!
        s.debug(LOG_TRACE, "pata: %6s drive %d read %d bytes from offset %d\n", s.name, s.selected, data_to_read, offset )
	fmt.Printf("pata: >>> %v\n", drive.data[0 : data_to_read])
	drive.data_amount   = data_to_read
	drive.data_pointer  = 0
	drive.status      &^= ST_BSY
	drive.status       |= ST_DRQ
	drive.state         = IDE_DATA_IN
	return
}

func (s *PATA) get_data_from_buffer() (byte) {
	var retval byte

        drive := &s.drive[s.selected]
	// XXX - any error?
	if drive.state != IDE_DATA_IN {
                s.debug(LOG_ERROR, "pata: %6s drive %d read from empty buffer\n", s.name, s.selected )
		return 0
	}

	
	retval = drive.data[drive.data_pointer]
        //s.debug(LOG_ERROR, "pata: %6s drive %d pointer %d value %d\n", s.name, s.selected, drive.data_pointer, retval )
        drive.data_pointer += 1
        if drive.data_pointer >= drive.data_amount {
		drive.status &^= (ST_BSY|ST_DRQ)
		drive.status  |= ST_DRDY
		drive.state    = IDE_IDLE
	}
	return retval
}

func (s *PATA) debug(level uint, format string, a ...any) {
        if level >= s.log_level {
                fmt.Printf("%-7s ", logLevelLookup[level])
                fmt.Printf(format, a...)
        }
}

func (s *PATA) Name(fn byte) string {
        return s.name
}

func (s *PATA) Size(fn byte) (uint32, uint32) {
        return 0x00, uint32(len(s.mem))
}

func (s *PATA) Clear() { 
}

func (s *PATA) Read(fn byte, addr uint32) (byte, error) {


        switch addr {
        case REG_PATA_DATA:
		val := s.get_data_from_buffer()
                //s.debug(LOG_TRACE, "pata: %6s drive %d read lo16 0x%02x from buffer\n", s.name, s.selected, val)
		return val, nil
        case REG_PATA_DATA+1:
		val := s.get_data_from_buffer()
                //s.debug(LOG_TRACE, "pata: %6s drive %d read hi16 0x%02x from buffer\n", s.name, s.selected, val)
		return val, nil
        case REG_PATA_CMD_STAT: // 0x0e - check status when read
                s.debug(LOG_TRACE, "pata: %6s read  0x%02x from %13s\n", s.name, s.drive[s.selected].status, REG[addr])
                return s.drive[s.selected].status, nil
        default:
                return 0, fmt.Errorf("pata: %6s Read  addr %6x is not implemented, 0 returned", s.name, addr)
        }
}

func (s *PATA) Write(fn byte, addr uint32, val byte) error {

        switch addr {
        case REG_PATA_CMD_STAT: // 0x0e - issue command when write

	        drive         := &s.drive[s.selected]
		drive.status &^= (ST_ERR|ST_DRDY)		// clear ERR and READY
		drive.status  |=  ST_BSY
                drive.err      = 0
                drive.state    = IDE_CMD
		drive.command  = val				// just for sake

                switch val {
                case 0x00: 
                        s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s (NOP)\n", s.name, val, REG[addr])
			drive.status  &^=  ST_BSY
			drive.status   |=  ST_DRDY
                case CMD_READ:  // 0x20
                        s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s (READ SECT)\n", s.name, val, REG[addr])
			s.cmd_read_sectors()
                default:
                        s.debug(LOG_ERROR, "pata: %6s write 0x%02x to   %-22s (unknown)\n", s.name, val, REG[addr])
                }

        case REG_PATA_SECT_CNT:  // 0x04
                s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s\n", s.name, val, REG[addr])
                s.drive[s.selected].sector_count = val

        case REG_PATA_SECT_SRT:  // 0x06
                s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s\n", s.name, val, REG[addr])
                s.drive[s.selected].lba0         = val

        case REG_PATA_CLDR_LO:   // 0x08
                s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s\n", s.name, val, REG[addr])
                s.drive[s.selected].lba1         = val

        case REG_PATA_CLDR_HI:   // 0x0a
                s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s\n", s.name, val, REG[addr])
                s.drive[s.selected].lba2         = val

        case REG_PATA_DEVH: // 0x0c
                s.debug(LOG_TRACE, "pata: %6s write 0x%02x to   %-22s\n", s.name, val, REG[addr])

                if (val & DEVH_DEV) > 0 {
                        s.selected = 1
                } else {
                        s.selected = 0
                }
                s.drive[s.selected].lba_mode = (val & DEVH_LBA) > 0
                s.drive[s.selected].lba3     =  val & DEVH_HEAD // bits 0:3


                s.debug(LOG_TRACE, "pata: %6s mode drive %d LBA %t lba3 %d\n", 
                        s.name, s.selected, s.drive[s.selected].lba_mode, s.drive[s.selected].lba3)

        default:
                return fmt.Errorf("pata: %6s Write addr %6x val %2x is not implemented", s.name, addr, val)
        }
        return nil
}

