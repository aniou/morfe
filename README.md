# go65c816
65c816 emulator in Go

*The project is at very early stage and lacks many of CPU and debugger features*

[![asciicast](https://asciinema.org/a/270744.svg)](https://asciinema.org/a/270744)

## Supported systems

Program was tested on:

* NetBSD 9 / Go 1.12
* Ubuntu 18.04 / Go 1.13
 
Program should work on:

* MS Windows / Go >=1.12

## Installing

### version A

```bash
git clone https://github.com/aniou/go65c816
cd go65c816
go run cmd/go65c816/main.go

# go to another terminal and run
go run cmd/netcon/main.go    
```

### version B

```bash
go get github.com/aniou/go65c816/cmd/go65c816
go get github.com/aniou/go65c816/cmd/netcon
~/go/bin/go65c816

# spawn another terminal emulator and run console
~/go/bin/netcon
```

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

## Input/Output

By default emulator provides simplest I/O via TCP socket opened at `localhost:12321`. Every byte written in emulator to addr `0xDF77` should be sent and every received byte is available from `0xDF75`. Buffer sizes for both directons are arbitraly set at 200 bytes.

Almost every program (telnet or nc) should work as client, but best results should be ahieved by client that sends data in character mode (i.e. every pressed key sends one byte). There is available simple client, called `netcon`.

## Memory map

Machine parameters may be tweaked by editing `emulator/platform/platform.go` file. Every memory area should be attached to internal bus, like in following example:

```go
        bus, _          := bus.New(logger)                     // new bus
        platform.CPU, _  = cpu65c816.New(bus)                  // new CPU
        console, _      := netconsole.NewNetConsole(logger)    // netconsole - first IO device
        ram, _          := memory.New(0x40000)                 // regular RAM, 256kB

        platform.CPU.Bus.Attach(ram,            "ram", 0x000000, 0x03FFFF)    // atach first region
        platform.CPU.Bus.Attach(console, "netconsole", 0x00DF00, 0x00DFFF)    // mask area 0xdf00-0xdfff
```

 * minimal area size: **16 bytes**
 * areas **must be** aligned at 4 bits (16 bytes)
 * areas are stacked

### TODO
 * interrupts
 * speed limit
 * real graphics
 * breakpoints
 * conditional breakpoints
 * additional commands
 * performance improvements

### Foreword

Project was inspired by [NES emulator](https://github.com/fogleman/nes) created by Michael Fogleman and [MOS 6502 emulator](https://github.com/pda/go6502) by Paul Annesley and contains files or concepts from both projects. Some algorithms and behaviours are modeled on the [C++ 65c816 emulator](https://github.com/andrew-jacobs/emu816) by Andrew Jacobs.
