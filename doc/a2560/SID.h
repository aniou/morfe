#ifndef SID_H_   /* Include guard */
#define SID_H_
// $00C20800..$00C008FF - Extern Left SID
// $00C20900..$00C009FF - Extern Right SID

// External SID Left Channel
unsigned char * SID_EXT_L_V1_FREQ_LO  	= (void *)0x00C20800;
unsigned char * SID_EXT_L_V1_FREQ_HI  	= (void *)0x00C20801;
unsigned char * SID_EXT_L_V1_PW_LO		= (void *)0x00C20802;
unsigned char * SID_EXT_L_V1_PW_HI		= (void *)0x00C20803;
unsigned char * SID_EXT_L_V1_CTRL		= (void *)0x00C20804;
unsigned char * SID_EXT_L_V1_ATCK_DECY  = (void *)0x00C20805;
unsigned char * SID_EXT_L_V1_SSTN_RLSE  = (void *)0x00C20806;
unsigned char * SID_EXT_L_V2_FREQ_LO    = (void *)0x00C20807;
unsigned char * SID_EXT_L_V2_FREQ_HI    = (void *)0x00C20808;
unsigned char * SID_EXT_L_V2_PW_LO      = (void *)0x00C20809;
unsigned char * SID_EXT_L_V2_PW_HI		= (void *)0x00C2080A;
unsigned char * SID_EXT_L_V2_CTRL       = (void *)0x00C2080B;
unsigned char * SID_EXT_L_V2_ATCK_DECY  = (void *)0x00C2080C;
unsigned char * SID_EXT_L_V2_SSTN_RLSE  = (void *)0x00C2080D;
unsigned char * SID_EXT_L_V3_FREQ_LO    = (void *)0x00C2080E;
unsigned char * SID_EXT_L_V3_FREQ_HI    = (void *)0x00C2080F;
unsigned char * SID_EXT_L_V3_PW_LO      = (void *)0x00C20810;
unsigned char * SID_EXT_L_V3_PW_HI		= (void *)0x00C20811;
unsigned char * SID_EXT_L_V3_CTRL       = (void *)0x00C20812;
unsigned char * SID_EXT_L_V3_ATCK_DECY  = (void *)0x00C20813;
unsigned char * SID_EXT_L_V3_SSTN_RLSE  = (void *)0x00C20814;
unsigned char * SID_EXT_L_FC_LO			= (void *)0x00C20815;
unsigned char * SID_EXT_L_FC_HI         = (void *)0x00C20816;
unsigned char * SID_EXT_L_RES_FILT      = (void *)0x00C20817;
unsigned char * SID_EXT_L_MODE_VOL      = (void *)0x00C20818;
unsigned char * SID_EXT_L_POT_X         = (void *)0x00C20819;
unsigned char * SID_EXT_L_POT_Y         = (void *)0x00C2081A;
unsigned char * SID_EXT_L_OSC3_RND      = (void *)0x00C2081B;
unsigned char * SID_EXT_L_ENV3          = (void *)0x00C2081C;
unsigned char * SID_EXT_L_NOT_USED0     = (void *)0x00C2081D;
unsigned char * SID_EXT_L_NOT_USED1     = (void *)0x00C2081E;
unsigned char * SID_EXT_L_NOT_USED2     = (void *)0x00C2081F;

