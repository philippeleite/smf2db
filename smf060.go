package main

import (
	"database/sql"
	"sync"
)

var pStmt060 *sql.Stmt
var once060 sync.Once
var count060 int

func smf060Init() {
	pStmt060 = createTable("smf060", `ID INT NOT NULL, SMF60TME DATETIME NOT NULL, SMF60DTE DATETIME NOT NULL,
		SMF60CPU TEXT NOT NULL, SMF60SUB TEXT NOT NULL, SMF60VER TEXT, SMF60PNM TEXT, SMF60JNM TEXT,
		SMF60RST DATETIME, SMF60RDT DATETIME, SMF60UID TEXT, SMF60CNM TEXT, SMF60TYP TEXT, SMF60ENM TEXT,
		SMF60CRC BLOB, SMF60RQT INT`)
}

// smf060Parse stores a VVDS update (VVR/NVR inserted, updated or deleted). SMF60CRC keeps the raw VVR or
// NVR, whose length is in its first two bytes.
func smf060Parse(rec []byte) {
	once060.Do(smf060Init)

	s := section(rec)
	crcLen := s.num(204, 2)
	count060 += 1
	insert(pStmt060,
		count060,
		s.time(2),           // SMF60TME
		s.date(6),           // SMF60DTE
		s.text(10, 4),       // SMF60CPU
		s.text(18, 2),       // SMF60SUB: UP, DE or IN
		s.text(36, 2),       // SMF60VER
		s.text(38, 8),       // SMF60PNM
		s.text(46, 8),       // SMF60JNM
		s.time(54),          // SMF60RST
		s.date(58),          // SMF60RDT
		s.text(62, 8),       // SMF60UID
		s.text(71, 44),      // SMF60CNM
		s.text(115, 1),      // SMF60TYP
		s.text(116, 44),     // SMF60ENM
		s.blob(204, crcLen), // SMF60CRC
		// Unnamed in the manual: with SMF60VER 2, the 8 bytes after the VVR/NVR hold the time, in
		// microseconds, needed to complete the VVDS request. NULL for older versions.
		s.u64(204+crcLen),
	)
}
