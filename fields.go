package main

import (
	"fmt"
	"math"
	"strings"
)

// field describes one fixed-offset field of an SMF section, so a section with many fields can be
// declared as a table (copied from the manual) instead of hand-written column and value lists.
type field struct {
	name string
	off  int
	kind byte // see fields.values
	n    int
}

type fields []field

// columns returns the column definitions for createTable, in field order. IBM names that contain
// characters SQL identifiers do not allow (such as R723MTV#) are double-quoted, so queries must quote them too.
func (f fields) columns() string {
	cols := make([]string, len(f))
	for i, fd := range f {
		typ := "INT"
		switch fd.kind {
		case 'c', 'x', 'a':
			typ = "TEXT"
		case 'b':
			typ = "BLOB"
		case 'f', 'i', 'w':
			typ = "REAL"
		case 't', 'T', 'h', 'H', 'd':
			typ = "DATETIME"
		}
		name := fd.name
		if strings.ContainsAny(name, "#@$") {
			name = `"` + name + `"`
		}
		cols[i] = name + " " + typ
	}
	return strings.Join(cols, ", ")
}

// values decodes the fields of s in field order; a field past the end of s is NULL. Kinds:
//
//	u unsigned binary   s signed binary    c EBCDIC text      x hex digits (packed without sign, raw codes)
//	b raw bytes         t TOD clock        f IBM hex float    p packed decimal integer (0..9 digits + sign)
//	h packed time 0hhmmssF                 d packed date 0cyydddF                    i packed interval mmsstttF, in seconds
//	H packed time hhmmsstF (tenths of a second)  T binary time in hundredths of a second since midnight (SMF header times)
//	a IPv4 address (4 bytes) as dotted text
//	M TOD clock as microseconds since 1900, an integer (keeps the precision todtime's text drops; subtract two for durations)
//	w unsigned binary wider than 8 bytes (such as 16-byte sums of squares), as an approximate float
func (f fields) values(s section) []interface{} {
	v := make([]interface{}, len(f))
	for i, fd := range f {
		b := s.sub(fd.off, fd.n)
		if b == nil {
			continue
		}
		switch fd.kind {
		case 'u':
			v[i] = s.num(fd.off, fd.n)
		case 's':
			v[i] = signed(b)
		case 'c':
			v[i] = ebcdic(b)
		case 'x':
			v[i] = fmt.Sprintf("%X", []byte(b))
		case 'b':
			v[i] = []byte(b)
		case 't':
			v[i] = todtime(b)
		case 'f':
			v[i] = hfp(b)
		case 'w':
			w := 0.0
			for _, c := range b {
				w = w*256 + float64(c)
			}
			v[i] = w
		case 'p':
			v[i] = packedInt(b)
		case 'h':
			d := packedDigits(b) // 0hhmmss
			v[i] = fmt.Sprintf("0000-00-00 %s:%s:%s.000", d[1:3], d[3:5], d[5:7])
		case 'T':
			v[i] = smftime(b)
		case 'a':
			v[i] = fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
		case 'M':
			v[i] = s.num(fd.off, 8) >> 12 & (1<<52 - 1) // bits 0-51 of the TOD clock, in microseconds
		case 'H':
			d := packedDigits(b) // hhmmsst
			v[i] = fmt.Sprintf("0000-00-00 %s:%s:%s.%s00", d[0:2], d[2:4], d[4:6], d[6:7])
		case 'd':
			v[i] = smfdate(b)
		case 'i':
			d := packedDigits(b) // mmssttt
			v[i] = float64(digits(d[0:2])*60+digits(d[2:4])) + float64(digits(d[4:7]))/1000
		}
	}
	return v
}

// signed converts a big-endian two's complement field of up to 8 bytes.
func signed(b []byte) int {
	v := int64(0)
	for _, c := range b {
		v = v<<8 | int64(c)
	}
	shift := 64 - 8*uint(len(b))
	return int(v << shift >> shift)
}

// packedDigits returns the decimal digits of a packed field, without the trailing sign nibble.
func packedDigits(b []byte) string {
	d := fmt.Sprintf("%X", b)
	return d[:len(d)-1]
}

// packedInt returns a packed decimal field as an int (non-digit nibbles count as 0).
func packedInt(b []byte) int {
	return digits(packedDigits(b))
}

func digits(d string) int {
	v := 0
	for _, c := range d {
		if c >= '0' && c <= '9' {
			v = v*10 + int(c-'0')
		} else {
			v *= 10
		}
	}
	return v
}

// hfp converts an IBM hexadecimal floating point number (4 or 8 bytes): sign bit, 7-bit excess-64
// base-16 exponent, then the fraction.
func hfp(b []byte) float64 {
	frac := uint64(0)
	for _, c := range b[1:] {
		frac = frac<<8 | uint64(c)
	}
	v := float64(frac) / math.Pow(2, float64(8*(len(b)-1))) * math.Pow(16, float64(int(b[0]&0x7f)-64))
	if b[0]&0x80 != 0 {
		v = -v
	}
	return v
}
