===========================
VICKY - the graphics engine
===========================

.. contents::

"Vicky" is the name of the graphics engine of the C256 Foenix.

"Vicky II" is the upgraded version of the graphic engines that resides in the
FMX Version.

The major difference in between VICKY and VICKY II is the external bus @ 32Bits
wide. (VICKY is only 16Bits) The access time has been doubled which allow the
chipset to have better access time and overall better performance. It also gave
enough room to implement larger resolutions with faster GUI block move.

Part Number: CFP9551

It supports:

 * 320x240@60FPS Max Resolution @ 256 Colors (1 Byte per pixel)
 * 400x300@60FPS Max Resolution @ 256 Colors (1 Byte per pixel)
 * 640x480@60FPS Max Resolution @ 256 Colors (1 Byte per pixel)
 * 800x600@60FPS Max Resolution @ 256 Colors (1 Byte per pixel) (GUI mode)

 * 32 Sprites with a resolution of 32x32 pixels
 * 4 layers of tiles with a resolution of 16x16 pixels for each tile
 * Text Mode

During the video layer composition, Vicky has to read 1 line of bitmap (640
pixels), 1 line of tiles for each layer (4x 640 pixels) and all the lines that
are part of the 32 Sprites (worst case scenario 32 lines of 32 pixels), that
are displayed on that line. Then a composition and priority encoding are done.
In order to establish a priority, in other words, to know which pixel will be
in front, Vicky needs to store 10 lines of 640 pixels then scan the lot and
determine which one will be displayed.

640 + 640 + 640 + 640 + 640 + (32 * 32) = 4224 Pixels are to be read during
a single line interval. The pixel rate is 200Mbytes/sec, 5ns, so the overall
operation takes 21.12us, without calculating the overhead in the FPGA to go
from one process to the other.

The pixel index value $00 is always transparent, regardless if it's bitmap,
tile or sprite. The respective values of the first 4 bytes represented in the
LUT are thus always ignored.

Despite the fact that there is an ALPHA value in the LUT, it is not
supported/used at all.

Vicky Global Memory Map: 

Vicky Global Memory Map
=======================

=== =============   ============  ======================= ===============
ok  Start Address   Stop Address  Register Description    Additional Info
=== =============   ============  ======================= ===============
    $AF:0000        $AF:0000      MASTER_CTRL_REG_L       `Master Control Registers`_
    $AF:0001        $AF:0001      MASTER_CTRL_REG_H
    $AF:0002        $AF:0002      VKY_RESERVED_00
    $AF:0003        $AF:0003      VKY_RESERVED_01
=== =============   ============  ======================= ===============

.. warning::
   still vicky I

=== =============   ============  ======================= ===============
ok  Start Address   Stop Address  Register Description    Additional Info
=== =============   ============  ======================= ===============
r   $AF:0004        $AF:0004      BORDER_CTRL_REG         `Border Control Registers`_
 w  $AF:0005        $AF:0005      BORDER_COLOR_B 
 w  $AF:0006        $AF:0006      BORDER_COLOR_G
 w  $AF:0007        $AF:0007      BORDER_COLOR_R
rw  $AF:0008        $AF:0008      BORDER_X_SIZE[5:0]
rw  $AF:0009        $AF:0009      BORDER_Y_SIZE[5:0]
    $AF:000A        $AF:000C      UNDEFINED - RESERVED
    $AF:000D        $AF:000D      BACKGROUND_COLOR_B      Background Control Registers
    $AF:000E        $AF:000E      BACKGROUND_COLOR_G
    $AF:000F        $AF:000F      BACKGROUND_COLOR_R
~~  $AF:0010        $AF:0010      VKY_TXT_CURSOR_CTRL_REG Cursor Control Registers
    $AF:0011        $AF:0011      VKY_TXT_RESERVED
