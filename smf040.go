package main

import (
	"database/sql"
	"sync"
)

// Type 40 (dynamic DD), written when a dynamically allocated data set is unallocated, concatenated or
// deconcatenated (TDDFLG). Laid out like type 4's EXCP part: smf040 holds the header, and smf040_dev
// one row per 8-byte device entry ((TDDVAR - 2) / 8 of them, starting at rec[62]).

var pStmt040, pStmt040Dev *sql.Stmt
var once040 sync.Once
var count040 int

var smf040Fields = fields{
	{"TDDRFLG", 0, 'u', 1},
	{"TDDRCDTS", 2, 'T', 4},
	{"TDDRCDTE", 6, 'd', 4},
	{"TDDCPUID", 10, 'c', 4},
	{"TDDUIF", 14, 'c', 8},
	{"TDDONTME", 22, 'T', 4},
	{"TDDONDTE", 26, 'd', 4},
	{"TDDUDATA", 30, 'c', 8},
	{"TDDINVSQ", 38, 'u', 1},
	{"TDDFLG", 39, 'u', 1},
	{"TDDRIN", 40, 'u', 2},
	{"TDDRCIND", 42, 'u', 2},
	{"TDDRCTOT", 44, 'u', 2},
	{"TDDVAR", 60, 'u', 2},
}

var smf040DevFields = fields{
	{"TDDDEVC", 0, 'u', 1},
	{"TDDUTYP", 1, 'u', 1},
	{"TDDCUAD", 2, 'u', 2},
	{"TDDNEXCP", 4, 'u', 4},
}

func smf040Init() {
	pStmt040 = createTable("smf040", "ID INT NOT NULL, "+smf040Fields.columns())
	pStmt040Dev = createTable("smf040_dev", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf040DevFields.columns())
}

func smf040Parse(rec []byte) {
	once040.Do(smf040Init)

	s := section(rec)
	count040 += 1
	insert(pStmt040, append([]interface{}{count040}, smf040Fields.values(s)...)...)
	for k := 0; k < (s.num(60, 2)-2)/8; k++ {
		dev := s.sub(62+8*k, 8)
		if dev == nil {
			break
		}
		insert(pStmt040Dev, append([]interface{}{count040, k + 1}, smf040DevFields.values(dev)...)...)
	}
}
