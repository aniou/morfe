#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <time.h>
#include <sys/time.h>
#include "m68k.h"

/* ROM and RAM sizes - exceeded addr is passed to go65c816 */
#define MAX_RAM 0xadffff

/* Read/write macros */
#define READ_BYTE(BASE, ADDR) (BASE)[ADDR]
#define READ_WORD(BASE, ADDR) (((BASE)[ADDR]<<8) |			\
							  (BASE)[(ADDR)+1])
#define READ_LONG(BASE, ADDR) (((BASE)[ADDR]<<24) |			\
							  ((BASE)[(ADDR)+1]<<16) |		\
							  ((BASE)[(ADDR)+2]<<8) |		\
							  (BASE)[(ADDR)+3])

#define WRITE_BYTE(BASE, ADDR, VAL) (BASE)[ADDR] = (VAL)&0xff
#define WRITE_WORD(BASE, ADDR, VAL) (BASE)[ADDR] = ((VAL)>>8) & 0xff;		\
									(BASE)[(ADDR)+1] = (VAL)&0xff
#define WRITE_LONG(BASE, ADDR, VAL) (BASE)[ADDR] = ((VAL)>>24) & 0xff;		\
									(BASE)[(ADDR)+1] = ((VAL)>>16)&0xff;	\
									(BASE)[(ADDR)+2] = ((VAL)>>8)&0xff;		\
									(BASE)[(ADDR)+3] = (VAL)&0xff


// forward declarations
unsigned int go_m68k_read_memory_8(unsigned int address);
void go_m68k_write_memory_8(unsigned int address, unsigned int value);
unsigned int go_m68k_read_memory_16(unsigned int address);
void go_m68k_write_memory_16(unsigned int address, unsigned int value);
unsigned int go_m68k_read_memory_32(unsigned int address);
void go_m68k_write_memory_32(unsigned int address, unsigned int value);

/* Prototypes */
void exit_error(char* fmt, ...);

/* variables */
unsigned char *g_ram;


/* Exit with an error message.  Use printf syntax. */
void exit_error(char* fmt, ...)
{
	static int guard_val = 0;
	char buff[100];
	unsigned int pc;
	va_list args;

	if(guard_val)
		return;
	else
		guard_val = 1;

	va_start(args, fmt);
	vfprintf(stderr, fmt, args);
	va_end(args);
	fprintf(stderr, "\n");
	pc = m68k_get_reg(NULL, M68K_REG_PPC);
	m68k_disassemble(buff, pc, M68K_CPU_TYPE_68000);
	fprintf(stderr, "At %04x: %s\n", pc, buff);

	exit(EXIT_FAILURE);
}

unsigned int m68k_read_memory_8(unsigned int address)
{
	if(address > MAX_RAM)
        return go_m68k_read_memory_8(address);

	return READ_BYTE(g_ram, address);
}

unsigned int m68k_read_memory_16(unsigned int address)
{
	if(address > MAX_RAM)
		return go_m68k_read_memory_16(address);

	return READ_WORD(g_ram, address);
}

unsigned int m68k_read_memory_32(unsigned int address)
{
	if(address > MAX_RAM)
		return go_m68k_read_memory_32(address);

	return READ_LONG(g_ram, address);
}

void m68k_write_memory_8(unsigned int address, unsigned int value)
{
    //printf("write8  at %x\n", address);
	if(address > MAX_RAM)
        go_m68k_write_memory_8(address, value);

	WRITE_BYTE(g_ram, address, value);
}

void m68k_write_memory_16(unsigned int address, unsigned int value)
{
    //printf("write16  at %x\n", address);
	if(address > MAX_RAM)
        go_m68k_write_memory_16(address, value);

	WRITE_WORD(g_ram, address, value);
}

void m68k_write_memory_32(unsigned int address, unsigned int value)
{
    //printf("write32  at %x\n", address);
	if(address > MAX_RAM)
        go_m68k_write_memory_32(address, value);

	WRITE_LONG(g_ram, address, value);
}

unsigned int m68k_read_disassembler_16(unsigned int address)
{
	return READ_WORD(g_ram, address);
}

unsigned int m68k_read_disassembler_32(unsigned int address)
{
	return READ_LONG(g_ram, address);
}

unsigned char* m68k_init_ram() {
    g_ram = calloc(MAX_RAM, sizeof(unsigned char));
	return(g_ram);
}

/*
 * code used previously for internal tests
 */

/*
//  Disassembler
void make_hex(char* buff, unsigned int pc, unsigned int length)
{
	char* ptr = buff;

	for(;length>0;length -= 2)
	{
		sprintf(ptr, "%04x", m68k_read_disassembler_16(pc));
		pc += 2;
		ptr += 4;
		if(length > 2)
			*ptr++ = ' ';
	}
}

void disassemble_program()
{
	unsigned int pc;
	unsigned int instr_size;
	char buff[100];
	char buff2[100];

	pc = m68k_read_disassembler_32(4);

	while(pc <= 0x16e)
	{
		instr_size = m68k_disassemble(buff, pc, M68K_CPU_TYPE_68000);
		make_hex(buff2, pc, instr_size);
		printf("%03x: %-20s: %s\n", pc, buff2, buff);
		pc += instr_size;
	}
	fflush(stdout);
}

void cpu_instr_callback(int pc)
{
	(void)pc;
    // The following code would print out instructions as they are executed
    
	static char buff[100];
	static char buff2[100];
	static unsigned int lpc;
	static unsigned int instr_size;

	lpc = m68k_get_reg(NULL, M68K_REG_PC);
	instr_size = m68k_disassemble(buff, lpc, M68K_CPU_TYPE_68000);
	make_hex(buff2, lpc, instr_size);
	printf("E %03x: %-20s: %s\n", lpc, buff2, buff);
	fflush(stdout);
    
}

void printCPUSpeed(int cycles) {
        if (cycles > 1000000) {
			printf("%i MHz\n", cycles / 1000000);
		} else if (cycles > 1000) {
			printf("%i kHz\n", cycles / 1000);
		} else {
			printf("%i Hz\n", cycles);
		}
}

void m68k_read_prog() {
	FILE* fhandle;
	if((fhandle = fopen("test", "rb")) == NULL)
		exit_error("Unable to open test");

	if(fread(g_ram, 1, MAX_ROM+1, fhandle) <= 0)
		exit_error("Error reading test");
	disassemble_program();
}

int m68k_step() {
    return m68k_execute(4000);
}

void main_inner()
{
    struct timeval stop, start;
    int cycles;

    cycles = 0;
	while(1)
	{
		// Our loop requires some interleaving to allow us to update the
		// input, output, and nmi devices.

		// Values to execute determine the interleave rate.
		// Smaller values allow for more accurate interleaving with multiple
		// devices/CPUs but is more processor intensive.
		// 100000 is usually a good value to start at, then work from there.

		// Note that I am not emulating the correct clock speed!
        gettimeofday(&start, NULL);
        while(1) {
		    cycles=cycles+m68k_execute(200);
            gettimeofday(&stop, NULL);
            if ((stop.tv_sec - start.tv_sec) >= 1) {
                printf("czas %li cycles %i ", start.tv_sec, cycles);
				printCPUSpeed(cycles);
                gettimeofday(&start, NULL);
                cycles=0;
            }
        }
	}

}

int main() {
	m68k_init_ram();
	m68k_init();
	m68k_set_cpu_type(M68K_CPU_TYPE_68EC030);
    m68k_read_prog();
    m68k_pulse_reset();
    main_inner();
    return(0);
}
*/