rw  $AF:0012        $AF:0012      VKY_TXT_CURSOR_CHAR_REG Char used as the Cursor
rw  $AF:0013        $AF:0013      VKY_TXT_CURSOR_COLR_REG Color Choice for the Cursor
rw  $AF:0014        $AF:0014      VKY_TXT_CURSOR_X_REG_L  Cursor Position X
rw  $AF:0015        $AF:0015      VKY_TXT_CURSOR_X_REG_H
rw  $AF:0016        $AF:0016      VKY_TXT_CURSOR_Y_REG_L  Cursor Position Y
rw  $AF:0017        $AF:0017      VKY_TXT_CURSOR_Y_REG_H
    $AF:0018        $AF:001B      UNDEFINED - RESERVED
    $AF:001C        $AF:001C      VKY_INFO_CHIP_NUM_L     Vicky Chip Part Number
    $AF:001D        $AF:001D      VKY_INFO_CHIP_NUM_H
    $AF:001E        $AF:001E      VKY_INFO_CHIP_VER_L     Vicky Chip Version
    $AF:001F        $AF:001F      VKY_INFO_CHIP_VER_H
    $AF:0020        $AF:00FF      UNDEFINED - RESERVED
    $AF:0100        $AF:0100      TL0_CONTROL_REG         Tile 0 Register Set
    $AF:0101        $AF:0101      TL0_START_ADDY_L        Start Address Within the Video Memory
    $AF:0102        $AF:0102      TL0_START_ADDY_M
    $AF:0103        $AF:0103      TL0_START_ADDY_H
    $AF:0104        $AF:0104      TL0_SCROLL_X
    $AF:0105        $AF:0105      TL0_SCROLL_Y
    $AF:0106        $AF:0107      UNDEFINED - RESERVED
    $AF:0108        $AF:0108      TL1_CONTROL_REG         Tile 1 Register Set
    $AF:0109        $AF:0109      TL1_START_ADDY_L        Start Address Within the Video Memory
    $AF:010A        $AF:010A      TL1_START_ADDY_M
    $AF:010B        $AF:010B      TL1_START_ADDY_H
    $AF:010C        $AF:010C      TL1_SCROLL_X
    $AF:010D        $AF:010D      TL1_SCROLL_Y
    $AF:010E        $AF:010F      UNDEFINED - RESERVED
    $AF:0110        $AF:0110      TL2_CONTROL_REG         Tile 2 Register Set
    $AF:0111        $AF:0111      TL2_START_ADDY_L        Start Address Within the Video Memory
    $AF:0112        $AF:0112      TL2_START_ADDY_M
    $AF:0113        $AF:0113      TL2_START_ADDY_H
    $AF:0114        $AF:0114      TL2_SCROLL_X
    $AF:0115        $AF:0115      TL2_SCROLL_Y
    $AF:0116        $AF:0117      UNDEFINED - RESERVED
    $AF:0118        $AF:0118      TL3_CONTROL_REG         Tile 3 Register Set
    $AF:0119        $AF:0119      TL3_START_ADDY_L        Start Address Within the Video Memory
    $AF:011A        $AF:011A      TL3_START_ADDY_M
    $AF:011B        $AF:011B      TL3_START_ADDY_H
    $AF:011C        $AF:011C      TL3_SCROLL_X
    $AF:011D        $AF:011D      TL3_SCROLL_Y
    $AF:011E        $AF:011F      UNDEFINED - RESERVED
    $AF:0120        $AF:013F      UNDEFINED - RESERVED
    $AF:0140        $AF:0140      BM_CONTROL_REG          Bitmap Registers Set
    $AF:0141        $AF:0141      BM_START_ADDY_L         Start Address Within the Video Memory
    $AF:0142        $AF:0142      BM_START_ADDY_M
    $AF:0143        $AF:0143      BM_START_ADDY_H
    $AF:0144        $AF:0144      BM_X_SIZE_L             Needs to be set to 640
    $AF:0145        $AF:0145      BM_X_SIZE_H
    $AF:0146        $AF:0146      BM_Y_SIZE_L             Needs to be set to 480
    $AF:0147        $AF:0147      BM_Y_SIZE_H
    $AF:0148        $AF:014F      BM_RESERVED
    $AF:0150        $AF:01FF      UNDEFINED - RESERVED
    $AF:0200        $AF:0200      SP00_CONTROL_REG        Sprite 0 (Highest Priority)
    $AF:0201        $AF:0201      SP00_ADDY_PTR_L         Start Address Within the Video Memory
    $AF:0202        $AF:0202      SP00_ADDY_PTR_M
    $AF:0203        $AF:0203      SP00_ADDY_PTR_H
    $AF:0204        $AF:0204      SP00_X_POS_L
    $AF:0205        $AF:0205      SP00_X_POS_H
    $AF:0206        $AF:0206      SP00_Y_POS_L
    $AF:0207        $AF:0207      SP00_Y_POS_H
    $AF:0208        $AF:020F      SP01 - Sprite 1
    $AF:0210        $AF:0217      SP02 - Sprite 2
    $AF:0218        $AF:021F      SP03 - Sprite 3
    $AF:0220        $AF:0227      SP04 - Sprite 4
    $AF:0228        $AF:022F      SP05 - Sprite 5
    $AF:0230        $AF:0237      SP06 - Sprite 6
    $AF:0238        $AF:023F      SP07 - Sprite 7
    $AF:0240        $AF:0247      SP08 - Sprite 8
    $AF:0248        $AF:024F      SP09 - Sprite 9
    $AF:0250        $AF:0257      SP10 - Sprite 10
    $AF:0258        $AF:025F      SP11 - Sprite 11
    $AF:0260        $AF:0267      SP12 - Sprite 12
    $AF:0268        $AF:026F      SP13 - Sprite 13
    $AF:0270        $AF:0277      SP14 - Sprite 14
    $AF:0278        $AF:027F      SP15 - Sprite 15
    $AF:0280        $AF:0287      SP16 - Sprite 16
    $AF:0288        $AF:028F      SP17 - Sprite 17
    $AF:0290        $AF:0297      SP18 - Sprite 18
    $AF:0298        $AF:029F      SP19 - Sprite 19
    $AF:02A0        $AF:02A7      SP20 - Sprite 20
    $AF:02A8        $AF:02AF      SP21 - Sprite 21
    $AF:02B0        $AF:02B7      SP22 - Sprite 22
    $AF:02B8        $AF:02BF      SP23 - Sprite 23
    $AF:02C0        $AF:02C7      SP24 - Sprite 24
    $AF:02C8        $AF:02CF      SP25 - Sprite 25
    $AF:02D0        $AF:02D7      SP26 - Sprite 26
    $AF:02D8        $AF:02DF      SP27 - Sprite 27
    $AF:02E0        $AF:02E7      SP28 - Sprite 28
    $AF:02E8        $AF:02EF      SP29 - Sprite 29
    $AF:02F0        $AF:02F7      SP30 - Sprite 30
    $AF:02F8        $AF:02FF      SP31 - Sprite 31
    $AF:0300        $AF:03FF      UNDEFINED - RESERVED    
    $AF:0400        $AF:0400      VDMA_CONTROL_REG        Video DMA Block
    $AF:0401        $AF:0401      VDMA_COUNT_REG_L
    $AF:0402        $AF:0402      VDMA_COUNT_REG_M
    $AF:0403        $AF:0403      VDMA_COUNT_REG_H
    $AF:0404        $AF:0404      VDMA_DATA_2_WRITE_L
    $AF:0405        $AF:0405      VDMA_DATA_2_WRITE_H
    $AF:0406        $AF:0406      VDMA_STRIDE_L
    $AF:0407        $AF:0407      VDMA_STRIDE_H
    $AF:0408        $AF:0408      VDMA_SRC_ADDY_L
    $AF:0409        $AF:0409      VDMA_SRC_ADDY_M
    $AF:040A        $AF:040A      VDMA_SRC_ADDY_H
    $AF:040B        $AF:040B      VDMA_RESERVED_0
    $AF:040C        $AF:040C      VDMA_DST_ADDY_L
    $AF:040D        $AF:040D      VDMA_DST_ADDY_M
    $AF:040E        $AF:040E      VDMA_DST_ADDY_H
    $AF:040F        $AF:040F      VDMA_RESERVED_1
    $AF:0410        $AF:04FF      UNDEFINED - RESERVED    
    $AF:0500        $AF:05FF      MOUSE_PTR_GRAPH0        16x16 Mem Block 0 for Mouse Pointer
    $AF:0600        $AF:06FF      MOUSE_PTR_GRAPH1        16x16 Mem Block 1 for Mouse Pointer
    $AF:0700        $AF:0700      MOUSE_PTR_CTRL_REG_L    Mouse Pointer Registers Block
    $AF:0701        $AF:0701      MOUSE_PTR_CTRL_REG_H
    $AF:0702        $AF:0702      MOUSE_PTR_X_POS_L       X Absolute Location of the Mouse
    $AF:0703        $AF:0703      MOUSE_PTR_X_POS_H       Presently Read Only
    $AF:0704        $AF:0704      MOUSE_PTR_Y_POS_L       Y Absolute Location of the Mouse
    $AF:0705        $AF:0705      MOUSE_PTR_Y_POS_H       Presently Read Only
    $AF:0706        $AF:0706      MOUSE_PTR_BYTE0         PS2 Mouse Packet Byte 0
    $AF:0707        $AF:0707      MOUSE_PTR_BYTE1         PS2 Mouse Packet Byte 1
    $AF:0708        $AF:0708      MOUSE_PTR_BYTE2         PS2 Mouse Packet Byte 2
    $AF:0709        $AF:070A      UNDEFINED MOUSE
    $AF:070B        $AF:070B      C256F_MODEL_MAJOR
    $AF:070C        $AF:070C      C256F_MODEL_MINOR
    $AF:070D        $AF:070D      FPGA_DOR                (Date of Release)
    $AF:070E        $AF:070E      FPGA_MOR                (Date of Release)
    $AF:070F        $AF:070F      FPGA_YOR                (Date of Release)
    $AF:0710        $AF:07FF      UNDEFINED - RESERVED
    $AF:0800        $AF:080F      RTC                     See the RTC Section for more details
    $AF:0810        $AF:0FFF      UNDEFINED - RESERVED
 !  $AF:1000        $AF:13FF      SUPERIO                 See the Super IO Section for more details
    $AF:1400        $AF:1F3F      UNDEFINED - RESERVED
 w  $AF:1F40        $AF:1F7F      FG_CHAR_LUT_PTR         Text Foreground Look-Up Table
 w  $AF:1F80        $AF:1FFF      BG_CHAR_LUT_PTR         Text Background Look-Up Table
    $AF:2000        $AF:23FF      GRPH_LUT0_PTR
    $AF:2400        $AF:27FF      GRPH_LUT1_PTR
    $AF:2800        $AF:2BFF      GRPH_LUT2_PTR
    $AF:2C00        $AF:2FFF      GRPH_LUT3_PTR
    $AF:3000        $AF:33FF      GRPH_LUT4_PTR           Not Implemented Yet
    $AF:3400        $AF:37FF      GRPH_LUT5_PTR           Not Implemented Yet
    $AF:3800        $AF:3BFF      GRPH_LUT6_PTR           Not Implemented Yet
    $AF:3C00        $AF:3FFF      GRPH_LUT7_PTR           Not Implemented Yet
    $AF:4000        $AF:40FF      GAMMA_B_LUT_PTR
    $AF:4100        $AF:41FF      GAMMA_G_LUT_PTR
    $AF:4200        $AF:42FF      GAMMA_R_LUT_PTR
    $AF:4300        $AF:4FFF      UNDEFINED - RESERVED
    $AF:5000        $AF:57FF      TILE_MAP0               Tile Map 0 Memory Block
    $AF:5800        $AF:5FFF      TILE_MAP1               Tile Map 1 Memory Block
    $AF:6000        $AF:67FF      TILE_MAP2               Tile Map 2 Memory Block
    $AF:6800        $AF:6FFF      TILE_MAP3               Tile Map 3 Memory Block
    $AF:7000        $AF:7FFF      UNDEFINED - RESERVED
