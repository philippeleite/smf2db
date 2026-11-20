package main

import (
	"database/sql"
	"encoding/binary"
	"sync"
)

// Type 42 (DFSMS statistics and configuration). smf042 holds the header and product section of every
// record; each decoded subtype adds its own smf042_* tables, linked by ID. Subtype parsers live in
// smf042io.go (5, 6) and smf042misc.go (1, 2, 3, 4, 9, 21, 24, 27). Other subtypes only get the header row.

var pStmt042 *sql.Stmt
var once042 sync.Once
var count042 int

var smf042Subtypes = map[int]func(rec []byte, id int){
	1:  smf042s01,
	2:  smf042s02,
	3:  smf042s03,
	4:  smf042s04,
	5:  smf042s05,
	6:  smf042s06,
	9:  smf042s09,
	21: smf042s21,
	24: smf042s24,
	27: smf042s27,
}

func smf042Init() {
	pStmt042 = createTable("smf042", `ID INT NOT NULL, SMF42TME DATETIME NOT NULL, SMF42DTE DATETIME NOT NULL,
		SMF42SID TEXT NOT NULL, SMF42SSI TEXT NOT NULL, SMF42STY INT NOT NULL,
		SMF42PDL TEXT, SMF42PDN TEXT, SMF42PSV INT, SMF42PTS DATETIME, SMF42PTE DATETIME`)
	smf042InitMisc()
	smf042InitIO()
}

func smf042Parse(rec []byte) {
	once042.Do(smf042Init)

	prod, _ := smf042Product(rec)
	count042 += 1
	insert(pStmt042,
		count042,
		smftime(rec[2:6]),                        // SMF42TME
		smfdate(rec[6:10]),                       // SMF42DTE
		ebcdic(rec[10:14]),                       // SMF42SID
		ebcdic(rec[14:18]),                       // SMF42SSI
		int(binary.BigEndian.Uint16(rec[18:20])), // SMF42STY
		prod.text(0, 8),                          // SMF42PDL
		prod.text(8, 10),                         // SMF42PDN
		prod.u8(18),                              // SMF42PSV
		prod.tod(20),                             // SMF42PTS
		prod.tod(28),                             // SMF42PTE
	)

	if f, ok := smf042Subtypes[int(binary.BigEndian.Uint16(rec[18:20]))]; ok {
		f(rec, count042)
	}
}

// smf042Product returns the product section (located by the triplet every subtype has at offset 28).
func smf042Product(rec []byte) (section, int) {
	off, length, n := tripletAt(rec, 24)
	return section(rec).sub(off, length), n
}

// smf042Triplet reads the i-th subtype triplet (they start at offset 36, right after the product
// triplet). It returns the offset of the first section within rec (the triplet counts from the RDW),
// the section length and the number of sections; n is 0 when the section is absent. SMF42NT is not
// used: most records set it to 0 even though the triplets are present.
func smf042Triplet(rec []byte, i int) (off, length, n int) {
	return tripletAt(rec, 32+8*i)
}

// tripletAt reads the triplet at rec[at:at+8] (4-byte offset from the RDW, 2-byte length, 2-byte number)
// and returns the offset within rec. Used by every record type with self-defining sections.
func tripletAt(rec []byte, at int) (off, length, n int) {
	t := section(rec).sub(at, 8)
	off, length, n = t.num(0, 4)-4, t.num(4, 2), t.num(6, 2)
	if off < 0 || n == 0 {
		return 0, 0, 0
	}
	return off, length, n
}

// smf042First returns the first section located by the i-th subtype triplet, or nil when absent.
func smf042First(rec []byte, i int) section {
	off, length, n := smf042Triplet(rec, i)
	if n == 0 {
		return nil
	}
	return section(rec).sub(off, length)
}

// smf042Each calls f for each of the consecutive sections located by the i-th subtype triplet.
func smf042Each(rec []byte, i int, f func(s section)) {
	off, length, n := smf042Triplet(rec, i)
	for k := 0; k < n; k++ {
		f(section(rec).sub(off+k*length, length))
	}
}

// section reads fields of a self-described SMF section. A field past the end of the section is NULL,
// so records from releases that wrote shorter sections still load.
type section []byte

// sub returns length bytes at off, or nil when they are not all inside the section.
func (s section) sub(off, length int) section {
	if off < 0 || length < 0 || off+length > len(s) {
		return nil
	}
	return s[off : off+length]
}

// num returns an n-byte (1 to 8) big-endian unsigned field as an int for navigation, 0 if out of range.
func (s section) num(off, n int) int {
	b := s.sub(off, n)
	if b == nil {
		return 0
	}
	v := 0
	for _, c := range b {
		v = v<<8 | int(c)
	}
	return v
}

func (s section) field(off, n int) interface{} {
	if s.sub(off, n) == nil {
		return nil
	}
	return s.num(off, n)
}

func (s section) u8(off int) interface{}  { return s.field(off, 1) }
func (s section) u16(off int) interface{} { return s.field(off, 2) }
func (s section) u24(off int) interface{} { return s.field(off, 3) }
func (s section) u32(off int) interface{} { return s.field(off, 4) }
func (s section) u64(off int) interface{} { return s.field(off, 8) }

func (s section) s32(off int) interface{} {
	if s.sub(off, 4) == nil {
		return nil
	}
	return int(int32(s.num(off, 4)))
}

// text converts an EBCDIC field in place, so each field must be read only once.
func (s section) text(off, n int) interface{} {
	if b := s.sub(off, n); b != nil {
		return ebcdic(b)
	}
	return nil
}

func (s section) blob(off, n int) interface{} {
	if b := s.sub(off, n); b != nil {
		return []byte(b)
	}
	return nil
}

func (s section) tod(off int) interface{} {
	if b := s.sub(off, 8); b != nil {
		return todtime(b)
	}
	return nil
}

func (s section) time(off int) interface{} {
	if b := s.sub(off, 4); b != nil {
		return smftime(b)
	}
	return nil
}

func (s section) date(off int) interface{} {
	if b := s.sub(off, 4); b != nil {
		return smfdate(b)
	}
	return nil
}

// u32s reads count consecutive 4-byte fields starting at off.
func (s section) u32s(off, count int) []interface{} {
	v := make([]interface{}, count)
	for i := range v {
		v[i] = s.u32(off + 4*i)
	}
	return v
}

// columns returns "PREFIXname TYPE" definitions, comma-separated, for a list of field name suffixes.
func columns(prefix, typ string, names []string) string {
	var b []byte
	for i, n := range names {
		if i > 0 {
			b = append(b, ", "...)
		}
		b = append(b, prefix+n+" "+typ...)
	}
	return string(b)
}
