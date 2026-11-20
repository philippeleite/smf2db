package main

import (
	"database/sql"
	"fmt"
	"sync"
)

// Type 33 (APPC/MVS TP accounting). smf033 holds the header, product and address space ID sections of every
// record; the conversation ID section of subtype 2 (conversation deallocated) goes to smf033_conv, together
// with the partner and local TP names and the LUWID, which the conversation section locates by offsets
// from the start of the record (SMF33CPO, SMF33CTO, SMF33CLO). SMF33RSD is documented as binary but holds a
// packed date. Not decoded, as the sample data has none: subtype 1 (TP identification and usage sections)
// and the user data field section.

var pStmt033, pStmt033Conv *sql.Stmt
var once033 sync.Once
var count033 int

var smf033HeaderFields = fields{
	{"SMF33TME", 2, 'T', 4},
	{"SMF33DTE", 6, 'd', 4},
	{"SMF33SID", 10, 'c', 4},
	{"SMF33WID", 14, 'c', 4},
	{"SMF33STP", 18, 'u', 2},
}

var smf033ProductFields = fields{
	{"SMF33TYP", 0, 'u', 2},
	{"SMF33RVN", 2, 'c', 2},
	{"SMF33PNM", 4, 'c', 8},
	{"SMF33OSL", 12, 'c', 8},
	{"SMF33SYN", 20, 'c', 8},
	{"SMF33SYP", 28, 'c', 8},
}

var smf033AsidFields = fields{
	{"SMF33JID", 0, 'c', 8},
	{"SMF33RST", 8, 'T', 4},
	{"SMF33RSD", 12, 'd', 4},
	{"SMF33STN", 16, 'c', 8},
}

var smf033ConvFields = fields{ // subtype 2 conversation ID section
	{"SMF33CID", 0, 'x', 8},
	{"SMF33CCO", 8, 'x', 8},
	{"SMF33CIO", 20, 'u', 1},
	{"SMF33CLR", 21, 'u', 1},
	{"SMF33CKD", 22, 'u', 1},
	{"SMF33CSL", 23, 'u', 1},
	{"SMF33CLL", 24, 'c', 17},
	{"SMF33CPL", 41, 'c', 17},
	{"SMF33CSH", 60, 'c', 8},
	{"SMF33CPU", 86, 'c', 10},
	{"SMF33CRT", 96, 't', 8},
	{"SMF33CQT", 104, 't', 8},
	{"SMF33CST", 112, 't', 8},
	{"SMF33CET", 120, 't', 8},
	{"SMF33CMN", 128, 'c', 8},
	{"SMF33CSN", 136, 'u', 4},
	{"SMF33CDS", 140, 'f', 8},
	{"SMF33CRE", 148, 'u', 4},
	{"SMF33CDR", 152, 'f', 8},
	{"SMF33CVB", 160, 'u', 4},
	{"SMF33CRC", 164, 'u', 4},
	{"SMF33CRS", 168, 'u', 4},
	{"SMF33CSA", 172, 'u', 4},
	{"SMF33CSS", 176, 't', 8},
	{"SMF33CSE", 184, 't', 8},
}

func smf033Init() {
	pStmt033 = createTable("smf033", "ID INT NOT NULL, "+smf033HeaderFields.columns()+", "+
		smf033ProductFields.columns()+", "+smf033AsidFields.columns())
	pStmt033Conv = createTable("smf033_conv", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf033ConvFields.columns()+
		", SMF33PARTNER_TPN TEXT, SMF33LOCAL_TPN TEXT, SMF33LUW TEXT")
}

func smf033Parse(rec []byte) {
	once033.Do(smf033Init)

	s := section(rec)
	count033 += 1
	first := func(at int) section {
		off, length, n := tripletAt(rec, at)
		if n == 0 {
			return nil
		}
		return s.sub(off, length)
	}
	row := append([]interface{}{count033}, smf033HeaderFields.values(s)...)
	row = append(row, smf033ProductFields.values(first(24))...)
	insert(pStmt033, append(row, smf033AsidFields.values(first(48))...)...)

	if s.num(18, 2) != 2 {
		return
	}
	// lenText reads a 2-byte length followed by EBCDIC text, at an offset counted from the RDW.
	lenText := func(off int) interface{} {
		if off < 4 {
			return nil
		}
		return s.text(off-4+2, s.num(off-4, 2))
	}
	off, length, n := tripletAt(rec, 56) // conversation ID sections
	for k := 0; k < n; k++ {
		c := s.sub(off+k*length, length)
		var luw interface{} // SMF33CLO: LUWID section, 2-byte length then the ID
		if lo := c.num(16, 4); lo >= 4 {
			if b := s.sub(lo-4+2, s.num(lo-4, 2)); b != nil {
				luw = fmt.Sprintf("%X", []byte(b))
			}
		}
		row := append([]interface{}{count033, k + 1}, smf033ConvFields.values(c)...)
		insert(pStmt033Conv, append(row, lenText(c.num(68, 4)), lenText(c.num(72, 4)), luw)...)
	}
}
