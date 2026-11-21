package main

import (
	"database/sql"
	"strings"
	"sync"
)

// Type 35 (TSO/E logoff), laid out like type 5. smf035 holds, in one row, the fixed header, the accounting
// fields (TLGNBRAC entries of a 1-byte length and EBCDIC text) and the relocate section, which has no offset
// field: it starts right after the last accounting field.

var pStmt035 *sql.Stmt
var once035 sync.Once
var count035 int

var smf035Fields = fields{
	{"TLGRCDTS", 2, 'T', 4},
	{"TLGRCDTE", 6, 'd', 4},
	{"TLGPUID", 10, 'c', 4},
	{"TLGUIF", 14, 'c', 8},
	{"TLGONTME", 22, 'T', 4},
	{"TLGONDTE", 26, 'd', 4},
	{"TLGUDATA", 30, 'c', 8},
	{"TLGSTPCT", 38, 'u', 1},
	{"TLGOUTCT", 43, 'u', 4},
	{"TLGINCT", 47, 'u', 4},
	{"TLGSTAT", 51, 'u', 2},
	{"TLGPRI", 53, 'u', 1},
	{"TLGNQTME", 54, 'T', 4},
	{"TLGNQDTE", 58, 'd', 4},
	{"TLGTRMI", 62, 'u', 1},
	{"TLGTRANT", 64, 'u', 4},
	{"TLGSPK", 72, 'u', 1},
	{"TLGSRBT", 73, 'u', 3},
	{"TLGTJS", 76, 'u', 4},
	{"TLGTTAT", 80, 'u', 4},
	{"TLGNTSN", 84, 'u', 4},
	{"TLGPGNO", 88, 'u', 2},
	{"TLGVAR", 92, 'u', 1},
	{"TLGCPUTM", 113, 'u', 3},
	{"TLGNBRAC", 116, 'u', 1},
}

var smf035RelocFields = fields{
	{"TLGCPUS", 0, 'u', 4},
	{"TLGIOCS", 4, 'u', 4},
	{"TLGMSOS", 8, 'u', 4},
	{"TLGSRBS", 12, 'u', 4},
	{"TLGTSN", 16, 'c', 8},
}

func smf035Init() {
	pStmt035 = createTable("smf035", "ID INT NOT NULL, "+smf035Fields.columns()+", TLGACFLD TEXT, "+smf035RelocFields.columns())
}

func smf035Parse(rec []byte) {
	once035.Do(smf035Init)

	s := section(rec)
	count035 += 1
	row := append([]interface{}{count035}, smf035Fields.values(s)...) // decode before the account fields convert in place

	// TLGACFLD: space-separated, as SMF5JSAF.
	var acct strings.Builder
	off := 117
	for k := 0; k < s.num(116, 1); k++ {
		n := s.num(off, 1)
		if t, ok := s.text(off+1, n).(string); ok {
			acct.WriteString(" " + t)
		}
		off += 1 + n
	}
	row = append(row, acct.String())
	insert(pStmt035, append(row, smf035RelocFields.values(s.sub(off, len(rec)-off))...)...)
}
