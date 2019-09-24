# go65c816
65c816 emulator in Go

## Keybindings

|Window  |Key       |Meaning|
---------|----------|--------
any      |TAB       |next window
any      |Ctrl+Space|run/stop CPU
any      |Ctrl+C    |exit emulator
any      |Ctrl+Q    |exit emulator
any      |Ctrl+P    |load hex file data/program.hex
any      |~ (tilde) |switch to/from command window
command  |Enter     |execute command
log, code|UP arrow  |cursor up
log, code|DOWN arrow|cursor down
code     |Space     |execute one step

## Commands

All values should be provided in hexadecimal form.

|Command           | Meaning |
-------------------|----------
|set mem [addr]    |set memory dump window do addr
|load hex [path]   |load program in hex format 
|run               |run/stop CPU
|peek, peek8 [addr]|peek one byte 
|peek16 [addr]     |peek word (without wraping at bank boundary) 
|peek24 [addr]     |peek 24-bit value (without wrapping at bank boundary)
|quit              |quit emulator

### TODO
 * interrupts
 * speed limit
 * real graphics
 * breakpoints
 * conditional breakpoints
 * additional commands
 * performance improvements
