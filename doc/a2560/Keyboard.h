#ifndef Keyboard_H_   /* Include guard */
#define Keyboard_H_

// Status
#define 	OUT_BUF_FULL    0x01
#define 	INPT_BUF_FULL	0x02
#define 	SYS_FLAG		0x04
#define 	CMD_DATA		0x08
#define 	KEYBD_INH       0x10
#define 	TRANS_TMOUT	    0x20
#define 	RCV_TMOUT		0x40
#define 	PARITY_EVEN		0x80
#define 	INH_KEYBOARD	0x10
#define 	KBD_ENA			0xAE
#define 	KBD_DIS			0xAD


// Keyboard
unsigned char * STATUS_PORT 	= (void *)0x00C02064;
unsigned char * KBD_OUT_BUF 	= (void *)0x00C02060;
unsigned char * KBD_INPT_BUF 	= (void *)0x00C02060;
unsigned char * KBD_CMD_BUF 	= (void *)0x00C02064;
unsigned char * KBD_DATA_BUF 	= (void *)0x00C02060;
unsigned char * PORT_A 			= (void *)0x00C02060;
unsigned char * PORT_B 			= (void *)0x00C02061;




#endif