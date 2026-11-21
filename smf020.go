package main

import (
	"database/sql"
	"strings"
	"sync"
)

// Type 20 (job initiation). smf020 holds, in one row, the fixed header, the accounting fields (SMF20NAF
// entries of a 1-byte length and EBCDIC text, joined with spaces as in SMF5JSAF) and the RACF relocate
// section, located by SMF20RLO, an offset from the flag byte (rec[0]) as in type 4.

var pStmt020 *sql.Stmt
var once020 sync.Once
var count020 int

var smf020Fields = fields{
	{"SMF20FLG", 0, 'u', 1},
	{"SMF20TME", 2, 'T', 4},
	{"SMF20DTE", 6, 'd', 4},
	{"SMF20SID", 10, 'c', 4},
	{"SMF20JBN", 14, 'c', 8},
	{"SMF20RST", 22, 'T', 4},
	{"SMF20RSD", 26, 'd', 4},
	{"SMF20UID", 30, 'c', 8},
	{"SMF20RLO", 38, 'u', 2},
	{"SMF20PGM", 40, 'c', 20},
	{"SMF20NAF", 60, 'u', 1},
}

var smf020RelocFields = fields{
	{"SMF20FLS", 0, 'u', 2},
	{"SMF20GRP", 2, 'c', 8},
	{"SMF20RUD", 10, 'c', 8},
	{"SMF20TID", 18, 'c', 8},
}

func smf020Init() {
	pStmt020 = createTable("smf020", "ID INT NOT NULL, "+smf020Fields.columns()+", SMF20ACT TEXT, "+smf020RelocFields.columns())
}

func smf020Parse(rec []byte) {
	once020.Do(smf020Init)

	s := section(rec)
	count020 += 1
	row := append([]interface{}{count020}, smf020Fields.values(s)...) // decode before the account fields convert in place

	var acct strings.Builder
	off := 61
	for k := 0; k < s.num(60, 1); k++ {
		n := s.num(off, 1)
		if t, ok := s.text(off+1, n).(string); ok {
			acct.WriteString(" " + t)
		}
		off += 1 + n
	}
	row = append(row, acct.String())
	rlo := s.num(38, 2)
	insert(pStmt020, append(row, smf020RelocFields.values(s.sub(rlo, s.num(rlo, 2)))...)...)
}
