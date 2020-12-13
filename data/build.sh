#!/bin/sh
#64tass --m65816 c256-shim.asm --long-address --flat --intel-hex -o c256-shim.hex --list c256-shim.lst

64tass -D TARGET=2 -D TARGET_SYS=1  -o graph5bm0.hex -L graph5bm0.lst --intel-hex  --m65816 graph5bm0.asm
