package main

import (
	"database/sql"
	"strings"
	"sync"
)

// Type 34 (TS-step termination: TSO/E logoff), laid out like type 4. smf034 holds, in one row, the fixed
// header, the accounting section that follows the EXCP entries (TCB time and account fields) and the
// relocate section at TIVRLCT, an offset from the flag byte (rec[0]), as in type 4. EXCP entries go to
// smf034_dev; none of the sample records has any (TIVVAR = 2), as in type 4.

var pStmt034, pStmt034Dev *sql.Stmt
var once034 sync.Once
var count034 int

var smf034Fields = fields{
	{"TIVRCDTS", 2, 'T', 4},
	{"TIVRCDTE", 6, 'd', 4},
	{"TIVCPUID", 10, 'c', 4},
	{"TIVUIF", 14, 'c', 8},
	{"TIVONTME", 22, 'T', 4},
	{"TIVONDTE", 26, 'd', 4},
	{"TIVUDATA", 30, 'c', 8},
	{"TIVINVSQ", 38, 'u', 1},
	{"TIVSIT", 39, 'T', 4},
	{"TIVOUTCT", 43, 'u', 4},
	{"TIVINCT", 47, 'u', 4},
	{"TIVSTAT", 51, 'u', 2},
	{"TIVPRI", 53, 'u', 1},
	{"TIVPRGNM", 54, 'c', 8},
	{"TIVINVNM", 62, 'c', 8},
	{"TIVSYST", 72, 'u', 2},
	{"TIVMCRE", 74, 'u', 2},
	{"TIVEFRGN", 78, 'u', 4},
	{"TIVSPK", 82, 'u', 1},
	{"TIVSTI", 83, 'u', 1},
	{"TIVAST", 86, 'T', 4},
	{"TIVPPST", 90, 'T', 4},
	{"TIVSRBT", 95, 'u', 3},
	{"TIVRIN", 98, 'u', 2},
}

var smf034RelocFields = fields{
	{"TIVPGIN", 0, 'u', 4},
	{"TIVPGOUT", 4, 'u', 4},
	{"TIVRGNS", 8, 'u', 4},
	{"TIVSIN", 12, 'u', 4},
	{"TIVSOUT", 16, 'u', 4},
	{"TIVVPI", 20, 'u', 4},
	{"TIVVPO", 24, 'u', 4},
	{"TIVSST", 28, 'u', 4},
	{"TIVACT", 32, 'u', 4},
	{"TIVPGNO", 36, 'u', 2},
	{"TIVTRANT", 38, 'u', 4},
	{"TIVCPM", 42, 'u', 4},
	{"TIVRCLAM", 46, 'u', 4},
	{"TIVCPGIN", 50, 'u', 4},
	{"TIVHSPI", 54, 'u', 4},
	{"TIVPGSTL", 58, 'u', 4},
	{"TIVPGSEC", 62, 'u', 8},
	{"TIVLPAI", 70, 'u', 4},
	{"TIVHSPO", 74, 'u', 4},
	{"TIVCPUS", 78, 'u', 4},
	{"TIVIOCS", 82, 'u', 4},
	{"TIVMSOS", 86, 'u', 4},
	{"TIVSRBS", 90, 'u', 4},
	{"TIVTSN", 94, 'c', 8},
}

var smf034DevFields = fields{
	{"TIVDEVC", 0, 'u', 1},
	{"TIVUTYP", 1, 'u', 1},
	{"TIVCUAD", 2, 'u', 2},
	{"TIVNEXCP", 4, 'u', 4},
}

func smf034Init() {
	pStmt034 = createTable("smf034", "ID INT NOT NULL, "+smf034Fields.columns()+
		", TIVCPUTM INT, TIVNBRAC INT, TIVACFLD TEXT, "+smf034RelocFields.columns())
	pStmt034Dev = createTable("smf034_dev", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf034DevFields.columns())
}

func smf034Parse(rec []byte) {
	once034.Do(smf034Init)

	s := section(rec)
	count034 += 1
	ndev := (s.num(102, 2) - 2) / 8 // TIVVAR: length of the EXCP entries plus this 2-byte field
	acct := s.sub(104+8*ndev, len(rec)-104-8*ndev)

	// TIVACFLD: TIVNBRAC fields, each a 1-byte length and EBCDIC text, space-separated as in SMF5JSAF.
	var fieldsText strings.Builder
	off := 5
	for k := 0; k < acct.num(4, 1); k++ {
		n := acct.num(off, 1)
		if t, ok := acct.text(off+1, n).(string); ok {
			fieldsText.WriteString(" " + t)
		}
		off += 1 + n
	}

	row := append([]interface{}{count034}, smf034Fields.values(s)...)
	row = append(row, acct.u24(1), acct.u8(4), fieldsText.String())
	reloc := s.num(100, 2) // TIVRLCT
	insert(pStmt034, append(row, smf034RelocFields.values(s.sub(reloc, len(rec)-reloc))...)...)

	for k := 0; k < ndev; k++ {
		insert(pStmt034Dev, append([]interface{}{count034, k + 1}, smf034DevFields.values(s.sub(104+8*k, 8))...)...)
	}
}
