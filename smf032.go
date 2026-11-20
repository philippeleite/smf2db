package main

import (
	"database/sql"
	"sync"
)

// Type 32 (TSO/E user work accounting). smf032 holds the header, product and identification sections of
// every record; the TSO/E command sections of subtypes 1 (interval) and 2 (session end), command name and
// count, go to smf032_cmd with ID and SEQ. The header triplets start at manual offset 24 (rec[20]), one
// field earlier than in RMF records. Not decoded, as the sample data has none: subtypes 3 and 4, whose
// command sections add CPU, terminal I/O and EXCP counts per command (TSO/E DETAIL accounting).

var pStmt032, pStmt032Cmd *sql.Stmt
var once032 sync.Once
var count032 int

var smf032HeaderFields = fields{
	{"SMF32TME", 2, 'T', 4},
	{"SMF32DTE", 6, 'd', 4},
	{"SMF32SID", 10, 'c', 4},
	{"SMF32WID", 14, 'c', 4},
	{"SMF32STP", 18, 'u', 2},
}

var smf032ProductFields = fields{
	{"SMF32TYP", 0, 'u', 2},
	{"SMF32RVN", 2, 'c', 2},
	{"SMF32PNM", 4, 'c', 8},
	{"SMF32OSL", 12, 'c', 8},
	{"SMF32SYN", 20, 'c', 8},
	{"SMF32SYP", 28, 'c', 8},
}

var smf032IdentFields = fields{
	{"SMF32JBN", 0, 'c', 8},
	{"SMF32PGM", 8, 'c', 8},
	{"SMF32STM", 16, 'c', 8},
	{"SMF32UIF", 24, 'c', 8},
	{"SMF32JNM", 32, 'c', 8},
	{"SMF32STN", 40, 'u', 2},
	{"SMF32JPT", 46, 'u', 2},
	{"SMF32AST", 48, 'T', 4},
	{"SMF32PPS", 52, 'T', 4},
	{"SMF32SIT", 56, 'T', 4},
	{"SMF32STD", 60, 'd', 4},
	{"SMF32RST", 64, 'T', 4},
	{"SMF32RSD", 68, 'd', 4},
	{"SMF32RET", 72, 'T', 4},
	{"SMF32RED", 76, 'd', 4},
	{"SMF32USR", 80, 'c', 20},
	{"SMF32GRP", 100, 'c', 8},
	{"SMF32RUD", 108, 'c', 8},
	{"SMF32TID", 116, 'c', 8},
}

var smf032CmdFields = fields{ // subtypes 1 and 2
	{"SMF32CMD", 0, 'c', 8},
	{"SMF32CNT", 8, 'u', 4},
}

func smf032Init() {
	pStmt032 = createTable("smf032", "ID INT NOT NULL, "+smf032HeaderFields.columns()+", "+
		smf032ProductFields.columns()+", "+smf032IdentFields.columns())
	pStmt032Cmd = createTable("smf032_cmd", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf032CmdFields.columns())
}

func smf032Parse(rec []byte) {
	once032.Do(smf032Init)

	s := section(rec)
	count032 += 1
	first := func(at int) section {
		off, length, n := tripletAt(rec, at)
		if n == 0 {
			return nil
		}
		return s.sub(off, length)
	}
	row := append([]interface{}{count032}, smf032HeaderFields.values(s)...)
	row = append(row, smf032ProductFields.values(first(20))...)
	insert(pStmt032, append(row, smf032IdentFields.values(first(28))...)...)

	if st := s.num(18, 2); st != 1 && st != 2 {
		return
	}
	off, length, n := tripletAt(rec, 36)
	for k := 0; k < n; k++ {
		insert(pStmt032Cmd, append([]interface{}{count032, k + 1}, smf032CmdFields.values(s.sub(off+k*length, length))...)...)
	}
}
