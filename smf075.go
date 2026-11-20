package main

import (
	"database/sql"
	"sync"
)

// Type 75 (RMF page data set activity): smf075 holds the header and RMF product section, smf075_page one row
// per page data set section, with ID and SEQ. The field table was generated from the manual with
// tools/genfields.py.

var pStmt075 *sql.Stmt
var once075 sync.Once
var count075 int

var smf075Sections = []*rmfSection{
	{subtype: 1, at: 32, table: "smf075_page", f: smf075PageFields}, // page data set
}

func smf075Init() {
	pStmt075 = createTable("smf075", rmfColumns("SMF75"))
	rmfCreateSections(smf075Sections)
}

func smf075Parse(rec []byte) {
	once075.Do(smf075Init)

	count075 += 1
	insert(pStmt075, rmfRow(rec, count075, "SMF75")...)
	rmfInsertSections(smf075Sections, rec, count075)
}

var smf075PageFields = fields{
	{"SMF75DSN", 0, 'c', 44},
	{"SMF75PST", 44, 'u', 1},
	{"SMF75FL2", 45, 'u', 1},
	{"SMF75TYP", 47, 'u', 4},
	{"SMF75CHA", 51, 'u', 2},
	{"SMF75VOL", 53, 'c', 6},
	{"SMF75SCS", 59, 'u', 1},
	{"SMF75SLA", 64, 'u', 4},
	{"SMF75MXU", 68, 'u', 4},
	{"SMF75MNU", 72, 'u', 4},
	{"SMF75AVU", 76, 'u', 4},
	{"SMF75BDS", 80, 'u', 4},
	{"SMF75USE", 84, 'u', 4},
	{"SMF75REQ", 88, 'u', 4}, // manual: 88 / 5A
	{"SMF75SIO", 92, 'u', 4},
	{"SMF75PGX", 96, 'u', 4},
	{"SMF75DEV", 100, 'c', 8},
	{"SMF75CU", 108, 'c', 8},
	{"SMF75LVU", 116, 'f', 8},
}
