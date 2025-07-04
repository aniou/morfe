# MORFE - Meyer's Own Re-combinable FrankenEmulator

> [!IMPORTANT]
> You probably should take a look at [MORFE/O](https://github.com/aniou/morfeo)
> that already is - or will be in short time - more feature rich and way faster
> than this implementation. Not only for m68k but also for C256 line!
> 
> This repo is rather abandoned and lost sync with current state of hardware.

---

MORFE is a kind-of-emulator software, created for my experiments and 
development for [Foenix](https://c256foenix.com/) machines.

At this moment You should also consider following factors:

* **MOST IMPORTANT:** it is unofficial work and features/lack 
  of features or emulator design does not correspond to features 
  or design of real FMX/U/GenX machine! 

  Do not make any unauthorized assumptions about real hardware!

* In this branch m68k is easily able to achieve 25Mhz, although now
  it is capped at 20Mhz. I have a plan for making all setting easily
  configurable, but first things (like debug facilities for both cpus)
  first...
  
## Build instructions for this branch

First, clone repo by ``git clone https://github.com/aniou/morfe``

Type ``make morfe`` for a 65c816-only version or ``make morfe-m68k`` 
for dual, 65c816/m68k one (but You should stick with 65c816 and use
[MORFE/O](https://github.com/aniou/morfeo) for m68k).

There is also ``make help`` that shows actual targets.

Typical session:

```shell
git clone https://github.com/aniou/morfe
cd morfe
make morfe
./morfe conf/c256.ini
```

### Note about m68k

For m68k a [Musashi](https://github.com/kstenerud/Musashi/) core
is used, built-in into emulator with some black magic around ``cgo``. 
Standard makefile should build all object files for You, although
working gcc will be necessary.

## Running

Emulator requires config file, that defines platform behaviour, files 
to be loaded at start and some other parameters.

Binaries should be run from project directory (all paths are relative
to top-directory of project), i.e.:

```shell
./morfe conf/c256.ini
./morfe-m68k conf/a2560k.ini
```

**Use ``F8`` to change active screen! By default [FoenixMCP](https://github.com/pweingar/FoenixMCP)
shows debug on head0 and console on head1**


## Built-in debugger

Debugger interface is available for m68k only. Press ``F9`` to call
debug window in terminal. Preferred terminal size is ``132x42`` 
(interface can be scaled on fly).

List of supported commands will be displayed in log frame.

## Compatibility status

### General

- a preliminary DIP-switch support exists, so far only DIP6 (HI/LO graphics mode
  selector) is implemented, see ``conf/c256.ini`` for examples

### Memory map

At this moment a sort-of FMX memory map is available, but GenX is on the horizont:
it is fast moving target, so stay tuned!

### Vicky II/III

See [here](https://wiki.c256foenix.com/index.php?title=VICKY_II) for VICKY II spec

- [x] preliminary PATA HDD support (r/o, only for a2560k, very rudimentary)
- [x] 640x480 mode
- [x] 800x600 mode (from 19.09.2021)
- [x] 1024x768, 640x400 - only for m68k-based machines
- [ ] double pixel mode
- [x] fullscreen mode
- [x] border support (partial, no scroll)
- [x] text mode 
- [x] text LUT
- [x] cursor 
- [x] fonts
- [x] bm0 bitmap
- [x] bm1 bitmap
- [x] bitmap LUT
- [x] overlay and background support
- [ ] tiles
- [ ] sprites
- [ ] GAMMA LUT
- [ ] 8-bit writes (Vicky writes are 8-bit even if A/X/Y are 16-bits wide)

### GABE

See [here](https://wiki.c256foenix.com/index.php?title=GABE) for GABE spec

- [x] math coprocessor
- [x] PS/2 keyboard input 
- [ ] mouse
- [ ] all other

### general features

- [x] IRQ (partial: only 65c816 mode)
- [ ] NMI
- [ ] reset button

## Keybindings

There are few keybindings now. 
*Warning:* following keys aren't passed to emulator!

|Key     |Effect
---------|---------------------------
F8       |Change active head in multi-head setups
F9       |Enter m68k debugger
F10      |- (nothing)
F11      |Toggle full-screen
F12      |Exit emulator

## Foreword

I owe thankful word for too many people. Excuse me if I omitted someone.

First at all: all hail to [Stefany Allaire](https://twitter.com/stefanyallaire), 
a Dark Mistress that brought to life all Foenix Family! We all praise her
brilliant work, persistence and vision!

Project was inspired by [NES emulator](https://github.com/fogleman/nes) 
created by Michael Fogleman and general layout as well as architectural
concepts are based on that project. I'm very grateful to Michael for 
inspiration and all things I learned from their code.

During development a 65c816 emulation I draw inspiration and concepts
from Michael's project as well as from [MOS 6502 emulator](https://github.com/pda/go6502) 
by Paul Annesley. Some algorithms and behaviours are modeled on the 
[C++ 65c816 emulator](https://github.com/andrew-jacobs/emu816) by Andrew Jacobs.

Project also draws inspirations, knowledge about Foenix's behaviour and even 
whole code snippets from [Foenix IDE](https://github.com/Trinity-11/FoenixIDE) 
by Daniel Tremblay.

When I was in doubt (usually) I was able to find solution and hints in 
[Foenix Kernel Code](https://github.com/Trinity-11/Kernel_FMX/) created and 
maintained by https://github.com/pweingar/

All code for Motorola is provided by [Musashi core](https://github.com/kstenerud/Musashi/).

Finally: daschewie - thanks for Your support!

