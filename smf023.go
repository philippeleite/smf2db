package main

import (
	"database/sql"
	"sync"
)

// Type 23 (SMF status), written at each SMF statistics interval. The header triplets (SMF23POF...,
// offsets from the RDW) start at rec[16]; the product, system and SMF statistics sections occur once,
// so smf023 holds them in one row. The statistics section has cumulative counters (SMF23NGR...) and
// their values for the interval (SMF23NGD...). Not decoded: the logstream statistics section (no sample
// data: SMF records to data sets there) and the spin lock and bind break sections ("for internal use only").

var pStmt023 *sql.Stmt
var once023 sync.Once
var count023 int

var smf023HeaderFields = fields{
	{"SMF23FLG", 0, 'u', 1},
	{"SMF23TME", 2, 'T', 4},
	{"SMF23DTE", 6, 'd', 4},
	{"SMF23SID", 10, 'c', 4},
}

var smf023ProductFields = fields{
	{"SMF23TID", 0, 'u', 2},
	{"SMF23RVN", 2, 'c', 2},
	{"SMF23PNM", 4, 'c', 8},
}

var smf023SystemFields = fields{
	{"SMF23INT", 0, 'c', 6},
	{"SMF23FLS", 6, 'c', 4},
	{"SMF23OSL", 10, 'c', 8},
	{"SMF23TOD", 18, 't', 8},
	{"SMF23SYN", 26, 'c', 8},
	{"SMF23SYP", 34, 'c', 8},
}

var smf023StatFields = fields{
	{"SMF23BFW", 0, 'u', 4},
	{"SMF23BFQ", 4, 'u', 4},
	{"SMF23RCW", 12, 'u', 4},
	{"SMF23BFA", 16, 'u', 4},
	{"SMF23BFT", 20, 'u', 4},
	{"SMF23BFH", 24, 'u', 4},
	{"SMF23BFM", 28, 'u', 4},
	{"SMF23BFL", 32, 'u', 4},
	{"SMF23SFG", 36, 'u', 4},
	{"SMF23NGR", 40, 'u', 8},
	{"SMF23PBG", 48, 'u', 8},
	{"SMF23NFR", 56, 'u', 8},
	{"SMF23PFX", 64, 'u', 8},
	{"SMF231RF", 72, 'u', 8},
	{"SMF23NRF", 80, 'u', 8},
	{"SMF23NIO", 88, 'u', 8},
	{"SMF23TCB", 96, 'u', 8},
	{"SMF23SRB", 104, 'u', 8},
	{"SMF23NGD", 112, 'u', 4},
	{"SMF23PBD", 116, 'u', 4},
	{"SMF23NFD", 120, 'u', 4},
	{"SMF23PFD", 124, 'u', 4},
	{"SMF231RD", 128, 'u', 4},
	{"SMF23NRD", 132, 'u', 4},
	{"SMF23NID", 136, 'u', 8},
	{"SMF23TCD", 144, 'u', 8},
	{"SMF23SRD", 152, 'u', 8},
	{"SMF23MBU", 160, 'u', 4},
}

func smf023Init() {
	pStmt023 = createTable("smf023", "ID INT NOT NULL, "+smf023HeaderFields.columns()+", "+smf023ProductFields.columns()+", "+
		smf023SystemFields.columns()+", "+smf023StatFields.columns())
}

func smf023Parse(rec []byte) {
	once023.Do(smf023Init)

	s := section(rec)
	first := func(at int) section {
		off, length, n := tripletAt(rec, at)
		if n == 0 {
			return nil
		}
		return s.sub(off, length)
	}
	count023 += 1
	row := append([]interface{}{count023}, smf023HeaderFields.values(s)...)
	row = append(row, smf023ProductFields.values(first(16))...)
	row = append(row, smf023SystemFields.values(first(24))...)
	insert(pStmt023, append(row, smf023StatFields.values(first(32))...)...)
}
