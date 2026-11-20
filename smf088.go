package main

import (
	"database/sql"
	"sync"
)

// Type 88 (System Logger). Each subtype 1 record describes one log stream, with one log stream, one event
// and one structure (interim storage) section, so smf088 holds all of them in a single row together with
// the header and product section. The 8-byte byte counts documented as EBCDIC (SMF88LWB, SMF88LDB,
// SMF88SWB, SMF88SIB, SMF88SAB) are IBM hex floating point in the records, SMF88LTD is a TOD clock and
// SMF88LIT a binary token. Not decoded, as the sample data has none: subtype 11 (structure alter).
// The field tables were generated from the manual with tools/genfields.py and the kinds above fixed by hand.

var pStmt088 *sql.Stmt
var once088 sync.Once
var count088 int

var smf088HeaderFields = fields{
	{"SMF88TME", 2, 'T', 4},
	{"SMF88DTE", 6, 'd', 4},
	{"SMF88SID", 10, 'c', 4},
	{"SMF88WID", 14, 'c', 4},
	{"SMF88STP", 18, 'u', 2},
}

var smf088ProductFields = fields{
	{"SMF88TYP", 0, 'u', 2},
	{"SMF88RVN", 2, 'c', 2},
	{"SMF88PNM", 4, 'c', 8},
	{"SMF88OSL", 12, 'c', 8},
}

func smf088Init() {
	pStmt088 = createTable("smf088", "ID INT NOT NULL, "+smf088HeaderFields.columns()+", "+smf088ProductFields.columns()+
		", "+smf088LogFields.columns()+", "+smf088EventFields.columns()+", "+smf088StructFields.columns())
}

func smf088Parse(rec []byte) {
	once088.Do(smf088Init)

	s := section(rec)
	count088 += 1
	// first returns the first section located by the triplet at rec[at:], or nil.
	first := func(at int) section {
		off, length, n := tripletAt(rec, at)
		if n == 0 {
			return nil
		}
		return s.sub(off, length)
	}
	row := append([]interface{}{count088}, smf088HeaderFields.values(s)...)
	row = append(row, smf088ProductFields.values(first(24))...)
	row = append(row, smf088LogFields.values(first(32))...)
	row = append(row, smf088EventFields.values(first(40))...)
	insert(pStmt088, append(row, smf088StructFields.values(first(48))...)...)
}

var smf088LogFields = fields{
	{"SMF88LIT", 0, 'x', 8},
	{"SMF88LSN", 8, 'c', 26},
	{"SMF88LFL", 34, 'u', 2},
	{"SMF88LTD", 36, 't', 8},
	{"SMF88LWI", 44, 'u', 4},
	{"SMF88LIB", 48, 'u', 4},
	{"SMF88LAB", 52, 'u', 4},
	{"SMF88LWB", 56, 'f', 8},
	{"SMF88LDB", 64, 'f', 8},
	{"SMF88LIO", 72, 'u', 4},
	{"SMF88LIS", 76, 'u', 4}, // manual: 76 / 52
	{"SMF88GRP", 80, 'c', 8},
}

var smf088EventFields = fields{
	{"SMF88EDS", 0, 'u', 4},
	{"SMF88ERI", 4, 'u', 4},
	{"SMF88ERC", 8, 'u', 4},
	{"SMF88ESF", 12, 'u', 4},
	{"SMF88ETT", 16, 'u', 4},
	{"SMF88ETF", 20, 'u', 4},
	{"SMF88EO", 24, 'u', 4},
	{"SMF88EFS", 28, 'u', 4},
	{"SMF88EDO", 32, 'u', 4},
	{"SMF88EAF", 36, 'u', 4},
	{"SMF88LBL", 40, 'u', 4},
}

var smf088StructFields = fields{
	{"SMF88STN", 0, 'c', 16},
	{"SMF88SWB", 16, 'f', 8},
	{"SMF88SIB", 24, 'f', 8},
	{"SMF88SAB", 32, 'f', 8},
	{"SMF88SII", 44, 'u', 4},
	{"SMF88SAI", 48, 'u', 4},
	{"SMF88SC1", 52, 'u', 4},
	{"SMF88SC2", 56, 'u', 4},
	{"SMF88SC3", 60, 'u', 4},
}