rw  $AF:8000        $AF:87FF      FONT_MEMORY_BANK0       FONT Character Graphic Mem
rw  $AF:8800        $AF:8FFF      FONT_MEMORY_BANK1       FONT Character Graphic Mem
    $AF:9000        $AF:9FFF      UNDEFINED - RESERVED
rw  $AF:A000        $AF:BFFF      CS_TEXT_MEM_PTR         Text Memory Block
rw  $AF:C000        $AF:DFFF      CS_COLOR_MEM_PTR        Color Text Memory Block
=== =============   ============  ======================= ===============


Master Control Registers
========================

.. note::
   See also `<https://github.com/Trinity-11/Kernel_FMX/blob/vicky-ii/src/vicky_ii_def.asm>`_

Modes are enabled and disabled via the Vicky Master Control Register at $AF:0000
via the control bits:

aa
--

cc
++

MASTER_CTRL_REG_L = $AF:0000
^^^^^^^^^^^^^^^^^^^^^^^^^^^^
 ::

    Mstr_Ctrl_Text_Mode_En  = $01   Enable the Text Mode
    Mstr_Ctrl_Text_Overlay  = $02   Enable the Overlay of the text mode on top of 
                                    Graphic Mode (the Background Color is ignored)
    Mstr_Ctrl_Graph_Mode_En = $04   Enable the Graphic Mode
    Mstr_Ctrl_Bitmap_En     = $08   Enable the Bitmap Module In Vicky
    Mstr_Ctrl_TileMap_En    = $10   Enable the Tile Module in Vicky
    Mstr_Ctrl_Sprite_En     = $20   Enable the Sprite Module in Vicky
    Mstr_Ctrl_GAMMA_En      = $40   Enable the GAMMA correction - The Analog and DVI have 
                                    different color values, the GAMMA is great to correct 
                                    the difference
    Mstr_Ctrl_Disable_Vid   = $80   This will disable the Scanning of the Video information 
                                    in the 4Meg of VideoRAM hence giving 100% bandwidth to 
                                    the CPU Bitmap Layer


    MASTER_CTRL_REG_H	    = $AF:0001

    Mstr_Ctrl_Video_Mode0   = $01   0 (bit cleared) - 640x480 (Clock @ 25.175Mhz) 
                                    1 (bit set)     - 800x600 (Clock @ 40Mhz)
    Mstr_Ctrl_Video_Mode1   = $02   0 (bit cleared) - No Pixel Doubling, 
                                    1 (bit set)     - Pixel Doubling (reduce the resolution by 2)


