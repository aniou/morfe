#ifndef VICKYIII_General_H_   /* Include guard */
#define VICKYIII_General_H_


// VICKY III 
// Channel A - Registers
unsigned int * MasterControlReg_A 		= (void *)0x00C40000;
unsigned int * BorderControlReg_L_A 	= (void *)0x00C40004;
unsigned int * BorderControlReg_H_A 	= (void *)0x00C40008;
unsigned int * BackGroundControlReg_A	= (void *)0x00C4000C;
unsigned int * CursorControlReg_L_A		= (void *)0x00C40010;
unsigned int * CursorControlReg_H_A		= (void *)0x00C40014;
unsigned short * LineInterrupt0_A		= (void *)0x00C40018;
unsigned short * LineInterrupt1_A		= (void *)0x00C4001A;
unsigned short * LineInterrupt2_A		= (void *)0x00C4001C;
unsigned short * LineInterrupt3_A		= (void *)0x00C4001E;

// Mouse Pointer
unsigned short * MousePointer_Mem_A		= (void *)0x00C40400;
unsigned short * MousePtr_A_CTRL_Reg	= (void *)0x00C40C00;

unsigned short * MousePtr_A_X_Pos		= (void *)0x00C40C02;
unsigned short * MousePtr_A_Y_Pos		= (void *)0x00C40C04;
unsigned short * MousePtr_A_Mouse0		= (void *)0x00C40C0A;
unsigned short * MousePtr_A_Mouse1		= (void *)0x00C40C0C;
unsigned short * MousePtr_A_Mouse2		= (void *)0x00C40C0E;

// Channel A - Memory Text Section
unsigned char * ScreenText_A			= (void *)0x00C60000;	// Text Memory
unsigned char * ColorText_A          	= (void *)0x00C68000;		// Color Memory 
unsigned short * FG_CLUT_A 				= (void *)0x00C6C400;		// Foreground LUT
unsigned short * BG_CLUT_A 				= (void *)0x00C6C440;		// Background LUT	

// Channel B - Registers
unsigned int * MasterControlReg_B 		= (void *)0x00C80000;
unsigned int * BorderControlReg_L_B 	= (void *)0x00C80004;
unsigned int * BorderControlReg_H_B 	= (void *)0x00C80008;
unsigned int * BackGroundControlReg_B	= (void *)0x00C8000C;
unsigned int * CursorControlReg_L_B		= (void *)0x00C80010;
unsigned int * CursorControlReg_H_B		= (void *)0x00C80014;
unsigned short * LineInterrupt0_B		= (void *)0x00C80018;
unsigned short * LineInterrupt1_B		= (void *)0x00C8001A;
unsigned short * LineInterrupt2_B		= (void *)0x00C8001C;
unsigned short * LineInterrupt3_B		= (void *)0x00C8001E;
// Mouse Pointer Screen B

unsigned short * MousePointer_Mem_B		= (void *)0x00C80400;
unsigned short * MousePtr_B_CTRL_Reg	= (void *)0x00C80C00;

unsigned short * MousePtr_B_X_Pos		= (void *)0x00C80C02;
unsigned short * MousePtr_B_Y_Pos		= (void *)0x00C80C04;
unsigned short * MousePtr_B_Mouse0		= (void *)0x00C80C0A;
unsigned short * MousePtr_B_Mouse1		= (void *)0x00C80C0C;
unsigned short * MousePtr_B_Mouse2		= (void *)0x00C80C0E;


// Channel A - Memory Text Section
unsigned char  * ScreenText_B			= (void *)0x00CA0000;		// Text Memory
unsigned char  * ColorText_B          	= (void *)0x00CA8000;		// Color Memory 
unsigned short * FG_CLUT_B 				= (void *)0x00CAC400;		// Foreground LUT
unsigned short * BG_CLUT_B 				= (void *)0x00CAC440;		// Background LUT	

unsigned int * BM0_Control_Reg			= (void *)0x00C80100;
unsigned int * BM0_Addy_Pointer_Reg     = (void *)0x00C80104;

unsigned char * LUT_0					= (void *)0x00C82000;
unsigned char * LUT_1					= (void *)0x00C82400;
unsigned char * LUT_2					= (void *)0x00C82800;
unsigned char * LUT_3					= (void *)0x00C82C00;
unsigned char * LUT_4					= (void *)0x00C83000;
unsigned char * LUT_5					= (void *)0x00C83400;
unsigned char * LUT_6					= (void *)0x00C83800;
unsigned char * LUT_7					= (void *)0x00C83C00;

unsigned char * VRAM_BANK0_Char			= (void *)0x00800000;
unsigned short * VRAM_BANK0_Short		= (void *)0x00800000;

unsigned char * VRAM_BANK1_Char			= (void *)0x00A00000;
unsigned short * VRAM_BANK1_Short		= (void *)0x00A00000;

// Sprite 
unsigned short * Sprite_0_CTRL			= (void *)0x00C81000;
unsigned short * Sprite_0_ADDY_HI		= (void *)0x00C81002;
unsigned short * Sprite_0_POS_X	     	= (void *)0x00C81004;
unsigned short * Sprite_0_POS_Y	     	= (void *)0x00C81006;






//  0xHHLL, 0xHHLL 
//  0xGGBB, 0xAARR
unsigned short fg_color_lut [32] = { 
	0x0000, 0xFF00,	// Black (transparent)
	0x0000, 0xFF80, // Mid-Tone Red 
	0x8000, 0xFF00, // Mid-Tone Green
	0x0080, 0xFF00, // Mid-Tone Blue
	0x8000, 0xFF80, // Mid-Tone Yellow
	0x8080, 0xFF00, // Mid-Tone Cian
	0x0080, 0xFF80, // Mid-Tone Purple 
	0x8080, 0xFF80, // 50% Grey 
	0x4500, 0xFFFF, // Orange? Brown?
	0x4513, 0xFF8B, // Orange? Brown?
	0x0000, 0xFF20, // 12.5% Red 
	0x2000, 0xFF00, // 12.5% Green
	0x0020, 0xFF00, // 12.5% Blue
	0x2020, 0xFF20, // 12.5% Grey
	0x4040, 0xFF40, // 25% Grey	
	0xFFFF, 0xFFFF 	// 100% Grey = White
	};

unsigned short bg_color_lut [32] = { 
	0x0000, 0xFF00,	// Black (transparent)
	0x0000, 0xFF80, // Mid-Tone Red 
	0x8000, 0xFF00, // Mid-Tone Green
	0x0080, 0xFF00, // Mid-Tone Blue
	0x2000, 0xFF20, // 12.5% Yellow
	0x2020, 0xFF00, // 12.5% Cian
	0x0020, 0xFF20, // 12.5% Purple 
	0x2020, 0xFF20, // 12.5% Grey 
	0x691E, 0xFFD2, // Orange? Brown?
	0x4513, 0xFF8B, // Orange? Brown?
	0x0000, 0xFF20, // 12.5% Red 
	0x2000, 0xFF00, // 12.5% Green
	0x0020, 0xFF00, // 12.5% Blue
	0x1010, 0xFF10, // 6.25% Grey
	0x4040, 0xFF40, // 25% Grey	
	0xFFFF, 0xFFFF 	// 100% Grey = White
	};

#endif