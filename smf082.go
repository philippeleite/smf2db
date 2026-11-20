package main

import (
	"database/sql"
	"sync"
)

// Type 82 (ICSF). The layout comes from z/OS Cryptographic Services ICSF System Programmer's Guide
// (SC14-7507, data/csfb200_icsf_spg_hcr77e0.pdf), not from the SMF manual. smf082 holds the header of every
// record; subtype 20 (cryptographic processor processing times), whose data follows the header directly,
// adds smf082_t20. Its NQAP, DQAP and WAIT clocks are stored as microseconds since 1900 (M kind), so
// DQAP - NQAP is the coprocessor time in microseconds. Not decoded, as the sample data has no other
// subtype: every other subtype, including those with the ICSF header (40 and above).

var pStmt082, pStmt082T20 *sql.Stmt
var once082 sync.Once
var count082 int

var smf082HeaderFields = fields{
	{"SMF82tme", 2, 'T', 4},
	{"SMF82dte", 6, 'd', 4},
	{"SMF82sid", 10, 'c', 4},
	{"SMF82ssi", 14, 'c', 4},
	{"SMF82sty", 18, 'u', 2},
}

// smf082T20Fields: subtype 20, offsets from the start of the subtype data (rec[20]).
var smf082T20Fields = fields{
	{"SMF82t20_FLAGS", 0, 'u', 4},
	{"SMF82t20_NQAP", 4, 'M', 8},
	{"SMF82t20_DQAP", 12, 'M', 8},
	{"SMF82t20_WAIT", 20, 'M', 8},
	{"SMF82t20_QLEN", 28, 'u', 4},
	{"SMF82t20_FUNC", 32, 'c', 2},
	{"SMF82t20_APIX", 34, 'u', 1},
	{"SMF82t20_APID", 35, 'c', 8},
	{"SMF82t20_DOMN", 43, 'u', 1},
	{"SMF82t20_REFR", 44, 'u', 1},
	{"SMF82t20_AP_TYPE", 45, 'u', 1},
}

func smf082Init() {
	pStmt082 = createTable("smf082", "ID INT NOT NULL, "+smf082HeaderFields.columns())
	pStmt082T20 = createTable("smf082_t20", "ID INT NOT NULL, "+smf082T20Fields.columns())
}

func smf082Parse(rec []byte) {
	once082.Do(smf082Init)

	s := section(rec)
	count082 += 1
	insert(pStmt082, append([]interface{}{count082}, smf082HeaderFields.values(s)...)...)
	if s.num(18, 2) == 20 {
		insert(pStmt082T20, append([]interface{}{count082}, smf082T20Fields.values(s.sub(20, len(rec)-20))...)...)
	}
}
