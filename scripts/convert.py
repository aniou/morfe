#!/usr/bin/python3

import fileinput

addr_mode = {
    'abs': 'm_Absolute',
    'abs,X': 'm_Absolute_X',
    'abs,Y': 'm_Absolute_Y',
    'acc': 'm_Accumulator',
    'imm': 'm_Immediate',
    'imp': 'm_Implied',
    'dir': 'm_DP',
    'dir,X': 'm_DP_X',
    'dir,Y': 'm_DP_Y',
    '(dir,X)': 'm_DP_X_Indirect',
    '(dir)': 'm_DP_Indirect',
    '[dir]': 'm_DP_Indirect_Long',
    '(dir),Y': 'm_DP_Indirect_Y',
    '[dir],Y': 'm_DP_Indirect_Long_Y',
    '(abs,X)': 'm_Absolute_X_Indirect',
    '(abs)': 'm_Absolute_Indirect',
    '[abs]': 'm_Absolute_Indirect_Long',
    'long': 'm_Absolute_Long',
    'long,X': 'm_Absolute_Long_X',
    'src,dest': 'm_BlockMove',
    'rel8': 'm_PC_Relative',
    'rel16': 'm_PC_Relative_Long',
    'stk,S': 'm_Stack_Relative',
    '(stk,S),Y': 'm_Stack_Relative_Indirect_Y'
}



counter = 0
for line in fileinput.input("cmd1.txt", openhook=fileinput.hook_encoded("utf-8")):
    t = line.split()
    if len(t) < 2:
        print("ERR: bad line %s in %s" % (line, dname))
        continue

    #print(t)

    opcode = int(t[0], 16)
    while counter < opcode:
        print("\t\t{0x%02x, \"XXX\", m_Implied,                   1, 1, c.xxx},\t// illegal/unknown opcode" % counter)
        counter+=1


    opname = t[6].lower()
    desc   = "// %s" % ' '.join(t[6:])
    cycles = t[2][0:1]

    if t[1].endswith('-m'):
        size    = int(t[1][:-2])
        mode    = 'm_Immediate_flagM'
    elif t[1].endswith('-x'):
        size    = int(t[1][:-2])
        mode    = 'm_Immediate_flagX'
    else:
        size    = int(t[1])
        mode    = addr_mode[t[3]]

    # some irregular fixes
    if opname=='pea':
        mode    = 'm_Stack_Implied'

    if opname=='pei':
        mode    = 'm_Stack_DP_Indirect'

    if opname=='pla':
        mode    = 'm_Stack_PC_Relative'

    counter+=1

    hexopcode="0x%02x" % opcode
    # {0x00, "brk", m_Implied, 1, 8, c.brk},
    mode=mode+','
    print("\t\t{0x%02x, \"%s\", %-28s %i, %s, c.%s},\t%s" % (opcode, opname, mode, size, cycles, opname, desc))
    #print(hexopcode, cycles, opname, size, mode, desc)