For example writing 0x03 to $AF:0001 should set 800x600 mode with 
doubled pixels.

.. note:: FMX Kernel
   Call `SETSIZES` at $00:112C to update the text screen size variables 
   based on the border and screen resolution


Border Control Registers
========================

 ::

  $AF:0004 

    Bit[0]     Enable (1 by default)  
    Bit[4..6]: X Scroll Offset (Will scroll Left) (Acceptable Value: 0..7),
	           i.e. by pixel


Modes
=====

.. warning::
   still vicky I

The bitmap is stored anywhere in $B0 bank memory. If the bitmap is supposed to
start at $B0:0000, the BM_START_ADDY has to be set to $00:0000.

test
  ddsds
  dsdsd

BM_CONTROL_REG = $AF0140
 ::

   Bit 0   = disable/enable
   Bit 1-3 = Target LUT address located at AF:2000 and up.

BM_START_ADDY_L = $AF0141
BM_START_ADDY_M = $AF0142
BM_START_ADDY_H = $AF0143
 ::
  
   Test


LUT
===

.. warning::
   still vicky I

A LUT, namely a Look-Up-Table, stores a selection of colors. 256 colors are
supported in the video composition, which are selectable out of 16.777.216
colors in the 24 Bit RGB color scheme. The LUT also contains an 8 Bit alpha
channel, though it isn't supported. In summary, the LUT has $400 (1024) bytes -
and the order for composing it is B -> G -> R -> A.