// External SID Right Channel
unsigned char * SID_EXT_R_V1_FREQ_LO  	= (void *)0x00C20900;
unsigned char * SID_EXT_R_V1_FREQ_HI  	= (void *)0x00C20901;
unsigned char * SID_EXT_R_V1_PW_LO		= (void *)0x00C20902;
unsigned char * SID_EXT_R_V1_PW_HI		= (void *)0x00C20903;
unsigned char * SID_EXT_R_V1_CTRL		= (void *)0x00C20904;
unsigned char * SID_EXT_R_V1_ATCK_DECY  = (void *)0x00C20905;
unsigned char * SID_EXT_R_V1_SSTN_RLSE  = (void *)0x00C20906;
unsigned char * SID_EXT_R_V2_FREQ_LO    = (void *)0x00C20907;
unsigned char * SID_EXT_R_V2_FREQ_HI    = (void *)0x00C20908;
unsigned char * SID_EXT_R_V2_PW_LO      = (void *)0x00C20909;
unsigned char * SID_EXT_R_V2_PW_HI		= (void *)0x00C2090A;
unsigned char * SID_EXT_R_V2_CTRL       = (void *)0x00C2090B;
unsigned char * SID_EXT_R_V2_ATCK_DECY  = (void *)0x00C2090C;
unsigned char * SID_EXT_R_V2_SSTN_RLSE  = (void *)0x00C2090D;
unsigned char * SID_EXT_R_V3_FREQ_LO    = (void *)0x00C2090E;
unsigned char * SID_EXT_R_V3_FREQ_HI    = (void *)0x00C2090F;
unsigned char * SID_EXT_R_V3_PW_LO      = (void *)0x00C20910;
unsigned char * SID_EXT_R_V3_PW_HI		= (void *)0x00C20911;
unsigned char * SID_EXT_R_V3_CTRL       = (void *)0x00C20912;
unsigned char * SID_EXT_R_V3_ATCK_DECY  = (void *)0x00C20913;
unsigned char * SID_EXT_R_V3_SSTN_RLSE  = (void *)0x00C20914;
unsigned char * SID_EXT_R_FC_LO			= (void *)0x00C20915;
unsigned char * SID_EXT_R_FC_HI         = (void *)0x00C20916;
unsigned char * SID_EXT_R_RES_FILT      = (void *)0x00C20917;
unsigned char * SID_EXT_R_MODE_VOL      = (void *)0x00C20918;
unsigned char * SID_EXT_R_POT_X         = (void *)0x00C20919;
unsigned char * SID_EXT_R_POT_Y         = (void *)0x00C2091A;
unsigned char * SID_EXT_R_OSC3_RND      = (void *)0x00C2091B;
unsigned char * SID_EXT_R_ENV3          = (void *)0x00C2091C;
unsigned char * SID_EXT_R_NOT_USED0     = (void *)0x00C2091D;
unsigned char * SID_EXT_R_NOT_USED1     = (void *)0x00C2091E;
unsigned char * SID_EXT_R_NOT_USED2     = (void *)0x00C2091F;

