.cpu "65816"

; Copyright 2019 Piotr Meyer <aniou@smutek.pl>
;
; Permission to use, copy, modify, and/or distribute this software for any
; purpose with or without fee is hereby granted, provided that the above 
; copyright notice and this permission notice appear in all copies.

; THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES 
; WITH REGARD TO THIS SOFTWARE  INCLUDING ALL IMPLIED WARRANTIES OF
; MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR 
; ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES 
; WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN 
; ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF 
; OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

; dead simple routines that reduces difference between C256 FMX kernel
; and go65c816 emulator in current shape

; custom INIT

    * = $1000

    clc
    xce
    ;jsr $1100
    ;jmp $2000
    ;WDM  #1
    jml $3a0000

; PUTC-like

	* = $1018

    ; php/pla in this place is somewhat redundant
    ; see of816 platform-specific routines
    ;wdm #$10
    php
    phb             ; preserve DBR too!
    sep #$20        ; set A short
    pha             ; preserve
    lda #$00        ; set DBR to $00
    pha
    plb
    pla             ; restore A
	sta $efe
    plb
    plp
	rtl

; GETC-like

	* = $104c

    ; php/pla in this place is somewhat redundant
    ; see of816 platform-specific routines
    php
    phb             ; preserve DBR too!
    sep #$20        ; set A short
    lda #$00        ; set DBR to $00
    pha
    plb
-	lda $f00
	beq -
    plb
    plp
	rtl

; BANNER-like

   * = $1100

   sep #$20        ; set A short
   .xl
   ldx #$0000

b0 lda banner, x
   beq b1
   jsl $1018
   inx
   bra b0
b1 jsl $104c
   rts

banner .text "press any key...", $0a, $00

; COMMAND PARSER Variables
; Command Parser Stuff between $000F00 -> $000F84 (see CMD_Parser.asm)
; KEY_BUFFER       = $000F00 ; 64 Bytes keyboard buffer
; KEY_BUFFER_SIZE  = $0080   ;128 Bytes (constant) keyboard buffer length
; KEY_BUFFER_END   = $000F7F ;  1 Byte  Last byte of keyboard buffer
; KEY_BUFFER_CMD   = $000F83 ;  1 Byte  Indicates the Command Process Status
; COMMAND_SIZE_STR = $000F84 ;  1 Byte
; COMMAND_COMP_TMP = $000F86 ;  2 Bytes
; KEYBOARD_SC_FLG  = $000F87 ;  1 Bytes that indicate the Status of Left Shift, Left CTRL, Left ALT, Right Shift
; KEYBOARD_SC_TMP  = $000F88 ;  1 Byte, Interrupt Save Scan Code while Processing
; KEYBOARD_LOCKS   = $000F89 ;  1 Byte, the status of the various lock keys
; KEYFLAG          = $000F8A ;  1 Byte, flag to indicate if CTRL-C has been pressed
; KEY_BUFFER_RPOS  = $000F8B ;  2 Byte, position of the character to read from the KEY_BUFFER
; KEY_BUFFER_WPOS  = $000F8D ;  2 Byte, position of the character to write to the KEY_BUFFER
;
; KERNEL_JMP_BEGIN = $001000 ; Reserved for the Kernel jump table
; KERNEL_JMP_END   = $001FFF

