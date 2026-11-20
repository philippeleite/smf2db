package main

import (
	"database/sql"
	"sync"
)

var pStmt065 *sql.Stmt
var once065 sync.Once
var count065 int

func smf065Init() {
	pStmt065 = createTable("smf065", `ID INT NOT NULL, SMF65TME DATETIME NOT NULL, SMF65DTE DATETIME NOT NULL,
		SMF65CPU TEXT NOT NULL, SMF65SUB TEXT NOT NULL, SMF65VER TEXT, SMF65PNM TEXT, SMF65JNM TEXT,
		SMF65RST DATETIME, SMF65RDT DATETIME, SMF65UID TEXT, SMF65FNC TEXT, SMF65CNM TEXT, SMF65TYP TEXT,
		SMF65ENM TEXT, SMF65CRC BLOB, SMF65ETS DATETIME`)
}

// smf065Parse stores a catalog delete (catalog record updated or deleted). Same layout as type 61, plus
// SMF65FNC. SMF65CRC keeps the raw catalog record, whose length is in its first two bytes.
func smf065Parse(rec []byte) {
	once065.Do(smf065Init)

	s := section(rec)
	crcLen := s.num(204, 2)
	var ets interface{} // SMF65ETS, STCKE timestamp after the catalog record; absent in older versions
	if b := s.sub(204+crcLen, 16); b != nil {
		ets = stcketime(b)
	}
	count065 += 1
	insert(pStmt065,
		count065,
		s.time(2),           // SMF65TME
		s.date(6),           // SMF65DTE
		s.text(10, 4),       // SMF65CPU
		s.text(18, 2),       // SMF65SUB: UP or DE
		s.text(36, 2),       // SMF65VER
		s.text(38, 8),       // SMF65PNM
		s.text(46, 8),       // SMF65JNM
		s.time(54),          // SMF65RST
		s.date(58),          // SMF65RDT
		s.text(62, 8),       // SMF65UID
		s.text(70, 1),       // SMF65FNC: S = data set scratched, U = only catalog entries modified
		s.text(71, 44),      // SMF65CNM
		s.text(115, 1),      // SMF65TYP, same values as SMF60TYP
		s.text(116, 44),     // SMF65ENM
		s.blob(204, crcLen), // SMF65CRC
		ets,                 // SMF65ETS
	)
}