// $00C21000..$00C211FF - Internal SID Left
// $00C21200..$00C213FF - Internal SID Right
// $00C21400..$00C215FF - Internal SID Neutral
// Internal SID Left Channel
unsigned char * SID_INT_L_V1_FREQ_LO  	= (void *)0x00C21000;
unsigned char * SID_INT_L_V1_FREQ_HI  	= (void *)0x00C21001;
unsigned char * SID_INT_L_V1_PW_LO		= (void *)0x00C21002;
unsigned char * SID_INT_L_V1_PW_HI		= (void *)0x00C21003;
unsigned char * SID_INT_L_V1_CTRL		= (void *)0x00C21004;
unsigned char * SID_INT_L_V1_ATCK_DECY  = (void *)0x00C21005;
unsigned char * SID_INT_L_V1_SSTN_RLSE  = (void *)0x00C21006;
unsigned char * SID_INT_L_V2_FREQ_LO    = (void *)0x00C21007;
unsigned char * SID_INT_L_V2_FREQ_HI    = (void *)0x00C21008;
unsigned char * SID_INT_L_V2_PW_LO      = (void *)0x00C21009;
unsigned char * SID_INT_L_V2_PW_HI		= (void *)0x00C2100A;
unsigned char * SID_INT_L_V2_CTRL       = (void *)0x00C2100B;
unsigned char * SID_INT_L_V2_ATCK_DECY  = (void *)0x00C2100C;
unsigned char * SID_INT_L_V2_SSTN_RLSE  = (void *)0x00C2100D;
unsigned char * SID_INT_L_V3_FREQ_LO    = (void *)0x00C2100E;
unsigned char * SID_INT_L_V3_FREQ_HI    = (void *)0x00C2100F;
unsigned char * SID_INT_L_V3_PW_LO      = (void *)0x00C21010;
unsigned char * SID_INT_L_V3_PW_HI		= (void *)0x00C21011;
unsigned char * SID_INT_L_V3_CTRL       = (void *)0x00C21012;
unsigned char * SID_INT_L_V3_ATCK_DECY  = (void *)0x00C21013;
unsigned char * SID_INT_L_V3_SSTN_RLSE  = (void *)0x00C21014;
unsigned char * SID_INT_L_FC_LO			= (void *)0x00C21015;
unsigned char * SID_INT_L_FC_HI         = (void *)0x00C21016;
unsigned char * SID_INT_L_RES_FILT      = (void *)0x00C21017;
unsigned char * SID_INT_L_MODE_VOL      = (void *)0x00C21018;
unsigned char * SID_INT_L_POT_X         = (void *)0x00C21019;
unsigned char * SID_INT_L_POT_Y         = (void *)0x00C2101A;
unsigned char * SID_INT_L_OSC3_RND      = (void *)0x00C2101B;
unsigned char * SID_INT_L_ENV3          = (void *)0x00C2101C;
unsigned char * SID_INT_L_NOT_USED0     = (void *)0x00C2101D;
unsigned char * SID_INT_L_NOT_USED1     = (void *)0x00C2101E;
unsigned char * SID_INT_L_NOT_USED2     = (void *)0x00C2101F;

// Internal SID Right Channel
unsigned char * SID_INT_R_V1_FREQ_LO  	= (void *)0x00C21200;
unsigned char * SID_INT_R_V1_FREQ_HI  	= (void *)0x00C21201;
unsigned char * SID_INT_R_V1_PW_LO		= (void *)0x00C21202;
unsigned char * SID_INT_R_V1_PW_HI		= (void *)0x00C21203;
unsigned char * SID_INT_R_V1_CTRL		= (void *)0x00C21204;
unsigned char * SID_INT_R_V1_ATCK_DECY  = (void *)0x00C21205;
unsigned char * SID_INT_R_V1_SSTN_RLSE  = (void *)0x00C21206;
unsigned char * SID_INT_R_V2_FREQ_LO    = (void *)0x00C21207;
unsigned char * SID_INT_R_V2_FREQ_HI    = (void *)0x00C21208;
unsigned char * SID_INT_R_V2_PW_LO      = (void *)0x00C21209;
unsigned char * SID_INT_R_V2_PW_HI		= (void *)0x00C2120A;
unsigned char * SID_INT_R_V2_CTRL       = (void *)0x00C2120B;
unsigned char * SID_INT_R_V2_ATCK_DECY  = (void *)0x00C2120C;
unsigned char * SID_INT_R_V2_SSTN_RLSE  = (void *)0x00C2120D;
unsigned char * SID_INT_R_V3_FREQ_LO    = (void *)0x00C2120E;
unsigned char * SID_INT_R_V3_FREQ_HI    = (void *)0x00C2120F;
unsigned char * SID_INT_R_V3_PW_LO      = (void *)0x00C21210;
unsigned char * SID_INT_R_V3_PW_HI		= (void *)0x00C21211;
unsigned char * SID_INT_R_V3_CTRL       = (void *)0x00C21212;
unsigned char * SID_INT_R_V3_ATCK_DECY  = (void *)0x00C21213;
unsigned char * SID_INT_R_V3_SSTN_RLSE  = (void *)0x00C21214;
unsigned char * SID_INT_R_FC_LO			= (void *)0x00C21215;
unsigned char * SID_INT_R_FC_HI         = (void *)0x00C21216;
unsigned char * SID_INT_R_RES_FILT      = (void *)0x00C21217;
unsigned char * SID_INT_R_MODE_VOL      = (void *)0x00C21218;
unsigned char * SID_INT_R_POT_X         = (void *)0x00C21219;
unsigned char * SID_INT_R_POT_Y         = (void *)0x00C2121A;
unsigned char * SID_INT_R_OSC3_RND      = (void *)0x00C2121B;
unsigned char * SID_INT_R_ENV3          = (void *)0x00C2121C;
unsigned char * SID_INT_R_NOT_USED0     = (void *)0x00C2121D;
unsigned char * SID_INT_R_NOT_USED1     = (void *)0x00C2121E;
unsigned char * SID_INT_R_NOT_USED2     = (void *)0x00C2121F;

