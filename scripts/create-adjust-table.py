#!/usr/bin/python3

"""
% grep -E '^[0-9a-fA-F]' commands.txt|cut -b 9-19|sort|uniq -c
     59
     14 -2*m
     14 -2*m+w
      2 -2*x+x*p
      3 -e
      1 +e*p
     55 -m
     53 -m+w
      7 -m+w-x+x*p
     15 -m-x+x*p
      8 +t+t*e*p
      1 +w
     14 -x
     10 -x+w

% grep '2\*m' commands.txt | cut -b 1-3 | xargs -I{} echo "{} 2" | wc -l
28

% grep '2\*m' commands.txt | cut -b 1-3 | xargs -I{} echo "{} 2" > file_m.txt
% tail -3 file_m.txt
6E  2
76  2
7E  2
^^  ^
|    cycles to adjust
|
operand

% create-adjust-table.py file_m.txt

"""

import fileinput
import sys

counter = 0
adjust  = [0] * 256
for line in fileinput.input(sys.argv[1], openhook=fileinput.hook_encoded("utf-8")):
    t = line.split()
    if len(t) < 2:
        print("ERR: bad line %s in %s" % (line, dname))
        continue

    adjust[int(t[0], 16)] = int(t[1])


for idx, val in enumerate(adjust):
    if idx % 16 == 0:
        print()
    print(val, end=", ")

