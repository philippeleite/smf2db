"""Generate a Go `fields` table (see fields.go) from the field rows of the IBM SMF manual.

Usage:
    pdftotext -raw data/ieag200_v3r1.pdf manual.txt
    python tools/genfields.py manual.txt smf070CpuFields FIRST_LINE LAST_LINE

FIRST_LINE/LAST_LINE delimit one section's rows in manual.txt (1-based, last excluded); find them with
grep -n on the section headings. Rows look like "12 C SMF70SER 3 packed CPU serial number". When the
manual's decimal and hex offsets disagree, the one that continues the sequence is used and the row is
annotated. Reserved fields and redefinitions are skipped. 8-byte binary fields described as STCK or TOD
format become TOD times (unless they are offsets). Review the output: signed fields, times whose description
wraps before the word STCK, and binary data documented as EBCDIC need their kind changed by hand.
"""
import re
import sys

path, var, start, end = sys.argv[1], sys.argv[2], int(sys.argv[3]), int(sys.argv[4])
# split on \n only: pdftotext puts form feeds at page breaks, which str.splitlines would also split on
lines = open(path, encoding='utf-8', errors='replace').read().split('\n')[start - 1:end - 1]
row = re.compile(r'^(\d+) ([0-9A-F]+) (\S+(?: factor)?) (\d+) '
                 r'(binary|Binary|signed|EBCDIC|char|character|packed|floating|float|l_float|s_float|CHAR|unsigned|bit string)\b(.*)$')
kinds = {'binary': 'u', 'unsigned': 'u', 'bit string': 'u', 'signed': 's', 'ebcdic': 'c', 'char': 'c',
         'character': 'c', 'packed': 'p', 'floating': 'f', 'float': 'f', 'l_float': 'f', 's_float': 'f'}
out, expect, seen = [], 0, set()
for line in lines:
    m = row.match(line.strip())
    if not m:
        continue
    dec, hx, name, n, fmt, desc = int(m[1]), int(m[2], 16), re.sub(r'\(\d+\)', '', m[3]).replace(' ', ''), int(m[4]), m[5].lower(), m[6].strip()
    off = dec if dec == hx else (hx if hx == expect else dec)
    note = '' if dec == hx else f' // manual: {dec} / {m[2]}'
    if name == '*' or desc.startswith('Reserved') or name in seen:
        expect = max(expect, off + n)
        continue
    if off < expect and out:
        out.append(f'\t// skipped {name} at {off}: redefines an earlier field')
        continue
    kind = kinds[fmt]
    if kind == 'u' and n == 8 and re.search(r'STCK|TOD format', desc) and 'ffset' not in desc:
        kind = 't'  # absolute TOD clock time (offsets in STCK format stay binary)
    elif kind == 'u' and n > 8:
        kind = 'w'  # wide counters, such as 16-byte sums of squares
    elif kind == 'u' and n not in (1, 2, 3, 4, 8):
        kind = 'b'
    seen.add(name)
    expect = off + n
    out.append(f'\t{{"{name}", {off}, \'{kind}\', {n}}},{note}')
print(f'var {var} = fields{{')
print('\n'.join(out))
print('}')
