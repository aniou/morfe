#ifndef Ethernet_H_   /* Include guard */
#define Ethernet_H_

unsigned short * ETH_ID_REV 				= (void *)0x00C00650;
unsigned short * ETH_IRQ_CFG 				= (void *)0x00C00654;
unsigned short * ETH_INT_STS 				= (void *)0x00C00658;		
unsigned short * ETH_INT_EN 				= (void *)0x00C0065C;
unsigned short * ETH_RESERVED0 				= (void *)0x00C00660;
unsigned short * ETH_BYTE_TEST				= (void *)0x00C00664;
unsigned short * ETH_FIFO_INT 				= (void *)0x00C00668;
unsigned short * ETH_RX_CFG 				= (void *)0x00C0066C;
unsigned short * ETH_TX_CFG			    	= (void *)0x00C00670;
unsigned short * ETH_HW_CFG					= (void *)0x00C00674;
unsigned short * ETH_RX_DP_CTL				= (void *)0x00C00678;
unsigned short * ETH_RX_FIFO_INF 			= (void *)0x00C0067C;
unsigned short * ETH_TX_FIFO_INF			= (void *)0x00C00680;
unsigned short * ETH_PMT_CTRL				= (void *)0x00C00684;
unsigned short * ETH_GPIO_CFG				= (void *)0x00C00688;
unsigned short * ETH_GPT_CFG 				= (void *)0x00C0068C;
unsigned short * ETH_GPT_CNT				= (void *)0x00C00690;
unsigned short * ETH_RESERVED1				= (void *)0x00C00694;
unsigned short * ETH_WORD_SWAP				= (void *)0x00C00698;
unsigned short * ETH_FREE_RUN				= (void *)0x00C0069C;
unsigned short * ETH_RX_DROP 				= (void *)0x00C006A0;
unsigned short * ETH_MAC_CSR_CMD			= (void *)0x00C006A4;
unsigned short * ETH_MAC_CSR_DATA			= (void *)0x00C006A8;
unsigned short * ETH_AFC_CFG				= (void *)0x00C006AC;
unsigned short * ETH_E2P_CMD				= (void *)0x00C006B0;
unsigned short * ETH_E2P_DATA				= (void *)0x00C006B4;


#endif