As an example, if you would want to compose a LUT of 16 base colors, it would
look like this:

==========  =================== =================
Address     Hex Values (BGRA)   Decimal RGB Value
==========  =================== =================
$AF:2000    -- Transparent --   -- Transparent --
$AF:2004    00 00 00 00         0, 0, 0
$AF:2008    FF FF FF 00         255, 255, 255
$AF:200C    00 00 88 00         0, 0, 136
$AF:2010    EE FF AA 00         238, 255, 170
$AF:2014    CC 44 CC 00         204, 68, 204
$AF:2018    55 CC 00 00         85, 204, 0
$AF:201C    AA 00 00 00         170, 0, 0
$AF:2020    77 EE EE 00         119, 238, 238
$AF:2024    55 88 DD 00         85, 136, 221
$AF:2028    00 44 66 00         0, 68, 102
$AF:202C    77 77 FF 00         119, 119, 255
$AF:2030    33 33 33 00         51, 51, 51
$AF:2034    77 77 77 00         119, 119, 119
$AF:2038    66 FF AA 00         102, 255, 170
$AF:203C    FF 88 00 00         255, 136, 0
$AF:2040    BB BB BB 00         187, 187, 187
==========  =================== =================

Addressing anything in the LUT is achieved by simply dividing the lower 10 bits
of target color address by 4. 


