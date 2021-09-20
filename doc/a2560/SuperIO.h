#ifndef SuperIO_H_   /* Include guard */
#define SuperIO_H_

// Init Super IO - PME - SMI
unsigned char * PME_STS_REG 	= (void *)0x00C02100;		//C22100
unsigned char * PME_EN_REG 		= (void *)0x00C02102;		//C22102

unsigned char * PME_STS1_REG	= (void *)0x00C02104;
unsigned char * PME_STS2_REG	= (void *)0x00C02105;
unsigned char * PME_STS3_REG	= (void *)0x00C02106;
unsigned char * PME_STS4_REG	= (void *)0x00C02107;
unsigned char * PME_STS5_REG	= (void *)0x00C02108;

unsigned char * PME_EN1_REG		= (void *)0x00C0210A;
unsigned char * PME_EN2_REG		= (void *)0x00C0210B;
unsigned char * PME_EN3_REG		= (void *)0x00C0210C;
unsigned char * PME_EN4_REG		= (void *)0x00C0210D;
unsigned char * PME_EN5_REG		= (void *)0x00C0210E;

unsigned char * SMI_STS1_REG	= (void *)0x00C02110;
unsigned char * SMI_STS2_REG	= (void *)0x00C02111;
unsigned char * SMI_STS3_REG	= (void *)0x00C02112;
unsigned char * SMI_STS4_REG	= (void *)0x00C02113;
unsigned char * SMI_STS5_REG	= (void *)0x00C02114;
				
unsigned char * SMI_EN1_REG		= (void *)0x00C02116;
unsigned char * SMI_EN2_REG		= (void *)0x00C02117;
unsigned char * SMI_EN3_REG		= (void *)0x00C02118;
unsigned char * SMI_EN4_REG		= (void *)0x00C02119;
unsigned char * SMI_EN5_REG		= (void *)0x00C0211A;

unsigned char * MSC_ST_REG			= (void *)0x00C0211C;
unsigned char * FORCE_DISK_CHANGE	= (void *)0x00C0211E;
unsigned char * FLOPPY_DATA_RATE	= (void *)0x00C0211F;

unsigned char * UART1_FIFO_CTRL_SHDW	= (void *)0x00C02120;
unsigned char * UART2_FIFO_CTRL_SHDW	= (void *)0x00C02121;
unsigned char * DEV_DISABLE_REG		= (void *)0x00C02122;

unsigned char * GP10_REG			= (void *)0x00C02123;
unsigned char * GP11_REG			= (void *)0x00C02124;
unsigned char * GP12_REG			= (void *)0x00C02125;
unsigned char * GP13_REG			= (void *)0x00C02126;
unsigned char * GP14_REG			= (void *)0x00C02127;
unsigned char * GP15_REG			= (void *)0x00C02128;
unsigned char * GP16_REG			= (void *)0x00C02129;
unsigned char * GP17_REG			= (void *)0x00C0212A;

unsigned char * GP20_REG			= (void *)0x00C0212B;
unsigned char * GP21_REG			= (void *)0x00C0212C;
unsigned char * GP22_REG			= (void *)0x00C0212D;
unsigned char * GP23_REG			= (void *)0x00C0212E;
unsigned char * GP24_REG			= (void *)0x00C0212F;
unsigned char * GP25_REG			= (void *)0x00C02130;
unsigned char * GP26_REG			= (void *)0x00C02131;
unsigned char * GP27_REG			= (void *)0x00C02132;

unsigned char * GP30_REG			= (void *)0x00C02133;
unsigned char * GP31_REG			= (void *)0x00C02134;
unsigned char * GP32_REG			= (void *)0x00C02135;
unsigned char * GP33_REG			= (void *)0x00C02136;
unsigned char * GP34_REG			= (void *)0x00C02137;
unsigned char * GP35_REG			= (void *)0x00C02138;
unsigned char * GP36_REG			= (void *)0x00C02139;
unsigned char * GP37_REG			= (void *)0x00C0213A;

unsigned char * GP40_REG			= (void *)0x00C0213B;
unsigned char * GP41_REG			= (void *)0x00C0213C;
unsigned char * GP42_REG			= (void *)0x00C0213D;
unsigned char * GP43_REG			= (void *)0x00C0213E;

unsigned char * GP50_REG			= (void *)0x00C0213F;
unsigned char * GP51_REG			= (void *)0x00C02140;
unsigned char * GP52_REG			= (void *)0x00C02141;
unsigned char * GP53_REG			= (void *)0x00C02142;
unsigned char * GP54_REG			= (void *)0x00C02143;
unsigned char * GP55_REG			= (void *)0x00C02144;
unsigned char * GP56_REG			= (void *)0x00C02145;
unsigned char * GP57_REG			= (void *)0x00C02146;

unsigned char * GP60_REG			= (void *)0x00C02147;
unsigned char * GP61_REG			= (void *)0x00C02148;

unsigned char * GP1_REG				= (void *)0x00C0214B;
unsigned char * GP2_REG				= (void *)0x00C0214C;
unsigned char * GP3_REG				= (void *)0x00C0214D;
unsigned char * GP4_REG				= (void *)0x00C0214E;
unsigned char * GP5_REG				= (void *)0x00C0214F;
unsigned char * GP6_REG				= (void *)0x00C02150;

unsigned char * FAN1_REG			= (void *)0x00C02156;
unsigned char * FAN2_REG			= (void *)0x00C02157;
unsigned char * FAN_CTRL_REG		= (void *)0x00C02158;
unsigned char * FAN1_TACH_REG		= (void *)0x00C02159;
unsigned char * FAN2_TACH_REG		= (void *)0x00C0215A;
unsigned char * FAN1_PRELOAD_REG	= (void *)0x00C0215B;
unsigned char * FAN2_PRELOAD_REG	= (void *)0x00C0215C;

unsigned char * LED1_REG			= (void *)0x00C0215D;
unsigned char * LED2_REG			= (void *)0x00C0215E;
unsigned char * KEYBOARD_SCAN_CODE	= (void *)0x00C0215F;

#endif