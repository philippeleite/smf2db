package main

import (
	"database/sql"
	"sync"
)

var pStmt061 *sql.Stmt
var once061 sync.Once
var count061 int

func smf061Init() {
	pStmt061 = createTable("smf061", `ID INT NOT NULL, SMF61TME DATETIME NOT NULL, SMF61DTE DATETIME NOT NULL,
		SMF61CPU TEXT NOT NULL, SMF61SUB TEXT NOT NULL, SMF61VER TEXT, SMF61PNM TEXT, SMF61JNM TEXT,
		SMF61RST DATETIME, SMF61RDT DATETIME, SMF61UID TEXT, SMF61CNM TEXT, SMF61TYP TEXT, SMF61ENM TEXT,
		SMF61CRC BLOB, SMF61ETS DATETIME`)
}

// smf061Parse stores a catalog define (catalog record inserted, updated or deleted). SMF61CRC keeps the
// raw catalog record, whose length is in its first two bytes.
func smf061Parse(rec []byte) {
	once061.Do(smf061Init)

	s := section(rec)
	crcLen := s.num(204, 2)
	var ets interface{} // SMF61ETS, STCKE timestamp after the catalog record; absent in older versions
	if b := s.sub(204+crcLen, 16); b != nil {
		ets = stcketime(b)
	}
	count061 += 1
	insert(pStmt061,
		count061,
		s.time(2),           // SMF61TME
		s.date(6),           // SMF61DTE
		s.text(10, 4),       // SMF61CPU
		s.text(18, 2),       // SMF61SUB: IN, DE or UP
		s.text(36, 2),       // SMF61VER
		s.text(38, 8),       // SMF61PNM
		s.text(46, 8),       // SMF61JNM
		s.time(54),          // SMF61RST
		s.date(58),          // SMF61RDT
		s.text(62, 8),       // SMF61UID
		s.text(71, 44),      // SMF61CNM
		s.text(115, 1),      // SMF61TYP, same values as SMF60TYP
		s.text(116, 44),     // SMF61ENM
		s.blob(204, crcLen), // SMF61CRC
		ets,                 // SMF61ETS
	)
}