// Internal SID Neutral Channel - When writting here, the value is written in R and L Channel at the same time
unsigned char * SID_INT_N_V1_FREQ_LO  	= (void *)0x00C41200;
unsigned char * SID_INT_N_V1_FREQ_HI  	= (void *)0x00C41201;
unsigned char * SID_INT_N_V1_PW_LO		= (void *)0x00C41202;
unsigned char * SID_INT_N_V1_PW_HI		= (void *)0x00C41203;
unsigned char * SID_INT_N_V1_CTRL		= (void *)0x00C41204;
unsigned char * SID_INT_N_V1_ATCK_DECY  = (void *)0x00C41205;
unsigned char * SID_INT_N_V1_SSTN_RLSE  = (void *)0x00C41206;
unsigned char * SID_INT_N_V2_FREQ_LO    = (void *)0x00C41207;
unsigned char * SID_INT_N_V2_FREQ_HI    = (void *)0x00C41208;
unsigned char * SID_INT_N_V2_PW_LO      = (void *)0x00C41209;
unsigned char * SID_INT_N_V2_PW_HI		= (void *)0x00C4120A;
unsigned char * SID_INT_N_V2_CTRL       = (void *)0x00C4120B;
unsigned char * SID_INT_N_V2_ATCK_DECY  = (void *)0x00C4120C;
unsigned char * SID_INT_N_V2_SSTN_RLSE  = (void *)0x00C4120D;
unsigned char * SID_INT_N_V3_FREQ_LO    = (void *)0x00C4120E;
unsigned char * SID_INT_N_V3_FREQ_HI    = (void *)0x00C4120F;
unsigned char * SID_INT_N_V3_PW_LO      = (void *)0x00C41210;
unsigned char * SID_INT_N_V3_PW_HI		= (void *)0x00C41211;
unsigned char * SID_INT_N_V3_CTRL       = (void *)0x00C41212;
unsigned char * SID_INT_N_V3_ATCK_DECY  = (void *)0x00C41213;
unsigned char * SID_INT_N_V3_SSTN_RLSE  = (void *)0x00C41214;
unsigned char * SID_INT_N_FC_LO			= (void *)0x00C41215;
unsigned char * SID_INT_N_FC_HI         = (void *)0x00C41216;
unsigned char * SID_INT_N_RES_FILT      = (void *)0x00C41217;
unsigned char * SID_INT_N_MODE_VOL      = (void *)0x00C41218;
unsigned char * SID_INT_N_POT_X         = (void *)0x00C41219;
unsigned char * SID_INT_N_POT_Y         = (void *)0x00C4121A;
unsigned char * SID_INT_N_OSC3_RND      = (void *)0x00C4121B;
unsigned char * SID_INT_N_ENV3          = (void *)0x00C4121C;
unsigned char * SID_INT_N_NOT_USED0     = (void *)0x00C4121D;
unsigned char * SID_INT_N_NOT_USED1     = (void *)0x00C4121E;
unsigned char * SID_INT_N_NOT_USED2     = (void *)0x00C4121F;



#endif