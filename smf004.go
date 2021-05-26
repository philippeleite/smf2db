package main

import (
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func smf004Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf004;
	CREATE TABLE smf004 (
		SMF4TME    DATETIME   NOT NULL,
		SMF4DTE    DATETIME   NOT NULL,
		SMF4SID    TEXT       NOT NULL,
		SMF4JBN    TEXT       NOT NULL,
		SMF4RST    DATETIME   NOT NULL,
		SMF4RSD    DATETIME   NOT NULL,
		SMF4UIF    TEXT       NOT NULL,
		SMF4STN    INT        NOT NULL,
		SMF4SIT    DATETIME   NOT NULL,
		SMF4STID   DATETIME   NOT NULL,
		SMF4NCI    INT        NOT NULL,
		SMF4SCC    INT        NOT NULL,
		SMF4PRTY   INT        NOT NULL,
		SMF4PGMN   TEXT       NOT NULL,
		SMF4STMN   TEXT       NOT NULL,
		SMF4SYST   INT        NOT NULL,
		SMF4HOST   INT        NOT NULL,
		SMF4RSH0   INT        NOT NULL,
		SMF4SPK    INT        NOT NULL,
		SMF4STI    INT        NOT NULL,
		SMF4AST    INT        NOT NULL,
		SMF4PPST   INT        NOT NULL,
		SMF4SRBT   INT        NOT NULL,
		SMF4RIN    INT        NOT NULL,
		SMF4RLCT   INT        NOT NULL,
		SMF4LENN   INT        NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

func smf004Parse(rec []byte) {
	sqlStmt := `
		INSERT INTO smf004
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
`
	prepStmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	defer prepStmt.Close()

	bcd2asc(rec[10:14], 4) // SMF4SID
	bcd2asc(rec[14:22], 8) // SMF4JBN
	bcd2asc(rec[30:38], 8) // SMF4UIF
	bcd2asc(rec[54:62], 8) // SMF4PGMN
	bcd2asc(rec[62:70], 8) // SMF4STMN
	_, err = prepStmt.Exec(
		smftime(rec[2:6]),                          // SMF4TME
		smfdate(rec[6:10]),                         // SMF4DTE
		string(rec[10:14]),                         // SMF4SID
		string(rec[14:22]),                         // SMF4JBN
		smftime(rec[22:26]),                        // SMF4RST
		smfdate(rec[26:30]),                        // SMF4RSD
		string(rec[30:38]),                         // SMF4UIF
		int(rec[38]),                               // SMF4STN
		smftime(rec[39:43]),                        // SMF4SIT
		smfdate(rec[43:47]),                        // SMF4STID
		int(binary.BigEndian.Uint32(rec[47:51])),   // SMF4NCI
		int(binary.BigEndian.Uint16(rec[51:53])),   // SMF4SCC
		int(rec[53]),                               // SMF4PRTY
		string(rec[54:62]),                         // SMF4PGMN
		string(rec[62:70]),                         // SMF4STMN
		int(binary.BigEndian.Uint16(rec[72:74])),   // SMF4SYST
		int(binary.BigEndian.Uint16(rec[74:76])),   // SMF4HOST
		int(binary.BigEndian.Uint32(rec[78:82])),   // SMF4RSH0
		int(rec[82]),                               // SMF4SPK
		int(rec[83]),                               // SMF4STI
		int(binary.BigEndian.Uint32(rec[86:90])),   // SMF4AST
		int(binary.BigEndian.Uint32(rec[90:94])),   // SMF4PPST
		b3toInt(rec[95:98]),                        // SMF4SRBT
		int(binary.BigEndian.Uint16(rec[98:100])),  // SMF4RIN
		int(binary.BigEndian.Uint16(rec[100:102])), // SMF4RLCT
		int(binary.BigEndian.Uint16(rec[102:104])), // SMF4LENN
	)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}
